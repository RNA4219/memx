# TASKS 運用ガイド（memx）

## 1. 命名規則
- Task Seed は **`TASK.<slug>-<MM-DD-YYYY>.md`** 形式で作成する。
- 日付は `date +%m-%d-%Y` の結果をそのまま使う（例: `03-03-2026`）。
- `slug` は英小文字・数字・ハイフンのみを使用する。


## 1-1. インシデントIDプレースホルダ移行手順（`docs/IN-<DATE8>-001.md` / `docs/IN-<DATE6xx>-001.md` 対応）
- 対象: Task Seed の `Source` / `Requirements` / 関連文書で、`IN-<DATE8>-001` や `IN-<DATE6xx>-001` のようなプレースホルダIDを参照している記述。
- 運用原則: 実運用の要件根拠として許可するIDは `docs/IN-<実日付>-<連番>.md` のみ。

### 移行ステップ
1. 棚卸し（dry-run）
   - `rg -n "IN-(Y{4}M{2}D{2}|[0-9]{6}xx)-[0-9]{3}" docs TASK.*.md` でプレースホルダ参照を列挙する。
2. 実IDへの置換先決定
   - 各プレースホルダ参照に対し、同一事象の実日付ID（`docs/IN-<YYYYMMDD>-001.md` 形式）を 1:1 で割り当てる。
   - 実日付ID文書が未作成の場合は先に `docs/IN-<実日付>-<連番>.md` を新規起票してから置換する。
3. 参照更新
   - Task Seed の `Source` / `Requirements` からプレースホルダIDを除去し、実日付IDへ更新する。
   - `Source` に `TBD` / テンプレートID が残存しないことを確認する。
4. テンプレート文書の扱い固定
   - `docs/IN-<DATE8>-001.md` / `docs/IN-<DATE6xx>-001.md` はテンプレート用途に限定し、実績証跡として参照禁止の注記を維持する。
5. 完了判定
   - `rg -n "IN-(Y{4}M{2}D{2}|[0-9]{6}xx)-[0-9]{3}|Source:.*TB[D]" docs TASK.*.md` の結果が 0 件であること。

## 推奨 front matter キー
- Task Seed 冒頭に YAML front matter を付与し、以下キーの記載を推奨する。
  - `priority`: `high` / `medium` / `low` のいずれか。
  - `owner`: 担当チームまたは担当者（例: `memx-core`）。
  - `deadline`: 期日を `YYYY-MM-DD` 形式で記載する。

### priority 記載ガイド（設計書作成タスク）
- `docs/design-docs-prioritization-spec.md` を正本として、次の4軸で判定する。
  - Blocker有無
  - REQ網羅率への影響
  - 契約差分 high 件数
  - Birdseye issue の有無
- 判定ルール:
  - `high`: 4軸のいずれかが high。
  - `medium`: high なし、かつ 1軸以上が medium。
  - `low`: 4軸すべて low。
- Task Seed には `Requirements` または `Dependencies` に判定根拠を1行で残す。
- 章別再評価（Phase完了時）で軸が変化した場合は `priority` を更新する。

## 2. Task Seed 必須項目
各 Task Seed には次の6項目を必須で含める。
- Incident 由来要件の転記ルールは `memx_spec_v3/docs/incident-to-task-traceability-spec.md` を参照する。

### Source
- 要件の出典を `path#Section` 形式で記載する。
- 例: `orchestration/memx-v1-bootstrap.md#Phase 2`
- 複数ある場合は箇条書きで列挙する。
- `orchestration/memx-design-docs-authoring.md` 由来の入力参照名を使う場合は、`memx_spec_v3/docs/design-reference-resolution-spec.md` の正規パスマッピングで解決した値のみを許可する。
- `HUB.codex.md` 工程2の運用ルールに従い、要件根拠として許可するインシデント記録は **`docs/IN-<実日付>-<連番>.md` のみ** とする（`docs/IN-BASELINE.md` は補助資料、`docs/IN-<YYYYMMDD>-001.md` などのテンプレートIDは不可）。
- `Source` にテンプレートID（例: `IN-<YYYYMMDD>-001`）または `TBD` を含む Task Seed は、`reviewing` を継続して差し戻す。

### Node IDs
- `docs/birdseye/index.json` の `node_id` を記載する。
- 依存グラフ（`depends_on`）と照合する対象タスクでは必須、照合対象外タスクでは任意とする。
- 記載形式は箇条書きで、依存元/依存先の対応がわかるように補足行を付ける。

### Objective
- タスクの目的を1〜3行で記載する。

### Requirements
- 満たすべき要件を箇条書きで記載する。
- 互換性や非機能制約がある場合はここに明記する。
- エラーコードを新規追加・変更するタスクでは、`memx_spec_v3/docs/requirements.md` の「6-4. エラーモデル」と `memx_spec_v3/docs/error-contract.md` を同時更新対象に含める。
- API/CLI の契約（request/response/error/`--json`）を変更するタスクでは、`memx_spec_v3/docs/contracts/openapi.yaml` と `memx_spec_v3/docs/contracts/cli-json.schema.json` の更新を必須とする。
- Phase 2〜4 対象タスクでは、`memx_spec_v3/docs/design-chapter-validation-spec.md` に準拠した章別検証サマリ作成（`chapter_id` / `req_coverage` / `contract_alignment_high_count` / `link_broken_count` / `birdseye_issue_count` / `evidence_paths`）を必須要件として記載する。
- Phase 2〜4 対象タスクでは、`Requirements` に `memx_spec_v3/docs/design-doc-dod-spec.md` 参照を必須で明記する（最終判定の正本）。

### Commands
- 実行・検証コマンドを列挙する。
- 実行順を持つ場合は上から順に並べる。
- 必須例として記載するのは、リポジトリ内で実在し運用対象のコマンドのみとする。
  - lint/type/test の基準は `docs/QUALITY_GATES.md` に従う。
  - 現行の必須最小構成は `go test ./...`（Go）で、Python/Node は対象外として扱う。
  - 仕様書作成・更新タスクでも同じ判定基準（`docs/QUALITY_GATES.md`）を `Commands` に記載し、Task 起票時の誤記載を防止する。
- Source/Commands の記載例（抽出表作成・検証の最小セット）
  - `mkdir -p memx_spec_v3/docs/reviews/inventory`
  - `date +%Y%m%d`
  - `test -f memx_spec_v3/docs/reviews/inventory/DESIGN-SOURCE-INVENTORY-$(date +%Y%m%d).md`
  - `rg -n "^\| .* \| REQ-" memx_spec_v3/docs/reviews/inventory/DESIGN-SOURCE-INVENTORY-$(date +%Y%m%d).md`
  - `rg -n "\| blocked \|" memx_spec_v3/docs/reviews/inventory/DESIGN-SOURCE-INVENTORY-$(date +%Y%m%d).md`
  - `rg -n "run_id|generated_at|source_commit|chapter_id|tool|status|severity_summary|evidence_paths" <artifact-path>`
  - `python -c 'import json,sys;d=json.load(open(sys.argv[1]));keys=["run_id","generated_at","source_commit","chapter_id","tool","status","severity_summary","evidence_paths"];print([k for k in keys if k not in d]);sys.exit(1 if any(k not in d for k in keys) else 0)' <artifact.json>`
- 品質ゲート参照スコープは「memx 本体（`memx_spec_v3/`）は `docs/QUALITY_GATES.md`、`workflow-cookbook/` は `workflow-cookbook/docs/QUALITY_GATES.md`」と明記する。
- Phase 2〜4 対象タスクでは、章別検証サマリの作成・更新・添付確認コマンド（または確認手順）を `Commands` に明記する。

### Dependencies
- 前提タスク、依存PR、外部条件を列挙する。
- 依存がない場合は `- none` と記載する。
- `Source` / `Dependencies` で `contracts.md` / `CONTRACTS.md` を参照する場合は、正規パス **`memx_spec_v3/docs/CONTRACTS.md`** のみを許可する。
- 次のいずれかを検出した場合は誤参照として `fail` 扱い（レビュー差し戻し）にする。
  - `Source` / `Dependencies` に `memx_spec_v3/docs/contracts.md`（小文字）または `contracts.md#...` が残存する。
  - `Source` / `Dependencies` に `docs/*.md`（root `docs/IN-*.md` を除く）を正本として扱う参照が残存する。
  - `Source` / `Dependencies` が `path#Section` 形式を満たしていても、`memx_spec_v3/docs/design-reference-resolution-spec.md` の正規パスマッピングへ解決されていない。
  - `Source` / `Dependencies` に `requirements.md` / `design.md` / `interfaces.md` / `traceability.md` / `CONTRACTS.md` / `EVALUATION.md` / `RUNBOOK.md` の曖昧名（相対名のみ）が残存する。

### Incident Trace
- Incident 由来の Task では、`docs/IN-<実日付>-<連番>.md#<対象章>` を記載し、`memx_spec_v3/docs/incident-to-task-traceability-spec.md` に従って `Requirements` / `Commands` / `Dependencies` へ転記する。
- `Commands` には Incident 検証コマンドを 1 件以上記載する。
- Incident 由来でない Task は `- none` と記載する。

### Release Note Draft
- `CHANGELOG.md` に転記する利用者影響の要約を1〜3行で記載する。
- 実装詳細ではなく「何が変わるか」「利用者への影響」を簡潔に記載する。

### Status
- ステータス語彙は HUB と同一の次のみを許可する。
  - planned
  - active
  - in_progress
  - reviewing
  - blocked
  - done
- `Status: done` へ遷移する条件として、Task Seed に `Release Note Draft` 記入済みであること。
- `Status: done` へ遷移する条件として、移送後に `Moved-to-CHANGES: YYYY-MM-DD` を追記済みであること。
- `Status: done` へ遷移する条件として、Phase 2〜4 対象タスクは最新の `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304-003.md` が更新済みで、対象章の `req_coverage` が `100%` または `集計対象外` で、`mapping_match_check` が `pass` であること。
- `docs/birdseye/index.json` と `nodes[].capsule` 実体の不整合を検知した場合は、対象 Task Seed の `Status` を `blocked` へ遷移し、欠落 capsule のパス・検知コマンド・暫定対処を Task Seed に記録する。
- Birdseye Readiness Check の `ready/degraded/blocked` は入力健全性の判定語彙であり、Task Seed `Status` には使用しない（`degraded` は不可）。

### 完了前チェック（`Release Note Draft` / `Status` / `Moved-to-CHANGES`）
- [ ] `Release Note Draft` を 1〜3 行で記載済み（`CHANGELOG.md` 反映内容と一致）
- [ ] `Status` が許可語彙（planned/active/in_progress/reviewing/blocked/done）のいずれかで、`done` へ遷移する根拠を記載済み
- [ ] `Moved-to-CHANGES: YYYY-MM-DD` を追記済み、または未移送理由を明記済み
- [ ] Phase 2〜4 対象タスクは `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304-003.md` の対象章行が更新済みで、`req_coverage = 100%` または `集計対象外` かつ `mapping_match_check = pass` を満たす
- [ ] 上記3項目は `memx_spec_v3/docs/design-deliverables-package-spec.md` の成果物パッケージ要件と矛盾しない

#### Phase Gate 判定との対応表
| Status | gate判定 | 使用条件 | 差し戻し時の更新 |
| --- | --- | --- | --- |
| planned | 未評価 | 起票直後で gate 列未記入 | gate 列初期化後に `active` へ更新 |
| active | hold/pass | 着手済み、実作業待ち | hold 要因を `Requirements`/`Commands` に追記 |
| in_progress | hold/pass | 実装・文書更新を進行中 | fail 検知時は即 `blocked` へ更新 |
| reviewing | hold | レビューで再判定中 | 指摘反映後に gate 列を再計算し `active` へ戻す |
| blocked | fail | gate 列に `high` が1つ以上 | fail 原因・再開条件・再検証コマンドを追記 |
| done | pass | gate 列が全て `low` かつ完了条件充足 | 差し戻し時は `reviewing` に戻し `Moved-to-CHANGES` を見直す |

#### 差し戻し時の更新手順（必須）
1. 差し戻しを検知した時点で `Status` を `reviewing` または `blocked` へ更新する（gate=hold は `reviewing`、gate=fail は `blocked`）。
2. `memx_spec_v3/docs/design-phase-gate-spec.md` の5列（`gate_blocker`/`gate_req_coverage`/`gate_contract_high`/`gate_birdseye_issue`/`gate_hub_source_coverage`）を再評価し、変更前後を Task Seed に追記する。
3. fail/hold の原因を `Requirements` と `Commands` に反映し、再実行コマンドを先頭に追記する。
4. 修正後に gate 再判定を実施し、`pass` になった場合のみ `active` → `in_progress` → `reviewing` → `done` の順で復帰する。

## 2-1. 破壊変更時の追記チェックリスト（追加必須）
CLI/API の既存必須フィールド削除、型変更、意味変更、既存コマンド/エンドポイント/エラーコード削除、`--json` 既定出力の非同型化を含む場合、Task Seed に次のチェックリストを追記する。

- [ ] `Source` は `path#Section` で記載済み
- [ ] `Node IDs` を記載済み（依存照合対象なら必須）
- [ ] `Requirements` に後方互換/非機能制約を明記済み
- [ ] エラーコード変更時は `memx_spec_v3/docs/requirements.md` と `memx_spec_v3/docs/error-contract.md` を更新対象に含めた
- [ ] 契約変更時は `memx_spec_v3/docs/contracts/openapi.yaml` と `memx_spec_v3/docs/contracts/cli-json.schema.json` を更新した
- [ ] `Commands` に検証コマンドを順序付きで記載済み
- [ ] `Release Note Draft` を記載済み
- [ ] `memx_spec_v3/CHANGES.md` と `CHANGELOG.md` への反映項目を記載済み
- [ ] `Status: done` 前に `Moved-to-CHANGES: YYYY-MM-DD` を追記する

## 2-1-1. Task Seed 起票時の参照解決チェック（常時必須）

- [ ] `Source` の全行が `path#Section` 形式で、`#Section` が空でない
- [ ] `requirements.md` / `design.md` / `interfaces.md` / `traceability.md` / `contracts.md` / `CONTRACTS.md` / `EVALUATION.md` / `RUNBOOK.md` / `docs/birdseye/index.json` を参照する場合、`memx_spec_v3/docs/design-reference-resolution-spec.md` の正規パスへ解決済み
- [ ] 相対名・曖昧名・複数候補解決が 0 件（1件でもあれば fail し、`reviewing` で差し戻し）
- [ ] Phase 1 抽出表（Design Source Inventory）の `source_path#section` にも、上記 3 チェックを同一条件で適用済み

## 2-1-2. Task Seed 起票時のトリガー判定チェック（変更案）

- [ ] `memx_spec_v3/docs/design-update-trigger-spec.md` のトリガー種別（TRG-REQ / TRG-OAS / TRG-CLI-SCHEMA / TRG-RUNBOOK / TRG-IN-PREVENTIVE）から該当IDを記載した
- [ ] トリガーごとの必須更新先（design/interfaces/traceability/EVALUATION/operations/レビュー記録）をチェックし、非該当は理由付きで明記した
- [ ] `CHANGELOG.md` / `memx_spec_v3/CHANGES.md` の反映要否を判定し、`Release Note Draft` と矛盾がない
- [ ] `Status: done` 遷移前チェックとして `Moved-to-CHANGES: YYYY-MM-DD` 要否を確定した

## 2-1-3. Phase 3（契約整合）チェック（必須）

- [ ] `memx_spec_v3/docs/contracts/reports/` 配下に `CONTRACT-ALIGN-YYYYMMDD-###.md` と `LATEST.md` が存在することを確認した
- [ ] レビュー証跡として `memx_spec_v3/docs/reviews/CONTRACT-ALIGN-<実日付>.md`（例: `CONTRACT-ALIGN-20260304.md`）の命名規則で実体を作成し、受け入れレポートから1ホップ参照できることを確認した
- [ ] 更新順序を `CONTRACT-ALIGN作成 -> LATEST更新 -> EVALUATION照合` で実施し、逆順更新がないことを確認した
- [ ] `memx_spec_v3/docs/contracts/reports/LATEST.md` の必須キー（`report_id/report_path/decision_date/high_count/phase3_status`）を schema 固定として全件記載されていることを確認した（1件でも欠落なら fail）
- [ ] `LATEST.md` は当該判定のたびに毎回更新し、`report_id/report_path/decision_date/high_count/phase3_status` の5キーを前回値のまま残置していないことを確認した
- [ ] `memx_spec_v3/docs/contracts/reports/LATEST.md` の必須キー値が最新レポートと整合していることを確認した
- [ ] `EVALUATION.md` の `report_id/report_path` 参照と、`memx_spec_v3/docs/contracts/reports/LATEST.md` の `report_id/report_path` が一致することを確認した
- [ ] `high_count > 0` またはレポートID不一致時は Phase 3 の `Status` を `reviewing/blocked` に固定し、`done` へ遷移していないことを確認した

### 2-1-3-0. Task Seed 分岐ルール（Phase 3 固定）

- `high_count = 0` かつ `LATEST.md` と `EVALUATION.md` の `report_id/report_path` 一致時のみ `Status: done` へ遷移可能。
- `high_count > 0` の場合は `Status: blocked` を基本とし、レビュー再判定中のみ `Status: reviewing` を許可する。
- `report_id` または `report_path` が不一致の場合は `Status: reviewing` 固定とし、`done` / `blocked` へ遷移しない。

## 2-1-3-1. レビュー実体ファイルSLAチェック（Task Seed close 条件連動・必須）

- [ ] 実体作成確認: `memx_spec_v3/docs/reviews/` 配下に対象実体（`DESIGN-REVIEW-<実日付>-<連番>.md` / `DESIGN-CHAPTER-VALIDATION-<実日付>.md` / `DESIGN-ACCEPTANCE-<実日付>.md`）が作成済みである
- [ ] SLA期限確認: 作成・再提出が `memx_spec_v3/docs/design-review-artifact-sla-spec.md` の提出期限（Phase遷移前必須、差戻し時2営業日以内）を満たしている
- [ ] 未提出時制約確認: 未提出時は Task Seed の `Status` を `reviewing` に維持し、`done` へ遷移していない
- [ ] close条件連動: 上記いずれか未充足の場合、Task Seed close（`Status: done`）を禁止し、差戻し理由と再提出期限を記録する

### 起票時タスク化提案（競合回避のため分離）
- 提案1: 「レビュー実体の存在・命名チェック（3成果物）」専用 Task Seed を先行起票する。
- 提案2: 「SLA期限/Status遷移制約チェック（close判定連動）」専用 Task Seed を後続起票する。

## 2-1-4. 受け入れレポート定型チェック（必須）

- [ ] DA-LC-03/04 手順確認: `memx_spec_v3/docs/design-acceptance-lifecycle-spec.md` の命名・作成タイミングに従い、`memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<実日付>.md` を新規作成した
- [ ] DA-LC-01/02 手順確認: テンプレート `DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用として扱い、実体判定・Status遷移・Release判定の証跡に使用せず、既存実体の上書き/改名もしていない
- [ ] 作成済みチェック: `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<実日付>.md` が存在する
- [ ] 命名正当性チェック: 実体ファイル名が `DESIGN-ACCEPTANCE-<実日付>.md` に一致し、`DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用として保持されている
- [ ] DA-LC-05 必須6項目チェック: 実体に「対象章/REQ網羅率/high差分件数/リンク不達件数/Birdseye issue件数/最終判定」が全て記載されている
- [ ] 最終判定一致チェック: 受け入れレポートの最終判定が `memx_spec_v3/docs/design-doc-dod-spec.md` の判定結果と一致する
- [ ] 証跡仕様準拠確認: Phase 対象の証跡（lint/type/test/link/contract/birdseye/coverage）が `memx_spec_v3/docs/design-gate-evidence-spec.md` の保存先・命名規則・最小記録粒度（実行日時・コマンド・結果・判定）に準拠している

### 2-1-4-1. design-doc-dod-spec 判定一致チェック手順（完了条件）
- [ ] `memx_spec_v3/docs/reviews/DESIGN-GATE-EVIDENCE-INDEX.md` を一次参照として使用し、5条件（`req_coverage` / `mapping_match_check` / `contract_align_report` / `design_acceptance` / `go.sum tracked`）の判定値・根拠ファイル・根拠見出しを確認した
- [ ] `memx_spec_v3/docs/design-doc-dod-spec.md` の 6 条件（REQ網羅率/high差分/リンク不達/Birdseye issue/レビュー記録完備/参照解決適合率）を入力証跡で照合した
- [ ] `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-<実日付>.md` の「6. 最終判定」が上記照合結果と一致することを記録した
- [ ] 不一致時は `Status: reviewing` を維持し、差戻し理由を Task Seed に追記した

### 起票時タスク化提案（競合回避のため分離）
- 提案1: 「受け入れレポート作成/命名検証」を行う Task Seed を先行起票し、DA-LC-03/04 を固定する。
- 提案2: 「最終判定一致の照合（DoD照合 + evidence_paths検証）」を行う後続 Task Seed を分離する。

### 起票時タスク化提案（2-1-3: 競合回避のため分離）
- 提案1: 「Trigger 判定のみ」を行う 0.5d Task Seed を先行起票し、該当 Trigger IDs と必須更新先を固定する。
- 提案2: 「文書反映 + レビュー記録 + CHANGES 転記」を行う後続 Task Seed を分離し、`done` 条件を満たした時点で統合する。

## 2-2. 共通チェックリスト

- [ ] `memx_spec_v3/go/go.mod` を変更した場合は、`memx_spec_v3/go/go.sum` を同一コミットで同時更新している
- [ ] 差分に `memx_spec_v3/go/go.mod` または `memx_spec_v3/go/go.sum` が含まれる場合、`git diff --name-only` で両方が差分に含まれていることを確認している


## 2-3. 変更タイプ別チェックリスト（requirements 0-0-4 整合）

### 互換維持変更
- [ ] `Requirements` に後方互換維持（CLI/API/`--json` 同型）を明記する
- [ ] `Commands` に `docs/QUALITY_GATES.md` で運用対象の lint/type/test を記載する（対象外言語は除外）
- [ ] `Release Note Draft` を記載する

### 破壊変更
- [ ] 本書「2-1. 破壊変更時の追記チェックリスト」を全件追記する
- [ ] 対象 I/F・移行先・移行期限・移行手順（2ステップ以上）を記載する
- [ ] `CHANGELOG.md` と `memx_spec_v3/CHANGES.md` の双方へ同日反映する

### 実験機能（feature flag 既定 OFF）
- [ ] feature flag 名・既定値 OFF・有効化条件を `Requirements` に明記する
- [ ] 既定挙動に影響しないことを `Requirements` に明記する
- [ ] 廃止/昇格条件（次マイナー or 次メジャー）を記載する

## 3. CHANGES 連携ルール（memx_spec_v3/CHANGES.md / CHANGELOG.md）
- 正本（canonical source）はリポジトリルートの `CHANGELOG.md` とする。
- `memx_spec_v3/CHANGES.md` は v3 仕様の履歴・互換性破壊テンプレート管理用の補助台帳として扱う。

完了タスクは次の手順で `memx_spec_v3/CHANGES.md` と `CHANGELOG.md` に反映する。

1. Task Seed の `Status` を `done` に更新する。
2. Task Seed の Objective/Requirements の要点を 1〜3 行で要約する。
3. `memx_spec_v3/CHANGES.md` の該当バージョン節に箇条書きで追記する。
4. 同内容を `CHANGELOG.md` に 1〜3 行の最小要約で追記する（重複エントリは禁止）。
5. 互換性破壊がある場合は `memx_spec_v3/CHANGES.md` の「互換性破壊時の記載テンプレート」を必ず併記する。
6. 移送後、Task Seed 側には `Moved-to-CHANGES: YYYY-MM-DD` を追記してトレース可能にする。

## 4. 変更タスク起票（Design Doc DoD 正本参照の追加）
以下は、Phase 2〜4 向け Requirements 例へ本仕様参照必須ルールを反映するための変更タスク案（起票用）。

- [ ] Task A: `docs/TASKS.md` の `Requirements` 節（Phase 2〜4 ルール）に `memx_spec_v3/docs/design-doc-dod-spec.md` の正本参照必須を追加する。
- [ ] Task B: `orchestration/memx-design-docs-authoring.md` の Phase 3/4 Done Criteria を本仕様参照へ置換する差分タスクを起票する（本文置換は別PR）。
- [ ] Task C: `memx_spec_v3/docs/design-acceptance-report-spec.md` / `memx_spec_v3/docs/design-review-spec.md` へ「最終判定の正本」相互参照を維持する運用チェックを Task Seed の `Requirements` に追加する。
  - Commands:
    - `date +%Y%m%d`
    - `test -f memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md` （テンプレート専用ファイルの存在確認のみ）
    - `rg -n "^## 1\. 対象章|^## 2\. REQ網羅率|^## 3\. high差分件数|^## 4\. リンク不達件数|^## 5\. Birdseye issue件数|^## 6\. 最終判定" memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md`
    - `rg -n "design-doc-dod-spec\.md" memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-YYYYMMDD.md`
    - `go test ./...`
  - 完了条件:
    - `DESIGN-ACCEPTANCE-YYYYMMDD.md` はテンプレート専用として必須6項目章立てのみを保持し、実体判定には使用しない。
    - 入力元6仕様（requirements-coverage / contract-alignment / link-integrity / birdseye / design-review / chapter-validation）の参照リンクが記載されている。
    - 判定ロジックの記述は `memx_spec_v3/docs/design-doc-dod-spec.md` 参照のみに統一され、重複ロジックを含まない。

## 5. Phase 1 抽出表の転記手順（Source/Requirements/Dependencies）
1. `memx_spec_v3/docs/reviews/inventory/DESIGN-SOURCE-INVENTORY-YYYYMMDD.md` から `blocked=0` 行のみ抽出する。
2. `source_path#section` を `Source`、`req_id` を `Requirements`、`depends_on` を `Dependencies` に 1:1 で転記する。
3. 転記時は `memx_spec_v3/docs/design-reference-resolution-spec.md` に従い、正規パス以外を差し戻す。
