---
intent_id: INT-001
owner: memx-resolver
status: active
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# Task Seed 運用ガイド

## 概要

Task Seed は `TASK.*-MM-DD-YYYY` 形式のファイルで、単一の作業単位を記述する。

## テンプレート

```markdown
---
task_id: TASK-001
status: planned
owner:
priority: high
created_at: 2026-03-10
deadline:
---

# Task Title

## Objective

- 何をするか

## Requirements

- [ ] 要件1
- [ ] 要件2

## Commands

```bash
# 実行コマンド
```

## Dependencies

- 依存タスクID

## Notes

- 補足情報
```

## ステータス遷移

```
planned → active → in_progress → reviewing → done
                ↓
              blocked → in_progress
```

## 現在のタスク一覧

| ID | タイトル | ステータス | 優先度 |
| --- | --- | --- | --- |
| TASK-001 | データモデル実装 | planned | high |
| TASK-002 | docs:ingest API | planned | high |
| TASK-003 | docs:resolve API | planned | high |
| TASK-004 | chunks:get API | planned | high |
| TASK-005 | reads:ack API | planned | medium |
| TASK-006 | stale-check API | planned | medium |
| TASK-007 | contracts:resolve API | planned | medium |

---

- 逆リンク: [HUB.codex.md](../HUB.codex.md)