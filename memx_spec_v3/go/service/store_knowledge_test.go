package service

import (
	"context"
	"path/filepath"
	"testing"

	"memx/db"
)

func TestIngestKnowledge(t *testing.T) {
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

	// 正常系
	note, err := svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "API設計方針",
		Body:         "RESTful APIの設計方針について",
		WorkingScope: "knowledge",
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}
	if note.ID == "" {
		t.Error("note.ID is empty")
	}
	if note.WorkingScope != "knowledge" {
		t.Errorf("WorkingScope = %q, want %q", note.WorkingScope, "knowledge")
	}

	// 異常系: working_scope 未指定
	_, err = svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title: "テスト",
		Body:  "本文",
	})
	if err == nil {
		t.Error("expected error for missing working_scope")
	}
}

func TestIngestKnowledge_AutoSummaryFailureLogsWarning(t *testing.T) {
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

	note, err := svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "Knowledge Title",
		Body:         "本文です",
		WorkingScope: "knowledge",
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}
	if note.Summary != "" {
		t.Errorf("expected empty summary on LLM failure, got %q", note.Summary)
	}

	assertAutoSummaryWarningLogged(t, logBuf.String(), "knowledge", "Knowledge Title", errTestLLM)
}

func TestGetKnowledge(t *testing.T) {
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
	created, err := svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "Design Pattern",
		Body:         "Singleton pattern description",
		WorkingScope: "patterns",
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}

	// 取得
	got, err := svc.GetKnowledge(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetKnowledge: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
}

func TestPinUnpinKnowledge(t *testing.T) {
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
	created, err := svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "Pinned Note",
		Body:         "This should be pinned",
		WorkingScope: "test",
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}

	// ピン留め
	err = svc.PinKnowledge(ctx, created.ID)
	if err != nil {
		t.Fatalf("PinKnowledge: %v", err)
	}

	// 確認
	got, err := svc.GetKnowledge(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetKnowledge: %v", err)
	}
	if !got.IsPinned {
		t.Error("expected IsPinned = true")
	}

	// ピン留め解除
	err = svc.UnpinKnowledge(ctx, created.ID)
	if err != nil {
		t.Fatalf("UnpinKnowledge: %v", err)
	}

	// 確認
	got, err = svc.GetKnowledge(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetKnowledge: %v", err)
	}
	if got.IsPinned {
		t.Error("expected IsPinned = false")
	}
}

func TestListPinnedKnowledge(t *testing.T) {
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

	// ピン留め付きでノート作成
	pinned, err := svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "Pinned Note",
		Body:         "Pinned content",
		WorkingScope: "test",
		IsPinned:     true,
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}

	// ピン留めなしでノート作成
	_, err = svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "Unpinned Note",
		Body:         "Unpinned content",
		WorkingScope: "test",
		IsPinned:     false,
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}

	// ピン留め一覧取得
	notes, err := svc.ListPinnedKnowledge(ctx, "test", 10)
	if err != nil {
		t.Fatalf("ListPinnedKnowledge: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 pinned note, got %d", len(notes))
	}
	if notes[0].ID != pinned.ID {
		t.Errorf("expected pinned note ID %q, got %q", pinned.ID, notes[0].ID)
	}
}

func TestSearchKnowledge(t *testing.T) {
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

	// 用語定義ノート作成
	_, err = svc.IngestKnowledge(ctx, IngestKnowledgeRequest{
		Title:        "マイクロサービス",
		Body:         "マイクロサービスは、アプリケーションを小さなサービスに分割するアーキテクチャパターンです",
		WorkingScope: "glossary",
	})
	if err != nil {
		t.Fatalf("IngestKnowledge: %v", err)
	}

	// 検索
	notes, err := svc.SearchKnowledge(ctx, "マイクロサービス", 10)
	if err != nil {
		t.Fatalf("SearchKnowledge: %v", err)
	}
	if len(notes) < 1 {
		t.Error("expected at least 1 result")
	}
}
