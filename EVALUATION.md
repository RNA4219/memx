---
intent_id: memx-acceptance-evaluation-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# EVALUATION

## v1 受け入れ基準
- 入出力互換: CLI→API の入出力マッピングが保持される。
- エラーコード整合: 入力不備は 400 系、内部障害は 500 系を返す。
- 最小性能目標（同一計測条件）:
  - `POST /v1/notes:ingest`: P50 <= 120ms, P95 <= 250ms
  - `POST /v1/notes:search`: P50 <= 80ms, P95 <= 180ms
  - `GET /v1/notes/{id}`: P50 <= 40ms, P95 <= 90ms

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

## 必須スコープ評価
- 必須コマンド: `mem in short`, `mem out search`, `mem out show`。
- 必須 API: `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}`。
- v1 非対象（将来機能）: GC / recall / working / tag / meta / lineage。

## Recall 評価条件
- 類似度閾値 `score >= 0.20` を適用。
- `top-k` は 1..50 に正規化（既定 8）。
- `range` は 0..10（既定 3）。
- `--stores` は trim/小文字/重複排除で正規化し、不正値は入力エラー。

## LLM 呼び出し評価条件
- 1リクエスト 15 秒タイムアウト。
- 最大 2 回リトライ（指数バックオフ）。
- 再試行可: ネットワーク障害、タイムアウト、HTTP 429/502/503/504。
- 再試行不可: 入力不正、認証/認可失敗、スキーマ不整合。

## インシデント対応要件
- インシデントは `IN-YYYYMMDD-XXX` 形式で記録。
- 初動記録に「検知」「影響」「5 Whys」「再発防止」「タイムライン」を必須記載。
