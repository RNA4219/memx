package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"memx/db"
)

// MemopediaNote は memopedia ストアのノート。
// working_scope が必須。知識ベース（用語定義・設計・方針）。
type MemopediaNote struct {
	Note
	WorkingScope string
	IsPinned     bool
}

// IngestMemopediaRequest は memopedia への投入リクエスト。
type IngestMemopediaRequest struct {
	Title        string
	Body         string
	Summary      string
	SourceType   string
	Origin       string
	SourceTrust  string
	Sensitivity  string
	Tags         []string
	WorkingScope string // 必須
	IsPinned     bool
	NoLLM        bool
}

func (r *IngestMemopediaRequest) normalize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Body = strings.TrimSpace(r.Body)
	r.SourceType = strings.TrimSpace(r.SourceType)
	r.Origin = strings.TrimSpace(r.Origin)
	r.SourceTrust = strings.TrimSpace(r.SourceTrust)
	r.Sensitivity = strings.TrimSpace(r.Sensitivity)
	r.WorkingScope = strings.TrimSpace(r.WorkingScope)

	if r.SourceType == "" {
		r.SourceType = "manual"
	}
	if r.SourceTrust == "" {
		r.SourceTrust = "user_input"
	}
	if r.Sensitivity == "" {
		r.Sensitivity = "internal"
	}
}

func (r *IngestMemopediaRequest) validate() error {
	if len(r.Title) > 500 {
		return fmt.Errorf("%w: title exceeds 500 characters", ErrInvalidArgument)
	}
	if len(r.Body) > 100000 {
		return fmt.Errorf("%w: body exceeds 100000 characters", ErrInvalidArgument)
	}

	validSourceTypes := map[string]bool{
		"web": true, "file": true, "chat": true, "agent": true, "manual": true,
	}
	if !validSourceTypes[r.SourceType] {
		return fmt.Errorf("%w: invalid source_type: %s", ErrInvalidArgument, r.SourceType)
	}

	validSourceTrust := map[string]bool{
		"trusted": true, "user_input": true, "untrusted": true,
	}
	if !validSourceTrust[r.SourceTrust] {
		return fmt.Errorf("%w: invalid source_trust: %s", ErrInvalidArgument, r.SourceTrust)
	}

	validSensitivity := map[string]bool{
		"public": true, "internal": true, "secret": true,
	}
	if !validSensitivity[r.Sensitivity] {
		return fmt.Errorf("%w: invalid sensitivity: %s", ErrInvalidArgument, r.Sensitivity)
	}

	// working_scope は必須
	if r.WorkingScope == "" {
		return fmt.Errorf("%w: working_scope is required for memopedia", ErrInvalidArgument)
	}

	return nil
}

// IngestMemopedia は memopedia.db にノートを保存する。
func (s *Service) IngestMemopedia(ctx context.Context, req IngestMemopediaRequest) (MemopediaNote, error) {
	req.normalize()
	if req.Title == "" || req.Body == "" {
		return MemopediaNote{}, fmt.Errorf("%w: title/body is required", ErrInvalidArgument)
	}

	if err := req.validate(); err != nil {
		return MemopediaNote{}, err
	}

	// Gatekeeper チェック
	if s.Gate != nil {
		decision, err := s.Gate.Check(ctx, db.GatekeeperCheckRequest{
			Kind:    db.GateKindMemoryStore,
			Profile: db.GateProfileNormal,
			Content: req.Title + "\n" + req.Body,
			Meta: db.GatekeeperMeta{
				SourceType:  req.SourceType,
				SourceTrust: req.SourceTrust,
				Sensitivity: req.Sensitivity,
				Store:       db.StoreMemopedia,
			},
		})
		if err != nil {
			return MemopediaNote{}, fmt.Errorf("gatekeeper check: %w", err)
		}
		if decision.Decision == db.DecisionDeny {
			return MemopediaNote{}, fmt.Errorf("%w: %s", ErrPolicyDenied, decision.Reason)
		}
		if decision.Decision == db.DecisionNeedsHuman {
			return MemopediaNote{}, fmt.Errorf("%w: %s", ErrNeedsHuman, decision.Reason)
		}
	}

	summary := req.Summary
	if !req.NoLLM && summary == "" && s.MiniLLM != nil {
		result, err := s.MiniLLM.Summarize(ctx, req.Title, req.Body)
		if err != nil {
			// 警告レベル、継続
		} else {
			summary = result.Summary
		}
	}

	id, err := newUUIDLike()
	if err != nil {
		return MemopediaNote{}, err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)

	isPinned := 0
	if req.IsPinned {
		isPinned = 1
	}

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return MemopediaNote{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
INSERT INTO memopedia.notes(
  id, title, summary, body,
  created_at, updated_at, last_accessed_at,
  access_count,
  source_type, origin, source_trust, sensitivity,
  working_scope, is_pinned
) VALUES(?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?)
`, id, req.Title, summary, req.Body, now, now, now, req.SourceType, req.Origin, req.SourceTrust, req.Sensitivity, req.WorkingScope, isPinned)
	if err != nil {
		return MemopediaNote{}, err
	}

	for _, t := range req.Tags {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if err := upsertTagAndBindMemopedia(ctx, tx, id, t, now); err != nil {
			return MemopediaNote{}, err
		}
	}

	_, _ = tx.ExecContext(ctx, `UPDATE memopedia.memopedia_meta SET value = CAST(CAST(value AS INTEGER) + 1 AS TEXT) WHERE key = 'note_count';`)

	if err := tx.Commit(); err != nil {
		return MemopediaNote{}, err
	}

	return MemopediaNote{
		Note: Note{
			ID:             id,
			Title:          req.Title,
			Summary:        summary,
			Body:           req.Body,
			CreatedAt:      now,
			UpdatedAt:      now,
			LastAccessedAt: now,
			AccessCount:    0,
			SourceType:     req.SourceType,
			Origin:         req.Origin,
			SourceTrust:    req.SourceTrust,
			Sensitivity:    req.Sensitivity,
		},
		WorkingScope: req.WorkingScope,
		IsPinned:     req.IsPinned,
	}, nil
}

func upsertTagAndBindMemopedia(ctx context.Context, tx *sql.Tx, noteID, tag, now string) error {
	_, err := tx.ExecContext(ctx, `
INSERT INTO memopedia.tags(name, route, parent_id, created_at, updated_at, usage_count)
VALUES(?, 'memopedia_only', NULL, ?, ?, 0)
ON CONFLICT(name) DO UPDATE SET updated_at = excluded.updated_at;
`, tag, now, now)
	if err != nil {
		return err
	}

	var tagID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM memopedia.tags WHERE name = ?;`, tag).Scan(&tagID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO memopedia.note_tags(note_id, tag_id) VALUES(?, ?);`, noteID, tagID); err != nil {
		return err
	}
	_, _ = tx.ExecContext(ctx, `UPDATE memopedia.tags SET usage_count = usage_count + 1 WHERE id = ?;`, tagID)
	return nil
}

// SearchMemopedia は memopedia ストアをFTS5で検索する。
func (s *Service) SearchMemopedia(ctx context.Context, query string, topK int) ([]MemopediaNote, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("%w: query is required", ErrInvalidArgument)
	}
	if len(query) > 1000 {
		return nil, fmt.Errorf("%w: query exceeds 1000 characters", ErrInvalidArgument)
	}
	if topK <= 0 {
		topK = 20
	}
	if topK > 100 {
		topK = 100
	}

	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT n.id, n.title, n.summary, n.body,
       n.created_at, n.updated_at, n.last_accessed_at, n.access_count,
       n.source_type, n.origin, n.source_trust, n.sensitivity,
       n.working_scope, n.is_pinned
FROM memopedia.notes n
WHERE n.rowid IN (
  SELECT rowid FROM memopedia.notes_fts WHERE notes_fts MATCH ?
)
ORDER BY n.created_at DESC
LIMIT ?;
`, query, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]MemopediaNote, 0, topK)
	for rows.Next() {
		var n MemopediaNote
		var isPinned int
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Summary, &n.Body,
			&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
			&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
			&n.WorkingScope, &isPinned,
		); err != nil {
			return nil, err
		}
		n.IsPinned = isPinned == 1
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMemopedia は id 指定で memopedia ノートを取得する。
func (s *Service) GetMemopedia(ctx context.Context, id string) (MemopediaNote, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return MemopediaNote{}, fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}
	if len(id) != 32 {
		return MemopediaNote{}, fmt.Errorf("%w: invalid id format", ErrInvalidArgument)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return MemopediaNote{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var n MemopediaNote
	var isPinned int
	err = tx.QueryRowContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity,
       working_scope, is_pinned
FROM memopedia.notes WHERE id = ?;
`, id).Scan(
		&n.ID, &n.Title, &n.Summary, &n.Body,
		&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
		&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
		&n.WorkingScope, &isPinned,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return MemopediaNote{}, ErrNotFound
		}
		return MemopediaNote{}, err
	}

	_, _ = tx.ExecContext(ctx, `
UPDATE memopedia.notes
SET last_accessed_at = ?, access_count = access_count + 1
WHERE id = ?;
`, now, id)

	if err := tx.Commit(); err != nil {
		return MemopediaNote{}, err
	}

	n.IsPinned = isPinned == 1
	n.LastAccessedAt = now
	n.AccessCount++
	return n, nil
}

// ListMemopediaByScope は working_scope でフィルタしてノートを取得する。
func (s *Service) ListMemopediaByScope(ctx context.Context, workingScope string, limit int) ([]MemopediaNote, error) {
	workingScope = strings.TrimSpace(workingScope)
	if workingScope == "" {
		return nil, fmt.Errorf("%w: working_scope is required", ErrInvalidArgument)
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity,
       working_scope, is_pinned
FROM memopedia.notes
WHERE working_scope = ?
ORDER BY created_at DESC
LIMIT ?;
`, workingScope, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]MemopediaNote, 0, limit)
	for rows.Next() {
		var n MemopediaNote
		var isPinned int
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Summary, &n.Body,
			&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
			&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
			&n.WorkingScope, &isPinned,
		); err != nil {
			return nil, err
		}
		n.IsPinned = isPinned == 1
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// ListPinnedMemopedia は is_pinned=1 のノートを取得する（Working Memory）。
func (s *Service) ListPinnedMemopedia(ctx context.Context, workingScope string, limit int) ([]MemopediaNote, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var rows *sql.Rows
	var err error

	if workingScope != "" {
		rows, err = s.Conn.DB.QueryContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity,
       working_scope, is_pinned
FROM memopedia.notes
WHERE is_pinned = 1 AND working_scope = ?
ORDER BY updated_at DESC
LIMIT ?;
`, workingScope, limit)
	} else {
		rows, err = s.Conn.DB.QueryContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity,
       working_scope, is_pinned
FROM memopedia.notes
WHERE is_pinned = 1
ORDER BY updated_at DESC
LIMIT ?;
`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]MemopediaNote, 0, limit)
	for rows.Next() {
		var n MemopediaNote
		var isPinned int
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Summary, &n.Body,
			&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
			&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
			&n.WorkingScope, &isPinned,
		); err != nil {
			return nil, err
		}
		n.IsPinned = isPinned == 1
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// PinMemopedia はノートをピン留めする（Working Memory化）。
func (s *Service) PinMemopedia(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}

	result, err := s.Conn.DB.ExecContext(ctx, `
UPDATE memopedia.notes SET is_pinned = 1, updated_at = ? WHERE id = ?;
`, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// UnpinMemopedia はノートのピン留めを解除する。
func (s *Service) UnpinMemopedia(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}

	result, err := s.Conn.DB.ExecContext(ctx, `
UPDATE memopedia.notes SET is_pinned = 0, updated_at = ? WHERE id = ?;
`, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}