package api

import "context"

func (c *InProcClient) DocsIngest(ctx context.Context, req DocsIngestRequest) (DocsIngestResponse, *Error) {
	doc, chunkCount, err := c.Svc.DocsIngest(ctx, req)
	if err != nil {
		return DocsIngestResponse{}, mapError(err)
	}
	return DocsIngestResponse{DocID: doc.DocID, Version: doc.Version, ChunkCount: chunkCount, Status: "ingested"}, nil
}

func (c *InProcClient) DocsResolve(ctx context.Context, req DocsResolveRequest) (DocsResolveResponse, *Error) {
	required, recommended, err := c.Svc.DocsResolve(ctx, req)
	if err != nil {
		return DocsResolveResponse{}, mapError(err)
	}
	return DocsResolveResponse{Required: required, Recommended: recommended}, nil
}

func (c *InProcClient) ChunksGet(ctx context.Context, req ChunksGetRequest) (ChunksGetResponse, *Error) {
	docID, chunks, err := c.Svc.ChunksGet(ctx, req)
	if err != nil {
		return ChunksGetResponse{}, mapError(err)
	}
	return ChunksGetResponse{DocID: docID, Chunks: chunks}, nil
}

func (c *InProcClient) DocsSearch(ctx context.Context, req DocsSearchRequest) (DocsSearchResponse, *Error) {
	results, err := c.Svc.DocsSearch(ctx, req)
	if err != nil {
		return DocsSearchResponse{}, mapError(err)
	}
	return DocsSearchResponse{Results: results}, nil
}

func (c *InProcClient) ReadsAck(ctx context.Context, req ReadsAckRequest) (ReadsAckResponse, *Error) {
	receipt, err := c.Svc.ReadsAck(ctx, req)
	if err != nil {
		return ReadsAckResponse{}, mapError(err)
	}
	return ReadsAckResponse{Receipt: receipt}, nil
}

func (c *InProcClient) DocsStaleCheck(ctx context.Context, req DocsStaleCheckRequest) (DocsStaleCheckResponse, *Error) {
	stale, err := c.Svc.DocsStaleCheck(ctx, req)
	if err != nil {
		return DocsStaleCheckResponse{}, mapError(err)
	}
	return DocsStaleCheckResponse{TaskID: req.TaskID, Stale: stale}, nil
}

func (c *InProcClient) ContractsResolve(ctx context.Context, req ContractsResolveRequest) (ContractsResolveResponse, *Error) {
	required, acceptance, forbidden, done, dependencies, err := c.Svc.ContractsResolve(ctx, req)
	if err != nil {
		return ContractsResolveResponse{}, mapError(err)
	}
	return ContractsResolveResponse{
		Required:           required,
		AcceptanceCriteria: acceptance,
		ForbiddenPatterns:  forbidden,
		DefinitionOfDone:   done,
		Dependencies:       dependencies,
	}, nil
}
