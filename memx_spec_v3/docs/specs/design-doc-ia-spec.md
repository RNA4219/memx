# Design Document IA Spec

## 目的
- 設計書作成開始時の参照入口を `design-doc-ia-spec.md` に統一し、正本・参照専用・生成物の更新順序を固定する。
- `requirements / design / interfaces / traceability / CONTRACTS / error-contract / operations` と各 `design-*-spec.md` の責務境界を明確化する。

## 相互参照
- 更新トリガー判定と必須更新先は [design-update-trigger-spec.md](./design-update-trigger-spec.md) を正本とする。
- 本書は IA（情報アーキテクチャ）観点の責務境界・更新優先順位を定義する。

## 責務境界マトリクス（基幹文書 × design-* spec）

| design-* spec | requirements | design | interfaces | traceability | CONTRACTS | error-contract | operations |
| --- | --- | --- | --- | --- | --- | --- | --- |
| design-doc-dod-spec.md | 設計完了判定条件の要件化 | DoD 判定の主対象 | I/F 記載完了条件 | 証跡完了条件 | 契約整合確認条件 | エラー条件の完了定義 | 運用観点の完了定義 |
| design-phase-gate-spec.md | ゲート判定要件 | 設計フェーズ進行条件 | I/F ゲート条件 | REQ 紐付けゲート | 契約整合ゲート | 失敗時ゲート基準 | 運用移行ゲート |
| design-review-spec.md | 要件観点レビュー | 設計レビュー本体 | I/F レビュー項目 | トレース妥当性 | 契約レビュー観点 | エラー設計レビュー | 運用レビュー観点 |
| design-chapter-validation-spec.md | 章構成の要件整合 | 章単位設計妥当性 | I/F 記述章の妥当性 | REQ 対応章の検証 | 契約記述章の検証 | エラー章の検証 | 運用章の検証 |
| design-chapter-node-mapping-spec.md | 要件ノード起点 | 設計章ノード主対象 | I/F ノード対応 | トレースノード対応 | 契約ノード対応 | エラーノード対応 | 運用ノード対応 |
| design-reference-resolution-spec.md | 要件参照解決 | 設計内参照解決 | I/F 参照解決 | トレース参照解決 | 契約参照解決 | エラー参照解決 | 運用参照解決 |
| design-reference-conformance-spec.md | 要件参照準拠 | 設計参照準拠 | I/F 参照準拠 | トレース参照準拠 | 契約参照準拠 | エラー参照準拠 | 運用参照準拠 |
| design-source-inventory-spec.md | 情報源の要件分類 | 設計情報源管理 | I/F 情報源管理 | トレース情報源管理 | 契約情報源管理 | エラー情報源管理 | 運用情報源管理 |
| design-source-inventory-operations-spec.md | 要件系運用情報源 | 設計運用情報源 | I/F 運用情報源 | トレース運用情報源 | 契約運用情報源 | エラー運用情報源 | 運用情報源の主対象 |
| design-evidence-schema-spec.md | 要件証跡スキーマ | 設計証跡スキーマ | I/F 証跡スキーマ | トレース証跡主対象 | 契約証跡スキーマ | エラー証跡スキーマ | 運用証跡スキーマ |
| design-evidence-template-spec.md | 要件証跡テンプレ | 設計証跡テンプレ | I/F 証跡テンプレ | トレース証跡テンプレ主対象 | 契約証跡テンプレ | エラー証跡テンプレ | 運用証跡テンプレ |
| design-acceptance-lifecycle-spec.md | 受入要件遷移 | 設計受入遷移 | I/F 受入遷移 | トレース受入遷移 | 契約受入遷移 | エラー受入遷移 | 運用受入遷移 |
| design-acceptance-report-spec.md | 受入結果の要件視点 | 設計受入レポート主対象 | I/F 受入結果記録 | トレース受入結果記録 | 契約受入結果記録 | エラー受入結果記録 | 運用受入結果記録 |
| design-update-trigger-spec.md | 要件変更トリガー | 設計更新トリガー | I/F 更新トリガー | トレース更新トリガー | 契約更新トリガー | エラー更新トリガー | 運用更新トリガー |

## 文書区分と更新優先順位

| 文書 | 区分 | 更新優先順位 | 更新ルール |
| --- | --- | --- | --- |
| `docs/requirements.md` | 正本 | 1 | 変更起点。関連文書更新前に確定。 |
| `docs/design.md` | 正本 | 2 | 要件反映後に更新。 |
| `docs/interfaces.md` | 正本 | 3 | 設計反映後に更新。 |
| `docs/traceability.md` | 正本 | 4 | 要件・設計・I/F 更新を同期反映。 |
| `docs/contracts/openapi.yaml` / `docs/contracts/cli-json.schema.json` | 正本 | 5 | 契約差分は traceability と同一変更で反映。 |
| `docs/error-contract.md` | 正本 | 6 | 契約・I/F 更新後にエラー契約を確定。 |
| `RUNBOOK.md` | 正本 | 7 | 運用手順への波及を最後に確定。 |
| `docs/CONTRACTS.md` | 参照専用 | 正本追随 | 正本契約から再生成・要約のみ。 |
| `docs/design-*-spec.md`（本書含む） | 参照専用 | 正本追随 | 正本運用ルール変更時のみ更新。 |
| `docs/reviews/*.md` / `docs/reviews/inventory/*.md` | 生成物（レビュー・受入レポート） | 最終 | 正本更新完了後に生成・記録。 |

## Trigger別 必須更新先マトリクス
凡例: `必須` / `条件付き` / `不要`

| Trigger | requirements | design | interfaces | traceability | CONTRACTS (machine-readable) | error-contract | operations | design-update-trigger-spec | design-doc-ia-spec |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TRG-REQ | 必須 | 必須 | 条件付き | 必須 | 条件付き | 条件付き | 条件付き | 必須 | 必須 |
| TRG-OAS | 条件付き | 条件付き | 必須 | 必須 | 必須 | 必須 | 条件付き | 必須 | 必須 |
| TRG-CLI-SCHEMA | 条件付き | 条件付き | 必須 | 必須 | 必須 | 条件付き | 条件付き | 必須 | 必須 |
| TRG-RUNBOOK | 条件付き | 条件付き | 条件付き | 条件付き | 不要 | 条件付き | 必須 | 必須 | 必須 |
| TRG-IN-PREVENTIVE | 条件付き | 条件付き | 条件付き | 必須 | 条件付き | 必須 | 必須 | 必須 | 必須 |

## 運用ルール
- Trigger 判定手順・Done 遷移条件は [design-update-trigger-spec.md](./design-update-trigger-spec.md) を参照する。
- 設計着手時は本書を最初に確認し、続けて `requirements.md` → `design.md` → `interfaces.md` → `traceability.md` の順で更新する。
