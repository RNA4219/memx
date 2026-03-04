---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
---

# design reference resolution spec

`orchestration/memx-design-docs-authoring.md` で使う入力参照名を、Task Seed 作成時に一意な正規パスへ解決するための仕様。

## 0. 適用スコープ

本仕様の参照解決ルールは、以下の成果物すべてに適用する。

- Task Seed
- Phase 1 抽出表（Design Source Inventory）
- 章ドラフト
- レビュー記録の `Source` 欄

## 1. 対象入力参照名と正規パスマッピング

| 入力参照名（非正規） | 正規パス（必須） |
| --- | --- |
| `requirements.md` | `memx_spec_v3/docs/requirements.md` |
| `design.md` | `memx_spec_v3/docs/design.md` |
| `interfaces.md` | `memx_spec_v3/docs/interfaces.md` |
| `traceability.md` | `memx_spec_v3/docs/traceability.md` |
| `EVALUATION.md` | `docs/birdseye/caps/EVALUATION.md.json` |
| `RUNBOOK.md` | `docs/birdseye/caps/RUNBOOK.md.json` |
| `docs/birdseye/index.json` | `docs/birdseye/index.json` |

## 2. `Source` 正規化ルール（`path#Section` 統一）

Task Seed / 章ドラフトの `Source` は、以下をすべて満たすこと。

1. `path#Section` 形式を必須とする（`#Section` が不要な場合も末尾に `#` ではなく章名を明示する）。
2. 相対名（例: `requirements.md#...`）を禁止し、上表の正規パスへ解決した絶対的なリポジトリ相対パスのみ許可する。
3. 曖昧名（例: `EVALUATION.md#...`、`RUNBOOK.md#...`）を禁止する。
4. 複数候補へ解決される参照名は自動補完せず fail とし、Task Seed を `reviewing` のまま差し戻す。

## 3. Task Seed 作成時の必須チェック

Task Seed 起票時に、以下チェックをすべて通過しなければならない。

- `Source` の全行が本仕様の正規パスマッピングに一致している。
- `Source` の全行が `path#Section` 形式で、`#Section` が空でない。
- 相対名・曖昧名・複数候補解決のいずれも 0 件。

## 4. 運用例（誤/正）

- 誤: `requirements.md`
  - 正: `memx_spec_v3/docs/requirements.md#6-4. エラーモデル`
- 誤: `design.md#API`
  - 正: `memx_spec_v3/docs/design.md#3. API設計`
- 誤: `EVALUATION.md#pass-fail`
  - 正: `docs/birdseye/caps/EVALUATION.md.json#pass_fail_rules`
- 誤: `RUNBOOK.md#rollback`
  - 正: `docs/birdseye/caps/RUNBOOK.md.json#rollback`
