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
- 最小性能目標: `ingest` / `search` / `show` がローカル単体で実用応答時間を維持する。

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
