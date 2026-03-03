# TASKS 運用ガイド（memx）

## 1. 命名規則
- Task Seed は **`TASK.<slug>-<MM-DD-YYYY>.md`** 形式で作成する。
- 日付は `date +%m-%d-%Y` の結果をそのまま使う（例: `03-03-2026`）。
- `slug` は英小文字・数字・ハイフンのみを使用する。

## 2. Task Seed 必須項目
各 Task Seed には次の5項目を必須で含める。

### Objective
- タスクの目的を1〜3行で記載する。

### Requirements
- 満たすべき要件を箇条書きで記載する。
- 互換性や非機能制約がある場合はここに明記する。

### Commands
- 実行・検証コマンドを列挙する。
- 実行順を持つ場合は上から順に並べる。

### Dependencies
- 前提タスク、依存PR、外部条件を列挙する。
- 依存がない場合は `- none` と記載する。

### Status
- ステータス語彙は HUB と同一の次のみを許可する。
  - planned
  - active
  - in_progress
  - reviewing
  - blocked
  - done

## 3. CHANGES 連携ルール（memx_spec_v3/CHANGES.md / CHANGELOG.md）
- 正本（canonical source）はリポジトリルートの `CHANGELOG.md` とする。
- `memx_spec_v3/CHANGES.md` は v3 仕様の履歴・互換性破壊テンプレート管理用の補助台帳として扱う。

完了タスクは次の手順で `memx_spec_v3/CHANGES.md` と `CHANGELOG.md` に反映する。

1. Task Seed の `Status` を `done` に更新する。
2. Task Seed の Objective/Requirements の要点を 1〜3 行で要約する。
3. `memx_spec_v3/CHANGES.md` の該当バージョン節に箇条書きで追記する。
4. 同内容を `CHANGELOG.md` に 1〜3 行の最小要約で追記する（重複エントリは禁止）。
5. 互換性破壊がある場合は `memx_spec_v3/CHANGES.md` の「互換性破壊時の記載テンプレート」を必ず併記する。
6. 移送後、Task Seed 側には `Moved-to-CHANGES: YYYY-MM-DD` を追記してトレース可能にする。
