# cookbook-resolver interfaces

## 1. 文書情報

- 文書名: cookbook-resolver interfaces
- 文書種別: interfaces
- 版: v0.1
- 作成日: 2026-03-10
- 状態: Draft

## 2. 目的

本書は、`cookbook-resolver` の外部インターフェースを定義する。

対象は以下とする。

- resolver API
- chunk 取得 API
- 読了記録 API
- stale 判定 API
- Skill 入口
- `agent-taskstate` との連携面

本書では、実装言語やフレームワークには依存せず、入出力契約と責務境界を定義する。

## 3. 設計方針

### 3.1 基本方針

- 文書原本は Markdown のまま保持する
- エージェントには全文ではなく必要情報を返す
- 参照解決と本文取得を分離する
- 読了結果は task と結びつけて保持する
- stale 判定は doc version と read receipt に基づいて行う

### 3.2 責務分離

- `cookbook-resolver`
  - 文書登録
  - 文書解決
  - chunk 取得
  - 読了記録の受付
  - stale 判定の問い合わせ窓口
- `memx-core`
  - 文書・chunk の保持
  - 検索・要約・取得
- `agent-taskstate`
  - task の正本
  - read receipt の正本
  - task dependency の正本
  - stale 状態の正本
- `tracker-bridge-materials`
  - 将来の外部 tracker / Birdseye 連携

## 4. 共通仕様

### 4.1 ID 規約

各エンティティは stable ID を持つこと。

想定例

- `doc:spec:memory-import`
- `doc:cookbook:hub-codex`
- `chunk:doc:spec:memory-import:001`
- `task:feature:local:123`
- `feature:memory-import`

### 4.2 version 規約

`version` は最低限、比較可能な文字列であること。

推奨形式

- 日付ベース
  - `2026-03-10`
- 日時ベース
  - `2026-03-10T13:00:00+09:00`
- リビジョン付き
  - `2026-03-10.1`

### 4.3 エラー応答方針

全 API は失敗時に以下を返せること。

- `error_code`
- `message`
- `details`

例

```json
{
  "error_code": "DOC_NOT_FOUND",
  "message": "document was not found",
  "details": {
    "doc_id": "doc:spec:missing"
  }
}
````

### 4.4 importance 値

chunk または doc の重要度として以下を使用する。

* `required`
* `recommended`
* `reference`

## 5. データ構造

## 5.1 Document

```json
{
  "doc_id": "doc:spec:memory-import",
  "doc_type": "spec",
  "title": "Memory Import Spec",
  "source_path": "docs/specs/memory-import.md",
  "version": "2026-03-10",
  "updated_at": "2026-03-10T09:00:00+09:00",
  "summary": "メモリ取り込み処理の仕様",
  "tags": ["memory", "import"],
  "feature_keys": ["memory-import"]
}
```

### 5.2 DocumentChunk

```json
{
  "chunk_id": "chunk:doc:spec:memory-import:001",
  "doc_id": "doc:spec:memory-import",
  "heading_path": ["Memory Import Spec", "Acceptance Criteria"],
  "ordinal": 1,
  "body": "受け入れ条件は...",
  "token_estimate": 280,
  "importance": "required"
}
```

### 5.3 ResolveEntry

```json
{
  "doc_id": "doc:spec:memory-import",
  "title": "Memory Import Spec",
  "version": "2026-03-10",
  "importance": "required",
  "reason": "core spec for this feature",
  "top_chunks": [
    "chunk:doc:spec:memory-import:001",
    "chunk:doc:spec:memory-import:002"
  ]
}
```

### 5.4 ReadReceipt

```json
{
  "task_id": "task:feature:local:123",
  "doc_id": "doc:spec:memory-import",
  "version": "2026-03-10",
  "chunk_ids": [
    "chunk:doc:spec:memory-import:001",
    "chunk:doc:spec:memory-import:002"
  ],
  "reader": "agent",
  "read_at": "2026-03-10T10:00:00+09:00"
}
```

### 5.5 StaleReason

```json
{
  "task_id": "task:feature:local:123",
  "doc_id": "doc:spec:memory-import",
  "previous_version": "2026-03-09",
  "current_version": "2026-03-10",
  "reason": "version_mismatch",
  "detected_at": "2026-03-10T10:10:00+09:00"
}
```

## 6. API 一覧

最低限の API は以下とする。

* `POST /v1/docs:ingest`
* `POST /v1/docs:resolve`
* `POST /v1/chunks:get`
* `POST /v1/docs:search`
* `POST /v1/reads:ack`
* `POST /v1/docs:stale-check`
* `POST /v1/contracts:resolve`

## 7. API 詳細

## 7.1 POST `/v1/docs:ingest`

### 概要

Markdown 文書を登録し、必要に応じて chunk 化する。

### リクエスト

```json
{
  "doc_type": "spec",
  "title": "Memory Import Spec",
  "source_path": "docs/specs/memory-import.md",
  "version": "2026-03-10",
  "updated_at": "2026-03-10T09:00:00+09:00",
  "tags": ["memory", "import"],
  "feature_keys": ["memory-import"],
  "summary": "メモリ取り込み処理の仕様",
  "body": "# Memory Import Spec\n...",
  "chunking": {
    "mode": "heading",
    "max_chars": 4000
  }
}
```

### バリデーション

* `doc_type` は必須
* `title` は必須
* `version` は必須
* `body` は必須
* `chunking.mode` は `heading` または `fixed` を許可
* 同一 `doc_id` 相当の再登録は version ルールに従うこと

### レスポンス

```json
{
  "doc_id": "doc:spec:memory-import",
  "version": "2026-03-10",
  "chunk_count": 4,
  "status": "ingested"
}
```

### 失敗例

* `INVALID_REQUEST`
* `UNSUPPORTED_DOC_TYPE`
* `VERSION_CONFLICT`

## 7.2 POST `/v1/docs:resolve`

### 概要

feature 名、task_id、topic のいずれかを入力に、読むべき doc を解決する。

### リクエスト

```json
{
  "feature": "memory-import",
  "task_id": "task:feature:local:123",
  "topic": null,
  "limit": 10
}
```

### 入力ルール

* `feature`、`task_id`、`topic` のいずれか1つ以上を必須とする
* `limit` は省略可能
* `task_id` が与えられた場合、task dependency 情報を参照してもよい

### レスポンス

```json
{
  "required": [
    {
      "doc_id": "doc:spec:memory-import",
      "title": "Memory Import Spec",
      "version": "2026-03-10",
      "importance": "required",
      "reason": "core spec",
      "top_chunks": [
        "chunk:doc:spec:memory-import:001",
        "chunk:doc:spec:memory-import:002"
      ]
    }
  ],
  "recommended": [
    {
      "doc_id": "doc:cookbook:import-patterns",
      "title": "Import Patterns",
      "version": "2026-03-09",
      "importance": "recommended",
      "reason": "implementation guidance",
      "top_chunks": [
        "chunk:doc:cookbook:import-patterns:003"
      ]
    }
  ]
}
```

### 解決優先順位

推奨順序は以下とする。

* task dependency に紐づく doc
* feature_keys が一致する doc
* tags が一致する doc
* topic に対する全文検索の上位 doc

## 7.3 POST `/v1/chunks:get`

### 概要

必要な chunk または本文の一部を返す。

### リクエスト例1 doc 単位

```json
{
  "doc_id": "doc:spec:memory-import",
  "limit": 3
}
```

### リクエスト例2 query 指定

```json
{
  "doc_id": "doc:spec:memory-import",
  "query": "acceptance criteria",
  "limit": 3
}
```

### リクエスト例3 heading 指定

```json
{
  "doc_id": "doc:spec:memory-import",
  "heading": "Acceptance Criteria",
  "limit": 10
}
```

### リクエスト例4 chunk_id 直接指定

```json
{
  "chunk_ids": [
    "chunk:doc:spec:memory-import:001",
    "chunk:doc:spec:memory-import:002"
  ]
}
```

### レスポンス

```json
{
  "doc_id": "doc:spec:memory-import",
  "chunks": [
    {
      "chunk_id": "chunk:doc:spec:memory-import:001",
      "heading_path": ["Memory Import Spec", "Acceptance Criteria"],
      "ordinal": 1,
      "importance": "required",
      "body": "受け入れ条件は..."
    }
  ]
}
```

### 取得ルール

* `chunk_ids` が指定された場合は最優先
* `heading` がある場合は heading 優先
* `query` がある場合は検索結果上位を返す
* いずれもなければ ordinal 順に先頭から返す

## 7.4 POST `/v1/docs:search`

### 概要

doc メタデータおよび本文に対して検索を行う。

### リクエスト

```json
{
  "query": "memory import acceptance criteria",
  "doc_types": ["spec", "cookbook"],
  "tags": ["memory"],
  "feature_keys": ["memory-import"],
  "limit": 10
}
```

### レスポンス

```json
{
  "results": [
    {
      "doc_id": "doc:spec:memory-import",
      "title": "Memory Import Spec",
      "version": "2026-03-10",
      "score": 0.93,
      "summary": "メモリ取り込み処理の仕様"
    }
  ]
}
```

## 7.5 POST `/v1/reads:ack`

### 概要

読了記録を登録する。

### リクエスト

```json
{
  "task_id": "task:feature:local:123",
  "doc_id": "doc:spec:memory-import",
  "version": "2026-03-10",
  "chunk_ids": [
    "chunk:doc:spec:memory-import:001",
    "chunk:doc:spec:memory-import:002"
  ],
  "reader": "agent"
}
```

### 処理ルール

* `task_id` は必須
* `doc_id` は必須
* `version` は必須
* `chunk_ids` は省略可能だが推奨
* 実際の正本保存先は `agent-taskstate` でもよい
* resolver 側で単独保持する場合でも、将来 `agent-taskstate` に移送可能な形式であること

### レスポンス

```json
{
  "status": "acknowledged",
  "task_id": "task:feature:local:123",
  "doc_id": "doc:spec:memory-import",
  "version": "2026-03-10"
}
```

### 失敗例

* `TASK_NOT_FOUND`
* `DOC_NOT_FOUND`
* `VERSION_REQUIRED`

## 7.6 POST `/v1/docs:stale-check`

### 概要

task に紐づく read receipt と最新 doc version を比較し、stale を判定する。

### リクエスト

```json
{
  "task_id": "task:feature:local:123"
}
```

### レスポンス fresh

```json
{
  "task_id": "task:feature:local:123",
  "status": "fresh",
  "stale_reasons": []
}
```

### レスポンス stale

```json
{
  "task_id": "task:feature:local:123",
  "status": "stale",
  "stale_reasons": [
    {
      "doc_id": "doc:spec:memory-import",
      "previous_version": "2026-03-09",
      "current_version": "2026-03-10",
      "reason": "version_mismatch"
    }
  ]
}
```

### 判定ルール

MVP では以下でよい。

* read receipt の `doc_id`
* read receipt の `version`
* 最新 doc の `version`
* 不一致なら stale

## 7.7 POST `/v1/contracts:resolve`

### 概要

feature または task に対する最低限の契約情報を返す。

### リクエスト

```json
{
  "feature": "memory-import",
  "task_id": "task:feature:local:123"
}
```

### レスポンス

```json
{
  "feature": "memory-import",
  "required_docs": [
    "doc:spec:memory-import",
    "doc:req:memory-core"
  ],
  "acceptance_criteria": [
    "import command succeeds",
    "invalid input is rejected"
  ],
  "forbidden_patterns": [
    "direct write without validation"
  ],
  "definition_of_done": [
    "required docs acknowledged",
    "tests updated"
  ]
}
```

## 8. Skill インターフェース

Skill は API の薄い入口として設計する。

## 8.1 `/resolve-docs`

### 入力

```json
{
  "feature": "memory-import",
  "task_id": "task:feature:local:123"
}
```

### 出力

```json
{
  "required": [
    {
      "doc_id": "doc:spec:memory-import",
      "version": "2026-03-10",
      "top_chunks": [
        "chunk:doc:spec:memory-import:001"
      ]
    }
  ],
  "recommended": []
}
```

## 8.2 `/read-chunks`

### 入力

```json
{
  "doc_id": "doc:spec:memory-import",
  "query": "acceptance criteria",
  "limit": 3
}
```

### 出力

```json
{
  "chunks": [
    {
      "chunk_id": "chunk:doc:spec:memory-import:001",
      "body": "受け入れ条件は..."
    }
  ]
}
```

## 8.3 `/ack-docs`

### 入力

```json
{
  "task_id": "task:feature:local:123",
  "doc_id": "doc:spec:memory-import",
  "version": "2026-03-10",
  "chunk_ids": [
    "chunk:doc:spec:memory-import:001"
  ]
}
```

### 出力

```json
{
  "status": "acknowledged"
}
```

## 8.4 `/stale-check`

### 入力

```json
{
  "task_id": "task:feature:local:123"
}
```

### 出力

```json
{
  "status": "stale",
  "stale_reasons": [
    {
      "doc_id": "doc:spec:memory-import",
      "previous_version": "2026-03-09",
      "current_version": "2026-03-10"
    }
  ]
}
```

## 8.5 `/resolve-contract`

### 入力

```json
{
  "feature": "memory-import"
}
```

### 出力

```json
{
  "required_docs": [
    "doc:spec:memory-import"
  ],
  "acceptance_criteria": [
    "import command succeeds"
  ],
  "forbidden_patterns": [
    "direct write without validation"
  ]
}
```

## 9. `agent-taskstate` 連携インターフェース

resolver は task 正本を持たない。
必要な状態は `agent-taskstate` に委譲可能であること。

## 9.1 read receipt 連携

### resolver から taskstate へ渡す想定データ

```json
{
  "task_id": "task:feature:local:123",
  "doc_id": "doc:spec:memory-import",
  "version": "2026-03-10",
  "chunk_ids": [
    "chunk:doc:spec:memory-import:001"
  ],
  "read_at": "2026-03-10T10:00:00+09:00",
  "reader": "agent"
}
```

## 9.2 task dependency 連携

### 例

```json
{
  "task_id": "task:feature:local:123",
  "ref_type": "doc",
  "ref_id": "doc:spec:memory-import",
  "expected_version": "2026-03-10"
}
```

## 9.3 stale reason 連携

### 例

```json
{
  "task_id": "task:feature:local:123",
  "ref_id": "doc:spec:memory-import",
  "previous_version": "2026-03-09",
  "current_version": "2026-03-10",
  "detected_at": "2026-03-10T10:10:00+09:00"
}
```

## 10. chunking 仕様

### 10.1 基本戦略

* 見出し単位を優先
* 長すぎる節のみ再分割
* summary chunk を任意で追加
* chunk 単独で意味が成立するようにする

### 10.2 mode

許可する `chunking.mode`

* `heading`
* `fixed`

MVP では `heading` を標準とする。

### 10.3 見出しパース

Markdown の見出し階層を `heading_path` として保持する。

例

```json
["BLUEPRINT", "Task Lifecycle", "Acceptance"]
```

## 11. 認可と安全性

MVP では高度な認可機構は必須としないが、以下を満たすこと。

* destructive operation を含めない
* 読み取り中心の API とする
* 文書登録は信頼できる実行経路からのみ行う
* read receipt の改ざん防止は将来課題とする

## 12. ログと監査

最低限以下を記録可能であること。

* いつ文書を登録したか
* いつ resolve したか
* どの task に対して何を ack したか
* stale がいつ検知されたか

推奨ログ項目

* timestamp
* operation
* actor
* target_id
* result

## 13. 互換性

* Markdown 原本を維持できること
* `memx-core` 単体でも最低限動くこと
* `agent-taskstate` がある場合は read receipt と stale を委譲できること
* 将来 `tracker-bridge-materials` とつなげられること

## 14. MVP 判定基準

以下を満たしたら interfaces MVP として成立とする。

* doc ingest ができる
* docs resolve ができる
* chunks get ができる
* reads ack ができる
* stale check ができる
* Skill 入口から同等操作ができる

## 15. 実装メモ

推奨初期実装順

* `POST /v1/docs:ingest`
* `POST /v1/docs:resolve`
* `POST /v1/chunks:get`
* `/resolve-docs`
* `/read-chunks`
* `POST /v1/reads:ack`
* `POST /v1/docs:stale-check`

## 16. 結論

本インターフェース群は、Markdown 群をそのままエージェントに読ませるのではなく、機能や task を入口にして必要文書と必要 chunk を解決し、その参照結果を契約として残すためのものである。

Skill は見た目の入口であり、本体は API 契約と状態保持の連携面にある。