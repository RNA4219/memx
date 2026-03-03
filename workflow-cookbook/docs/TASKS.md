---
intent_id: DOC-LEGACY
owner: docs-core
status: active
last_reviewed_at: 2025-10-28
next_review_due: 2025-11-28
---

# Task Seeds 運用ガイド

## 1. 目的

- `HUB.codex.md` や `BLUEPRINT.md` などで抽出された課題を、着手可能な Task Seed として整理する。
- インシデント記録（`docs/IN-*.md`）や品質基準（`EVALUATION.md`）と突合し、再発防止と受け入れ条件を明文化する。
- `CHECKLISTS.md` の出荷基準を満たすための実施ログを残し、後続レビューで追跡できるようにする。

## 2. 記入テンプレート

Task Seed は `TASK.codex.md` に定義されたテンプレートに準拠し、以下の要点を押さえる。

```markdown
---
intent_id: INT-xxx
owner: contributor-handle
status: planned|active|...
last_reviewed_at: YYYY-MM-DD
next_review_due: YYYY-MM-DD
---

# Task Seed Title

## Objective
- 目的を一文で記す。

## Scope
- In: 影響範囲（ディレクトリや機能名）
- Out: 明示的に触れない領域

## Requirements
- Behavior / I/O / Constraints / Acceptance Criteria を `TASK.codex.md` の粒度で箇条書き。
- セキュリティ関連タスクでは `REQ-SEC-001` / `REQ-RET-001` / `REQ-SEC-AUD-001` / `REQ-SEC-AUD-002` / `REQ-SEC-GRD-001` の参照を明記する。
- Lint/Type/Test のゼロエラーを必須条件として明記する。

## Affected Paths
- グロブ表記で変更予定ファイル群を列挙。

## Local Commands
- 利用スタックに応じたゲートコマンドを `TASK.codex.md` の雛形から抜粋。

## 3. 登録済み Task Seeds

- [docs/tasks/task-autosave-project-locks.md](tasks/task-autosave-project-locks.md): AutoSave と Merge のロック整合性を段階導入で検証するタスク。

## Deliverables
- PR 要約・Intent 明記・必要なドキュメント差分を記載。
```

> **用語補足**: `Objective`・`Scope`・`Requirements` の定義は [`docs/addenda/A_Glossary.md`](addenda/A_Glossary.md) を参照し、
> `TASK.codex.md` の章立てと整合させる。

## 3. 検証ログ（TDD 前提）

1. **テスト設計を先行**: 着手前に必要なユニット/統合テストを列挙し、期待する失敗/成功条件を `Tests` セクションへ記す。
2. **実行コマンドの記録**: `Tests` もしくは `Commands` セクションに、実際に走らせたコマンドと結果（例: `pytest -q` → fail/pass）を時系列で追記する。
3. **インシデントとの連携**: 再発防止策が `docs/IN-*.md` に存在する場合、該当節を参照し、テストケースや検証ログにリンクを残す。
4. **チェックリスト照合**: ゲート通過後は `CHECKLISTS.md` の該当項目を確認し、未完了項目があれば Follow-up へ移す。

> 代表的な検証手順は [`docs/addenda/I_Test_Cases.md`](addenda/I_Test_Cases.md) を参照する。

## 4. フォローアップ手順

- **未解決事項**: 実装後も残るリスクや TODO は `Follow-ups` セクションに列挙し、必要なら新規 Task Seed を起票する。
- **情報同期**: 追記した Task Seed は `HUB.codex.md` の分類に基づき、関連ドキュメント（Blueprint /
  Guardrails / Incident）とのリンクを整備する。
- **レビュー結果の反映**: レビュアーからの追加要求は `Notes` に記録し、着手が別タスクになる場合は Task Seed ID を採番して紐付ける。
- **完了判定**: `CHECKLISTS.md` と `EVALUATION.md` の条件を満たし、検証ログがすべてグリーンであることを確認して `status: done` へ更新する。
- **Changelog 通番**: 変更履歴を編集する場合は [README.md の変更履歴の更新ルール](../README.md#changelog-update-rules) に従い、
  既存の最大通番に 1 を加えて 4 桁ゼロ埋めで記録する。
- **成果の転記**: 完了した Task Seed の成果差分は `[Unreleased](../CHANGELOG.md#unreleased)` に通番付きで記録し、
  当該 Task Seed からリンクを張って追跡できるようにする。

---

更新日: 2025-10-14
