---
intent_id: INT-001
owner: your-handle
status: active   # draft|active|deprecated
last_reviewed_at: 2025-10-14
next_review_due: 2025-11-14
---

# Evaluation

## Acceptance Criteria

- 必須要件（フォーマット・件数・整合性など）
- PR本文に Priority Score（値と根拠）が記録されていること。
- governance/policy.yaml の forbidden_paths を変更しないこと。
- インシデント発生時は docs/IN-YYYYMMDD-XXX.md を作成し、該当PRおよびRUNBOOKから相互リンクする
- 要件カバレッジが `req_coverage=100%` を満たすこと（[docs/qa/req_coverage_report.md](docs/qa/req_coverage_report.md)）。
- 検収時は [docs/qa/contract_align_report.md](docs/qa/contract_align_report.md) を必須参照とし、契約整合（aligned / mismatch）と是正タスクの最新化を確認する。

## KPIs

| 指標 | 目的 | 収集方法 | 目標値 |
| --- | --- | --- | --- |
| `checklist_compliance_rate` | ドキュメント出荷時に必須チェックリストへ準拠できた割合を可視化し、ヒューマンエラーの早期検知につなげる。 | `python -m tools.perf.collect_metrics --suite qa --metrics-url <Prometheus URL> --log-path <StructuredLogger 等の運用ログ>` で `.ga/qa-metrics.json` を生成し、`checklist_compliance_rate` を参照する。必要に応じて構成管理ログ（例: `docs/logs/docops.log`）で完了数と対象総数を突合する。詳細は [RUNBOOK.md#Observability](RUNBOOK.md#observability)。 | 週次平均で 0.95 以上。 |
| `task_seed_cycle_time_minutes` | Task Seed の受付から初回処理完了までの所要時間を把握し、着手遅延を抑制する。 | `.ga/qa-metrics.json` に正規化される `task_seed_cycle_time_*` 系イベントを参照する（収集 CLI は [RUNBOOK.md#Observability](RUNBOOK.md#observability) に準拠）。 | 1440 分（24 時間）以下を維持。 |
| `birdseye_refresh_delay_minutes` | Birdseye ダッシュボードの更新遅延を監視し、情報可視化の鮮度を保証する。 | Prometheus の `birdseye_refresh_delay_*` 系メトリクスを CLI が平均化し `.ga/qa-metrics.json` に書き出す。必要に応じてジョブ監視ログで遅延の有無を確認する。 | 60 分以下を維持。 |
| `review_latency` | レビュー待機時間を定量化し、ボトルネックを可視化する。 | Prometheus の `workflow_review_latency_*` / `legacy_review_latency_*` を CLI が正規化した値を `.ga/qa-metrics.json` から取得する。詳細は [RUNBOOK.md#Observability](RUNBOOK.md#observability)。 | 12 時間以下を維持。 |
| `compress_ratio` | トリミング後コンテキストの圧縮率を測定し、情報損失を防ぐ。 | `tools.perf.metrics_registry.MetricsRegistry.observe_trim` を通じてエクスポートし、収集 CLI が `.ga/qa-metrics.json` に書き出した値を確認する。 | 0.60 以下を維持し、過剰圧縮を回避。 |
| `semantic_retention` | コンテキストトリミング後に保持された意味情報の割合を監視し、質の劣化を検知する。 | 埋め込みログ（例: `StructuredLogger` 経由の `semantic_retention`）を CLI が統合し `.ga/qa-metrics.json` へ出力する。手順は [RUNBOOK.md#Observability](RUNBOOK.md#observability) を参照。 | 0.85 以上を維持。 |
| `reopen_rate` | 再オープン率を追跡し、運用完了後の手戻りを抑制する。 | Prometheus の `workflow_reopen_rate_*` → `docops_reopen_rate` → `reopen_rate` を収集 CLI が正規化し `.ga/qa-metrics.json` に出力する。 | 5% 以下を維持。 |
| `spec_completeness` | スペック充足率を定量化し、要求事項の欠落を防ぐ。 | Prometheus の `workflow_spec_completeness_*` と運用ログ（例: `StructuredLogger` の `spec_completeness_*`）を CLI が統合し `.ga/qa-metrics.json` で欠損や乖離を確認する。 | 90% 以上を維持。 |

> KPI の収集手順と CLI オプションは常に [`RUNBOOK.md#Observability`](RUNBOOK.md#observability) と同期し、差異があれば双方を更新する。

## Test Outline

- 単体: 入力→出力の例テーブル（[ケース I-03](docs/addenda/I_Test_Cases.md#i-03-task-seed-tdd-例)）
- 結合: 代表シナリオ（[ケース I-02](docs/addenda/I_Test_Cases.md#i-02-birdseye-再生成確認)）
- 回帰: 重要パス再確認（[ケース I-01](docs/addenda/I_Test_Cases.md#i-01-チェックリスト突合)）

> 検証手順の詳細は [`docs/addenda/I_Test_Cases.md`](docs/addenda/I_Test_Cases.md) を参照する。

## Verification Checklist

- [ ] 主要フローが動作する（手動確認）
- [ ] エラー時挙動が明示されている
- [ ] 依存関係が再現できる環境である
