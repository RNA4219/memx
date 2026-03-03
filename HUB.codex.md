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

## 言語ポリシー
- デフォルト言語は日本語。
- コード識別子（変数名・関数名・型名・CLI フラグ・JSON キー）は英語を維持する。
- 外部仕様や既存 API 名は原文尊重で改変しない。

## 出力例（YAML）
タスク化時は、追跡可能な最小単位として次の YAML 形式を採用する。

```yaml
task_id: TASK.sync-hub-yaml-03-03-2026
source: orchestration/roadmap.md#Phase2
objective: HUB 出力契約に YAML 例と転記規約を追加する
requirements:
  - task_id/source/objective/requirements/commands/dependencies/status を必須化する
  - source は orchestration/...#Phase... 形式で記載する
commands:
  - rg "出力例（YAML）" HUB.codex.md
  - rg "対応表" HUB.codex.md
dependencies:
  - none
status: in_progress
```

- `source` は `orchestration/<file>.md#Phase<N>` 形式を正とし、作業起点を一意に追跡可能にする。

## `docs/TASKS.md` 必須項目との対応表

| YAML キー | 転記先（`docs/TASKS.md`） | 固定ルール |
| --- | --- | --- |
| `task_id` | 1. 命名規則（Task Seed ファイル名） | `TASK.<slug>-<MM-DD-YYYY>.md` を生成し、識別子として維持 |
| `source` | Dependencies | `orchestration/...#Phase...` を依存・起点情報として先頭に記載 |
| `objective` | Objective | 1〜3 行でそのまま転記 |
| `requirements` | Requirements | 箇条書きで順序を維持して転記 |
| `commands` | Commands | 実行順を維持して転記 |
| `dependencies` | Dependencies | `- none` を含め原文維持で転記 |
| `status` | Status | `planned/active/in_progress/reviewing/blocked/done` のみ許可 |

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

### 非採用（workflow-cookbookとの差分）
- workflow-cookbook 側の詳細テンプレート本文・運用例・CI 手順の全文移植は非採用。
- 理由: memx では導線統一を優先し、詳細規定は BLUEPRINT / RUNBOOK / GUARDRAILS / EVALUATION を正本とするため。
