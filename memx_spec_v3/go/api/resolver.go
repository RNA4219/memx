package api

import (
	"context"
	"fmt"
)

// Resolver は typed_ref からエンティティを解決するインターフェース。
// P4 Phase 3B: current memx-core adapter 用の最小実装。
type Resolver interface {
	// ResolveRef は単一の typed_ref を解決する。
	ResolveRef(ctx context.Context, ref TypedRef) (ResolvedRef, error)

	// ResolveMany は複数の typed_ref を一括解決する。
	ResolveMany(ctx context.Context, refs []TypedRef) (*ResolveReport, error)

	// LoadSummary は要約を取得する（summary-first retrieval）。
	LoadSummary(ctx context.Context, ref TypedRef) (*SummaryPayload, error)

	// LoadSelectedRaw は必要時のみ raw データを取得する。
	LoadSelectedRaw(ctx context.Context, ref TypedRef, selector RawSelector) (*RawPayload, error)
}

// ResolvedRef は解決結果を表す。
type ResolvedRef struct {
	Ref      TypedRef   `json:"ref"`
	Status   RefStatus  `json:"status"`
	Summary  string     `json:"summary,omitempty"`
	Metadata RefMetadata `json:"metadata,omitempty"`
}

// RefStatus は解決状態を表す。
type RefStatus string

const (
	RefStatusResolved   RefStatus = "resolved"
	RefStatusUnresolved RefStatus = "unresolved"
	RefStatusUnsupported RefStatus = "unsupported"
)

// RefMetadata は参照先のメタデータ。
type RefMetadata struct {
	Title       string `json:"title,omitempty"`
	SourceType  string `json:"source_type,omitempty"`
	Sensitivity string `json:"sensitivity,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// ResolveReport は一括解決結果を表す。
type ResolveReport struct {
	Resolved   []ResolvedRef `json:"resolved"`
	Unresolved []TypedRef    `json:"unresolved"`
	Unsupported []TypedRef   `json:"unsupported"`
	Diagnostics ResolverDiagnostics `json:"diagnostics"`
}

// ResolverDiagnostics は解決診断情報。
type ResolverDiagnostics struct {
	MissingRefs    []TypedRef `json:"missing_refs"`
	UnsupportedRefs []TypedRef `json:"unsupported_refs"`
	ResolverWarnings []string `json:"resolver_warnings"`
	PartialBundle   bool      `json:"partial_bundle"`
}

// SummaryPayload は要約取得結果。
type SummaryPayload struct {
	Ref     TypedRef `json:"ref"`
	Summary string   `json:"summary"`
	Exists  bool     `json:"exists"`
}

// RawSelector は raw データ取得時のセレクタ。
type RawSelector struct {
	IncludeBody   bool     `json:"include_body"`
	IncludeTags   bool     `json:"include_tags"`
	Fields        []string `json:"fields,omitempty"`
}

// RawPayload は raw データ取得結果。
type RawPayload struct {
	Ref   TypedRef `json:"ref"`
	Raw   string   `json:"raw,omitempty"`
	Found bool     `json:"found"`
}

// ErrUnsupportedRef は未対応の typed_ref が指定された場合のエラー。
type ErrUnsupportedRef struct {
	Ref TypedRef
}

func (e *ErrUnsupportedRef) Error() string {
	return fmt.Sprintf("unsupported typed_ref: %s", e.Ref)
}

// ErrUnresolvedRef は解決できない typed_ref が指定された場合のエラー。
type ErrUnresolvedRef struct {
	Ref TypedRef
	Reason string
}

func (e *ErrUnresolvedRef) Error() string {
	return fmt.Sprintf("unresolved typed_ref: %s (%s)", e.Ref, e.Reason)
}

// ShortNoteResolver は short ストアを使った Resolver の実装。
// apiパッケージ内で完結する実装（循環参照回避）。
type ShortNoteResolver struct {
	searchFunc func(ctx context.Context, query string, topK int) ([]Note, error)
	showFunc   func(ctx context.Context, id string) (*Note, error)
}

// NewShortNoteResolver は ShortNoteResolver を作成する。
func NewShortNoteResolver(
	searchFunc func(ctx context.Context, query string, topK int) ([]Note, error),
	showFunc func(ctx context.Context, id string) (*Note, error),
) *ShortNoteResolver {
	return &ShortNoteResolver{
		searchFunc: searchFunc,
		showFunc:   showFunc,
	}
}

// ResolveRef は単一の typed_ref を解決する。
func (r *ShortNoteResolver) ResolveRef(ctx context.Context, ref TypedRef) (ResolvedRef, error) {
	// memx ドメインのみ対応
	if ref.Domain != DomainMemx {
		return ResolvedRef{}, &ErrUnsupportedRef{Ref: ref}
	}

	// 有効なエンティティタイプかチェック
	switch ref.Type {
	case EntityTypeEvidence, EntityTypeEvidenceChunk, EntityTypeKnowledge, EntityTypeArtifact, EntityTypeLineage:
		// OK
	default:
		return ResolvedRef{}, &ErrUnsupportedRef{Ref: ref}
	}

	// ID で参照解決
	note, err := r.showFunc(ctx, ref.ID)
	if err != nil {
		return ResolvedRef{
			Ref:    ref,
			Status: RefStatusUnresolved,
		}, nil
	}
	if note == nil {
		return ResolvedRef{
			Ref:    ref,
			Status: RefStatusUnresolved,
		}, nil
	}

	return ResolvedRef{
		Ref:     ref,
		Status:  RefStatusResolved,
		Summary: note.Summary,
		Metadata: RefMetadata{
			Title:       note.Title,
			SourceType:  note.SourceType,
			Sensitivity: note.Sensitivity,
			CreatedAt:   note.CreatedAt,
			UpdatedAt:   note.UpdatedAt,
		},
	}, nil
}

// ResolveMany は複数の typed_ref を一括解決する。
func (r *ShortNoteResolver) ResolveMany(ctx context.Context, refs []TypedRef) (*ResolveReport, error) {
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

	// 何か問題があれば partial bundle
	report.Diagnostics.PartialBundle = len(report.Unresolved) > 0 || len(report.Unsupported) > 0

	return report, nil
}

// LoadSummary は要約を取得する（summary-first retrieval）。
func (r *ShortNoteResolver) LoadSummary(ctx context.Context, ref TypedRef) (*SummaryPayload, error) {
	// memx ドメインのみ対応
	if ref.Domain != DomainMemx {
		return &SummaryPayload{Ref: ref, Exists: false}, &ErrUnsupportedRef{Ref: ref}
	}

	note, err := r.showFunc(ctx, ref.ID)
	if err != nil {
		return &SummaryPayload{Ref: ref, Exists: false}, nil
	}

	return &SummaryPayload{
		Ref:     ref,
		Summary: note.Summary,
		Exists:  true,
	}, nil
}

// LoadSelectedRaw は必要時のみ raw データを取得する。
func (r *ShortNoteResolver) LoadSelectedRaw(ctx context.Context, ref TypedRef, selector RawSelector) (*RawPayload, error) {
	// memx ドメインのみ対応
	if ref.Domain != DomainMemx {
		return &RawPayload{Ref: ref, Found: false}, &ErrUnsupportedRef{Ref: ref}
	}

	note, err := r.showFunc(ctx, ref.ID)
	if err != nil {
		return &RawPayload{Ref: ref, Found: false}, nil
	}

	// selector に応じて raw を構築
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

// ValidateTypedRefForResolve は解決可能な typed_ref かどうかを検証する。
func ValidateTypedRefForResolve(ref TypedRef) error {
	if !ref.IsValid() {
		return fmt.Errorf("invalid typed_ref: %s", ref)
	}

	// 現状は memx ドメインのみ対応
	if ref.Domain != DomainMemx {
		return &ErrUnsupportedRef{Ref: ref}
	}

	// provider は local のみ対応
	if ref.Provider != ProviderLocal {
		return &ErrUnsupportedRef{Ref: ref}
	}

	return nil
}