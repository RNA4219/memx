package db

import "context"

// EmbeddingClient はテキストから埋め込みベクトルを生成する。
type EmbeddingClient interface {
	EmbedText(ctx context.Context, texts []string) ([][]float32, error)
}

// TagsAndScores はタグと各種スコアの推定結果。
type TagsAndScores struct {
	Tags             []string
	Relevance        float64
	Quality          float64
	Novelty          float64
	ImportanceStatic float64
	Sensitivity      string
}

// SummarizeResult は要約生成の結果。
type SummarizeResult struct {
	Summary string
}

// MiniLLMClient はタグ付けとスコアリングを担当する軽量モデルのインターフェース。
type MiniLLMClient interface {
	TagAndScore(ctx context.Context, noteBody string) (TagsAndScores, error)
	Summarize(ctx context.Context, title, body string) (SummarizeResult, error)
}

// ClusterInput は ReflectLLM が観測ノート群から要約を生成する際の入力。
type ClusterInput struct {
	NoteIDs []string
	Body    string
}

// PageUpdateInput は Memopedia ページ更新時の入力。
type PageUpdateInput struct {
	PageID          string
	ExistingContent string
	NewObservations []string
}

// ReflectLLMClient は観測ノート要約と Memopedia ページ更新を担当する。
type ReflectLLMClient interface {
	SummarizeCluster(ctx context.Context, cluster ClusterInput) (string, error)
	UpdateMemopediaPage(ctx context.Context, input PageUpdateInput) (string, error)
}
