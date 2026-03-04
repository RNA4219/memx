# Design Source Inventory Spec

## 目的
- Phase 1（情報収集）で、情報源（固定入力+拡張入力）の抽出結果を同一フォーマットで管理し、`docs/TASKS.md` へ転記可能な粒度に正規化する。

## 対象入力（固定）
- `memx_spec_v3/docs/requirements.md`
- `memx_spec_v3/docs/traceability.md`
- `memx_spec_v3/docs/design.md`
- `memx_spec_v3/docs/interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`
- `docs/IN-*.md`（実績インシデントのみ。テンプレート除外）
- `orchestration/*.md`
- 必要に応じて `docs/INCIDENT_TEMPLATE.md`（参照のみ。実績証跡扱い不可）

## 関連仕様
- 運用ルールは `design-source-inventory-operations-spec.md` を参照する。

## インベントリ行フォーマット
抽出結果は 1 行 1 項目で記録し、以下の必須列をすべて持つこと。

| 列名 | 必須 | 説明 |
| --- | --- | --- |
| `source_type` | 必須 | 情報源種別。`Blueprint` / `Runbook` / `Incident` / `Orchestration` / `Reference` などの固定語彙を使う |
| `source_path#section` | 必須 | 情報源ファイルと章/セクション。例: `memx_spec_v3/docs/design.md#3.2 Data Flow` |
| `req_id` | 必須 | REQ-ID。複数ある場合は 1 行 1 ID に分割 |
| `contract_ref` | 必須 | 契約参照（OpenAPI/CLI schema/本文見出しなど） |
| `node_id` | 必須 | `docs/birdseye/index.json` の `node_id` |
| `depends_on` | 必須 | 先行 node_id または前提要件。無ければ `none` |
| `owner` | 必須 | 抽出結果の責任者（GitHub ID またはチーム名） |
| `reviewed_at` | 必須 | 最終レビュー日（`YYYY-MM-DD`） |
| `node_resolution_status` | 必須 | `source_path#section` から node を解決した結果（`resolved` / `ambiguous` / `missing`） |

## 欠損時の扱い
以下のいずれかに該当した行は `Status: blocked` とする。

- `source_path#section` が不正（対象入力外、またはセクション未特定）
- `source_type` が未記入、または `source_path#section` と整合しない（例: `docs/IN-*.md` を `Incident` 以外で記録）
- `req_id` が未記入、または `memx_spec_v3/docs/requirements.md` に存在しない
- `contract_ref` が未記入で契約整合の追跡不能
- `node_id` が `docs/birdseye/index.json` に存在しない
- `node_resolution_status` が `ambiguous` または `missing`
- `owner` / `reviewed_at` が未設定
- Incident 情報源（`docs/IN-*.md`）の取り込みが 0 件
- `orchestration/*.md` 由来の `depends_on` が未解決（Task Seed 化不可）

差し戻し条件（Phase 1 Done 不可）は以下。

- `blocked` 行が 1 件でも残っている
- `node_resolution_status` が `ambiguous` または `missing` の行が 1 件でもある（Phase 1 exit 不可）
- `depends_on` が未解決で Task Seed 化（<=0.5d）できない行がある
- Incident 未取り込み、または Orchestration 依存未解決の `blocked` 条件が 1 件でもある
- 必須列が欠けた行を含む状態で `docs/TASKS.md` へ転記しようとしている
