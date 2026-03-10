package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTP_DocsEndpoints(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	ctx := context.Background()
	_, apiErr := server.InProc.DocsIngest(ctx, DocsIngestRequest{
		DocType:     "spec",
		Title:       "Resolver API Spec",
		Version:     "2026-03-10",
		FeatureKeys: []string{"resolver-api"},
		TaskIDs:     []string{"task:feature:local:resolver"},
		Body: `# Resolver API Spec

## Acceptance Criteria
- returns required docs`,
	})
	if apiErr != nil {
		t.Fatalf("DocsIngest: %s", apiErr.Message)
	}

	jsonBody, _ := json.Marshal(DocsResolveRequest{Feature: "resolver-api"})
	resolveReq := httptest.NewRequest(http.MethodPost, "/v1/docs:resolve", bytes.NewReader(jsonBody))
	resolveReq.Header.Set("Content-Type", "application/json")
	resolveRec := httptest.NewRecorder()
	server.Handler().ServeHTTP(resolveRec, resolveReq)
	if resolveRec.Code != http.StatusOK {
		t.Fatalf("resolve status=%d body=%s", resolveRec.Code, resolveRec.Body.String())
	}
	var resolveResp DocsResolveResponse
	if err := json.Unmarshal(resolveRec.Body.Bytes(), &resolveResp); err != nil {
		t.Fatalf("decode resolve response: %v", err)
	}
	if len(resolveResp.Required) != 1 {
		t.Fatalf("unexpected resolve response: %#v", resolveResp)
	}

	ackBody, _ := json.Marshal(ReadsAckRequest{TaskID: "task:feature:local:resolver", DocID: resolveResp.Required[0].DocID})
	ackReq := httptest.NewRequest(http.MethodPost, "/v1/reads:ack", bytes.NewReader(ackBody))
	ackReq.Header.Set("Content-Type", "application/json")
	ackRec := httptest.NewRecorder()
	server.Handler().ServeHTTP(ackRec, ackReq)
	if ackRec.Code != http.StatusOK {
		t.Fatalf("ack status=%d body=%s", ackRec.Code, ackRec.Body.String())
	}

	staleBody, _ := json.Marshal(DocsStaleCheckRequest{TaskID: "task:feature:local:resolver"})
	staleReq := httptest.NewRequest(http.MethodPost, "/v1/docs:stale-check", bytes.NewReader(staleBody))
	staleReq.Header.Set("Content-Type", "application/json")
	staleRec := httptest.NewRecorder()
	server.Handler().ServeHTTP(staleRec, staleReq)
	if staleRec.Code != http.StatusOK {
		t.Fatalf("stale status=%d body=%s", staleRec.Code, staleRec.Body.String())
	}
	var staleResp DocsStaleCheckResponse
	if err := json.Unmarshal(staleRec.Body.Bytes(), &staleResp); err != nil {
		t.Fatalf("decode stale response: %v", err)
	}
	if len(staleResp.Stale) != 0 {
		t.Fatalf("expected no stale docs, got %#v", staleResp.Stale)
	}
}
