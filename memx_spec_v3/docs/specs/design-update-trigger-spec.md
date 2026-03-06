# Design Update Trigger Spec

## 目的
- Task Seed 起票前に「どの変更が、どの文書更新を必須化するか」を統一判定し、`Status: done` 遷移条件と `CHANGELOG.md` / `memx_spec_v3/CHANGES.md` 反映可否を一致させる。

## トリガー種別（判定対象）
- **TRG-REQ**: `requirements.md` の REQ 追加/変更/削除。
- **TRG-OAS**: OpenAPI 差分（endpoint / request / response / error / code 変更）。
- **TRG-CLI-SCHEMA**: CLI `--json` schema 差分。
- **TRG-RUNBOOK**: `RUNBOOK.md` の運用手順・監視・復旧導線の変更。
- **TRG-IN-PREVENTIVE**: `docs/IN-*.md` 再発防止策の新規追加・更新。

## 必須更新先マトリクス
凡例: `必須` / `条件付き` / `不要`

| Trigger | design (`design.md`) | interfaces (`interfaces.md`) | traceability (`traceability.md`) | EVALUATION (`EVALUATION.md`) | operations (`RUNBOOK.md`) | レビュー記録 (`memx_spec_v3/docs/reviews/`) |
| --- | --- | --- | --- | --- | --- | --- |
| TRG-REQ | 必須 | 条件付き（I/F影響時） | 必須 | 条件付き（評価観点変更時） | 条件付き（運用影響時） | 必須 |
| TRG-OAS | 条件付き（設計影響時） | 必須 | 必須 | 条件付き（評価観点変更時） | 条件付き（運用導線変更時） | 必須 |
| TRG-CLI-SCHEMA | 条件付き（設計影響時） | 必須 | 必須 | 条件付き（評価観点変更時） | 条件付き（運用導線変更時） | 必須 |
| TRG-RUNBOOK | 条件付き（設計前提変更時） | 条件付き（I/F記述へ波及時） | 条件付き（REQ対応の証跡更新時） | 条件付き（運用評価軸変更時） | 必須 | 必須 |
| TRG-IN-PREVENTIVE | 条件付き（恒久対策が設計変更を伴う時） | 条件付き（契約変更を伴う時） | 必須 | 必須 | 必須 | 必須 |

## IA仕様との対応
- IA の責務境界・文書区分・更新優先順位は [design-doc-ia-spec.md](./design-doc-ia-spec.md) を参照する。
- 本書は Trigger 判定と Done 遷移条件の正本、`design-doc-ia-spec.md` は文書責務境界の正本として相互参照する。

## 必須更新先マトリクス（IA連携）
凡例: `必須` / `条件付き` / `不要`

| Trigger | requirements | design | interfaces | traceability | CONTRACTS (machine-readable) | error-contract | operations | design-doc-ia-spec |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TRG-REQ | 必須 | 必須 | 条件付き | 必須 | 条件付き | 条件付き | 条件付き | 必須 |
| TRG-OAS | 条件付き | 条件付き | 必須 | 必須 | 必須 | 必須 | 条件付き | 必須 |
| TRG-CLI-SCHEMA | 条件付き | 条件付き | 必須 | 必須 | 必須 | 条件付き | 条件付き | 必須 |
| TRG-RUNBOOK | 条件付き | 条件付き | 条件付き | 条件付き | 不要 | 条件付き | 必須 | 必須 |
| TRG-IN-PREVENTIVE | 条件付き | 条件付き | 条件付き | 必須 | 条件付き | 必須 | 必須 | 必須 |

## CHANGELOG / CHANGES 反映要否判定

### 判定ルール
- **要反映（両方）**: 利用者影響あり（API/CLI/`--json`/運用手順/互換性）またはインシデント起因の恒久対策。
- **要反映（`memx_spec_v3/CHANGES.md` のみ）**: v3設計内部の補助情報更新で、利用者影響なし。
- **反映不要**: 誤字修正・リンク修正のみで、仕様意味・運用意味が不変。

### 破壊変更の強制ルール
- CLI/API/`--json` の互換性破壊がある場合、`CHANGELOG.md` と `memx_spec_v3/CHANGES.md` の**双方同日反映を必須**とする。

## `Status: done` 遷移条件との整合
Task Seed は次をすべて満たした場合のみ `done` へ遷移できる。

1. 起票時にトリガー判定（複数可）を記録済み。
2. 判定されたトリガーに対し、本仕様の必須更新先がすべて更新済み。
3. `Release Note Draft` を記載済み。
4. 本仕様の「CHANGELOG / CHANGES 反映要否判定」に従って転記済み。
5. 転記対象がある場合、`Moved-to-CHANGES: YYYY-MM-DD` を追記済み。

## Task Seed 記載フォーマット（推奨）
```md
### Trigger 判定
- Trigger IDs: TRG-REQ, TRG-OAS
- Impact: 利用者影響あり（API response 変更）
- Required Updates:
  - [x] memx_spec_v3/docs/design.md
  - [x] memx_spec_v3/docs/interfaces.md
  - [x] memx_spec_v3/docs/traceability.md
  - [ ] EVALUATION.md（今回非該当: 評価観点変更なし）
  - [x] RUNBOOK.md
  - [x] memx_spec_v3/docs/reviews/DESIGN-REVIEW-YYYYMMDD-###.md
- Changelog Decision: CHANGELOG.md + memx_spec_v3/CHANGES.md
```
