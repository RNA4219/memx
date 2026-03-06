# Design Acceptance Lifecycle Spec

## 1. 目的
本仕様は、設計受け入れレポートのテンプレートと実体ファイルの責務、命名、差戻し条件を定義する。作成タイミング・提出期限・ステータス遷移制約のSLAは `memx_spec_v3/docs/design-review-artifact-sla-spec.md` を正本とする。

## 2. テンプレートと実体ファイルの責務分離
### 2.1 テンプレート（チェックID: DA-LC-01）
- `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用とする。
- テンプレートは章構成・必須キー・記載形式の定義のみを保持し、判定実績を記録しない。

### 2.2 実体ファイル（チェックID: DA-LC-02）
- 受け入れ判定の実績は `memx_spec_v3/docs/reviews/` 配下の実体ファイルへ記録する。
- 実体ファイルの新規作成を必須とし、既存実体ファイルの上書き・テンプレートの直接運用・改名運用を禁止する。

## 3. 実体命名規則（チェックID: DA-LC-03）
- 実体ファイル名は `DESIGN-ACCEPTANCE-<実日付>.md` とする。
- `<実日付>` はローカル日付 `YYYYMMDD` を使用する。
- 例: `DESIGN-ACCEPTANCE-20260304.md`

## 4. 作成タイミング（チェックID: DA-LC-04）
`memx_spec_v3/docs/design-review-artifact-sla-spec.md` の `3.3 DESIGN-ACCEPTANCE 実体` に従う。

## 5. 差戻し条件（チェックID: DA-LC-05）
以下のいずれかに該当した場合、受け入れ判定を差し戻す。
1. 必須6項目（対象章/REQ網羅率/high差分件数/リンク不達件数/Birdseye issue件数/最終判定）の欠落
2. `DESIGN-ACCEPTANCE-YYYYMMDD.md`（テンプレート専用）を実体判定・Status遷移・Release判定の証跡として利用した
3. `evidence_paths` の不整合（未存在パス、誤パス、参照不能パスを含む）
4. gate 列が `high`

## 6. 参照
- 入力元・判定規則: `memx_spec_v3/docs/design-acceptance-report-spec.md`
- 作成SLA（作成トリガー/owner/提出期限/未提出時制約）: `memx_spec_v3/docs/design-review-artifact-sla-spec.md`
- レビュー記録運用: `memx_spec_v3/docs/reviews/README.md`
- 最終判定正本: `memx_spec_v3/docs/design-doc-dod-spec.md`
