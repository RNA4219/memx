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

## 5. `docs/birdseye/index.json` 更新時の追随ルール

### 5.1 差分検知
- 更新時は旧版/新版の `node_id` と `depends_on` を比較し、以下を区分する。
  - 追加: 新規 `node_id`
  - 変更: 既存 `node_id` の `depends_on` 変更
  - 削除: 旧版にのみ存在する `node_id`
- 差分検知結果は対応表の `review_note` に要約を残す。

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
