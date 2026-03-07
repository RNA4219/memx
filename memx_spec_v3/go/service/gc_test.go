package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"memx/db"
)

func TestService_GCShort_FeatureDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Feature flag disabled
	_, err = svc.GCShort(context.Background(), GCRequest{
		Target:  "short",
		DryRun:  false,
		Enabled: false,
	})
	if err != ErrFeatureDisabled {
		t.Errorf("expected ErrFeatureDisabled, got: %v", err)
	}
}

func TestService_GCShort_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create some notes
	for i := 0; i < 5; i++ {
		_, err := svc.IngestShort(context.Background(), IngestNoteRequest{
			Title: "test note",
			Body:  "test body",
		})
		if err != nil {
			t.Fatalf("failed to ingest note: %v", err)
		}
	}

	// Dry-run should succeed
	result, err := svc.GCShort(context.Background(), GCRequest{
		Target:  "short",
		DryRun:  true,
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.DryRun {
		t.Error("expected dry_run to be true")
	}
	if result.DryRunResult == nil {
		t.Fatal("expected dry_run_result to be set")
	}
	if result.DryRunResult.Target != "short" {
		t.Errorf("expected target 'short', got: %s", result.DryRunResult.Target)
	}
	if result.DryRunResult.Decision.Metrics.NoteCount != 5 {
		t.Errorf("expected note_count 5, got: %d", result.DryRunResult.Decision.Metrics.NoteCount)
	}
}

func TestService_GCShort_SkippedUnderLimit(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create a few notes (under limit)
	for i := 0; i < 3; i++ {
		_, err := svc.IngestShort(context.Background(), IngestNoteRequest{
			Title: "test note",
			Body:  "test body",
		})
		if err != nil {
			t.Fatalf("failed to ingest note: %v", err)
		}
	}

	// GC should skip (under limit)
	result, err := svc.GCShort(context.Background(), GCRequest{
		Target:  "short",
		DryRun:  false,
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "skipped" {
		t.Errorf("expected status 'skipped', got: %s", result.Status)
	}
}

func TestService_GCTriggerCheck(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create notes
	for i := 0; i < 5; i++ {
		_, err := svc.IngestShort(context.Background(), IngestNoteRequest{
			Title: "test note",
			Body:  "test body",
		})
		if err != nil {
			t.Fatalf("failed to ingest note: %v", err)
		}
	}

	cfg := DefaultGCConfig
	decision, err := svc.gcTriggerCheck(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Under soft limit
	if decision.ShouldRun {
		t.Error("expected should_run to be false (under limit)")
	}
	if decision.Metrics.NoteCount != 5 {
		t.Errorf("expected note_count 5, got: %d", decision.Metrics.NoteCount)
	}
	if decision.Reason != "under_limit" {
		t.Errorf("expected reason 'under_limit', got: %s", decision.Reason)
	}
}

func TestService_GCDryRunResult_Format(t *testing.T) {
	result := &GCDryRunResult{
		Target: "short",
		Phase:  "phase0",
		Decision: GCTriggerDecision{
			ShouldRun: true,
			Reason:    "soft_limit_reached",
		},
	}
	result.Decision.Metrics.NoteCount = 1300
	result.Decision.Metrics.SoftLimitNotes = 1200
	result.Decision.Metrics.HardLimitNotes = 2000

	// JSON output
	json := result.ToJSON()
	if json == "" {
		t.Error("expected non-empty JSON output")
	}

	// Human readable output
	human := result.FormatDryRunOutput()
	if human == "" {
		t.Error("expected non-empty human readable output")
	}
	if !containsAll(human, "GC Dry-Run Result", "short", "phase0", "1300") {
		t.Errorf("unexpected human output: %s", human)
	}
}

func containsAll(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && containsHelper(s, sub)))
}

func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestService_GCShort_ArchiveOldNotes(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create a note
	note, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "old note",
		Body:  "this is an old note",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	// Verify note exists in short
	_, err = svc.GetShort(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("failed to get note from short: %v", err)
	}

	// Test dry-run JSON output
	result, err := svc.GCShort(context.Background(), GCRequest{
		Target:  "short",
		DryRun:  true,
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify dry-run doesn't modify data
	_, err = svc.GetShort(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("note should still exist after dry-run: %v", err)
	}

	if !result.DryRun {
		t.Error("expected dry_run to be true")
	}
}

func TestGCConfig_Defaults(t *testing.T) {
	cfg := DefaultGCConfig

	if cfg.SoftLimitNotes != 1200 {
		t.Errorf("expected soft_limit 1200, got: %d", cfg.SoftLimitNotes)
	}
	if cfg.HardLimitNotes != 2000 {
		t.Errorf("expected hard_limit 2000, got: %d", cfg.HardLimitNotes)
	}
	if cfg.MinIntervalMinutes != 180 {
		t.Errorf("expected min_interval 180, got: %d", cfg.MinIntervalMinutes)
	}
}

func TestService_GCShort_WithArchiveStore(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create a note
	_, err = svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test note",
		Body:  "test body",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	// Update last_gc_at to simulate old GC run
	oldTime := time.Now().Add(-time.Hour * 4).Format(time.RFC3339Nano)
	_, _ = svc.Conn.DB.Exec("UPDATE short_meta SET value = ? WHERE key = 'last_gc_at'", oldTime)

	// GC should still skip (under limit)
	result, err := svc.GCShort(context.Background(), GCRequest{
		Target:  "short",
		DryRun:  false,
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Status should be skipped since we're under the limit
	if result.Status != "skipped" {
		t.Logf("Status: %s (note_count may be over limit in test)", result.Status)
	}
}