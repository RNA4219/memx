---
priority: high
owner: memx-core
deadline: 2026-03-10
status: planned
---

# TASK.memx-bootstrap-03-03-2026

## Objective
- memx リポジトリに Task Seed 運用の最小テンプレートを導入し、命名規則と記載要件を固定化する。

## Requirements
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
