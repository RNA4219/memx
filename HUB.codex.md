---
intent_id: hub-codex-contract
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# HUB.codex.md

`HUB_SCOPE_DECLARATION`: 本ファイルの適用範囲はリポジトリルート配下全体。ただし `workflow-cookbook/` 配下の作業では `workflow-cookbook/HUB.codex.md` を優先する。

## 目的
エージェント応答の契約を統一し、実行環境差異（tools 可用性差）に依存しない最小互換出力を保証する。

## 必須出力（固定順）
以下 5 セクションを必須とする（通常時）。

1. `plan`
   - 実施方針を箇条書きで記載。
   - 変更対象・非対象を明記。
   - 各タスクに `node_id` / `role` / `source_caps` を必須で含める。
   - `source_caps` は `docs/birdseye/caps/*.json` の参照パスを保持する。
2. `patch`
   - 変更ファイルと要点差分を記載。
   - 実変更がない場合は `no-op` と明記。
3. `tests`
   - 実行した検証コマンドと結果（pass/fail/warn）を記載。
   - 未実行の場合は理由を 1 行で記載。
4. `commands`
   - 実行・提案コマンドを列挙。
   - **正本（canonical source）** は `memx_spec_v3/docs/quickstart.md` とし、コマンド表記は**リポジトリルート起点**の `go run ./memx_spec_v3/go/cmd/mem ...` に統一する。
   - commands セクションには最低限以下を短く再掲する。
     - ingest: `go run ./memx_spec_v3/go/cmd/mem in short --title "..." --stdin --api-url http://127.0.0.1:7766`
     - search: `go run ./memx_spec_v3/go/cmd/mem out search "..." --api-url http://127.0.0.1:7766`
     - show: `go run ./memx_spec_v3/go/cmd/mem out show <NOTE_ID> --api-url http://127.0.0.1:7766`
   - `cd memx_spec_v3/go` + `go run ./cmd/mem ...` は代替表記としてのみ許容し、正本扱いしない。
   - 差異チェック手順（更新時は必須）:
     1. `rg -n "go run ./memx_spec_v3/go/cmd/mem|go run ./cmd/mem|cd memx_spec_v3/go" memx_spec_v3/docs/quickstart.md RUNBOOK.md HUB.codex.md`
     2. 正本が quickstart のみであること、RUNBOOK/HUB が同一規約で追従していることを確認する。
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

## Birdseye Bootstrap
依存関係解決へ進む前に、以下を**この順序で**参照する。

1. `docs/BIRDSEYE.md`（用語/成果物理解）
2. `docs/birdseye/index.json`
3. 必要最小限の `docs/birdseye/caps/*.json`
4. `docs/birdseye/hot.json`（`optional`。未運用の間は欠損を許容）
5. 不足時のみ `tools/codemap/update.py`

`index.json` / `caps` が揃っていればタスク化を開始できる。`hot.json` が欠損していても継続実行可とし、`notes.readiness_status=ready` を維持したまま `notes.missing_files` に `docs/birdseye/hot.json` を記録する。
この順序（`optional` を含む）を満たさない場合はタスク化を開始しない。
Birdseye 復旧が必要な場合の遷移先は [`RUNBOOK.md` の「Birdseye 鮮度不足時の復旧手順」](RUNBOOK.md#birdseye-鮮度不足時の復旧手順) のみとする。

## 自動タスク分割フロー
以下を入力として順序どおり処理し、Task Seed（`TASK.<slug>-MM-DD-YYYY.md`）へ写像する。

1. 工程1: オーケストレーション要件抽出
   - 入力: `orchestration/*.md`
   - 処理: 実施ステップ、制約、役割分担を抽出して実行単位へ分割する。
   - 出力（Task Seed への写像）: Objective 候補と Requirements 候補を生成する。

2. 工程2: 実装指示の正規化
   - 入力: `docs/IN-*.md`
   - 運用ルール: `docs/IN-BASELINE.md` は補助資料として扱い、`docs/IN-<実日付>-<連番>.md` の実インシデント記録を優先入力する。
   - 処理: 指示の重複/競合を正規化し、受け入れ条件と検証観点を統一する。
   - 出力（Task Seed への写像）: Requirements を確定し、Commands 候補を生成する。

3. 工程3: 依存関係グラフ解決
   - 前提条件: 「Birdseye Bootstrap」の参照順完了。
   - 入力: `docs/birdseye/index.json`（必要最小限の `docs/birdseye/caps/*.json` を含む。`docs/birdseye/hot.json` は `optional`）
   - 処理: 参照リンクと依存順を解決し、着手順とブロッカーを特定する。
   - 出力（Task Seed への写像）: Dependencies を確定し、Status 初期値（planned）を設定する。

4. 工程4: 既存 Task Seed との差分統合
   - 入力: `TASK.*-MM-DD-YYYY.md`
   - 処理: 既存タスクとの重複・競合を検出し、追記/新規作成の方針を決定する。
   - レビュー規則: `Source` 未記載の Task Seed は差し戻し（reviewing 継続）とし、出典追記後に再レビューする。
   - 出力（Task Seed への写像）: Source/Objective/Requirements/Commands/Dependencies/Status を差分反映する。

5. 工程5: TASKS 形式への最終落とし込み
   - 入力: 工程1〜4の統合結果
   - 処理: `docs/TASKS.md` の必須項目規約に整形し、語彙・命名・並び順を検証する。
   - レビュー規則: `Source` が欠落した Task Seed は承認せず、工程4へ差し戻す。
   - 出力（Task Seed への写像）: 最終的に `Source / Objective / Requirements / Commands / Dependencies / Status` を満たす Task Seed を確定する。

6. 工程6: 完了タスクの履歴反映
   - 入力: `Status: done` の Task Seed と `memx_spec_v3/CHANGES.md` の差分
   - 処理: 完了内容を重複排除して `CHANGELOG.md` に最小要約で反映し、必要に応じて `memx_spec_v3/CHANGES.md` に互換性情報を追記する。
   - 出力: `CHANGELOG.md`（正本）更新と Task Seed への `Moved-to-CHANGES: YYYY-MM-DD` 記録。

## 出力例（YAML）
### 最小 YAML スキーマ（Output Contract 準拠）
```yaml
plan:
  - task_id: string
    source: string
    objective: string
    node_id: string
    role: string
    source_caps:
      - string
```

### Task Seed 例
```yaml
task_id: TASK.normalize-hub-yaml-03-03-2026
source:
  - orchestration/plan.md#Phase1
  - orchestration/implementation.md#Phase2
birdseye:
  node_id: node-orchestration-plan-phase1
  role: orchestrator
  caps_ref: docs/birdseye/caps/requirements.json
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
| `source` | `Source` | 追跡元として `- <path>#Phase...` 形式で列挙 |
| `objective` | `Objective` | 1〜3 行で要約せず原文転記 |
| `requirements` | `Requirements` | 箇条書きで順序維持して転記 |
| `commands` | `Commands` | 実行順を維持して上から転記 |
| `dependencies` | `Dependencies` | 依存なしは `- none` に正規化 |
| `status` | `Status` | `docs/TASKS.md` の許可語彙のみ受理 |

## 言語ポリシー
- デフォルト言語は日本語。
- コード識別子（変数名・関数名・型名・CLI フラグ・JSON キー）は英語を維持する。
- 外部仕様や既存 API 名は原文尊重で改変しない。

## 運用メモ（Birdseye 鮮度判定）
- `docs/birdseye/index.json.generated_at` は UTC RFC3339 であること。
- 判定時刻から **7日以内** なら鮮度 OK、**7日超** は鮮度不足として扱う。
- 鮮度不足時は [`RUNBOOK.md` の「Birdseye 鮮度不足時の復旧手順」](RUNBOOK.md#birdseye-鮮度不足時の復旧手順) を実施する（遷移先は本リンクに統一）。

## Birdseye Access Preflight（開始前チェック）
Birdseye 入力の健全性を事前判定し、`ready | degraded | blocked` の 3 値で `notes.readiness_status` を必ず記録する。

### Birdseye Readiness Preflight

本節は **判定のみ** を扱う。復旧コマンド列は [`RUNBOOK.md` の「Birdseye 鮮度不足時の復旧手順」](RUNBOOK.md#birdseye-鮮度不足時の復旧手順) を正本とし、HUB には重複掲載しない。

#### 最低限の実行コマンド例（判定用）
1. 存在確認
   ```bash
   test -f docs/birdseye/index.json
   ```
2. JSON パース確認
   ```bash
   python -m json.tool docs/birdseye/index.json >/dev/null
   ```
3. `caps` 存在確認
   ```bash
   python - <<'PY'
   import json
   from pathlib import Path

   idx = json.loads(Path('docs/birdseye/index.json').read_text())
   nodes = idx.get('nodes', {})
   items = nodes.values() if isinstance(nodes, dict) else nodes
   missing = []
   for n in items:
       cap = n.get('caps')
       if cap and not Path(cap).exists():
           missing.append((n.get('node_id', 'unknown'), cap))
   print('missing_caps=', missing)
   raise SystemExit(1 if missing else 0)
   PY
   ```
4. `generated_at` 鮮度確認（7日以内）
   ```bash
   python - <<'PY'
   import datetime as dt
   import json
   from pathlib import Path

   idx = json.loads(Path('docs/birdseye/index.json').read_text())
   gen = dt.datetime.fromisoformat(idx['generated_at'].replace('Z', '+00:00'))
   age_days = (dt.datetime.now(dt.timezone.utc) - gen).total_seconds() / 86400
   print(f'age_days={age_days:.2f}')
   raise SystemExit(1 if age_days > 7 else 0)
   PY
   ```

#### チェック項目（固定順）
1. `docs/birdseye/index.json` 存在確認。  
   判定条件: 存在なら `ready`、欠損なら `blocked`。
2. JSON パース可否（`docs/birdseye/index.json`）。  
   判定条件: パース成功なら `ready`、失敗なら `blocked`。
3. `index.nodes` 形式妥当性（`object` / `array` の互換入力を許容）。  
   判定条件: `object|array` のいずれかなら `ready`、それ以外は `blocked`。
4. 各ノードの `caps` 参照先ファイルの存在。  
   判定条件: 全件存在なら `ready`、一部欠損なら `degraded`、全件欠損または走査不能なら `blocked`。
5. `generated_at` の7日以内判定（本ファイル「運用メモ（Birdseye 鮮度判定）」に従う）。  
   判定条件: 7日以内なら `ready`、7日超なら `degraded`、値欠損/不正なら `blocked`。
6. `nodes[].capsule` 実体存在確認。  
   判定条件: 全件存在なら `ready`、一部欠損なら `degraded`、全件欠損または走査不能なら `blocked`。

---

### ケース別ハンドリング

- **ケースA: `index.json` 欠損 / 破損**
  - 判定: `blocked`
  - 処置: 新規タスク生成を停止し、復旧完了まで Task Seed の `Status` を `planned/active` に進めない。

- **ケースB: `caps` 部分欠損**
  - 判定: `degraded`
  - 処置:
    - 既知ノードのみ処理継続
    - 欠損 `caps` 依存ノードは `Status: blocked` として分離管理
  - 復旧参照先: [`RUNBOOK.md` の「Birdseye 鮮度不足時の復旧手順」](RUNBOOK.md#birdseye-鮮度不足時の復旧手順)

- **ケースC: `hot.json` 欠損**
  - 判定: `ready`（警告付き）
  - 処置: 継続実行可。`notes.missing_files` に欠損を記録する。

---

### 欠損時 `notes` 必須キー

欠損または劣化を検知した場合、`notes` に以下キーを必須で記載する。

- `readiness_status`
- `missing_files`
- `impacted_node_ids`
- `provisional_decision`
- `regen_request_to`

#### `notes` 記載例（実コマンド結果との対応）

| コマンド結果（例） | `notes` に記載するキー/値（例） |
| --- | --- |
| `test -f docs/birdseye/index.json` が非ゼロ終了 | `readiness_status: blocked` / `missing_files: [docs/birdseye/index.json]` |
| `python -m json.tool docs/birdseye/index.json` が失敗 | `readiness_status: blocked` / `provisional_decision: stop task seeding and request regen via RUNBOOK` |
| `missing_caps=[('requirements','docs/birdseye/caps/requirements.json')]` | `readiness_status: degraded` / `missing_files: [docs/birdseye/caps/requirements.json]` / `impacted_node_ids: [requirements]` |
| `age_days=8.20`（7日超） | `readiness_status: degraded` / `provisional_decision: continue with unaffected nodes and run RUNBOOK recovery` |
| 復旧担当未アサイン | `regen_request_to: unassigned` |

```yaml
notes:
  readiness_status: degraded
  missing_files:
    - docs/birdseye/caps/requirements.json
  impacted_node_ids:
    - requirements
  provisional_decision: continue with unaffected nodes and run RUNBOOK.md#birdseye-鮮度不足時の復旧手順
  regen_request_to: unassigned
```

- `ready | degraded | blocked` は Birdseye 入力健全性専用語彙として `notes.readiness_status` に限定して使用する。
- Task Seed の `Status` は `docs/TASKS.md` 許可語彙（`planned/active/in_progress/reviewing/blocked/done`）のみを使用する。
- 語彙衝突回避のため、`degraded` は `notes.readiness_status` のみに限定し、Task Seed `Status` に流用しない。

## Birdseyeアクセス異常時の再生成依頼テンプレート

復旧手順本文は `RUNBOOK.md` の「Birdseye 鮮度不足時の復旧手順」と
`./workflow-cookbook/tools/codemap/README.md`（別名パス `tools/codemap/README.md` は注記扱い）の「Birdseyeアクセス異常時の復旧手順」を正本とし、
実行コマンドの canonical path は `workflow-cookbook/tools/codemap/update.py` に固定する。本節は依頼フォーマットのみを定義する。

```yaml
birdseye_regen_request:
  requested_to_role: "<依頼先ロール>"
  execute_emit: "index|caps|hot"
  recovery_runbook: "RUNBOOK.md#birdseye-鮮度不足時の復旧手順"
  execute_command_alias_note: "`tools/codemap/update.py` は別名パス（注記扱い）"
  post_checks:
    generated_files:
      - "docs/birdseye/index.json"
      - "docs/birdseye/hot.json"
      - "docs/birdseye/caps/*.json"
    missing_resolved: true
    generated_at_updated: true
  re_escalation_if:
    - "実行コマンドが非ゼロ終了した"
    - "必要成果物が生成されない、またはJSONパースに失敗した"
    - "欠損が解消しない、または generated_at が更新されない"
```

## ノード抽出ルール
- ノード抽出単位は `##` 見出しとし、`###` 以下は同一ノード内の補足情報として扱う。
- ノード内の `- [ ]` チェック項目はタスク分解対象とし、各項目を独立した実行タスクとして展開する。
- `[Blocker]` を含む見出し（例: `## [Blocker] ...`）は優先度を最上位に昇格し、通常ノードより先に着手する。
- 依存解決時は、依存先ノードを `active` へ昇格して先行処理し、解決後に依存元ノードを `in_progress` へ復帰させる。
- ステータス遷移の標準フローは `planned → active → in_progress → reviewing → done` とする。
- 例外として、ブロッカー発生時のみ `in_progress → blocked → in_progress` を許可する。
- ステータス語彙の正本は [`docs/TASKS.md`](docs/TASKS.md) とし、差異が出ないように本節と同一語彙で運用する。
- front matter が未設定の Task Seed は `priority=medium` / `owner=unassigned` / `deadline=tbd` / `status=planned` をフォールバック値として補完する。

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
