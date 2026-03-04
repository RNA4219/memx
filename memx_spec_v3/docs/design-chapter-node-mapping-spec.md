# Design Chapter Node Mapping Spec

## 1. 目的
本仕様は `chapter_id` の命名規則と `chapter_id -> node_id` 対応表の標準フォーマットを定義し、`docs/birdseye/index.json` 更新時に章ドラフト・検証サマリ・レビュー記録の追随を非破壊で実施できる状態を維持する。

## 2. 適用範囲
- `memx_spec_v3/docs/design.md` を起点とする章ドラフト運用。
- `orchestration/memx-design-docs-authoring.md` の Phase 1/2/3。
- `memx_spec_v3/docs/design-chapter-validation-spec.md` で参照する章別検証サマリ。
- `docs/birdseye/index.json` 由来の `node_id`・`depends_on` 追随運用。

## 3. `chapter_id` 命名規則

### 3.1 形式（安定ID）
- `chapter_id` は `path#anchor_slug` 形式で固定する。
- `path` は `memx_spec_v3/docs/design-reference-resolution-spec.md` で正規化済みの相対パスを使う。
- `anchor_slug` は見出し表示名を直接使わず、初回採番時に確定した slug を継続利用する。
- 例: `memx_spec_v3/docs/design.md#chapter-03-data-flow`

### 3.2 表示名変更時の非破壊方針
- 見出し表示名を変更しても `chapter_id` は変更しない。
- 追跡用に対応表へ `display_title`（現在表示名）を更新し、`chapter_id` は据え置く。
- 表示名変更で slug を再計算してはならない。

### 3.3 廃止時の扱い
- 章を廃止する場合は対応表から即時削除せず、`status: deprecated` を設定する。
- 廃止章の `node_id` は `replacement_chapter_id` がある場合のみ引き継ぎ可能とし、履歴行を残す。
- 完全削除は「2リリース経過かつ参照ゼロ」を満たしたレビュー承認後に実施する。

## 4. `chapter_id -> node_id` 対応表フォーマット

### 4.1 管理形式
- 章対応表は Markdown テーブルで管理し、列順を固定する。
- 必須列:
  1. `chapter_id`
  2. `display_title`
  3. `node_id`
  4. `depends_on`
  5. `status` (`active` / `deprecated`)
  6. `last_verified_at` (UTC RFC3339)
  7. `review_note`

### 4.2 最小テンプレート
| chapter_id | display_title | node_id | depends_on | status | last_verified_at | review_note |
| --- | --- | --- | --- | --- | --- | --- |
| memx_spec_v3/docs/design.md#chapter-03-data-flow | 3. データフロー | design-dataflow | requirements-core,interfaces-contract | active | 2026-03-04T00:00:00Z | initial |


### 4.3 章対応表（2026-03-04 再検証）
| chapter_id | display_title | node_id | depends_on | status | last_verified_at | review_note |
| --- | --- | --- | --- | --- | --- | --- |
| `memx_spec_v3/docs/design.md#1. レイヤ構成` | `1. レイヤ構成` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/design.md#2. DB 責務分割` | `2. DB 責務分割` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/design.md#3. 移行戦略` | `3. 移行戦略` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/design.md#4. ユースケース設計` | `4. ユースケース設計` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/design.md#5. ADR参照運用ルール` | `5. ADR参照運用ルール` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/design.md#6. 設計→契約→検証 導線（要件ID単位）` | `6. 設計→契約→検証 導線（要件ID単位）` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/design.md#7. design-template 段階移行チェックリスト（章単位）` | `7. design-template 段階移行チェックリスト（章単位）` | `design` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#0. 文書の位置づけ` | `0. 文書の位置づけ` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#1. CLI I/O（v1 必須）` | `1. CLI I/O（v1 必須）` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#2. API I/O（v1 必須）` | `2. API I/O（v1 必須）` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#3. 互換ルール` | `3. 互換ルール` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#4. エラー面` | `4. エラー面` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#5. 契約変更手順（更新順序固定）` | `5. 契約変更手順（更新順序固定）` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |
| `memx_spec_v3/docs/interfaces.md#6. 付録: RUNBOOK連携 I/F ID（v1運用）` | `6. 付録: RUNBOOK連携 I/F ID（v1運用）` | `interfaces` | `requirements` | `active` | `2026-03-04T08:40:55Z` | `chapter coverage recalculation` |

## 5. `docs/birdseye/index.json` 更新時の追随ルール

### 5.0 `source_path#section` から `chapter_id`・`node_id` を解決する標準手順
`source_path#section` からの解決は、以下の手順を**順番固定**で実施する。

1. `source_path` 正規化
   - `memx_spec_v3/docs/design-reference-resolution-spec.md` に従い、相対パスへ正規化する。
2. `chapter_id` 候補抽出
   - 章対応表の `chapter_id` 先頭（`path` 部分）が、正規化後 `source_path` と一致する行を候補とする。
3. `chapter_id` 確定
   - `source_path#section` の `section` が既存 `chapter_id` の `anchor_slug` に対応する場合は当該行を採用する。
   - 同一 `source_path` 配下で複数候補になる場合は、`display_title` が `section` と最も一致する行を採用する。
   - なお単一候補のみの場合は、その候補を採用する。
4. `node_id` 解決（`docs/birdseye/index.json`）
   - `chapter_id` 確定後、`index.json` から `node_id` を探索する。
   - 利用フィールドと優先順位は **5.1.1** を正本とする。
5. 結果判定
   - 一意に解決できた場合のみ `resolved`。
   - 複数候補が残る場合は `ambiguous`。
   - 候補ゼロ、または `index.json` に該当 `node_id` が存在しない場合は `missing`。
   - `ambiguous` / `missing` は Phase 1 exit 不可（`design-source-inventory-spec.md` 参照）。

### 5.1 差分検知
- 更新時は旧版/新版の `node_id` と `depends_on` を比較し、以下を区分する。
  - 追加: 新規 `node_id`
  - 変更: 既存 `node_id` の `depends_on` 変更
  - 削除: 旧版にのみ存在する `node_id`
- 差分検知結果は対応表の `review_note` に要約を残す。

### 5.1.1 `docs/birdseye/index.json` の node 解決フィールド優先順位
`chapter_id` に対応する `node_id` 解決時は、`index.json` のノード情報に対して以下の順で照合する。

1. `source_path` 完全一致
   - ノード側 `source_path` が `chapter_id` の `path` と完全一致するものを最優先候補にする。
2. `section` / `anchor` 一致
   - 1 の候補内で、ノード側 `section` または `anchor` が `chapter_id` の `anchor_slug` に一致するものを優先する。
3. `title` 近似一致
   - 2 で一意化できない場合、ノード側 `title` と章対応表 `display_title` の一致度で順位付けする。
4. `node_id` 安定継承
   - 3 でも同率の場合、既存章対応表の `node_id` と一致する候補を優先する。

上記で一意化できない場合は `ambiguous`、候補ゼロの場合は `missing` とする。

### 5.1.2 fail 条件（未解決時）
以下のいずれかに該当した場合、node 解決は fail とし、Phase 1 の完了判定を通してはならない。

- `chapter_id` が確定できず `ambiguous` または `missing` になる。
- `chapter_id` は確定したが `index.json` の候補が複数で一意化できず `ambiguous` になる。
- `chapter_id` は確定したが `index.json` に該当ノードがなく `missing` になる。
- fail 行を残したまま `docs/TASKS.md` の `Node IDs` へ転記しようとする。

### 5.2 互換維持
- 既存 `chapter_id` は維持し、`node_id` 変更時も `chapter_id` を再採番しない。
- `node_id` 削除時は対象行を `deprecated` 化し、代替がある場合のみ `review_note` に後継 `node_id` を明記する。
- `depends_on` 変更のみの場合は `last_verified_at` と `review_note` のみ更新する。

### 5.3 レビュー観点
- `chapter_id` の再採番・削除が発生していないこと。
- `deprecated` 行に廃止理由または後継情報が記載されていること。
- 章別検証サマリ（`design-chapter-validation-spec` 準拠）の `chapter_id` と対応表が一致すること。
- `docs/TASKS.md` の `Node IDs` への転記値が対応表と一致すること。

## 6. 運用チェックリスト
- [ ] Phase 1 抽出時に章対応表の初版を更新した。
- [ ] Phase 2 章ドラフト更新時に章対応表の `display_title` / `node_id` を再確認した。
- [ ] Birdseye 差分（追加/変更/削除）を `review_note` に記録した。
- [ ] `deprecated` 行の扱い（後継・維持期間）をレビュー記録へ反映した。
