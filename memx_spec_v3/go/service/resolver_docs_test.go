package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"memx/db"
)

func newResolverServiceForTest(t *testing.T) *Service {
	t.Helper()
	tmpDir := t.TempDir()
	svc, err := New(db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return svc
}

func TestResolverDocsLifecycle(t *testing.T) {
	svc := newResolverServiceForTest(t)
	defer func() { _ = svc.Close() }()
	ctx := context.Background()

	doc, chunkCount, err := svc.DocsIngest(ctx, DocsIngestRequest{
		DocType:     "spec",
		Title:       "Memory Import Spec",
		SourcePath:  "docs/specs/memory-import.md",
		Version:     "2026-03-10",
		FeatureKeys: []string{"memory-import"},
		TaskIDs:     []string{"task:feature:local:123"},
		Tags:        []string{"memory", "import"},
		Body: `# Memory Import Spec

## Acceptance Criteria
- imports can be replayed safely

## Forbidden Patterns
- do not skip stale check

## Definition of Done
- contract response contains required docs

## Dependencies
- doc:spec:shared-schema`,
	})
	if err != nil {
		t.Fatalf("DocsIngest: %v", err)
	}
	if chunkCount < 4 {
		t.Fatalf("expected multiple chunks, got %d", chunkCount)
	}
	if doc.DocID == "" {
		t.Fatal("expected generated doc id")
	}

	required, recommended, err := svc.DocsResolve(ctx, DocsResolveRequest{Feature: "memory-import"})
	if err != nil {
		t.Fatalf("DocsResolve: %v", err)
	}
	if len(required) != 1 || required[0].DocID != doc.DocID {
		t.Fatalf("unexpected required docs: %#v", required)
	}
	if len(recommended) != 0 {
		t.Fatalf("unexpected recommended docs: %#v", recommended)
	}

	_, chunks, err := svc.ChunksGet(ctx, ChunksGetRequest{DocID: doc.DocID, Heading: "Acceptance", Limit: 5})
	if err != nil {
		t.Fatalf("ChunksGet: %v", err)
	}
	if len(chunks) != 1 || chunks[0].Importance != "required" {
		t.Fatalf("unexpected chunks: %#v", chunks)
	}

	receipt, err := svc.ReadsAck(ctx, ReadsAckRequest{TaskID: "task:feature:local:123", DocID: doc.DocID, ChunkIDs: []string{chunks[0].ChunkID}})
	if err != nil {
		t.Fatalf("ReadsAck: %v", err)
	}
	if receipt.Version != "2026-03-10" {
		t.Fatalf("unexpected receipt version: %#v", receipt)
	}

	_, _, err = svc.DocsIngest(ctx, DocsIngestRequest{
		DocID:       doc.DocID,
		DocType:     "spec",
		Title:       "Memory Import Spec",
		SourcePath:  "docs/specs/memory-import.md",
		Version:     "2026-03-11",
		FeatureKeys: []string{"memory-import"},
		TaskIDs:     []string{"task:feature:local:123"},
		Body: `# Memory Import Spec

## Acceptance Criteria
- imports can be replayed safely`,
	})
	if err != nil {
		t.Fatalf("DocsIngest update: %v", err)
	}

	stale, err := svc.DocsStaleCheck(ctx, DocsStaleCheckRequest{TaskID: "task:feature:local:123"})
	if err != nil {
		t.Fatalf("DocsStaleCheck: %v", err)
	}
	if len(stale) != 1 || stale[0].CurrentVersion != "2026-03-11" {
		t.Fatalf("unexpected stale response: %#v", stale)
	}

	required, acceptance, forbidden, done, dependencies, err := svc.ContractsResolve(ctx, ContractsResolveRequest{TaskID: "task:feature:local:123"})
	if err != nil {
		t.Fatalf("ContractsResolve: %v", err)
	}
	if len(required) != 1 {
		t.Fatalf("unexpected required docs: %#v", required)
	}
	if len(acceptance) != 1 || acceptance[0] != "imports can be replayed safely" {
		t.Fatalf("unexpected acceptance criteria: %#v", acceptance)
	}
	if len(forbidden) != 0 {
		t.Fatalf("expected no forbidden patterns after update, got %#v", forbidden)
	}
	if len(done) != 0 {
		t.Fatalf("expected no definition of done after update, got %#v", done)
	}
	if len(dependencies) != 0 {
		t.Fatalf("expected no dependencies after update, got %#v", dependencies)
	}
}

func TestDocsIngestRejectsOlderVersion(t *testing.T) {
	svc := newResolverServiceForTest(t)
	defer func() { _ = svc.Close() }()
	ctx := context.Background()

	doc, _, err := svc.DocsIngest(ctx, DocsIngestRequest{DocType: "spec", Title: "Versioned Spec", Version: "2026-03-11", Body: "# Spec"})
	if err != nil {
		t.Fatalf("initial DocsIngest: %v", err)
	}

	_, _, err = svc.DocsIngest(ctx, DocsIngestRequest{DocID: doc.DocID, DocType: "spec", Title: "Versioned Spec", Version: "2026-03-10", Body: "# Older"})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestDocsStaleCheckUsesLatestReceiptPerDoc(t *testing.T) {
	svc := newResolverServiceForTest(t)
	defer func() { _ = svc.Close() }()
	ctx := context.Background()

	doc, _, err := svc.DocsIngest(ctx, DocsIngestRequest{DocType: "spec", Title: "Latest Receipt Spec", Version: "2026-03-10", TaskIDs: []string{"task:latest"}, Body: "# Spec"})
	if err != nil {
		t.Fatalf("DocsIngest: %v", err)
	}
	if _, err := svc.ReadsAck(ctx, ReadsAckRequest{TaskID: "task:latest", DocID: doc.DocID, Version: "2026-03-09"}); err != nil {
		t.Fatalf("ReadsAck old: %v", err)
	}
	if _, err := svc.ReadsAck(ctx, ReadsAckRequest{TaskID: "task:latest", DocID: doc.DocID, Version: "2026-03-10"}); err != nil {
		t.Fatalf("ReadsAck latest: %v", err)
	}

	stale, err := svc.DocsStaleCheck(ctx, DocsStaleCheckRequest{TaskID: "task:latest"})
	if err != nil {
		t.Fatalf("DocsStaleCheck: %v", err)
	}
	if len(stale) != 0 {
		t.Fatalf("expected latest receipt to clear stale, got %#v", stale)
	}
}

func TestDocsIngestFixedChunking(t *testing.T) {
	svc := newResolverServiceForTest(t)
	defer func() { _ = svc.Close() }()
	ctx := context.Background()

	doc, chunkCount, err := svc.DocsIngest(ctx, DocsIngestRequest{
		DocType: "spec",
		Title:   "Fixed Chunk Spec",
		Version: "2026-03-10",
		Body:    "0123456789ABCDEFGHIJ0123456789ABCDEFGHIJ",
		Chunking: ChunkingOptions{
			Mode:     "fixed",
			MaxChars: 10,
		},
	})
	if err != nil {
		t.Fatalf("DocsIngest fixed: %v", err)
	}
	if chunkCount < 4 {
		t.Fatalf("expected fixed chunks, got %d", chunkCount)
	}
	_, chunks, err := svc.ChunksGet(ctx, ChunksGetRequest{DocID: doc.DocID})
	if err != nil {
		t.Fatalf("ChunksGet: %v", err)
	}
	if chunks[0].Heading != "Fixed Chunk Spec" {
		t.Fatalf("expected fixed chunk heading to stay on title, got %#v", chunks[0])
	}
	if len(chunks[0].Body) > 10 {
		t.Fatalf("expected fixed chunk body length <= 10, got %d", len(chunks[0].Body))
	}
}
