package api

import "memx/service"

type ResolverDocument = service.ResolverDocument

type ResolverChunk = service.ResolverChunk

type ResolveEntry = service.ResolveEntry

type ResolverReadReceipt = service.ResolverReadReceipt

type ResolverStaleReason = service.ResolverStaleReason

type ChunkingOptions = service.ChunkingOptions

type DocsIngestRequest = service.DocsIngestRequest

type DocsResolveRequest = service.DocsResolveRequest

type ChunksGetRequest = service.ChunksGetRequest

type DocsSearchRequest = service.DocsSearchRequest

type ReadsAckRequest = service.ReadsAckRequest

type DocsStaleCheckRequest = service.DocsStaleCheckRequest

type ContractsResolveRequest = service.ContractsResolveRequest

type DocsIngestResponse struct {
	DocID      string `json:"doc_id"`
	Version    string `json:"version"`
	ChunkCount int    `json:"chunk_count"`
	Status     string `json:"status"`
}

type DocsResolveResponse struct {
	Required    []ResolveEntry `json:"required"`
	Recommended []ResolveEntry `json:"recommended"`
}

type ChunksGetResponse struct {
	DocID  string          `json:"doc_id"`
	Chunks []ResolverChunk `json:"chunks"`
}

type DocsSearchResponse struct {
	Results []ResolveEntry `json:"results"`
}

type ReadsAckResponse struct {
	Receipt ResolverReadReceipt `json:"receipt"`
}

type DocsStaleCheckResponse struct {
	TaskID string                `json:"task_id"`
	Stale  []ResolverStaleReason `json:"stale"`
}

type ContractsResolveResponse struct {
	Required           []ResolveEntry `json:"required"`
	AcceptanceCriteria []string       `json:"acceptance_criteria"`
	ForbiddenPatterns  []string       `json:"forbidden_patterns"`
	DefinitionOfDone   []string       `json:"definition_of_done"`
	Dependencies       []string       `json:"dependencies"`
}
