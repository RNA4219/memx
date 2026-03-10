package api

import "context"

func (c *HTTPClient) DocsIngest(ctx context.Context, req DocsIngestRequest) (DocsIngestResponse, *Error) {
	var out DocsIngestResponse
	if err := c.post(ctx, "/v1/docs:ingest", req, &out); err != nil {
		return DocsIngestResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) DocsResolve(ctx context.Context, req DocsResolveRequest) (DocsResolveResponse, *Error) {
	var out DocsResolveResponse
	if err := c.post(ctx, "/v1/docs:resolve", req, &out); err != nil {
		return DocsResolveResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) ChunksGet(ctx context.Context, req ChunksGetRequest) (ChunksGetResponse, *Error) {
	var out ChunksGetResponse
	if err := c.post(ctx, "/v1/chunks:get", req, &out); err != nil {
		return ChunksGetResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) DocsSearch(ctx context.Context, req DocsSearchRequest) (DocsSearchResponse, *Error) {
	var out DocsSearchResponse
	if err := c.post(ctx, "/v1/docs:search", req, &out); err != nil {
		return DocsSearchResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) ReadsAck(ctx context.Context, req ReadsAckRequest) (ReadsAckResponse, *Error) {
	var out ReadsAckResponse
	if err := c.post(ctx, "/v1/reads:ack", req, &out); err != nil {
		return ReadsAckResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) DocsStaleCheck(ctx context.Context, req DocsStaleCheckRequest) (DocsStaleCheckResponse, *Error) {
	var out DocsStaleCheckResponse
	if err := c.post(ctx, "/v1/docs:stale-check", req, &out); err != nil {
		return DocsStaleCheckResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) ContractsResolve(ctx context.Context, req ContractsResolveRequest) (ContractsResolveResponse, *Error) {
	var out ContractsResolveResponse
	if err := c.post(ctx, "/v1/contracts:resolve", req, &out); err != nil {
		return ContractsResolveResponse{}, err
	}
	return out, nil
}
