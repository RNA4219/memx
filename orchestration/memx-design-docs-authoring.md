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

> 注記: Dependencies では非正規名の記載を許可するが、Task Seed / Phase 1 抽出表 / 章ドラフト / レビュー記録へ転記する時点で `memx_spec_v3/docs/design-reference-resolution-spec.md` に従い正規パスへ正規化すること。

- [ ] Birdseye 検証を `docs/birdseye/memx-birdseye-validation-spec.md` に従って実行し、caps 実体欠落を issue 出力する（Task Seed 1件、<=0.5d）
- [ ] `requirements.md` から要件ID一覧を抽出し、重複・欠番を洗い出す（Task Seed 1件、<=0.5d）
- [ ] `traceability.md` を入力成果物として主要 REQ-ID の設計/I/F/評価/契約マッピングを確認する（Task Seed 1件、<=0.5d）
- [ ] `design.md` の章見出しと要件IDの参照有無を対応表にする（Task Seed 1件、<=0.5d）
- [ ] `interfaces.md` の入出力契約を要件ID単位で列挙する（Task Seed 1件、<=0.5d）
- [ ] `EVALUATION.md` の評価観点を要件IDに紐づける（Task Seed 1件、<=0.5d）
- [ ] `RUNBOOK.md` の運用手順と契約依存箇所を抽出する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/index.json` から node_id と depends_on を取得し、章候補との対応表を作る（Task Seed 1件、<=0.5d）

### Done Criteria
- 情報源7ファイルの抽出結果が要件ID/契約ID/node_id単位で一覧化されている（参照仕様: `../memx_spec_v3/docs/design-source-inventory-spec.md`）
- 抽出結果を Task Seed 化可能な粒度（1項目=1タスク、<=0.5d）で分割している
- `docs/TASKS.md` 必須項目のうち `Source` / `Node IDs` / `Requirements` へ直接転記できる状態になっている

### Phase 1 参照先追記計画（情報源7ファイルの抽出結果）
- [ ] `design-source-inventory-spec.md` の必須列（`source_path#section`, `req_id`, `contract_ref`, `node_id`, `depends_on`, `owner`, `reviewed_at`）で抽出表を作成する
- [ ] 抽出表で `blocked` 条件と差し戻し条件を判定し、未解決行を 0 件にする
- [ ] Phase 1 Done Criteria の判定時に、上記仕様書への準拠確認をチェック項目として記録する

## Phase 2: 章別ドラフト
### Dependencies
- `requirements.md`
- `design.md`
- `interfaces.md`
- `EVALUATION.md`
- `RUNBOOK.md`
- `docs/birdseye/index.json`
- `memx_spec_v3/docs/design-acceptance-report-spec.md`

- [ ] 章ごとに `Objective`（1〜3行）を作成する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Source` を `path#Section` 形式で記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Node IDs` を `docs/birdseye/index.json` 参照で付与する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Requirements` を要件ID箇条書きで記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Commands`（仕様整合チェック用）を記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Dependencies` と `Status: planned` を記述する（Task Seed 1件/章、<=0.5d）

### Done Criteria
- 全章ドラフトが `Source/Node IDs/Objective/Requirements/Commands/Dependencies/Status` を満たし、`memx_spec_v3/docs/design-template.md` に準拠している
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

> 注記: Dependencies では非正規名の記載を許可するが、Task Seed / Phase 1 抽出表 / 章ドラフト / レビュー記録へ転記する時点で `memx_spec_v3/docs/design-reference-resolution-spec.md` に従い正規パスへ正規化すること。

- [ ] 要件ID網羅率を算出し、章別ドラフトの欠落IDを補完する（Task Seed 1件、<=0.5d）
- [ ] `design.md` と `interfaces.md` の契約差分を比較し、相違を解消する（`memx_spec_v3/docs/contract-alignment-spec.md` に従う）（Task Seed 1件、<=0.5d）
- [ ] `traceability.md` を入力成果物として参照し、要件ID網羅率と契約対応の欠落を補完する（Task Seed 1件、<=0.5d）
- [ ] `design.md` と `interfaces.md` の契約差分を比較し、相違を解消する（Task Seed 1件、<=0.5d）
- [ ] `RUNBOOK.md` の手順リンクが契約記述と一致するか確認する（Task Seed 1件、<=0.5d）
- [ ] `EVALUATION.md` の評価項目リンクが最新章へ到達するか確認する（Task Seed 1件、<=0.5d）
- [ ] Birdseye 検証 issue（node_id 参照切れを含む）を `docs/birdseye/memx-birdseye-validation-spec.md` に従って修正する（Task Seed 1件、<=0.5d）
- [ ] 章間リンクの相対パス/アンカー健全性を確認する（Task Seed 1件、<=0.5d）
- [ ] Phase 3 Done Criteria 判定時に `memx_spec_v3/docs/link-integrity-spec.md` を参照する運用タスクを登録し、章別 Task Seed へ反映する（Task Seed 1件、<=0.5d）

### Done Criteria
- 仕様整合チェックを満たす（要件ID網羅率 100%）
- 仕様整合チェックを満たす（契約同期: `memx_spec_v3/docs/contract-alignment-spec.md` に基づく判定で high=0 件）
- 仕様整合チェックを満たす（リンク健全性: `memx_spec_v3/docs/link-integrity-spec.md` に基づき、章内/章間/運用リンクの不達 0 件）
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
- [ ] レビュー記録は [テンプレート](../memx_spec_v3/docs/reviews/TEMPLATE.md) から作成し、[保存先ルール](../memx_spec_v3/docs/reviews/README.md) に従って `memx_spec_v3/docs/reviews/` へ保存する（Task Seed 1件、<=0.5d）
- [ ] 章ごとにレビューコメントを反映し再チェックする（Task Seed 1件/章、<=0.5d）
- [ ] 仕様整合チェック結果（要件ID網羅率/契約同期/リンク健全性）をレビュー記録に添付する（Task Seed 1件、<=0.5d）
- [ ] 統合受け入れレポートを `memx_spec_v3/docs/design-acceptance-report-spec.md` に従って作成する（Task Seed 1件、<=0.5d）
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
- 全章のレビュー記録が `memx_spec_v3/docs/design-review-spec.md` に準拠している
  - 保存先: `memx_spec_v3/docs/reviews/`
  - 命名: `DESIGN-REVIEW-YYYYMMDD-###.md`
  - 必須項目: 対象章 / 関連 REQ-ID / Node IDs / 指摘一覧（重大度付き） / 再確認結果 / 判定（pass/fail/waiver）
  - 判定根拠: `EVALUATION.md` pass/fail ルール参照を必須化
  - 完了条件: `docs/TASKS.md` の `Release Note Draft` / `Status` / `Moved-to-CHANGES` を確認済み
- 統合受け入れレポートが `memx_spec_v3/docs/design-acceptance-report-spec.md` に準拠している
  - 保存先: `memx_spec_v3/docs/reviews/`
  - 命名: `DESIGN-ACCEPTANCE-YYYYMMDD.md`
  - 必須項目: 対象章 / REQ網羅率 / high差分件数 / リンク不達件数 / Birdseye issue件数 / 最終判定
  - 判定規則: `high>0` または `REQ網羅率<100%` または `リンク不達件数>0` または `Birdseye issue件数>0` で fail
- `lint/type/test` ではなく、仕様整合チェック（要件ID網羅率・契約同期・リンク健全性）で受け入れ判定する運用が明文化されている
- Phase 1〜4 のチェックボックスがすべて完了している
