---
intent_id: INT-001
owner: memx-resolver
status: draft
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# 実装計画（Implementation Plan）

本計画は [README.md](../README.md) および [HUB.codex.md](../HUB.codex.md) の導入指針に従い、cookbook-resolver 機能の段階導入を支える最小単位の意思決定と依存関係を整理する。

## フラグ方針

- `resolver.enabled` フラグで resolver 機能を段階的に有効化する
- 未完了状態では強制的にオフ

## 依存関係

- resolver 機能は `docs/requirements.md` の要件に従う
- API は `docs/interfaces.md` の定義に従う
- 実装は `docs/design.md` の方針に従う
- テストは `EVALUATION.md` の受け入れ基準に従う

## 段階導入チェックリスト

### Phase 1: データモデルとDB

1. [ ] `resolver_documents` テーブル作成
2. [ ] `resolver_chunks` テーブル作成
3. [ ] `resolver_document_links` テーブル作成
4. [ ] `resolver_read_receipts` テーブル作成

### Phase 2: 文書登録とChunk生成

1. [ ] `POST /v1/docs:ingest` API実装
2. [ ] Chunk生成ロジック実装（見出し優先）
3. [ ] 文書メタデータ保存

### Phase 3: 文書解決

1. [ ] `POST /v1/docs:resolve` API実装
2. [ ] feature/task/topic からの解決ロジック実装
3. [ ] required/recommended 分類ロジック実装

### Phase 4: Chunk取得

1. [ ] `POST /v1/chunks:get` API実装
2. [ ] doc_id/query/heading 指定取得実装
3. [ ] `POST /v1/docs:search` API実装

### Phase 5: 読了記録とStale判定

1. [ ] `POST /v1/reads:ack` API実装
2. [ ] `POST /v1/docs:stale-check` API実装
3. [ ] Stale判定ロジック実装（version比較）

### Phase 6: 契約解決

1. [ ] `POST /v1/contracts:resolve` API実装
2. [ ] acceptance_criteria/forbidden_patterns/DoD抽出

### Phase 7: Skill実装

1. [ ] `/resolve-docs` Skill実装
2. [ ] `/read-chunks` Skill実装
3. [ ] `/ack-docs` Skill実装
4. [ ] `/stale-check` Skill実装
5. [ ] `/resolve-contract` Skill実装

## 優先順位

| Phase | 優先度 | 依存 |
| --- | --- | --- |
| Phase 1 | 高 | なし |
| Phase 2 | 高 | Phase 1 |
| Phase 3 | 高 | Phase 2 |
| Phase 4 | 高 | Phase 2 |
| Phase 5 | 中 | Phase 3, Phase 4 |
| Phase 6 | 中 | Phase 3 |
| Phase 7 | 低 | Phase 1-6 |

---

- 逆リンク: [README.md](../README.md) / [HUB.codex.md](../HUB.codex.md)