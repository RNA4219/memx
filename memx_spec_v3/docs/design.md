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

## 5. 設計→契約→検証 導線
参照順は以下の固定順序とする。
1. `interfaces.md`（I/F 設計の確認）
2. `contracts/openapi.yaml`（HTTP 契約の確認）
3. `contracts/cli-json.schema.json`（CLI `--json` 契約の確認）
4. `RUNBOOK.md`（検証手順・トレース実行）
