# Incident to Task Traceability Spec

## 目的
- `docs/IN-<実日付>-<連番>.md` 由来の再発防止要件を、Task Seed の `Requirements` / `Commands` / `Dependencies` へ欠損なく転記し、検証可能な状態で管理する。

## 対象入力
- 入力元は `docs/IN-<実日付>-<連番>.md` の対象章とする。
- 対象章の例: `再発防止テスト観点`。
- `docs/IN-<YYYYMMDD>-001.md` などテンプレートID文書は対象外とする。

## 転記ルール（Task Seed）

### Requirements への転記
- Incident 対象章から抽出した再発防止要件を、`Requirements` に箇条書きで転記する。
- 各行に `docs/IN-<実日付>-<連番>.md#<対象章>` を併記し、出典追跡可能にする。

### Commands への転記
- Incident 要件の充足確認に使う検証コマンドを `Commands` へ転記する。
- 実行順がある場合は上から順に記載する。
- **最低 1 件以上の検証コマンド**を必須とする。

### Dependencies への転記
- Incident 要件の前提条件（先行Task、外部入力、依存PR）がある場合は `Dependencies` に転記する。
- 依存がない場合は `- none` を明記する。

## 完了判定
- Incident 対象章の転記先が `Requirements` / `Commands` / `Dependencies` の3節で確認できる。
- `Commands` に検証コマンドが 1 件以上記載されている。
- Task Seed の `Source` に Incident の実ID（`docs/IN-<実日付>-<連番>.md#<対象章>`）が記載されている。
