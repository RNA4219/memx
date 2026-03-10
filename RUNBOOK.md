---
intent_id: INT-001
owner: memx-resolver
status: active
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# Runbook

## Environments

- Local: `mem` CLI直接実行、SQLiteファイルは `*.db`
- CI: テスト実行環境、一時DB使用
- Prod: API サーバーモード、`mem api serve --addr 127.0.0.1:7766`

## Execute

### 準備
```bash
cd memx_spec_v3/go
go build ./cmd/mem
```

### 文書登録
```bash
mem in docs --title "Memory Import Spec" --doc-type spec --body "$(cat docs/specs/memory-import.md)"
```

### 文書解決
```bash
mem out resolve --feature memory-import
mem out resolve --task-id task:feature:local:123
```

### Chunk取得
```bash
mem out chunks --doc-id doc:spec:memory-import
mem out chunks --doc-id doc:spec:memory-import --query "acceptance criteria"
```

### 読了記録
```bash
mem ack docs --task-id task:feature:local:123 --doc-id doc:spec:memory-import --version 2026-03-10
```

### Stale確認
```bash
mem stale-check --task-id task:feature:local:123
```

## Observability

- ログ: 標準出力/標準エラー出力
- メトリクス: 確認方法は開発中
- インシデント: `docs/IN-YYYYMMDD-XXX.md` に記録

## Confirm

- Execute結果を主要出力と突き合わせ
- CHECKLISTS.md で整合性確認

## Rollback / Retry

- DBファイルのバックアップから復元
- インシデント記録を `docs/IN-YYYYMMDD-XXX.md` に追記