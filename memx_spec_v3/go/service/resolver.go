package service

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

// TypedRef は typed_ref の service パッケージ内表現。
// api.TypedRef との変換は境界層で行う。
type TypedRef struct {
	Domain     string
	Type       string
	Provider   string
	ID         string
}

func (r TypedRef) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", r.Domain, r.Type, r.Provider, r.ID)
}

func (r TypedRef) IsValid() bool {
	return r.Domain != "" && r.Type != "" && r.Provider != "" && r.ID != ""
}

// ResolvedRef は解決結果を表す。
type ResolvedRef struct {
	Ref      TypedRef
	Status   RefStatus
	Summary  string
	Metadata RefMetadata
}

// RefStatus は解決状態を表す。
type RefStatus string

const (
	RefStatusResolved    RefStatus = "resolved"
	RefStatusUnresolved  RefStatus = "unresolved"
	RefStatusUnsupported RefStatus = "unsupported"
)

// RefMetadata は参照先のメタデータ。
type RefMetadata struct {
	Title       string
	SourceType  string
	Sensitivity string
	CreatedAt   string
	UpdatedAt   string
}

// ResolveReport は一括解決結果を表す。
type ResolveReport struct {
	Resolved    []ResolvedRef
	Unresolved  []TypedRef
	Unsupported []TypedRef
	Diagnostics ResolverDiagnostics
}

// ResolverDiagnostics は解決診断情報。
type ResolverDiagnostics struct {
	MissingRefs     []TypedRef
	UnsupportedRefs []TypedRef
	ResolverWarnings []string
	PartialBundle   bool
}

// SummaryPayload は要約取得結果。
type SummaryPayload struct {
	Ref     TypedRef
	Summary string
	Exists  bool
}

// RawSelector は raw データ取得時のセレクタ。
type RawSelector struct {
	IncludeBody bool
	IncludeTags bool
	Fields      []string
}

// RawPayload は raw データ取得結果。
type RawPayload struct {
	Ref   TypedRef
	Raw   string
	Found bool
}

// ShortNoteResolver は short ストアを使った Resolver の実装。
// P4 Phase 3B: 現時点の memx-core 実装に合わせた最小 adapter。
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
	if ref.Domain != "memx" {
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
	if ref.Domain != "memx" {
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
	if ref.Domain != "memx" {
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

// ErrUnsupportedRef は未対応の typed_ref が指定された場合のエラー。
type ErrUnsupportedRef struct {
	Ref TypedRef
}

func (e *ErrUnsupportedRef) Error() string {
	return fmt.Sprintf("unsupported typed_ref: %s", e.Ref)
}

// ErrUnresolvedRef は解決できない typed_ref が指定された場合のエラー。
type ErrUnresolvedRef struct {
	Ref    TypedRef
	Reason string
}

func (e *ErrUnresolvedRef) Error() string {
	return fmt.Sprintf("unresolved typed_ref: %s (%s)", e.Ref, e.Reason)
}

// ValidateTypedRefForResolve は解決可能な typed_ref かどうかを検証する。
func ValidateTypedRefForResolve(ref TypedRef) error {
	if !ref.IsValid() {
		return fmt.Errorf("invalid typed_ref: %s", ref)
	}

	// 現状は memx ドメインのみ対応
	if ref.Domain != "memx" {
		return &ErrUnsupportedRef{Ref: ref}
	}

	// provider は local のみ対応
	if ref.Provider != "local" {
		return &ErrUnsupportedRef{Ref: ref}
	}

	return nil
}