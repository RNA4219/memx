---
priority_phase_1: high
priority_phase_2: high
priority_phase_3: high
priority_phase_4: medium
owner: memx-core
deadline: 2026-03-31
status: planned
---

# memx Design Docs Authoring Orchestration

## Phase 1: 情報収集
### Dependencies
- `requirements.md`
- `traceability.md`
- `design.md`
- `interfaces.md`
- `EVALUATION.md`
- `RUNBOOK.md`
- `docs/birdseye/index.json`

- [ ] `requirements.md` から要件ID一覧を抽出し、重複・欠番を洗い出す（Task Seed 1件、<=0.5d）
- [ ] `traceability.md` を入力成果物として主要 REQ-ID の設計/I/F/評価/契約マッピングを確認する（Task Seed 1件、<=0.5d）
- [ ] `design.md` の章見出しと要件IDの参照有無を対応表にする（Task Seed 1件、<=0.5d）
- [ ] `interfaces.md` の入出力契約を要件ID単位で列挙する（Task Seed 1件、<=0.5d）
- [ ] `EVALUATION.md` の評価観点を要件IDに紐づける（Task Seed 1件、<=0.5d）
- [ ] `RUNBOOK.md` の運用手順と契約依存箇所を抽出する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/index.json` から node_id と depends_on を取得し、章候補との対応表を作る（Task Seed 1件、<=0.5d）

### Done Criteria
- 情報源7ファイルの抽出結果が要件ID/契約ID/node_id単位で一覧化されている
- 抽出結果を Task Seed 化可能な粒度（1項目=1タスク、<=0.5d）で分割している
- `docs/TASKS.md` 必須項目のうち `Source` / `Node IDs` / `Requirements` へ直接転記できる状態になっている

## Phase 2: 章別ドラフト
### Dependencies
- `requirements.md`
- `design.md`
- `interfaces.md`
- `EVALUATION.md`
- `RUNBOOK.md`
- `docs/birdseye/index.json`

- [ ] 章ごとに `Objective`（1〜3行）を作成する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Source` を `path#Section` 形式で記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Node IDs` を `docs/birdseye/index.json` 参照で付与する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Requirements` を要件ID箇条書きで記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Commands`（仕様整合チェック用）を記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Dependencies` と `Status: planned` を記述する（Task Seed 1件/章、<=0.5d）

### Done Criteria
- 全章ドラフトが `Source/Node IDs/Objective/Requirements/Commands/Dependencies/Status` を満たす
- 各章が HUB ノード抽出ルールでそのまま Task Seed 化できる
- 章ごとの未解決事項が 0.5d 以内の追加タスクへ分解済みである

## Phase 3: 契約整合
### Dependencies
- `requirements.md`
- `traceability.md`
- `design.md`
- `interfaces.md`
- `EVALUATION.md`
- `RUNBOOK.md`
- `docs/birdseye/index.json`

- [ ] 要件ID網羅率を算出し、章別ドラフトの欠落IDを補完する（Task Seed 1件、<=0.5d）
- [ ] `traceability.md` を入力成果物として参照し、要件ID網羅率と契約対応の欠落を補完する（Task Seed 1件、<=0.5d）
- [ ] `design.md` と `interfaces.md` の契約差分を比較し、相違を解消する（Task Seed 1件、<=0.5d）
- [ ] `RUNBOOK.md` の手順リンクが契約記述と一致するか確認する（Task Seed 1件、<=0.5d）
- [ ] `EVALUATION.md` の評価項目リンクが最新章へ到達するか確認する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/index.json` の node_id 参照切れを修正する（Task Seed 1件、<=0.5d）
- [ ] 章間リンクの相対パス/アンカー健全性を確認する（Task Seed 1件、<=0.5d）

### Done Criteria
- 仕様整合チェックを満たす（要件ID網羅率 100%）
- 仕様整合チェックを満たす（契約同期: `design.md` と `interfaces.md` の差分 0 件）
- 仕様整合チェックを満たす（リンク健全性: 章内/章間/運用リンクの不達 0 件）
- 仕様整合チェック結果を各章 Task Seed の `Requirements` と `Commands` に反映済みである

## Phase 4: 受け入れレビュー
### Dependencies
- `requirements.md`
- `design.md`
- `interfaces.md`
- `EVALUATION.md`
- `RUNBOOK.md`
- `docs/birdseye/index.json`

- [ ] 章ごとに `Release Note Draft`（1〜3行）を作成する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとにレビューコメントを反映し再チェックする（Task Seed 1件/章、<=0.5d）
- [ ] 仕様整合チェック結果（要件ID網羅率/契約同期/リンク健全性）をレビュー記録に添付する（Task Seed 1件、<=0.5d）
- [ ] `Status` を `reviewing` から `done` へ更新する条件を確認する（Task Seed 1件、<=0.5d）
- [ ] `Moved-to-CHANGES: YYYY-MM-DD` の追記対象を確定する（Task Seed 1件、<=0.5d）

### Done Criteria
- 全章が `docs/TASKS.md` 必須項目フォーマットへマッピング済みである
  - `Source`: `path#Section`
  - `Node IDs`: `docs/birdseye/index.json` の node_id
  - `Objective`: 1〜3行
  - `Requirements`: 要件IDと整合条件
  - `Commands`: 仕様整合チェック実行手順
  - `Dependencies`: 前提タスク/外部条件
  - `Release Note Draft`: 利用者影響の要約
  - `Status`: 許可語彙（planned/active/in_progress/reviewing/blocked/done）
- `lint/type/test` ではなく、仕様整合チェック（要件ID網羅率・契約同期・リンク健全性）で受け入れ判定する運用が明文化されている
- Phase 1〜4 のチェックボックスがすべて完了している
