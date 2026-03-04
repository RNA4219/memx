# Design Phase Gate Spec

## 1. 目的と正本範囲
- 本書は、Design Docs オーサリングの **Phase 1〜4 の gate 判定正本** である。
- 各 Phase の gate 判定は、本書で定義する entry/exit criteria・fail 条件・遷移条件に従う。
- `orchestration/memx-design-docs-authoring.md` は実施手順の正本とし、判定ロジックを重複定義しない。
- 優先度評価の4軸（Blocker / REQ網羅率 / 契約差分 high 件数 / Birdseye issue）に `HUB入力カバレッジ` を加えた5列を、各 gate の判定列として必須入力とする（`docs/design-docs-prioritization-spec.md` 準拠）。

## 2. Gate 判定列（全Phase共通）

| 列名 | 値 | 判定基準 |
| --- | --- | --- |
| `gate_blocker` | `high/medium/low` | 後続Phase停止なら `high`、遅延のみなら `medium`、影響なしは `low` |
| `gate_req_coverage` | `high/medium/low` | REQ網羅率100%未達確定なら `high`、低下可能性なら `medium`、影響なしは `low` |
| `gate_contract_high` | `high/medium/low` | 契約差分 `high` が1件以上なら `high`、`medium/low` 差分のみは `medium`、差分なしは `low` |
| `gate_birdseye_issue` | `high/medium/low` | node_id参照切れ/caps欠落は `high`、軽微issueは `medium`、issueなしは `low` |
| `gate_hub_source_coverage` | `high/medium/low` | HUB必須入力（検索キー/対象パス）で Incident または Orchestration の欠落があれば `high`、欠落はあるが Incident/Orchestration 欠落はない場合は `medium`、欠落なしは `low` |

### 2.1 HUB入力カバレッジ判定の固定ルール
- `gate_hub_source_coverage` の判定根拠は以下の検索キー・対象パスに固定する。
  - 検索キー: `Incident`, `Orchestration`, `TASK`
  - 対象パス: `docs/IN-*.md`, `orchestration/*.md`, `TASK.*`
- 判定時は上記3系統をすべて確認し、レビュー時に別解釈（任意キーワード/任意パス追加）を持ち込まない。

### 2.2 総合判定ルール
1. 5列のいずれかが `high` の場合は gate `fail`。
2. `high` がなく `medium` が1つ以上ある場合は gate `hold`（差し戻し解消後に再判定）。
3. 5列すべて `low` の場合のみ gate `pass`。

## 3. Phase Gate 定義

## Phase 1（情報収集）
### Entry criteria
- `memx_spec_v3/docs/design-update-trigger-spec.md` で対象 Trigger IDs が確定済み。
- 入力ソース（requirements/design/interfaces/traceability/EVALUATION/RUNBOOK/Birdseye index）が参照可能。

### 入力成果物
- `memx_spec_v3/docs/requirements.md`
- `memx_spec_v3/docs/traceability.md`
- `memx_spec_v3/docs/design.md`
- `memx_spec_v3/docs/interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`

### 必須チェック
- Design Source Inventory を `memx_spec_v3/docs/design-source-inventory-spec.md` の必須列で作成する。
- `source_path#section` は `memx_spec_v3/docs/design-reference-resolution-spec.md` に従って正規化する。
- gate 判定列（5軸）を記録する。

### fail条件
- 入力成果物の未取得が1件以上。
- Source 正規化未完了または曖昧参照が1件以上。
- gate 判定列のいずれかが `high`。

### 出力成果物
- `memx_spec_v3/docs/reviews/inventory/DESIGN-SOURCE-INVENTORY-YYYYMMDD.md`
- Task Seed 化可能な抽出結果（1項目=1タスク、<=0.5d）

### 次Phase遷移条件
- gate 判定が `pass`。
- 抽出表の未解決行（blocked/reviewing）が 0 件。

## Phase 2（章別ドラフト）
### Entry criteria
- Phase 1 gate `pass`。
- 章対応表（chapter_id -> node_id）の更新対象が確定済み。

### 入力成果物
- Phase 1 の抽出表
- `memx_spec_v3/docs/design-template.md`
- `memx_spec_v3/docs/design-chapter-node-mapping-spec.md`

### 必須チェック
- 各章に `Source/Node IDs/Objective/Requirements/Commands/Dependencies/Status` を記載する。
- 各章で gate 判定列（5軸）を更新する。
- `docs/TASKS.md` の語彙（planned/active/in_progress/reviewing/blocked/done）を使用する。

### fail条件
- 必須項目欠落の章が1件以上。
- 章単位 gate 判定列のいずれかが `high`。
- `Status` 語彙逸脱が1件以上。

### 出力成果物
- 章別 Task Seed ドラフト
- 更新済み章対応表（chapter_id -> node_id）

### 次Phase遷移条件
- 全章の gate 判定が `pass`。
- 未解決事項が 0.5d 単位タスクへ分解済み。

## Phase 3（契約整合）
### Entry criteria
- Phase 2 gate `pass`。
- 章別ドラフトに契約参照（OpenAPI/CLI schema/requirements）が付与済み。

### 入力成果物
- Phase 2 の章別 Task Seed
- `memx_spec_v3/docs/contract-alignment-spec.md`
- `memx_spec_v3/docs/requirements-coverage-spec.md`
- `memx_spec_v3/docs/link-integrity-spec.md`
- `docs/birdseye/memx-birdseye-validation-spec.md`

### 必須チェック
- REQ網羅率、契約差分 high 件数、リンク不達件数、Birdseye issue 件数を章別に算出する。
- 算出結果を gate 判定列（5軸）へ反映する。
- 判定結果は `memx_spec_v3/docs/design-doc-dod-spec.md` と矛盾しない。

### fail条件
- 契約差分 `high` が1件以上。
- REQ網羅率 100% 未達。
- リンク不達または Birdseye issue が未解消。
- gate 判定列のいずれかが `high`。

### 出力成果物
- 章別検証サマリ（`chapter_id/req_coverage/contract_alignment_high_count/link_broken_count/birdseye_issue_count/evidence_paths`）
- `memx_spec_v3/docs/contracts/reports/CONTRACT-ALIGN-YYYYMMDD-###.md`
- `memx_spec_v3/docs/contracts/reports/LATEST.md`
- 契約整合の修正差分

### 次Phase遷移条件
- 全章で gate 判定 `pass`。
- high差分 0 件、REQ網羅率 100%、リンク不達 0 件、Birdseye issue 0 件。

## Phase 4（受け入れレビュー）
### Entry criteria
- Phase 3 gate `pass`。
- レビュー対象章と証跡ファイルが確定済み。

### 入力成果物
- Phase 3 の章別検証サマリ
- `memx_spec_v3/docs/design-review-spec.md`
- `memx_spec_v3/docs/design-acceptance-report-spec.md`
- `memx_spec_v3/docs/reviews/` 配下のレビュー記録

### 必須チェック
- 受け入れレポートに必須6項目（対象章/REQ網羅率/high差分件数/リンク不達件数/Birdseye issue件数/最終判定）を記載する。
- 最終判定前に gate 判定列（5軸）を再計算する。
- `CHANGELOG.md` / `memx_spec_v3/CHANGES.md` への反映要否を確定する。

### fail条件
- 必須6項目の欠落。
- gate 判定列のいずれかが `high`。
- 最終判定が `fail` なのに `Status: done` へ遷移した場合。

### 出力成果物
- `memx_spec_v3/docs/reviews/DESIGN-REVIEW-YYYYMMDD-XXX.md`
- `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md`
- CHANGES 反映メモ（Task Seed 側 `Moved-to-CHANGES` 含む）

### 次Phase遷移条件
- 最終 gate 判定 `pass`。
- `docs/TASKS.md` の完了条件（Release Note Draft / Moved-to-CHANGES）を満たす。
