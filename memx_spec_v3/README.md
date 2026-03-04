# memx_spec v3

要求/仕様/設計/I/F 契約と、CLI→API→Service→DB の最小実装をまとめたディレクトリです。

## 構成

## 設計書作成開始入口（IA）
- 設計着手時の単一入口: [docs/design-doc-ia-spec.md](./docs/design-doc-ia-spec.md)

正本は `docs/spec.md` の役割分担に従って参照してください。

- `docs/spec.md`
  仕様インデックス（正本/補助の役割分担、API/CLI/エラー/NFR/運用要件の参照導線）。
- `docs/requirements.md`
  システム全体の要件定義の正本（目的・アーキテクチャ・CLI・GC・LLM役割など、レビュー内容を反映）。
- `docs/contracts/openapi.yaml`
  API 契約の正本。
- `docs/contracts/cli-json.schema.json`
  CLI `--json` 契約の正本。
- `docs/traceability.md`
  要求・契約・実装の対応関係（トレーサビリティ）の正本。
- `docs/design.md`
  設計（レイヤ構成、DB責務分割、移行戦略）。
- `docs/interfaces.md`
  CLI/API I/O、互換ルール、エラー面。
- `docs/CONTRACTS.md`
  API/CLI の機械可読契約一覧（フィールド単位）。
- `schema/short.sql`
  `short.db` 用の CREATE TABLE スキーマ（FTS5 トリガ修正・カラム制約調整・user_version 追加済み）。

### Go（モジュールは `go/` 配下）

- `go/go.mod`
  - Go モジュール定義（依存は最小）。
- 運用ルール
  - `go.mod` を変更した場合は `go.sum` を同時更新し、同一コミットで管理する。
  - 依存ロック監査ログ: [`docs/reviews/GO-LOCK-AUDIT-20260304.md`](./docs/reviews/GO-LOCK-AUDIT-20260304.md)
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
