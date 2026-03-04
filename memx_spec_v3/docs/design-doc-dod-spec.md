# Design Doc Definition of Done (DoD) Spec

## 1. 目的
本仕様は、設計ドキュメント群の完成判定（Definition of Done）を単一の正本として固定し、Phase 3/4 の重複判定定義を統一する。

## 2. 対象スコープ（必須）
本仕様の判定対象は次の成果物に固定する。

1. `memx_spec_v3/docs/design.md`
2. `memx_spec_v3/docs/interfaces.md`
3. `memx_spec_v3/docs/traceability.md`
4. `memx_spec_v3/docs/design-chapter-validation-spec.md` に準拠した章別検証サマリ
5. `memx_spec_v3/docs/design-acceptance-report-spec.md` に準拠した統合受け入れレポート

## 3. 完成判定ルール（固定）
最終判定を `pass` とする条件は次の全件充足のみとし、1件でも未達なら `fail` とする。

1. REQ網羅率が `100%`（`memx_spec_v3/docs/requirements-coverage-spec.md`）
2. contract alignment の `severity: high` 件数が `0`（`memx_spec_v3/docs/contract-alignment-spec.md`）
3. リンク不達件数が `0`（`memx_spec_v3/docs/link-integrity-spec.md`）
4. Birdseye issue 未解決件数が `0`（`docs/birdseye/memx-birdseye-validation-spec.md`）
5. レビュー記録が必須項目を満たして完備（`memx_spec_v3/docs/design-review-spec.md`）
6. 参照解決適合率が `100%`（`memx_spec_v3/docs/design-reference-conformance-spec.md`）

## 4. 判定実行順（推奨固定順）
1. 章別検証サマリを更新する。
2. 統合受け入れレポートを更新する。
3. 本仕様 3 章の5条件を照合する。
4. 判定結果を `pass` / `fail` で確定する。

## 5. 関連仕様への適用方針
- `orchestration/memx-design-docs-authoring.md` の Phase 3/4 Done Criteria で重複する個別判定文は、本仕様参照へ置換する。
- `memx_spec_v3/docs/design-acceptance-report-spec.md` と `memx_spec_v3/docs/design-review-spec.md` は、最終判定の正本が本仕様である旨を相互参照として明記する。
- `docs/TASKS.md` の Phase 2〜4 向け Task Seed 要件例には、本仕様参照を必須ルールとして追加する。
