# Design Chapter Validation Summary Spec

## 1. 目的
本仕様は、Phase 2〜4 で共通利用する章別検証サマリの必須フィールド・参照先・判定利用ルールを固定し、章単位の整合確認結果を一貫した形式で記録するための共通フォーマットを定義する。

## 2. 適用範囲
- `orchestration/memx-design-docs-authoring.md` の Phase 2 / Phase 3 / Phase 4 で作成・更新する章別検証結果。
- `memx_spec_v3/docs/design-review-spec.md` に基づくレビュー記録。
- `memx_spec_v3/docs/design-acceptance-report-spec.md` に基づく統合受け入れレポート。

## 3. 章別検証サマリの必須フィールド
章別検証サマリは、1章につき 1 レコードで以下の 8 フィールドを必須で含める。

1. **`chapter_id`**
   - 章識別子。`path#section` 形式で記載する。
   - 例: `memx_spec_v3/docs/design.md#3. データフロー`
2. **`req_coverage`（%）**
   - 対象章に割り当てられた要件IDの網羅率（0〜100）。
3. **`contract_alignment_high_count`**
   - 契約整合チェックで検出した `severity: high` 件数。
4. **`link_broken_count`**
   - 対象章内および章間リンクの不達件数。
5. **`birdseye_issue_count`**
   - `docs/birdseye/memx-birdseye-validation-spec.md` 準拠で未解決の issue 件数。
6. **`evidence_paths`**
   - レビュー記録および受け入れレポートへの参照パス配列。
   - 最低 2 件（レビュー記録 1 件 + 受け入れレポート 1 件）を必須とする。
7. **`mapping_spec_ref`**
   - `chapter_id -> node_id` 対応規則の参照先。
   - 固定値: `memx_spec_v3/docs/design-chapter-node-mapping-spec.md`。
8. **`mapping_match_check`**
   - 章対応表との一致確認手順の実施結果。
   - `pass` / `fail` とし、`pass` 判定には以下を満たすこと。
     - `chapter_id` が章対応表に存在する。
     - 対応する `node_id` が `docs/birdseye/index.json` の最新値と一致する。

## 4. 記録ルール
- 章別検証サマリは、章単位で追記・更新可能な Markdown テーブルまたは YAML 配列で管理する。
- `chapter_id` は同一文書内で一意とする。
- `chapter_id` の命名と廃止運用は `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` に従う。
- `req_coverage` は `%` 付き表記（例: `100%`）または数値（例: `100`）のいずれかに統一する。
- `evidence_paths` は実在ファイルのみを許可し、テンプレートパスや `TBD` を禁止する。
- `mapping_match_check` の判定ログ（比較日時/比較対象）をレビュー記録または受け入れレポートへ残す。

## 4.1 `req_coverage` / `mapping_match_check` 算出手順（算出元の固定）
`req_coverage` と `mapping_match_check` は、以下の算出元を固定して記録する。

1. `req_coverage` の算出元
   - 要件トレース: `memx_spec_v3/docs/traceability.md`
   - 章定義: `memx_spec_v3/docs/design.md` / `memx_spec_v3/docs/interfaces.md`
   - 受け入れ・レビュー根拠: `memx_spec_v3/docs/reviews/` 配下の `DESIGN-REVIEW-<実日付>-<連番>.md` / `DESIGN-ACCEPTANCE-<実日付>.md`
2. `mapping_match_check` の算出元
   - 章対応表: `memx_spec_v3/docs/design-chapter-node-mapping-spec.md`
   - Birdseye ノード実体: `docs/birdseye/index.json`
   - 実施ログ: `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-<実日付>.md`
3. 実施順序（必須）
   - 先に `traceability.md` を使って章ごとの `REQ-ID` カバレッジを再計算し、`req_coverage` を更新する。
   - 次に章対応表と Birdseye index を突合して `mapping_match_check` を更新する。
   - 最後に契約整合レポート（`CONTRACT-ALIGN-<実日付>-<連番>.md` / `LATEST.md`）と章別検証レポートの整合を確認し、`evidence_paths` を確定する。

## 5. Phase 判定での利用ルール
- **Phase 2**: 章別ドラフト完成時に、初期値として章別検証サマリを作成する。
- **Phase 3**: 契約整合・リンク健全性・Birdseye 修正結果を章別検証サマリへ反映する。
- **Phase 4**: レビュー記録と統合受け入れレポートの参照を `evidence_paths` に追記し、最終判定根拠として使用する。

## 6. 最小テンプレート
```yaml
- chapter_id: memx_spec_v3/docs/design.md#3. データフロー
  req_coverage: 100
  contract_alignment_high_count: 0
  link_broken_count: 0
  birdseye_issue_count: 0
  evidence_paths:
    - memx_spec_v3/docs/reviews/DESIGN-REVIEW-20260304-001.md
    - memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-20260304.md
  mapping_spec_ref: memx_spec_v3/docs/design-chapter-node-mapping-spec.md
  mapping_match_check: pass
```
