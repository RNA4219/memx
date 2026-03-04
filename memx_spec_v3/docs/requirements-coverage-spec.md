# Requirements Coverage Spec

## 1. 目的
本仕様は、`requirements.md` に定義された有効 `REQ-*` に対する `traceability.md` のマッピング充足率（coverage）を算出・判定する手順を固定し、Phase 3 の完了判定を一意にする。

## 2. 用語定義
- **total_req（母数）**: `memx_spec_v3/docs/requirements.md` に存在する有効 `REQ-*` の総数。
- **mapped_req（分子）**: `memx_spec_v3/docs/traceability.md` において、対象 `REQ-*` の行が存在し、かつ 5 列（`Source / Design Mapping / Interface Mapping / Evaluation Mapping / Contract Mapping`）がすべて非空である件数。
- **coverage**: `mapped_req / total_req`。

## 3. 算出ルール（必須）
1. `requirements.md` から `REQ-*` を抽出し、有効 REQ 一覧を作成する。
2. 有効 REQ ごとに `traceability.md` の該当行を確認し、5 列完備時のみ `mapped_req` に加算する。
3. 次式で coverage を算出する。

```text
coverage = mapped_req / total_req
```

- `total_req = 0` は入力不備として `Status: blocked` とする。

## 4. 除外ルール（廃止REQ・waiver中REQ）
- **廃止REQ**:
  - `req-id-lifecycle-spec.md` 4章の必須記録（置換先REQ-ID・移行期限・`CHANGES.md` 連携）を満たす場合に限り、coverage 集計対象から除外してよい。
  - 必須記録が欠ける廃止REQは無効とし、有効REQとして `total_req` に含める。
- **waiver中REQ**:
  - waiver は「判定の一時猶予」であり、マッピング免除ではない。
  - waiver 中であっても有効REQとして `total_req` に含め、5 列マッピング完備でなければ `mapped_req` に含めない。

## 5. 判定閾値と Status 遷移
- Phase 3 完了条件は **coverage = 100%**（`mapped_req == total_req`）で固定する。
- 以下のいずれかを満たす場合、`Status: blocked` へ遷移する。
  1. `coverage < 100%`
  2. 有効REQに対して `traceability.md` の行欠落が 1 件以上ある。
  3. 有効REQに対して 5 列の空欄が 1 セル以上ある。
  4. `total_req = 0` または入力ソース不整合（参照ファイル欠落・破損）がある。

## 6. 証跡フォーマット
判定証跡は YAML と Markdown テーブルの両方を許容する。少なくとも以下項目を含める。

### 6.1 YAML 例
```yaml
coverage_report:
  snapshot_at: "2026-03-04T00:00:00Z"
  source_requirements: "memx_spec_v3/docs/requirements.md"
  source_traceability: "memx_spec_v3/docs/traceability.md"
  total_req: 24
  mapped_req: 24
  coverage: 1.0
  threshold: 1.0
  phase3_gate: pass
  status: done
  excluded:
    deprecated: []
    waiver: []
  blocked_reasons: []
```

### 6.2 Markdown テーブル例
| snapshot_at | total_req | mapped_req | coverage | threshold | phase3_gate | status | blocked_reasons |
| --- | ---: | ---: | ---: | ---: | --- | --- | --- |
| 2026-03-04T00:00:00Z | 24 | 24 | 100% | 100% | pass | done | - |

## 7. `design-review-spec.md` 参照規約
- `memx_spec_v3/docs/design-review-spec.md` の判定根拠セクションでは、本仕様に基づく coverage 証跡（YAML または Markdown テーブル）への参照を必須とする。
- 判定記録には、最低限 `total_req`、`mapped_req`、`coverage`、`status`、`blocked_reasons` を残す。
