---
priority: medium
owner: memx-core
deadline: 2026-03-14
status: planned
---

# TASK.reassess-design-docs-priority-03-03-2026

## Source
| Source | Purpose |
| --- | --- |
| `docs/design-docs-prioritization-spec.md#3-優先度決定ルール` | 設計書タスク優先度の判定基準 |
| `orchestration/memx-design-docs-authoring.md#phase-1-情報収集` | Phaseチェック項目の優先度付与運用 |
| `docs/TASKS.md#priority-記載ガイド設計書作成タスク` | Task Seed front matter 記載ルール |

## Node IDs
- requirements: 運用基準ノード

## Objective
- 既存 Task Seed（`TASK.*-03-03-2026.md`）の `priority` を新仕様で再評価し、判定根拠を追記して優先度表記を統一する。

## Requirements
- 対象は `TASK.memx-bootstrap-03-03-2026.md` / `TASK.replace-incident-placeholder-ids-03-03-2026.md` / `TASK.migrate-other-ddl-order-03-03-2026.md` / `TASK.gc-trigger-dryrun-03-03-2026.md` / `TASK.recall-query-normalization-03-03-2026.md` とする。
- 各Taskについて4軸（Blocker有無、REQ網羅率影響、契約差分 high 件数、Birdseye issue 有無）を判定し、根拠を `Requirements` または `Dependencies` に1行追記する。
- 判定結果に応じて front matter の `priority` を更新する（high/medium/low）。
- `priority` 更新時も Objective/実装要件は変更しない。

## Commands
- `for f in TASK.*-03-03-2026.md; do echo "### $f"; sed -n '1,20p' "$f"; done`
- `rg -n "priority:|Blocker|REQ網羅率|契約差分|Birdseye" TASK.*-03-03-2026.md`
- `git status --short`

## Dependencies
- `docs/design-docs-prioritization-spec.md` が mainline に反映済みであること

## Release Note Draft
- 既存 Task Seed の優先度を新しい設計書タスク判定仕様で再評価し、優先度判断の根拠を追跡可能にする。

## Status
- planned
