---
priority: high
owner: memx-core
deadline: 2026-03-12
status: planned
---

# TASK.gc-trigger-dryrun-03-03-2026

## Objective
- `memx_spec_v3/go/db/gc.go` の GC trigger 判定と dry-run JSON を固定化し、dry-run で副作用が発生しないことを保証する。

## Requirements
- 正常系: trigger 条件を満たす入力で GC 対象が抽出され、dry-run JSON が期待スキーマで出力されることを検証する。
- 入力エラー: 不正な trigger パラメータ（負値、未定義モード）で、既存のエラーハンドリング方針に一致することを検証する。
- 境界値: trigger 閾値ちょうど一致時に実行/非実行判定が仕様通りであることを検証する。
- 副作用有無: dry-run 実行時は DB 更新/削除が発生せず、実行モード時のみ副作用が発生することを検証する。
- JSON 互換: 既存 CLI/JSON 契約を維持し、キー名・型の後方互換を確認する。

## Commands
- `go test ./memx_spec_v3/go/db -run GC -count=1`
- `go test ./memx_spec_v3/go/db -run DryRun -count=1`
- `git status --short`

## Dependencies
- `memx_spec_v3/docs/requirements.md` の GC ポリシー節
- `TASK.memx-bootstrap-03-03-2026.md`

## Status
- planned
