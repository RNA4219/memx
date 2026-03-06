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

## Incident → Task Seed 変換時の必須フィールド（固定）
- 本節の必須フィールドは固定とし、省略・別名化を禁止する。
- 入力 Incident は `docs/IN-<実日付>-<連番>.md` のみを許可し、`docs/IN-YYYYMMDD-001.md` や `docs/IN-BASELINE.md` を証跡入力として扱わない。

### 固定フィールド
- `Source`
  - `docs/IN-<実日付>-<連番>.md#<対象章>` を必須記録する。
- `Requirements`
  - Incident 対象章から抽出した再発防止要件を転記する。
  - 関連要件IDを必須記録する。
  - waiver 運用時は waiver対象要件IDを必須記録する。
- `Commands`
  - Incident 要件の検証コマンドを 1 件以上記録する。
- `Dependencies`
  - Incident 由来依存を記録する。依存がない場合は `- none` を明記する。

### 入力品質要件
- Incident 実体ファイルに以下が欠落する場合、Task Seed 変換を `blocked` とする。
  - ID項目（`インシデントID` / `Task Seed 参照ID`）
  - 時刻項目（`発生日` / `起票日` / `検知日時` / `暫定復旧完了日時` / `恒久復旧完了日時`）
  - 要件ID項目（`関連要件ID`、waiver時は `waiver対象要件ID`）
  - waiver項目（waiver時）
  - 証跡パス（`関連証跡パス`）
