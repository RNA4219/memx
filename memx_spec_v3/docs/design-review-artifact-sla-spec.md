# Design Review Artifact SLA Spec

## 1. 目的
本仕様は、`memx_spec_v3/docs/reviews/` 配下のレビュー実体ファイル（`DESIGN-REVIEW-*` / `DESIGN-CHAPTER-VALIDATION-*` / `DESIGN-ACCEPTANCE-*`）に対する作成SLAを定義する。

## 2. 対象
- 対象ディレクトリ: `memx_spec_v3/docs/reviews/`
- 対象成果物（実体ファイルのみ）:
  - `DESIGN-REVIEW-<実日付>-<連番>.md`
  - `DESIGN-CHAPTER-VALIDATION-<実日付>.md`
  - `DESIGN-ACCEPTANCE-<実日付>.md`
- テンプレート（例: `DESIGN-ACCEPTANCE-YYYYMMDD.md`）は対象外。

## 3. 成果物別SLA

### 3.1 DESIGN-REVIEW 実体
- 成果物: `DESIGN-REVIEW-<実日付>-<連番>.md`
- 作成トリガー:
  1. Phase 3 entry（レビュー開始時）
  2. Phase 3 exit（最終判定確定前）
  3. 差戻し時（再レビュー開始時）
- 作成責任者（owner）: レビュー担当者（design reviewer）
- 提出期限:
  - Phase 3 entry 分: entry 当日中（遷移前必須）
  - Phase 3 exit 分: exit 判定前必須
  - 差戻し再提出: 差戻し起票から **2営業日以内**
- 未提出時のステータス遷移制約:
  - `Status` は `reviewing` を維持する。
  - `Status: done` への遷移を禁止する。

### 3.2 DESIGN-CHAPTER-VALIDATION 実体
- 成果物: `DESIGN-CHAPTER-VALIDATION-<実日付>.md`
- 作成トリガー:
  1. Phase 2 exit（章別検証完了時）
  2. Phase 3 entry（レビュー投入前の章別再計測時）
  3. 差戻し時（対象章を再検証した時点）
- 作成責任者（owner）: 章別検証担当者（chapter validation owner）
- 提出期限:
  - Phase 遷移前必須（Phase 2 exit / Phase 3 entry）
  - 差戻し再提出: 差戻し起票から **2営業日以内**
- 未提出時のステータス遷移制約:
  - `Status` は `reviewing` を維持する。
  - `Status: done` への遷移を禁止する。

### 3.3 DESIGN-ACCEPTANCE 実体
- 成果物: `DESIGN-ACCEPTANCE-<実日付>.md`
- 作成トリガー:
  1. Phase 4 entry（受け入れレビュー開始時）
  2. Phase 4 exit（最終受け入れ判定前）
  3. 差戻し時（再判定時）
- 作成責任者（owner）: 受け入れ判定責任者（acceptance owner）
- 提出期限:
  - Phase 4 entry 分: entry 当日中（遷移前必須）
  - Phase 4 exit 分: exit 判定前必須
  - 差戻し再提出: 差戻し起票から **2営業日以内**
- 未提出時のステータス遷移制約:
  - `Status` は `reviewing` を維持する。
  - `Status: done` への遷移を禁止する。

## 4. 運用ルール
- 実体未作成・命名不正・テンプレート誤用のいずれかを検知した場合、差戻しとして扱う。
- 差戻し中は関連 Task Seed に不足成果物名・検知コマンド・再提出期限を必ず記録する。
- 本SLAは `design-acceptance-lifecycle-spec.md` と `design-review-remediation-spec.md` から参照される正本とする。

## 5. 参照
- `memx_spec_v3/docs/design-acceptance-lifecycle-spec.md`
- `memx_spec_v3/docs/design-review-remediation-spec.md`
- `docs/TASKS.md`
