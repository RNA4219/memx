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
