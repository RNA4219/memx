package api

import "context"

// Client は CLI / Tool / Agent が利用する API クライアント。
// v1.4 では HTTP と in-proc の両方を提供する。
type Client interface {
	// Short store
	NotesIngest(ctx context.Context, req NotesIngestRequest) (NotesIngestResponse, *Error)
	NotesSearch(ctx context.Context, req NotesSearchRequest) (NotesSearchResponse, *Error)
	NotesGet(ctx context.Context, id string) (Note, *Error)
	GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error)
	// 要約機能
	Summarize(ctx context.Context, id string) (SummarizeResponse, *Error)
	SummarizeBatch(ctx context.Context, req SummarizeBatchRequest) (SummarizeBatchResponse, *Error)

	// Journal store
	JournalIngest(ctx context.Context, req JournalIngestRequest) (JournalIngestResponse, *Error)
	JournalSearch(ctx context.Context, req JournalSearchRequest) (JournalSearchResponse, *Error)
	JournalGet(ctx context.Context, id string) (JournalNote, *Error)
	JournalListByScope(ctx context.Context, req JournalListByScopeRequest) (JournalListByScopeResponse, *Error)

	// Knowledge store
	KnowledgeIngest(ctx context.Context, req KnowledgeIngestRequest) (KnowledgeIngestResponse, *Error)
	KnowledgeSearch(ctx context.Context, req KnowledgeSearchRequest) (KnowledgeSearchResponse, *Error)
	KnowledgeGet(ctx context.Context, id string) (KnowledgeNote, *Error)
	KnowledgeListByScope(ctx context.Context, req KnowledgeListByScopeRequest) (KnowledgeListByScopeResponse, *Error)
	KnowledgeListPinned(ctx context.Context, req KnowledgeListPinnedRequest) (KnowledgeListPinnedResponse, *Error)
	KnowledgePin(ctx context.Context, id string) (PinResponse, *Error)
	KnowledgeUnpin(ctx context.Context, id string) (UnpinResponse, *Error)

	// Archive store
	ArchiveGet(ctx context.Context, id string) (ArchiveNote, *Error)
	ArchiveList(ctx context.Context, req ArchiveListRequest) (ArchiveListResponse, *Error)
	ArchiveRestore(ctx context.Context, id string) (ArchiveRestoreResponse, *Error)

	// Recall
	Recall(ctx context.Context, req RecallRequest) (RecallResponse, *Error)

	// Resolver API (P4)
	ResolveRef(ctx context.Context, req ResolveRefRequest) (ResolveRefResponse, *Error)
	ResolveMany(ctx context.Context, req ResolveManyRequest) (ResolveManyResponse, *Error)
	LoadSummary(ctx context.Context, req LoadSummaryRequest) (LoadSummaryResponse, *Error)
	LoadSelectedRaw(ctx context.Context, req LoadSelectedRawRequest) (LoadSelectedRawResponse, *Error)
	BuildBundle(ctx context.Context, req BuildBundleRequest) (BuildBundleResponse, *Error)

	// Docs resolver
	DocsIngest(ctx context.Context, req DocsIngestRequest) (DocsIngestResponse, *Error)
	DocsResolve(ctx context.Context, req DocsResolveRequest) (DocsResolveResponse, *Error)
	ChunksGet(ctx context.Context, req ChunksGetRequest) (ChunksGetResponse, *Error)
	DocsSearch(ctx context.Context, req DocsSearchRequest) (DocsSearchResponse, *Error)
	ReadsAck(ctx context.Context, req ReadsAckRequest) (ReadsAckResponse, *Error)
	DocsStaleCheck(ctx context.Context, req DocsStaleCheckRequest) (DocsStaleCheckResponse, *Error)
	ContractsResolve(ctx context.Context, req ContractsResolveRequest) (ContractsResolveResponse, *Error)
}
