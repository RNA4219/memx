---
intent_id: INT-001
owner: your-handle
status: active   # draft|active|deprecated
last_reviewed_at: 2025-10-14
next_review_due: 2025-11-14
---

# Agent Tool Policy — Dual Stack

## Runtimes

- Native function-calling tools are registered (OpenAI/Gemini/Vertex).
- No native tools; an external orchestrator parses JSON blocks in text.

## Rules

1. If native tools exist, CALL them (function calling) using the tool names
   below.
2. Always MIRROR each call as a JSON envelope so non-native runtimes can parse:

   ```tool_request
   {"name":"web.search","arguments":{"q":"...", "recency":30}}
   ```

3. Never fabricate tool results. If tools are unavailable, emit `plan` and JSON
   envelopes only.
4. Platform-specific macros remain VERBATIM (do not expand).
5. Default language: Japanese unless code identifiers dictate otherwise.

## Logical Tool Names

- web.search{q, recency?, domains?}
- web.open{url}
- drive.search{query, owner?, modified_after?}
- gmail.search{query, max_results?}
- calendar.search{time_min?, time_max?, query?}

## Output Contract

`plan`/`patch`/`tests`/`commands`/`notes`

- `plan` は各タスクに `node_id` / `role` / `source_caps` を必須で埋め込む。
- `source_caps` は `docs/birdseye/caps/*.json` のファイル名（例: `caps/core.json`）を保持する。
- JSON/YAML の最小フィールドセットは下記を満たす。

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

## HUB.codex.md

`HUB_SCOPE_DECLARATION`: 本ファイルの適用範囲は `workflow-cookbook/` ツリー。

リポジトリ内の仕様・運用MDを集約し、エージェントがタスクを自動分割できるようにするハブ定義。
`BLUEPRINT.md` など既存ファイルに加えて、オーケストレーション専用のMD（例: `orchestration/*.md`）も取り込む。

## 1. 目的

- リポジトリ配下の計画資料から作業ユニットを抽出し、優先度順に配列
- オーケストレーションMD（ワークフロー全体の段取り記載）を検出し、必要な子タスクへ展開
- 生成されたタスクリストを `TASK.*-MM-DD-YYYY` 形式の Task Seed へマッピング

## 2. 入力ファイル分類

- **Blueprint** (`BLUEPRINT.md`): 要件・制約・背景。優先順: 高。
  備考: 最上位方針。
- **Runbook** (`RUNBOOK.md`): 実行手順・コマンド。優先順: 中。
  備考: 具体的操作。
- **Guardrails** (`GUARDRAILS.md`): ガードレール/行動指針。優先順: 高。
  備考: 全メンバー必読。
- **Incident Logs** (`docs/IN-*.md`): インシデント記録（影響・再発防止など）。優先順: 高。
  備考: 再発防止策とフォローアップ抽出。
- **Evaluation** (`EVALUATION.md`): 受け入れ基準・品質指標。優先順: 中。
  備考: 検収条件。
- **Checklist** (`CHECKLISTS.md`): リリース/レビュー確認項目。優先順: 低。
  備考: 後工程。
- **Orchestration** (`orchestration/*.md`): ワークフロー構成・依存関係。優先順: 可変。
  備考: 最優先のブロッカーを提示。
- **Birdseye Map** (`docs/birdseye/index.json` など): 依存トポロジと役割を把握。優先順: 高。
  備考: `plan` 出力にノードID/役割を埋め込む基準面。
- **Task Seeds** (`TASK.*-MM-DD-YYYY`): 既存タスクドラフト。優先順: 高。
  備考: 未着手タスクの候補。

補完資料一覧:

- `README.md`: リポジトリ概要と参照リンク
- `CHANGELOG.md`: 完了タスクと履歴の記録
- `.github/PULL_REQUEST_TEMPLATE.md`: PR 作成時のチェック項目（Intent/リスク/Canary連携）
- `.github/ISSUE_TEMPLATE/bug.yml`: Intent ID と自動ゲート確認を必須化した不具合報告フォーム
- `governance/policy.yaml`: QA が管理する自己改変境界・カナリア中止条件・SLO
- `governance/prioritization.yaml`: 設計更新の優先度スコア計算ルール
- `docs/IN-*.md`: インシデントログ本体。Blueprint/Evaluation との相互リンクを維持し、再発防止策の同期を確認
- `docs/INCIDENT_TEMPLATE.md`: 検知/影響/5Whys/再発防止/タイムラインのインシデント雛形
- `docs/TASKS.md`: Task Seed 運用ガイドとテンプレートの要点
- `docs/ADR/README.md`: 最新 ADR 索引と更新手順
- `CODE_OF_CONDUCT.md`: コントリビューター行動規範と `maintainers@workflow-cookbook.example` / `security@workflow-cookbook.example` の連絡窓口
- `SECURITY.md`: 脆弱性報告窓口と連絡手順
- `CODEOWNERS`: `/governance/**` とインシデント雛形を QA 管轄とする宣言
- `LICENSE`: OSS としての配布条件（Apache-2.0）
- `.github/release-drafter.yml`: リリースノート自動整形のテンプレート
- `.github/workflows/release-drafter.yml`: Release Drafter の CI 設定
- `docs/UPSTREAM.md`: Workflow Cookbook 派生リポからの知見取り込み手順と評価基準
- `docs/UPSTREAM_WEEKLY_LOG.md`: Upstream 差分確認の週次ログテンプレート
- `docs/addenda/A_Glossary.md`: 用語定義を参照するための補足資料
- `docs/addenda/D_Context_Trimming.md`: コンテキストトリミング指標・検証フローの詳細ガイド
- `docs/addenda/G_Security_Privacy.md`: SAC 原則に準拠したキー管理・ログマスキング等の運用ディテールを参照
- `datasets/README.md`: データセット取得履歴とハッシュを管理。データ保持レビュー時は本表で収集状況を確認

更新日: 2025-10-24

## 3. 自動タスク分割フロー

0. **Birdseye Readiness Check**:
   - **`index.json` 検証**: `docs/birdseye/index.json` の存在確認と JSON 妥当性検証を行う。
   - **`nodes` 形式検証**: `index.json.nodes` が map/object であることを確認する。
   - **`caps` パス検証**: `index.json.nodes` の各要素に対して `caps` パスの存在確認を行う。
   - **互換注意**: 将来 `nodes` が配列形式へ変更される可能性を考慮し、object/array の両対応判定を実装する。
   - **鮮度検証**: `index.json.generated_at` と対象ドキュメント更新時刻の鮮度条件を確認し、判定基準は [`GUARDRAILS.md` の「鮮度管理（Staleness Handling）」](GUARDRAILS.md#鮮度管理staleness-handling) を正とする。
   - **判定値**: 結果は `ready | degraded | blocked` の3値で扱う。`degraded` は既知ノード限定の分割継続、`blocked` は新規分割停止を意味する。失敗時（`degraded` / `blocked`）は `notes.readiness_status` へ判定結果を必ず記録する。
1. **スキャン**: ルートと `orchestration/` 配下を再帰探索し、Markdown front matter
   (`---`) を含むファイルを優先取得。
   
### Birdseye Bootstrap

1. `workflow-cookbook/docs/BIRDSEYE.md` を参照して用語と成果物を理解
2. `workflow-cookbook/docs/birdseye/index.json` を読み込みノード存在を確認
3. 必要な `workflow-cookbook/docs/birdseye/caps/*.json` を最小読込
4. 不足時のみ `workflow-cookbook/tools/codemap/update.py` へ遷移

この順序を満たさない場合はタスク化を開始しない。

2. **Birdseye 専用サブステップ**:
   - **読込順固定**: Birdseye JSON を第一読者として、`docs/birdseye/index.json` → `docs/birdseye/caps/*.json` → `docs/birdseye/hot.json` の順で必ず読み込む。
   - **対象抽出条件**: 対象ファイルの `node_id` 起点で ±2 hop を抽出し、未解決ノードは `hot.json` の hot list で補完する。
   - **埋込必須項目**: 各候補タスクに `node_id` / `role` / `source_caps` を付与し、GUARDRAILS の `plan` 出力要件（ノードID明示）を初期段階で満たす。
3. **ノード生成**: 各ファイルから `##` レベルの節をノード化し、`Priority`
   `Dependencies` などのキーワードを抽出。
4. **依存解決**: Orchestrationノードに含まれる依存パスを解析し、該当セクションを子ノードとして連結。
5. **インシデント抽出**: `docs/IN-*.md` のインシデントセクションを走査。
   再発防止やテスト強化の箇条書きを Task Seed 候補としてタグ付け。
6. **粒度調整**: ノード内の ToDo / 箇条書きを単位作業へ分割し、`<= 0.5d`
   を目安にまとめ直し。
7. **テンプレート投影**: 各作業ユニットを `TASK.*-MM-DD-YYYY` 形式の Task Seed
   (`Objective` `Requirements` `Commands`) へ変換し、欠損フィールドは元資料の該当行を引用。
8. **出力整形**: 優先度、依存、担当の有無でソートし、GitHub Issue もしくは
   PR下書きとしてJSON/YAMLに整形。
9. **タスク化**: タスクは独立性が保てる粒度まで分割し、責務の重複(コンフリクト)を避ける。
　 変更は小さく・短時間で終わるブランチとして切り、早めのrebaseで常に最新に追従する。
   リスクがある、タスクが重なっている場合は**Task Seeds** (`TASK.*-MM-DD-YYYY`)に記載を行うこと。

## 4. ノード抽出ルール

- Front matter内の `priority`, `owner`, `deadline` を最優先で採用
- 節タイトルに `[Blocker]` を含む場合は依存解決フェーズで最上位へ昇格
- 箇条書きのうち `[]` or `[ ]` 形式はチェックリスト扱い、`- [ ]` はタスク分解対象。詳細ステータスは後述`Task Status & Blockers`参照
- コードブロックはコマンドサンプルとして `Commands` セクションに集約

- **Task Status & Blockers**

```yaml
許容ステータス（Allowed）
- `[]` or `[ ]` or `- [ ]`：未着手・未割り振り
- planned：バックログ。着手順待ち
- active：受付済/優先キュー入り（担当/期日が付いた状態）
- in_progress：着手中
- reviewing：見直し中（レビュー/ふりかえり/承認待ち）
- blocked：ブロック中（外的依存で進められない）
- done：完了

遷移例（標準）
planned → active → in_progress → reviewing → done
ブロック例（例外）
in_progress → blocked → in_progress（解除後に戻す）
```

## 5. 出力例（擬似）

```yaml
- task_id: 20240401-01
  source: orchestration/api-rollout.md#Phase1
  birdseye:
    node_id: node-orchestration-api-rollout-phase1
    role: orchestration
    hops: 0
    caps_ref: docs/birdseye/caps/orchestration.api-rollout.json
  objective: API Gateway ルーティング切替の段階実行
  scope:
    in: [infra/aws/apigw]
    out: [legacy/cli]
  requirements:
    behavior:
      - Blue/Green 切替時にダウンタイム0
    constraints:
      - 既存API破壊禁止
  commands:
    - terraform plan -target=module.api_gateway
  dependencies:
    - 20240331-ops-01
```

受け入れ条件（最低限）:

- Task Seed に記録された `birdseye.node_id` から、`docs/birdseye/index.json` の元ノードを一意に逆引きできること。

## 6. 運用メモ

- Orchestration MD には `## Phase` `## Stage` 等の段階名を揃える
- タスク自動生成ツールはドライランでJSON出力を確認後にIssue化
- 生成後は `CHANGELOG.md` へ反映済みタスクを移すことで履歴が追える
- Birdseye 鮮度: `docs/birdseye/index.json.generated_at` が最新コミットより古ければ再収集を要求。
  該当 Capsule も同時更新。
- 鮮度判断と人間エスカレーション条件は [`GUARDRAILS.md` の「鮮度管理（Staleness Handling）」](GUARDRAILS.md#鮮度管理staleness-handling) を正とし、
  `index.generated_at` 逆転 / Caps 不在 / 対象ノード未登録 / `codemap.update` 未実装時は同節の依頼フローへ遷移する。
- `codemap.update` は Birdseye 再生成時のみ実行。
  Dual Stack では関数呼び出し→`tool_request` ミラーを同一内容で送る。

### Birdseyeアクセス異常時ハンドリング

- **ケースA（`index.json` 不在/破損）**
  - 判定: `docs/birdseye/index.json` が読めない、または JSON パース失敗。
  - Readiness: `blocked`。
  - 処理: **タスク生成を停止**し、新規分割を行わない。
  - ステータス: 全候補を `blocked`。
  - 次アクション: `codemap.update` による `index+caps` 再生成を人間へ依頼（`GUARDRAILS.md` の鮮度管理フローに準拠）。
- **ケースB（`caps` 部分欠損）**
  - 判定: `index.json` は読めるが、参照先 `docs/birdseye/caps/*.json` の一部が欠損/破損。
  - Readiness: `degraded`。
  - 処理: **既知ノードのみ限定タスク化**し、欠損ノード依存のタスクは生成しても `blocked` に固定。
  - ステータス: 既知ノード=`planned`/`active` 可、欠損ノード=`blocked`。
  - 次アクション: 欠損 Capsule のみを対象に `codemap.update`（`emit:"caps"`）再生成を依頼。
- **ケースC（`hot.json` のみ欠損）**
  - 判定: `index.json` と必要 `caps` は利用可能で、`docs/birdseye/hot.json` のみ欠損/破損。
  - Readiness: `ready`（警告付き）。
  - 処理: **`index`/`caps` ベースで継続処理**し、ホットリスト最適化のみ無効化。
  - ステータス: 通常どおり（警告付き）でタスク生成を継続。
  - 次アクション: 警告を `notes` に残し、次回 Birdseye 再生成時に `hot.json` を復旧。

#### `notes` 必須記録項目（Birdseye欠損時）

- `readiness_status`: `ready` / `degraded` / `blocked` のいずれか（必須）。
- `missing_files`: 欠損/破損ファイルの相対パス（例: `docs/birdseye/index.json`）。
- `impacted_node_ids`: 影響を受けたノードID一覧（不明な場合は `unknown` を明記）。
- `provisional_decision`: 暫定判断（停止 / 限定継続 / 警告付き継続）と理由。
- `regen_request_to`: 再生成依頼先（担当チーム/ロール名、または `human-operator`）。

#### 運用テンプレート（Birdseyeアクセス異常時の再生成依頼）

> 復旧手順の詳細は [`tools/codemap/README.md` の「Birdseyeアクセス異常時の復旧手順」](tools/codemap/README.md#birdseyeアクセス異常時の復旧手順) を参照。

```yaml
birdseye_regen_request:
  requested_to_role: "<依頼先ロール>"
  execute_emit: "index|caps|hot"
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
