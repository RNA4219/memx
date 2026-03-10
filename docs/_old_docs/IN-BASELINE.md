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


## 4. `docs/IN-*.md` 最小監査項目（固定要件）

- 適用要件: `REQ-NFR-006`
- すべての `docs/IN-*.md` は次を必須記載とする。
  1. 事象識別子（ID/発生日/起票日/重大度/ステータス）
  2. 要件トレーサビリティ（関連要件ID/要件違反有無/違反要件）
  3. 時間監査（検知/暫定復旧/恒久復旧）
  4. 復旧行動監査（再試行回数/ロールバック実施有無/再計画チケットID）
  5. 影響監査（影響対象/期間/規模/CIA）
  6. 証跡ファイル保存先（`artifacts/ops/incident-summary.json` / `artifacts/ops/recovery-log.ndjson` など）

## waiver運用時の必須記録項目（REQ-NFR-006 連動）

- 運用ルール: waiver 記録は必ず `docs/IN-<実日付>-<連番>.md` に残し、`EVALUATION.md` の運用NFR合否判定で証跡として参照可能な状態にする。
- 対象: `REQ-NFR-001` / `REQ-RET-001` / 運用NFR（`REQ-NFR-002`〜`REQ-NFR-005`）で fail を一時許容する場合。
- `docs/IN-<実日付>-<連番>.md` に必須記録すること:
  1. waiver対象要件ID（例: `REQ-NFR-005`）
  2. waiver理由（技術的制約/外部依存/緊急運用の別）
  3. 期限（UTC、失効日時）
  4. 暫定リスク受容者（承認者）
  5. 代替統制（監視強化・手動運用手順・追加検証）
  6. 解除条件（どの証跡が揃えば waiver を解消するか）
  7. 関連証跡パス（`artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson` など）
