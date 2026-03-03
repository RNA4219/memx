# IN-YYYYMMDD-001（雛形インシデント）

- インシデントID: `IN-YYYYMMDD-001`
- 発生日（UTC）: `YYYY-MM-DD`
- 起票日（UTC）: `YYYY-MM-DD`
- 重大度: `SEV3`
- ステータス: `Open`
- 関連要件: [`memx 要件定義: 0.目的とスコープ`](../memx_spec_v3/docs/requirements.md#0-目的とスコープ)
- 使用テンプレート: [`docs/INCIDENT_TEMPLATE.md`](./INCIDENT_TEMPLATE.md)

## 1. 検知（Detection）

- 検知日時: `YYYY-MM-DDThh:mm:ssZ`
- 検知経路: 例）ユーザー報告
- 初動担当: `TBD`
- 事象概要: 例）誤ったメタデータ設定により検索結果の一部が欠落

## 2. 影響（Impact）

- 影響対象: 例）`mem out search`
- 影響期間: `TBD`
- 影響規模: `TBD`
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
