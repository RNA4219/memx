package api

import (
	"context"
)

// MultiStoreResolver は複数ストアを横断して typed_ref を解決する Resolver。
// P4 Phase 3C: short/journal/knowledge/archive を統合。
type MultiStoreResolver struct {
	shortResolver    *ShortNoteResolver
	journalResolver  *ShortNoteResolver
	knowledgeResolver *ShortNoteResolver
	archiveResolver  *ShortNoteResolver
}

// NewMultiStoreResolver は MultiStoreResolver を作成する。
func NewMultiStoreResolver(
	shortSearch func(ctx context.Context, query string, topK int) ([]Note, error),
	shortShow func(ctx context.Context, id string) (*Note, error),
	journalSearch func(ctx context.Context, query string, topK int) ([]JournalNote, error),
	journalShow func(ctx context.Context, id string) (*JournalNote, error),
	knowledgeSearch func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error),
	knowledgeShow func(ctx context.Context, id string) (*KnowledgeNote, error),
	archiveShow func(ctx context.Context, id string) (*ArchiveNote, error),
) *MultiStoreResolver {
	return &MultiStoreResolver{
		shortResolver: NewShortNoteResolver(shortSearch, shortShow),
		journalResolver: NewShortNoteResolver(
			adaptJournalSearch(journalSearch),
			adaptJournalShow(journalShow),
		),
		knowledgeResolver: NewShortNoteResolver(
			adaptKnowledgeSearch(knowledgeSearch),
			adaptKnowledgeShow(knowledgeShow),
		),
		archiveResolver: NewShortNoteResolver(
			nil, // archive has no search
			adaptArchiveShow(archiveShow),
		),
	}
}

// adaptJournalSearch は JournalNote の検索関数を Note に適応させる。
func adaptJournalSearch(search func(ctx context.Context, query string, topK int) ([]JournalNote, error)) func(ctx context.Context, query string, topK int) ([]Note, error) {
	return func(ctx context.Context, query string, topK int) ([]Note, error) {
		notes, err := search(ctx, query, topK)
		if err != nil {
			return nil, err
		}
		result := make([]Note, len(notes))
		for i, n := range notes {
			result[i] = Note(n.NoteBase)
		}
		return result, nil
	}
}

// adaptJournalShow は JournalNote の取得関数を Note に適応させる。
func adaptJournalShow(show func(ctx context.Context, id string) (*JournalNote, error)) func(ctx context.Context, id string) (*Note, error) {
	return func(ctx context.Context, id string) (*Note, error) {
		n, err := show(ctx, id)
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, nil
		}
		note := Note(n.NoteBase)
		return &note, nil
	}
}

// adaptKnowledgeSearch は KnowledgeNote の検索関数を Note に適応させる。
func adaptKnowledgeSearch(search func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error)) func(ctx context.Context, query string, topK int) ([]Note, error) {
	return func(ctx context.Context, query string, topK int) ([]Note, error) {
		notes, err := search(ctx, query, topK)
		if err != nil {
			return nil, err
		}
		result := make([]Note, len(notes))
		for i, n := range notes {
			result[i] = Note(n.NoteBase)
		}
		return result, nil
	}
}

// adaptKnowledgeShow は KnowledgeNote の取得関数を Note に適応させる。
func adaptKnowledgeShow(show func(ctx context.Context, id string) (*KnowledgeNote, error)) func(ctx context.Context, id string) (*Note, error) {
	return func(ctx context.Context, id string) (*Note, error) {
		n, err := show(ctx, id)
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, nil
		}
		note := Note(n.NoteBase)
		return &note, nil
	}
}

// adaptArchiveShow は ArchiveNote の取得関数を Note に適応させる。
func adaptArchiveShow(show func(ctx context.Context, id string) (*ArchiveNote, error)) func(ctx context.Context, id string) (*Note, error) {
	return func(ctx context.Context, id string) (*Note, error) {
		n, err := show(ctx, id)
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, nil
		}
		note := Note(*n)
		return &note, nil
	}
}

// ResolveRef は単一の typed_ref を解決する。
// エンティティタイプに応じて適切なストアを使用する。
func (r *MultiStoreResolver) ResolveRef(ctx context.Context, ref TypedRef) (ResolvedRef, error) {
	// memx ドメインのみ対応
	if ref.Domain != DomainMemx {
		return ResolvedRef{}, &ErrUnsupportedRef{Ref: ref}
	}

	// エンティティタイプに応じてルーティング
	var resolver *ShortNoteResolver
	switch ref.Type {
	case EntityTypeEvidence, EntityTypeEvidenceChunk:
		// evidence は short ストアを優先、見つからなければ archive
		resolved, err := r.shortResolver.ResolveRef(ctx, ref)
		if err == nil && resolved.Status == RefStatusResolved {
			return resolved, nil
		}
		// short になければ archive を試す
		return r.archiveResolver.ResolveRef(ctx, ref)

	case EntityTypeKnowledge:
		resolver = r.knowledgeResolver

	case EntityTypeArtifact:
		// artifact は現状 knowledge として扱う（将来分離）
		resolver = r.knowledgeResolver

	case EntityTypeLineage:
		// lineage は全ストアを検索
		for _, res := range []*ShortNoteResolver{
			r.shortResolver,
			r.journalResolver,
			r.knowledgeResolver,
			r.archiveResolver,
		} {
			if res == nil {
				continue
			}
			resolved, err := res.ResolveRef(ctx, ref)
			if err == nil && resolved.Status == RefStatusResolved {
				return resolved, nil
			}
		}
		return ResolvedRef{
			Ref:    ref,
			Status: RefStatusUnresolved,
		}, nil

	default:
		return ResolvedRef{}, &ErrUnsupportedRef{Ref: ref}
	}

	if resolver == nil {
		return ResolvedRef{
			Ref:    ref,
			Status: RefStatusUnresolved,
		}, nil
	}

	return resolver.ResolveRef(ctx, ref)
}

// ResolveMany は複数の typed_ref を一括解決する。
func (r *MultiStoreResolver) ResolveMany(ctx context.Context, refs []TypedRef) (*ResolveReport, error) {
	report := &ResolveReport{
		Resolved:    []ResolvedRef{},
		Unresolved:  []TypedRef{},
		Unsupported: []TypedRef{},
		Diagnostics: ResolverDiagnostics{
			MissingRefs:     []TypedRef{},
			UnsupportedRefs: []TypedRef{},
			ResolverWarnings: []string{},
			PartialBundle:   false,
		},
	}

	for _, ref := range refs {
		resolved, err := r.ResolveRef(ctx, ref)
		if err != nil {
			if _, ok := err.(*ErrUnsupportedRef); ok {
				report.Unsupported = append(report.Unsupported, ref)
				report.Diagnostics.UnsupportedRefs = append(report.Diagnostics.UnsupportedRefs, ref)
			} else {
				report.Unresolved = append(report.Unresolved, ref)
				report.Diagnostics.MissingRefs = append(report.Diagnostics.MissingRefs, ref)
			}
			continue
		}

		switch resolved.Status {
		case RefStatusResolved:
			report.Resolved = append(report.Resolved, resolved)
		case RefStatusUnresolved:
			report.Unresolved = append(report.Unresolved, ref)
			report.Diagnostics.MissingRefs = append(report.Diagnostics.MissingRefs, ref)
		case RefStatusUnsupported:
			report.Unsupported = append(report.Unsupported, ref)
			report.Diagnostics.UnsupportedRefs = append(report.Diagnostics.UnsupportedRefs, ref)
		}
	}

	report.Diagnostics.PartialBundle = len(report.Unresolved) > 0 || len(report.Unsupported) > 0

	return report, nil
}

// LoadSummary は要約を取得する（summary-first retrieval）。
func (r *MultiStoreResolver) LoadSummary(ctx context.Context, ref TypedRef) (*SummaryPayload, error) {
	resolved, err := r.ResolveRef(ctx, ref)
	if err != nil {
		return &SummaryPayload{Ref: ref, Exists: false}, err
	}
	if resolved.Status != RefStatusResolved {
		return &SummaryPayload{Ref: ref, Exists: false}, nil
	}
	return &SummaryPayload{
		Ref:     ref,
		Summary: resolved.Summary,
		Exists:  true,
	}, nil
}

// LoadSelectedRaw は必要時のみ raw データを取得する。
func (r *MultiStoreResolver) LoadSelectedRaw(ctx context.Context, ref TypedRef, selector RawSelector) (*RawPayload, error) {
	// memx ドメインのみ対応
	if ref.Domain != DomainMemx {
		return &RawPayload{Ref: ref, Found: false}, &ErrUnsupportedRef{Ref: ref}
	}

	// エンティティタイプに応じて適切なストアから raw を取得
	var showFunc func(ctx context.Context, id string) (*Note, error)

	switch ref.Type {
	case EntityTypeEvidence, EntityTypeEvidenceChunk:
		// short を優先、なければ archive
		showFunc = r.shortResolver.showFunc
		if showFunc == nil {
			showFunc = r.archiveResolver.showFunc
		}
	case EntityTypeKnowledge, EntityTypeArtifact:
		showFunc = r.knowledgeResolver.showFunc
	case EntityTypeLineage:
		// 最初に見つかったストアから取得
		for _, res := range []*ShortNoteResolver{
			r.shortResolver,
			r.journalResolver,
			r.knowledgeResolver,
			r.archiveResolver,
		} {
			if res != nil && res.showFunc != nil {
				showFunc = res.showFunc
				break
			}
		}
	}

	if showFunc == nil {
		return &RawPayload{Ref: ref, Found: false}, nil
	}

	note, err := showFunc(ctx, ref.ID)
	if err != nil {
		return &RawPayload{Ref: ref, Found: false}, nil
	}

	var raw string
	if selector.IncludeBody {
		raw = note.Body
	} else {
		raw = note.Summary
	}

	return &RawPayload{
		Ref:   ref,
		Raw:   raw,
		Found: true,
	}, nil
}