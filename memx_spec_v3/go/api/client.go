package api

import "context"

// Client は CLI / Tool / Agent が利用する API クライアント。
// v1.3 では HTTP と in-proc の両方を提供する。
type Client interface {
	NotesIngest(ctx context.Context, req NotesIngestRequest) (NotesIngestResponse, *Error)
	NotesSearch(ctx context.Context, req NotesSearchRequest) (NotesSearchResponse, *Error)
	NotesGet(ctx context.Context, id string) (Note, *Error)
	GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error)
	// 要約機能
	Summarize(ctx context.Context, id string) (SummarizeResponse, *Error)
	SummarizeBatch(ctx context.Context, req SummarizeBatchRequest) (SummarizeBatchResponse, *Error)
}
