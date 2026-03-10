# memx Birdseye Validation Spec

## 1. 目的
`docs/birdseye/index.json` の `nodes[].capsule` と実ファイル群の整合性を検証し、Task Seed へ即転記可能な Issue を生成する。

## 2. 検証対象
- インデックス: `docs/birdseye/index.json`
- 検証キー:
  - `nodes[].node_id`
  - `nodes[].depends_on`
  - `nodes[].capsule`
- 実体確認対象:
  - `nodes[].capsule` で指定されたファイル
  - `depends_on` で参照される `node_id` の存在

## 3. 失敗条件
以下を検出した場合は失敗とする。

1. capsule 未存在
   - `nodes[].capsule` のパスに対応する実ファイルが存在しない。
2. node_id 重複
   - `nodes[].node_id` に重複がある。
3. depends_on 循環
   - `depends_on` グラフに循環参照がある。
4. リンク先欠損
   - `depends_on` が未定義の `node_id` を参照している。

## 4. 出力形式（Task Seed 転記用）
検知結果は 1 issue = 1 レコードで出力し、以下のフィールドを必須とする。

| field | type | 説明 |
| --- | --- | --- |
| `issue_id` | string | 一意ID（例: `BIRDSEYE-VAL-001`） |
| `node_id` | string | 問題対象の node_id（全体問題時は代表 node_id または `global`） |
| `path` | string | 問題対象パス（例: capsule パス、`docs/birdseye/index.json`） |
| `severity` | enum | `high` / `medium` / `low` |
| `temporary_action` | string | 恒久対応までの暫定運用手順 |

### severity の目安
- `high`: 依存解決不能（循環、リンク先欠損）
- `medium`: 参照重複や capsule 欠落で章生成が停止/不安定
- `low`: 運用で回避可能だが追補が必要


## 4.1 共通メタキー追記ルール
`memx_spec_v3/docs/design-evidence-schema-spec.md` 準拠で、Issue レコードと集計サマリの両方に次のキーを必須とする。

- `run_id`
- `generated_at`
- `source_commit`
- `chapter_id`
- `tool`
- `status`
- `severity_summary`
- `evidence_paths`

`chapter_id` は章横断チェックの場合 `global` を使用し、`evidence_paths` は `docs/birdseye/index.json` と検証ログ実体を含める。

## 5. ステータス連携ルール
- 検証で issue を 1 件でも検知した時点で、該当 Task の `Status` を `blocked` に遷移させる。
- 複数 issue がある場合、Task は全 issue 解消まで `blocked` を維持する。
- 再検証で issue が 0 件になった時のみ、通常フロー（`planned` / `active` / `in_progress` / `reviewing`）へ復帰可能とする。

## 6. orchestration との統合ポイント
本仕様は以下の既存チェック項目を統合して運用する。

- `orchestration/memx-design-docs-authoring.md` Phase 1 先頭チェック項目
  - 「`docs/birdseye/index.json` の `nodes[].capsule` 参照先を全件確認し、caps 実体欠落を抽出」
- `orchestration/memx-design-docs-authoring.md` Phase 3 項目
  - 「`docs/birdseye/index.json` の node_id 参照切れを修正」

上記 2 項目は個別運用ではなく、本仕様に定義した検証実行と issue 出力・`blocked` 遷移ルールで一元管理する。
