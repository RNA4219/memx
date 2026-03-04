# Design Docs Prioritization Spec

## 1. 目的
- 設計書作成タスク（Task Seed）の `priority` を、`governance/prioritization.yaml` の重み付き評価方針に整合する判定軸で標準化する。
- 判定軸は設計書オーサリングで必須となる次の4項目に限定する。
  - Blocker有無
  - REQ網羅率への影響
  - 契約差分 high 件数
  - Birdseye issue の有無

## 2. 評価軸定義

HUB入力カバレッジ判定の根拠は、検索キー `Incident` / `Orchestration` / `TASK` と対象パス `docs/IN-*.md` / `orchestration/*.md` / `TASK.*` に固定する。

| 評価軸 | 判定条件 | high | medium | low |
| --- | --- | --- | --- | --- |
| Blocker有無 | Task未実施時に後続Phaseやレビューが停止するか | 停止する（blocked発生） | 一部遅延するが並行可能 | 停止しない |
| REQ網羅率への影響 | Task未実施時の要件ID網羅率低下度合い | 100%未達が確定 | 低下の可能性あり | 影響なし |
| 契約差分 high 件数 | high差分（契約不整合）への寄与 | high差分の新規/未解消が1件以上 | medium/low差分のみ | 差分なし |
| Birdseye issue の有無 | Birdseye検証 issue の有無と影響 | issueあり（node_id参照切れ・caps欠落等） | 軽微issueあり（文言/リンク揺れ） | issueなし |

## 3. 優先度決定ルール
1. 次のいずれかを満たす場合、`priority: high`。
   - Blocker有無 = high
   - REQ網羅率への影響 = high
   - 契約差分 high 件数 = high
   - Birdseye issue の有無 = high
2. 上記に該当せず、4軸のうち1つ以上が medium の場合、`priority: medium`。
3. 4軸すべてが low の場合のみ、`priority: low`。

4. HUB入力漏れ（上記固定キー/固定パスでの未充足）がある場合、4軸評価への影響として優先度を引き上げる。
   - HUB入力漏れがある場合の最低優先度は `medium`。
   - `Incident` または `Orchestration` の漏れを含む場合は `priority: high`。

## 4. 運用手順
1. Task Seed 起票時に4軸を判定し、`Requirements` か `Dependencies` に根拠を1行ずつ残す（加えて固定キー/固定パスで HUB入力漏れ有無を確認する）。
2. `priority` は本仕様のルールで決定し、`docs/TASKS.md` の書式に従って記載する。
3. Phase完了判定（orchestration）ごとに再評価し、軸が変化した場合は `priority` を更新する。
4. `status: blocked` へ遷移した Task は、再開時に4軸を再判定してから `planned/active` へ戻す。

## 5. governance/prioritization.yaml との整合
- `governance/prioritization.yaml` の `impact/urgency/risk/recovery_cost` は、設計書タスクでは次の対応で解釈する。
  - impact: REQ網羅率への影響 + 契約差分 high 件数
  - urgency: Blocker有無
  - risk: Birdseye issue の有無 + 契約差分 high 件数
  - recovery_cost: 修正対象章数・再レビュー工数（Task Seed の工数見積で補足）
- 本仕様はTask Seed用の簡易判定であり、詳細スコアリングが必要な場合は `governance/prioritization.yaml` を正本として併用する。

## 6. Phase Gate 判定列への取り込み
- `memx_spec_v3/docs/design-phase-gate-spec.md` の gate 判定では、本書4軸を次の列名で必ず記録する。
  - `gate_blocker`（Blocker有無）
  - `gate_req_coverage`（REQ網羅率への影響）
  - `gate_contract_high`（契約差分 high 件数）
  - `gate_birdseye_issue`（Birdseye issue の有無）
- 各列の値は `high/medium/low` の3値のみ許可する。
- gate の総合判定は「high が1つでもあれば fail、high なしで medium があれば hold、全列 low のみ pass」とする。
- Task Seed の `priority` は gate 判定列の値を流用して再計算し、Phase 境界で差分が出た場合は更新する。
- HUB入力漏れが検出された場合は本書 3.4 を優先して適用し、最低 `medium`（Incident/Orchestration漏れは `high`）へ補正する。
