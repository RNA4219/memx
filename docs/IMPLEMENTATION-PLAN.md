---
intent_id: INT-001
owner: memx-resolver
status: active
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

### Phase 1: データモデルとDB ✅

1. [x] `resolver_documents` テーブル作成
2. [x] `resolver_chunks` テーブル作成
3. [x] `resolver_document_links` テーブル作成
4. [x] `resolver_read_receipts` テーブル作成

### Phase 2: 文書登録とChunk生成 ✅

1. [x] `POST /v1/docs:ingest` API実装
2. [x] Chunk生成ロジック実装（見出し優先）
3. [x] 文書メタデータ保存

### Phase 3: 文書解決 ✅

1. [x] `POST /v1/docs:resolve` API実装
2. [x] feature/task/topic からの解決ロジック実装
3. [x] required/recommended 分類ロジック実装

### Phase 4: Chunk取得 ✅

1. [x] `POST /v1/chunks:get` API実装
2. [x] doc_id/query/heading 指定取得実装
3. [x] `POST /v1/docs:search` API実装

### Phase 5: 読了記録とStale判定 ✅

1. [x] `POST /v1/reads:ack` API実装
2. [x] `POST /v1/docs:stale-check` API実装
3. [x] Stale判定ロジック実装（version比較）

### Phase 6: 契約解決 ✅

1. [x] `POST /v1/contracts:resolve` API実装
2. [x] acceptance_criteria/forbidden_patterns/DoD抽出

### Phase 7: CLI実装 ✅

1. [x] `mem docs ingest` CLI実装
2. [x] `mem docs resolve` CLI実装
3. [x] `mem docs chunks` CLI実装
4. [x] `mem docs search` CLI実装
5. [x] `mem docs ack` CLI実装
6. [x] `mem docs stale` CLI実装
7. [x] `mem docs contract` CLI実装

## 優先順位

| Phase | 優先度 | 依存 | 状態 |
| --- | --- | --- | --- |
| Phase 1 | 高 | なし | ✅ 完了 |
| Phase 2 | 高 | Phase 1 | ✅ 完了 |
| Phase 3 | 高 | Phase 2 | ✅ 完了 |
| Phase 4 | 高 | Phase 2 | ✅ 完了 |
| Phase 5 | 中 | Phase 3, Phase 4 | ✅ 完了 |
| Phase 6 | 中 | Phase 3 | ✅ 完了 |
| Phase 7 | 低 | Phase 1-6 | ✅ 完了 |

---

- 逆リンク: [README.md](../README.md) / [HUB.codex.md](../HUB.codex.md)