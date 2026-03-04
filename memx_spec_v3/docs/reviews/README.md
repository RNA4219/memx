# Design Review Records

## 保存先
- 設計レビュー記録は `memx_spec_v3/docs/reviews/` に保存する。

## 命名規則
- ファイル名は `DESIGN-REVIEW-YYYYMMDD-###.md` とする。
  - `YYYYMMDD`: レビュー実施日（ローカル日付）
  - `###`: 同日内の 001 始まり連番
- 例: `DESIGN-REVIEW-20260304-001.md`

## 更新ルール
- 新規レビューごとに `TEMPLATE.md` を複製して新規ファイルを作成する（既存記録の上書き禁止）。
- 記録作成時は `design-review-spec.md` の記録テンプレート（`# DESIGN REVIEW: <title>` で始まるテンプレート）を必ず利用し、必須 6 項目（対象章/関連 REQ-ID/Node IDs/指摘一覧（重大度付き）/再確認結果/判定）を記入する。
- 判定欄には `EVALUATION.md` の該当ルール参照と証跡を必ず記載する。
- `docs/TASKS.md` 連携項目（`Release Note Draft` / `Status` / `Moved-to-CHANGES`）を記録クローズ前に更新する。

## レビュー起票トリガー（必須）
- 以下のいずれかに差分がある PR は、設計レビュー記録の新規起票を必須とする。
  - `memx_spec_v3/docs/requirements.md`
  - `memx_spec_v3/docs/design.md`
  - `memx_spec_v3/docs/interfaces.md`
  - `memx_spec_v3/docs/CONTRACTS.md`
- 上記に該当しない差分でも、`REQ-*` の追加/更新、設計判断、外部 I/F 契約変更を含む場合は起票する。
- 上記トリガー差分に `memx_spec_v3/docs/contracts.md`（小文字）が含まれる場合は、正規パスへ未修正の誤参照として `fail` 扱いで差し戻す（正: `memx_spec_v3/docs/CONTRACTS.md`）。

## レビュー単位と必須添付
- レビュー単位は **章単位** を標準とし、必要に応じて **PR 単位** で集約してよい。
- いずれの単位でも、以下を記録に必須添付する。
  - 実行コマンドと結果（例: `git diff --name-only <base>...HEAD`、検証コマンド出力）
  - 証跡リンク（ログ、コメント ID、成果物パス等）

## 判定語彙とエスカレーション基準
- 判定語彙は `pass` / `fail` / `waiver` の 3 種のみを許可する。
- エスカレーション（`fail` または `waiver` 必須）条件:
  - `critical` 指摘が 1 件以上残存する。
  - `major` 指摘が 2 件以上残存する。
  - 必須証跡（コマンド結果/証跡リンク）が不足する。
  - `requirements.md` と `traceability.md` のマッピング不整合が 1 件以上ある。
- `waiver` を選ぶ場合は `design-review-spec.md` の waiver 条件に従い、`docs/IN-<実日付>-<連番>.md` を必須参照とする。


## 受け入れレポート運用ルール（必須）
- `DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用であり、受け入れ判定の実体ファイルとしては使用しない。
- リリース判定時は毎回 `DESIGN-ACCEPTANCE-<実日付>.md` を `memx_spec_v3/docs/reviews/` に新規作成する（既存ファイル上書き・テンプレート流用禁止）。
- 受け入れレポートでは、必須6項目（対象章/REQ網羅率/high差分件数/リンク不達件数/Birdseye issue件数/最終判定）に加え、`evidence_paths` が実在ファイルのみを指すことの検証結果を必須記載とする。

## クローズ条件
- 記録クローズは以下を全て満たした時点とする。
  - 判定が確定し、根拠参照（`EVALUATION.md` + 証跡）が記載済み。
  - `docs/TASKS.md` の `Release Note Draft` / `Status` が更新済み。
  - `docs/TASKS.md` の `Moved-to-CHANGES` が反映済み（`Moved-to-CHANGES: YYYY-MM-DD`）。
