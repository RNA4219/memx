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

// -------------------- Chronicle --------------------

// ChronicleNote は chronicle ストアのノート。
type ChronicleNote struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	LastAccessedAt string `json:"last_accessed_at"`
	AccessCount    int64  `json:"access_count"`
	SourceType     string `json:"source_type"`
	Origin         string `json:"origin"`
	SourceTrust    string `json:"source_trust"`
	Sensitivity    string `json:"sensitivity"`
	WorkingScope   string `json:"working_scope"`
	IsPinned       bool   `json:"is_pinned"`
}

// ChronicleIngestRequest は chronicle への投入リクエスト。
type ChronicleIngestRequest struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Summary      string   `json:"summary,omitempty"`
	SourceType   string   `json:"source_type,omitempty"`
	Origin       string   `json:"origin,omitempty"`
	SourceTrust  string   `json:"source_trust,omitempty"`
	Sensitivity  string   `json:"sensitivity,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	WorkingScope string   `json:"working_scope"`
	IsPinned     bool     `json:"is_pinned,omitempty"`
}

// ChronicleIngestResponse は chronicle への投入レスポンス。
type ChronicleIngestResponse struct {
	Note ChronicleNote `json:"note"`
}

// ChronicleSearchRequest は chronicle 検索リクエスト。
type ChronicleSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

// ChronicleSearchResponse は chronicle 検索レスポンス。
type ChronicleSearchResponse struct {
	Notes []ChronicleNote `json:"notes"`
}

// ChronicleListByScopeRequest は scope 指定リストリクエスト。
type ChronicleListByScopeRequest struct {
	WorkingScope string `json:"working_scope"`
	Limit        int    `json:"limit,omitempty"`
}

// ChronicleListByScopeResponse は scope 指定リストレスポンス。
type ChronicleListByScopeResponse struct {
	Notes []ChronicleNote `json:"notes"`
}

// -------------------- Memopedia --------------------

// MemopediaNote は memopedia ストアのノート。
type MemopediaNote struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	LastAccessedAt string `json:"last_accessed_at"`
	AccessCount    int64  `json:"access_count"`
	SourceType     string `json:"source_type"`
	Origin         string `json:"origin"`
	SourceTrust    string `json:"source_trust"`
	Sensitivity    string `json:"sensitivity"`
	WorkingScope   string `json:"working_scope"`
	IsPinned       bool   `json:"is_pinned"`
}

// MemopediaIngestRequest は memopedia への投入リクエスト。
type MemopediaIngestRequest struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Summary      string   `json:"summary,omitempty"`
	SourceType   string   `json:"source_type,omitempty"`
	Origin       string   `json:"origin,omitempty"`
	SourceTrust  string   `json:"source_trust,omitempty"`
	Sensitivity  string   `json:"sensitivity,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	WorkingScope string   `json:"working_scope"`
	IsPinned     bool     `json:"is_pinned,omitempty"`
}

// MemopediaIngestResponse は memopedia への投入レスポンス。
type MemopediaIngestResponse struct {
	Note MemopediaNote `json:"note"`
}

// MemopediaSearchRequest は memopedia 検索リクエスト。
type MemopediaSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

// MemopediaSearchResponse は memopedia 検索レスポンス。
type MemopediaSearchResponse struct {
	Notes []MemopediaNote `json:"notes"`
}

// MemopediaListByScopeRequest は scope 指定リストリクエスト。
type MemopediaListByScopeRequest struct {
	WorkingScope string `json:"working_scope"`
	Limit        int    `json:"limit,omitempty"`
}

// MemopediaListByScopeResponse は scope 指定リストレスポンス。
type MemopediaListByScopeResponse struct {
	Notes []MemopediaNote `json:"notes"`
}

// MemopediaListPinnedRequest はピン留めノート一覧リクエスト。
type MemopediaListPinnedRequest struct {
	WorkingScope string `json:"working_scope,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

// MemopediaListPinnedResponse はピン留めノート一覧レスポンス。
type MemopediaListPinnedResponse struct {
	Notes []MemopediaNote `json:"notes"`
}

// PinRequest はピン留めリクエスト。
type PinRequest struct {
	ID string `json:"id"`
}

// PinResponse はピン留めレスポンス。
type PinResponse struct {
	Success bool `json:"success"`
}

// UnpinRequest はピン解除リクエスト。
type UnpinRequest struct {
	ID string `json:"id"`
}

// UnpinResponse はピン解除レスポンス。
type UnpinResponse struct {
	Success bool `json:"success"`
}

// -------------------- Archive --------------------

// ArchiveNote は archive ストアのノート。
type ArchiveNote struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	Body           string `json:"body"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	LastAccessedAt string `json:"last_accessed_at"`
	AccessCount    int64  `json:"access_count"`
	SourceType     string `json:"source_type"`
	Origin         string `json:"origin"`
	SourceTrust    string `json:"source_trust"`
	Sensitivity    string `json:"sensitivity"`
}

// ArchiveListRequest は archive 一覧リクエスト。
type ArchiveListRequest struct {
	Limit int `json:"limit,omitempty"`
}

// ArchiveListResponse は archive 一覧レスポンス。
type ArchiveListResponse struct {
	Notes []ArchiveNote `json:"notes"`
}

// ArchiveRestoreRequest は archive 復元リクエスト。
type ArchiveRestoreRequest struct {
	ID string `json:"id"`
}

// ArchiveRestoreResponse は archive 復元レスポンス。
type ArchiveRestoreResponse struct {
	Note Note `json:"note"`
}
