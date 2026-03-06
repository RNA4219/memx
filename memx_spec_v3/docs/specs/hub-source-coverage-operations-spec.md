# HUB Source Coverage Operations Spec

## 1. 目的
- 本仕様は `memx_spec_v3/docs/design-hub-source-coverage-spec.md` で定義した判定ロジックの運用手順を固定する。
- 実行タイミング・記録形式・承認責務を統一し、`gate_hub_source_coverage` の運用ぶれを防止する。

## 2. 適用範囲
- 対象 gate: `gate_hub_source_coverage`
- 判定キー・対象ソース・写像規則（`pass/warn/fail` → `high/medium/low`）は `memx_spec_v3/docs/design-hub-source-coverage-spec.md` を正本とする。
- 本書は「いつ実行するか」「どこにどう記録するか」「誰が承認するか」の運用のみを定義する。

## 3. 実行タイミング（必須）

### 3.1 Task Seed 起票時
- `docs/TASKS.md` の Task Seed 新規起票時に 1 回実行する。
- 目的は、着手前に `Incident` / `Orchestration` / `TASK` の入力欠落を検知すること。

### 3.2 Phase 遷移時
- `Phase 1→2`、`Phase 2→3`、`Phase 3→4` の各遷移直前に実行する。
- 判定結果は各 gate 判定に添付し、`design-phase-gate-spec.md` の総合判定へ反映する。

### 3.3 `Status: done` 直前
- Task Seed の `Status: done` 更新直前に最終実行する。
- 最終結果が `fail` または `hold` の場合は `done` への更新を禁止する。

## 4. 記録先・記録形式（固定）

### 4.1 記録ファイル
- 保存先は `memx_spec_v3/docs/reviews/` 配下に固定する。
- ファイル名は次に固定する。
  - `HUB-SOURCE-COVERAGE-<YYYYMMDD>-<task_or_phase_id>.md`

### 4.2 必須フィールド
各実行記録は、判定キーごとに次のフィールドを必須とする。

| フィールド | 必須値 |
| --- | --- |
| `decision` | `pass` / `warn` / `fail` / `hold` |
| `target_source` | 判定対象ソース（例: `docs/IN-*.md`） |
| `search_key` | `Incident` / `Orchestration` / `TASK` |
| `operator` | 実行者（GitHub ID または運用上の一意名） |
| `executed_at` | 実行日時（ISO 8601, JST 明記） |

### 4.3 記録テンプレート（最小）
```md
## HUB Source Coverage Run
- task_or_phase_id: <ID>
- operator: <name>
- executed_at: <YYYY-MM-DDTHH:mm:ss+09:00>

| search_key | target_source | decision | notes |
| --- | --- | --- | --- |
| Incident | docs/IN-*.md | pass | ... |
| Orchestration | orchestration/*.md | warn | ... |
| TASK | TASK.* | fail | ... |
```

## 5. 承認責務

### 5.1 必須役割
- `owner`（必須）: 判定実行と記録作成の責任者。
- `reviewer`（必須）: 判定妥当性と証跡完全性の承認責任者。
- `owner` と `reviewer` の兼務は禁止する。

### 5.2 承認条件
- `reviewer` は、3 判定キーすべての記録と必須フィールド充足を確認して承認する。
- 承認後にのみ、Phase 遷移または `Status: done` 更新を許可する。

### 5.3 fail/hold 時の差し戻し条件
- 次のいずれかに該当する場合は `reviewer` が差し戻す。
  1. いずれかの判定キーが `fail`。
  2. `hold` 判定（追加調査待ち・証跡不足・検索未完了）が 1 件でもある。
  3. 必須フィールド欠落、または記録ファイル命名規則違反。
- 差し戻し後は、補完実施と再記録を行い、同一 Task/Phase ID で再承認する。

## 6. 既存仕様との接続
- gate 判定列と総合判定への反映先: `memx_spec_v3/docs/design-phase-gate-spec.md`
- Task Seed 起票・`Status: done` 運用との接続: `docs/TASKS.md`
- Incident 起点のトレース要件との整合: `memx_spec_v3/docs/incident-to-task-traceability-spec.md`

## 7. 運用チェックリスト
- [ ] Task Seed 起票時・Phase 遷移時・`Status: done` 直前の 3 タイミングで実行した
- [ ] 記録を `memx_spec_v3/docs/reviews/HUB-SOURCE-COVERAGE-<YYYYMMDD>-<id>.md` に保存した
- [ ] `decision/target_source/search_key/operator/executed_at` を全判定キーで記録した
- [ ] `owner` と `reviewer` の分離を満たした
- [ ] `fail/hold` の場合に差し戻しと再承認を実施した
