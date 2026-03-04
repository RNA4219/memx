# INCIDENT TEMPLATE

- インシデントID: `IN-<YYYYMMDD>-<連番3桁>`
- Task Seed 参照ID（必須）: `IN-<YYYYMMDD>-<連番3桁>`（Task Seed の Source/Requirements から逆引きするため、インシデントIDと同値で記載）
- 発生日（UTC）:
- 起票日（UTC）:
- 重大度: `SEV1 | SEV2 | SEV3 | SEV4`
- ステータス: `Open | Mitigated | Resolved | Closed`
- 関連要件: [`memx 要件定義: 0.目的とスコープ`](../memx_spec_v3/docs/requirements.md#0-目的とスコープ)

## 最小監査項目（REQ-NFR-006）

- 適用要件: `REQ-NFR-006`
- 必須記録項目:
  - 事象識別子（ID/発生日/起票日/重大度/ステータス）
  - 要件トレーサビリティ（関連要件ID/要件違反有無/違反要件）
  - 時間監査（検知/暫定復旧/恒久復旧）
  - 復旧行動監査（再試行回数/ロールバック実施有無/再計画チケットID）
  - 影響監査（影響対象/期間/規模/CIA）
  - 証跡ファイル保存先（ログ/メトリクス/判定結果）

## waiver運用時の必須記録項目（REQ-NFR-006 連動）

- 運用ルール: waiver 記録は必ず `docs/IN-<実日付>-<連番>.md` に残し、`EVALUATION.md` の運用NFR合否判定で証跡として参照可能な状態にする。
- `docs/IN-<実日付>-<連番>.md` に必須記録すること:
  1. waiver対象要件ID（例: `REQ-NFR-005`）
  2. waiver理由（技術的制約/外部依存/緊急運用の別）
  3. 期限（UTC、失効日時）
  4. 暫定リスク受容者（承認者）
  5. 代替統制（監視強化・手動運用手順・追加検証）
  6. 解除条件（どの証跡が揃えば waiver を解消するか）
  7. 関連証跡パス（`artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson` など）

## 1. 検知（Detection）

- 検知日時:
- 検知経路（監視/ユーザー報告/レビュー等）:
- 初動担当:
- 事象概要:
- 暫定復旧完了日時:
- 恒久復旧完了日時:
- 要件違反の有無: `Yes | No`
- 違反した要件ID/節（例: `requirements.md#2-7-1`）:

## 2. 影響（Impact）

- 影響対象（機能/データ/利用者）:
- 対象データ分類（`public | internal | secret`）:
- 漏えい有無: `有 | 無 | 調査中`
- 影響期間:
- 影響規模（件数・範囲）:
- 再試行回数:
- ロールバック実施有無: `Yes | No`
- 再計画チケットID:
- 機密性/完全性/可用性への影響:

## 3. 原因分析（5 Whys）

1. Why1:
2. Why2:
3. Why3:
4. Why4:
5. Why5:

## 4. 再発防止策（Preventive Actions）

- 再発防止項目ID（必須・複数可）: `PA-01`, `PA-02` ...（Task Seed Requirements から参照するID）
- 恒久対策:
- 検知強化:
- 運用改善:
- オーナー:
- 期限:

## 5. タイムライン（Timeline）

| 時刻（UTC） | 事実 | 判断/対応 | 担当 |
| --- | --- | --- | --- |
|  |  |  |  |


## 6. 証跡（Evidence）

- 証跡ファイル:
  - `artifacts/ops/incident-summary.json`
  - `artifacts/ops/recovery-log.ndjson`
  - その他関連ログ/メトリクス:
