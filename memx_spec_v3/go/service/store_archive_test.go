package service

import (
	"context"
	"path/filepath"
	"testing"

	"memx/db"
)

func TestArchiveNoteFromShort(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal: filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 1. short にノート作成
	shortNote, err := svc.IngestShort(ctx, IngestNoteRequest{
		Title: "Archive Test",
		Body:  "This will be archived",
	})
	if err != nil {
		t.Fatalf("IngestShort: %v", err)
	}

	// 2. archive へ退避
	archiveNote, err := svc.ArchiveNoteFromShort(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("ArchiveNoteFromShort: %v", err)
	}

	if archiveNote.ID != shortNote.ID {
		t.Errorf("ID mismatch: got %q, want %q", archiveNote.ID, shortNote.ID)
	}
	if archiveNote.Title != shortNote.Title {
		t.Errorf("Title mismatch: got %q, want %q", archiveNote.Title, shortNote.Title)
	}

	// 3. short からは削除されていることを確認
	_, err = svc.GetShort(ctx, shortNote.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for short note, got: %v", err)
	}

	// 4. archive から取得できることを確認
	got, err := svc.GetArchive(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("GetArchive: %v", err)
	}
	if got.ID != shortNote.ID {
		t.Errorf("GetArchive ID mismatch: got %q, want %q", got.ID, shortNote.ID)
	}
}

func TestRestoreFromArchive(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal: filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 1. short にノート作成 → archive へ退避
	shortNote, err := svc.IngestShort(ctx, IngestNoteRequest{
		Title: "Restore Test",
		Body:  "This will be restored",
	})
	if err != nil {
		t.Fatalf("IngestShort: %v", err)
	}

	_, err = svc.ArchiveNoteFromShort(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("ArchiveNoteFromShort: %v", err)
	}

	// 2. archive から復元
	restoredNote, err := svc.RestoreFromArchive(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("RestoreFromArchive: %v", err)
	}

	if restoredNote.ID != shortNote.ID {
		t.Errorf("ID mismatch: got %q, want %q", restoredNote.ID, shortNote.ID)
	}
	if restoredNote.Title != shortNote.Title {
		t.Errorf("Title mismatch: got %q, want %q", restoredNote.Title, shortNote.Title)
	}

	// 3. short に戻っていることを確認
	got, err := svc.GetShort(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("GetShort after restore: %v", err)
	}
	if got.ID != shortNote.ID {
		t.Errorf("GetShort ID mismatch: got %q, want %q", got.ID, shortNote.ID)
	}

	// 4. archive からは削除されていることを確認
	_, err = svc.GetArchive(ctx, shortNote.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for archive note after restore, got: %v", err)
	}
}

func TestGetArchive(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal: filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 存在しないID
	_, err = svc.GetArchive(ctx, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// 無効なID形式
	_, err = svc.GetArchive(ctx, "invalid-id")
	if err == nil {
		t.Error("expected error for invalid id format")
	}
}

func TestListArchive(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal: filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// 複数ノートをアーカイブ
	for i := 0; i < 3; i++ {
		shortNote, err := svc.IngestShort(ctx, IngestNoteRequest{
			Title: "Archive Test",
			Body:  "Test body",
		})
		if err != nil {
			t.Fatalf("IngestShort: %v", err)
		}
		_, err = svc.ArchiveNoteFromShort(ctx, shortNote.ID)
		if err != nil {
			t.Fatalf("ArchiveNoteFromShort: %v", err)
		}
	}

	// 一覧取得
	notes, err := svc.ListArchive(ctx, 10)
	if err != nil {
		t.Fatalf("ListArchive: %v", err)
	}
	if len(notes) != 3 {
		t.Errorf("expected 3 notes, got %d", len(notes))
	}
}

func TestGetArchiveLineage(t *testing.T) {
	ctx := context.Background()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal: filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer svc.Close()

	// short → archive
	shortNote, err := svc.IngestShort(ctx, IngestNoteRequest{
		Title: "Lineage Test",
		Body:  "Test body",
	})
	if err != nil {
		t.Fatalf("IngestShort: %v", err)
	}

	_, err = svc.ArchiveNoteFromShort(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("ArchiveNoteFromShort: %v", err)
	}

	// lineage 確認
	lineage, err := svc.GetArchiveLineage(ctx, shortNote.ID)
	if err != nil {
		t.Fatalf("GetArchiveLineage: %v", err)
	}
	if len(lineage) == 0 {
		t.Error("expected at least one lineage record")
		return
	}

	found := false
	for _, l := range lineage {
		if l.Relation == "archived_from" && l.SrcNoteID == shortNote.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected archived_from lineage record")
	}
}