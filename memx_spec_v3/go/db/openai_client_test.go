package db

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoadOpenAIConfigFromEnv(t *testing.T) {
	t.Setenv("MEMX_LLM_PROVIDER", "openai")
	t.Setenv("OPENAI_API_KEY", "sk-test")
	t.Setenv("MEMX_OPENAI_MODEL", "gpt-5-mini")
	t.Setenv("MEMX_OPENAI_REFLECT_MODEL", "gpt-5")
	t.Setenv("MEMX_OPENAI_BASE_URL", "https://example.test/v1")
	t.Setenv("MEMX_OPENAI_TIMEOUT_SECONDS", "12")
	t.Setenv("OPENAI_PROJECT", "proj_123")
	t.Setenv("OPENAI_ORGANIZATION", "org_123")

	cfg, ok, err := LoadOpenAIConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadOpenAIConfigFromEnv: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if cfg.APIKey != "sk-test" {
		t.Fatalf("unexpected api key: %q", cfg.APIKey)
	}
	if cfg.MiniModel != "gpt-5-mini" {
		t.Fatalf("unexpected mini model: %q", cfg.MiniModel)
	}
	if cfg.ReflectModel != "gpt-5" {
		t.Fatalf("unexpected reflect model: %q", cfg.ReflectModel)
	}
	if cfg.BaseURL != "https://example.test/v1" {
		t.Fatalf("unexpected base url: %q", cfg.BaseURL)
	}
	if cfg.Timeout != 12*time.Second {
		t.Fatalf("unexpected timeout: %v", cfg.Timeout)
	}
	if cfg.Project != "proj_123" {
		t.Fatalf("unexpected project: %q", cfg.Project)
	}
	if cfg.Organization != "org_123" {
		t.Fatalf("unexpected organization: %q", cfg.Organization)
	}
}

func TestLoadAlibabaConfigFromEnv(t *testing.T) {
	t.Setenv("MEMX_LLM_PROVIDER", "alibaba")
	t.Setenv("DASHSCOPE_API_KEY", "dash-test")
	t.Setenv("MEMX_ALIBABA_MODEL", "qwen3-max")
	t.Setenv("MEMX_ALIBABA_REFLECT_MODEL", "qwen-max")
	t.Setenv("MEMX_ALIBABA_REGION", "beijing")
	t.Setenv("MEMX_ALIBABA_TIMEOUT_SECONDS", "15")

	cfg, ok, err := LoadOpenAIConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadOpenAIConfigFromEnv: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if cfg.APIKey != "dash-test" {
		t.Fatalf("unexpected api key: %q", cfg.APIKey)
	}
	if cfg.MiniModel != "qwen3-max" {
		t.Fatalf("unexpected mini model: %q", cfg.MiniModel)
	}
	if cfg.ReflectModel != "qwen-max" {
		t.Fatalf("unexpected reflect model: %q", cfg.ReflectModel)
	}
	if cfg.BaseURL != defaultAlibabaBeijingBaseURL {
		t.Fatalf("unexpected base url: %q", cfg.BaseURL)
	}
	if cfg.Timeout != 15*time.Second {
		t.Fatalf("unexpected timeout: %v", cfg.Timeout)
	}
	if !cfg.UseChatCompletions {
		t.Fatal("expected chat completions mode for Alibaba config")
	}
}

func TestLoadAlibabaConfigEnablesInlineInstructions(t *testing.T) {
	t.Setenv("MEMX_LLM_PROVIDER", "alibaba")
	t.Setenv("DASHSCOPE_API_KEY", "dash-test")

	cfg, ok, err := LoadOpenAIConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadOpenAIConfigFromEnv: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if !cfg.InlineInstructions {
		t.Fatal("expected inline instructions for Alibaba config")
	}
}
func TestOpenAIClientSummarize(t *testing.T) {
	type requestPayload struct {
		Model           string `json:"model"`
		Input           string `json:"input"`
		Instructions    string `json:"instructions"`
		MaxOutputTokens int    `json:"max_output_tokens"`
		Store           bool   `json:"store"`
		Reasoning       struct {
			Effort string `json:"effort"`
		} `json:"reasoning"`
		Text struct {
			Verbosity string `json:"verbosity"`
			Format    struct {
				Type string `json:"type"`
			} `json:"format"`
		} `json:"text"`
	}

	var captured requestPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test" {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"output":[{"type":"message","content":[{"type":"output_text","text":"要約テキストです。"}]}]}`))
	}))
	defer server.Close()

	client, err := NewOpenAIClient(OpenAIConfig{
		APIKey:    "sk-test",
		BaseURL:   server.URL + "/v1",
		MiniModel: "gpt-5-mini",
	})
	if err != nil {
		t.Fatalf("NewOpenAIClient: %v", err)
	}

	result, err := client.Summarize(context.Background(), "title", "body")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}
	if result.Summary != "要約テキストです。" {
		t.Fatalf("unexpected summary: %q", result.Summary)
	}
	if captured.Model != "gpt-5-mini" {
		t.Fatalf("unexpected model: %q", captured.Model)
	}
	if captured.MaxOutputTokens != defaultSummaryMaxTokens {
		t.Fatalf("unexpected max_output_tokens: %d", captured.MaxOutputTokens)
	}
	if captured.Store {
		t.Fatal("expected store=false")
	}
	if captured.Reasoning.Effort != "minimal" {
		t.Fatalf("unexpected reasoning effort: %q", captured.Reasoning.Effort)
	}
	if captured.Text.Verbosity != "low" || captured.Text.Format.Type != "text" {
		t.Fatalf("unexpected text config: %+v", captured.Text)
	}
}

func TestOpenAIClientSummarizeClusterUsesReflectModel(t *testing.T) {
	var model string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Model string `json:"model"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		model = payload.Model
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"output_text":"cluster summary"}`))
	}))
	defer server.Close()

	client, err := NewOpenAIClient(OpenAIConfig{
		APIKey:       "sk-test",
		BaseURL:      server.URL + "/v1",
		MiniModel:    "gpt-5-mini",
		ReflectModel: "gpt-5",
	})
	if err != nil {
		t.Fatalf("NewOpenAIClient: %v", err)
	}

	summary, err := client.SummarizeCluster(context.Background(), ClusterInput{
		NoteIDs: []string{"n1", "n2"},
		Body:    "body1\nbody2",
	})
	if err != nil {
		t.Fatalf("SummarizeCluster: %v", err)
	}
	if summary != "cluster summary" {
		t.Fatalf("unexpected cluster summary: %q", summary)
	}
	if model != "gpt-5" {
		t.Fatalf("unexpected reflect model: %q", model)
	}
}

func TestOpenAIClientErrorIncludesAPIMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"error":{"message":"bad api key"}}`))
	}))
	defer server.Close()

	client, err := NewOpenAIClient(OpenAIConfig{
		APIKey:  "sk-test",
		BaseURL: server.URL + "/v1",
	})
	if err != nil {
		t.Fatalf("NewOpenAIClient: %v", err)
	}

	_, err = client.Summarize(context.Background(), "title", "body")
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "openai responses api: bad api key" {
		t.Fatalf("unexpected error: %s", got)
	}
}

func TestOpenAIClientAlibabaUsesChatCompletions(t *testing.T) {
	var captured struct {
		Model     string `json:"model"`
		MaxTokens int    `json:"max_tokens"`
		Messages  []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	client, err := NewOpenAIClient(OpenAIConfig{
		APIKey:             "dash-test",
		BaseURL:            server.URL + "/v1",
		MiniModel:          "qwen3-max",
		InlineInstructions: true,
		UseChatCompletions: true,
	})
	if err != nil {
		t.Fatalf("NewOpenAIClient: %v", err)
	}

	result, err := client.Summarize(context.Background(), "title", "body")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}
	if result.Summary != "ok" {
		t.Fatalf("unexpected summary: %q", result.Summary)
	}
	if captured.Model != "qwen3-max" {
		t.Fatalf("unexpected model: %q", captured.Model)
	}
	if captured.MaxTokens != defaultSummaryMaxTokens {
		t.Fatalf("unexpected max_tokens: %d", captured.MaxTokens)
	}
	if len(captured.Messages) != 1 {
		t.Fatalf("expected a single user message, got %d", len(captured.Messages))
	}
	if captured.Messages[0].Role != "user" {
		t.Fatalf("unexpected role: %q", captured.Messages[0].Role)
	}
	if !containsAll(captured.Messages[0].Content, "Instructions:", "Task Input:", "Title:", "Body:") {
		t.Fatalf("expected combined prompt, got %q", captured.Messages[0].Content)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
