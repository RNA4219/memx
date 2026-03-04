# Design Review Record Spec

## 1. 目的
本仕様は、設計レビュー記録の必須項目・判定根拠・保存規約を固定し、`EVALUATION.md` と `docs/TASKS.md` に整合した受け入れ判断を行うための共通フォーマットを定義する。

## 2. 保存先と命名規則（必須）
- 保存先は **`memx_spec_v3/docs/reviews/`** に固定する。
- ファイル名は **`DESIGN-REVIEW-YYYYMMDD-###.md`** とする。
  - `YYYYMMDD`: レビュー実施日（ローカル日付）
  - `###`: 同日内の 001 始まり連番
- 例: `memx_spec_v3/docs/reviews/DESIGN-REVIEW-20260304-001.md`

## 3. レビュー記録の必須フィールド
各レビュー記録は、以下 6 項目を必須で含む。

1. **対象章**
   - 例: `memx_spec_v3/docs/design.md#3. データフロー`
2. **関連 REQ-ID**
   - 章で検証した `REQ-*` を列挙する。
3. **Node IDs**
   - `docs/birdseye/index.json` の `node_id` を列挙する。
4. **指摘一覧（重大度付き）**
   - 各指摘に `severity` を必須付与（`critical` / `major` / `minor`）。
5. **再確認結果**
   - 指摘ごとの再確認結果を `resolved` / `remaining` で明記する。
6. **判定（pass/fail/waiver）**
   - レビュー全体の最終判定を 1 つ記載する。

## 4. `EVALUATION.md` 連携（判定根拠参照を必須化）
- `pass` / `fail` / `waiver` のいずれの判定でも、**判定根拠として `EVALUATION.md` の該当 pass/fail ルール参照を必須**とする。
- 根拠参照は次を最低限含む。
  - 対応 `REQ-ID`
  - `EVALUATION.md` 内の該当節またはアンカー（例: `#req-cli-001-passfail`）
  - 判定に使用した証跡（コマンド結果、計測結果、レビューコメント ID など）
- 要件マッピング充足率の判定根拠は `memx_spec_v3/docs/requirements-coverage-spec.md` に従い、coverage 証跡（YAML または Markdown テーブル）参照を必須とする。
- `waiver` の場合は、`EVALUATION.md` の waiver 条件に従い `docs/IN-<実日付>-<連番>.md` 参照を必須とする。

## 5. `docs/TASKS.md` 連携（レビュー完了条件）
レビュー完了（記録クローズ）条件に、以下 3 点を必須で含める。

1. **Release Note Draft**
   - 対応 Task Seed に利用者影響の要約（1〜3 行）が記入済みであること。
2. **Status**
   - 対応 Task Seed の `Status` が `reviewing` から `done` へ遷移可能条件を満たしていること。
3. **Moved-to-CHANGES**
   - 変更移送後、Task Seed に `Moved-to-CHANGES: YYYY-MM-DD` が追記済みであること。

## 6. 差分レビュー時の未マッピングREQ検出手順（必須）
- 目的: `requirements.md` で追加・変更された `REQ-*` が `traceability.md` に未反映のままマージされることを防止する。
- 参照仕様: `memx_spec_v3/docs/req-id-lifecycle-spec.md` 5章「`traceability.md` 同時更新必須ルール（未マッピング時は fail）」。
- 確認対象差分:
  1. `git diff --name-only <base>...HEAD` で `memx_spec_v3/docs/requirements.md` が含まれるか確認する。
  2. 含まれる場合、`requirements.md` 差分で追加/変更された `REQ-*` を列挙する。
  3. `traceability.md` の表（主要REQ + 主要REQ以外）に同一 ID が 1 行ずつ存在することを確認する。
- 判定観点（レビューコメントに残す）:
  - `memx_spec_v3/docs/req-id-lifecycle-spec.md` 5章に従い、`Source / Design Mapping / Interface Mapping / Evaluation Mapping / Contract Mapping` の 5 列が埋まっている。
  - `Design Mapping` が必ず `memx_spec_v3/docs/design.md#...` を指している。
  - 同一 PR で `requirements.md` と `traceability.md` の更新が同時に入っている。
- 未マッピングが 1 件でもある場合は `fail` とし、解消後に再レビューする。

## 7. 記録テンプレート（最小）
```md
# DESIGN REVIEW: <title>
- Review ID: DESIGN-REVIEW-YYYYMMDD-###
- 対象章: <path#section>
- 関連 REQ-ID:
  - REQ-...
- Node IDs:
  - node-...

## 指摘一覧（重大度付き）
- [DR-001] severity: major
  - 指摘: ...
  - 対応: ...
  - 再確認結果: resolved

## 判定
- 判定: pass|fail|waiver
- 根拠参照:
  - EVALUATION: <REQ-ID / anchor>
  - 証跡: <command-log / artifact / comment-id>
  - （waiver時）IN: docs/IN-<実日付>-<連番>.md

## TASKS 連携確認
- Release Note Draft: done|pending
- Status: reviewing|done
- Moved-to-CHANGES: YYYY-MM-DD|pending
```
