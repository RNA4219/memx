---
intent_id: memx-acceptance-evaluation-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# EVALUATION

## v1 受け入れ基準（Release Scope Matrix 準拠）
- 判定対象（MUST (v1)）: `mem in short` / `mem out search` / `mem out show` と `POST /v1/notes:ingest` / `POST /v1/notes:search` / `GET /v1/notes/{id}`。
- 入出力互換: MUST (v1) の CLI→API 入出力マッピングが保持される。
- エラーコード整合: MUST (v1) の入力不備は 400 系、内部障害は 500 系を返す。
- 最小性能目標（同一計測条件、MUST (v1) API のみ）:
  - `POST /v1/notes:ingest`: P50 <= 120ms, P95 <= 250ms
  - `POST /v1/notes:search`: P50 <= 80ms, P95 <= 180ms
  - `GET /v1/notes/{id}`: P50 <= 40ms, P95 <= 90ms

## 要件IDトレーサビリティ（判定基準との相互参照）
| Requirement ID | 判定基準 | 判定ルール参照 |
| --- | --- | --- |
| `REQ-CLI-001` | pass/fail | 本書「v1 受け入れ基準（Release Scope Matrix 準拠）」・`RUNBOOK.md` manual 手順 |
| `REQ-API-001` | pass/fail | 本書「v1 受け入れ基準（Release Scope Matrix 準拠）」・`RUNBOOK.md` manual 手順 |
| `REQ-GC-001` | pass/fail | `RUNBOOK.md` の `mem gc short` 手順および manual コマンド |
| `REQ-SEC-001` | pass/fail | `requirements.md` 2-7-1 の `sensitivity` 契約 |
| `REQ-RET-001` | pass/fail/waiver | 本書「性能合否基準（fail / waiver）」の運用を準用し、保持期限逸脱は waiver 記録必須 |
| `REQ-ERR-001` | pass/fail | 本書「エラーコード整合」および `requirements.md` 6-4 |
| `REQ-NFR-001` | pass/fail/waiver | 本書「性能合否基準（fail / waiver）」および「REQ-NFR-001 合否判定ルール」 |

## 性能合否基準（fail / waiver）
- 判定対象データセットは `short` 10,000 件 / 1件 約500文字（UTF-8）、ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）、ウォームアップ20回、本計測200回で固定する。
- 合格（pass）: `REQ-NFR-001` の `ingest` / `search` / `show` の **p50/p95 が全て閾値以内**。
- 不合格（fail）: `REQ-NFR-001` でいずれか 1 指標でも閾値超過、または計測条件が不一致。
- 例外承認（waiver）:
  - `REQ-NFR-001` / `REQ-RET-001` の fail を一時的に許容する場合のみ適用できる。
  - 最低限、`docs/IN-YYYYMMDD-XXX.md` のインシデント記録、超過理由、是正期限、暫定運用策、責任者を明記する。
  - waiver は期限付きとし、期限超過時は自動的に fail 扱いへ戻す。

## `REQ-NFR-001` 合否判定ルール
- 判定対象は `results.ingest` / `results.search` / `results.show` の `p50_ms` / `p95_ms` のみ。
- fail 条件（閾値超過）は以下のいずれか 1 つでも満たした場合とする。
  - `results.ingest.p50_ms > 120` または `results.ingest.p95_ms > 250`
  - `results.search.p50_ms > 80` または `results.search.p95_ms > 180`
  - `results.show.p50_ms > 40` または `results.show.p95_ms > 90`
- pass 条件は上記 6 指標がすべて閾値以内（`<=`）であること。
- `p50_ms` / `p95_ms` の欠損、または計測条件不一致は fail 扱いとする。


## 運用NFR（可用性/復旧/整合性回復）合否基準

- 対象要件: `REQ-NFR-002` / `REQ-NFR-003` / `REQ-NFR-004` / `REQ-NFR-005` / `REQ-NFR-006`

### 証跡ファイル（必須）
- `artifacts/ops/incident-summary.json`
  - 必須キー: `incident_id`, `detected_at`, `mitigated_at`, `resolved_at`, `rto_minutes`, `rpo_minutes`, `retry_count`
- `artifacts/ops/recovery-log.ndjson`
  - 必須イベント: `detect`, `retry`, `rollback`（実施時）, `replan`（実施時）, `mitigate`, `resolve`
- `docs/IN-*.md`
  - `REQ-NFR-006` で定義した最小監査項目を満たすこと

### 判定ロジック
- pass 条件（全件必須）:
  1. `incident-summary.json.rto_minutes <= 30`（`REQ-NFR-002`）
  2. `incident-summary.json.rpo_minutes <= 5`（`REQ-NFR-002`）
  3. `mitigated_at - detected_at <= 15分`（`REQ-NFR-003`）
  4. `retry_count <= 2`（`REQ-NFR-004`）
  5. `recovery-log.ndjson` に整合性回復イベント（`rollback` または `replan`）が存在し、30 分以内に収束または `docs/IN-*.md` 起票済み（`REQ-NFR-005`）
  6. 対応する `docs/IN-*.md` が最小監査項目を欠落なく記載（`REQ-NFR-006`）
- fail 条件:
  - 上記 1〜6 のいずれか 1 つでも不成立
  - 証跡ファイル欠損、または時刻/回数の整合が取れない

## governance/metrics.yaml 同期運用ルール
- `governance/metrics.yaml` は本書（`EVALUATION.md`）を正本として同期する。
- 少なくとも性能項目 `ingest` / `search` / `show` の **項目名** と **閾値文字列** は完全一致させる。
- 不一致を検知したレビュー/CI は fail とし、同一コミットで差分を解消する。

## 閾値項目と RUNBOOK 出力項目の対応
- 判定に使用する値は `RUNBOOK.md` の `artifacts/perf/perf-result.json` に含まれる `results.<endpoint>.p50_ms` / `results.<endpoint>.p95_ms` とする。
- 対応は 1 対 1 で固定し、別名指標は使用しない。
  - `POST /v1/notes:ingest` ↔ `results.ingest.p50_ms` / `results.ingest.p95_ms`
  - `POST /v1/notes:search` ↔ `results.search.p50_ms` / `results.search.p95_ms`
  - `GET /v1/notes/{id}` ↔ `results.show.p50_ms` / `results.show.p95_ms`

## 性能計測条件
- ノート件数: 10,000 件（short ストア）
- 本文長: 1 ノートあたり約 500 文字（UTF-8 プレーンテキスト）
- 実行環境: ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）
- ウォームアップ有無: あり（各エンドポイント 20 リクエスト）
- 計測回数: 各エンドポイント 200 リクエスト（ウォームアップ除外）

## スコープ別評価ポリシー
- MUST (v1): 合否判定対象（本ドキュメントの全受け入れ判定に使用）。
- SHOULD (v1.x): 判定対象外。`mem.features.gc_short=true` で有効化した実験運用時のみ、参考指標として別レポート化する。
- FUTURE (v1.1+): 判定対象外（v1 の受け入れ結果に影響させない）。

## 明示的な判定対象外
- CLI: `mem gc short`（SHOULD 実験機能）, `mem out recall`, `mem working`, `mem tag`, `mem meta`, `mem lineage`, `mem distill`, `mem out context`。
- API: `POST /v1/gc:run`（SHOULD 実験機能）, Recall/Working/Tag/Meta/Lineage 系 API。
- よって Recall 評価条件は v1 受け入れ判定に使用しない（必要時は別紙評価）。

## LLM 呼び出し評価条件
- 1リクエスト 15 秒タイムアウト。
- 最大 2 回リトライ（指数バックオフ）。
- 再試行可: ネットワーク障害、タイムアウト、HTTP 429/502/503/504。
- 再試行不可: 入力不正、認証/認可失敗、スキーマ不整合。

## インシデント対応要件
- インシデントは `IN-YYYYMMDD-XXX` 形式で記録。
- 初動記録に「検知」「影響」「5 Whys」「再発防止」「タイムライン」を必須記載。
