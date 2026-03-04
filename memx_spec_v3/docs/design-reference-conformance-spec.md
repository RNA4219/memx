---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
---

# design reference conformance spec

`memx_spec_v3/docs/design-reference-resolution-spec.md` で正規化された参照が、成果物上で適合しているかを判定する仕様。

## 0. 役割分担

- 解決ルール（入力参照名 -> 正規パス変換）は `memx_spec_v3/docs/design-reference-resolution-spec.md` を正本とする。
- 本仕様は、成果物に記載された参照文字列の適合判定（warn/fail）を正本とする。

## 1. 許可参照形式（必須）

- 許可形式は `path#section` のみとする。
- `path` はリポジトリ相対の正規パスを必須とし、解決は `memx_spec_v3/docs/design-reference-resolution-spec.md` のマッピングに従う。
- `section` は空文字を禁止し、対象仕様の見出しまたはキーを明示する。

## 2. 禁止形式

以下を禁止し、検出時は本仕様 3 章の判定レベルを適用する。

- 曖昧名（例: `requirements.md#...`, `EVALUATION.md#...`, `RUNBOOK.md#...`）
- 非正規パス（例: `memx_spec_v3/docs/EVALUATION.md#...`, `memx_spec_v3/docs/RUNBOOK.md#...`）
- `path#section` 以外（例: `path` 単体、`#section` 欠落、section 空）

## 3. 判定レベル（warn/fail）

- `warn`: 参照候補は一意だが表記が非正規（例: 曖昧名、別名、拡張子違い）。修正タスクを同一 Phase 内で必須化する。
- `fail`: 解決不能、複数候補、非正規パス固定違反、または `path#section` 不成立。対象成果物の完了判定を停止する。

## 4. Phase 別適用（必須）

| Phase | 対象成果物 | 判定タイミング | fail 時の扱い |
| --- | --- | --- | --- |
| Phase 1 | 抽出表（Design Source Inventory） | 抽出表更新時 | 抽出表を未完了のまま差し戻し |
| Phase 2 | 章ドラフト | 章ドラフト作成/更新時 | 当該章 Task Seed を `reviewing` 維持 |
| Phase 4 | レビュー記録 | レビュー記録作成時 | レビュー判定を `fail` 固定 |

## 5. 完了条件

- Phase 1 抽出表・Phase 2 章ドラフト・Phase 4 レビュー記録の参照が、すべて `path#section` で記録されている。
- 上記成果物における参照解決適合率（適合件数/総参照件数）が `100%` である。
