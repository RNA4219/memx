---
intent_id: INT-001
owner: your-handle
status: active   # draft|active|deprecated
last_reviewed_at: 2025-10-14
next_review_due: 2025-11-14
---

# Task Seed Template

## メタデータ

```yaml
task_id: YYYYMMDD-xx
repo: https://github.com/owner/repo
base_branch: main
work_branch: feat/short-slug
priority: P1|P2|P3
langs: [auto]   # auto | python | typescript | go | rust | etc.
```

## Objective

{{一文で目的}}

## Scope

- In: {{対象(ディレクトリ/機能/CLI)を箇条書き}}
- Out: {{非対象(触らない領域)を箇条書き}}

## Requirements

- Behavior:
  - {{期待挙動1}}
  - {{期待挙動2}}
- I/O Contract:
  - Input: {{型/例}}
  - Output: {{型/例}}
- Constraints:
  - 既存API破壊なし / 不要な依存追加なし
  - Lint/Type/Test はゼロエラー
- Acceptance Criteria:
  - {{検収条件1}}
  - {{検収条件2}}
- Security Requirement IDs（該当時は必須記載）:
  - `REQ-SEC-001` / `REQ-RET-001`
  - `REQ-SEC-AUD-001` / `REQ-SEC-AUD-002`
  - `REQ-SEC-GRD-001`

## Affected Paths

- {{glob例: backend/src/**, frontend/src/hooks/**, tools/*.sh}}

## Local Commands（存在するものだけ実行）

```bash
## Python
ruff check . && black --check . && mypy --strict . && pytest -q

## TypeScript/Node
pnpm lint && pnpm typecheck && pnpm test
npm run lint && npm run typecheck && npm test

## Go
go vet ./... && go test ./...

## Rust
cargo fmt --check && cargo clippy -- -D warnings && cargo test

## Fallback
make ci || true
```

## Deliverables

- PR: タイトル/要約/影響/ロールバックに加え、本文へ `Intent: INT-xxx` と `## EVALUATION` アンカーを明記
  - 必要なら `Priority Score: <number>` を追記
- Artifacts: 変更パッチ、テスト、必要ならREADME/CHANGELOG差分

---

## Plan

### Steps

1) 現状把握（対象ファイル列挙、既存テストとI/O確認）
2) 小さな差分で仕様を満たす実装
3) sample::fail の再現手順/前提/境界値を洗い出し、必要な工程を増補
4) テスト追加/更新（先に/同時）
5) コマンド群でゲート通過
6) ドキュメント最小更新（必要なら）

## Patch

***Provide a unified diff. Include full paths. New files must be complete.***

## Tests

### Outline

- Unit:
  - {{case-1: 入力→出力の最小例}}
  - {{case-2: エッジ/エラー例}}
- Integration:
  - {{代表シナリオ1つ}}

## Commands

### Run gates

- （上の "Local Commands" から該当スタックを選んで実行）

## Notes

### Rationale

- {{設計判断を1～2行}}

### Risks

- {{既知の制約/互換性リスク}}

### Follow-ups

- {{後続タスクあれば}}
