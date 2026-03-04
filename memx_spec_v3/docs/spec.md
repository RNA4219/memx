---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
priority: high
---

# memx 仕様インデックス（正本/補助の定義）

本書は `memx_spec_v3/docs/` の仕様参照導線を一本化するためのインデックスである。

## 1. 文書の役割分担（正本/補助）

### 1-1. 正本（Normative）

- `requirements.md`
  - memx v1.3 の仕様正本。
  - 要件ID、スコープ定義、API/CLI/NFR/運用要件の判定根拠は本書を優先する。
- `contracts/openapi.yaml`
  - HTTP API 契約とエラー応答スキーマの正本。
- `contracts/cli-json.schema.json`
  - CLI `--json` 出力契約の正本。
- `traceability.md`
  - 主要 REQ-ID の要件→設計→I/F→評価→契約の対応正本。

### 1-2. 補助（Secondary）

- `interfaces.md`
  - 人間可読の I/O 説明と互換方針を示す補助仕様。
- `CONTRACTS.md`
  - 正本スキーマ（`contracts/openapi.yaml` / `contracts/cli-json.schema.json`）への索引・抜粋。
  - 重複定義（フィールド型・required・制約の再定義）を持たない。
- `error-contract.md`
  - エラー契約の運用要約。正本は `contracts/openapi.yaml`。
- `quickstart.md`
  - 起動/疎通確認手順。契約定義の正本ではない。
- `operations-spec.md`
  - 運用章（インシデント起票・waiver・RTO/RPO・補償収束）の参照固定用サマリ。正本は `requirements.md`。

## 2. 要件別の参照導線

### 2-1. API 契約

- 正本: `contracts/openapi.yaml`
- 仕様背景と受け入れ条件: `requirements.md` の「6. API 要件（v1.3 追加）」

### 2-2. CLI 契約

- 正本: `contracts/cli-json.schema.json`
- CLI 要件: `requirements.md` の「3. CLI 要件」

### 2-3. エラー契約

- 正本: `contracts/openapi.yaml`（`components.responses` / `components.schemas.Error*`）
- 運用要約: `error-contract.md`
- 要件側の基準: `requirements.md` の「6-4. エラーモデル」

### 2-4. NFR（非機能要件）

- 正本: `requirements.md` の「5. 非機能要件」

### 2-5. 運用要件

- 正本: `requirements.md` の「11. インシデント対応要件（運用）」
- 運用章の固定参照先: `operations-spec.md`

## 3. 更新順序（契約変更時・固定）

1. `requirements.md` を更新する。
2. `traceability.md` を更新する。
3. 正本スキーマ（`contracts/openapi.yaml` / `contracts/cli-json.schema.json`）を更新する。
4. `interfaces.md` と `CONTRACTS.md` を更新する。
5. `EVALUATION*` / `operations-spec.md`（RUNBOOK 相当）を更新する。
6. 必要に応じて `memx_spec_v3/README.md` の導線を更新する。
---

# memx 仕様（spec）

## 1. 対象ユースケース
- ローカル環境での個人メモ投入・検索・参照。
- LLM/Agent からの機械呼び出し（CLI/API）での短期記憶運用。
- 将来の蒸留・昇格・GC を見据えた 4 ストア構成（short/chronicle/memopedia/archive）。

## 2. スコープ境界
### In Scope（v1）
- CLI: `mem in short` / `mem out search` / `mem out show`。
- API: `POST /v1/notes:ingest` / `POST /v1/notes:search` / `GET /v1/notes/{id}`。
- ローカル SQLite を前提にした単体運用。
- CLI `--json` と API レスポンスの同型維持。

### Out of Scope（v1）
- Web UI。
- マルチユーザー運用、認証・認可・監査基盤の本格提供。
- 常駐必須プロセス設計。
- 完全自律エージェントランタイム。

## 3. 非ゴール
- GUI での操作体験最適化。
- クラウド前提の水平分散や外部ベクターDB必須化。
- v1 内での破壊的 API/CLI 変更。

## 4. 受け入れ観点
- 互換性: v1 必須 I/F の後方互換を維持する。
- エラー: 入力不備は 400 系、内部障害は 500 系で返す。
- 品質: ingest/search/show がローカル単体で実用応答時間を満たす。
- 安全性: fail-closed 方針に従い、機密入力は保存拒否できる。
