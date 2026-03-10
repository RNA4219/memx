# cookbook-resolver design

## 1. 文書情報

- 文書名: cookbook-resolver design
- 文書種別: design
- 版: v0.1
- 作成日: 2026-03-10
- 状態: Draft

## 2. 目的

本書は `docs/requirements.md` と `docs/interfaces.md` を満たすための最小実装方針を定義する。

対象は以下とする。

- 文書登録
- chunk 化
- 文書解決
- chunk 取得
- 読了記録
- stale 判定
- 契約解決

## 3. 実装方針

### 3.1 配置

最小実装は `memx_spec_v3/go` 配下に追加し、既存の `memx-core` 系 API と同じプロセスで提供する。

- DB: `short.db` に resolver 用テーブルを追加する
- service: resolver 用 usecase を追加する
- api: `/v1/docs:*` などの HTTP API を追加する
- client: in-proc / HTTP client の両方から同じ API を呼べるようにする

### 3.2 agent-taskstate 連携

本段階では `agent-taskstate` に直接書き込まない。

代わりに以下を行う。

- `read_receipt` 相当を `memx-resolver` 内に正規化して保持する
- `task_id` をキーに stale 判定を返す API を提供する
- 将来 `agent-taskstate` へ移譲しやすいように payload 形状を `docs/interfaces.md` に合わせる

### 3.3 chunking

chunk 生成は見出し優先とする。

- Markdown 見出しを検出して section 化する
- section が長すぎる場合のみ固定長で再分割する
- `importance` は見出し名と `doc_type` から推定する

### 3.4 解決ロジック

文書解決は軽量な決定的ロジックを採用する。

- `task_id` 一致を最優先
- 次に `feature_keys` 一致
- 次に `tags` / `title` / `summary` / `body` の部分一致
- `required` / `recommended` は doc importance で振り分ける

### 3.5 契約情報

契約情報は以下の 2 系統から集約する。

- ingest 時に明示入力された配列
- 見出し名から抽出した section

抽出対象は以下とする。

- Acceptance Criteria
- Forbidden Patterns
- Definition of Done
- Dependencies

## 4. データモデル

### 4.1 resolver_documents

文書本体と契約系メタデータを保持する。

主な項目:

- `doc_id`
- `doc_type`
- `title`
- `source_path`
- `version`
- `updated_at`
- `summary`
- `body`
- `tags_json`
- `feature_keys_json`
- `task_ids_json`
- `acceptance_criteria_json`
- `forbidden_patterns_json`
- `definition_of_done_json`
- `dependencies_json`
- `importance`

### 4.2 resolver_chunks

文書の参照単位を保持する。

主な項目:

- `chunk_id`
- `doc_id`
- `heading`
- `heading_path_json`
- `ordinal`
- `body`
- `token_estimate`
- `importance`

### 4.3 resolver_document_links

文書間依存を保持する。

主な項目:

- `src_doc_id`
- `dst_doc_id`
- `link_type`

### 4.4 resolver_read_receipts

task と文書参照の対応を保持する。

主な項目:

- `task_id`
- `doc_id`
- `version`
- `chunk_ids_json`
- `reader`
- `read_at`

## 5. 既知の制約

- version 比較は最小実装として完全な順序比較を行わず、文字列不一致を stale とみなす
- task dependency は外部正本に問い合わせず、ローカル保持の `task_ids_json` を優先利用する
- 全文検索は resolver 専用 FTS を作らず、最小実装では LIKE ベースとする

## 6. 完了条件

- `docs/interfaces.md` に記載した最小 API が呼べる
- 文書を登録すると chunk が生成される
- feature / task / topic から required / recommended docs を返せる
- read receipt 登録と stale 判定が動く
- contract resolve が acceptance / forbidden / DoD / dependencies を返せる
