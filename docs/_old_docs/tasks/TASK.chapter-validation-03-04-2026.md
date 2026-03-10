---
priority: high
owner: memx-core
deadline: 2026-03-11
status: reviewing
---

# TASK.chapter-validation-03-04-2026

## Source
- `orchestration/memx-design-docs-authoring.md#phase-3`
- `memx_spec_v3/docs/design-chapter-validation-spec.md#章別検証サマリ`
- `docs/TASKS.md#2-task-seed-必須項目`

## Node IDs
- design-phase3: chapter validation gate

## Objective
- chapter別検証で `req_coverage` と mapping 整合を同時に確定し、Phase 3/4 の受け入れ判定の入力を固定する。

## Requirements
- `chapter_id` ごとの検証結果を一覧化し、全 chapter 行に `req_coverage` と `mapping_match_check` を記録する。
- 検証結果は `memx_spec_v3/docs/reviews/inventory/` 配下に成果物として保存し、後続 Task から参照可能にする。
- 完了条件（固定）: **全 chapter 行で `req_coverage=100%` かつ `mapping_match_check=pass`**。
- `reviewing -> done` 判定は本Task単体で完結させ、他Taskの編集有無を条件にしない。

## Commands
- `mkdir -p memx_spec_v3/docs/reviews/inventory`
- `date +%Y%m%d`
- `rg -n "req_coverage|mapping_match_check|chapter_id" memx_spec_v3/docs/reviews/inventory`
- `go test ./...`

## Dependencies
- なし

## Release Note Draft
- chapter別検証の完了判定を `req_coverage` / `mapping_match_check` に固定し、後続タスク参照用の成果物出力を標準化。

## Status
- reviewing
- reviewing 継続条件: chapter単位の検証表に未記入行、または `req_coverage<100%` / `mapping_match_check!=pass` が1件でも残る。
- done 遷移条件: **全 chapter 行で `req_coverage=100%` かつ `mapping_match_check=pass`** を満たし、成果物パスを記録済み。
