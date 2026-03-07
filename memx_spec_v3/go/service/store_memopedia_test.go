package service

import (
	"context"
	"path/filepath"
	"testing"

	"memx/db"
)

func TestIngestMemopedia(t *testing.T) {
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

	// 正常系
	note, err := svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title:        "API設計方針",
		Body:         "RESTful APIの設計方針について",
		WorkingScope: "knowledge",
	})
	if err != nil {
		t.Fatalf("IngestMemopedia: %v", err)
	}
	if note.ID == "" {
		t.Error("note.ID is empty")
	}
	if note.WorkingScope != "knowledge" {
		t.Errorf("WorkingScope = %q, want %q", note.WorkingScope, "knowledge")
	}

	// 異常系: working_scope 未指定
	_, err = svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title: "テスト",
		Body:  "本文",
	})
	if err == nil {
		t.Error("expected error for missing working_scope")
	}
}

func TestGetMemopedia(t *testing.T) {
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
	created, err := svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title:        "Design Pattern",
		Body:         "Singleton pattern description",
		WorkingScope: "patterns",
	})
	if err != nil {
		t.Fatalf("IngestMemopedia: %v", err)
	}

	// 取得
	got, err := svc.GetMemopedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetMemopedia: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
}

func TestPinUnpinMemopedia(t *testing.T) {
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
	created, err := svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title:        "Pinned Note",
		Body:         "This should be pinned",
		WorkingScope: "test",
	})
	if err != nil {
		t.Fatalf("IngestMemopedia: %v", err)
	}

	// ピン留め
	err = svc.PinMemopedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("PinMemopedia: %v", err)
	}

	// 確認
	got, err := svc.GetMemopedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetMemopedia: %v", err)
	}
	if !got.IsPinned {
		t.Error("expected IsPinned = true")
	}

	// ピン留め解除
	err = svc.UnpinMemopedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("UnpinMemopedia: %v", err)
	}

	// 確認
	got, err = svc.GetMemopedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetMemopedia: %v", err)
	}
	if got.IsPinned {
		t.Error("expected IsPinned = false")
	}
}

func TestListPinnedMemopedia(t *testing.T) {
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

	// ピン留め付きでノート作成
	pinned, err := svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title:        "Pinned Note",
		Body:         "Pinned content",
		WorkingScope: "test",
		IsPinned:     true,
	})
	if err != nil {
		t.Fatalf("IngestMemopedia: %v", err)
	}

	// ピン留めなしでノート作成
	_, err = svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title:        "Unpinned Note",
		Body:         "Unpinned content",
		WorkingScope: "test",
		IsPinned:     false,
	})
	if err != nil {
		t.Fatalf("IngestMemopedia: %v", err)
	}

	// ピン留め一覧取得
	notes, err := svc.ListPinnedMemopedia(ctx, "test", 10)
	if err != nil {
		t.Fatalf("ListPinnedMemopedia: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 pinned note, got %d", len(notes))
	}
	if notes[0].ID != pinned.ID {
		t.Errorf("expected pinned note ID %q, got %q", pinned.ID, notes[0].ID)
	}
}

func TestSearchMemopedia(t *testing.T) {
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

	// 用語定義ノート作成
	_, err = svc.IngestMemopedia(ctx, IngestMemopediaRequest{
		Title:        "マイクロサービス",
		Body:         "マイクロサービスは、アプリケーションを小さなサービスに分割するアーキテクチャパターンです",
		WorkingScope: "glossary",
	})
	if err != nil {
		t.Fatalf("IngestMemopedia: %v", err)
	}

	// 検索
	notes, err := svc.SearchMemopedia(ctx, "マイクロサービス", 10)
	if err != nil {
		t.Fatalf("SearchMemopedia: %v", err)
	}
	if len(notes) < 1 {
		t.Error("expected at least 1 result")
	}
}