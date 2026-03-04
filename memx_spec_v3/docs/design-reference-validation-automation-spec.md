---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
---

# design reference validation automation spec

設計書作成前に実行する参照解決の静的検査仕様。`path#section` 解決・参照名正規化・fail 判定を自動化する。

## 1. 入力対象（固定）

参照検査は次の Markdown を入力対象とする。

- `TASK.*.md`
- `orchestration/*.md`
- `memx_spec_v3/docs/design-*.md`

上記以外のファイルは本仕様の自動 fail 判定対象外とし、必要時は別タスクで対象追加する。

## 2. 検査観点

1. `path#section` 形式適合
   - `memx_spec_v3/docs/design-reference-conformance-spec.md` 準拠で検査する。
2. 参照名の正規化
   - `memx_spec_v3/docs/design-reference-resolution-spec.md` の正規パスマッピングへ解決する。
   - alias 判定は `memx_spec_v3/docs/design-reference-alias-spec.md` の辞書（`resolution_rule` / `failure_action`）を正本として適用する。
3. 評価・運用参照の canonical 強制
   - `EVALUATION.md` / `RUNBOOK.md` は `docs/birdseye/caps/*.json` 経由に正規化する。
4. `docs/IN-*.md` 相対解決
   - `docs/IN-*.md#section` 形式のみ許可し、`docs/` 配下の非 `IN-` ファイル参照は alias 辞書に従って warn/fail 判定する。

## 3. fail 条件（自動停止）

次のいずれか 1 件でも検出した場合は fail とする。

1. `memx_spec_v3/docs/EVALUATION.md` 参照
   - `memx_spec_v3/docs/EVALUATION.md` または `memx_spec_v3/docs/EVALUATION.md#...` が入力対象に存在する。
2. `path#section` 非準拠
   - `path` のみ、`#section` 欠落、または許可外パスの混在。
3. テンプレート ID 混入
   - `IN-YYYYMMDD-001` などテンプレート識別子を実績 ID として参照している。
4. 未解決 `contracts.md` 表記ゆれ
   - `memx_spec_v3/docs/contracts.md`（小文字）や `contracts.md` が、`memx_spec_v3/docs/CONTRACTS.md` へ解決されず残存している。
5. 禁止パターン参照
   - alias spec で禁止された `../`、絶対パス、`memx_spec_v3/docs/EVALUATION.md` / `RUNBOOK.md` を検出した。

## 3-1. alias 辞書の自動検証観点

- 入力参照文字列ごとに alias 辞書の `resolution_rule` を適用し、完全一致・相対解決・禁止パターンを分類する。
- 分類結果に対して alias 辞書の `failure_action`（warn/fail）をそのまま判定へ反映する。
- `warn` は issue レコードを残して継続、`fail` は即時停止とする。

## 4. `docs/birdseye/memx-birdseye-validation-spec.md` との責務境界

責務重複を避けるため、検証責務を次の通り固定する。

- 本仕様（design reference validation automation）
  - Markdown 内の参照文字列検査（`path#section`、canonical 解決、表記ゆれ検出）。
  - 対象: `TASK.*.md` / `orchestration/*.md` / `memx_spec_v3/docs/design-*.md`。
- `docs/birdseye/memx-birdseye-validation-spec.md`
  - `docs/birdseye/index.json` の `node_id` / `depends_on` / `capsule` の整合検証。
  - 対象: index/caps の node・caps 実体。

本仕様は node/caps 自体の正当性を判定しない。Birdseye 側は Markdown 参照解決の妥当性を判定しない。

## 5. 出力と運用

- 出力は issue レコード（1 fail = 1 レコード）で記録する。
- fail が 1 件以上の場合、設計書作成タスクの `Status` は `blocked` に遷移する。
- 再検査で fail 0 件を確認した時のみ `planned/active/in_progress/reviewing` に復帰できる。
