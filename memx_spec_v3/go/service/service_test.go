package service

import (
	"context"
	"path/filepath"
	"testing"

	"memx/db"
)

func TestService_IngestShort_Validation(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	tests := []struct {
		name    string
		req     IngestNoteRequest
		wantErr error
	}{
		{
			name:    "empty title",
			req:     IngestNoteRequest{Body: "test"},
			wantErr: ErrInvalidArgument,
		},
		{
			name:    "empty body",
			req:     IngestNoteRequest{Title: "test"},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "title too long",
			req: IngestNoteRequest{
				Title: string(make([]byte, 501)),
				Body:  "test",
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "body too long",
			req: IngestNoteRequest{
				Title: "test",
				Body:  string(make([]byte, 100001)),
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "invalid source_type",
			req: IngestNoteRequest{
				Title:      "test",
				Body:       "test",
				SourceType: "invalid",
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "invalid source_trust",
			req: IngestNoteRequest{
				Title:       "test",
				Body:        "test",
				SourceTrust: "invalid",
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "invalid sensitivity",
			req: IngestNoteRequest{
				Title:       "test",
				Body:        "test",
				Sensitivity: "invalid",
			},
			wantErr: ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.IngestShort(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr.Error() && err.Error()[:len("invalid argument")] != "invalid argument" {
				// Check if error contains the expected error message
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestService_IngestShort_PolicyDenied(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// secret sensitivity は deny される
	req := IngestNoteRequest{
		Title:       "secret note",
		Body:        "this is secret",
		Sensitivity: "secret",
	}

	_, err = svc.IngestShort(context.Background(), req)
	if err == nil {
		t.Fatal("expected policy denied error, got nil")
	}
	if err.Error() != ErrPolicyDenied.Error() && err.Error()[:len("policy denied")] != "policy denied" {
		t.Errorf("expected policy denied error, got: %v", err)
	}
}

func TestService_IngestShort_Success(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	req := IngestNoteRequest{
		Title:       "test note",
		Body:        "this is a test note",
		SourceType:  "manual",
		SourceTrust: "user_input",
		Sensitivity: "internal",
		Tags:        []string{"test", "example"},
	}

	note, err := svc.IngestShort(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if note.Title != req.Title {
		t.Errorf("expected title %q, got: %q", req.Title, note.Title)
	}
	if note.Body != req.Body {
		t.Errorf("expected body %q, got: %q", req.Body, note.Body)
	}
	if note.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestService_SearchShort_Validation(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// 空のクエリ
	_, err = svc.SearchShort(context.Background(), "", 10)
	if err == nil {
		t.Fatal("expected error for empty query")
	}

	// 長すぎるクエリ
	longQuery := string(make([]byte, 1001))
	_, err = svc.SearchShort(context.Background(), longQuery, 10)
	if err == nil {
		t.Fatal("expected error for long query")
	}
}

func TestService_GetShort_Validation(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	tests := []struct {
		name string
		id   string
	}{
		{"empty id", ""},
		{"wrong length", "abc"},
		{"non-hex", "ghijklmnopqrstuvwxyz123456789012"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetShort(context.Background(), tt.id)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestService_SummarizeNote_NoLLM(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create a note first
	note, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test",
		Body:  "test body",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	// Summarize without LLM configured should fail
	_, err = svc.SummarizeNote(context.Background(), note.ID)
	if err == nil {
		t.Fatal("expected error when MiniLLM not configured")
	}
}

func TestService_SummarizeNotes_NoLLM(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Create notes first
	note1, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test1",
		Body:  "test body 1",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}
	note2, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test2",
		Body:  "test body 2",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	// Summarize batch without ReflectLLM configured should fail
	_, err = svc.SummarizeNotes(context.Background(), []string{note1.ID, note2.ID})
	if err == nil {
		t.Fatal("expected error when ReflectLLM not configured")
	}
}

// mockMiniLLM is a mock implementation of MiniLLMClient for testing.
type mockMiniLLM struct {
	summary string
	err     error
}

func (m *mockMiniLLM) TagAndScore(ctx context.Context, noteBody string) (db.TagsAndScores, error) {
	return db.TagsAndScores{}, nil
}

func (m *mockMiniLLM) Summarize(ctx context.Context, title, body string) (db.SummarizeResult, error) {
	if m.err != nil {
		return db.SummarizeResult{}, m.err
	}
	return db.SummarizeResult{Summary: m.summary}, nil
}

func TestService_SummarizeNote_WithMockLLM(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Set mock LLM
	svc.SetMiniLLM(&mockMiniLLM{summary: "This is a generated summary."})

	// Create a note first
	note, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test note",
		Body:  "This is a long test body that should be summarized by the LLM.",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	// Summarize
	updatedNote, err := svc.SummarizeNote(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("failed to summarize note: %v", err)
	}

	if updatedNote.Summary != "This is a generated summary." {
		t.Errorf("expected summary %q, got: %q", "This is a generated summary.", updatedNote.Summary)
	}
}

func TestService_IngestShort_WithAutoSummary(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Set mock LLM
	svc.SetMiniLLM(&mockMiniLLM{summary: "Auto-generated summary."})

	// Ingest without explicit summary - should auto-generate
	note, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test note",
		Body:  "Test body for auto-summary.",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	if note.Summary != "Auto-generated summary." {
		t.Errorf("expected auto-generated summary, got: %q", note.Summary)
	}
}

func TestService_IngestShort_AutoSummaryFailureLogsWarning(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	logBuf, logger := newBufferedTestLogger()
	svc.SetLogger(logger)
	svc.SetMiniLLM(&mockMiniLLM{err: errTestLLM})

	note, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test note",
		Body:  "Test body for failed auto-summary.",
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}
	if note.Summary != "" {
		t.Errorf("expected empty summary on LLM failure, got: %q", note.Summary)
	}

	assertAutoSummaryWarningLogged(t, logBuf.String(), "short", "test note", errTestLLM)
}

func TestService_IngestShort_NoLLMFlag(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short: filepath.Join(tmpDir, "short.db"),
	}

	svc, err := New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	// Set mock LLM
	svc.SetMiniLLM(&mockMiniLLM{summary: "Should not be used."})

	// Ingest with NoLLM=true - should skip auto-summary
	note, err := svc.IngestShort(context.Background(), IngestNoteRequest{
		Title: "test note",
		Body:  "Test body.",
		NoLLM: true,
	})
	if err != nil {
		t.Fatalf("failed to ingest note: %v", err)
	}

	if note.Summary != "" {
		t.Errorf("expected empty summary with NoLLM=true, got: %q", note.Summary)
	}
}
