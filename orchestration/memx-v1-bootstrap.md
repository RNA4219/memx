# memx v1 Bootstrap Orchestration

## Phase 1: Alignment & Scope Freeze
### Dependencies
- `BLUEPRINT.md`
- `RUNBOOK.md`
- `memx_spec_v3/docs/requirements.md`

- [ ] 主要ユースケースを 30 分以内で棚卸しし、対象/非対象を 1 ページで確定する
- [ ] 要件の曖昧点を 3 件以内に絞って確認事項化する
- [ ] 既存運用手順との衝突ポイントを列挙し、優先度を付ける

## Phase 2: Task Seed Decomposition
### Dependencies
- `CHECKLISTS.md`
- `TASK.memx-bootstrap-03-03-2026.md`
- `memx_spec_v3/docs/requirements.md`

- [ ] フェーズ横断タスクを 0.5 日以内の単位へ分割する
- [ ] 各タスクに Done 条件（検証コマンド or 成果物）を 1 つ以上付与する
- [ ] 依存順に並べ替え、並行実行可能タスクを明示する

## Phase 3: Execution Readiness
### Dependencies
- `RUNBOOK.md`
- `HUB.codex.md`
- `memx_spec_v3/docs/quickstart.md`

- [ ] 実行コマンドのテンプレートを確認し、手戻りリスクを事前に記録する
- [ ] 検証観点（lint/type/test）をタスク単位で割り当てる
- [ ] 初回実装バッチを 1 日未満で完了する順序へ最終調整する
