---
intent_id: INT-001
owner: docs-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-04-04
---

# Contract Align Report

## 対象契約一覧（`docs/CONTRACTS.md` 節名）

- Artifacts
- Config
- Conventions

## feature detection 前提の確認結果

| 契約節 | 確認対象 | 未提供時の期待挙動 | 確認結果 |
| --- | --- | --- | --- |
| Artifacts | `.ga/qa-metrics.json` | メトリクス拡張ファイルが未提供でも Cookbook 本体は継続動作する。必要時のみ生成・取り込みされる。 | 設計どおり（aligned） |
| Config | `governance/predictor.yaml` | 設定ファイルが未提供でも既定値で実行される。 | 設計どおり（aligned） |
| Conventions | feature detection（存在検出）全般 | 外部拡張がなくても正常動作し、提供時のみ拡張を有効化する。 | 設計どおり（aligned） |

## 関連 Runbook / Checklist との整合

- Runbook: `RUNBOOK.md#Observability` で `.ga/qa-metrics.json` の生成確認と欠損時ハンドリングが定義されており、Artifacts 契約と整合。 
- Checklist: `CHECKLISTS.md#Development` で Runbook 連携と TDD フロー確認が要求され、契約逸脱時の差分反映導線がある。 
- Checklist: `CHECKLISTS.md#Pull-Request--Review` に本レポート更新確認を追加し、レビュー時に契約整合を必須確認化した。 

## 判定と是正タスク

- 判定: **aligned**
- 是正タスク:
  1. `EVALUATION.md` に本レポートへの必須参照を維持し、検収で契約整合確認を固定化する。
  2. PR / Review 時に `contract_align_report` 更新確認を継続し、契約変更の見落としを防止する。

