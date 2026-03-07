package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ArchiveNote は archive ストアのノート。
// アーカイブ済みのノード。FTS/embeddingなし。
type ArchiveNote struct {
	Note
}

// GetArchive は id 指定で archive ノートを取得する。
// archive は通常検索対象外だが、バックトラック用に取得可能。
func (s *Service) GetArchive(ctx context.Context, id string) (ArchiveNote, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return ArchiveNote{}, fmt.Errorf("%w: id is required", ErrInvalidArgument)
	}
	if len(id) != 32 {
		return ArchiveNote{}, fmt.Errorf("%w: invalid id format", ErrInvalidArgument)
	}

	var n ArchiveNote
	err := s.Conn.DB.QueryRowContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity
FROM archive.notes WHERE id = ?;
`, id).Scan(
		&n.ID, &n.Title, &n.Summary, &n.Body,
		&n.CreatedAt, &n.UpdatedAt, &n.LastAccessedAt, &n.AccessCount,
		&n.SourceType, &n.Origin, &n.SourceTrust, &n.Sensitivity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ArchiveNote{}, ErrNotFound
		}
		return ArchiveNote{}, err
	}

	return n, nil
}

// ListArchive は archive ノートを一覧する（管理用）。
func (s *Service) ListArchive(ctx context.Context, limit int) ([]ArchiveNote, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT id, title, summary, body,
       created_at, updated_at, last_accessed_at, access_count,
       source_type, origin, source_trust, sensitivity
FROM archive.notes
ORDER BY created_at DESC
LIMIT ?;
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ArchiveNote, 0, limit)
	for rows.Next() {
		var n ArchiveNote
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

// ArchiveNoteRequest は short→archive 退避用のリクエスト。
// GC から呼び出される。
type ArchiveNoteRequest struct {
	Note
	SrcStore   string // 元のストア名（通常 "short"）
	SrcNoteID  string // 元のノートID
	Relation   string // lineage 用の関係（通常 "archived_from"）
}

// ArchiveNoteFromShort は short ノートを archive に退避する。
// GC（Phase 3）から呼び出される。
// 補償設計：先に archive へ保存 → lineage 記録 → 最後に short 削除。
func (s *Service) ArchiveNoteFromShort(ctx context.Context, srcID string) (ArchiveNote, error) {
	srcID = strings.TrimSpace(srcID)
	if srcID == "" {
		return ArchiveNote{}, fmt.Errorf("%w: src_id is required", ErrInvalidArgument)
	}

	// 1. short からノート取得
	srcNote, err := s.GetShort(ctx, srcID)
	if err != nil {
		return ArchiveNote{}, fmt.Errorf("get source note: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return ArchiveNote{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// 2. archive へ挿入
	_, err = tx.ExecContext(ctx, `
INSERT INTO archive.notes(
  id, title, summary, body,
  created_at, updated_at, last_accessed_at,
  access_count,
  source_type, origin, source_trust, sensitivity
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, srcNote.ID, srcNote.Title, srcNote.Summary, srcNote.Body,
		srcNote.CreatedAt, now, srcNote.LastAccessedAt,
		srcNote.AccessCount,
		srcNote.SourceType, srcNote.Origin, srcNote.SourceTrust, srcNote.Sensitivity)
	if err != nil {
		return ArchiveNote{}, fmt.Errorf("insert to archive: %w", err)
	}

	// 3. lineage に記録
	_, err = tx.ExecContext(ctx, `
INSERT INTO archive.lineage(src_store, src_note_id, dest_store, dest_note_id, relation, created_at)
VALUES(?, ?, 'archive', ?, 'archived_from', ?);
`, "short", srcID, srcNote.ID, now)
	if err != nil {
		return ArchiveNote{}, fmt.Errorf("insert lineage: %w", err)
	}

	// 4. short から削除
	_, err = tx.ExecContext(ctx, `DELETE FROM notes WHERE id = ?;`, srcID)
	if err != nil {
		return ArchiveNote{}, fmt.Errorf("delete from short: %w", err)
	}

	// 5. short.meta 更新
	_, _ = tx.ExecContext(ctx, `UPDATE short_meta SET value = CAST(CAST(value AS INTEGER) - 1 AS TEXT) WHERE key = 'note_count';`)

	// 6. archive.meta 更新
	_, _ = tx.ExecContext(ctx, `UPDATE archive.archive_meta SET value = CAST(CAST(value AS INTEGER) + 1 AS TEXT) WHERE key = 'note_count';`)

	if err := tx.Commit(); err != nil {
		return ArchiveNote{}, err
	}

	return ArchiveNote{Note: srcNote}, nil
}

// RestoreFromArchive は archive から short へ復元する。
func (s *Service) RestoreFromArchive(ctx context.Context, archiveID string) (Note, error) {
	archiveID = strings.TrimSpace(archiveID)
	if archiveID == "" {
		return Note{}, fmt.Errorf("%w: archive_id is required", ErrInvalidArgument)
	}

	// 1. archive からノート取得
	archiveNote, err := s.GetArchive(ctx, archiveID)
	if err != nil {
		return Note{}, fmt.Errorf("get archive note: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// 2. short へ挿入
	_, err = tx.ExecContext(ctx, `
INSERT INTO notes(
  id, title, summary, body,
  created_at, updated_at, last_accessed_at,
  access_count,
  source_type, origin, source_trust, sensitivity
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, archiveNote.ID, archiveNote.Title, archiveNote.Summary, archiveNote.Body,
		archiveNote.CreatedAt, now, now,
		0, // access_count リセット
		archiveNote.SourceType, archiveNote.Origin, archiveNote.SourceTrust, archiveNote.Sensitivity)
	if err != nil {
		return Note{}, fmt.Errorf("insert to short: %w", err)
	}

	// 3. lineage に記録
	_, err = tx.ExecContext(ctx, `
INSERT INTO lineage(src_store, src_note_id, dest_store, dest_note_id, relation, created_at)
VALUES('archive', ?, 'short', ?, 'restored_to', ?);
`, archiveID, archiveNote.ID, now)
	if err != nil {
		return Note{}, fmt.Errorf("insert lineage: %w", err)
	}

	// 4. archive から削除
	_, err = tx.ExecContext(ctx, `DELETE FROM archive.notes WHERE id = ?;`, archiveID)
	if err != nil {
		return Note{}, fmt.Errorf("delete from archive: %w", err)
	}

	// 5. meta 更新
	_, _ = tx.ExecContext(ctx, `UPDATE short_meta SET value = CAST(CAST(value AS INTEGER) + 1 AS TEXT) WHERE key = 'note_count';`)
	_, _ = tx.ExecContext(ctx, `UPDATE archive.archive_meta SET value = CAST(CAST(value AS INTEGER) - 1 AS TEXT) WHERE key = 'note_count';`)

	if err := tx.Commit(); err != nil {
		return Note{}, err
	}

	return archiveNote.Note, nil
}

// GetArchiveLineage は archive ノートの系譜を取得する。
func (s *Service) GetArchiveLineage(ctx context.Context, archiveID string) ([]LineageRecord, error) {
	archiveID = strings.TrimSpace(archiveID)
	if archiveID == "" {
		return nil, fmt.Errorf("%w: archive_id is required", ErrInvalidArgument)
	}

	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT id, src_store, src_note_id, dest_store, dest_note_id, relation, created_at
FROM archive.lineage
WHERE dest_note_id = ? OR src_note_id = ?
ORDER BY created_at DESC;
`, archiveID, archiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]LineageRecord, 0)
	for rows.Next() {
		var r LineageRecord
		if err := rows.Scan(
			&r.ID, &r.SrcStore, &r.SrcNoteID, &r.DestStore, &r.DestNoteID, &r.Relation, &r.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// LineageRecord は lineage テーブルのレコード。
type LineageRecord struct {
	ID          int64
	SrcStore    string
	SrcNoteID   string
	DestStore   string
	DestNoteID  string
	Relation    string
	CreatedAt   string
}