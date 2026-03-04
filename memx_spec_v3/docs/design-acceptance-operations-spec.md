# Design Acceptance Operations Spec

## 1. 目的
- 本仕様は、`memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<YYYYMMDD>.md` の作成・判定・差戻し運用を単一手順として固定する。
- 判定の入力元、判定順序、責務分担を統一し、Phase 4 の受け入れ判定ぶれを防止する。

## 2. 適用範囲
- 対象は `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<YYYYMMDD>.md` の**実体ファイルのみ**とする。
- `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用とし、章構造・必須キー定義のみを保持する。
- 実体ファイルは判定値（実測値・判定結果・根拠パス）を記録する責務を持つ。

## 3. 固定入力元
判定に使用する入力元は次に固定し、追加・代替する場合は本仕様改定を必須とする。
1. `memx_spec_v3/docs/requirements-coverage-spec.md`
2. `memx_spec_v3/docs/contract-alignment-spec.md`
3. `memx_spec_v3/docs/link-integrity-spec.md`
4. `docs/birdseye/memx-birdseye-validation-spec.md`
5. 章別検証サマリ（`memx_spec_v3/docs/design-chapter-validation-spec.md` 準拠）
6. レビュー記録（`memx_spec_v3/docs/design-review-spec.md` 準拠）

## 4. 判定手順（固定順序）

### 4.1 収集
- 第3章の固定入力元から、判定に必要な最新値を収集する。
- 収集結果は `evidence_paths` として実体ファイルへ列挙する。

### 4.2 整合チェック
- 収集した値と `evidence_paths` の参照先が一致することを確認する。
- 欠落・参照不能・版不一致がある場合は `fail` とする。

### 4.3 閾値判定
- `requirements-coverage-spec.md` / `contract-alignment-spec.md` / `link-integrity-spec.md` / `docs/birdseye/memx-birdseye-validation-spec.md` の判定規則に従い、受け入れ判定を実施する。
- いずれかが fail 条件に該当した場合、最終判定は `fail` とする。

### 4.4 fail 時差戻し
- `fail` 判定時は `docs/TASKS.md` の対象タスクを `reviewing` で維持し、`done` へ遷移させない。
- 差戻し理由（不足証跡、閾値未達、整合不一致）を実体ファイルに追記し、再作業を要求する。

### 4.5 再作成
- 差戻し後の再判定では、`memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<YYYYMMDD>.md` を新規作成する。
- 既存実体ファイルの上書き更新は禁止する。

## 5. 責務分担（owner / reviewer）

### 5.1 owner（レポート作成責任）
- 固定入力元の収集、整合チェック、閾値判定実行、実体ファイル作成を担当する。
- `docs/TASKS.md` の `Status: reviewing` 期間中に受け入れレポートを作成・更新する。

### 5.2 reviewer（承認責任）
- owner が作成した実体ファイルの証跡整合、閾値判定、差戻し要否を承認する。
- 承認条件を満たした場合のみ `docs/TASKS.md` の `Status` を `reviewing` から `done` へ遷移させる。

### 5.3 兼務制約
- `owner` と `reviewer` の兼務は禁止する。

## 6. 関連仕様
- 判定値の記録フォーマット: `memx_spec_v3/docs/design-acceptance-report-spec.md`
- テンプレート/実体ライフサイクル: `memx_spec_v3/docs/design-acceptance-lifecycle-spec.md`
- Task Status 遷移運用: `docs/TASKS.md`
