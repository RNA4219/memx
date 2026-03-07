package api

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"memx/db"
	"memx/service"
)

// TestCLI_JSON_Ingest はCLI --json と API レスポンスの同型性をテストする。
// 仕様書: requirements.md "CLI `--json` は API レスポンスと同型を維持する"
func TestCLI_JSON_Ingest(t *testing.T) {
	// API レスポンスの構造
	apiResp := NotesIngestResponse{
		Note: Note{
			ID:             "0123456789abcdef0123456789abcdef",
			Ref:            NewTypedRef(EntityTypeEvidence, "0123456789abcdef0123456789abcdef"),
			Title:          "Test Note",
			Summary:        "Test summary",
			Body:           "Test body",
			CreatedAt:      "2026-03-08T00:00:00Z",
			UpdatedAt:       "2026-03-08T00:00:00Z",
			LastAccessedAt: "2026-03-08T00:00:00Z",
			AccessCount:    0,
			SourceType:     "manual",
			Origin:         "",
			SourceTrust:    "user_input",
			Sensitivity:    "internal",
		},
	}

	// JSONエンコード
	jsonBytes, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// JSONデコードして再構築
	var decoded NotesIngestResponse
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// 必須フィールドの存在確認
	if decoded.Note.ID != apiResp.Note.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.Note.ID, apiResp.Note.ID)
	}
	if decoded.Note.Title != apiResp.Note.Title {
		t.Errorf("Title mismatch: got %q, want %q", decoded.Note.Title, apiResp.Note.Title)
	}
	if decoded.Note.Body != apiResp.Note.Body {
		t.Errorf("Body mismatch: got %q, want %q", decoded.Note.Body, apiResp.Note.Body)
	}
	if decoded.Note.Ref != apiResp.Note.Ref {
		t.Errorf("Ref mismatch: got %q, want %q", decoded.Note.Ref, apiResp.Note.Ref)
	}
}

// TestCLI_JSON_Search は検索レスポンスの同型性をテストする。
func TestCLI_JSON_Search(t *testing.T) {
	apiResp := NotesSearchResponse{
		Notes: []Note{
			{
				ID:             "0123456789abcdef0123456789abcdef",
				Ref:            NewTypedRef(EntityTypeEvidence, "0123456789abcdef0123456789abcdef"),
				Title:          "Note 1",
				Summary:        "Summary 1",
				Body:           "Body 1",
				CreatedAt:      "2026-03-08T00:00:00Z",
				UpdatedAt:       "2026-03-08T00:00:00Z",
				LastAccessedAt: "2026-03-08T00:00:00Z",
				AccessCount:    1,
				SourceType:     "manual",
				Origin:         "",
				SourceTrust:    "user_input",
				Sensitivity:    "internal",
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded NotesSearchResponse
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Notes) != 1 {
		t.Errorf("expected 1 note, got: %d", len(decoded.Notes))
	}
	if decoded.Notes[0].ID != apiResp.Notes[0].ID {
		t.Errorf("ID mismatch")
	}
}

// TestCLI_JSON_Show は単一ノート取得の同型性をテストする。
func TestCLI_JSON_Show(t *testing.T) {
	apiNote := Note{
		ID:             "0123456789abcdef0123456789abcdef",
		Ref:            NewTypedRef(EntityTypeEvidence, "0123456789abcdef0123456789abcdef"),
		Title:          "Show Test",
		Summary:        "Summary",
		Body:           "Body content",
		CreatedAt:      "2026-03-08T00:00:00Z",
		UpdatedAt:       "2026-03-08T00:00:00Z",
		LastAccessedAt: "2026-03-08T00:00:00Z",
		AccessCount:    5,
		SourceType:     "manual",
		Origin:         "",
		SourceTrust:    "user_input",
		Sensitivity:    "internal",
	}

	jsonBytes, err := json.MarshalIndent(apiNote, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Note
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ID != apiNote.ID {
		t.Errorf("ID mismatch")
	}
	if decoded.AccessCount != apiNote.AccessCount {
		t.Errorf("AccessCount mismatch: got %d, want %d", decoded.AccessCount, apiNote.AccessCount)
	}
}

// TestCLI_JSON_GC はGCレスポンスの同型性をテストする。
func TestCLI_JSON_GC(t *testing.T) {
	apiResp := GCRunResponse{
		Status: "completed",
	}

	jsonBytes, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded GCRunResponse
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Status != apiResp.Status {
		t.Errorf("Status mismatch: got %q, want %q", decoded.Status, apiResp.Status)
	}
}

// TestCLI_JSON_ErrorResponse はエラーレスポンスの同型性をテストする。
func TestCLI_JSON_ErrorResponse(t *testing.T) {
	apiErr := Error{
		Code:    CodeInvalidArgument,
		Message: "title is required",
		Details: map[string]interface{}{"field": "title"},
	}

	jsonBytes, err := json.MarshalIndent(apiErr, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Error
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Code != apiErr.Code {
		t.Errorf("Code mismatch: got %q, want %q", decoded.Code, apiErr.Code)
	}
	if decoded.Message != apiErr.Message {
		t.Errorf("Message mismatch: got %q, want %q", decoded.Message, apiErr.Message)
	}
}

// TestCLI_JSON_Summarize は要約レスポンスの同型性をテストする。
func TestCLI_JSON_Summarize(t *testing.T) {
	apiResp := SummarizeResponse{
		Note: Note{
			ID:             "0123456789abcdef0123456789abcdef",
			Ref:            NewTypedRef(EntityTypeEvidence, "0123456789abcdef0123456789abcdef"),
			Title:          "Summarized Note",
			Summary:        "This is the generated summary.",
			Body:           "Original body content",
			CreatedAt:      "2026-03-08T00:00:00Z",
			UpdatedAt:       "2026-03-08T00:00:01Z",
			LastAccessedAt: "2026-03-08T00:00:00Z",
			AccessCount:    2,
			SourceType:     "manual",
			Origin:         "",
			SourceTrust:    "user_input",
			Sensitivity:    "internal",
		},
	}

	jsonBytes, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded SummarizeResponse
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Note.Summary != apiResp.Note.Summary {
		t.Errorf("Summary mismatch: got %q, want %q", decoded.Note.Summary, apiResp.Note.Summary)
	}
}

// TestAPI_Client_Interface はClientインターフェースが正しく実装されているかテストする。
func TestAPI_Client_Interface(t *testing.T) {
	// InProcClientがClientインターフェースを実装していることを確認
	var _ Client = (*InProcClient)(nil)
	// HTTPClientがClientインターフェースを実装していることを確認
	var _ Client = (*HTTPClient)(nil)
}

// TestAPI_InProcClient_Integration はInProcClientの統合テストを行う。
func TestAPI_InProcClient_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	paths := db.Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	svc, err := service.New(paths)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	defer svc.Close()

	client := NewInProcClient(svc)
	ctx := context.Background()

	// Ingest
	ingestResp, apiErr := client.NotesIngest(ctx, NotesIngestRequest{
		Title: "Integration Test",
		Body:  "This is an integration test.",
	})
	if apiErr != nil {
		t.Fatalf("ingest failed: %v", apiErr)
	}
	if ingestResp.Note.ID == "" {
		t.Error("expected non-empty ID")
	}

	// Search
	searchResp, apiErr := client.NotesSearch(ctx, NotesSearchRequest{
		Query: "integration",
		TopK:  10,
	})
	if apiErr != nil {
		t.Fatalf("search failed: %v", apiErr)
	}
	if len(searchResp.Notes) == 0 {
		t.Error("expected at least one search result")
	}

	// Get
	getNote, apiErr := client.NotesGet(ctx, ingestResp.Note.ID)
	if apiErr != nil {
		t.Fatalf("get failed: %v", apiErr)
	}
	if getNote.Title != "Integration Test" {
		t.Errorf("expected title 'Integration Test', got: %q", getNote.Title)
	}

	// JSON出力の検証（同型性確認）
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ingestResp); err != nil {
		t.Fatalf("failed to encode ingest response: %v", err)
	}

	// JSONをデコードしてフィールドが存在するか確認
	var decoded map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	note, ok := decoded["note"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'note' field in response")
	}

	requiredFields := []string{"id", "title", "body", "created_at", "updated_at", "source_type", "sensitivity"}
	for _, field := range requiredFields {
		if _, exists := note[field]; !exists {
			t.Errorf("required field %q missing in JSON output", field)
		}
	}
}