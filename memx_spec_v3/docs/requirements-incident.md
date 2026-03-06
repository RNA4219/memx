---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - インシデント対応

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 11. インシデント対応要件（運用）

- セキュリティ/品質インシデントは `docs/INCIDENT_TEMPLATE.md` に従って記録する。
- 実運用インシデントの受入対象は `docs/IN-<実日付>-<連番>.md` 形式のみとする。
- `docs/IN-YYYYMMDD-001.md` / `docs/IN-202603xx-001.md` はテンプレートであり、要件根拠・実績証跡としては扱わない。
- 初動時点で「検知」「影響」「5 Whys」「再発防止」「タイムライン」を最低限記載する。
- テンプレート: [`docs/IN-YYYYMMDD-001.md`](../../docs/IN-YYYYMMDD-001.md)
- 実在インシデント参照例: [`docs/IN-20260303-002.md`](../../docs/IN-20260303-002.md)
- 注記: `Source` にはテンプレートID（`IN-YYYYMMDD-001` / `IN-202603xx-001`）を記載せず、`docs/IN-<実日付>-<連番>.md` のみを記載する。