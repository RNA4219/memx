# IN-BASELINE（初期運用ベースライン）

- 最終更新日（UTC）: `2026-03-03`
- 対象: memx v3 運用

## 1. 現時点の既知インシデント有無

- 既知インシデント: **なし**
- 補足: 新規発生時は `docs/INCIDENT_TEMPLATE.md` を基に `docs/IN-<YYYYMMDD>-<連番>.md` を起票する。

## 2. 起票トリガ（どの指標逸脱で起票するか）

`governance/metrics.yaml` の必須指標で breach が発生した場合に起票する。

- `response_time` の閾値超過（ingest/search/show のいずれかで `threshold` を超過）
- `compatibility` の閾値未達
- `error_classification` の逸脱
- `recall_threshold` の閾値未達

※ breach の語彙は `governance/metrics.yaml` と同一（`threshold` / `action_on_breach`）で運用する。

## 3. 起票責任者と SLA

- 起票責任者: 当番オペレーション担当（Primary On-call）
- 代理起票者: memx-core レビュー担当
- 起票 SLA:
  - SEV1: 検知から **15分以内**
  - SEV2: 検知から **30分以内**
  - SEV3: 検知から **4時間以内**
  - SEV4: 検知から **1営業日以内**
