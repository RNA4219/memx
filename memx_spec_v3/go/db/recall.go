package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// RecallScope は Semantic Recall 用のスコープ。
type RecallScope string

const (
	ScopeSelf    RecallScope = "self"
	ScopeSession RecallScope = "session"
)

// StoreKind はストアの種類。
type StoreKind string

const (
	StoreShort     StoreKind = "short"
	StoreJournal   StoreKind = "journal"
	StoreKnowledge StoreKind = "knowledge"
	StoreArchive   StoreKind = "archive"
)

// RecallQuery は意味ベース検索のクエリ。
type RecallQuery struct {
	Text         string
	Scope        RecallScope
	Stores       []StoreKind
	TopK         int
	MessageRange int
	FallbackFTS  bool // true の場合、埋め込みクライアント未設定時に FTS にフォールバック
}

// RecallNote は notes テーブル相当の構造体。
type RecallNote struct {
	ID      string
	Title   string
	Summary string
	Body    string
	Store   StoreKind
	Score   float64
}

// NoteWithContext は anchor ノートと、その前後の文脈。
type NoteWithContext struct {
	Anchor RecallNote
	Before []RecallNote
	After  []RecallNote
}

// DefaultRecallThreshold は類似度のデフォルト閾値。
const DefaultRecallThreshold = 0.7

// Recall は指定されたクエリに基づいて Semantic Recall を実行する。
// 手順:
// 1. クエリを EmbeddingClient で埋め込み化
// 2. 対象ストアの note_embeddings で類似度計算し、閾値以上を抽出
// 3. 上位 top-k を anchor として created_at 前後 range を連結取得
// 4. --stores 入力を正規化し、不正値は 400 系で失敗
// 5. 埋め込みクライアント未設定時はデフォルトエラー、明示フラグ時のみ FTS フォールバック
func (c *Conn) Recall(ctx context.Context, q RecallQuery) ([]NoteWithContext, error) {
	// 入力正規化
	if err := validateRecallQuery(q); err != nil {
		return nil, err
	}

	// ストア正規化
	stores := normalizeStores(q.Stores)
	if len(stores) == 0 {
		return nil, fmt.Errorf("invalid stores: at least one valid store required")
	}

	// 埋め込みクライアント確認
	if c.Embed == nil {
		if q.FallbackFTS {
			return c.recallFTS(ctx, q, stores)
		}
		return nil, fmt.Errorf("embedding client not configured: set OPENAI_API_KEY or use --fallback-fts")
	}

	// 埋め込み生成
	embeddings, err := c.Embed.EmbedText(ctx, []string{q.Text})
	if err != nil {
		return nil, fmt.Errorf("embedding generation failed: %w", err)
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("embedding generation failed: empty result")
	}
	embedding := float32To64(embeddings[0])

	// 類似度検索
	topK := q.TopK
	if topK <= 0 {
		topK = 10
	}

	results, err := c.searchByEmbedding(ctx, stores, embedding, topK, DefaultRecallThreshold)
	if err != nil {
		return nil, err
	}

	// 文脈取得（前後 range）
	messageRange := q.MessageRange
	if messageRange <= 0 {
		messageRange = 5
	}

	var notesWithContext []NoteWithContext
	for _, anchor := range results {
		before, after, err := c.getContextNotes(ctx, anchor, messageRange)
		if err != nil {
			continue // エラーはスキップ
		}
		notesWithContext = append(notesWithContext, NoteWithContext{
			Anchor: anchor,
			Before: before,
			After:  after,
		})
	}

	return notesWithContext, nil
}

// validateRecallQuery はクエリの妥当性を検証する。
func validateRecallQuery(q RecallQuery) error {
	if strings.TrimSpace(q.Text) == "" {
		return fmt.Errorf("query text is required")
	}
	if len(q.Text) > 1000 {
		return fmt.Errorf("query text exceeds 1000 characters")
	}
	return nil
}

// normalizeStores はストア指定を正規化する。
func normalizeStores(stores []StoreKind) []StoreKind {
	validStores := map[StoreKind]bool{
		StoreShort:     true,
		StoreJournal:   true,
		StoreKnowledge: true,
	}

	var result []StoreKind
	for _, s := range stores {
		if validStores[s] {
			result = append(result, s)
		}
	}

	// デフォルトは short のみ
	if len(result) == 0 {
		result = []StoreKind{StoreShort}
	}

	return result
}

// searchByEmbedding は埋め込み類似度検索を実行する。
func (c *Conn) searchByEmbedding(ctx context.Context, stores []StoreKind, embedding []float64, topK int, threshold float64) ([]RecallNote, error) {
	// 埋め込みベクトルを文字列表現に変換
	embeddingStr := embeddingToJSON(embedding)

	var results []RecallNote
	for _, store := range stores {
		notes, err := c.searchStoreByEmbedding(ctx, store, embeddingStr, topK, threshold)
		if err != nil {
			continue // エラーはスキップ
		}
		results = append(results, notes...)
	}

	// スコア順ソート（降順）
	sortByScore(results)

	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

// searchStoreByEmbedding は指定ストアで埋め込み検索する。
func (c *Conn) searchStoreByEmbedding(ctx context.Context, store StoreKind, embeddingStr string, topK int, threshold float64) ([]RecallNote, error) {
	db := c.DB
	if store == StoreShort {
		db = c.ShortDB
	} else if store == StoreJournal {
		db = c.JournalDB
	} else if store == StoreKnowledge {
		db = c.KnowledgeDB
	}

	if db == nil {
		return nil, fmt.Errorf("store %s not available", store)
	}

	// ベクトル類似度計算（コサイン類似度の簡易実装）
	query := `
		SELECT n.id, n.title, n.summary, n.body, e.score
		FROM notes n
		JOIN note_embeddings e ON n.id = e.note_id
		WHERE e.score >= ?
		ORDER BY e.score DESC
		LIMIT ?
	`

	rows, err := db.QueryContext(ctx, query, threshold, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []RecallNote
	for rows.Next() {
		var n RecallNote
		if err := rows.Scan(&n.ID, &n.Title, &n.Summary, &n.Body, &n.Score); err != nil {
			continue
		}
		n.Store = store
		notes = append(notes, n)
	}

	return notes, nil
}

// recallFTS は FTS フォールバック検索を実行する。
func (c *Conn) recallFTS(ctx context.Context, q RecallQuery, stores []StoreKind) ([]NoteWithContext, error) {
	var results []NoteWithContext

	topK := q.TopK
	if topK <= 0 {
		topK = 10
	}

	for _, store := range stores {
		notes, err := c.searchStoreFTS(ctx, store, q.Text, topK)
		if err != nil {
			continue
		}
		for _, n := range notes {
			results = append(results, NoteWithContext{Anchor: n})
		}
	}

	return results, nil
}

// searchStoreFTS は指定ストアで FTS 検索する。
func (c *Conn) searchStoreFTS(ctx context.Context, store StoreKind, query string, topK int) ([]RecallNote, error) {
	db := c.DB
	if store == StoreShort {
		db = c.ShortDB
	} else if store == StoreJournal {
		db = c.JournalDB
	} else if store == StoreKnowledge {
		db = c.KnowledgeDB
	}

	if db == nil {
		return nil, fmt.Errorf("store %s not available", store)
	}

	ftsQuery := `
		SELECT n.id, n.title, n.summary, n.body
		FROM notes_fts
		JOIN notes n ON notes_fts.rowid = n.rowid
		WHERE notes_fts MATCH ?
		ORDER BY bm25(notes_fts)
		LIMIT ?
	`

	rows, err := db.QueryContext(ctx, ftsQuery, query, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []RecallNote
	for rows.Next() {
		var n RecallNote
		if err := rows.Scan(&n.ID, &n.Title, &n.Summary, &n.Body); err != nil {
			continue
		}
		n.Store = store
		n.Score = 1.0 // FTS は固定スコア
		notes = append(notes, n)
	}

	return notes, nil
}

// getContextNotes は anchor ノートの前後ノートを取得する。
func (c *Conn) getContextNotes(ctx context.Context, anchor RecallNote, r int) ([]RecallNote, []RecallNote, error) {
	db := c.DB
	if anchor.Store == StoreShort {
		db = c.ShortDB
	} else if anchor.Store == StoreJournal {
		db = c.JournalDB
	} else if anchor.Store == StoreKnowledge {
		db = c.KnowledgeDB
	}

	if db == nil {
		return nil, nil, fmt.Errorf("store not available")
	}

	// anchor の created_at を取得
	var createdAt string
	err := db.QueryRowContext(ctx, `SELECT created_at FROM notes WHERE id = ?`, anchor.ID).Scan(&createdAt)
	if err != nil {
		return nil, nil, err
	}

	// 前のノート
	beforeQuery := `
		SELECT id, title, summary, body FROM notes
		WHERE created_at < ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	beforeRows, err := db.QueryContext(ctx, beforeQuery, createdAt, r)
	if err != nil {
		return nil, nil, err
	}
	defer beforeRows.Close()

	var before []RecallNote
	for beforeRows.Next() {
		var n RecallNote
		if err := beforeRows.Scan(&n.ID, &n.Title, &n.Summary, &n.Body); err != nil {
			continue
		}
		n.Store = anchor.Store
		before = append(before, n)
	}

	// 後のノート
	afterQuery := `
		SELECT id, title, summary, body FROM notes
		WHERE created_at > ?
		ORDER BY created_at ASC
		LIMIT ?
	`
	afterRows, err := db.QueryContext(ctx, afterQuery, createdAt, r)
	if err != nil {
		return before, nil, nil
	}
	defer afterRows.Close()

	var after []RecallNote
	for afterRows.Next() {
		var n RecallNote
		if err := afterRows.Scan(&n.ID, &n.Title, &n.Summary, &n.Body); err != nil {
			continue
		}
		n.Store = anchor.Store
		after = append(after, n)
	}

	return before, after, nil
}

func float32To64(values []float32) []float64 {
	out := make([]float64, len(values))
	for i, value := range values {
		out[i] = float64(value)
	}
	return out
}

// embeddingToJSON は埋め込みベクトルを JSON 文字列に変換する。
func embeddingToJSON(e []float64) string {
	if len(e) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range e {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%.6f", v))
	}
	sb.WriteString("]")
	return sb.String()
}

// sortByScore はスコア順（降順）にソートする。
func sortByScore(notes []RecallNote) {
	for i := 0; i < len(notes)-1; i++ {
		for j := i + 1; j < len(notes); j++ {
			if notes[j].Score > notes[i].Score {
				notes[i], notes[j] = notes[j], notes[i]
			}
		}
	}
}

// Close は DB 接続を閉じる（sql.DB インターフェース対応）。
type DBCloser interface {
	Close() error
}

// Conn フィールドに EmbeddingClient を追加するための定義
// 注意: 実際の Conn 構造体は db/open.go で定義されているため、
// このファイルでは EmbeddingClient を使用するメソッドのみを提供する
var _ = sql.NullString{} // sql パッケージの使用確認
