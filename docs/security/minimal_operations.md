# memx最小版: セキュリティ運用補足

> 対象: ローカル運用の最小構成（単一運用者または小規模チーム）

## 関連ドキュメント
- [SECURITY.md](../../SECURITY.md)
- [付録 G: Security & Privacy](../addenda/G_Security_Privacy.md)

## 1. データ分類ごとの運用参照手順

1. 正本要件 [requirements 2-7-1](../../memx_spec_v3/docs/requirements.md#2-7-1-sensitivity-別の保存可否マスキング保持期間) を開き、対象データの `sensitivity` を確認する。
2. 保存可否・マスキング・保持期間は要件本文の値をそのまま採用する（本書へ数値を再定義しない）。
3. 判定が `deny` または `needs_human` の場合は [§3 手順](#3-gatekeeper-deny--needs_human-発生時の運用手順) を実施する。
4. サンプル/テストは `.env.example` 等のダミー値のみ使用し、実秘密はリポジトリへ含めない。

## 2. ログ保持期間と削除ルール（手順）

1. 対象ログ（アプリ実行ログ / Gatekeeper 判定ログ / 運用判断メモ）を収集する。
2. 保持期間を要件本文（2-7-1）で確認する（internal/secret の最小保持 90 日を含む）。
3. 期限超過分を週次で削除する（手動または定期ジョブ）。
4. `archive` 退避/物理削除を実行する場合は、要件本文（2-7-3 / 2-7-4）の actor/approval/audit に従う。
5. 監査ログ項目・要件IDは要件本文（2-7-2 / 2-7-3 / 2-7-4）を参照して確認する（`archive_move` / `archive_purge` は固定項目 `result` / `reason` / `retryable` / `owner` / `next_attempt_at` を必須確認）。
6. バックアップを保持する場合も要件本文と同一保持期間を適用する。

## 3. Gatekeeper `deny` / `needs_human` 発生時の運用手順

1. **即時停止**: 当該入力の保存/出力処理を止める（再送しない）。
2. **記録**: `docs/incidents/` 配下のテンプレートに準拠して、以下を残す。
   - 発生日時、コマンド/機能、判定（`deny` or `needs_human`）
   - 入力種別（内容そのものは必要最小限。secretは記録しない）
   - 実施した遮断措置、暫定対応、担当者
3. **再試行条件**:
   - 誤検知の根拠（ルール確認・入力修正）がある場合のみ1回再試行可。
   - 同一条件で再度 `deny`/`needs_human` の場合は再試行停止。
4. **エスカレーション**:
   - 同日内に2回以上発生、または secret 混入疑いがある場合は Security Champion（不在時 Maintainer）へ即時連絡。
   - 必要に応じて RUNBOOK のインシデント手順へ移行する。
