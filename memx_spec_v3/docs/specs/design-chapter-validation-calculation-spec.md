# Design Chapter Validation Calculation Spec

## 1. 目的
本仕様は、`DESIGN-CHAPTER-VALIDATION` レポートの各列の算出式・入力抽出キー・再計算トリガー・出力更新ルールを固定し、章単位の検証結果を再現可能にする。

## 2. 適用範囲
- `memx_spec_v3/docs/design-chapter-validation-spec.md` で定義する章別検証サマリ。
- `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` を参照する `mapping_match_check` 判定。
- 出力先 `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-<YYYYMMDD>.md`。

## 3. 入力ソースと抽出キー（固定）

### 3.1 抽出キー
抽出キーは以下に固定する。
- `chapter_id`
- `REQ-ID`
- `contract_ref`
- `node_id`
- `evidence_paths`

### 3.2 読み取り対象ファイル
| 抽出キー | 主読取元 | 補助読取元 | 用途 |
| --- | --- | --- | --- |
| `chapter_id` | `memx_spec_v3/docs/design.md` / `memx_spec_v3/docs/interfaces.md` | `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` | 章単位レコードの識別と対応表照合 |
| `REQ-ID` | `memx_spec_v3/docs/requirements.md` | `memx_spec_v3/docs/traceability.md` | 章に割当済み要件集合の特定 |
| `contract_ref` | `memx_spec_v3/docs/interfaces.md` | `memx_spec_v3/docs/traceability.md` | 契約整合（high）件数算出 |
| `node_id` | `docs/birdseye/index.json` | `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` | 章対応ノード一致判定 |
| `evidence_paths` | `memx_spec_v3/docs/reviews/` 配下のレビュー/受入レポート | - | 必須証跡の存在確認 |


## 3.3 `DESIGN-CHAPTER-VALIDATION-<実日付>.md` 列別入力ソース固定
`memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-<実日付>.md` の列は、以下の入力元以外を使用してはならない。

| 列 | 入力ソース（固定） | 参照節（固定） |
| --- | --- | --- |
| `chapter_id` | `memx_spec_v3/docs/design.md` / `memx_spec_v3/docs/interfaces.md` | 章見出し（`path#section`） |
| `req_coverage` | `memx_spec_v3/docs/traceability.md` | 対象 `REQ-ID` の設計反映状態行 |
| `contract_alignment_high_count` | `memx_spec_v3/docs/reviews/CONTRACT-ALIGN-<実日付>-<連番>.md` または `LATEST.md` | `severity: high` 一覧 |
| `link_broken_count` | `memx_spec_v3/docs/reviews/DESIGN-REVIEW-<実日付>-<連番>.md` | リンク検証結果節（broken一覧） |
| `birdseye_issue_count` | `docs/birdseye/index.json` | 対象 `node_id` の unresolved issue 一覧 |
| `evidence_paths` | `memx_spec_v3/docs/reviews/DESIGN-REVIEW-<実日付>-<連番>.md` / `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<実日付>.md` | ファイル実体 |
| `mapping_spec_ref` | `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` | 章対応表定義 |
| `mapping_match_check` | `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` + `docs/birdseye/index.json` | 章対応表行 / node解決結果 |
| `updated_at` | 当日更新ログ | レポートヘッダ時刻 |
| `calculation_basis` | 当日計算根拠 | 実行コミットIDまたは実行時刻 |

## 4. 列ごとの算出式

### 4.1 `req_coverage`
- 定義:
  - `chapter_req_total` = 対象 `chapter_id` に割当済みの `REQ-ID` 総数。
  - `chapter_req_covered` = 上記のうち、`traceability.md` 上で「設計反映済み」と判定できる `REQ-ID` 数。
- 算出式:
  - `req_coverage = (chapter_req_covered / chapter_req_total) * 100`
- 端数処理:
  - 小数第1位を四捨五入し整数 `%` で記録。
- 例外:
  - `chapter_req_total = 0` の場合は `0%`。

### 4.2 `contract_alignment_high_count`
- 定義:
  - 対象 `chapter_id` に紐づく `contract_ref` を抽出し、契約整合チェック結果から `severity: high` を集計。
- 算出式:
  - `contract_alignment_high_count = count(severity == "high" and chapter_id == target)`

### 4.3 `link_broken_count`
- 定義:
  - 対象 `chapter_id` 内のリンクと、当該章から参照する章間リンクのうち不達（404/未存在アンカー/未存在ファイル）件数。
- 算出式:
  - `link_broken_count = count(link_status == "broken" and source_chapter_id == target)`

### 4.4 `birdseye_issue_count`
- 定義:
  - `docs/birdseye/index.json` 起点で対象 `node_id` に紐づく未解決 issue 件数。
- 算出式:
  - `birdseye_issue_count = count(issue_status != "resolved" and node_id == mapped_node_id)`

### 4.5 `mapping_match_check`
- 判定値: `pass` / `fail`
- `pass` 条件（`design-chapter-node-mapping-spec.md` 整合）:
  1. `chapter_id` が章対応表で一意に確定できる（`ambiguous` / `missing` でない）。
  2. 確定した `chapter_id` に対し `node_id` が `docs/birdseye/index.json` で一意に解決できる。
  3. 解決した `node_id` が章対応表の `node_id` と一致する。
- `fail` 条件:
  - 上記いずれかを満たさない場合。
  - 特に `design-chapter-node-mapping-spec.md` 5.1.2 の fail 条件（`ambiguous` / `missing` / 一意化不可）に該当する場合。

## 5. 再計算トリガー
以下の更新を検知した時点で、全 `chapter_id` を再計算対象にする。
- `memx_spec_v3/docs/requirements.md` 更新時
- `memx_spec_v3/docs/design.md` 更新時
- `memx_spec_v3/docs/interfaces.md` 更新時
- `memx_spec_v3/docs/traceability.md` 更新時
- `docs/birdseye/index.json`（Birdseye）更新時

## 6. 出力先・必須更新フィールド・差分レビュー観点

### 6.1 出力先（固定）
- `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-<YYYYMMDD>.md`

### 6.2 必須更新フィールド
再計算時は対象日付ファイルで以下を必須更新とする。
- `chapter_id`
- `req_coverage`
- `contract_alignment_high_count`
- `link_broken_count`
- `birdseye_issue_count`
- `evidence_paths`
- `mapping_spec_ref`
- `mapping_match_check`
- `updated_at`（レポートヘッダまたは各行メタ）
- `calculation_basis`（対象コミット or 実行時刻）

### 6.3 差分レビュー観点
- 数値差分:
  - `req_coverage` 低下、`contract_alignment_high_count` 増加、`link_broken_count` 増加、`birdseye_issue_count` 増加を優先確認。
- 判定差分:
  - `mapping_match_check: pass -> fail` は即時レビュー対象。
- 証跡差分:
  - `evidence_paths` が実在ファイルを参照していること。
- 対応表整合:
  - `mapping_spec_ref` が固定値 `memx_spec_v3/docs/design-chapter-node-mapping-spec.md` であること。


## 7. `0%` / `fail` のまま完了できない条件（close禁止）
以下のいずれかに該当する場合、`DESIGN-CHAPTER-VALIDATION-<実日付>.md` は完了扱いにしてはならず、Task Seed の `Status: done` への遷移を禁止する。

1. `req_coverage = 0%` かつ `chapter_req_total > 0`。
2. `mapping_match_check = fail`。
3. `contract_alignment_high_count > 0`。
4. `link_broken_count > 0`。
5. `birdseye_issue_count > 0`。
6. `evidence_paths` が 2 件未満、または実在しないファイルを含む。

差し戻し時は、未充足列・入力元ファイル・再計算日時を `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-<実日付>.md` に追記する。
