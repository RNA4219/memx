# IN-YYYYMMDD-001（雛形インシデント）

> **テンプレート専用ファイル**: このファイルは記入例付きの雛形です。実インシデントの記録には使用せず、`docs/IN-<実日付>-<連番>.md` を新規作成して管理してください。

- インシデントID: `IN-YYYYMMDD-001`
- 発生日（UTC）: `YYYY-MM-DD`
- 起票日（UTC）: `YYYY-MM-DD`
- 重大度: `SEV3`
- ステータス: `Open`
- 関連要件: [`memx 要件定義: 0.目的とスコープ`](../memx_spec_v3/docs/requirements.md#0-目的とスコープ)
- 使用テンプレート: [`docs/INCIDENT_TEMPLATE.md`](./INCIDENT_TEMPLATE.md)

## 最小監査項目（REQ-NFR-006）

- 適用要件: `REQ-NFR-006`
- 必須記録項目:
  - 事象識別子（ID/発生日/起票日/重大度/ステータス）
  - 要件トレーサビリティ（関連要件ID/要件違反有無/違反要件）
  - 時間監査（検知/暫定復旧/恒久復旧）
  - 復旧行動監査（再試行回数/ロールバック実施有無/再計画チケットID）
  - 影響監査（影響対象/期間/規模/CIA）
  - 証跡ファイル保存先（ログ/メトリクス/判定結果）

## 1. 検知（Detection）

- 検知日時: `YYYY-MM-DDThh:mm:ssZ`
- 検知経路: 例）ユーザー報告
- 初動担当: `TBD`
- 事象概要: 例）誤ったメタデータ設定により検索結果の一部が欠落
- 暫定復旧完了日時: `YYYY-MM-DDThh:mm:ssZ`
- 恒久復旧完了日時: `YYYY-MM-DDThh:mm:ssZ`

## 2. 影響（Impact）

- 影響対象: 例）`mem out search`
- 影響期間: `TBD`
- 影響規模: `TBD`
- 再試行回数: `0`
- ロールバック実施有無: `No`
- 再計画チケットID: `TBD`
- CIA影響: 可用性（低）

## 3. 原因分析（5 Whys）

1. Why1: 設定値が誤っていたため
2. Why2: レビュー時に差分検知できなかったため
3. Why3: 変更時の確認チェックが不足していたため
4. Why4: 運用手順が未整備だったため
5. Why5: インシデント起因変更の記録ルールがなかったため

## 4. 再発防止策（Preventive Actions）

- 恒久対策: `CHANGES.md` にタグ付き記録ルールを追加
- 検知強化: 影響範囲チェックリストをテンプレート化
- 運用改善: 事後レビューを 1 週間以内に実施
- オーナー: `TBD`
- 期限: `YYYY-MM-DD`

## 5. タイムライン（Timeline）

| 時刻（UTC） | 事実 | 判断/対応 | 担当 |
| --- | --- | --- | --- |
| YYYY-MM-DDThh:mm:ssZ | 事象を検知 | 一次切り分け開始 | TBD |
| YYYY-MM-DDThh:mm:ssZ | 原因候補を特定 | 暫定回避策を適用 | TBD |

## 6. 証跡（Evidence）

- 証跡ファイル:
  - `artifacts/ops/incident-summary.json`
  - `artifacts/ops/recovery-log.ndjson`
  - その他関連ログ/メトリクス:
