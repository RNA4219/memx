---
priority: high
owner: memx-core
deadline: 2026-03-11
status: planned
---

# TASK.recall-query-normalization-03-03-2026

## Objective
- `memx_spec_v3/go/db/recall.go` の検索入力正規化と閾値判定を安定化し、`top-k`/`range`/`stores` の解釈差分と埋め込み未設定時の挙動を明確化する。

## Requirements
- 正常系: `top-k`/`range`/`stores` を含む有効入力で、期待件数・期待ストア集合が返ることを検証する。
- 入力エラー: 不正な `top-k`（0未満、非数値）または `range` 逆転入力で、既存ポリシーに沿ったエラー応答になることを検証する。
- 境界値: `top-k=1`、`range` の下限/上限一致、`stores` 空配列時の正規化結果を検証する。
- 閾値適用: score threshold の有効/無効境界（等値含む）でフィルタ結果が仕様通りであることを検証する。
- 埋め込み未設定: embedding 未設定時にフォールバックまたは明示エラーのいずれか既存仕様へ揃え、互換性を維持する。

## Commands
- `go test ./memx_spec_v3/go/db -run Recall -count=1`
- `go test ./memx_spec_v3/go/db -run QueryNormalization -count=1`
- `git status --short`

## Dependencies
- `memx_spec_v3/docs/requirements.md` の検索仕様確定
- `TASK.memx-bootstrap-03-03-2026.md`

## Status
- planned
