package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"memx/db"
	"memx/service"
)

// setupHTTPTestServer はテスト用のHTTPサーバーを作成する。
func setupHTTPTestServer(t *testing.T) (*HTTPServer, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal:   filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := service.New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	server := NewHTTPServer(svc)
	cleanup := func() { svc.Close() }

	return server, cleanup
}

// TestHTTP_Healthz は /healthz エンドポイントのテスト。
func TestHTTP_Healthz(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got: %q", rec.Body.String())
	}
}

// TestHTTP_NotesIngest_Success は正常なノート保存のテスト。
func TestHTTP_NotesIngest_Success(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	body := NotesIngestRequest{
		Title:       "Test Note",
		Body:        "This is a test note body.",
		SourceType:  "manual",
		SourceTrust: "user_input",
		Sensitivity: "internal",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp NotesIngestResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Note.ID == "" {
		t.Error("expected non-empty note ID")
	}
	if resp.Note.Title != "Test Note" {
		t.Errorf("expected title 'Test Note', got: %q", resp.Note.Title)
	}
	if resp.Note.SourceType != "manual" {
		t.Errorf("expected source_type 'manual', got: %q", resp.Note.SourceType)
	}
}

// TestHTTP_NotesIngest_ValidationErrors はバリデーションエラーのテスト。
func TestHTTP_NotesIngest_ValidationErrors(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		body       NotesIngestRequest
		wantStatus int
		wantCode   ErrorCode
	}{
		{
			name:       "empty title",
			body:       NotesIngestRequest{Body: "test"},
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidArgument,
		},
		{
			name:       "empty body",
			body:       NotesIngestRequest{Title: "test"},
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidArgument,
		},
		{
			name:       "invalid source_type",
			body:       NotesIngestRequest{Title: "t", Body: "b", SourceType: "invalid"},
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidArgument,
		},
		{
			name:       "invalid source_trust",
			body:       NotesIngestRequest{Title: "t", Body: "b", SourceTrust: "invalid"},
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidArgument,
		},
		{
			name:       "invalid sensitivity",
			body:       NotesIngestRequest{Title: "t", Body: "b", Sensitivity: "invalid"},
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got: %d", tt.wantStatus, rec.Code)
			}

			var errResp Error
			if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
				t.Fatalf("failed to parse error response: %v", err)
			}
			if errResp.Code != tt.wantCode {
				t.Errorf("expected error code %q, got: %q", tt.wantCode, errResp.Code)
			}
		})
	}
}

// TestHTTP_NotesIngest_PolicyDenied はポリシー拒否のテスト。
func TestHTTP_NotesIngest_PolicyDenied(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	body := NotesIngestRequest{
		Title:       "Secret Note",
		Body:        "This is secret.",
		Sensitivity: "secret",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got: %d", rec.Code)
	}

	var errResp Error
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Code != CodeGatekeepDeny {
		t.Errorf("expected error code %q, got: %q", CodeGatekeepDeny, errResp.Code)
	}
}

// TestHTTP_NotesIngest_ConfidentialSensitivity はconfidential機密度のテスト。
// 仕様書: AGENT_GUIDE.md - sensitivity に confidential が含まれる
func TestHTTP_NotesIngest_ConfidentialSensitivity(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	body := NotesIngestRequest{
		Title:       "Confidential Note",
		Body:        "This is confidential information.",
		Sensitivity: "confidential",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	// confidential は許可される（secret のみ拒否）
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp NotesIngestResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Note.Sensitivity != "confidential" {
		t.Errorf("expected sensitivity 'confidential', got: %q", resp.Note.Sensitivity)
	}
}

// TestHTTP_NotesIngest_AllValidSensitivities は全ての有効なsensitivity値をテストする。
func TestHTTP_NotesIngest_AllValidSensitivities(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	validSensitivities := []string{"public", "internal", "confidential"}

	for _, sens := range validSensitivities {
		t.Run(sens, func(t *testing.T) {
			body := NotesIngestRequest{
				Title:       "Test Note",
				Body:        "Test body.",
				Sensitivity: sens,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("sensitivity %q: expected status 200, got: %d", sens, rec.Code)
			}
		})
	}
}

// TestHTTP_NotesSearch_Success は正常な検索のテスト。
func TestHTTP_NotesSearch_Success(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	// まずノートを保存
	ingestBody := NotesIngestRequest{
		Title: "Search Test Note",
		Body:  "This note contains unique keyword xyz123.",
	}
	jsonIngest, _ := json.Marshal(ingestBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonIngest))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	// 検索実行
	searchBody := NotesSearchRequest{
		Query: "xyz123",
		TopK:  10,
	}
	jsonSearch, _ := json.Marshal(searchBody)

	req = httptest.NewRequest(http.MethodPost, "/v1/notes:search", bytes.NewReader(jsonSearch))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", rec.Code)
	}

	var resp NotesSearchResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Notes) == 0 {
		t.Error("expected at least one search result")
	}
}

// TestHTTP_NotesSearch_ValidationError は検索バリデーションエラーのテスト。
func TestHTTP_NotesSearch_ValidationError(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		body       NotesSearchRequest
		wantStatus int
	}{
		{
			name:       "empty query",
			body:       NotesSearchRequest{Query: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "query too long (over 1000 chars)",
			body:       NotesSearchRequest{Query: string(make([]byte, 1001))},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/v1/notes:search", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got: %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

// TestHTTP_NotesSearch_TopKLimit はtop_kパラメータの制限テスト。
// 仕様書: requirements-api.md - top_k 上限は100
func TestHTTP_NotesSearch_TopKLimit(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	// テスト用に複数のノートを作成
	for i := 0; i < 5; i++ {
		ingestBody := NotesIngestRequest{
			Title: "Test Note",
			Body:  "unique_keyword_xyz test note",
		}
		jsonIngest, _ := json.Marshal(ingestBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonIngest))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.Handler().ServeHTTP(rec, req)
	}

	tests := []struct {
		name        string
		topK        int
		wantStatus  int
		expectCount int // 期待される結果数（上限チェック）
	}{
		{
			name:        "top_k default (0 should become 20)",
			topK:        0,
			wantStatus:  http.StatusOK,
			expectCount: 5, // 5件作成したので全て返る
		},
		{
			name:        "top_k within limit",
			topK:        50,
			wantStatus:  http.StatusOK,
			expectCount: 5,
		},
		{
			name:        "top_k at limit (100)",
			topK:        100,
			wantStatus:  http.StatusOK,
			expectCount: 5,
		},
		{
			name:        "top_k over limit (150 should be capped to 100)",
			topK:        150,
			wantStatus:  http.StatusOK,
			expectCount: 5, // 実際のデータは5件しかない
		},
		{
			name:        "negative top_k should become default (20)",
			topK:        -1,
			wantStatus:  http.StatusOK,
			expectCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchBody := NotesSearchRequest{
				Query: "unique_keyword_xyz",
				TopK:  tt.topK,
			}
			jsonSearch, _ := json.Marshal(searchBody)

			req := httptest.NewRequest(http.MethodPost, "/v1/notes:search", bytes.NewReader(jsonSearch))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got: %d", tt.wantStatus, rec.Code)
				return
			}

			var resp NotesSearchResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if len(resp.Notes) > 100 {
				t.Errorf("top_k limit exceeded: got %d results, max is 100", len(resp.Notes))
			}
		})
	}
}

// TestHTTP_NotesGet_Success は正常なノート取得のテスト。
func TestHTTP_NotesGet_Success(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	// まずノートを保存
	ingestBody := NotesIngestRequest{
		Title: "Get Test Note",
		Body:  "This is a test note for GET endpoint.",
	}
	jsonIngest, _ := json.Marshal(ingestBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonIngest))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	var ingestResp NotesIngestResponse
	json.Unmarshal(rec.Body.Bytes(), &ingestResp)
	noteID := ingestResp.Note.ID

	// GET実行
	req = httptest.NewRequest(http.MethodGet, "/v1/notes/"+noteID, nil)
	rec = httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", rec.Code)
	}

	var note Note
	if err := json.Unmarshal(rec.Body.Bytes(), &note); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if note.ID != noteID {
		t.Errorf("expected ID %q, got: %q", noteID, note.ID)
	}
	if note.Title != "Get Test Note" {
		t.Errorf("expected title 'Get Test Note', got: %q", note.Title)
	}
}

// TestHTTP_NotesGet_NotFound は存在しないノート取得のテスト。
func TestHTTP_NotesGet_NotFound(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/v1/notes/00000000000000000000000000000000", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", rec.Code)
	}

	var errResp Error
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Code != CodeNotFound {
		t.Errorf("expected error code %q, got: %q", CodeNotFound, errResp.Code)
	}
}

// TestHTTP_NotesGet_ValidationError はGETバリデーションエラーのテスト。
func TestHTTP_NotesGet_ValidationError(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "empty id",
			id:         "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id format",
			id:         "invalid-id",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/notes/"+tt.id, nil)
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got: %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

// TestHTTP_GCRun_DryRun はGC dry-runのテスト。
func TestHTTP_GCRun_DryRun(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	body := GCRunRequest{
		Target:  "short",
		Options: GCOptions{DryRun: true},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/gc:run", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", rec.Code)
	}

	var resp GCRunResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// dry-runの場合、StatusにJSON形式の結果が入る
	if resp.Status == "" {
		t.Error("expected non-empty status")
	}
}

// TestHTTP_GCRun_FeatureDisabled はGC機能無効時のテスト。
// 注: API層ではEnabled=trueが固定で設定される。Feature flagのチェックはCLI層で行われる。
func TestHTTP_GCRun_FeatureDisabled(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	// API層ではEnabled=trueが固定のため、dry_run=falseでも正常に実行される
	// Feature flagのチェックはCLI層で行われる
	body := GCRunRequest{
		Target:  "short",
		Options: GCOptions{DryRun: false},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/gc:run", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	// GC実行が試みられる（soft_limit未満の場合はskippedになる）
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", rec.Code)
	}
}

// TestHTTP_MethodNotAllowed は許可されないメソッドのテスト。
func TestHTTP_MethodNotAllowed(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"GET ingest", http.MethodGet, "/v1/notes:ingest"},
		{"GET search", http.MethodGet, "/v1/notes:search"},
		{"POST get", http.MethodPost, "/v1/notes/test-id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status 405, got: %d", rec.Code)
			}
		})
	}
}

type staticMiniLLM struct {
	summary string
}

func (m *staticMiniLLM) TagAndScore(ctx context.Context, noteBody string) (db.TagsAndScores, error) {
	return db.TagsAndScores{}, nil
}

func (m *staticMiniLLM) Summarize(ctx context.Context, title, body string) (db.SummarizeResult, error) {
	return db.SummarizeResult{Summary: m.summary}, nil
}

func TestHTTP_NotesIngest_NoLLMSkipsAutoSummary(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	server.InProc.Svc.SetMiniLLM(&staticMiniLLM{summary: "should not be used"})

	body := NotesIngestRequest{
		Title:       "No LLM Note",
		Body:        "This should skip auto-summary.",
		SourceType:  "manual",
		SourceTrust: "user_input",
		Sensitivity: "internal",
		NoLLM:       true,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:ingest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp NotesIngestResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Note.Summary != "" {
		t.Fatalf("expected empty summary when no_llm=true, got: %q", resp.Note.Summary)
	}
}

func TestHTTP_NotesRecall_Success(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	ctx := context.Background()
	shortResp, apiErr := server.InProc.NotesIngest(ctx, NotesIngestRequest{
		Title:       "Recall Short",
		Body:        "shared recall marker short",
		Sensitivity: "internal",
		NoLLM:       true,
	})
	if apiErr != nil {
		t.Fatalf("NotesIngest: %s", apiErr.Message)
	}
	journalResp, apiErr := server.InProc.JournalIngest(ctx, JournalIngestRequest{
		Title:        "Recall Journal",
		Body:         "shared recall marker journal",
		WorkingScope: "project:memx",
		Sensitivity:  "internal",
		NoLLM:        true,
	})
	if apiErr != nil {
		t.Fatalf("JournalIngest: %s", apiErr.Message)
	}
	knowledgeResp, apiErr := server.InProc.KnowledgeIngest(ctx, KnowledgeIngestRequest{
		Title:        "Recall Knowledge",
		Body:         "shared recall marker knowledge",
		WorkingScope: "glossary",
		Sensitivity:  "internal",
		NoLLM:        true,
	})
	if apiErr != nil {
		t.Fatalf("KnowledgeIngest: %s", apiErr.Message)
	}

	jsonBody, _ := json.Marshal(RecallRequest{
		Query:       "shared recall marker",
		Stores:      []string{"short", "journal", "knowledge"},
		FallbackFTS: true,
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/notes:recall", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp RecallResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	gotIDs := map[string]struct{}{}
	for _, result := range resp.Results {
		gotIDs[result.Anchor.ID] = struct{}{}
	}
	for _, wantID := range []string{shortResp.Note.ID, journalResp.Note.ID, knowledgeResp.Note.ID} {
		if _, ok := gotIDs[wantID]; !ok {
			t.Fatalf("expected %s in recall response: %#v", wantID, gotIDs)
		}
	}
}

func TestHTTP_BuildBundle_Success(t *testing.T) {
	server, cleanup := setupHTTPTestServer(t)
	defer cleanup()

	ctx := context.Background()
	journalResp, apiErr := server.InProc.JournalIngest(ctx, JournalIngestRequest{
		Title:        "Bundle Journal",
		Body:         "bundle journal body",
		Summary:      "bundle journal summary",
		WorkingScope: "project:memx",
		Sensitivity:  "internal",
		NoLLM:        true,
	})
	if apiErr != nil {
		t.Fatalf("JournalIngest: %s", apiErr.Message)
	}
	knowledgeResp, apiErr := server.InProc.KnowledgeIngest(ctx, KnowledgeIngestRequest{
		Title:        "Bundle Knowledge",
		Body:         "bundle knowledge body",
		Summary:      "bundle knowledge summary",
		WorkingScope: "glossary",
		Sensitivity:  "internal",
		NoLLM:        true,
	})
	if apiErr != nil {
		t.Fatalf("KnowledgeIngest: %s", apiErr.Message)
	}

	jsonBody, _ := json.Marshal(BuildBundleRequest{
		Purpose:    "verify",
		SourceRefs: []TypedRef{journalResp.Note.Ref, knowledgeResp.Note.Ref},
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/bundle:build", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp BuildBundleResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Bundle.ID == "" {
		t.Fatal("expected non-empty bundle id")
	}
	if resp.Bundle.Purpose != "verify" {
		t.Fatalf("expected purpose verify, got %q", resp.Bundle.Purpose)
	}
	if len(resp.Bundle.SourceRefs) != 2 {
		t.Fatalf("expected 2 source refs, got %d", len(resp.Bundle.SourceRefs))
	}
	if len(resp.Bundle.EvidenceRefs) != 1 || resp.Bundle.EvidenceRefs[0].ID != journalResp.Note.ID {
		t.Fatalf("expected journal note evidence ref, got %#v", resp.Bundle.EvidenceRefs)
	}
	if resp.Bundle.RawIncluded {
		t.Fatal("expected raw_included=false by default")
	}
	if resp.Bundle.Diagnostics.PartialBundle {
		t.Fatal("expected complete bundle")
	}
}
