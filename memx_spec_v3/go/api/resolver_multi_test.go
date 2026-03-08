package api

import (
	"context"
	"testing"
)

func TestMultiStoreResolver_ResolveRef(t *testing.T) {
	// Create test data
	shortNotes := []Note{
		{ID: "short-1", Title: "Short Note", Body: "Body", Summary: "Summary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}},
	}
	journalNotes := []JournalNote{
		{NoteBase: NoteBase{ID: "journal-1", Title: "Journal Note", Body: "JBody", Summary: "JSummary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "journal-1"}}, WorkingScope: "project:test"},
	}
	knowledgeNotes := []KnowledgeNote{
		{NoteBase: NoteBase{ID: "knowledge-1", Title: "Knowledge Note", Body: "KBody", Summary: "KSummary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeKnowledge, ID: "knowledge-1"}}, WorkingScope: "glossary"},
	}
	archiveNotes := []ArchiveNote{
		{ID: "archive-1", Title: "Archive Note", Body: "ABody", Summary: "ASummary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "archive-1"}},
	}

	shortSearch := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return shortNotes, nil
	}
	shortShow := func(ctx context.Context, id string) (*Note, error) {
		for i := range shortNotes {
			if shortNotes[i].ID == id {
				return &shortNotes[i], nil
			}
		}
		return nil, nil
	}
	journalSearch := func(ctx context.Context, query string, topK int) ([]JournalNote, error) {
		return journalNotes, nil
	}
	journalShow := func(ctx context.Context, id string) (*JournalNote, error) {
		for i := range journalNotes {
			if journalNotes[i].ID == id {
				return &journalNotes[i], nil
			}
		}
		return nil, nil
	}
	knowledgeSearch := func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error) {
		return knowledgeNotes, nil
	}
	knowledgeShow := func(ctx context.Context, id string) (*KnowledgeNote, error) {
		for i := range knowledgeNotes {
			if knowledgeNotes[i].ID == id {
				return &knowledgeNotes[i], nil
			}
		}
		return nil, nil
	}
	archiveShow := func(ctx context.Context, id string) (*ArchiveNote, error) {
		for i := range archiveNotes {
			if archiveNotes[i].ID == id {
				return &archiveNotes[i], nil
			}
		}
		return nil, nil
	}

	resolver := NewMultiStoreResolver(shortSearch, shortShow, journalSearch, journalShow, knowledgeSearch, knowledgeShow, archiveShow)

	tests := []struct {
		name       string
		ref        TypedRef
		wantStatus RefStatus
		wantErr    bool
	}{
		{
			name:       "resolve evidence from short store",
			ref:        TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"},
			wantStatus: RefStatusResolved,
			wantErr:    false,
		},
		{
			name:       "resolve evidence from journal store",
			ref:        TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "journal-1"},
			wantStatus: RefStatusResolved,
			wantErr:    false,
		},
		{
			name:       "resolve evidence from archive (not in short or journal)",
			ref:        TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "archive-1"},
			wantStatus: RefStatusResolved,
			wantErr:    false,
		},
		{
			name:       "resolve knowledge from knowledge store",
			ref:        TypedRef{Domain: DomainMemx, Type: EntityTypeKnowledge, ID: "knowledge-1"},
			wantStatus: RefStatusResolved,
			wantErr:    false,
		},
		{
			name:       "resolve artifact from knowledge store",
			ref:        TypedRef{Domain: DomainMemx, Type: EntityTypeArtifact, ID: "knowledge-1"},
			wantStatus: RefStatusResolved,
			wantErr:    false,
		},
		{
			name:       "unresolved evidence (not found)",
			ref:        TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "nonexistent"},
			wantStatus: RefStatusUnresolved,
			wantErr:    false,
		},
		{
			name:       "unsupported domain",
			ref:        TypedRef{Domain: "other", Type: EntityTypeEvidence, ID: "short-1"},
			wantStatus: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved, err := resolver.ResolveRef(context.Background(), tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resolved.Status != tt.wantStatus {
				t.Errorf("ResolveRef() status = %v, want %v", resolved.Status, tt.wantStatus)
			}
		})
	}
}

func TestMultiStoreResolver_ResolveMany(t *testing.T) {
	shortNotes := []Note{
		{ID: "short-1", Title: "Short Note", Body: "Body", Summary: "Summary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}},
	}
	knowledgeNotes := []KnowledgeNote{
		{NoteBase: NoteBase{ID: "knowledge-1", Title: "Knowledge Note", Body: "KBody", Summary: "KSummary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeKnowledge, ID: "knowledge-1"}}, WorkingScope: "glossary"},
	}

	shortSearch := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return shortNotes, nil
	}
	shortShow := func(ctx context.Context, id string) (*Note, error) {
		for i := range shortNotes {
			if shortNotes[i].ID == id {
				return &shortNotes[i], nil
			}
		}
		return nil, nil
	}
	knowledgeSearch := func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error) {
		return knowledgeNotes, nil
	}
	knowledgeShow := func(ctx context.Context, id string) (*KnowledgeNote, error) {
		for i := range knowledgeNotes {
			if knowledgeNotes[i].ID == id {
				return &knowledgeNotes[i], nil
			}
		}
		return nil, nil
	}
	noJournalSearch := func(ctx context.Context, query string, topK int) ([]JournalNote, error) {
		return nil, nil
	}
	noJournalShow := func(ctx context.Context, id string) (*JournalNote, error) {
		return nil, nil
	}
	noArchiveShow := func(ctx context.Context, id string) (*ArchiveNote, error) {
		return nil, nil
	}

	resolver := NewMultiStoreResolver(shortSearch, shortShow, noJournalSearch, noJournalShow, knowledgeSearch, knowledgeShow, noArchiveShow)

	refs := []TypedRef{
		{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"},
		{Domain: DomainMemx, Type: EntityTypeKnowledge, ID: "knowledge-1"},
		{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "nonexistent"},
		{Domain: "other", Type: EntityTypeEvidence, ID: "unknown"},
	}

	report, err := resolver.ResolveMany(context.Background(), refs)
	if err != nil {
		t.Fatalf("ResolveMany() error = %v", err)
	}

	if len(report.Resolved) != 2 {
		t.Errorf("Expected 2 resolved, got %d", len(report.Resolved))
	}
	if len(report.Unresolved) != 1 {
		t.Errorf("Expected 1 unresolved, got %d", len(report.Unresolved))
	}
	if len(report.Unsupported) != 1 {
		t.Errorf("Expected 1 unsupported, got %d", len(report.Unsupported))
	}
	if !report.Diagnostics.PartialBundle {
		t.Error("Expected PartialBundle to be true")
	}
}

func TestMultiStoreResolver_LoadSummary(t *testing.T) {
	shortNotes := []Note{
		{ID: "short-1", Title: "Test", Body: "Body", Summary: "Test Summary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}},
	}

	shortSearch := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return shortNotes, nil
	}
	shortShow := func(ctx context.Context, id string) (*Note, error) {
		for i := range shortNotes {
			if shortNotes[i].ID == id {
				return &shortNotes[i], nil
			}
		}
		return nil, nil
	}
	noSearch := func(ctx context.Context, query string, topK int) ([]JournalNote, error) {
		return nil, nil
	}
	noShow := func(ctx context.Context, id string) (*JournalNote, error) {
		return nil, nil
	}
	noKnowledgeSearch := func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error) {
		return nil, nil
	}
	noKnowledgeShow := func(ctx context.Context, id string) (*KnowledgeNote, error) {
		return nil, nil
	}
	noArchiveShow := func(ctx context.Context, id string) (*ArchiveNote, error) {
		return nil, nil
	}

	resolver := NewMultiStoreResolver(shortSearch, shortShow, noSearch, noShow, noKnowledgeSearch, noKnowledgeShow, noArchiveShow)

	t.Run("existing ref", func(t *testing.T) {
		ref := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}
		payload, err := resolver.LoadSummary(context.Background(), ref)
		if err != nil {
			t.Fatalf("LoadSummary() error = %v", err)
		}
		if !payload.Exists {
			t.Error("Expected Exists to be true")
		}
		if payload.Summary != "Test Summary" {
			t.Errorf("Expected summary 'Test Summary', got '%s'", payload.Summary)
		}
	})

	t.Run("non-existing ref", func(t *testing.T) {
		ref := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "nonexistent"}
		payload, err := resolver.LoadSummary(context.Background(), ref)
		if err != nil {
			t.Fatalf("LoadSummary() error = %v", err)
		}
		if payload.Exists {
			t.Error("Expected Exists to be false")
		}
	})
}

func TestMultiStoreResolver_LoadSelectedRaw(t *testing.T) {
	shortNotes := []Note{
		{ID: "short-1", Title: "Test", Body: "Full Body Content", Summary: "Summary Content", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}},
	}
	journalNotes := []JournalNote{
		{NoteBase: NoteBase{ID: "journal-1", Title: "Journal", Body: "Journal Body", Summary: "Journal Summary", Ref: TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "journal-1"}}, WorkingScope: "project:test"},
	}

	shortSearch := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return shortNotes, nil
	}
	shortShow := func(ctx context.Context, id string) (*Note, error) {
		for i := range shortNotes {
			if shortNotes[i].ID == id {
				return &shortNotes[i], nil
			}
		}
		return nil, nil
	}
	journalSearch := func(ctx context.Context, query string, topK int) ([]JournalNote, error) {
		return journalNotes, nil
	}
	journalShow := func(ctx context.Context, id string) (*JournalNote, error) {
		for i := range journalNotes {
			if journalNotes[i].ID == id {
				return &journalNotes[i], nil
			}
		}
		return nil, nil
	}
	noKnowledgeSearch := func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error) {
		return nil, nil
	}
	noKnowledgeShow := func(ctx context.Context, id string) (*KnowledgeNote, error) {
		return nil, nil
	}
	noArchiveShow := func(ctx context.Context, id string) (*ArchiveNote, error) {
		return nil, nil
	}

	resolver := NewMultiStoreResolver(shortSearch, shortShow, journalSearch, journalShow, noKnowledgeSearch, noKnowledgeShow, noArchiveShow)

	t.Run("with body", func(t *testing.T) {
		ref := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}
		payload, err := resolver.LoadSelectedRaw(context.Background(), ref, RawSelector{IncludeBody: true})
		if err != nil {
			t.Fatalf("LoadSelectedRaw() error = %v", err)
		}
		if !payload.Found {
			t.Error("Expected Found to be true")
		}
		if payload.Raw != "Full Body Content" {
			t.Errorf("Expected 'Full Body Content', got '%s'", payload.Raw)
		}
	})

	t.Run("without body (summary only)", func(t *testing.T) {
		ref := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "short-1"}
		payload, err := resolver.LoadSelectedRaw(context.Background(), ref, RawSelector{IncludeBody: false})
		if err != nil {
			t.Fatalf("LoadSelectedRaw() error = %v", err)
		}
		if !payload.Found {
			t.Error("Expected Found to be true")
		}
		if payload.Raw != "Summary Content" {
			t.Errorf("Expected 'Summary Content', got '%s'", payload.Raw)
		}
	})

	t.Run("journal fallback", func(t *testing.T) {
		ref := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, ID: "journal-1"}
		payload, err := resolver.LoadSelectedRaw(context.Background(), ref, RawSelector{IncludeBody: true})
		if err != nil {
			t.Fatalf("LoadSelectedRaw() error = %v", err)
		}
		if !payload.Found {
			t.Error("Expected Found to be true")
		}
		if payload.Raw != "Journal Body" {
			t.Errorf("Expected 'Journal Body', got '%s'", payload.Raw)
		}
	})
}
