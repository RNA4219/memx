---
priority_phase_1: high
priority_phase_2: high
priority_phase_3: medium
owner: memx-core
deadline: 2026-03-17
status: planned
---

# memx v1 Bootstrap Orchestration

## Phase 1: Alignment & Scope Freeze
### Dependencies
- `BLUEPRINT.md`
- `RUNBOOK.md`
- `memx_spec_v3/docs/requirements.md`

- [ ] 主要ユースケースを 30 分以内で棚卸しし、対象/非対象を 1 ページで確定する
- [ ] 要件の曖昧点を 3 件以内に絞って確認事項化する
- [ ] 既存運用手順との衝突ポイントを列挙し、優先度を付ける

### Done Criteria
- `orchestration/memx-v1-scope-freeze.md` を作成し、目的・非目的を明記している
- レビュー承認条件として `owner: memx-core` を含む 2 名以上の Approve が記録されている
- Phase 1 のチェックボックスがすべて完了している

## Phase 2: Task Seed Decomposition
### Dependencies
- `CHECKLISTS.md`
- `TASK.memx-bootstrap-03-03-2026.md`
- `memx_spec_v3/docs/requirements.md`

- [ ] フェーズ横断タスクを 0.5 日以内の単位へ分割する
- [ ] 各タスクに Done 条件（検証コマンド or 成果物）を 1 つ以上付与する
- [ ] 依存順に並べ替え、並行実行可能タスクを明示する
- [ ] 個別 Task Seed を作成する（`TASK.recall-query-normalization-03-03-2026.md` / `TASK.gc-trigger-dryrun-03-03-2026.md` / `TASK.migrate-other-ddl-order-03-03-2026.md`）

### Done Criteria
- Task Seed が 3 件以上作成され、`TASK.*.md` としてリポジトリに存在する
- 依存解決済み条件として、各 Task Seed の `Depends on` が循環参照なしでトポロジカル順に並んでいる
- 各 Task Seed に必須検証コマンド（最低 1 つの `lint` / `type` / `test`）が明記されている
- Phase 2 のチェックボックスがすべて完了している

## Phase 3: Execution Readiness
### Dependencies
- `RUNBOOK.md`
- `HUB.codex.md`
- `memx_spec_v3/docs/quickstart.md`

- [ ] 実行コマンドのテンプレートを確認し、手戻りリスクを事前に記録する
- [ ] 検証観点（lint/type/test）をタスク単位で割り当てる
- [ ] 初回実装バッチを 1 日未満で完了する順序へ最終調整する

### Done Criteria
- 初回実装バッチ着手前に `ruff check .` / `mypy --strict .` / `pytest` のゲートを定義し、全通過を実行条件として明記している
- 上記ゲート失敗時は「直前の安定コミットへロールバック」または「Task Seed を再分割して再計画」のいずれかを必須実施として記載している
- 再計画時は失敗要因と再試行条件を `TASK.memx-bootstrap-03-03-2026.md` に追記する運用が明記されている
- Phase 3 のチェックボックスがすべて完了している
