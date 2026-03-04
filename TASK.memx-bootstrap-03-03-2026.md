---
priority: high
owner: memx-core
deadline: 2026-03-10
status: planned
---

# TASK.memx-bootstrap-03-03-2026

## Source
| Source | Purpose |
| --- | --- |
| `memx_spec_v3/docs/requirements.md#task-seed-source-fixed` | Task Seed の Source/Requirements 直接参照の固定表（REQ-*） |
| `docs/TASKS.md#2-task-seed-必須項目` | Task Seed 必須項目の運用基準 |

## Objective
- memx リポジトリに Task Seed 運用の最小テンプレートを導入し、命名規則と記載要件を固定化する。

## Requirements
- インシデント再発防止: [`docs/IN-20260303-001.md`](docs/IN-20260303-001.md) の `TP-01/TP-02` に従い、実インシデント由来条件を要件とテストへ明示的に転記する。
- `docs/TASKS.md` に必須項目（Objective/Requirements/Commands/Dependencies/Status）を定義する。
- Task Seed の命名規則を `TASK.<slug>-<MM-DD-YYYY>.md` に統一する。
- 完了タスクを `memx_spec_v3/CHANGES.md` へ移送する手順を明文化する。
- Status 語彙を `planned/active/in_progress/reviewing/blocked/done` に統一する。

## Commands
- `date +%m-%d-%Y`
- `git status --short`

## Dependencies
- none

## Status
- planned
