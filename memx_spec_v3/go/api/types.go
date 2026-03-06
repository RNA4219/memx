package api

// v1.3: API はツール/AI 向けの安定 I/F。

type ErrorCode string

const (
	CodeInvalidArgument ErrorCode = "INVALID_ARGUMENT"
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeGatekeepDeny    ErrorCode = "GATEKEEP_DENY"
	CodeInternal        ErrorCode = "INTERNAL"
)

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Note は API の返却モデル。文字列は RFC3339 の UTC を想定。
type Note struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	LastAccessedAt string `json:"last_accessed_at"`
	AccessCount    int64  `json:"access_count"`

	SourceType  string `json:"source_type"`
	Origin      string `json:"origin"`
	SourceTrust string `json:"source_trust"`
	Sensitivity string `json:"sensitivity"`
}

type NotesIngestRequest struct {
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	Summary     string   `json:"summary,omitempty"`
	SourceType  string   `json:"source_type,omitempty"`
	Origin      string   `json:"origin,omitempty"`
	SourceTrust string   `json:"source_trust,omitempty"`
	Sensitivity string   `json:"sensitivity,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type NotesIngestResponse struct {
	Note Note `json:"note"`
}

type NotesSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

type NotesSearchResponse struct {
	Notes []Note `json:"notes"`
}

type GCOptions struct {
	DryRun bool `json:"dry_run,omitempty"`
}

type GCRunRequest struct {
	Target  string    `json:"target"` // v1: "short" のみ想定
	Options GCOptions `json:"options,omitempty"`
}

type GCRunResponse struct {
	Status string `json:"status"` // "ok"
}

// SummarizeRequest は単一ノートの要約リクエスト。
type SummarizeRequest struct {
	ID string `json:"id"`
}

// SummarizeResponse は単一ノートの要約レスポンス。
type SummarizeResponse struct {
	Note Note `json:"note"`
}

// SummarizeBatchRequest は複数ノートの統合要約リクエスト。
type SummarizeBatchRequest struct {
	IDs []string `json:"ids"`
}

// SummarizeBatchResponse は複数ノートの統合要約レスポンス。
type SummarizeBatchResponse struct {
	Summary   string `json:"summary"`
	NoteCount int    `json:"note_count"`
}
