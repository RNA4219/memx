---
priority: high
owner: memx-core
deadline: 2026-03-08
status: planned
---

# TASK.replace-incident-placeholder-ids-03-03-2026

## Source
| Source | Purpose |
| --- | --- |
| `HUB.codex.md#自動タスク分割フロー` | タスク分割フロー準拠 |
| `docs/TASKS.md#2-task-seed-必須項目` | Task Seed 必須項目の運用基準 |
| `memx_spec_v3/docs/requirements.md#task-seed-source-fixed` | REQ-* の直接参照固定表 |

## Objective
- 既存 Task Seed（`TASK.memx-bootstrap-03-03-2026.md` ほか）に残る `docs/IN-20260303-001.md` 参照を、実在するインシデントIDへ置換する。
- テンプレートID/TBD を参照したままレビュー通過しないよう、差し戻し基準に沿って是正を完了する。

## Requirements
- 対象 Task Seed を棚卸しし、`docs/IN-20260303-001.md` 参照箇所を一覧化する。
- 実在する `docs/IN-<実日付>-<連番>.md` を特定し、各 Task Seed の `Requirements`/`Source` を置換する。
- 置換後に `Source` からテンプレートIDや未確定値を除去できていることを確認する。
- 変更は Task Seed のトレーサビリティ修正のみに限定し、Objective/実装要件は改変しない。

## Commands
- `rg -n "IN-20260303-001|IN-<YYYYMMDD>-001|TBD" TASK.*-03-03-2026.md`
- `rg -n "^## Source|IN-" TASK.*-03-03-2026.md`
- `git status --short`

## Dependencies
- `docs/IN-<実日付>-<連番>.md` 形式の実インシデントが記録済みであること
- `TASK.memx-bootstrap-03-03-2026.md`

## Release Note Draft
- Task Seed のインシデント参照を実在IDへ統一し、再発防止要件のトレーサビリティを改善する。

## Status
- planned
