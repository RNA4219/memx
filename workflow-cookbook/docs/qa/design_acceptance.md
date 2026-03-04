---
intent_id: DOC-LEGACY
owner: docs-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-04-04
---

# Design Acceptance Checklist

`docs/design.md` と `docs/interfaces.md` の設計記述に対して、受入判定を揃えるための公式チェックリストです。
判定は各観点で `accepted` または `rework` を選び、根拠リンク（設計節 / テスト・評価節）を必ず併記してください。

## 判定ルール

- `accepted`: 観点を満たし、設計節と検証節の双方に根拠がある。
- `rework`: 観点を満たさない、または根拠リンクが不足している。

## 受入観点

| 観点 | チェック内容 | 判定 (`accepted/rework`) | 設計節リンク | テスト/評価節リンク | 判定メモ |
| :-- | :-- | :-- | :-- | :-- | :-- |
| 責務境界 | 機能ごとの提供物/入力/依存が境界として明示され、重複責務がない。 | [ ] accepted / [ ] rework | [docs/interfaces.md#boundary-map](../interfaces.md#boundary-map) | [EVALUATION.md#verification-checklist](../../EVALUATION.md#verification-checklist) | |
| 互換性 | 既存ドキュメント・テンプレートとの整合があり、破壊的変更の意図が明示されている。 | [ ] accepted / [ ] rework | [docs/design.md#design](../design.md#design) | [EVALUATION.md#acceptance-criteria](../../EVALUATION.md#acceptance-criteria) | |
| 例外方針 | 失敗時の扱い（検知、是正、再実行）が設計上で定義され、運用チェックと接続されている。 | [ ] accepted / [ ] rework | [docs/design.md#design](../design.md#design) | [EVALUATION.md#verification-checklist](../../EVALUATION.md#verification-checklist) / [docs/security/Security_Review_Checklist.md](../security/Security_Review_Checklist.md) | |
| 非機能影響 | CI/レビュー/運用で必要な検証項目（品質・運用負荷・監視）への影響が追跡可能である。 | [ ] accepted / [ ] rework | [docs/design.md#design](../design.md#design) / [docs/interfaces.md#boundary-map](../interfaces.md#boundary-map) | [EVALUATION.md#kpis](../../EVALUATION.md#kpis) / [EVALUATION.md#test-outline](../../EVALUATION.md#test-outline) | |

## 記録

- 実施日:
- 実施者:
- 対象変更 (PR / Task):
- 結論:
