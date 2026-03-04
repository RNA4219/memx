# Design Source Inventory Spec

## 目的
- Phase 1（情報収集）で、情報源7ファイルの抽出結果を同一フォーマットで管理し、`docs/TASKS.md` へ転記可能な粒度に正規化する。

## 対象入力（固定）
- `memx_spec_v3/docs/requirements.md`
- `memx_spec_v3/docs/traceability.md`
- `memx_spec_v3/docs/design.md`
- `memx_spec_v3/docs/interfaces.md`
- `docs/birdseye/caps/EVALUATION.md.json`
- `docs/birdseye/caps/RUNBOOK.md.json`
- `docs/birdseye/index.json`

## インベントリ行フォーマット
抽出結果は 1 行 1 項目で記録し、以下の必須列をすべて持つこと。

| 列名 | 必須 | 説明 |
| --- | --- | --- |
| `source_path#section` | 必須 | 情報源ファイルと章/セクション。例: `memx_spec_v3/docs/design.md#3.2 Data Flow` |
| `req_id` | 必須 | REQ-ID。複数ある場合は 1 行 1 ID に分割 |
| `contract_ref` | 必須 | 契約参照（OpenAPI/CLI schema/本文見出しなど） |
| `node_id` | 必須 | `docs/birdseye/index.json` の `node_id` |
| `depends_on` | 必須 | 先行 node_id または前提要件。無ければ `none` |
| `owner` | 必須 | 抽出結果の責任者（GitHub ID またはチーム名） |
| `reviewed_at` | 必須 | 最終レビュー日（`YYYY-MM-DD`） |

## 欠損時の扱い
以下のいずれかに該当した行は `Status: blocked` とする。

- `source_path#section` が不正（対象入力外、またはセクション未特定）
- `req_id` が未記入、または `memx_spec_v3/docs/requirements.md` に存在しない
- `contract_ref` が未記入で契約整合の追跡不能
- `node_id` が `docs/birdseye/index.json` に存在しない
- `owner` / `reviewed_at` が未設定

差し戻し条件（Phase 1 Done 不可）は以下。

- `blocked` 行が 1 件でも残っている
- `depends_on` が未解決で Task Seed 化（<=0.5d）できない行がある
- 必須列が欠けた行を含む状態で `docs/TASKS.md` へ転記しようとしている
