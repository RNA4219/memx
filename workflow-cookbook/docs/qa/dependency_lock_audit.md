---
intent_id: QA-LOCK-AUDIT
owner: qa-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-04-04
---

# Dependency Lock Audit

`go.mod` / `go.sum` の整合性を監査した記録を残す。

## 実施手順

1. `go mod tidy` が差分を生まないことを確認する。
2. `go mod verify` が成功することを確認する。
3. `memx_spec_v3/go/go.sum` が tracked（`git ls-files` に含まれる）ことを確認する。

## 監査記録

| 実施日 | 実施者 | 対象ディレクトリ | go.mod チェック | go.sum チェック | `go mod tidy` 差分 | `go mod verify` | tracked 確認 | 備考 |
| :-- | :-- | :-- | :-- | :-- | :-- | :-- | :-- | :-- |
| YYYY-MM-DD | @handle | `memx_spec_v3/go` | OK/NG | OK/NG | なし/あり | OK/NG | OK/NG | |
