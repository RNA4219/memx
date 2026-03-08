package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClient は /v1/* を叩くクライアント。
// BaseURL 例: http://127.0.0.1:7766

type HTTPClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *HTTPClient) NotesIngest(ctx context.Context, req NotesIngestRequest) (NotesIngestResponse, *Error) {
	var out NotesIngestResponse
	if err := c.post(ctx, "/v1/notes:ingest", req, &out); err != nil {
		return NotesIngestResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) NotesSearch(ctx context.Context, req NotesSearchRequest) (NotesSearchResponse, *Error) {
	var out NotesSearchResponse
	if err := c.post(ctx, "/v1/notes:search", req, &out); err != nil {
		return NotesSearchResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) NotesGet(ctx context.Context, id string) (Note, *Error) {
	var out Note
	if err := c.get(ctx, "/v1/notes/"+id, &out); err != nil {
		return Note{}, err
	}
	return out, nil
}

func (c *HTTPClient) GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error) {
	var out GCRunResponse
	if err := c.post(ctx, "/v1/gc:run", req, &out); err != nil {
		return GCRunResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) Summarize(ctx context.Context, id string) (SummarizeResponse, *Error) {
	var out SummarizeResponse
	if err := c.post(ctx, "/v1/notes:summarize", SummarizeRequest{ID: id}, &out); err != nil {
		return SummarizeResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) SummarizeBatch(ctx context.Context, req SummarizeBatchRequest) (SummarizeBatchResponse, *Error) {
	var out SummarizeBatchResponse
	if err := c.post(ctx, "/v1/notes:summarize-batch", req, &out); err != nil {
		return SummarizeBatchResponse{}, err
	}
	return out, nil
}

// -------------------- Journal --------------------

func (c *HTTPClient) JournalIngest(ctx context.Context, req JournalIngestRequest) (JournalIngestResponse, *Error) {
	var out JournalIngestResponse
	if err := c.post(ctx, "/v1/journal:ingest", req, &out); err != nil {
		return JournalIngestResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) JournalSearch(ctx context.Context, req JournalSearchRequest) (JournalSearchResponse, *Error) {
	var out JournalSearchResponse
	if err := c.post(ctx, "/v1/journal:search", req, &out); err != nil {
		return JournalSearchResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) JournalGet(ctx context.Context, id string) (JournalNote, *Error) {
	var out JournalNote
	if err := c.get(ctx, "/v1/journal/"+id, &out); err != nil {
		return JournalNote{}, err
	}
	return out, nil
}

func (c *HTTPClient) JournalListByScope(ctx context.Context, req JournalListByScopeRequest) (JournalListByScopeResponse, *Error) {
	var out JournalListByScopeResponse
	if err := c.post(ctx, "/v1/journal:list-by-scope", req, &out); err != nil {
		return JournalListByScopeResponse{}, err
	}
	return out, nil
}

// -------------------- Knowledge --------------------

func (c *HTTPClient) KnowledgeIngest(ctx context.Context, req KnowledgeIngestRequest) (KnowledgeIngestResponse, *Error) {
	var out KnowledgeIngestResponse
	if err := c.post(ctx, "/v1/knowledge:ingest", req, &out); err != nil {
		return KnowledgeIngestResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) KnowledgeSearch(ctx context.Context, req KnowledgeSearchRequest) (KnowledgeSearchResponse, *Error) {
	var out KnowledgeSearchResponse
	if err := c.post(ctx, "/v1/knowledge:search", req, &out); err != nil {
		return KnowledgeSearchResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) KnowledgeGet(ctx context.Context, id string) (KnowledgeNote, *Error) {
	var out KnowledgeNote
	if err := c.get(ctx, "/v1/knowledge/"+id, &out); err != nil {
		return KnowledgeNote{}, err
	}
	return out, nil
}

func (c *HTTPClient) KnowledgeListByScope(ctx context.Context, req KnowledgeListByScopeRequest) (KnowledgeListByScopeResponse, *Error) {
	var out KnowledgeListByScopeResponse
	if err := c.post(ctx, "/v1/knowledge:list-by-scope", req, &out); err != nil {
		return KnowledgeListByScopeResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) KnowledgeListPinned(ctx context.Context, req KnowledgeListPinnedRequest) (KnowledgeListPinnedResponse, *Error) {
	var out KnowledgeListPinnedResponse
	if err := c.post(ctx, "/v1/knowledge:list-pinned", req, &out); err != nil {
		return KnowledgeListPinnedResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) KnowledgePin(ctx context.Context, id string) (PinResponse, *Error) {
	var out PinResponse
	if err := c.post(ctx, "/v1/knowledge/"+id+":pin", nil, &out); err != nil {
		return PinResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) KnowledgeUnpin(ctx context.Context, id string) (UnpinResponse, *Error) {
	var out UnpinResponse
	if err := c.post(ctx, "/v1/knowledge/"+id+":unpin", nil, &out); err != nil {
		return UnpinResponse{}, err
	}
	return out, nil
}

// -------------------- Archive --------------------

func (c *HTTPClient) ArchiveGet(ctx context.Context, id string) (ArchiveNote, *Error) {
	var out ArchiveNote
	if err := c.get(ctx, "/v1/archive/"+id, &out); err != nil {
		return ArchiveNote{}, err
	}
	return out, nil
}

func (c *HTTPClient) ArchiveList(ctx context.Context, req ArchiveListRequest) (ArchiveListResponse, *Error) {
	var out ArchiveListResponse
	if err := c.get(ctx, "/v1/archive", &out); err != nil {
		return ArchiveListResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) ArchiveRestore(ctx context.Context, id string) (ArchiveRestoreResponse, *Error) {
	var out ArchiveRestoreResponse
	if err := c.post(ctx, "/v1/archive/"+id+":restore", nil, &out); err != nil {
		return ArchiveRestoreResponse{}, err
	}
	return out, nil
}

// -------------------- Recall --------------------

func (c *HTTPClient) Recall(ctx context.Context, req RecallRequest) (RecallResponse, *Error) {
	var out RecallResponse
	if err := c.post(ctx, "/v1/notes:recall", req, &out); err != nil {
		return RecallResponse{}, err
	}
	return out, nil
}

// -------------------- Resolver API (P4) --------------------

func (c *HTTPClient) ResolveRef(ctx context.Context, req ResolveRefRequest) (ResolveRefResponse, *Error) {
	var out ResolveRefResponse
	if err := c.post(ctx, "/v1/resolve", req, &out); err != nil {
		return ResolveRefResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) ResolveMany(ctx context.Context, req ResolveManyRequest) (ResolveManyResponse, *Error) {
	var out ResolveManyResponse
	if err := c.post(ctx, "/v1/resolve-many", req, &out); err != nil {
		return ResolveManyResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) LoadSummary(ctx context.Context, req LoadSummaryRequest) (LoadSummaryResponse, *Error) {
	var out LoadSummaryResponse
	if err := c.post(ctx, "/v1/summary", req, &out); err != nil {
		return LoadSummaryResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) LoadSelectedRaw(ctx context.Context, req LoadSelectedRawRequest) (LoadSelectedRawResponse, *Error) {
	var out LoadSelectedRawResponse
	if err := c.post(ctx, "/v1/raw", req, &out); err != nil {
		return LoadSelectedRawResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) BuildBundle(ctx context.Context, req BuildBundleRequest) (BuildBundleResponse, *Error) {
	var out BuildBundleResponse
	if err := c.post(ctx, "/v1/bundle:build", req, &out); err != nil {
		return BuildBundleResponse{}, err
	}
	return out, nil
}

func (c *HTTPClient) post(ctx context.Context, path string, in interface{}, out interface{}) *Error {
	b, _ := json.Marshal(in)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return &Error{Code: CodeInternal, Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return &Error{Code: CodeInternal, Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return decodeAPIError(resp)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return &Error{Code: CodeInternal, Message: "failed to decode response"}
	}
	return nil
}

func (c *HTTPClient) get(ctx context.Context, path string, out interface{}) *Error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return &Error{Code: CodeInternal, Message: err.Error()}
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return &Error{Code: CodeInternal, Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return decodeAPIError(resp)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return &Error{Code: CodeInternal, Message: "failed to decode response"}
	}
	return nil
}

func decodeAPIError(resp *http.Response) *Error {
	b, _ := io.ReadAll(resp.Body)
	var e Error
	if err := json.Unmarshal(b, &e); err == nil && e.Code != "" {
		return &e
	}
	return &Error{Code: CodeInternal, Message: fmt.Sprintf("http %d: %s", resp.StatusCode, string(b))}
}
