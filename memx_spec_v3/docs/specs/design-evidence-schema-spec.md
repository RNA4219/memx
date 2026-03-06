# Design Evidence Schema Spec

## 1. 目的
本仕様は、検証成果物の証跡メタデータを共通化し、Phase 2〜4 の判定根拠を同一キーで機械検証できる状態に固定する。

## 2. 適用対象
以下の検証成果物は、本仕様の共通必須キーを必ず含む。

1. `memx_spec_v3/docs/requirements-coverage-spec.md`
2. `memx_spec_v3/docs/contract-alignment-spec.md`
3. `memx_spec_v3/docs/link-integrity-spec.md`
4. `docs/birdseye/memx-birdseye-validation-spec.md`

## 3. 共通必須キー
各成果物の集計出力・詳細出力の両方で、次のキーを必須とする。

| key | type | 必須 | 説明 |
| --- | --- | --- | --- |
| `run_id` | string | yes | 実行単位の一意ID（例: `RC-20260304-001`） |
| `generated_at` | string(datetime) | yes | 生成時刻（UTC, ISO8601） |
| `source_commit` | string | yes | 検証対象コミットSHA（短縮可） |
| `chapter_id` | string | yes | 対象章ID（章横断集計時は `global`） |
| `tool` | string | yes | 実行ツール識別子（例: `rg`, `python`, `custom-checker`） |
| `status` | enum | yes | `pass` / `fail` / `blocked` |
| `severity_summary` | object | yes | `high` / `medium` / `low` 件数の集計 |
| `evidence_paths` | string[] | yes | 判定根拠ファイル群（実在パスのみ） |

## 4. 追記ルール
- 既存の個別キーは維持し、先頭または末尾に共通必須キーを追加する。
- `chapter_id` は `design-chapter-validation-spec.md` の `chapter_id` と同値で管理する。
- `evidence_paths` は `TBD`・テンプレート値を禁止し、解決可能な相対パスのみ許可する。
- `severity_summary` は対象仕様の重大度語彙に揃える（未使用重大度は `0` を明示）。

## 5. 検証コマンド最小要件
- Markdown証跡は `rg -n "run_id|generated_at|source_commit|chapter_id|tool|status|severity_summary|evidence_paths" <artifact>` で存在確認する。
- JSON証跡を採用する場合は同一キーを JSON オブジェクト直下に持つことを検証する。

## 6. workflow-cookbook/schemas 配置方針
将来CIでの自動検証に備え、配置は **JSON Schema 併設** を採用する。

- 規約文書: 本仕様（Markdown）を正本とする。
- 機械検証: `workflow-cookbook/schemas/` に JSON Schema を追加し、CI から参照する。
- 運用方針: Markdown 変更時は同一PRで JSON Schema も同期更新する。
