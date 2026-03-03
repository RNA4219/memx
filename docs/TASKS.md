# TASKS 運用ガイド（memx）

## 1. 命名規則
- Task Seed は **`TASK.<slug>-<MM-DD-YYYY>.md`** 形式で作成する。
- 日付は `date +%m-%d-%Y` の結果をそのまま使う（例: `03-03-2026`）。
- `slug` は英小文字・数字・ハイフンのみを使用する。

## 推奨 front matter キー
- Task Seed 冒頭に YAML front matter を付与し、以下キーの記載を推奨する。
  - `priority`: `high` / `medium` / `low` のいずれか。
  - `owner`: 担当チームまたは担当者（例: `memx-core`）。
  - `deadline`: 期日を `YYYY-MM-DD` 形式で記載する。

## 2. Task Seed 必須項目
各 Task Seed には次の6項目を必須で含める。

### Source
- 要件の出典を `path#Section` 形式で記載する。
- 例: `orchestration/memx-v1-bootstrap.md#Phase 2`
- 複数ある場合は箇条書きで列挙する。
- `HUB.codex.md` 工程2の運用ルールに従い、要件根拠として許可するインシデント記録は **`docs/IN-<実日付>-<連番>.md` のみ** とする（`docs/IN-BASELINE.md` は補助資料、`docs/IN-YYYYMMDD-001.md` などのテンプレートIDは不可）。
- `Source` にテンプレートID（例: `IN-YYYYMMDD-001`）または `TBD` を含む Task Seed は、`reviewing` を継続して差し戻す。

### Node IDs
- `docs/birdseye/index.json` の `node_id` を記載する。
- 依存グラフ（`depends_on`）と照合する対象タスクでは必須、照合対象外タスクでは任意とする。
- 記載形式は箇条書きで、依存元/依存先の対応がわかるように補足行を付ける。

### Objective
- タスクの目的を1〜3行で記載する。

### Requirements
- 満たすべき要件を箇条書きで記載する。
- 互換性や非機能制約がある場合はここに明記する。
- エラーコードを新規追加・変更するタスクでは、`memx_spec_v3/docs/requirements.md` の「6-4. エラーモデル」と `memx_spec_v3/docs/error-contract.md` を同時更新対象に含める。

### Commands
- 実行・検証コマンドを列挙する。
- 実行順を持つ場合は上から順に並べる。

### Dependencies
- 前提タスク、依存PR、外部条件を列挙する。
- 依存がない場合は `- none` と記載する。

### Release Note Draft
- `CHANGELOG.md` に転記する利用者影響の要約を1〜3行で記載する。
- 実装詳細ではなく「何が変わるか」「利用者への影響」を簡潔に記載する。

### Status
- ステータス語彙は HUB と同一の次のみを許可する。
  - planned
  - active
  - in_progress
  - reviewing
  - blocked
  - done
- `Status: done` へ遷移する条件として、Task Seed に `Release Note Draft` 記入済みであること。
- `Status: done` へ遷移する条件として、移送後に `Moved-to-CHANGES: YYYY-MM-DD` を追記済みであること。

## 2-1. 破壊変更時の追記チェックリスト（追加必須）
CLI/API の既存必須フィールド削除、型変更、意味変更、既存コマンド/エンドポイント/エラーコード削除、`--json` 既定出力の非同型化を含む場合、Task Seed に次のチェックリストを追記する。

- [ ] `Source` は `path#Section` で記載済み
- [ ] `Node IDs` を記載済み（依存照合対象なら必須）
- [ ] `Requirements` に後方互換/非機能制約を明記済み
- [ ] エラーコード変更時は `memx_spec_v3/docs/requirements.md` と `memx_spec_v3/docs/error-contract.md` を更新対象に含めた
- [ ] `Commands` に検証コマンドを順序付きで記載済み
- [ ] `Release Note Draft` を記載済み
- [ ] `memx_spec_v3/CHANGES.md` と `CHANGELOG.md` への反映項目を記載済み
- [ ] `Status: done` 前に `Moved-to-CHANGES: YYYY-MM-DD` を追記する

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
