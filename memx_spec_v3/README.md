# memx_spec v3

要件定義（v1.3）と、CLI→API→Service→DB のレイヤリングを反映した最小実装をまとめた ZIP です。

## 構成

- `docs/requirements.md`  
  システム全体の要件定義（目的・アーキテクチャ・CLI・GC・LLM役割など、レビュー内容を反映）。
- `schema/short.sql`  
  `short.db` 用の CREATE TABLE スキーマ（FTS5 トリガ修正・カラム制約調整・user_version 追加済み）。

### Go（モジュールは `go/` 配下）

- `go/go.mod`
  - Go モジュール定義（依存は最小）。
- `go/db/*`
  - DB 接続・マイグレーション・LLM/Gatekeeper の注入口（インフラ層）。
- `go/service/*`
  - Usecase 層（短期ストアへの ingest/search/get の最小実装）。
- `go/api/*`
  - ツール/AI 向け API（HTTP と in-proc クライアント）。
- `go/cmd/mem/main.go`
  - 人間向け CLI。API の薄いラッパ。

## Status（v1必須 / v1.1以降）

### 必須コマンド

| 区分 | コマンド |
| --- | --- |
| v1必須 | `mem in short` |
| v1必須 | `mem out search` |
| v1必須 | `mem out show` |

### 必須API

| 区分 | API |
| --- | --- |
| v1必須 | `POST /v1/notes:ingest` |
| v1必須 | `POST /v1/notes:search` |
| v1必須 | `GET /v1/notes/{id}` |

### 非対象（v1時点）

| 区分 | 対象外項目 |
| --- | --- |
| v1.1以降 | GC |
| v1.1以降 | recall |
| v1.1以降 | working |
| v1.1以降 | tag |
| v1.1以降 | meta |
| v1.1以降 | lineage |

### 受け入れ条件

| 区分 | 条件 |
| --- | --- |
| v1必須 | 入出力互換（CLI→API の入出力マッピングが保持されること） |
| v1必須 | エラーコード（入力不備: 400系 / 内部障害: 500系 を返すこと） |
| v1必須 | 最小性能目標（`ingest`/`search`/`show` がローカル単体で実用応答時間を維持すること） |
