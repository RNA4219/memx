---
intent_id: memx-governance-blueprint-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# BLUEPRINT

## 目的
- memx は個人用・ローカル運用の知識＋記憶管理基盤を提供する。
- v1 は後方互換を最優先し、CLI/API/SQLite の安定運用を重視する。

## 正本への導線
- 要求事項と ID 定義: [memx_spec_v3/docs/requirements.md](./memx_spec_v3/docs/requirements.md)
- 仕様（ユースケース/スコープ/非ゴール/受け入れ）: [memx_spec_v3/docs/spec.md](./memx_spec_v3/docs/spec.md)
- 設計（レイヤ/DB責務/移行）: [memx_spec_v3/docs/design.md](./memx_spec_v3/docs/design.md)
- I/F（CLI/API, 互換, エラー）: [memx_spec_v3/docs/interfaces.md](./memx_spec_v3/docs/interfaces.md)
- 機械可読契約（フィールド単位）: [memx_spec_v3/docs/CONTRACTS.md](./memx_spec_v3/docs/CONTRACTS.md)

## 互換性ポリシー
- v1 の必須 I/F は破壊変更を禁止する。
- 破壊変更は `FUTURE(v2+)` に隔離し、段階移行で導入する。
- 実験機能は feature flag 既定 OFF で追加する。
