package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"memx/db"
)

// Service は memx の Usecase 層。
// DB/LLM/Gatekeeper などの実装詳細を隠蔽し、API から呼び出される。
type Service struct {
	Conn       *db.Conn
	Gate       db.Gatekeeper
	MiniLLM    db.MiniLLMClient // 任意（nil の場合は要約生成なし）
	ReflectLLM db.ReflectLLMClient
}

func New(paths db.Paths) (*Service, error) {
	c, err := db.OpenAll(paths)
	if err != nil {
		return nil, err
	}
	return &Service{
		Conn: c,
		Gate: db.NewDefaultGatekeeper(db.GateProfileNormal), // デフォルトは NORMAL
	}, nil
}

func (s *Service) Close() error {
	if s == nil || s.Conn == nil {
		return nil
	}
	return s.Conn.Close()
}

// SetMiniLLM は MiniLLM クライアントを設定する。
func (s *Service) SetMiniLLM(client db.MiniLLMClient) {
	s.MiniLLM = client
}

// SetReflectLLM は ReflectLLM クライアントを設定する。
func (s *Service) SetReflectLLM(client db.ReflectLLMClient) {
	s.ReflectLLM = client
}

// NewResolver は typed_ref 解決用の Resolver を作成する。
// P4 Phase 3B: 現時点の memx-core 実装に合わせた ShortNoteResolver を返す。
func (s *Service) NewResolver() *ShortNoteResolver {
	return NewShortNoteResolver(
		s.searchShortInternal,
		s.showShortInternal,
	)
}

// searchShortInternal は Resolver 内部用の検索関数。
func (s *Service) searchShortInternal(ctx context.Context, query string, topK int) ([]Note, error) {
	return s.SearchShort(ctx, query, topK)
}

// showShortInternal は Resolver 内部用の取得関数。
func (s *Service) showShortInternal(ctx context.Context, id string) (*Note, error) {
	n, err := s.GetShort(ctx, id)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// IngestNoteRequest は short への投入（最小）。
// v1.3 では CLI はこれを API に渡すだけ。
type IngestNoteRequest struct {
	Title       string
	Body        string
	Summary     string
	SourceType  string
	Origin      string
	SourceTrust string
	Sensitivity string
	Tags        []string
	NoLLM       bool // true の場合は LLM による要約・タグ生成をスキップ
}

func (r *IngestNoteRequest) normalize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Body = strings.TrimSpace(r.Body)
	r.SourceType = strings.TrimSpace(r.SourceType)
	r.Origin = strings.TrimSpace(r.Origin)
	r.SourceTrust = strings.TrimSpace(r.SourceTrust)
	r.Sensitivity = strings.TrimSpace(r.Sensitivity)

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

// validate は入力値の妥当性を検証する。
func (r *IngestNoteRequest) validate() error {
	// タイトル長制限（最大 500 文字）
	if len(r.Title) > 500 {
		return fmt.Errorf("%w: title exceeds 500 characters", ErrInvalidArgument)
	}

	// 本文長制限（最大 100000 文字 = 約 50KB）
	if len(r.Body) > 100000 {
		return fmt.Errorf("%w: body exceeds 100000 characters", ErrInvalidArgument)
	}

	// source_type の列挙値チェック
	validSourceTypes := map[string]bool{
		"web":    true,
		"file":   true,
		"chat":   true,
		"agent":  true,
		"manual": true,
	}
	if !validSourceTypes[r.SourceType] {
		return fmt.Errorf("%w: invalid source_type: %s", ErrInvalidArgument, r.SourceType)
	}

	// source_trust の列挙値チェック
	validSourceTrust := map[string]bool{
		"trusted":    true,
		"user_input": true,
		"untrusted":  true,
	}
	if !validSourceTrust[r.SourceTrust] {
		return fmt.Errorf("%w: invalid source_trust: %s", ErrInvalidArgument, r.SourceTrust)
	}

	// sensitivity の列挙値チェック
	validSensitivity := map[string]bool{
		"public":       true,
		"internal":     true,
		"confidential": true,
		"secret":       true,
	}
	if !validSensitivity[r.Sensitivity] {
		return fmt.Errorf("%w: invalid sensitivity: %s", ErrInvalidArgument, r.Sensitivity)
	}

	return nil
}

// IngestShort は short.db にノートを保存する。
// Gatekeeper によるポリシーチェックを行い、deny の場合は保存しない。
// MiniLLM が設定されている場合、Summary が空なら自動で要約を生成する。
func (s *Service) IngestShort(ctx context.Context, req IngestNoteRequest) (Note, error) {
	req.normalize()
	if req.Title == "" || req.Body == "" {
		return Note{}, fmt.Errorf("%w: title/body is required", ErrInvalidArgument)
	}

	// 入力値検証
	if err := req.validate(); err != nil {
		return Note{}, err
	}

	// Gatekeeper チェック
	if s.Gate != nil {
		decision, err := s.Gate.Check(ctx, db.GatekeeperCheckRequest{
			Kind:    db.GateKindMemoryStore,
			Profile: db.GateProfileNormal, // TODO: 設定可能にする
			Content: req.Title + "\n" + req.Body,
			Meta: db.GatekeeperMeta{
				SourceType:  req.SourceType,
				SourceTrust: req.SourceTrust,
				Sensitivity: req.Sensitivity,
				Store:       db.StoreShort,
			},
		})
		if err != nil {
			return Note{}, fmt.Errorf("gatekeeper check: %w", err)
		}
		switch decision.Decision {
		case db.DecisionDeny:
			return Note{}, fmt.Errorf("%w: %s", ErrPolicyDenied, decision.Reason)
		case db.DecisionNeedsHuman:
			// v1.3 では needs_human もエラーとして扱う
			// 将来的には保留キューに入れる等の処理
			return Note{}, fmt.Errorf("%w: %s", ErrNeedsHuman, decision.Reason)
		}
	}

	// 自動要約生成（NoLLM=false, Summary空, MiniLLM設定済みの場合）
	summary := req.Summary
	if !req.NoLLM && summary == "" && s.MiniLLM != nil {
		result, err := s.MiniLLM.Summarize(ctx, req.Title, req.Body)
		if err != nil {
			// 要約生成失敗は警告レベルとし、空の要約で保存を継続
			// TODO: ログ出力
		} else {
			summary = result.Summary
		}
	}

	id, err := newUUIDLike()
	if err != nil {
		return Note{}, err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
INSERT INTO notes(
  id, title, summary, body,
  created_at, updated_at, last_accessed_at,
  access_count,
  source_type, origin, source_trust, sensitivity
) VALUES(?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?)
`, id, req.Title, summary, req.Body, now, now, now, req.SourceType, req.Origin, req.SourceTrust, req.Sensitivity)
	if err != nil {
		return Note{}, err
	}

	for _, t := range req.Tags {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if err := upsertTagAndBind(ctx, tx, id, t, now); err != nil {
			return Note{}, err
		}
	}

	// meta（近似でOK）
	_, _ = tx.ExecContext(ctx, `UPDATE short_meta SET value = CAST(CAST(value AS INTEGER) + 1 AS TEXT) WHERE key = 'note_count';`)

	if err := tx.Commit(); err != nil {
		return Note{}, err
	}

	return Note{
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
	}, nil
}

func upsertTagAndBind(ctx context.Context, tx *sql.Tx, noteID, tag, now string) error {
	// route は最小実装として short_only。
	_, err := tx.ExecContext(ctx, `
INSERT INTO tags(name, route, parent_id, created_at, updated_at, usage_count)
VALUES(?, 'short_only', NULL, ?, ?, 0)
ON CONFLICT(name) DO UPDATE SET updated_at = excluded.updated_at;
`, tag, now, now)
	if err != nil {
		return err
	}

	var tagID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM tags WHERE name = ?;`, tag).Scan(&tagID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO note_tags(note_id, tag_id) VALUES(?, ?);`, noteID, tagID); err != nil {
		return err
	}
	_, _ = tx.ExecContext(ctx, `UPDATE tags SET usage_count = usage_count + 1 WHERE id = ?;`, tagID)
	return nil
}

// SearchShort は FTS5 を使った検索。
func (s *Service) SearchShort(ctx context.Context, query string, topK int) ([]Note, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("%w: query is required", ErrInvalidArgument)
	}
	// クエリ長制限（最大 1000 文字）
	if len(query) > 1000 {
		return nil, fmt.Errorf("%w: query exceeds 1000 characters", ErrInvalidArgument)
	}
	if topK <= 0 {
		topK = 20
	}
	// topK 上限
	if topK > 100 {
		topK = 100
	}

	like := likePattern(query)
	rows, err := queryRowsWithFallback(ctx, s.Conn.DB,
		`
SELECT n.id, n.title, n.summary, n.body,
       n.created_at, n.updated_at, n.last_accessed_at, n.access_count,
       n.source_type, n.origin, n.source_trust, n.sensitivity
FROM notes_fts
JOIN notes n ON notes_fts.rowid = n.rowid
WHERE notes_fts MATCH ?
ORDER BY bm25(notes_fts)
LIMIT ?;
`, []interface{}{query, topK},
		`
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity
FROM notes
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

	out := make([]Note, 0, topK)
	for rows.Next() {
		var n Note
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Summary, &n.Body,
			&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
			&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
		); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetShort は id 指定でノートを取得し、アクセス情報を更新する。
func (s *Service) GetShort(ctx context.Context, id string) (Note, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Note{}, fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}
	// ID 形式チェック（32文字のhex）
	if len(id) != 32 {
		return Note{}, fmt.Errorf("%w: invalid id format", ErrInvalidArgument)
	}
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return Note{}, fmt.Errorf("%w: invalid id format: must be hex", ErrInvalidArgument)
		}
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var n Note
	err = tx.QueryRowContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity
FROM notes WHERE id = ?;
`, id).Scan(
		&n.ID, &n.Title, &n.Summary, &n.Body,
		&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
		&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}

	_, _ = tx.ExecContext(ctx, `
UPDATE notes
SET last_accessed_at = ?, access_count = access_count + 1
WHERE id = ?;
`, now, id)

	if err := tx.Commit(); err != nil {
		return Note{}, err
	}

	n.LastAccessedAt = now
	n.AccessCount++
	return n, nil
}

func newUUIDLike() (string, error) {
	// 32 hex = 128-bit。UUIDv4 に厳密に合わせる必要は v1 ではない。
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SummarizeNoteRequest は既存ノートの要約再生成リクエスト。
type SummarizeNoteRequest struct {
	ID string
}

// SummarizeNote は既存ノートの要約を再生成して更新する。
func (s *Service) SummarizeNote(ctx context.Context, id string) (Note, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Note{}, fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}

	// ノート取得
	note, err := s.GetShort(ctx, id)
	if err != nil {
		return Note{}, err
	}

	// LLM で要約生成
	if s.MiniLLM == nil {
		return Note{}, fmt.Errorf("%w: MiniLLM is not configured", ErrInvalidArgument)
	}

	result, err := s.MiniLLM.Summarize(ctx, note.Title, note.Body)
	if err != nil {
		return Note{}, fmt.Errorf("summarize: %w", err)
	}

	// 更新
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err = s.Conn.DB.ExecContext(ctx, `
UPDATE notes SET summary = ?, updated_at = ? WHERE id = ?;
`, result.Summary, now, id)
	if err != nil {
		return Note{}, err
	}

	note.Summary = result.Summary
	note.UpdatedAt = now
	return note, nil
}

// SummarizeNotesRequest は複数ノートの統合要約リクエスト。
type SummarizeNotesRequest struct {
	IDs []string
}

// SummarizeNotesResult は複数ノート統合要約の結果。
type SummarizeNotesResult struct {
	Summary   string
	NoteCount int
}

// SummarizeNotes は複数ノートを統合して1つの要約を生成する。
func (s *Service) SummarizeNotes(ctx context.Context, ids []string) (SummarizeNotesResult, error) {
	if len(ids) == 0 {
		return SummarizeNotesResult{}, fmt.Errorf("%w: at least one id is required", ErrInvalidArgument)
	}

	if s.ReflectLLM == nil {
		return SummarizeNotesResult{}, fmt.Errorf("%w: ReflectLLM is not configured", ErrInvalidArgument)
	}

	// ノート本文を収集
	var bodies []string
	for _, id := range ids {
		note, err := s.GetShort(ctx, strings.TrimSpace(id))
		if err != nil {
			return SummarizeNotesResult{}, fmt.Errorf("get note %s: %w", id, err)
		}
		bodies = append(bodies, note.Body)
	}

	// 統合要約
	input := db.ClusterInput{
		NoteIDs: ids,
		Body:    strings.Join(bodies, "\n\n---\n\n"),
	}

	summary, err := s.ReflectLLM.SummarizeCluster(ctx, input)
	if err != nil {
		return SummarizeNotesResult{}, fmt.Errorf("summarize cluster: %w", err)
	}

	return SummarizeNotesResult{
		Summary:   summary,
		NoteCount: len(ids),
	}, nil
}
