# CONTRACT ALIGN REPORT: CONTRACT-ALIGN-20260304

- report_id: `CONTRACT-ALIGN-20260304`
- run_id: `contract-align-20260304-001`
- generated_at: `2026-03-04T09:20:00Z`
- source_commit: `HEAD`
- tool: `manual-review`
- status: `pass`
- source_inputs:
  - `memx_spec_v3/docs/design.md`
  - `memx_spec_v3/docs/interfaces.md`
  - `memx_spec_v3/docs/contracts/openapi.yaml`
  - `memx_spec_v3/docs/contracts/cli-json.schema.json`
  - `memx_spec_v3/docs/error-contract.md`
  - `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304-003.md`

## 1. diff_summary（最小必須）
- high_count: `0`
- medium_count: `0`
- low_count: `0`
- diff_summary: `差分なし（全件整合）`

## 2. chapter 集計（contract_alignment_high_count）
`DESIGN-CHAPTER-VALIDATION-20260304-003.md` の章別表と一致。

| chapter_id | contract_alignment_high_count |
| --- | ---: |
| `memx_spec_v3/docs/design.md#1. レイヤ構成` | `0` |
| `memx_spec_v3/docs/design.md#2. DB 責務分割` | `0` |
| `memx_spec_v3/docs/design.md#3. 移行戦略` | `0` |
| `memx_spec_v3/docs/design.md#4. ユースケース設計` | `0` |
| `memx_spec_v3/docs/design.md#5. ADR参照運用ルール` | `0` |
| `memx_spec_v3/docs/design.md#6. 設計→契約→検証 導線（要件ID単位）` | `0` |
| `memx_spec_v3/docs/design.md#7. design-template 段階移行チェックリスト（章単位）` | `0` |
| `memx_spec_v3/docs/interfaces.md#0. 文書の位置づけ` | `0` |
| `memx_spec_v3/docs/interfaces.md#1. CLI I/O（v1 必須）` | `0` |
| `memx_spec_v3/docs/interfaces.md#2. API I/O（v1 必須）` | `0` |
| `memx_spec_v3/docs/interfaces.md#3. 互換ルール` | `0` |
| `memx_spec_v3/docs/interfaces.md#4. エラー面` | `0` |
| `memx_spec_v3/docs/interfaces.md#5. 契約変更手順（更新順序固定）` | `0` |
| `memx_spec_v3/docs/interfaces.md#6. 付録: RUNBOOK連携 I/F ID（v1運用）` | `0` |

## 3. final_decision（最小必須）
- final_decision: `pass`
- decision_rule: `high_count = 0 のため Phase 3 契約整合は pass`
- evidence_paths:
  - `memx_spec_v3/docs/reviews/CONTRACT-ALIGN-20260304.md`
  - `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304-003.md`
