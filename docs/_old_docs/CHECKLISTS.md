---
intent_id: memx-implementation-checklists-v1
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# CHECKLISTS

## 仕様反映チェック
- [ ] 4ストア（short/journal/knowledge/archive）の責務分離を維持している。
- [ ] CLI は API ラッパに徹し、DB 直接アクセスを持たない。
- [ ] Service を唯一の業務ロジック入口にしている。

## DB/マイグレーションチェック
- [ ] `schema/*.sql` 適用方式に統一している。
- [ ] DDL 適用順序（notes→fts→tags→embeddings→user_version）を守っている。
- [ ] `PRAGMA user_version` 運用ルール（初期1、破壊的変更のみ+1）を守っている。

## API/CLI 互換性チェック
- [ ] v1 必須コマンド（in short / out search / out show）が動作する。
- [ ] v1 必須 API（ingest / search / get）が動作する。
- [ ] v1 内で必須フィールド削除・意味変更・成功レスポンス破壊を行っていない。
- [ ] `--json` 出力互換を維持している。
- [ ] エラー契約（`INVALID_ARGUMENT` / `NOT_FOUND` / `INTERNAL`）に破壊変更がない。

## Safety/Policy チェック
- [ ] `docs/security/minimal_operations.md` のデータ分類（public/internal/secret）に沿って保存・出力・マスキングを確認している。
- [ ] Gatekeeper `deny`/`needs_human` 発生時の記録先・再試行条件・エスカレーション手順を確認している。
- [ ] Gatekeeper の `needs_human` を deny 相当で扱っている。
- [ ] 保存前/出力前フックを維持している。
- [ ] secrets ブロック方針（policy + gatekeeper）を守っている。

## 運用・品質チェック
- [ ] Birdseye鮮度確認（`docs/birdseye/index.json.generated_at` が7日以内、超過時は RUNBOOK 手順を実施）を完了している。
- [ ] PRテンプレ項目の記入完了を確認している。
- [ ] ドラフトノート確認（担当: リリース担当者、タイミング: PR マージ後〜リリースタグ作成前）を完了している。
- [ ] GC トリガ判定は `memory_policy.yaml.gc.short` のみを参照している。
- [ ] `short_meta` はヒント運用とし、GC 前に正確値再計算を行う。
- [ ] インシデント記録をテンプレート要件に沿って作成できる。
- [ ] 性能SLO達成確認（EVALUATION.md の計測条件で ingest/search/show の P50/P95 が閾値内）を完了している。
- [ ] `governance/metrics.yaml` と `EVALUATION.md` の必須指標定義が同期している。
- [ ] 設定ファイル不在時は `docs/QUALITY_GATES.md` に基づき対象外判定し、対象外言語の lint/type/test は必須扱いにしない（例: `pyproject.toml` / `pytest.ini` / `package.json` 不在）。
