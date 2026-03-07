package api

import "context"

// Client は CLI / Tool / Agent が利用する API クライアント。
// v1.3 では HTTP と in-proc の両方を提供する。
type Client interface {
	// Short store
	NotesIngest(ctx context.Context, req NotesIngestRequest) (NotesIngestResponse, *Error)
	NotesSearch(ctx context.Context, req NotesSearchRequest) (NotesSearchResponse, *Error)
	NotesGet(ctx context.Context, id string) (Note, *Error)
	GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error)
	// 要約機能
	Summarize(ctx context.Context, id string) (SummarizeResponse, *Error)
	SummarizeBatch(ctx context.Context, req SummarizeBatchRequest) (SummarizeBatchResponse, *Error)

	// Chronicle store
	ChronicleIngest(ctx context.Context, req ChronicleIngestRequest) (ChronicleIngestResponse, *Error)
	ChronicleSearch(ctx context.Context, req ChronicleSearchRequest) (ChronicleSearchResponse, *Error)
	ChronicleGet(ctx context.Context, id string) (ChronicleNote, *Error)
	ChronicleListByScope(ctx context.Context, req ChronicleListByScopeRequest) (ChronicleListByScopeResponse, *Error)

	// Memopedia store
	MemopediaIngest(ctx context.Context, req MemopediaIngestRequest) (MemopediaIngestResponse, *Error)
	MemopediaSearch(ctx context.Context, req MemopediaSearchRequest) (MemopediaSearchResponse, *Error)
	MemopediaGet(ctx context.Context, id string) (MemopediaNote, *Error)
	MemopediaListByScope(ctx context.Context, req MemopediaListByScopeRequest) (MemopediaListByScopeResponse, *Error)
	MemopediaListPinned(ctx context.Context, req MemopediaListPinnedRequest) (MemopediaListPinnedResponse, *Error)
	MemopediaPin(ctx context.Context, id string) (PinResponse, *Error)
	MemopediaUnpin(ctx context.Context, id string) (UnpinResponse, *Error)

	// Archive store
	ArchiveGet(ctx context.Context, id string) (ArchiveNote, *Error)
	ArchiveList(ctx context.Context, req ArchiveListRequest) (ArchiveListResponse, *Error)
	ArchiveRestore(ctx context.Context, id string) (ArchiveRestoreResponse, *Error)
}
