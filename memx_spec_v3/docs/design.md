---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 設計（design）

## 1. レイヤ構成
```
CLI -> API -> Service(Usecase) -> DB / LLM / Gatekeeper
```
- CLI: 入力整形と表示のみを担当。
- API: 安定 JSON I/F を提供。
- Service: ビジネスロジックの唯一入口。
- DB/LLM/Gatekeeper: 副作用を持つインフラ層。

## 2. DB 責務分割
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

| store | Requirement ID | DB責務 | 禁止変更 | 許可変更 |
| --- | --- | --- | --- | --- |
| short | [`REQ-STORE-SHORT-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-SHORT-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 一次投入先として `short.notes` 保存、GC dry-run の候補抽出（非更新） | 必須保存項目の破壊、既定挙動での自動削除、`--json` 出力キー破壊 | 任意列追加、任意 CLI オプション追加、dry-run 診断項目追加 |
| chronicle | [`REQ-STORE-CHR-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-CHR-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 時系列ログを `working_scope` 必須で保持し、期間/タグ抽出に応答 | `working_scope` 必須性撤回、既定ソート反転、既存 ID 体系変更 | 非破壊インデックス追加、任意フィルタ追加、任意メタ列追加 |
| memopedia | [`REQ-STORE-MP-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-MP-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 用語/方針ノート保持、reflect 時の版管理保持 | 本文自動上書き、同一 ID 再利用、必須列削除 | 任意セクション追加、参照メタ追加、互換な検索キー追加 |
| archive | [`REQ-STORE-ARC-001`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive), [`REQ-STORE-ARC-002`](./requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive) | 退避ノート保持、lineage 記録後のみ short 削除、purge dry-run 候補提示 | 監査ログなし purge、retention 無視削除、lineage 未記録削除 | 保持メタ列追加、監査ログ項目追加、dry-run 出力列追加 |

### Source
- `memx_spec_v3/docs/requirements.md#1-2-3-store-別要求short--chronicle--memopedia--archive`

### Dependencies
- `memx_spec_v3/docs/interfaces.md`
- `memx_spec_v3/docs/contracts/openapi.yaml`
- `memx_spec_v3/docs/contracts/cli-json.schema.json`

## 2.2 Security/Retention 設計

| Requirement ID | 保存可否 | 監査証跡 | 保持期間 | 削除条件 |
| --- | --- | --- | --- | --- |
| [`REQ-SEC-001`](./requirements.md#2-7-security--retention-requirements) | `sensitivity=secret` は fail-closed で保存禁止 | 拒否時に判定理由と actor を監査ログに残す | ポリシー判定ログは運用監査期間に従う | 保存禁止のためデータ削除フロー対象外 |
| [`REQ-RET-001`](./requirements.md#2-7-security--retention-requirements) | `archive` 退避後のみ保持継続 | 退避/削除の実行結果を `archive_move` / `archive_purge` で記録 | 期限は `retention_until` 基準で管理 | 期限超過かつ承認済み、監査項目充足時のみ削除 |
| [`REQ-SEC-AUD-001`](./requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5) | `archive_move` 実行時のみ保存状態遷移許可 | 固定フィールド（actor/approval/evidence path など）を必須保存 | 監査証跡は改ざん不可で保持 | 必須監査キー欠落時は遷移拒否 |
| [`REQ-SEC-AUD-002`](./requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5) | `archive_purge` は承認フロー通過時のみ許可 | purge 監査ログ固定項目を必須保存 | purge 証跡は retention 監査期間満了まで保持 | 監査コンテキスト不足時は削除禁止 |

### Source
- `memx_spec_v3/docs/requirements.md#2-7-security--retention-requirements`
- `memx_spec_v3/docs/requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5`

### Dependencies
- `GUARDRAILS.md`
- `docs/security/minimal_operations.md`
- `RUNBOOK.md`

## 3. 移行戦略
- マイグレーションは `schema/*.sql` を正本として適用する。
- `PRAGMA user_version` を採用し、破壊的/非互換 DDL のみバージョンを進める。
- v1 では後方互換を最優先し、破壊変更は v2+（`FUTURE`）へ隔離する。
- 実験機能は feature flag 既定 OFF で段階導入する。

## 4. ユースケース設計

### 4.1 Ingest
- 入力/出力
  - 入力: CLI `mem in short --content <text> [--json]` / API `POST /v1/notes:ingest`
  - 出力: `id`, `store=short`, `created_at`（CLI `--json` は API と同型）
- レイヤ通過順（CLI/API/Service/DB/LLM/Gatekeeper）
  - CLI（入力整形） → API（契約検証） → Service（ingest実行） → Gatekeeper（機密/ポリシー判定） → DB（shortへ保存）
- 失敗分岐（INVALID_ARGUMENT / POLICY_DENIED / INTERNAL）
  - `INVALID_ARGUMENT`: `content` 空/不正
  - `POLICY_DENIED`: Gatekeeper が fail-closed で拒否
  - `INTERNAL`: DB 書き込み失敗、予期しない実行時エラー
- 再試行可否（retryable true/false）
  - `INVALID_ARGUMENT`: `false`
  - `POLICY_DENIED`: `false`
  - `INTERNAL`: `true`

**REQ 対応表**

| REQ ID | リンク |
| --- | --- |
| `REQ-CLI-001` | [requirements.md#3-cli-要件](./requirements.md#3-cli-要件) |
| `REQ-API-001` | [requirements.md#6-api-要件v13-追加](./requirements.md#6-api-要件v13-追加) |
| `REQ-ERR-001` | [requirements.md#6-4-エラーモデル](./requirements.md#6-4-エラーモデル) |
| `REQ-SEC-001` | [requirements.md#2-7-security--retention-requirements](./requirements.md#2-7-security--retention-requirements) |

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

**REQ 対応表**

| REQ ID | リンク |
| --- | --- |
| `REQ-CLI-001` | [requirements.md#3-cli-要件](./requirements.md#3-cli-要件) |
| `REQ-API-001` | [requirements.md#6-api-要件v13-追加](./requirements.md#6-api-要件v13-追加) |
| `REQ-ERR-001` | [requirements.md#6-4-エラーモデル](./requirements.md#6-4-エラーモデル) |
| `REQ-SEC-001` | [requirements.md#2-7-security--retention-requirements](./requirements.md#2-7-security--retention-requirements) |

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

**REQ 対応表**

| REQ ID | リンク |
| --- | --- |
| `REQ-CLI-001` | [requirements.md#3-cli-要件](./requirements.md#3-cli-要件) |
| `REQ-API-001` | [requirements.md#6-api-要件v13-追加](./requirements.md#6-api-要件v13-追加) |
| `REQ-ERR-001` | [requirements.md#6-4-エラーモデル](./requirements.md#6-4-エラーモデル) |
| `REQ-SEC-001` | [requirements.md#2-7-security--retention-requirements](./requirements.md#2-7-security--retention-requirements) |

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

## 5. 設計→契約→検証 導線（要件ID単位）

| Requirement ID | Design Section | Interface ID | Evaluation項目 |
| --- | --- | --- | --- |
| [`REQ-CLI-001`](./requirements.md#3-cli-要件) | 4.1 / 4.2 / 4.3 | CLI `mem in short`, `mem out search`, `mem out show` | CLI `--json` と API 同型、必須項目一致 |
| [`REQ-API-001`](./requirements.md#6-api-要件v13-追加) | 4.1 / 4.2 / 4.3 | `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}` | HTTP ステータス、レスポンス型、互換性 |
| [`REQ-GC-001`](./requirements.md#3-5-mem-gc-shortobserver--reflector) | 4.4 | `mem gc short --dry-run`, `POST /v1/gc:run` | dry-run DB 非更新、候補件数/対象ID整合 |
| [`REQ-ERR-001`](./requirements.md#6-4-エラーモデル) | 4.1〜4.4 | 共通 ErrorCode 契約 | `INVALID_ARGUMENT` / `POLICY_DENIED` / `INTERNAL` の retryable 整合 |
| [`REQ-SEC-001`](./requirements.md#2-7-security--retention-requirements) | 2.2, 4.1〜4.4 | Gatekeeper 判定（ingest/search/show/gc） | fail-closed 拒否、監査ログ記録 |
| [`REQ-NFR-001`](./requirements.md#5-1-性能目標v1必須3エンドポイント) | 5.1 | ingest/search/show | p95 閾値の達成（計測プロトコル準拠） |
| [`REQ-NFR-002`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | 5.1 | 運用復旧フロー | RTO/RPO 同時成立 |
| [`REQ-NFR-005`](./requirements.md#5-3-整合性回復要件archive-補償フロー) | 5.1 | short→archive 補償 | 30分以内収束または IN 起票 |

## 5.1 NFR設計（性能/復旧/整合性回復）

| Requirement ID | 判定入力（ログ/成果物） | 判定方法 |
| --- | --- | --- |
| [`REQ-NFR-001`](./requirements.md#5-1-性能目標v1必須3エンドポイント) | `RUNBOOK.md` の trace 実行ログ、性能計測結果（p95） | ingest/search/show の 3 エンドポイントのみ対象に閾値比較 |
| [`REQ-NFR-002`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | インシデント記録（`docs/IN-*.md`）、運用タイムライン | `rto_minutes <= 30` かつ `rpo_minutes <= 5` を同時充足 |
| [`REQ-NFR-003`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | 検知時刻/暫定復旧時刻ログ | 検知〜暫定復旧が 15 分以内 |
| [`REQ-NFR-004`](./requirements.md#5-2-可用性復旧整合性回復運用nfr) | 再試行履歴、ジョブログ | 1 リクエストあたり再処理 2 回以内 |
| [`REQ-NFR-005`](./requirements.md#5-3-整合性回復要件archive-補償フロー) | 補償フロー実行ログ、`docs/IN-*.md` 起票有無 | 30分以内収束、未収束時は IN 起票 |

### Source
- `memx_spec_v3/docs/requirements.md#5-1-性能目標v1必須3エンドポイント`
- `memx_spec_v3/docs/requirements.md#5-2-可用性復旧整合性回復運用nfr`
- `memx_spec_v3/docs/requirements.md#5-3-整合性回復要件archive-補償フロー`

### Dependencies
- `RUNBOOK.md`
- `EVALUATION.md`
- `docs/IN-*.md`
