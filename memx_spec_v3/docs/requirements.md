---
owner: memx-core
status: active
last_reviewed_at: 2026-03-06
next_review_due: 2026-06-06
priority: high
---

# memx 要求事項（MUST / SHOULD / FUTURE）

> **注意**: 本書は要件の索引・概要を記載します。詳細は分割ファイルを参照してください。

## 仕様文書の正本範囲

- 本書（`requirements.md`）は memx v1.3 の要件正本の索引です。
- API 契約の正本は `contracts/openapi.yaml`、CLI `--json` 契約の正本は `contracts/cli-json.schema.json` とする。
- `error-contract.md` と `quickstart.md` は本書に従属する補助文書です。
- 参照起点は `spec.md` とします。

## 0. 目的とスコープ

設計詳細は `design.md`、I/F 詳細は `interfaces.md` を参照し、契約詳細は正本の `contracts/openapi.yaml` と `contracts/cli-json.schema.json` を参照してください。

## ADR参照運用ルール

- 要件に影響する設計判断を追加/変更する場合、該当節へ ADR リンクを追記してください。
- `design.md` と本書の該当節リンクは同一PRで同期更新してください。
- ADR は `docs/ADR/README.md` 索引更新を必須とし、未更新のままマージしません。

## 1. MUST（v1）

- CLI: `mem in short` / `mem out search` / `mem out show` を提供する。
- API: `POST /v1/notes:ingest` / `POST /v1/notes:search` / `GET /v1/notes/{id}` を提供する。
- CLI `--json` は API レスポンスと同型を維持する。
- 入力不備は 400 系、内部障害は 500 系で返す。
- fail-closed 方針で機密入力を拒否できること。
- ingest/search/show の最小性能目標を満たすこと。

## 2. SHOULD（v1.x）

- `POST /v1/gc:run` は「公開可否」と「実行可否」を分離して扱う。
- `mem.features.gc_short=true` 時のみ `mem gc short` / `POST /v1/gc:run` の実行を有効化する。
- SHOULD 機能は feature flag 既定 OFF で提供し、既定挙動を壊さない。

## 3. FUTURE（v2+）

- Recall/Working/Tag/Meta/Lineage/Distill 系 CLI/API の正式導入。
- 破壊変更を伴う再設計（段階移行前提）。

## 4. ID 定義

主要要件ID:

| 領域 | ID | 概要 |
| --- | --- | --- |
| CLI | `REQ-CLI-001` | v1必須3コマンド |
| API | `REQ-API-001` | v1必須3エンドポイント |
| GC | `REQ-GC-001` | GC dry-run |
| Security | `REQ-SEC-001` | fail-closed |
| Error | `REQ-ERR-001` | エラーモデル |
| NFR | `REQ-NFR-001`〜`REQ-NFR-006` | 非機能要件 |

---

## 分割ファイル一覧

| ファイル | 内容 | 元セクション |
| --- | --- | --- |
| [requirements-release.md](./requirements-release.md) | Release Scope Matrix / バージョニング / トレーサビリティ | 0-1, 0-2, 0-3 |
| [requirements-architecture.md](./requirements-architecture.md) | 全体アーキテクチャ / ストア構成 / Gatekeeper | 1-0 〜 1-5 |
| [requirements-data-model.md](./requirements-data-model.md) | データモデル / Security & Retention | 2-1 〜 2-7 |
| [requirements-cli.md](./requirements-cli.md) | CLI 要件 | 3-1 〜 3-6 |
| [requirements-llm.md](./requirements-llm.md) | LLM 戦略 | 4-1 〜 4-4 |
| [requirements-nfr.md](./requirements-nfr.md) | 非機能要件 | 5-1 〜 5-4 |
| [requirements-api.md](./requirements-api.md) | API 要件 | 6-1 〜 6-5 |
| [requirements-incident.md](./requirements-incident.md) | インシデント対応要件 | 11 |

---

## クイックリファレンス

### 主要要件ID一覧

| ID | 説明 | 詳細 |
| --- | --- | --- |
| `REQ-CLI-001` | CLI v1必須3コマンド | [requirements-cli.md](./requirements-cli.md) |
| `REQ-API-001` | API v1必須3エンドポイント | [requirements-api.md](./requirements-api.md) |
| `REQ-GC-001` | GC dry-run | [requirements-cli.md#3-5](./requirements-cli.md#3-5-mem-gc-shortobserver--reflector) |
| `REQ-SEC-001` | fail-closed security | [requirements-data-model.md#2-7](./requirements-data-model.md#2-7-security--retention-requirements) |
| `REQ-ERR-001` | エラーモデル | [requirements-api.md#6-4](./requirements-api.md#6-4-エラーモデル) |
| `REQ-NFR-001` | 性能目標 | [requirements-nfr.md#5-1](./requirements-nfr.md#5-1-性能目標v1必須3エンドポイント) |