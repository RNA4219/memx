package service

import (
	"context"
	"path/filepath"
	"testing"

	"memx/db"
)

func TestIngestJournal(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 正常系: working_scope 必須
	note, err := svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "テストノート",
		Body:         "本文です",
		WorkingScope: "project:memx",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}
	if note.ID == "" {
		t.Error("note.ID is empty")
	}
	if note.WorkingScope != "project:memx" {
		t.Errorf("WorkingScope = %q, want %q", note.WorkingScope, "project:memx")
	}

	// 異常系: working_scope 未指定
	_, err = svc.IngestJournal(ctx, IngestJournalRequest{
		Title: "テスト",
		Body:  "本文",
	})
	if err == nil {
		t.Error("expected error for missing working_scope")
	}
}

func TestIngestJournal_AutoSummaryFailureLogsWarning(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	logBuf, logger := newBufferedTestLogger()
	svc.SetLogger(logger)
	svc.SetMiniLLM(&mockMiniLLM{err: errTestLLM})

	note, err := svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Journal Title",
		Body:         "本文です",
		WorkingScope: "project:memx",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}
	if note.Summary != "" {
		t.Errorf("expected empty summary on LLM failure, got %q", note.Summary)
	}

	assertAutoSummaryWarningLogged(t, logBuf.String(), "journal", "Journal Title", errTestLLM)
}

func TestIngestJournal_SecretDeny(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// secret は deny される
	_, err = svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Secret Note",
		Body:         "This is secret",
		Sensitivity:  "secret",
		WorkingScope: "test",
	})
	if err == nil {
		t.Error("expected error for secret sensitivity")
	}
}

func TestGetJournal(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// ノート作成
	created, err := svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Test Note",
		Body:         "Test Body",
		WorkingScope: "test",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}

	// 取得
	got, err := svc.GetJournal(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetJournal: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
	if got.WorkingScope != "test" {
		t.Errorf("WorkingScope = %q, want %q", got.WorkingScope, "test")
	}
	if got.AccessCount != 1 {
		t.Errorf("AccessCount = %d, want 1", got.AccessCount)
	}

	// 存在しないID（有効なhex形式）
	_, err = svc.GetJournal(ctx, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != ErrNotFound {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestSearchJournal(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 複数ノート作成
	_, err = svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Go Programming Language",
		Body:         "Go is a programming language developed by Google",
		WorkingScope: "dev",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}

	_, err = svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Python Programming",
		Body:         "Python is a popular scripting language",
		WorkingScope: "dev",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}

	// 検索
	notes, err := svc.SearchJournal(ctx, "Go", 10)
	if err != nil {
		t.Fatalf("SearchJournal: %v", err)
	}
	if len(notes) < 1 {
		t.Error("expected at least 1 result")
	}
}

// TestFTSExistence はFTSテーブルが正しく作成されているか確認する
func TestFTSExistence(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// FTSテーブルの存在確認
	var tableName string
	err = svc.Conn.DB.QueryRowContext(ctx, `
SELECT name FROM journal.sqlite_master WHERE type='table' AND name='notes_fts';
`).Scan(&tableName)
	if err != nil {
		t.Logf("FTS table query error: %v", err)
	} else {
		t.Logf("FTS table exists: %s", tableName)
	}

	// すべてのテーブルを表示
	rows, err := svc.Conn.DB.QueryContext(ctx, `
SELECT name, type FROM journal.sqlite_master WHERE type IN ('table', 'virtual table');
`)
	if err != nil {
		t.Fatalf("query sqlite_master: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name, typ string
		if err := rows.Scan(&name, &typ); err != nil {
			t.Fatalf("scan: %v", err)
		}
		t.Logf("journal table: %s (%s)", name, typ)
	}
}

func TestListJournalByScope(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 異なるスコープでノート作成
	_, err = svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Note 1",
		Body:         "Body 1",
		WorkingScope: "project:A",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}

	_, err = svc.IngestJournal(ctx, IngestJournalRequest{
		Title:        "Note 2",
		Body:         "Body 2",
		WorkingScope: "project:B",
	})
	if err != nil {
		t.Fatalf("IngestJournal: %v", err)
	}

	// スコープでフィルタ
	notes, err := svc.ListJournalByScope(ctx, "project:A", 10)
	if err != nil {
		t.Fatalf("ListJournalByScope: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Title != "Note 1" {
		t.Errorf("Title = %q, want %q", notes[0].Title, "Note 1")
	}
}
