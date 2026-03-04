---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
---

# design.md 章テンプレート（必須）

`memx_spec_v3/docs/design.md` の各章は、以下テンプレートを必須で満たすこと。
`docs/TASKS.md` の必須項目・語彙を直接反映し、Task Seed へ転記可能な粒度で記述する。

## 章テンプレート

### Objective
- 目的を 1〜3 行で記述する。

### Source
- `path#Section` 形式で記述する。
- 複数ある場合は箇条書きで列挙する。
- `TBD` やテンプレートID（例: `IN-<YYYYMMDD>-001`）は記載しない。

### Node IDs
- `docs/birdseye/index.json` の `node_id` を箇条書きで記述する。
- 依存照合対象の章では必須。必要に応じて `depends_on` 対応を補足する。

### Requirements
- REQ-ID を箇条書きで記述する（例: `REQ-API-001`）。
- 要件IDは `memx_spec_v3/docs/requirements.md` の該当節と一致させる。
- 契約変更・エラー変更を伴う場合は `docs/TASKS.md` の追加更新条件（契約/エラーモデル同時更新）を Requirements 内へ明示する。

### Commands
- 仕様整合チェック用コマンドを実行順に列挙する。
- 最低限、要件ID網羅・契約同期・リンク健全性の確認手順を含める。

### Dependencies
- 前提タスク、依存PR、外部条件を箇条書きで記述する。
- 依存がない場合は `- none` と記載する。

### Status
- 以下語彙のみ使用可（`docs/TASKS.md` 準拠）。
  - `planned`
  - `active`
  - `in_progress`
  - `reviewing`
  - `blocked`
  - `done`


## `docs/TASKS.md` 必須項目・語彙の反映

- `Source`: `path#Section` 形式を必須化（本テンプレートの `Source` 見出しで充足）。
- `Node IDs`: `docs/birdseye/index.json` の `node_id` を記載（本テンプレートの `Node IDs` 見出しで充足）。
- `Objective`: 1〜3行（本テンプレートの `Objective` 見出しで充足）。
- `Requirements`: REQ-ID 箇条書き（本テンプレートの `Requirements` 見出しで充足）。
- `Commands`: 検証コマンド列挙（本テンプレートの `Commands` 見出しで充足）。
- `Dependencies`: 前提条件列挙（本テンプレートの `Dependencies` 見出しで充足）。
- `Status`: 許可語彙 `planned/active/in_progress/reviewing/blocked/done` のみ使用。
- `Release Note Draft`: `design.md` 章本文では任意。Task Seed 化・完了化（`done`）時に必須で補完する。

## 記入例（雛形）

```md
### Objective
- <1〜3行で章の目的>

### Source
- <path#Section>

### Node IDs
- <node_id>

### Requirements
- <REQ-ID>

### Commands
- <仕様整合チェックコマンド>

### Dependencies
- <前提タスク/依存PR/外部条件>

### Status
- planned
```
