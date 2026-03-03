---
intent_id: memx-operations-runbook-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# RUNBOOK

## v1 必須運用フロー
1. **ingest**: `mem in short`（または `POST /v1/notes:ingest`）で short へ投入。
2. **search**: `mem out search`（または `POST /v1/notes:search`）で FTS 検索。
3. **show**: `mem out show`（または `GET /v1/notes/{id}`）で単一ノート参照。

## `mem in short` 手順
1. CLI で `file` または stdin から本文を読み込み、request を生成。
2. API へ request を送信（in-proc または HTTP）。
3. Service 側で Gatekeeper(`memory_store`)・ノート保存・タグ/埋め込み更新・`short_meta` 更新を実行。
4. CLI が note id 等を整形表示。

## `mem out recall` 手順
1. クエリを EmbeddingClient で埋め込み化。
2. 対象ストアの `note_embeddings` で類似度計算し、閾値以上を抽出。
3. 上位 `top-k` を anchor として `created_at` 前後 `range` を連結取得。
4. `--stores` 入力を正規化し、不正値は 400 系で失敗。
5. 埋め込みクライアント未設定時はデフォルトエラー、明示フラグ時のみ FTS フォールバック。

## `mem gc short` 手順
1. `short_meta` と `memory_policy.yaml.gc.short` からトリガ判定。
2. soft/hard limit 条件と `min_interval_minutes` で実行可否を決定。
3. 実行時は正確値を再計算して閾値を再確認。
4. `--dry-run` は DB 変更せず予定操作 JSON のみ返却。

## LLM クライアント運用
- `EmbeddingClient`: 埋め込み生成。
- `MiniLLMClient`: タグ・スコア・機密度推定。
- `ReflectLLMClient`: Observer/Reflector 要約更新。
- タイムアウト 15 秒、最大 2 回リトライ（指数バックオフ）、再試行可/不可を区別して実装する。
