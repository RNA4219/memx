---
priority: high
owner: memx-core
deadline: 2026-03-12
status: reviewing
---

# TASK.phase3-4-artifacts-03-04-2026

## Source
- `orchestration/memx-design-docs-authoring.md#phase-4`
- `memx_spec_v3/docs/design-acceptance-lifecycle-spec.md#da-lc-03-04`
- `docs/TASKS.md#2-1-3-phase-3-4-成果物存在トリガー必須`

## Node IDs
- design-phase4: artifact gate

## Objective
- Phase 3/4 の必須成果物（contract align / acceptance / tracked file）の存在判定を固定し、受け入れ審査の差戻し条件を明確化する。

## Requirements
- `CONTRACT-ALIGN-*` と `LATEST.md` の存在を確認し、確認ログを保存する。
- `DESIGN-ACCEPTANCE-<実日付>.md` の存在を確認し、命名が実日付であることを明示する。
- `go.sum` が tracked 状態であることを確認する。
- 完了条件（固定）: **`CONTRACT-ALIGN-*` と `LATEST.md` 存在、`DESIGN-ACCEPTANCE-<実日付>.md` 存在、`go.sum` tracked**。
- 5条件チェックは `memx_spec_v3/docs/reviews/DESIGN-GATE-EVIDENCE-INDEX.md` を一次参照とし、`req_coverage` / `design_acceptance` / `mapping_match_check` の参照版を `20260304-003` 系列（acceptance は `20260304`）で一致させる。

## Commands
- `rg --files memx_spec_v3/docs/reviews | rg 'CONTRACT-ALIGN-|LATEST\.md|DESIGN-ACCEPTANCE-'`
- `git ls-files go.sum`
- `git status --short go.sum`
- `go test ./...`

## Dependencies
- `TASK.chapter-validation-03-04-2026.md` の成果物（chapter validation inventory）の参照のみを許可する。
- 同一ファイル同時編集を避けるため、本TaskはA成果物の読み取り専用とし、A対象ファイルへ追記しない。


## Notes
- 差分説明: Evidence Index の `mapping_match_check` 根拠を `...-003` へ更新し、DoD の5条件一次参照（`docs/TASKS.md#2-1-4-1`）との整合を維持。

## Release Note Draft
- Phase 3/4 受け入れの必須成果物チェックを独立タスク化し、完了条件と差戻し条件を明文化。

## Status
- reviewing
- reviewing 継続条件: 必須成果物の欠落、命名不一致、または `go.sum` untracked/未追跡。
- done 遷移条件: **`CONTRACT-ALIGN-*` と `LATEST.md` 存在、`DESIGN-ACCEPTANCE-<実日付>.md` 存在、`go.sum` tracked** を満たし、A成果物参照のみで判定できること。
