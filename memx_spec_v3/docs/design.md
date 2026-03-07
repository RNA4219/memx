---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 設計（design）

## 1. レイヤ構成

### Objective
- CLI/API/Service/Infra の責務境界を固定し、設計変更時の影響範囲を限定する。

### Source
- `memx_spec_v3/docs/design.md#1. レイヤ構成`
- `memx_spec_v3/docs/requirements.md#1-2-2-v1-スコープ境界（normative）`

### Node IDs
- `design`
  - depends_on: `requirements`

### Requirements
- `REQ-CLI-001`
- `REQ-API-001`
- `REQ-ERR-001`

### Commands
- 要件ID網羅: `rg -n "REQ-CLI-001|REQ-API-001|REQ-ERR-001" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "mem in short|POST /v1/notes:ingest|ErrorCode" memx_spec_v3/docs/design.md memx_spec_v3/docs/contracts/openapi.yaml memx_spec_v3/docs/contracts/cli-json.schema.json`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `memx_spec_v3/docs/contracts/openapi.yaml`
- `memx_spec_v3/docs/contracts/cli-json.schema.json`

### Status
- active

```
CLI -> API -> Service(Usecase) -> DB / LLM / Gatekeeper
```
- CLI: 入力整形と表示のみを担当。
- API: 安定 JSON I/F を提供。
- Service: ビジネスロジックの唯一入口。
- DB/LLM/Gatekeeper: 副作用を持つインフラ層。

## 2. DB 責務分割

### Objective
- 4DB 分割の責務境界を明示し、保存先と更新規則を一貫させる。

### Source
- `memx_spec_v3/docs/design.md#2. DB 責務分割`
- `docs/ADR/ADR-0001-4db-boundary.md#ADR-0001: 4DB分割と責務境界`
- `memx_spec_v3/docs/requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive`

### Node IDs
- `design`
  - depends_on: `requirements`

### Requirements
- `REQ-STORE-SHORT-001`
- `REQ-STORE-SHORT-002`
- `REQ-STORE-CHR-001`
- `REQ-STORE-CHR-002`
- `REQ-STORE-MP-001`
- `REQ-STORE-MP-002`
- `REQ-STORE-ARC-001`
- `REQ-STORE-ARC-002`

### Commands
- 要件ID網羅: `rg -n "REQ-STORE-(SHORT|CHR|MP|ARC)-00[12]" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "short|chronicle|memopedia|archive" memx_spec_v3/docs/design.md memx_spec_v3/docs/interfaces.md`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `docs/ADR/ADR-0001-4db-boundary.md`
- `memx_spec_v3/docs/interfaces.md`

### Status
- active

- ADR: [ADR-0001: 4DB分割と責務境界](../../docs/ADR/ADR-0001-4db-boundary.md)
- `short.db`: 一次投入先。短期メモ、GC 対象の起点。
- `chronicle.db`: 時系列ログ（出来事・進捗）。
- `memopedia.db`: 抽象知識（定義・方針）。
- `archive.db`: 退避保管（通常検索対象外）。

共通責務:
- `notes`, `tags`, `note_tags`, `note_embeddings`, `notes_fts`（archive は一部省略可）。

short 固有:
- `short_meta`: GC 判定メタ。
- `lineage`: 蒸留/昇格/退避の系譜。

## 2.1 store別設計詳細（short/chronicle/memopedia/archive）

### Objective
- store ごとの必須要件・禁止変更・許可変更を固定し、後方互換を維持する。

### Source
- `memx_spec_v3/docs/design.md#2.1 store別設計詳細（short/chronicle/memopedia/archive）`
- `memx_spec_v3/docs/requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive`

### Node IDs
- `design`
  - depends_on: `requirements`

### Requirements
- `REQ-STORE-SHORT-001`
- `REQ-STORE-SHORT-002`
- `REQ-STORE-CHR-001`
- `REQ-STORE-CHR-002`
- `REQ-STORE-MP-001`
- `REQ-STORE-MP-002`
- `REQ-STORE-ARC-001`
- `REQ-STORE-ARC-002`

### Commands
- 要件ID網羅: `rg -n "REQ-STORE-(SHORT|CHR|MP|ARC)-00[12]" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "short.notes|chronicle.notes|memopedia.notes|archive.notes" memx_spec_v3/docs/design.md memx_spec_v3/docs/interfaces.md`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `memx_spec_v3/docs/interfaces.md`
- `memx_spec_v3/docs/contracts/openapi.yaml`
- `memx_spec_v3/docs/contracts/cli-json.schema.json`

### Status
- active

| store | Requirement ID | DB責務 | 禁止変更 | 許可変更 |
| --- | --- | --- | --- | --- |
| short | [`REQ-STORE-SHORT-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-SHORT-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 一次投入先として `short.notes` 保存、GC dry-run の候補抽出（非更新） | 必須保存項目の破壊、既定挙動での自動削除、`--json` 出力キー破壊 | 任意列追加、任意 CLI オプション追加、dry-run 診断項目追加 |
| chronicle | [`REQ-STORE-CHR-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-CHR-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 時系列ログを `working_scope` 必須で保持し、期間/タグ抽出に応答 | `working_scope` 必須性撤回、既定ソート反転、既存 ID 体系変更 | 非破壊インデックス追加、任意フィルタ追加、任意メタ列追加 |
| memopedia | [`REQ-STORE-MP-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-MP-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 用語/方針ノート保持、reflect 時の版管理保持 | 本文自動上書き、同一 ID 再利用、必須列削除 | 任意セクション追加、参照メタ追加、互換な検索キー追加 |
| archive | [`REQ-STORE-ARC-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-ARC-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 退避ノート保持、lineage 記録後のみ short 削除、purge dry-run 候補提示 | 監査ログなし purge、retention 無視削除、lineage 未記録削除 | 保持メタ列追加、監査ログ項目追加、dry-run 出力列追加 |

## 2.2 Security/Retention 設計

### Objective
- fail-closed 判定、監査証跡、retention 削除条件を要件準拠で固定する。

### Source
- `memx_spec_v3/docs/design.md#2.2 Security/Retention 設計`
- `memx_spec_v3/docs/requirements.md#2-7-security--retention-requirements`
- `memx_spec_v3/docs/requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5`
- `memx_spec_v3/docs/requirements.md#2-7-5-guardrails-fail-closed-との整合チェック要件`

### Node IDs
- `design`
  - depends_on: `requirements`
- `guardrails`
  - depends_on: `requirements`, `governance_policy`

### Requirements
- `REQ-SEC-001`
- `REQ-RET-001`
- `REQ-SEC-AUD-001`
- `REQ-SEC-AUD-002`
- `REQ-SEC-GRD-001`

### Commands
- 要件ID網羅: `rg -n "REQ-(SEC-001|RET-001|SEC-AUD-001|SEC-AUD-002|SEC-GRD-001)" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "POLICY_DENIED|archive_move|archive_purge" memx_spec_v3/docs/design.md memx_spec_v3/docs/contracts/openapi.yaml GUARDRAILS.md`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `GUARDRAILS.md`
- `docs/security/minimal_operations.md`
- `RUNBOOK.md`

### Status
- active

| Requirement ID | 保存可否 | 監査証跡 | 保持期間 | 削除条件 |
| --- | --- | --- | --- | --- |
| [`REQ-SEC-001`](./requirements.md#2-7-security--retention-requirements) | `sensitivity=secret` は fail-closed で保存禁止 | 拒否時に判定理由と actor を監査ログに残す | ポリシー判定ログは運用監査期間に従う | 保存禁止のためデータ削除フロー対象外 |
| [`REQ-RET-001`](./requirements.md#2-7-security--retention-requirements) | `archive` 退避後のみ保持継続 | 退避/削除の実行結果を `archive_move` / `archive_purge` で記録 | 期限は `retention_until` 基準で管理 | 期限超過かつ承認済み、監査項目充足時のみ削除 |
| [`REQ-SEC-AUD-001`](./requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5) | `archive_move` 実行時のみ保存状態遷移許可 | 固定フィールド（actor/approval/evidence path など）を必須保存 | 監査証跡は改ざん不可で保持 | 必須監査キー欠落時は遷移拒否 |
| [`REQ-SEC-AUD-002`](./requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5) | `archive_purge` は承認フロー通過時のみ許可 | purge 監査ログ固定項目を必須保存 | purge 証跡は retention 監査期間満了まで保持 | 監査コンテキスト不足時は削除禁止 |

## 3. 移行戦略

### Objective
- スキーマ移行の後方互換方針と feature flag 導入境界を固定する。

### Source
- `memx_spec_v3/docs/design.md#3. 移行戦略`
- `memx_spec_v3/docs/requirements.md#6-2-v1-互換性方針`

### Node IDs
- `design`
  - depends_on: `requirements`

### Requirements
- `REQ-API-001`
- `REQ-ERR-001`

### Commands
- 要件ID網羅: `rg -n "REQ-API-001|REQ-ERR-001" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "user_version|feature flag|後方互換" memx_spec_v3/docs/design.md memx_spec_v3/docs/contracts/openapi.yaml`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `schema/*.sql`
- `memx_spec_v3/docs/contracts/openapi.yaml`

### Status
- active

- マイグレーションは `schema/*.sql` を正本として適用する。
- `PRAGMA user_version` を採用し、破壊的/非互換 DDL のみバージョンを進める。
- v1 では後方互換を最優先し、破壊変更は v2+（`FUTURE`）へ隔離する。
- 実験機能は feature flag 既定 OFF で段階導入する。

## 4. ユースケース設計

### Objective
- v1 必須ユースケース（ingest/search/show/gc dry-run）のI/O契約と失敗分岐を固定する。

### Source
- `memx_spec_v3/docs/design.md#4. ユースケース設計`
- `memx_spec_v3/docs/requirements.md#3-cli-要件`
- `memx_spec_v3/docs/requirements.md#6-api-要件v13-追加`
- `memx_spec_v3/docs/requirements.md#6-4-エラーモデル`

### Node IDs
- `design`
  - depends_on: `requirements`
- `api`
  - depends_on: `service`, `quickstart`

### Requirements
- `REQ-CLI-001`
- `REQ-API-001`
- `REQ-GC-001`
- `REQ-ERR-001`
- `REQ-SEC-001`

### Commands
- 要件ID網羅: `rg -n "REQ-(CLI-001|API-001|GC-001|ERR-001|SEC-001)" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "POST /v1/notes:ingest|POST /v1/notes:search|GET /v1/notes/\{id\}|POST /v1/gc:run" memx_spec_v3/docs/design.md memx_spec_v3/docs/contracts/openapi.yaml`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `memx_spec_v3/docs/contracts/openapi.yaml`
- `memx_spec_v3/docs/contracts/cli-json.schema.json`
- `memx_spec_v3/docs/interfaces.md`

### Status
- active

### 4.1 Ingest
- 入力/出力
  - 入力: CLI `mem in short --content <text> [--json]` / API `POST /v1/notes:ingest`
  - 出力: `id`, `store=short`, `created_at`（CLI `--json` は API と同型）
- レイヤ通過順（CLI/API/Service/DB/LLM/Gatekeeper）
  - CLI（入力整形） → API（契約検証） → Service（ingest実行） → Gatekeeper（機密/ポリシー判定） → DB（shortへ保存）
- 失敗分岐（INVALID_ARGUMENT / POLICY_DENIED / INTERNAL）
  - `INVALID_ARGUMENT`: `content` 空/不正、title/body長制限超過、enum値不正
  - `POLICY_DENIED`: Gatekeeper が fail-closed で拒否（`sensitivity=secret` 等）
  - `NEEDS_HUMAN`: Gatekeeper が要人間判定（v1.3 ではエラーとして扱う）
  - `INTERNAL`: DB 書き込み失敗、予期しない実行時エラー
- 再試行可否（retryable true/false）
  - `INVALID_ARGUMENT`: `false`
  - `POLICY_DENIED`: `false`
  - `NEEDS_HUMAN`: `false`
  - `INTERNAL`: `true`

### 4.2 Search
- 入力/出力
  - 入力: CLI `mem out search --query <text> [--json]` / API `POST /v1/notes:search`
  - 出力: `items[]`（`id`,`content`,`score`,`store`）, `next_cursor`（任意）
- レイヤ通過順（CLI/API/Service/DB/LLM/Gatekeeper）
  - CLI（入力整形） → API（契約検証） → Service（検索戦略決定） → Gatekeeper（検索可否判定） → DB（FTS/メタ検索） → LLM（要約/再ランキングが有効時のみ）
- 失敗分岐（INVALID_ARGUMENT / POLICY_DENIED / INTERNAL）
  - `INVALID_ARGUMENT`: `query` 空/不正
  - `POLICY_DENIED`: Gatekeeper が検索要求を拒否
  - `INTERNAL`: DB クエリエラー、LLM 呼び出し失敗（有効時）
- 再試行可否（retryable true/false）
  - `INVALID_ARGUMENT`: `false`
  - `POLICY_DENIED`: `false`
  - `INTERNAL`: `true`

### 4.3 Show
- 入力/出力
  - 入力: CLI `mem out show --id <note_id> [--json]` / API `GET /v1/notes/{id}`
  - 出力: `id`,`content`,`store`,`created_at`,`updated_at`,`tags[]`
- レイヤ通過順（CLI/API/Service/DB/LLM/Gatekeeper）
  - CLI（入力整形） → API（契約検証） → Service（取得ロジック） → Gatekeeper（閲覧可否判定） → DB（主キー参照）
- 失敗分岐（INVALID_ARGUMENT / POLICY_DENIED / INTERNAL）
  - `INVALID_ARGUMENT`: `id` 形式不正
  - `POLICY_DENIED`: Gatekeeper が閲覧を拒否
  - `INTERNAL`: DB 参照失敗、想定外エラー
- 再試行可否（retryable true/false）
  - `INVALID_ARGUMENT`: `false`
  - `POLICY_DENIED`: `false`
  - `INTERNAL`: `true`

### 4.4 GC dry-run
- 入力/出力
  - 入力: CLI `mem gc short --dry-run --threshold <n>` / API `POST /v1/gc:run`（`dry_run=true`）
  - 出力: `candidates[]`,`summary`,`dry_run=true`（DB 非更新）
- レイヤ通過順（CLI/API/Service/DB/LLM/Gatekeeper）
  - CLI（入力整形） → API（契約検証） → Service（GC 判定） → Gatekeeper（実行可否判定） → DB（候補抽出のみ、更新禁止）
- 失敗分岐（INVALID_ARGUMENT / POLICY_DENIED / INTERNAL）
  - `INVALID_ARGUMENT`: `threshold` 範囲不正、`dry_run` 指定不整合
  - `POLICY_DENIED`: Gatekeeper または feature flag 条件で拒否
  - `INTERNAL`: DB 読み取り失敗、想定外エラー
- 再試行可否（retryable true/false）
  - `INVALID_ARGUMENT`: `false`
  - `POLICY_DENIED`: `false`
  - `INTERNAL`: `true`

**REQ 対応表**

| REQ ID | リンク |
| --- | --- |
| `REQ-CLI-001` | [requirements.md#3-cli-要件](./requirements.md#3-cli-要件) |
| `REQ-API-001` | [requirements.md#6-api-要件v13-追加](./requirements.md#6-api-要件v13-追加) |
| `REQ-GC-001` | [requirements.md#3-5-mem-gc-shortobserver--reflector](./requirements.md#3-5-mem-gc-shortobserver--reflector) |
| `REQ-ERR-001` | [requirements.md#6-4-エラーモデル](./requirements.md#6-4-エラーモデル) |
| `REQ-SEC-001` | [requirements.md#2-7-security--retention-requirements](./requirements.md#2-7-security--retention-requirements) |

## 5. ADR参照運用ルール

### Objective
- 設計判断の固定を ADR 起点で管理し、requirements/design の同時更新を徹底する。

### Source
- `memx_spec_v3/docs/design.md#5. ADR参照運用ルール`
- `docs/ADR/ADR-0002-v1-required-endpoints.md#ADR-0002: v1必須3エンドポイント`
- `docs/ADR/ADR-0003-errorcode-retryable-boundary.md#ADR-0003: ErrorCode/retryable 境界`

### Node IDs
- `design`
  - depends_on: `requirements`

### Requirements
- `REQ-API-001`
- `REQ-ERR-001`

### Commands
- 要件ID網羅: `rg -n "REQ-API-001|REQ-ERR-001" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "ADR-0002|ADR-0003|ErrorCode|retryable" memx_spec_v3/docs/design.md docs/ADR/ADR-0002-v1-required-endpoints.md docs/ADR/ADR-0003-errorcode-retryable-boundary.md`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `docs/ADR/ADR-0002-v1-required-endpoints.md`
- `docs/ADR/ADR-0003-errorcode-retryable-boundary.md`

### Status
- active

- 本書で設計判断を追加/変更する場合は、該当節に ADR リンクを追記する。
- ADR 未作成で判断を固定しない。最小でも `Context / Decision / Consequences / Status / Date` を満たす ADR を先に作成する。
- 本書の該当節リンクと `requirements.md` の対応節リンクは同一PRで更新する。
- v1必須3エンドポイント関連は [ADR-0002](../../docs/ADR/ADR-0002-v1-required-endpoints.md) を参照する。
- ErrorCode / retryable 境界は [ADR-0003](../../docs/ADR/ADR-0003-errorcode-retryable-boundary.md) を参照する。

## 6. 設計→契約→検証 導線（要件ID単位）

### Objective
- 要件IDごとに design / contracts / evaluation の導線を固定し、追跡可能性を維持する。

### Source
- `memx_spec_v3/docs/design.md#6. 設計→契約→検証 導線（要件ID単位）`
- `memx_spec_v3/docs/requirements.md#1-2-4-要件トレーサビリティ（normative）`
- `RUNBOOK.md#traceability`
- `EVALUATION.md#req-gates`

### Node IDs
- `design`
  - depends_on: `requirements`
- `evaluation`
  - depends_on: `runbook`, `checklists`

### Requirements
- `REQ-CLI-001`
- `REQ-API-001`
- `REQ-GC-001`
- `REQ-ERR-001`
- `REQ-SEC-001`
- `REQ-NFR-001`
- `REQ-NFR-002`
- `REQ-NFR-005`

### Commands
- 要件ID網羅: `rg -n "REQ-(CLI-001|API-001|GC-001|ERR-001|SEC-001|NFR-001|NFR-002|NFR-005)" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "POST /v1/notes:ingest|POST /v1/notes:search|GET /v1/notes/\{id\}|POST /v1/gc:run|ErrorCode" memx_spec_v3/docs/design.md memx_spec_v3/docs/contracts/openapi.yaml`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `RUNBOOK.md`
- `EVALUATION.md`
- `memx_spec_v3/docs/contracts/openapi.yaml`

### Status
- active

| Requirement ID | Design Section | Interface ID | Evaluation項目 |
| --- | --- | --- | --- |
| [`REQ-CLI-001`](./requirements.md#3-cli-要件) | 4.1 / 4.2 / 4.3 | CLI `mem in short`, `mem out search`, `mem out show` | CLI `--json` と API 同型、必須項目一致 |
| [`REQ-API-001`](./requirements.md#6-api-要件v13-追加) | 4.1 / 4.2 / 4.3 | `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}` | HTTP ステータス、レスポンス型、互換性 |
| [`REQ-GC-001`](./requirements.md#3-5-mem-gc-shortobserver--reflector) | 4.4 | `mem gc short --dry-run`, `POST /v1/gc:run` | dry-run DB 非更新、候補件数/対象ID整合 |
| [`REQ-ERR-001`](./requirements.md#6-4-エラーモデル) | 4.1〜4.4 | 共通 ErrorCode 契約 | `INVALID_ARGUMENT` / `POLICY_DENIED` / `INTERNAL` の retryable 整合 |
| [`REQ-SEC-001`](./requirements.md#2-7-security--retention-requirements) | 2.2, 4.1〜4.4 | Gatekeeper 判定（ingest/search/show/gc） | fail-closed 拒否、監査ログ記録 |
| [`REQ-NFR-001`](./requirements.md#5-1-性能目標v1必須3エンドポイント) | 6.1 | ingest/search/show | p95 閾値の達成（計測プロトコル準拠） |
| [`REQ-NFR-002`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | 6.1 | 運用復旧フロー | RTO/RPO 同時成立 |
| [`REQ-NFR-005`](./requirements.md#5-3-整合性回復要件archive-補償フロー) | 6.1 | short→archive 補償 | 30分以内収束または IN 起票 |
| [`REQ-NFR-006`](./requirements.md#5-4-インシデント記録docsin-md最小監査項目) | 6.2 | インシデント監査記録/waiver 記録 | 必須監査項目 + waiver 項目欠落なし |

## 6.1 NFR設計（性能 / 復旧 / 整合性回復）

### Objective
- NFR の計測入力・判定方法・失敗時運用を要件ID単位で明確化する。

### Source
- `memx_spec_v3/docs/design.md#6.1 NFR設計（性能 / 復旧 / 整合性回復）`
- `memx_spec_v3/docs/requirements.md#5-1-性能目標v1必須3エンドポイント`
- `memx_spec_v3/docs/requirements.md#5-2-可用性復旧整合性回復運用nfr`
- `memx_spec_v3/docs/requirements.md#5-3-整合性回復要件archive-補償フロー`

### Node IDs
- `design`
  - depends_on: `requirements`
- `runbook`
  - depends_on: `requirements`, `guardrails`, `governance_prioritization`
- `evaluation`
  - depends_on: `runbook`, `checklists`

### Requirements
- `REQ-NFR-001`
- `REQ-NFR-002`
- `REQ-NFR-003`
- `REQ-NFR-004`
- `REQ-NFR-005`
- `REQ-NFR-006`

### Commands
- 要件ID網羅: `rg -n "REQ-NFR-00[1-6]" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md`
- 契約同期: `rg -n "p95|RTO|RPO|30分以内|15分以内|再処理" memx_spec_v3/docs/design.md RUNBOOK.md EVALUATION.md`
- リンク健全性: `python scripts/check_links.py memx_spec_v3/docs/design.md`

### Dependencies
- `RUNBOOK.md`
- `EVALUATION.md`
- `docs/IN-*.md`

### Status
- active

| Requirement ID | 判定入力（ログ/成果物） | 判定方法 |
| --- | --- | --- |
| [`REQ-NFR-001`](./requirements.md#5-1-性能目標v1必須3エンドポイント) | `RUNBOOK.md` の trace 実行ログ、性能計測結果（p95） | ingest/search/show の 3 エンドポイントのみ対象に閾値比較 |
| [`REQ-NFR-002`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | インシデント記録（`docs/IN-*.md`）、運用タイムライン | `rto_minutes <= 30` かつ `rpo_minutes <= 5` を同時充足 |
| [`REQ-NFR-003`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | 検知時刻 / 暫定復旧時刻ログ | 検知〜暫定復旧が 15 分以内 |
| [`REQ-NFR-004`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | 再試行履歴、ジョブログ | 1 リクエストあたり再処理 2 回以内 |
| [`REQ-NFR-005`](./requirements.md#5-3-整合性回復要件archive-補償フロー) | 補償フロー実行ログ、`docs/IN-*.md` 起票有無 | 30分以内収束、未収束時は IN 起票 |

## 6.2 REQ-NFR-006（監査記録 / waiver）責務境界

### Objective
- `REQ-NFR-006` の監査記録と waiver 要件について、設計責務（記録対象/参照導線）と運用責務（起票/承認/維持）の境界を固定する。

### Source
- `memx_spec_v3/docs/design.md#6-2-req-nfr-006監査記録--waiver-責務境界`
- `memx_spec_v3/docs/requirements.md#5-4-インシデント記録docsin-md最小監査項目`
- `memx_spec_v3/docs/requirements.md#5-4-1-waiver-時の必須記録docsin-md-運用連動`
- `memx_spec_v3/docs/operations-spec.md#2-waiver-記録必須項目req-nfr-006`
- `memx_spec_v3/docs/operations-spec.md#6-必須証跡ファイル一覧とキー定義`

### Node IDs
- `design`
  - depends_on: `requirements`, `traceability`
- `operations`
  - depends_on: `requirements`, `runbook`
- `evidence`
  - depends_on: `operations`

### Requirements
- `REQ-NFR-006`

### Commands
- 要件ID網羅: `rg -n "REQ-NFR-006" memx_spec_v3/docs/design.md memx_spec_v3/docs/requirements.md memx_spec_v3/docs/traceability.md EVALUATION.md`
- 証跡導線確認: `rg -n "docs/IN-\*\.md|artifacts/ops/incident-summary.json|artifacts/ops/recovery-log.ndjson" memx_spec_v3/docs/design.md memx_spec_v3/docs/operations-spec.md`

### Dependencies
- `memx_spec_v3/docs/operations-spec.md`
- `EVALUATION.md`
- `docs/IN-*.md`
- `artifacts/ops/incident-summary.json`
- `artifacts/ops/recovery-log.ndjson`

### Status
- active

| 責務層 | 固定責務 | 非責務（運用へ委譲） | 参照固定先 |
| --- | --- | --- | --- |
| 設計（本書） | 必須監査項目と waiver 項目のデータ境界を定義し、評価/証跡への参照導線を固定する。 | 個別インシデントの起票判断、承認実施、期限更新の実務判断。 | `memx_spec_v3/docs/operations-spec.md#2-waiver-記録必須項目req-nfr-006`, `memx_spec_v3/docs/operations-spec.md#6-必須証跡ファイル一覧とキー定義` |
| 運用（operations-spec / RUNBOOK） | `docs/IN-<実日付>-<連番>.md` 起票、waiver 7 項目の記録・期限管理、証跡ファイル更新を実行する。 | 要件自体の閾値変更や判定基準の再定義。 | `memx_spec_v3/docs/operations-spec.md#7-waiver-運用発動条件期限解除条件未解除時エスカレーション`, `RUNBOOK.md#障害時手順要件id紐付け` |
| 評価（EVALUATION） | `REQ-NFR-006` 合否を必須監査項目欠落/waiver期限切れ/再計画チケット欠落で機械判定する。 | 証跡の作成・補完作業。 | `EVALUATION.md#運用nfr可用性復旧整合性回復合否基準` |

## 7. design-template 段階移行チェックリスト（章単位）

- [x] 1. レイヤ構成 を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 2. DB 責務分割 を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 2.1 store別設計詳細（short/chronicle/memopedia/archive） を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 2.2 Security/Retention 設計 を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 3. 移行戦略 を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 4. ユースケース設計 を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 5. ADR参照運用ルール を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 6. 設計→契約→検証 導線（要件ID単位） を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）
- [x] 6.1 NFR設計（性能 / 復旧 / 整合性回復） を `memx_spec_v3/docs/design-template.md` 準拠へ移行（Objective/Source/Node IDs/Requirements/Commands/Dependencies/Status）

## 8. 実装状況（2026-03-06 更新）

### 8.1 マイグレーション層
| ファイル | 状態 | 説明 |
| --- | --- | --- |
| `go/db/migrate_short.go` | 完了 | short.db スキーマ適用 |
| `go/db/migrate_other.go` | 完了 | chronicle/memopedia/archive スキーマ適用 |
| `go/db/open.go` | 完了 | 4DB ATTACH + マイグレーション |

**実装詳細**:
- 各ストアを個別に開いてマイグレーション後、ATTACH する方式
- `PRAGMA user_version` による冪等性保証
- スキーマ構成は `requirements.md#1-2-1` 参照

**スキーマ構成（各ストア共通）**:

| テーブル | short | chronicle | memopedia | archive |
| --- | :---: | :---: | :---: | :---: |
| `notes` | ✅ | ✅ | ✅ | ✅ |
| `notes_fts` (FTS5) | ✅ | ✅ | ✅ | ❌ |
| FTS同期トリガー | ✅ | ✅ | ✅ | ❌ |
| `tags` / `note_tags` | ✅ | ✅ | ✅ | ✅ |
| `note_embeddings` | ✅ | ✅ | ✅ | ❌ |
| `lineage` | ✅ | ✅ | ✅ | ✅ |
| `*_meta` (GC用) | ✅ | ✅ | ✅ | ✅ |
| `working_scope` 列 | ❌ | ✅ | ✅ | ❌ |
| `is_pinned` 列 | ❌ | ✅ | ✅ | ❌ |

**インデックス**（全ストア共通）:
- `idx_notes_created_at`, `idx_notes_last_accessed_at`
- `idx_notes_source_trust`, `idx_notes_sensitivity`
- `idx_tags_name`, `idx_tags_parent`
- `idx_note_tags_tag_id`, `idx_note_tags_note_id`
- `idx_lineage_src`, `idx_lineage_dest`

### 8.2 Gatekeeper 層
| ファイル | 状態 | 説明 |
| --- | --- | --- |
| `go/db/gatekeeper.go` | 完了 | インターフェース定義 |
| `go/db/gatekeeper_impl.go` | 完了 | DefaultGatekeeper 実装 |

**実装詳細**:
- `DefaultGatekeeper`: プロファイルベースのルール判定
- `AllowAllGatekeeper` / `DenyAllGatekeeper`: テスト用ヘルパー
- プロファイル: `STRICT` / `NORMAL` / `DEV`
- fail-closed: `sensitivity=secret` は常に拒否

### 8.3 Service 層
| ファイル | 状態 | 説明 |
| --- | --- | --- |
| `go/service/service.go` | 完了 | IngestShort/SearchShort/GetShort |
| `go/service/errors.go` | 完了 | エラー定義（ErrPolicyDenied, ErrNeedsHuman 追加） |
| `go/service/models.go` | 完了 | Note モデル |

**実装詳細**:
- `IngestShort`: Gatekeeper.Check 呼び出し、入力バリデーション
- 入力バリデーション: title/body 長制限、enum 値チェック

### 8.4 API 層
| ファイル | 状態 | 説明 |
| --- | --- | --- |
| `go/api/types.go` | 完了 | API 型定義 |
| `go/api/errors.go` | 完了 | エラーマッピング（GATEKEEP_DENY 対応） |
| `go/api/http_server.go` | 完了 | HTTP サーバー |
| `go/api/http_client.go` | 完了 | HTTP クライアント |
| `go/api/inproc_client.go` | 完了 | in-proc クライアント |

### 8.5 GC機能
| ファイル | 状態 | 説明 |
| --- | --- | --- |
| `go/service/gc.go` | 完了 | Phase0トリガ判定、Phase3 Archive退避 |
| `go/service/gc_test.go` | 完了 | GC機能テスト |

**実装詳細**:
- Phase0: ノート数とlast_gc_atに基づくトリガ判定
  - soft_limit: 1200ノード (interval経過後実行)
  - hard_limit: 2000ノード (即時実行)
  - min_interval: 180分
- Phase3: アクセス数0で30日以上経過のノートをarchiveへ退避
- dry-run: DB更新せず判定結果のみ返却
- feature flag: `--enable-gc` で有効化（デフォルト無効）

### 8.6 テスト
| パッケージ | テスト数 | 状態 |
| --- | --- | --- |
| `go/db` | 8 | 全て PASS |
| `go/service` | 13 | 全て PASS |
