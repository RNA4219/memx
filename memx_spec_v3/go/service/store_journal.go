package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"memx/db"
)

// JournalNote は journal ストアのノート。
// working_scope が必須。
type JournalNote struct {
	Note
	WorkingScope string
	IsPinned     bool
}

// IngestJournalRequest は journal への投入リクエスト。
type IngestJournalRequest struct {
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

func (r *IngestJournalRequest) normalize() {
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

func (r *IngestJournalRequest) validate() error {
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
		return fmt.Errorf("%w: working_scope is required for journal", ErrInvalidArgument)
	}

	return nil
}

// IngestJournal は journal.db にノートを保存する。
func (s *Service) IngestJournal(ctx context.Context, req IngestJournalRequest) (JournalNote, error) {
	req.normalize()
	if req.Title == "" || req.Body == "" {
		return JournalNote{}, fmt.Errorf("%w: title/body is required", ErrInvalidArgument)
	}

	if err := req.validate(); err != nil {
		return JournalNote{}, err
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
				Store:       db.StoreJournal,
			},
		})
		if err != nil {
			return JournalNote{}, fmt.Errorf("gatekeeper check: %w", err)
		}
		if decision.Decision == db.DecisionDeny {
			return JournalNote{}, fmt.Errorf("%w: %s", ErrPolicyDenied, decision.Reason)
		}
		if decision.Decision == db.DecisionNeedsHuman {
			return JournalNote{}, fmt.Errorf("%w: %s", ErrNeedsHuman, decision.Reason)
		}
	}

	summary := req.Summary
	if !req.NoLLM && summary == "" && s.MiniLLM != nil {
		result, err := s.MiniLLM.Summarize(ctx, req.Title, req.Body)
		if err != nil {
			s.warnAutoSummaryFailure("journal", req.Title, err)
		} else {
			summary = result.Summary
		}
	}

	id, err := newUUIDLike()
	if err != nil {
		return JournalNote{}, err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)

	isPinned := 0
	if req.IsPinned {
		isPinned = 1
	}

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return JournalNote{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
INSERT INTO journal.notes(
  id, title, summary, body,
  created_at, updated_at, last_accessed_at,
  access_count,
  source_type, origin, source_trust, sensitivity,
  working_scope, is_pinned
) VALUES(?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?)
`, id, req.Title, summary, req.Body, now, now, now, req.SourceType, req.Origin, req.SourceTrust, req.Sensitivity, req.WorkingScope, isPinned)
	if err != nil {
		return JournalNote{}, err
	}

	for _, t := range req.Tags {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if err := upsertTagAndBindJournal(ctx, tx, id, t, now); err != nil {
			return JournalNote{}, err
		}
	}

	_, _ = tx.ExecContext(ctx, `UPDATE journal.journal_meta SET value = CAST(CAST(value AS INTEGER) + 1 AS TEXT) WHERE key = 'note_count';`)

	if err := tx.Commit(); err != nil {
		return JournalNote{}, err
	}

	return JournalNote{
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

func upsertTagAndBindJournal(ctx context.Context, tx *sql.Tx, noteID, tag, now string) error {
	_, err := tx.ExecContext(ctx, `
INSERT INTO journal.tags(name, route, parent_id, created_at, updated_at, usage_count)
VALUES(?, 'journal_only', NULL, ?, ?, 0)
ON CONFLICT(name) DO UPDATE SET updated_at = excluded.updated_at;
`, tag, now, now)
	if err != nil {
		return err
	}

	var tagID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM journal.tags WHERE name = ?;`, tag).Scan(&tagID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO journal.note_tags(note_id, tag_id) VALUES(?, ?);`, noteID, tagID); err != nil {
		return err
	}
	_, _ = tx.ExecContext(ctx, `UPDATE journal.tags SET usage_count = usage_count + 1 WHERE id = ?;`, tagID)
	return nil
}

// SearchJournal は journal ストアをFTS5で検索する。
func (s *Service) SearchJournal(ctx context.Context, query string, topK int) ([]JournalNote, error) {
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

	like := likePattern(query)
	rows, err := queryRowsWithFallback(ctx, s.Conn.DB,
		`
SELECT n.id, n.title, n.summary, n.body,
       n.created_at, n.updated_at, n.last_accessed_at, n.access_count,
       n.source_type, n.origin, n.source_trust, n.sensitivity,
       n.working_scope, n.is_pinned
FROM journal.notes n
WHERE n.rowid IN (
  SELECT rowid FROM journal.notes_fts WHERE notes_fts MATCH ?
)
ORDER BY n.created_at DESC
LIMIT ?;
`, []interface{}{query, topK},
		`
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity,
       working_scope, is_pinned
FROM journal.notes
WHERE title LIKE ?
   OR summary LIKE ?
   OR body LIKE ?
ORDER BY created_at DESC
LIMIT ?;
`, []interface{}{like, like, like, topK})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]JournalNote, 0, topK)
	for rows.Next() {
		var n JournalNote
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

// GetJournal は id 指定で journal ノートを取得する。
func (s *Service) GetJournal(ctx context.Context, id string) (JournalNote, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return JournalNote{}, fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}
	if len(id) != 32 {
		return JournalNote{}, fmt.Errorf("%w: invalid id format", ErrInvalidArgument)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	var n JournalNote
	var isPinned int
	err := s.Conn.DB.QueryRowContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity,
       working_scope, is_pinned
FROM journal.notes WHERE id = ?;
`, id).Scan(
		&n.ID, &n.Title, &n.Summary, &n.Body,
		&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
		&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
		&n.WorkingScope, &isPinned,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return JournalNote{}, ErrNotFound
		}
		return JournalNote{}, err
	}

	_, _ = s.Conn.DB.ExecContext(ctx, `
UPDATE journal.notes
SET last_accessed_at = ?, access_count = access_count + 1
WHERE id = ?;
`, now, id)

	n.IsPinned = isPinned == 1
	n.LastAccessedAt = now
	n.AccessCount++
	return n, nil
}

// ListJournalByScope は working_scope でフィルタしてノートを取得する。
func (s *Service) ListJournalByScope(ctx context.Context, workingScope string, limit int) ([]JournalNote, error) {
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
FROM journal.notes
WHERE working_scope = ?
ORDER BY created_at DESC
LIMIT ?;
`, workingScope, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]JournalNote, 0, limit)
	for rows.Next() {
		var n JournalNote
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
