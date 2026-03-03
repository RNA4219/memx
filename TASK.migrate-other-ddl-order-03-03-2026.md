---
priority: high
owner: memx-core
deadline: 2026-03-13
status: planned
---

# TASK.migrate-other-ddl-order-03-03-2026

## Objective
- `memx_spec_v3/go/db/migrate_other.go` の `chronicle`/`memopedia`/`archive` DDL 適用順と `user_version` 更新運用を明確化し、再実行安全性を担保する。

## Requirements
- 正常系: 初回マイグレーションで `chronicle` → `memopedia` → `archive` の順に DDL が適用され、`user_version` が期待値へ更新されることを検証する。
- 入力エラー: 想定外 `user_version` または欠損スキーマ時に、既存例外方針に沿って中断/通知されることを検証する。
- 境界値: 既に最新 `user_version` の DB へ再適用した場合に no-op で終了することを検証する。
- 順序保証: DDL 実行順の依存関係が壊れないよう、順序を固定するテストを追加する。
- 互換性: `user_version` 運用ルールを維持し、既存 DB への後方互換を確認する。

## Commands
- `go test ./memx_spec_v3/go/db -run MigrateOther -count=1`
- `go test ./memx_spec_v3/go/db -run UserVersion -count=1`
- `git status --short`

## Dependencies
- `memx_spec_v3/docs/requirements.md` のマイグレーション方針
- `TASK.memx-bootstrap-03-03-2026.md`

## Status
- planned
