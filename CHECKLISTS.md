---
intent_id: INT-001
owner: memx-resolver
status: active
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# Checklists

## Hygiene

- [ ] コードフォーマット確認（gofmt / goimports）
- [ ] Lint確認（golangci-lint）
- [ ] 型チェック確認（go vet）
- [ ] テスト実行（go test ./...）
- [ ] ドキュメント更新（必要に応じて）

## Review

- [ ] PR本文にIntent: INT-xxxを含める
- [ ] EVALUATION見出しへのリンクを本文に明示
- [ ] 変更の影響範囲を確認
- [ ] 破壊的変更の有無を確認

## Release

- [ ] 全テストが通過
- [ ] CHANGELOG.mdを更新
- [ ] バージョンタグを作成
- [ ] リリースノートを作成

## Incident

- [ ] docs/IN-YYYYMMDD-XXX.mdを作成
- [ ] 影響範囲を記録
- [ ] 再発防止策を記録
- [ ] 関連PR/チケットへリンク