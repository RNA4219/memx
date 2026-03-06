---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - LLM 戦略

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 4. LLM 戦略

### 4-1. 役割分離

LLM は少なくとも 3 役割に分離する：

1. EmbeddingClient
   - テキスト → ベクター 変換。
   - Semantic Recall で使用。
2. MiniLLMClient
   - タグ生成・スコアリング（`relevance / quality / novelty / importance`）・`sensitivity` 推定。
   - 軽量モデル（1B〜3B）を想定。
3. ReflectLLMClient
   - クラスタ要約（Observer）・Memopedia ページ更新（Reflector）。
   - 7B〜27B クラスのモデルを想定。

Go 側では `go/db/llm_client.go` に interface を定義し、`db.Conn` にこれらをフィールドとして注入して使う。

### 4-2. 同期／非同期の扱い

- `mem in` 実行時に全ての LLM を同期で呼ぶとレイテンシが伸びる。
- v1 では：
  - `mem in` では最低限のフィールド（title/body/source_*）だけ即時保存し、
  - タグ付け・スコアリング・埋め込み生成はオプションで非同期キューに積む実装も許容範囲とする。
- CLI としては：
  - `--no-llm` で完全に LLM を使わない形
  - デフォルトでは同期処理（ただし後でオプションで非同期化も検討）
- 要件レベルでは、「LLM を使うか／どのタイミングで使うか」を `mem in` のフラグと設定ファイルで切り替え可能にする。

### 4-3. 設定例

`config.yaml` のイメージ：

```yaml
llm:
  embed:
    provider: local
    endpoint: "http://localhost:8000/embed"
  mini:
    provider: local
    endpoint: "http://localhost:8001/generate"
  reflect:
    provider: local
    endpoint: "http://localhost:8002/generate"
```

`memory_policy.yaml`（GC 関連キー雛形）：

```yaml
version: 1

gc:
  short:
    soft_limit_notes: 1200
    hard_limit_notes: 2000
    min_interval_minutes: 180
    target_delete_batch_size: 200
    max_archive_retries: 1
```

### 4-4. 各クライアント共通の呼び出し契約

対象：`EmbeddingClient` / `MiniLLMClient` / `ReflectLLMClient` / Gatekeeper 呼び出し。

- タイムアウト：**1 リクエスト 15 秒**（`context.WithTimeout` 等で必須化）。
- 最大リトライ回数：**2 回**（初回 + リトライ 2 = 最大 3 試行、指数バックオフ）。
- 再試行可エラー：
  - ネットワーク断、接続リセット、タイムアウト
  - HTTP 429 / 502 / 503 / 504
- 再試行不可エラー：
  - 入力不正（HTTP 400 相当）
  - 認証/認可失敗（HTTP 401/403 相当）
  - モデル仕様不一致（JSON スキーマ不整合、必須フィールド欠落）
- ingest 時の部分失敗ポリシー：
  - `notes` 保存成功をコミット境界の最小単位とし、ノート本体保存は継続。
  - タグ生成・`note_tags`・埋め込み生成の失敗は後追い再実行対象として記録し、ingest 全体は成功扱いにできる。
  - ただし Gatekeeper が `deny` / `needs_human` を返した場合は fail-closed で ingest 全体を失敗にする。
- エラーコードのマッピング方針（`go/api/errors.go` 整合）：
  - API 返却コードは `INVALID_ARGUMENT` / `NOT_FOUND` / `INTERNAL` を最小集合として保証する。
  - クライアント個別エラーは service 層で sentinel error に正規化し、`go/api/errors.go` の `mapError` に 1 箇所で集約マップする。
  - 未分類エラーは互換性維持のため常に `INTERNAL` にフォールバックする。