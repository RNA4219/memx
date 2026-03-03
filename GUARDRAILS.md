---
intent_id: memx-safety-guardrails-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# GUARDRAILS

## セーフティ・判定
- Gatekeeper 判定は `allow` / `deny` / `needs_human`。
- v1.3 では `needs_human` を deny 相当として fail-closed 運用する。
- 保存前(`memory_store`)と出力前(`memory_output`)で必ずフック可能な構造を維持する。

## エラー・互換性
- API 最小保証コードは `INVALID_ARGUMENT` / `NOT_FOUND` / `INTERNAL`。
- 未分類エラーは互換維持のため `INTERNAL` にフォールバック。
- v1 内で禁止:
  - 必須フィールド削除
  - 既存フィールド意味変更
  - 既存成功レスポンスの型/構造破壊

## スキーマ/移行
- `migrate_other.go` は `schema/*.sql` 適用に統一し、部分適用失敗時はロールバック。
- DDL 適用順序は notes → notes_fts(採用時) → tags/note_tags → note_embeddings(採用時) → user_version。
- `user_version` は初期 1、破壊的/非互換 DDL のみ +1。

## データ一貫性
- ATTACH 跨ぎ完全原子性は前提にせず、「データ喪失より重複許容」で設計する。
- `lineage` により追跡・再蒸留可能性を確保する。

## セキュリティ
- APIキーや秘密情報は `memory_policy.yaml` と Gatekeeper で保存前ブロックを行う。
- v1 はローカル運用前提とし、認証・権限・監査を伴う公開運用は対象外。
