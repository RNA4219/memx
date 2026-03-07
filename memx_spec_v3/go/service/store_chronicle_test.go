package service

import (
	"context"
	"path/filepath"
	"testing"

	"memx/db"
)

func TestIngestChronicle(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 正常系: working_scope 必須
	note, err := svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "テストノート",
		Body:         "本文です",
		WorkingScope: "project:memx",
	})
	if err != nil {
		t.Fatalf("IngestChronicle: %v", err)
	}
	if note.ID == "" {
		t.Error("note.ID is empty")
	}
	if note.WorkingScope != "project:memx" {
		t.Errorf("WorkingScope = %q, want %q", note.WorkingScope, "project:memx")
	}

	// 異常系: working_scope 未指定
	_, err = svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title: "テスト",
		Body:  "本文",
	})
	if err == nil {
		t.Error("expected error for missing working_scope")
	}
}

func TestIngestChronicle_SecretDeny(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// secret は deny される
	_, err = svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "Secret Note",
		Body:         "This is secret",
		Sensitivity:  "secret",
		WorkingScope: "test",
	})
	if err == nil {
		t.Error("expected error for secret sensitivity")
	}
}

func TestGetChronicle(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// ノート作成
	created, err := svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "Test Note",
		Body:         "Test Body",
		WorkingScope: "test",
	})
	if err != nil {
		t.Fatalf("IngestChronicle: %v", err)
	}

	// 取得
	got, err := svc.GetChronicle(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetChronicle: %v", err)
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
	_, err = svc.GetChronicle(ctx, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != ErrNotFound {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestSearchChronicle(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 複数ノート作成
	_, err = svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "Go言語入門",
		Body:         "GoはGoogleが開発したプログラミング言語です",
		WorkingScope: "dev",
	})
	if err != nil {
		t.Fatalf("IngestChronicle: %v", err)
	}

	_, err = svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "Python入門",
		Body:         "Pythonは人気のあるスクリプト言語です",
		WorkingScope: "dev",
	})
	if err != nil {
		t.Fatalf("IngestChronicle: %v", err)
	}

	// 検索
	notes, err := svc.SearchChronicle(ctx, "Go", 10)
	if err != nil {
		t.Fatalf("SearchChronicle: %v", err)
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
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
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
SELECT name FROM chronicle.sqlite_master WHERE type='table' AND name='notes_fts';
`).Scan(&tableName)
	if err != nil {
		t.Logf("FTS table query error: %v", err)
	} else {
		t.Logf("FTS table exists: %s", tableName)
	}

	// すべてのテーブルを表示
	rows, err := svc.Conn.DB.QueryContext(ctx, `
SELECT name, type FROM chronicle.sqlite_master WHERE type IN ('table', 'virtual table');
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
		t.Logf("chronicle table: %s (%s)", name, typ)
	}
}

func TestListChronicleByScope(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 異なるスコープでノート作成
	_, err = svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "Note 1",
		Body:         "Body 1",
		WorkingScope: "project:A",
	})
	if err != nil {
		t.Fatalf("IngestChronicle: %v", err)
	}

	_, err = svc.IngestChronicle(ctx, IngestChronicleRequest{
		Title:        "Note 2",
		Body:         "Body 2",
		WorkingScope: "project:B",
	})
	if err != nil {
		t.Fatalf("IngestChronicle: %v", err)
	}

	// スコープでフィルタ
	notes, err := svc.ListChronicleByScope(ctx, "project:A", 10)
	if err != nil {
		t.Fatalf("ListChronicleByScope: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Title != "Note 1" {
		t.Errorf("Title = %q, want %q", notes[0].Title, "Note 1")
	}
}