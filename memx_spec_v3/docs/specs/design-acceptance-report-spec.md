# Design Acceptance Report Spec

## 1. 目的
本仕様は、設計受け入れレビューの統合レポートの入力元・必須項目・判定規則・保存規約を固定し、Phase 4 の完了判定を一意化する。

## 1-1. 正本仕様の参照（最終判定）
- 最終判定（`pass` / `fail`）の正本は `memx_spec_v3/docs/design-doc-dod-spec.md` とする。
- 本仕様の判定規則は `memx_spec_v3/docs/design-doc-dod-spec.md` と矛盾してはならない。

## 2. 入力元（必須）
統合レポートは、以下 6 仕様の結果を入力として集約する。

1. `memx_spec_v3/docs/requirements-coverage-spec.md`
2. `memx_spec_v3/docs/contract-alignment-spec.md`
3. `memx_spec_v3/docs/link-integrity-spec.md`
4. `docs/birdseye/memx-birdseye-validation-spec.md`
5. `memx_spec_v3/docs/design-review-spec.md`
6. `memx_spec_v3/docs/design-chapter-validation-spec.md`

## 3. 統合レポートのテンプレート定義（正本参照）
- 統合レポート（`DESIGN-ACCEPTANCE-YYYYMMDD.md`）の必須セクション・必須キー・許可値・命名規則・保存先は `memx_spec_v3/docs/design-evidence-template-spec.md` を正本とする。
- `DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用ファイルとして扱い、リリース判定の実体記録に流用してはならない。
- 本仕様は Phase 4 の判定規則と入力要件のみを定義する。

## 3.1 章別検証サマリ参照（必須）
- 統合レポートは、対象章ごとに `memx_spec_v3/docs/design-chapter-validation-spec.md` の章別検証サマリ参照を必須で含める。
- 各章参照には次のフィールドを最低限含める。
  - `chapter_id`
  - `req_coverage`
  - `contract_alignment_high_count`
  - `link_broken_count`
  - `birdseye_issue_count`
  - `evidence_paths`
- 章別検証サマリ参照が欠落している章が 1 件でもある場合、最終判定は `fail` とする。


## 3.2 入力エビデンスの共通メタスキーマ検証（必須）
統合レポート作成前に、2章で定義した全入力エビデンスが `memx_spec_v3/docs/design-evidence-schema-spec.md` の共通必須キーに準拠していることを検証する。

- 検証対象キー: `run_id` / `generated_at` / `source_commit` / `chapter_id` / `tool` / `status` / `severity_summary` / `evidence_paths`
- 1キーでも欠落がある入力は不受理とし、統合レポートの最終判定は `fail` とする。
- `evidence_paths` は列挙された全パスが実在ファイルであることを必須とし、1件でも未存在パスがある入力は不受理（最終判定 `fail`）とする。
- lint/type/test/link/contract/birdseye/coverage の保存先・命名規則・最小記録粒度は `memx_spec_v3/docs/design-gate-evidence-spec.md` を正本として参照し、本仕様で重複定義しない。
- 受理可否は Task Seed の `Commands` に記録した検証コマンド結果で追跡可能であること。

## 3.3 変更計画（導入ステップ）
1. 既存4仕様の出力節に共通メタキー追記ルールを反映する。
2. 各検証成果物テンプレートへ共通キーを追記し、サンプルを更新する。
3. Phase 4 統合時に、入力エビデンスのキー存在チェックを必須ゲート化する。

## 3.4 実体ファイル記載テンプレート（必須6項目固定）
`memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<実日付>.md` の実体は、次の6項目をこの順序で必須記載とする。

1. 対象章
2. REQ網羅率
3. high差分件数
4. リンク不達件数
5. Birdseye issue件数
6. 最終判定

上記6項目はテンプレート専用ファイル（`memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md`）の章見出しとして固定し、名称変更・省略・順序入替を禁止する。テンプレートは実体判定には使用せず、実体は `DESIGN-ACCEPTANCE-<実日付>.md` のみを作成・更新する。

## 4. 判定規則（固定）
- `memx_spec_v3/docs/design-review-spec.md` と合わせて、最終判定の正本は `memx_spec_v3/docs/design-doc-dod-spec.md` とする。

最終判定は以下の固定ルールで算出する。

- `high差分件数 > 0` の場合は `fail`
- `REQ網羅率 < 100%` の場合は `fail`
- `リンク不達件数 > 0` の場合は `fail`
- `Birdseye issue件数 > 0` の場合は `fail`
- 上記いずれにも該当しない場合のみ `pass`

## 5. 保存場所・命名規則・テンプレート
- 保存場所・命名規則・作成タイミング・差戻し条件は `memx_spec_v3/docs/design-acceptance-lifecycle-spec.md` を正本とする。
- テンプレート本文（章構成・必須キー）は `memx_spec_v3/docs/design-evidence-template-spec.md` を参照する。
- 本仕様では受け入れレポートのライフサイクル運用を重複定義せず、`memx_spec_v3/docs/design-acceptance-lifecycle-spec.md` のチェックID（DA-LC-01〜05）へ参照集約する。

## 5.1 存在判定の正本ルール（必須）
- 受け入れレポートの `exist` 判定は、正本命名規約 `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<YYYYMMDD>.md` に合致した**実体ファイル**の実在確認のみを有効とする。
- `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md`（テンプレート専用）は `exist` 判定の対象外とする。
