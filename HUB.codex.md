---
intent_id: hub-codex-contract
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# HUB.codex.md

## 目的
エージェント応答の契約を統一し、実行環境差異（tools 可用性差）に依存しない最小互換出力を保証する。

## 必須出力（固定順）
以下 5 セクションを必須とする（通常時）。

1. `plan`
   - 実施方針を箇条書きで記載。
   - 変更対象・非対象を明記。
2. `patch`
   - 変更ファイルと要点差分を記載。
   - 実変更がない場合は `no-op` と明記。
3. `tests`
   - 実行した検証コマンドと結果（pass/fail/warn）を記載。
   - 未実行の場合は理由を 1 行で記載。
4. `commands`
   - 実行・提案コマンドを列挙。
   - **正本（canonical source）** は `memx_spec_v3/docs/quickstart.md` の API 起動/投入/検索/表示の各コマンド例とする。
   - 参照リンク: [`memx_spec_v3/docs/quickstart.md`](memx_spec_v3/docs/quickstart.md)
5. `notes`
   - 判断理由、制約、未解決事項を最小限で記載。
   - 競合解消がある場合は「双方の意図をどう最小統合したか」を 1 行で記載。

## 失敗時出力契約
- ツール未使用・利用不能・実行基盤制約で処理不能な場合は、
  - `plan` と要求情報（tool 呼び出しの `request envelope` JSON）のみを返す。
  - `patch/tests/commands/notes` は省略可。
- 通常時（ツール実行成功時）は `plan/patch/tests/commands/notes` の 5 セクションを必須とする。
- 禁止事項: ツール結果の推測・捏造。

## 実行環境差異の統一方針
- Native function-calling 可能環境: ツール呼び出しを実行し、同内容の request envelope を併記可。
- 非対応環境: request envelope を一次成果物として返却。
- どの環境でも、最終的な契約解釈は本ファイルを優先する。

- オーケストレーション入力ソースとして `orchestration/*.md` を参照する。

## 自動タスク分割フロー
以下を入力として順序どおり処理し、Task Seed（`TASK.<slug>-MM-DD-YYYY.md`）へ写像する。

1. 工程1: オーケストレーション要件抽出
   - 入力: `orchestration/*.md`
   - 処理: 実施ステップ、制約、役割分担を抽出して実行単位へ分割する。
   - 出力（Task Seed への写像）: Objective 候補と Requirements 候補を生成する。

2. 工程2: 実装指示の正規化
   - 入力: `docs/IN-*.md`
   - 処理: 指示の重複/競合を正規化し、受け入れ条件と検証観点を統一する。
   - 出力（Task Seed への写像）: Requirements を確定し、Commands 候補を生成する。

3. 工程3: 依存関係グラフ解決
   - 入力: `docs/birdseye/index.json`
   - 処理: 参照リンクと依存順を解決し、着手順とブロッカーを特定する。
   - 出力（Task Seed への写像）: Dependencies を確定し、Status 初期値（planned）を設定する。

4. 工程4: 既存 Task Seed との差分統合
   - 入力: `TASK.*-MM-DD-YYYY.md`
   - 処理: 既存タスクとの重複・競合を検出し、追記/新規作成の方針を決定する。
   - 出力（Task Seed への写像）: Objective/Requirements/Commands/Dependencies/Status を差分反映する。

5. 工程5: TASKS 形式への最終落とし込み
   - 入力: 工程1〜4の統合結果
   - 処理: `docs/TASKS.md` の必須項目規約に整形し、語彙・命名・並び順を検証する。
   - 出力（Task Seed への写像）: 最終的に `Objective / Requirements / Commands / Dependencies / Status` を満たす Task Seed を確定する。

6. 工程6: 完了タスクの履歴反映
   - 入力: `Status: done` の Task Seed と `memx_spec_v3/CHANGES.md` の差分
   - 処理: 完了内容を重複排除して `CHANGELOG.md` に最小要約で反映し、必要に応じて `memx_spec_v3/CHANGES.md` に互換性情報を追記する。
   - 出力: `CHANGELOG.md`（正本）更新と Task Seed への `Moved-to-CHANGES: YYYY-MM-DD` 記録。

## 出力例（YAML）
```yaml
task_id: TASK.normalize-hub-yaml-03-03-2026
source:
  - orchestration/plan.md#Phase1
  - orchestration/implementation.md#Phase2
objective: docs/TASKS.md 準拠の Task Seed 生成規約を固定化する
requirements:
  - source は orchestration/...#Phase... の追跡可能形式を使う
  - docs/TASKS.md の必須項目へ欠落なく写像できること
commands:
  - rg "^## " orchestration/*.md
  - python scripts/task_seed_validate.py TASK.normalize-hub-yaml-03-03-2026.md
dependencies:
  - TASK.extract-orchestration-03-03-2026
status: planned
```

### `docs/TASKS.md` 必須項目との対応表

| YAML キー | 転記先（Task Seed） | 転記ルール |
| --- | --- | --- |
| `task_id` | ファイル名 `TASK.<slug>-<MM-DD-YYYY>.md` | `task_id` の `TASK.` 以降をファイル名として使用 |
| `source` | `Requirements` 補足行 | 追跡元として `- source: <path>#Phase...` 形式で列挙 |
| `objective` | `Objective` | 1〜3 行で要約せず原文転記 |
| `requirements` | `Requirements` | 箇条書きで順序維持して転記 |
| `commands` | `Commands` | 実行順を維持して上から転記 |
| `dependencies` | `Dependencies` | 依存なしは `- none` に正規化 |
| `status` | `Status` | `docs/TASKS.md` の許可語彙のみ受理 |

## 言語ポリシー
- デフォルト言語は日本語。
- コード識別子（変数名・関数名・型名・CLI フラグ・JSON キー）は英語を維持する。
- 外部仕様や既存 API 名は原文尊重で改変しない。

## ノード抽出ルール
- ノード抽出単位は `##` 見出しとし、`###` 以下は同一ノード内の補足情報として扱う。
- ノード内の `- [ ]` チェック項目はタスク分解対象とし、各項目を独立した実行タスクとして展開する。
- `[Blocker]` を含む見出し（例: `## [Blocker] ...`）は優先度を最上位に昇格し、通常ノードより先に着手する。
- 依存解決時は、依存先ノードを `active` へ昇格して先行処理し、解決後に依存元ノードを `in_progress` へ復帰させる。
- ステータス遷移の標準フローは `planned → active → in_progress → reviewing → done` とする。
- 例外として、ブロッカー発生時のみ `in_progress → blocked → in_progress` を許可する。
- ステータス語彙の正本は [`docs/TASKS.md`](docs/TASKS.md) とし、差異が出ないように本節と同一語彙で運用する。

## memx側で採用する補完資料一覧

workflow-cookbook の補完資料をそのまま複製せず、memx の運用最小セットとして以下を採用する。

### 採用
- `docs/ADR/README.md`（ADR 運用入口）
- `docs/UPSTREAM.md`（upstream 取り込み方針）
- `docs/UPSTREAM_WEEKLY_LOG.md`（upstream 週次ログ）
- `docs/addenda/A_Glossary.md`（用語統一）
- `docs/addenda/D_Context_Trimming.md`（コンテキスト削減基準）
- `docs/addenda/G_Security_Privacy.md`（セキュリティ/プライバシー基準）
- `datasets/README.md`（データセット台帳）
- `CHANGELOG.md`（完了タスクの利用者向け変更履歴の正本）

### 非採用（workflow-cookbookとの差分）
- workflow-cookbook 側の詳細テンプレート本文・運用例・CI 手順の全文移植は非採用。
- 理由: memx では導線統一を優先し、詳細規定は BLUEPRINT / RUNBOOK / GUARDRAILS / EVALUATION を正本とするため。
