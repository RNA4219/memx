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

## 3. 統合レポート必須項目
統合レポート（`DESIGN-ACCEPTANCE-YYYYMMDD.md`）には以下 6 項目を必須で含める。

1. **対象章**
   - 受け入れ対象の章一覧（`path#section` 形式）
2. **REQ網羅率**
   - `coverage_rate`（%）
3. **high差分件数**
   - 契約同期結果の `severity: high` 件数
4. **リンク不達件数**
   - `link_unreachable_count`
5. **Birdseye issue件数**
   - `memx-birdseye-validation-spec.md` に基づく未解決 issue 件数
6. **最終判定**
   - `pass` または `fail`


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

## 4. 判定規則（固定）
- `memx_spec_v3/docs/design-review-spec.md` と合わせて、最終判定の正本は `memx_spec_v3/docs/design-doc-dod-spec.md` とする。

最終判定は以下の固定ルールで算出する。

- `high差分件数 > 0` の場合は `fail`
- `REQ網羅率 < 100%` の場合は `fail`
- `リンク不達件数 > 0` の場合は `fail`
- `Birdseye issue件数 > 0` の場合は `fail`
- 上記いずれにも該当しない場合のみ `pass`

## 5. 保存場所・命名規則（必須）
- 保存先は `memx_spec_v3/docs/reviews/` に固定する。
- ファイル名は `DESIGN-ACCEPTANCE-YYYYMMDD.md` とする。
- 例: `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-20260304.md`

## 6. 記録テンプレート（最小）
```md
# DESIGN ACCEPTANCE REPORT: <title>
- Report ID: DESIGN-ACCEPTANCE-YYYYMMDD
- 対象章:
  - memx_spec_v3/docs/design.md#...
  - memx_spec_v3/docs/interfaces.md#...

## メトリクス
- REQ網羅率: 100%
- high差分件数: 0
- リンク不達件数: 0
- Birdseye issue件数: 0

## 最終判定
- 判定: pass|fail
- 根拠:
  - requirements-coverage: <artifact/path>
  - contract-alignment: <artifact/path>
  - link-integrity: <artifact/path>
  - birdseye-validation: <artifact/path>
  - design-review: <artifact/path>
  - chapter-validation: <artifact/path>
```
