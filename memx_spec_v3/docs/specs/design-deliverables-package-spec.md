# Design Deliverables Package Spec

## 1. 目的
本仕様は、**設計書更新1件あたりの必須成果物**を固定し、作成漏れと重複記述を防止する。

## 2. 対象範囲
- 対象は「設計書更新1件あたりの必須成果物」のみ。
- 判定単位は 1 つの設計更新タスク（Task Seed 1件）とする。

## 3. 必須成果物パッケージ（固定）

| 必須ファイル | 必須項目 | 作成/更新タイミング |
| --- | --- | --- |
| `memx_spec_v3/docs/design.md` | 更新対象章、関連 `REQ-ID`、関連 `node_id` | Phase 2（章別ドラフト）で更新、Phase 3/4 で整合確認 |
| `memx_spec_v3/docs/interfaces.md` | 変更対象 I/F、関連 `REQ-ID`、関連 `node_id`、判定に使う契約差分の有無 | Phase 2 で更新、Phase 3 で契約整合を反映 |
| `memx_spec_v3/docs/traceability.md` | `REQ-ID`↔設計/I/F/評価の対応、関連 `node_id` | Phase 1 で更新対象を確定、Phase 2〜3 で更新、Phase 4 で最終確認 |
| `memx_spec_v3/docs/reviews/*.md` | 判定（`pass`/`fail`/`waiver`）、対象 `REQ-ID`、対象 `node_id`、`evidence_paths` | Phase 4 で新規作成/更新（レビュー実施時） |
| `CHANGELOG.md` | `Release Note Draft` を反映した利用者影響要約（1〜3行） | Phase 4 の `Status: done` 直前に更新 |
| `memx_spec_v3/CHANGES.md` | `Release Note Draft` と整合する変更要約、必要時の互換性破壊記録 | Phase 4 の `Status: done` 直前に更新 |

## 4. 運用ルール
- Phase 4 完了判定では、上表の全ファイルが対象タスクに対して更新済みまたは「非該当理由」明記済みであること。
- `docs/TASKS.md` では `Release Note Draft` / `Status` / `Moved-to-CHANGES` を本仕様と整合する形で確認すること。
