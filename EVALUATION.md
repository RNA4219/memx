---
intent_id: INT-001
owner: memx-resolver
status: active
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# Evaluation

## Acceptance Criteria

- [ ] workflow-cookbookの主要文書をdocとして登録できる
- [ ] feature名からrequired / recommended docsを返せる
- [ ] docから必要chunkを取得できる
- [ ] 読了したdoc versionとchunk_idsをtaskに紐づけて記録できる
- [ ] doc version更新時にstale判定できる
- [ ] docs/interfaces.md に記載した最小APIが呼べる
- [ ] 文書を登録するとchunkが生成される
- [ ] contract resolveがacceptance / forbidden / DoD / dependenciesを返せる

## KPIs

| 指標 | 目的 | 目標値 |
| --- | --- | --- |
| doc解決精度 | feature/taskから適切なdocを返す割合 | 90%以上 |
| chunk取得時間 | chunk取得のレスポンス時間 | 100ms以下 |
| stale検知率 | 更新されたdocのstale検知率 | 100% |

## Test Outline

- 単体: 各APIの入出力テスト
- 結合: doc登録→解決→chunk取得→読了記録→stale判定のフロー
- 回帰: 既存機能への影響確認

## Verification Checklist

- [ ] 主要フローが動作する（手動確認）
- [ ] エラー時挙動が明示されている
- [ ] 依存関係が再現できる環境である