package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultOpenAIBaseURL          = "https://api.openai.com/v1"
	defaultOpenAIMiniModel        = "gpt-5-mini"
	defaultOpenAIMaxTimeout       = 30 * time.Second
	defaultAlibabaIntlBaseURL     = "https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1"
	defaultAlibabaBeijingBaseURL  = "https://dashscope.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1"
	defaultAlibabaCompatibleModel = "qwen3-max"
	defaultSummaryMaxTokens       = 180
	defaultClusterMaxTokens       = 280
	defaultTagScoreMaxTokens      = 220
	defaultKnowledgeMaxTokens     = 400
)

// OpenAIConfig は memx で使う OpenAI 互換 LLM 接続設定。
type OpenAIConfig struct {
	APIKey             string
	BaseURL            string
	MiniModel          string
	ReflectModel       string
	Project            string
	Organization       string
	Timeout            time.Duration
	InlineInstructions bool
}

// OpenAIClient は OpenAI Responses API 互換の LLM クライアント。
type OpenAIClient struct {
	apiKey             string
	baseURL            string
	miniModel          string
	reflectModel       string
	project            string
	organization       string
	inlineInstructions bool
	httpClient         *http.Client
}

// NewOpenAIClient は設定済みの OpenAI 互換クライアントを作成する。
func NewOpenAIClient(cfg OpenAIConfig) (*OpenAIClient, error) {
	apiKey := strings.TrimSpace(cfg.APIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("openai-compatible api key is required")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultOpenAIBaseURL
	}

	miniModel := strings.TrimSpace(cfg.MiniModel)
	if miniModel == "" {
		miniModel = defaultOpenAIMiniModel
	}

	reflectModel := strings.TrimSpace(cfg.ReflectModel)
	if reflectModel == "" {
		reflectModel = miniModel
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = defaultOpenAIMaxTimeout
	}

	return &OpenAIClient{
		apiKey:             apiKey,
		baseURL:            baseURL,
		miniModel:          miniModel,
		reflectModel:       reflectModel,
		project:            strings.TrimSpace(cfg.Project),
		organization:       strings.TrimSpace(cfg.Organization),
		inlineInstructions: cfg.InlineInstructions,
		httpClient:         &http.Client{Timeout: timeout},
	}, nil
}

// NewOpenAIClientFromEnv は環境変数から OpenAI 互換クライアントを作成する。
// OpenAI または Alibaba Cloud Model Studio の設定を自動検出する。
func NewOpenAIClientFromEnv() (*OpenAIClient, bool, error) {
	cfg, ok, err := LoadOpenAIConfigFromEnv()
	if err != nil || !ok {
		return nil, ok, err
	}
	client, err := NewOpenAIClient(cfg)
	if err != nil {
		return nil, false, err
	}
	return client, true, nil
}

// LoadOpenAIConfigFromEnv は環境変数から OpenAI 互換設定を読み出す。
// MEMX_LLM_PROVIDER=openai|alibaba で明示指定できる。未指定時は OpenAI → Alibaba の順に自動検出する。
func LoadOpenAIConfigFromEnv() (OpenAIConfig, bool, error) {
	switch normalizeProvider(firstNonEmptyEnv("MEMX_LLM_PROVIDER")) {
	case "":
		if cfg, ok, err := loadOpenAICompatibleConfigFromEnv(); ok || err != nil {
			return cfg, ok, err
		}
		return loadAlibabaCompatibleConfigFromEnv()
	case "openai":
		return loadOpenAICompatibleConfigFromEnv()
	case "alibaba":
		return loadAlibabaCompatibleConfigFromEnv()
	default:
		return OpenAIConfig{}, false, fmt.Errorf("unsupported MEMX_LLM_PROVIDER: %q", firstNonEmptyEnv("MEMX_LLM_PROVIDER"))
	}
}

func loadOpenAICompatibleConfigFromEnv() (OpenAIConfig, bool, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("MEMX_OPENAI_API_KEY"))
	}
	if apiKey == "" {
		return OpenAIConfig{}, false, nil
	}

	sharedModel := firstNonEmptyEnv("MEMX_OPENAI_MODEL")
	miniModel := firstNonEmpty(firstNonEmptyEnv("MEMX_OPENAI_MINI_MODEL"), sharedModel)
	reflectModel := firstNonEmpty(firstNonEmptyEnv("MEMX_OPENAI_REFLECT_MODEL"), miniModel, sharedModel)
	baseURL := firstNonEmpty(firstNonEmptyEnv("MEMX_OPENAI_BASE_URL"), firstNonEmptyEnv("OPENAI_BASE_URL"), defaultOpenAIBaseURL)
	timeout, err := loadTimeoutFromEnv("MEMX_OPENAI_TIMEOUT_SECONDS")
	if err != nil {
		return OpenAIConfig{}, false, err
	}

	return OpenAIConfig{
		APIKey:             apiKey,
		BaseURL:            baseURL,
		MiniModel:          firstNonEmpty(miniModel, defaultOpenAIMiniModel),
		ReflectModel:       firstNonEmpty(reflectModel, miniModel, defaultOpenAIMiniModel),
		Project:            firstNonEmpty(firstNonEmptyEnv("OPENAI_PROJECT"), firstNonEmptyEnv("MEMX_OPENAI_PROJECT")),
		Organization:       firstNonEmpty(firstNonEmptyEnv("OPENAI_ORGANIZATION"), firstNonEmptyEnv("MEMX_OPENAI_ORGANIZATION")),
		Timeout:            timeout,
		InlineInstructions: false,
	}, true, nil
}

func loadAlibabaCompatibleConfigFromEnv() (OpenAIConfig, bool, error) {
	apiKey := firstNonEmpty(firstNonEmptyEnv("DASHSCOPE_API_KEY"), firstNonEmptyEnv("MEMX_ALIBABA_API_KEY"), firstNonEmptyEnv("MEMX_DASHSCOPE_API_KEY"))
	if apiKey == "" {
		return OpenAIConfig{}, false, nil
	}

	sharedModel := firstNonEmpty(firstNonEmptyEnv("MEMX_ALIBABA_MODEL"), firstNonEmptyEnv("MEMX_DASHSCOPE_MODEL"))
	miniModel := firstNonEmpty(firstNonEmptyEnv("MEMX_ALIBABA_MINI_MODEL"), firstNonEmptyEnv("MEMX_DASHSCOPE_MINI_MODEL"), sharedModel, defaultAlibabaCompatibleModel)
	reflectModel := firstNonEmpty(firstNonEmptyEnv("MEMX_ALIBABA_REFLECT_MODEL"), firstNonEmptyEnv("MEMX_DASHSCOPE_REFLECT_MODEL"), sharedModel, miniModel)
	baseURL := firstNonEmpty(firstNonEmptyEnv("MEMX_ALIBABA_BASE_URL"), firstNonEmptyEnv("MEMX_DASHSCOPE_BASE_URL"), firstNonEmptyEnv("DASHSCOPE_BASE_URL"))
	if baseURL == "" {
		var err error
		baseURL, err = alibabaCompatibleBaseURL(firstNonEmptyEnv("MEMX_ALIBABA_REGION", "MEMX_DASHSCOPE_REGION", "DASHSCOPE_REGION"))
		if err != nil {
			return OpenAIConfig{}, false, err
		}
	}
	timeout, err := loadTimeoutFromEnv("MEMX_ALIBABA_TIMEOUT_SECONDS", "MEMX_DASHSCOPE_TIMEOUT_SECONDS")
	if err != nil {
		return OpenAIConfig{}, false, err
	}

	return OpenAIConfig{
		APIKey:             apiKey,
		BaseURL:            baseURL,
		MiniModel:          miniModel,
		ReflectModel:       reflectModel,
		Timeout:            timeout,
		InlineInstructions: true,
	}, true, nil
}

func loadTimeoutFromEnv(keys ...string) (time.Duration, error) {
	timeoutSeconds := firstNonEmptyEnv(keys...)
	if timeoutSeconds == "" {
		return defaultOpenAIMaxTimeout, nil
	}
	seconds, err := strconv.Atoi(timeoutSeconds)
	if err != nil || seconds <= 0 {
		return 0, fmt.Errorf("invalid timeout seconds in %v: %q", keys, timeoutSeconds)
	}
	return time.Duration(seconds) * time.Second, nil
}

func alibabaCompatibleBaseURL(region string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(region)) {
	case "", "sg", "singapore", "intl", "international":
		return defaultAlibabaIntlBaseURL, nil
	case "cn", "beijing", "china", "china-beijing":
		return defaultAlibabaBeijingBaseURL, nil
	default:
		return "", fmt.Errorf("unsupported Alibaba region for Responses API: %q", region)
	}
}

func normalizeProvider(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

// TagAndScore はメモ本文からタグとスコアを推定する。
func (c *OpenAIClient) TagAndScore(ctx context.Context, noteBody string) (TagsAndScores, error) {
	const instructions = "You extract memory tags and scores. Return JSON only with keys: tags, relevance, quality, novelty, importance_static, sensitivity. tags is an array of up to 5 short strings. Scores are numbers between 0 and 1. sensitivity must be one of public, internal, confidential, secret."

	raw, err := c.runResponses(ctx, c.miniModel, instructions, noteBody, defaultTagScoreMaxTokens)
	if err != nil {
		return TagsAndScores{}, err
	}

	var parsed struct {
		Tags             []string `json:"tags"`
		Relevance        float64  `json:"relevance"`
		Quality          float64  `json:"quality"`
		Novelty          float64  `json:"novelty"`
		ImportanceStatic float64  `json:"importance_static"`
		Sensitivity      string   `json:"sensitivity"`
	}
	if err := json.Unmarshal([]byte(stripCodeFence(raw)), &parsed); err != nil {
		return TagsAndScores{}, fmt.Errorf("parse tag-and-score response: %w", err)
	}

	return TagsAndScores{
		Tags:             parsed.Tags,
		Relevance:        parsed.Relevance,
		Quality:          parsed.Quality,
		Novelty:          parsed.Novelty,
		ImportanceStatic: parsed.ImportanceStatic,
		Sensitivity:      normalizeSensitivity(parsed.Sensitivity),
	}, nil
}

// Summarize は単一メモの要約を生成する。
func (c *OpenAIClient) Summarize(ctx context.Context, title, body string) (SummarizeResult, error) {
	const instructions = "You summarize memory notes for a local agent memory system. Return only plain summary text in 1 to 3 sentences. Preserve the original language when clear; if mixed or unclear, default to Japanese. No markdown, bullets, or headings."

	input := fmt.Sprintf("Title:\n%s\n\nBody:\n%s", strings.TrimSpace(title), strings.TrimSpace(body))
	summary, err := c.runResponses(ctx, c.miniModel, instructions, input, defaultSummaryMaxTokens)
	if err != nil {
		return SummarizeResult{}, err
	}
	return SummarizeResult{Summary: summary}, nil
}

// SummarizeCluster は複数ノートの統合要約を生成する。
func (c *OpenAIClient) SummarizeCluster(ctx context.Context, cluster ClusterInput) (string, error) {
	const instructions = "You summarize a cluster of memory notes into one compact synthesis. Return only plain summary text in 2 to 5 sentences. Preserve the original language when clear; if mixed or unclear, default to Japanese. Merge duplicate facts and keep important actions, decisions, and constraints."

	input := fmt.Sprintf("Note IDs: %s\n\nCluster Body:\n%s", strings.Join(cluster.NoteIDs, ", "), strings.TrimSpace(cluster.Body))
	return c.runResponses(ctx, c.reflectModel, instructions, input, defaultClusterMaxTokens)
}

// UpdateKnowledgePage は既存の Knowledge ページに新しい観測を統合する。
func (c *OpenAIClient) UpdateKnowledgePage(ctx context.Context, input PageUpdateInput) (string, error) {
	const instructions = "You update a knowledge page using new observations. Return only the full updated markdown document. Preserve the original language when clear; if mixed or unclear, default to Japanese. Keep stable facts, merge duplicates, and integrate only substantiated new observations."

	prompt := fmt.Sprintf(
		"Page ID: %s\n\nExisting Content:\n%s\n\nNew Observations:\n- %s",
		strings.TrimSpace(input.PageID),
		strings.TrimSpace(input.ExistingContent),
		strings.Join(input.NewObservations, "\n- "),
	)
	return c.runResponses(ctx, c.reflectModel, instructions, prompt, defaultKnowledgeMaxTokens)
}

func (c *OpenAIClient) runResponses(ctx context.Context, model, instructions, input string, maxOutputTokens int) (string, error) {
	if c.inlineInstructions && strings.TrimSpace(instructions) != "" {
		input = fmt.Sprintf("Follow these instructions carefully.\n\nInstructions:\n%s\n\nTask Input:\n%s", strings.TrimSpace(instructions), strings.TrimSpace(input))
		instructions = ""
	}

	reqBody := openAIResponsesRequest{
		Model:           model,
		Input:           input,
		Instructions:    instructions,
		MaxOutputTokens: maxOutputTokens,
		Store:           false,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal responses request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create responses request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	if c.project != "" {
		req.Header.Set("OpenAI-Project", c.project)
	}
	if c.organization != "" {
		req.Header.Set("OpenAI-Organization", c.organization)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call responses api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read responses api body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var apiErr openAIErrorResponse
		if err := json.Unmarshal(body, &apiErr); err == nil && strings.TrimSpace(apiErr.Error.Message) != "" {
			return "", fmt.Errorf("openai responses api: %s", strings.TrimSpace(apiErr.Error.Message))
		}
		return "", fmt.Errorf("openai responses api: status %d", resp.StatusCode)
	}

	var decoded openAIResponsesResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", fmt.Errorf("decode responses api body: %w", err)
	}

	text := extractResponseText(decoded)
	if text == "" {
		return "", fmt.Errorf("openai responses api returned empty text")
	}
	return text, nil
}

func extractResponseText(resp openAIResponsesResponse) string {
	if text := strings.TrimSpace(resp.OutputText); text != "" {
		return text
	}

	var parts []string
	for _, item := range resp.Output {
		if text := strings.TrimSpace(item.Text); text != "" {
			parts = append(parts, text)
		}
		for _, content := range item.Content {
			if content.Type != "output_text" && content.Type != "text" {
				continue
			}
			if text := strings.TrimSpace(content.Text); text != "" {
				parts = append(parts, text)
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func stripCodeFence(s string) string {
	trimmed := strings.TrimSpace(s)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}

	lines := strings.Split(trimmed, "\n")
	if len(lines) == 0 {
		return trimmed
	}
	lines = lines[1:]
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
		lines = lines[:len(lines)-1]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func normalizeSensitivity(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "public", "internal", "confidential", "secret":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "internal"
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}
	return ""
}

type openAIResponsesRequest struct {
	Model           string `json:"model"`
	Input           string `json:"input"`
	Instructions    string `json:"instructions,omitempty"`
	MaxOutputTokens int    `json:"max_output_tokens,omitempty"`
	Store           bool   `json:"store"`
}

type openAIResponsesResponse struct {
	OutputText string               `json:"output_text"`
	Output     []openAIResponseItem `json:"output"`
}

type openAIResponseItem struct {
	Type    string                  `json:"type"`
	Text    string                  `json:"text"`
	Content []openAIResponseContent `json:"content"`
}

type openAIResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
