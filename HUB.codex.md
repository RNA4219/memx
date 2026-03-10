---
intent_id: INT-001
owner: memx-resolver
status: active
last_reviewed_at: 2026-03-10
next_review_due: 2026-04-10
---

# HUB.codex.md

`HUB_SCOPE_DECLARATION`: 本ファイルの適用範囲は `memx-resolver/` ツリー。

リポジトリ内の仕様・運用MDを集約し、エージェントがタスクを自動分割できるようにするハブ定義。

## 1. 目的

- リポジトリ配下の計画資料から作業ユニットを抽出し、優先度順に配列
- 生成されたタスクリストを `TASK.*-MM-DD-YYYY` 形式の Task Seed へマッピング

## 2. 入力ファイル分類

- **Blueprint** (`BLUEPRINT.md`): 要件・制約・背景。優先順: 高。
- **Runbook** (`RUNBOOK.md`): 実行手順・コマンド。優先順: 中。
- **Guardrails** (`GUARDRAILS.md`): ガードレール/行動指針。優先順: 高。
- **Evaluation** (`EVALUATION.md`): 受け入れ基準・品質指標。優先順: 中。
- **Checklist** (`CHECKLISTS.md`): リリース/レビュー確認項目。優先順: 低。
- **Birdseye Map** (`docs/birdseye/index.json`): 依存トポロジと役割。優先順: 高。
- **Requirements** (`docs/requirements.md`): 要件定義。優先順: 高。
- **Design** (`docs/design.md`): 設計書。優先順: 高。
- **Interfaces** (`docs/interfaces.md`): インターフェース定義。優先順: 高。

補完資料一覧:

- `README.md`: リポジトリ概要と参照リンク
- `CHANGELOG.md`: 完了タスクと履歴の記録
- `AGENT_GUIDE.md`: AIエージェント向けガイド
- `memx_spec_v3/`: 実装コード

更新日: 2026-03-10

## 3. 自動タスク分割フロー

0. **Birdseye Readiness Check**:
   - `docs/birdseye/index.json` の存在確認と JSON 妥当性検証
   - `nodes` 形式検証、`caps` パス検証
   - 鮮度検証: `generated_at` と対象ドキュメント更新時刻の比較
   - 判定値: `ready | degraded | blocked`

1. **スキャン**: ルート配下を再帰探索し、Markdown front matterを含むファイルを優先取得

2. **Birdseye 専用サブステップ**:
   - 読込順固定: `index.json` → `caps/*.json` → `hot.json`
   - 対象抽出条件: 対象ファイルの `node_id` 起点で ±2 hop を抽出

3. **ノード生成**: 各ファイルから `##` レベルの節をノード化し、`Priority` `Dependencies` を抽出

4. **依存解決**: インポート関係を解析し、該当セクションを子ノードとして連結

5. **粒度調整**: ノード内の ToDo / 箇条書きを単位作業へ分割し、`<= 0.5d` を目安にまとめ直し

6. **テンプレート投影**: 各作業ユニットを `TASK.*-MM-DD-YYYY` 形式の Task Seed へ変換

7. **出力整形**: 優先度、依存、担当の有無でソート

## 4. ノード抽出ルール

- Front matter内の `priority`, `owner`, `deadline` を最優先で採用
- 節タイトルに `[Blocker]` を含む場合は依存解決フェーズで最上位へ昇格
- 箇条書きのうち `[]` or `[ ]` 形式はチェックリスト扱い

### Task Status & Blockers

```yaml
許容ステータス（Allowed）
- `[]` or `[ ]` or `- [ ]`：未着手・未割り振り
- planned：バックログ
- active：受付済/優先キュー入り
- in_progress：着手中
- reviewing：見直し中
- blocked：ブロック中
- done：完了

遷移例（標準）
planned → active → in_progress → reviewing → done
```

## 5. 出力例

```yaml
- task_id: 20260310-01
  source: docs/requirements.md#機能要件
  birdseye:
    node_id: docs/requirements.md
    role: requirements
    hops: 0
    caps_ref: docs/birdseye/caps/docs.requirements.md.json
  objective: 文書登録APIの実装
  scope:
    in: [memx_spec_v3/go/api]
    out: []
  requirements:
    behavior:
      - Markdown文書を登録できる
    constraints:
      - SQLiteを使用
  commands:
    - go test ./api/...
  dependencies: []
```

## 6. 運用メモ

- タスク自動生成ツールはドライランでJSON出力を確認後にIssue化
- 生成後は `CHANGELOG.md` へ反映済みタスクを移す
- Birdseye 鮮度: `generated_at` が最新コミットより古ければ再収集を要求