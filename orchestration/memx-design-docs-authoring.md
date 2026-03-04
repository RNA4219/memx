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

## 判定責務の分離
- 本書は Phase 1〜4 の実施手順（チェック実行順・タスク分解・証跡作成）を定義する。
- gate の entry/exit criteria・fail 条件・次Phase遷移条件の正本は `memx_spec_v3/docs/design-phase-gate-spec.md` とする。
- gate 判定に用いる5軸（Blocker / REQ網羅率 / 契約差分 high / Birdseye issue / HUB入力カバレッジ）は `docs/design-docs-prioritization-spec.md` と整合させる。

## Phase 1: 情報収集
### Preprocessing: トリガー判定
- [ ] `memx_spec_v3/docs/design-update-trigger-spec.md` を参照し、今回変更の Trigger IDs（TRG-REQ / TRG-OAS / TRG-CLI-SCHEMA / TRG-RUNBOOK / TRG-IN-PREVENTIVE）を先に確定する（Task Seed 1件、<=0.5d）
- [ ] Trigger IDs に対応する必須更新先（design/interfaces/traceability/EVALUATION/operations/レビュー記録）を Phase 1 で更新対象として固定する（Task Seed 1件、<=0.5d）
- [ ] `CHANGELOG.md` / `memx_spec_v3/CHANGES.md` の反映要否を前処理で判定し、Task Seed の `Release Note Draft` / `Status: done` 条件へ反映する（Task Seed 1件、<=0.5d）

### Priority Label Rule
- Phase 1 の各チェック項目は `docs/design-docs-prioritization-spec.md` の4軸で判定する。
- `high`: Blocker発生、REQ網羅率100%未達確定、Birdseye issue（caps欠落/node_id参照切れ）あり。
- `medium`: Blockerなしだが REQ網羅率低下の可能性または軽微 Birdseye issue がある。
- `low`: 4軸すべて低リスクで、情報整理のみ。
### Dependencies
- `requirements.md`
- `traceability.md`
- `design.md`
- `interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`
- `docs/IN-*.md`（実績インシデントのみ。テンプレート除外）
- `orchestration/*.md`
- `docs/INCIDENT_TEMPLATE.md`（必要時の参照のみ。実績証跡扱い不可）
- `docs/incident-record-operations-spec.md`（Incident source evidence 採用可否条件の正本）

> 注記: Dependencies では非正規名の記載を許可するが、Task Seed / Phase 1 抽出表 / 章ドラフト / レビュー記録へ転記する時点で `memx_spec_v3/docs/design-reference-resolution-spec.md` に従い正規パスへ正規化すること。
> 正規化例（Task Seed転記時）: `requirements.md` → `memx_spec_v3/docs/requirements.md#<section>`、`CONTRACTS.md` → `memx_spec_v3/docs/CONTRACTS.md#<section>`、`EVALUATION.md` → `docs/birdseye/caps/EVALUATION.md.json#<section_key>`。
> 曖昧名（例: `requirements.md` / `contracts.md` / `EVALUATION.md` のみ）を残したまま Task Seed を起票してはならない。

- [ ] Birdseye 検証を `docs/birdseye/memx-birdseye-validation-spec.md` に従って実行し、caps 実体欠落を issue 出力する（Task Seed 1件、<=0.5d）
- [ ] `requirements.md` から要件ID一覧を抽出し、重複・欠番を洗い出す（Task Seed 1件、<=0.5d）
- [ ] `traceability.md` を入力成果物として主要 REQ-ID の設計/I/F/評価/契約マッピングを確認する（Task Seed 1件、<=0.5d）
- [ ] `design.md` の章見出しと要件IDの参照有無を対応表にする（Task Seed 1件、<=0.5d）
- [ ] `interfaces.md` の入出力契約を要件ID単位で列挙する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/caps/EVALUATION.md.json` の評価観点を要件IDに紐づける（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/caps/RUNBOOK.md.json` の運用手順と契約依存箇所を抽出する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/index.json` から node_id と depends_on を取得し、`memx_spec_v3/docs/design-chapter-node-mapping-spec.md` 準拠の章対応表（chapter_id -> node_id）を更新する（Task Seed 1件、<=0.5d）
- [ ] `docs/IN-*.md`（実績インシデント）から再発防止要件・運用証跡を抽出し、テンプレート（`docs/INCIDENT_TEMPLATE.md`）混入を除外する（Task Seed 1件、<=0.5d）
- [ ] Incident source evidence は `docs/incident-record-operations-spec.md` に従い、必須項目（ID/時刻/要件ID/waiver項目/証跡パス）充足時のみ採用する（Task Seed 1件、<=0.5d）
- [ ] `orchestration/*.md` の依存関係（depends_on 相当）を抽出し、未解決依存を `blocked` 条件として明示する（Task Seed 1件、<=0.5d）
- [ ] `memx_spec_v3/docs/design-reference-conformance-spec.md` に従い、Phase 1 抽出表の参照解決適合（`path#section` / warn/fail 判定）を確認する（Task Seed 1件、<=0.5d）

### Done Criteria
- Phase 1 Done Criteria は `../memx_spec_v3/docs/design-source-inventory-spec.md` と `../memx_spec_v3/docs/design-chapter-node-mapping-spec.md` を正本として判定する（情報源7ファイル一覧化、Task Seed 粒度、`docs/TASKS.md` 転記可否、node 解決成否を含む）
- `gate_hub_source_coverage`（`high/medium/low`）を必須入力として記録し、判定ロジック・写像規則・証跡要件は `memx_spec_v3/docs/design-hub-source-coverage-spec.md` を正本として適用する。あわせて Incident 転記完了チェック（`memx_spec_v3/docs/incident-to-task-traceability-spec.md` 準拠で `Requirements` / `Commands` / `Dependencies` への転記完了、かつ検証コマンド1件以上）を必須化する
- 必須成果物の対象範囲と更新タイミングは `memx_spec_v3/docs/design-deliverables-package-spec.md` を正本として適用する

### Phase 1 参照先追記計画（情報源（固定入力+拡張入力）の抽出結果）
- [ ] `design-source-inventory-spec.md` の必須列（`source_type`, `source_path#section`, `req_id`, `contract_ref`, `node_id`, `depends_on`, `owner`, `reviewed_at`, `node_resolution_status`）で抽出表を作成する
- [ ] `design-source-inventory-spec.md` と `design-chapter-node-mapping-spec.md` の fail 条件（`ambiguous` / `missing` 含む）を判定し、未解決行を 0 件にする
- [ ] Phase 1 Done Criteria 判定ログは上記2仕様への準拠確認のみを正本として記録する

## Phase 2: 章別ドラフト
### Priority Label Rule
- Phase 2 の各チェック項目は章ごとに4軸判定し、Task Seed の front matter `priority` へ反映する。
- `high`: 当該章の未作成により Phase 3 で Blocker 化する、または REQ網羅率100%未達が確定する。
- `medium`: 章ドラフトはあるが契約差分 high 予兆や軽微 Birdseye issue が残る。
- `low`: 追記のみで REQ網羅率・契約差分・Birdseye issue に影響しない。
### Dependencies
- `requirements.md`
- `design.md`
- `interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`
- `memx_spec_v3/docs/design-acceptance-report-spec.md`

- [ ] 章ごとに `Objective`（1〜3行）を作成する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Source` を `path#Section` 形式で記述する（Task Seed 1件/章、<=0.5d）
- [ ] `memx_spec_v3/docs/design-reference-conformance-spec.md` に従い、章ドラフト `Source` の参照解決適合（`path#section` / warn/fail 判定）を確認する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Node IDs` を `docs/birdseye/index.json` 参照で付与し、章対応表（chapter_id -> node_id）を更新する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Requirements` を要件ID箇条書きで記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Commands`（仕様整合チェック用）を記述する（Task Seed 1件/章、<=0.5d）
- [ ] 章ごとに `Dependencies` と `Status: planned` を記述する（Task Seed 1件/章、<=0.5d）

### Done Criteria
- 全章ドラフトが `Source/Node IDs/Objective/Requirements/Commands/Dependencies/Status` を満たし、`memx_spec_v3/docs/design-template.md` に準拠している
- 章対応表（`memx_spec_v3/docs/design-chapter-node-mapping-spec.md`）が Phase 2 更新内容に追随している
- 各章が HUB ノード抽出ルールでそのまま Task Seed 化できる
- 章ごとの未解決事項が 0.5d 以内の追加タスクへ分解済みである
- 必須成果物の対象範囲と更新タイミングは `memx_spec_v3/docs/design-deliverables-package-spec.md` を正本として適用する

## Phase 3: 契約整合
### Priority Label Rule
- Phase 3 の各チェック項目は契約差分 high 件数を最優先軸として判定する。
- `high`: 契約差分 high が1件以上、REQ網羅率100%未達、または Birdseye issue 未解消。
- `medium`: high差分はないが medium/low差分またはリンク修正が残る。
- `low`: 差分なしで確認作業のみ。
### Dependencies
- `requirements.md`
- `traceability.md`
- `design.md`
- `interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`

> 注記: Dependencies では非正規名の記載を許可するが、Task Seed / Phase 1 抽出表 / 章ドラフト / レビュー記録へ転記する時点で `memx_spec_v3/docs/design-reference-resolution-spec.md` に従い正規パスへ正規化すること。
> 正規化例（Task Seed転記時）: `requirements.md` → `memx_spec_v3/docs/requirements.md#<section>`、`CONTRACTS.md` → `memx_spec_v3/docs/CONTRACTS.md#<section>`、`EVALUATION.md` → `docs/birdseye/caps/EVALUATION.md.json#<section_key>`。
> 曖昧名（例: `requirements.md` / `contracts.md` / `EVALUATION.md` のみ）を残したまま Task Seed を起票してはならない。

- [ ] 要件ID網羅率を算出し、章別ドラフトの欠落IDを補完する（Task Seed 1件、<=0.5d）
- [ ] `design.md` と `interfaces.md` の契約差分を比較し、相違を解消する（`memx_spec_v3/docs/contract-alignment-spec.md` に従う）（Task Seed 1件、<=0.5d）
- [ ] `traceability.md` を入力成果物として参照し、要件ID網羅率と契約対応の欠落を補完する（Task Seed 1件、<=0.5d）
- [ ] `design.md` と `interfaces.md` の契約差分を比較し、相違を解消する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/caps/RUNBOOK.md.json` の手順リンクが契約記述と一致するか確認する（Task Seed 1件、<=0.5d）
- [ ] `docs/birdseye/caps/EVALUATION.md.json` の評価項目リンクが最新章へ到達するか確認する（Task Seed 1件、<=0.5d）
- [ ] Birdseye 検証 issue（node_id 参照切れを含む）を `docs/birdseye/memx-birdseye-validation-spec.md` に従って修正する（Task Seed 1件、<=0.5d）
- [ ] 章間リンクの相対パス/アンカー健全性を確認する（Task Seed 1件、<=0.5d）
- [ ] Phase 3 Done Criteria 判定時に `memx_spec_v3/docs/link-integrity-spec.md` を参照する運用タスクを登録し、章別 Task Seed へ反映する（Task Seed 1件、<=0.5d）

### Done Criteria
- 仕様整合チェックの完成判定は `memx_spec_v3/docs/design-doc-dod-spec.md` を正本として実施する
- 仕様整合チェック結果を各章 Task Seed の `Requirements` と `Commands` に反映済みである
- 各章の検証結果を `memx_spec_v3/docs/design-chapter-validation-spec.md` 準拠の章別検証サマリとして作成済みである（参照: `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304.md`）
- 必須成果物の対象範囲と更新タイミングは `memx_spec_v3/docs/design-deliverables-package-spec.md` を正本として適用する

### Phase 3 Done Criteria 追記案（本文編集は別タスク）
- 重複する個別判定文（REQ網羅率/契約同期/link 健全性）は、`memx_spec_v3/docs/design-doc-dod-spec.md#3. 完成判定ルール（固定）` 参照へ置換する。
- 置換後の判定文案: 「Phase 3 の完成判定は `memx_spec_v3/docs/design-doc-dod-spec.md` を正本として実施する。」

## Phase 4: 受け入れレビュー
### Priority Label Rule
- Phase 4 の各チェック項目は最終判定阻害の有無で付与する。
- `high`: `Status: done` への遷移を阻害する欠落（REQ網羅率<100%、high差分>0、Birdseye issue>0）がある。
- `medium`: 判定阻害はないがレビュー指摘反映や記録不足が残る。
- `low`: 記録整備・移送準備のみ。
### Dependencies
- `requirements.md`
- `design.md`
- `interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`
- `docs/incident-record-operations-spec.md`（Phase 4 source evidence 最終判定条件）

- [ ] 章ごとに `Release Note Draft`（1〜3行）を作成する（Task Seed 1件/章、<=0.5d）
- [ ] レビュー記録は [テンプレート](../memx_spec_v3/docs/reviews/TEMPLATE.md) から作成し、[保存先ルール](../memx_spec_v3/docs/reviews/README.md) に従って `memx_spec_v3/docs/reviews/` へ保存する（Task Seed 1件、<=0.5d）
- [ ] 章ごとにレビューコメントを反映し再チェックする（Task Seed 1件/章、<=0.5d）
- [ ] 仕様整合チェック結果（要件ID網羅率/契約同期/リンク健全性）をレビュー記録に添付する（Task Seed 1件、<=0.5d）
- [ ] 統合受け入れレポートを `memx_spec_v3/docs/design-acceptance-report-spec.md` に従って作成する（Task Seed 1件、<=0.5d）
- [ ] `Status` を `reviewing` から `done` へ更新する条件を確認する（Task Seed 1件、<=0.5d）
- [ ] `Moved-to-CHANGES: YYYY-MM-DD` の追記対象を確定する（Task Seed 1件、<=0.5d）

### Done Criteria
- 全章が `docs/TASKS.md` 必須項目フォーマットへマッピング済みである（`Release Note Draft` / `Status` / `Moved-to-CHANGES` を含む）
- 最終 gate 再計算時に `gate_hub_source_coverage`（`high/medium/low`）を必須入力として記録し、判定根拠は `docs/IN-*.md`・`orchestration/*.md`・`TASK.*` を対象に検索キー `Incident` / `Orchestration` / `TASK` で固定する
- Incident を source evidence として採用する場合、`docs/incident-record-operations-spec.md` 準拠で実体ファイル（`docs/IN-<実日付>-<連番>.md`）のみを許可し、テンプレート/ベースライン参照が残るレビューは `done` 遷移不可とする
- 全章のレビュー記録が `memx_spec_v3/docs/design-review-spec.md` に準拠している
  - 保存先: `memx_spec_v3/docs/reviews/`
  - 命名: `DESIGN-REVIEW-YYYYMMDD-###.md`
  - 必須項目: 対象章 / 関連 REQ-ID / Node IDs / 指摘一覧（重大度付き） / 再確認結果 / 判定（pass/fail/waiver）
  - 判定根拠: `docs/birdseye/caps/EVALUATION.md.json` pass/fail ルール参照を必須化
- 統合受け入れレポートが `memx_spec_v3/docs/design-acceptance-report-spec.md` と `memx_spec_v3/docs/design-acceptance-lifecycle-spec.md` に準拠している
  - 判定手順・入力元・差戻し・再作成・責務分担は `memx_spec_v3/docs/design-acceptance-operations-spec.md` を唯一の参照先として適用する
  - DA-LC-01: テンプレートと実体ファイルの責務分離
  - DA-LC-02: 実体ファイルの新規作成（上書き・テンプレート流用禁止）
  - DA-LC-03: 実体命名規則 `DESIGN-ACCEPTANCE-<実日付>.md`
  - DA-LC-04: 作成タイミング（Phase 4 entry / 再レビュー / 最終判定）
  - DA-LC-05: 差戻し条件（必須6項目欠落 / `evidence_paths` 不整合 / gate列 `high`）
- 章別検証サマリ（`memx_spec_v3/docs/design-chapter-validation-spec.md`）が作成済みで、レビュー記録・受け入れレポートの参照を添付済みである（参照: `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304.md`）
- `lint/type/test` ではなく、仕様整合チェック（要件ID網羅率・契約同期・リンク健全性）で受け入れ判定する運用が明文化されている
- Phase 1〜4 のチェックボックスがすべて完了している
- 必須成果物の対象範囲と更新タイミングは `memx_spec_v3/docs/design-deliverables-package-spec.md` を正本として適用する

### Phase 4 Done Criteria 追記案（本文編集は別タスク）
- 「REQ網羅率/high差分/リンク不達/Birdseye issue」の閾値判定文は `memx_spec_v3/docs/design-doc-dod-spec.md#3. 完成判定ルール（固定）` へ集約参照する。
- 置換後の判定文案: 「最終判定は `memx_spec_v3/docs/design-doc-dod-spec.md` を唯一の正本として `pass` / `fail` を決定する。」
