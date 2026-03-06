---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
---

# design reference alias spec

Phase 1 入力で許容するエイリアスと正規パスの対応、および解決時の判定ルールを定義する。

## 1. alias 辞書（Phase 1 入力）

| alias | canonical_path | resolution_rule | failure_action |
| --- | --- | --- | --- |
| `requirements.md` | `memx_spec_v3/docs/requirements.md` | 完全一致 | fail |
| `design.md` | `memx_spec_v3/docs/design.md` | 完全一致 | fail |
| `interfaces.md` | `memx_spec_v3/docs/interfaces.md` | 完全一致 | fail |
| `traceability.md` | `memx_spec_v3/docs/traceability.md` | 完全一致 | fail |
| `contracts.md` | `memx_spec_v3/docs/CONTRACTS.md` | 完全一致（大文字小文字ゆれ吸収） | fail |
| `CONTRACTS.md` | `memx_spec_v3/docs/CONTRACTS.md` | 完全一致 | fail |
| `EVALUATION.md` | `docs/birdseye/caps/EVALUATION.md.json` | 完全一致 | fail |
| `RUNBOOK.md` | `docs/birdseye/caps/RUNBOOK.md.json` | 完全一致 | fail |
| `docs/birdseye/index.json` | `docs/birdseye/index.json` | 完全一致 | fail |
| `docs/IN-*.md` | `docs/IN-*.md` | 相対解決（`docs/` 直下の `IN-` プレフィックスのみ許可） | warn |

## 2. 禁止パターン

次の入力は alias 解決前に禁止パターンとして扱う。

- `memx_spec_v3/docs/contracts.md`（小文字ファイル名）
- `memx_spec_v3/docs/EVALUATION.md` / `memx_spec_v3/docs/RUNBOOK.md`
- `../` を含む相対参照
- 絶対パス（`/` 開始）

禁止パターンの `failure_action` はすべて `fail` とする。

## 3. `docs/IN-*.md` の相対解決ルール

- 入力が `docs/IN-*.md#<section>` に一致する場合のみ許可する。
- `canonical_path` は入力 path をそのまま使用する（展開・置換しない）。
- `docs/` 配下でも `IN-` 以外（例: `docs/notes.md`）は alias 辞書外として `warn`。
- 実在しない `docs/IN-*.md` 参照は、存在確認フェーズで `fail`。

## 4. 運用

- 参照解決の実処理順序は `memx_spec_v3/docs/design-reference-resolution-spec.md` に従う。
- 自動検査時の alias 検証観点は `memx_spec_v3/docs/design-reference-validation-automation-spec.md` に従う。
