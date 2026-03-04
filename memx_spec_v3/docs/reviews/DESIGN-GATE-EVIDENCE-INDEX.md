# DESIGN GATE EVIDENCE INDEX

5条件の参照先固定（レビュー時の一次参照）を本書で定義する。

## req_coverage
- 判定値: `100%` | 根拠ファイル: `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-20260304.md` | 根拠箇所（見出し名）: `## 2. REQ網羅率`

## mapping_match_check
- 判定値: `pass`（全章） | 根拠ファイル: `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304-002.md` | 根拠箇所（見出し名）: `## mapping_match_check 判定ログ`

## contract_align_report
- 判定値: `high_count=0, phase3_status=done` | 根拠ファイル: `memx_spec_v3/docs/contracts/reports/LATEST.md` | 根拠箇所（見出し名）: `先頭キー（report_id/report_path/decision_date/high_count/phase3_status）`

## design_acceptance
- 判定値: `pass` | 根拠ファイル: `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-20260304.md` | 根拠箇所（見出し名）: `## 6. 最終判定`

## go.sum tracked
- 判定値: `tracked` | 根拠ファイル: `memx_spec_v3/go/go.sum` | 根拠箇所（見出し名）: `N/A（ファイル実体）` | 運用確認コマンド: `git ls-files memx_spec_v3/go/go.sum`（出力: `memx_spec_v3/go/go.sum`）
- 補足（継続確認）: 単発確認ではなく、監査ログ `memx_spec_v3/docs/reviews/GO-LOCK-AUDIT-20260304.md` を更新し継続記録する。
