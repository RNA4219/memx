# Design Evidence Template Spec

## 1. 目的
本仕様は、設計エビデンス成果物テンプレートの正本定義（必須セクション・必須キー・許可値・命名規則・保存先）を一本化する。

## 2. 適用対象（3成果物限定）
- `DESIGN-SOURCE-INVENTORY-YYYYMMDD.md`
- `DESIGN-REVIEW-YYYYMMDD-###.md`
- `DESIGN-ACCEPTANCE-YYYYMMDD.md`

## 3. 共通必須キー（`design-evidence-schema-spec.md` 整合）
以下 6 キーは 3 成果物すべてで必須とする。

| key | type | 必須 | 許可値/制約 |
| --- | --- | --- | --- |
| `run_id` | string | yes | `^[A-Z]+-[0-9]{8}-[0-9]{3}$` 例: `RC-20260304-001` |
| `generated_at` | string(datetime) | yes | UTC ISO8601（例: `2026-03-04T12:34:56Z`） |
| `source_commit` | string | yes | Git SHA（短縮可） |
| `chapter_id` | string | yes | `design-chapter-validation-spec.md` と同値、章横断は `global` |
| `status` | enum | yes | `pass` / `fail` / `blocked` |
| `evidence_paths` | string[] | yes | 解決可能な相対パスのみ（`TBD` 禁止） |

## 4. テンプレート定義（正本）

### 4.1 `DESIGN-SOURCE-INVENTORY-YYYYMMDD.md`

| 区分 | 定義 |
| --- | --- |
| 保存先 | `memx_spec_v3/docs/reviews/inventory/` |
| 命名規則 | `DESIGN-SOURCE-INVENTORY-YYYYMMDD.md`（`YYYYMMDD` は JST 作業日） |
| 必須セクション | `# DESIGN SOURCE INVENTORY` / `## Metadata` / `## Inventory Table` / `## Approval` |
| 必須キー | `run_id`, `generated_at`, `source_commit`, `chapter_id`, `status`, `evidence_paths`, `req_id`, `source_path`, `reviewed_at` |
| 許可値 | `status`: `pass` / `fail` / `blocked`; `reviewed_at`: `YYYY-MM-DD`; `source_path`: `path#section` |

### 4.2 `DESIGN-REVIEW-YYYYMMDD-###.md`

| 区分 | 定義 |
| --- | --- |
| 保存先 | `memx_spec_v3/docs/reviews/` |
| 命名規則 | `DESIGN-REVIEW-YYYYMMDD-###.md`（`###` は同日 001 始まり連番） |
| 必須セクション | `# DESIGN REVIEW` / `## Metadata` / `## Findings` / `## Re-check` / `## Decision` |
| 必須キー | `run_id`, `generated_at`, `source_commit`, `chapter_id`, `status`, `evidence_paths`, `review_id`, `req_ids`, `node_ids`, `decision` |
| 許可値 | `status`: `pass` / `fail` / `blocked`; `decision`: `pass` / `fail` / `waiver`; `severity`: `critical` / `major` / `minor` |

### 4.3 `DESIGN-ACCEPTANCE-YYYYMMDD.md`

| 区分 | 定義 |
| --- | --- |
| 保存先 | `memx_spec_v3/docs/reviews/` |
| 命名規則 | `DESIGN-ACCEPTANCE-YYYYMMDD.md` |
| 必須セクション | `# DESIGN ACCEPTANCE REPORT` / `## Metadata` / `## Metrics` / `## Chapter Summary` / `## Final Decision` |
| 必須キー | `run_id`, `generated_at`, `source_commit`, `chapter_id`, `status`, `evidence_paths`, `coverage_rate`, `contract_alignment_high_count`, `link_unreachable_count`, `birdseye_issue_count`, `final_decision` |
| 許可値 | `status`: `pass` / `fail` / `blocked`; `final_decision`: `pass` / `fail`; `coverage_rate`: `0..100`（%） |

## 5. 運用ルール
- テンプレート定義の正本は本仕様とし、個別仕様は本仕様への参照リンクのみを保持する。
- 個別仕様でテンプレートを拡張する場合、まず本仕様を更新してから参照先仕様を更新する。
