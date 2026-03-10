---
priority: high
owner: memx-core
deadline: 2026-03-12
status: planned
---

# TASK.gc-trigger-dryrun-03-03-2026

## Source
| Source | Purpose |
| --- | --- |
| `orchestration/memx-v1-bootstrap.md#phase-2` | 実装対象フェーズの起点 |
| `memx_spec_v3/docs/requirements.md#task-seed-source-fixed` | REQ-* の直接参照固定表 |

## Node IDs
- requirements: 仕様出典（requirements ノード）
- service: 実装対象ノード

## Objective
- `memx_spec_v3/go/db/gc.go` の GC trigger 判定と dry-run JSON を固定化し、dry-run で副作用が発生しないことを保証する。

## Requirements
- インシデント再発防止: [`docs/IN-20260303-002.md`](docs/IN-20260303-002.md) の `TP-01/TP-02` に従い、実インシデント由来条件を要件とテストへ明示的に転記する。
- 正常系: trigger 条件を満たす入力で GC 対象が抽出され、dry-run JSON が期待スキーマで出力されることを検証する。
- 入力エラー: 不正な trigger パラメータ（負値、未定義モード）で、既存のエラーハンドリング方針に一致することを検証する。
- 境界値: trigger 閾値ちょうど一致時に実行/非実行判定が仕様通りであることを検証する。
- 副作用有無: dry-run 実行時は DB 更新/削除が発生せず、実行モード時のみ副作用が発生することを検証する。
- JSON 互換: 既存 CLI/JSON 契約を維持し、キー名・型の後方互換を確認する。
- 無効時レスポンス固定依存: `mem.features.gc_short=false` かつ route 公開時の `POST /v1/gc:run` は `HTTP 500` + `INTERNAL` エラー（FAILED_PRECONDITION 相当フォールバック）を返す正本契約に従う。

## Commands
- `go test ./memx_spec_v3/go/db -run GC -count=1`
- `go test ./memx_spec_v3/go/db -run DryRun -count=1`
- `git status --short`

## Dependencies
- `memx_spec_v3/docs/requirements.md` の「0-1. Release Scope Matrix」および GC ポリシー節
- `TASK.memx-bootstrap-03-03-2026.md`

## Release Note Draft
- GC trigger 判定と dry-run JSON の挙動を固定し、dry-run 実行時にDB副作用が起きないことを利用者向けに保証する。

## Status
- planned
