---
intent_id: DOC-LEGACY
owner: docs-core
status: active
last_reviewed_at: 2025-10-28
next_review_due: 2025-11-28
---

# Roadmap & Specs ハブ

<!-- markdownlint-disable MD013 -->

本書は仕様・運用資料を横断的に束ねる更新ハブであり、各ドキュメントの役割と連携経路を即座に確認できます。タスク着手前やレビュー時、リリース判定前など節目の判断で参照し、関係資料の鮮度と整合を維持してください。
代表ドキュメントは `BLUEPRINT.md`・`RUNBOOK.md`・`CHECKLISTS.md`・`EVALUATION.md` を中心に、`HUB.codex.md` や `SECURITY.md` など運用基盤資料が支えます。いずれかを更新する際は本ハブと相互リンクを再点検し、全資料の記述差異がないよう同期してください。

## 上位ドキュメント索引

| レイヤ | ドキュメント | 主要内容 | 関連ワークフロー |
| :-- | :-- | :-- | :-- |
| 骨子 | [BLUEPRINT.md](../BLUEPRINT.md) | 課題・スコープ・I/O 契約を1ページに集約。 | 設計着手前にインテントを固め、`HUB.codex.md` の[入力ファイル分類](../HUB.codex.md#2-入力ファイル分類)へ渡す。 |
| 実行 | [RUNBOOK.md](../RUNBOOK.md) | 環境差分と準備→実行→確認の手順を列挙。 | インシデント時に [IN-20250115-001.md](IN-20250115-001.md) を参照し、`CHECKLISTS.md` の [Release](../CHECKLISTS.md#release) を更新。 |
| 検収 | [EVALUATION.md](../EVALUATION.md) | 受入基準・KPI・検証チェックリストを定義。 | ガバナンス CI で `governance/policy.yaml` の[禁止パス](../governance/policy.yaml)を監視。 |
| 行動指針 | [GUARDRAILS.md](../GUARDRAILS.md) | 行動原則と Birdseye 読込手順を規定。 | 「[Birdseye / Minimal Context Intake Guardrails](../GUARDRAILS.md#birdseye--minimal-context-intake-guardrails鳥観図最小読込)」を参照し、タスク分割前に対象ノードを確定。 |
| タスク統制 | [HUB.codex.md](../HUB.codex.md) | 仕様・運用ドキュメントを束ね、タスク化を自動分配。 | `TASK.*` シード生成時に [入力ファイル分類](../HUB.codex.md#2-入力ファイル分類) と [自動タスク分割フロー](../HUB.codex.md#3-自動タスク分割フロー) を参照。 |
| リリース判定 | [CHECKLISTS.md](../CHECKLISTS.md) | 日次・リリース・衛生のチェック項目。 | リリース前に `EVALUATION.md` の[Verification Checklist](../EVALUATION.md#verification-checklist)と突き合わせ。 |
| フォーク統制 | [docs/FORK_NOTES.md](FORK_NOTES.md) | 派生リポの独自差分・逆輸入方針・未解決課題を一元管理。 | フォーク更新時およびリリース前に `CHECKLISTS.md` の[Hygiene](../CHECKLISTS.md#hygiene)と突合。 |
| セキュリティ統制 | [SECURITY.md](../SECURITY.md) / [security/SAC.md](security/SAC.md) | 報告窓口と SAC 原則の適用範囲・是正判断を集約。 | セキュリティレビュー準備で [security/Security_Review_Checklist.md](security/Security_Review_Checklist.md) と `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)を突合し、対応責務を同期。 |
| 依存グラフ | [docs/BIRDSEYE.md](BIRDSEYE.md) / [docs/birdseye/index.json](birdseye/index.json) / [caps/*](birdseye/caps/) | ノード一覧とカプセルで Birdseye トポロジを提供し、Guardrails からの参照起点を集約。 | `codemap.update` で再生成し、`GUARDRAILS.md` の[鮮度管理](../GUARDRAILS.md#鮮度管理staleness-handling)に従って鮮度を監視。 |
| ガバナンス | [governance/policy.yaml](../governance/policy.yaml) / [prioritization.yaml](../governance/prioritization.yaml) / [metrics.yaml](../governance/metrics.yaml) | セルフモディフィケーション制御、優先度算出基準、定常メトリクス。 | `HUB.codex.md` の優先度判定および `CHECKLISTS.md` の衛生チェックで参照。 |
| 設計判断 | [docs/ADR/README.md](ADR/README.md) | ADR 一覧と作成手順、判断変更時のレビュー連携を統括。 | 設計変更 PR で更新・新規 ADR を提出し、レビューテンプレに添付して承認後にマージ。 |
| 仕様 | [docs/spec.md](spec.md) | レシピ仕様の原則と更新手続きを集約。 | テンプレ更新時に `TASK.codex.md` の[Task Seed Template](../TASK.codex.md#task-seed-template)と整合性を確認。 |
| 設計 | [docs/design.md](design.md) | ディレクトリ構成とアーキテクチャ意図を整理。 | 設計レビューで `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)項目と照合。 |
| 設計受入 | [docs/qa/design_acceptance.md](qa/design_acceptance.md) | 設計レイヤの公式受入観点（責務境界・互換性・例外方針・非機能影響）を判定付きで管理。 | 設計レビューとリリース判定で `docs/design.md`・`docs/interfaces.md`・`EVALUATION.md` の根拠リンクを突合。 |
| 要件 | [docs/requirements.md](requirements.md) | 要件トレーサビリティと受入観点を提示。 | ガバナンス確認で `EVALUATION.md` の[Acceptance Criteria](../EVALUATION.md#acceptance-criteria)とリンクを確認。 |
| I/O 契約 | [docs/CONTRACTS.md](CONTRACTS.md) | 外部連携の I/O 契約と feature detection の扱いを定義。 | 拡張実装時に `RUNBOOK.md` の[Execute](../RUNBOOK.md#execute)手順と突き合わせ。 |
| 境界定義 | [docs/interfaces.md](interfaces.md) | 機能境界・受け渡し契約をテーブル形式で管理。 | 並行開発時に責務衝突を避けるため、機能追加ごとに更新してレビューへ添付。 |
| セキュリティ審査 | [docs/security/Security_Review_Checklist.md](security/Security_Review_Checklist.md) | リリース前セキュリティチェック項目をフェーズ別に整理。 | 審査会議前に `SECURITY.md` と照合し、`CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)と結果を同期。 |

## ADR サマリ

「上位ドキュメント索引」で ADR の位置付けを確認したうえで、本節から `docs/ADR/README.md` のカテゴリ別整理へ遷移できます。判断の更新は下記の代表 ADR を起点にし、索引テーブルと相互リンクを点検してください。

- **ADR-001 基盤方針**（[概要](ADR/README.md#adr-core-policy)）：`BLUEPRINT.md` の[Constraints / Assumptions](../BLUEPRINT.md#3-constraints--assumptions)と同期し、アーキテクチャ前提の差異を把握します。更新時は `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)で証跡を確認し、`docs/ADR/README.md` の[レビューフロー](ADR/README.md#レビューフロー)と [PR テンプレート](../.github/pull_request_template.md)へ ADR 番号と関連資料を記録してからマージします。
- **ADR-010 マイルストーン DoD**（[概要](ADR/README.md#adr-milestone-dod)）：ロードマップの Done 条件と KPI エビデンスを整理し、進捗レポートと同期します。更新時は `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)で証跡を整え、`docs/ADR/README.md` の[レビューフロー](ADR/README.md#レビューフロー)と [PR テンプレート](../.github/pull_request_template.md)で共有します。
- **ADR-020 運用ガバナンス**（[概要](ADR/README.md#adr-operations-governance)）：`SECURITY.md` や `RUNBOOK.md` の差分を通じて統制判断を履歴化します。更新時は `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)で統制結果を確認し、`docs/ADR/README.md` の[レビューフロー](ADR/README.md#レビューフロー)と [PR テンプレート](../.github/pull_request_template.md)に差分リンクと影響範囲を添付してください。

## 実装ディレクトリ↔仕様対応

| ディレクトリ | 紐付く仕様ドキュメント | 備考 |
| :-- | :-- | :-- |
| `docs/birdseye/` | [GUARDRAILS.md](../GUARDRAILS.md#birdseye--minimal-context-intake-guardrails鳥観図最小読込) / [tools/codemap/README.md](../tools/codemap/README.md#実行手順) | Birdseye トポロジーを生成し鮮度を保つ[^birdseye] |
| `docs/security/` | [SECURITY.md](../SECURITY.md) / [security/Security_Review_Checklist.md](security/Security_Review_Checklist.md) | セキュリティレビューと SAC を同期[^security] |
| `examples/` | [design.md](design.md) / [spec.md](spec.md) | レシピの参照実装と設計・仕様を突き合わせる |
| `styles/` | [design.md](design.md) / [requirements.md](requirements.md) | QA ルールに基づく表記統一・禁止用語チェック[^styles] |
| `tools/` | [design.md](design.md) / [../RUNBOOK.md](../RUNBOOK.md#execute) | ドキュメント同期用スクリプトと運用手順を連結 |
| `tools/codemap/` | [tools/codemap/README.md](../tools/codemap/README.md) / [../CHECKLISTS.md](../CHECKLISTS.md#hygiene) | Birdseye 再生成 CLI と衛生チェックを担保[^codemap] |
| `governance/` | [../EVALUATION.md](../EVALUATION.md#acceptance-criteria) / [../governance/policy.yaml](../governance/policy.yaml) | 受入基準と禁止パス・優先度設定を同期[^governance] |
| `tests/` | [../EVALUATION.md](../EVALUATION.md#test-outline) / [birdseye/caps/](birdseye/caps/) | テストアウトラインと Birdseye カプセルを連携更新[^tests] |

[^birdseye]: `python tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps` を実行して `docs/birdseye/index.json`・`docs/birdseye/hot.json`・`caps/*` を再生成し、`GUARDRAILS.md` の鮮度管理基準を維持する。
[^security]: レビュー結果と是正策は `docs/security/SAC.md` に記録し、チェックリストでトレースする。
[^styles]: `styles/qa/QA.yml` の用語統一・禁止用語ルールをレビュー時に適用する。
[^codemap]: CLI 実行後は `CHECKLISTS.md` の Hygiene セクションで差分確認を行う。
[^governance]: 方針改定時は `EVALUATION.md` の Acceptance Criteria と照合し、必要に応じて優先度・禁止パスを更新する。
[^tests]: テスト追加時は `EVALUATION.md` の Test Outline と Birdseye カプセルの `tests` フィールドを更新する。

## ロードマップ

- **Guardrails 更新の定期レビュー**
  - [GUARDRAILS.md](../GUARDRAILS.md) の `next_review_due` を起点に方針を棚卸し。
  - 更新内容は `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release) でレビューし、`HUB.codex.md` の[Rules](../HUB.codex.md#rules)へ波及。
- **Birdseye 鮮度管理の自動化**
  - `docs/birdseye/index.json` の `generated_at` 監視と、`tools/codemap/update.py` の `--targets` / `--emit` フローを `RUNBOOK.md` の[Execute](../RUNBOOK.md#execute)へ追記予定。
  - ステータスは `GUARDRAILS.md` の[鮮度管理](../GUARDRAILS.md#鮮度管理staleness-handling)と[codemap 未実装時の暫定手順](../GUARDRAILS.md#codemap-未実装時の暫定手順)で追跡。
- **CI テンプレ整備**
  - `ci-config.md` の設定例を `README.md` の[再利用CIの呼び出し例](../README.md#再利用ciの呼び出し例下流リポ側)と整合させる。
  - `governance/policy.yaml` の `ci.required_jobs` を基準に、`CHECKLISTS.md` の[Daily](../CHECKLISTS.md#daily)で稼働状況をモニタ。
- **CI 段階的導入フレームの実装**
  - [docs/ci_phased_rollout_requirements.md](ci_phased_rollout_requirements.md) に定義した Phase 0〜3 の要件を順次満たし、CI 負荷を段階的に引き上げる。
  - フェーズ移行ごとに `governance/policy.yaml`・`RUNBOOK.md`・`CHECKLISTS.md` を更新し、`CHANGELOG.md` の `[Unreleased]` へ通番付きで記録する。

## 参照クイックリンク

- [docs/ci-config.md](ci-config.md)：CI プリセットの分岐条件と再利用手順を集約。**利用シーン**：CI 設定変更前に `CHECKLISTS.md` の[Daily](../CHECKLISTS.md#daily)で運用要件をクロスチェック。
- [docs/BIRDSEYE.md](BIRDSEYE.md) / [docs/birdseye/index.json](birdseye/index.json) / [docs/birdseye/hot.json](birdseye/hot.json) / [birdseye/caps/](birdseye/caps/) / [tools/codemap/README.md#実行手順](../tools/codemap/README.md#実行手順)：Birdseye トポロジーの参照起点と生成結果、運用手順を一括で把握。**利用シーン**：1. `BIRDSEYE.md` で確認手順とリンクを把握。2. `generated_at`（必要に応じてホットリスト項目の `last_verified_at`）を確認し鮮度閾値を超えた場合は同期対象にする。3. README の手順通り `python tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps` を実行し `docs/birdseye/index.json`・`docs/birdseye/hot.json`・`caps/*` を再生成。4. `CHECKLISTS.md` の[Hygiene](../CHECKLISTS.md#hygiene)と `GUARDRAILS.md` の[鮮度管理](../GUARDRAILS.md#鮮度管理staleness-handling)を突き合わせて差分と期限を監視。
- [docs/interfaces.md](interfaces.md)：機能境界と受け渡し契約をテーブル化。**利用シーン**：境界整理や責務調整時に `docs/CONTRACTS.md` と `RUNBOOK.md` の[Execute](../RUNBOOK.md#execute)を並行確認。
- [docs/INCIDENT_TEMPLATE.md](INCIDENT_TEMPLATE.md)：インシデント報告テンプレートとエスカレーション導線を定義。**利用シーン**：インシデント対応の初動で `RUNBOOK.md` の[Confirm](../RUNBOOK.md#confirm)を基点にメトリクス照合・記録更新・運用チャネル報告を完了し、`CHECKLISTS.md` の[Hygiene](../CHECKLISTS.md#hygiene)で未完了項目を洗い出す。
- [docs/ADR/README.md](ADR/README.md)：設計判断の記録・改訂フローを統括。**利用シーン**：設計変更 PR に更新・新規 ADR を添付し、レビューテンプレと連携。
- [docs/security/Security_Review_Checklist.md](security/Security_Review_Checklist.md)：セキュリティ審査項目と証跡収集ポイントを整理。**利用シーン**：リリース前審査で `SECURITY.md`・`docs/security/SAC.md`・`CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)を同期。
- [docs/qa/dependency_lock_audit.md](qa/dependency_lock_audit.md)：`go.mod`/`go.sum` の整合確認ログを記録。**利用シーン**：依存更新時やリリース前に lockfile 監査を実施し、`CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)・[Hygiene](../CHECKLISTS.md#hygiene)と監査結果を同期。
- [SECURITY.md](../SECURITY.md) / [docs/security/SAC.md](security/SAC.md)：報告窓口の連絡経路と SAC 拘束事項を統合し、対応判断の前提条件を提示。**利用シーン**：セキュリティ審査で通知窓口と拘束条件を共有し、インシデント初動前の前提確認で運用責務の充足を判定。
- [CHANGELOG.md](../CHANGELOG.md)：リリース差分と意思決定の履歴を集約し、更新ルールを一元管理。**利用シーン**：`CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)完了後に `README.md` の[使い方（最短）](../README.md#使い方最短)手順と照合してガバナンス記録を反映。
- [LICENSE](../LICENSE) / [CHECKLISTS.md#release](../CHECKLISTS.md#release) / [CHANGELOG.md](../CHANGELOG.md)：Workflow Cookbook 公開版と同様に配布物へ必須ライセンス・変更履歴・監査観点を束ねる。**利用シーン**：リリース成果物へ `LICENSE` を同梱し、`CHECKLISTS.md#release` の配布物チェックを踏まえて `CHANGELOG.md` の公開内容と突合する。
- [README.md](../README.md#変更履歴の更新ルール)：`CHANGELOG.md` の更新トリガー・書式・`CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)を用いた突合フローを整理。**利用シーン**：リリース確定後にチェックリスト→変更履歴→再確認の流れを短時間でなぞり、記録漏れを防ぐ。
- [EVALUATION.md#Test Outline](../EVALUATION.md#test-outline) / [tests/](../tests/)：評価指標とテストケース集を束ねた TDD 前提の検証ハブ。**利用シーン**：テスト追加前のチェックで指標・ケース網羅を見直し、Birdseye カプセル同期の要否を判断して TDD フローを開始。
- [docs/qa/mapping_match_check.md](qa/mapping_match_check.md) / [docs/qa/mapping_matrix.yaml](qa/mapping_matrix.yaml)：REQ ごとの design/contract/test 接続状況と判定結果を記録する QA マッピング。**利用シーン**：仕様レビュー時に未接続 REQ の有無を即時確認し、gap が 0 件であることをリリース判定に反映。

Guardrails 連動資料は行動原則と更新判断の基準を担い、本節は運用ドキュメントの即時参照に特化するため、改訂時は前述の整合チェック先と `GUARDRAILS.md` の[実装原則](../GUARDRAILS.md#実装原則)の適用範囲を併せて確認する。

## 更新フロー

1. **Guardrails ドキュメント改訂**
   - 方針変更は `BLUEPRINT.md` の[Constraints / Assumptions](../BLUEPRINT.md#3-constraints--assumptions)で前提を更新後、`GUARDRAILS.md` の[実装原則](../GUARDRAILS.md#実装原則)と[生成物に関する要求](../GUARDRAILS.md#生成物に関する要求出力契約)へ反映。
   - 改訂後、`EVALUATION.md` の[Acceptance Criteria](../EVALUATION.md#acceptance-criteria)および `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)で相互リンクを確認。
2. **Birdseye 再生成**
   - `GUARDRAILS.md` の[鮮度管理](../GUARDRAILS.md#鮮度管理staleness-handling)に沿って再生成条件を判定。
   - `tools/codemap/README.md` の[実行手順](../tools/codemap/README.md#実行手順)通り `python tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps` を実行し、`docs/birdseye/index.json`・`docs/birdseye/hot.json`・`caps/*` を更新。
   - ツール未整備時は `GUARDRAILS.md` の[codemap 未実装時の暫定手順](../GUARDRAILS.md#codemap-未実装時の暫定手順)に従って手動更新を依頼し、結果を `HUB.codex.md` の[Output Contract](../HUB.codex.md#output-contract)へ反映。

3. **lockfile 監査**
   - [docs/qa/dependency_lock_audit.md](qa/dependency_lock_audit.md) を起点に `go.mod`/`go.sum` の整合を確認し、監査記録テーブルへ結果を追記。
   - `CHECKLISTS.md` の[Release](../CHECKLISTS.md#release)・[Hygiene](../CHECKLISTS.md#hygiene)項目と突合し、`memx_spec_v3/go/go.sum` の tracked 状態を確認。

<!-- markdownlint-enable MD013 -->
