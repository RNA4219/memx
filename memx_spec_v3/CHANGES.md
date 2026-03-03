# Changes

## v3 (requirements v1.3)

- CLI と API を分離：CLI は API の薄いラッパ。
- API（HTTP + in-proc）を追加：`/v1/notes:ingest`, `/v1/notes:search`, `/v1/notes/{id}` など。
- Service(usecase) 層を追加：短期ストアの ingest/search/get を最小実装。
- DB 層を `go/db` に整理（OpenAll を追加、MustOpenAll は互換）。

## 互換性破壊時の記載テンプレート

- 対象: （API/CLI `--json`/エンドポイント名）
- 変更種別: （削除/型変更/意味変更）
- 影響範囲: （影響するクライアント・コマンド）
- 移行先: （新エンドポイント or `/v2`）
- 移行期限: （YYYY-MM-DD または 次メジャー）
- 移行手順:
  1. ...
  2. ...
- 互換期間中の挙動: （並行提供/警告表示など）

## インシデント起因の修正記録ルール（タグ付け規則）

- インシデントに起因する変更は、変更記録の先頭に `[#incident:<インシデントID>]` を付与する。
  - 例: `[#incident:IN-20260303-001] search フィルタ条件の修正`
- 重大度を追記する場合は `[#sev:SEV1|SEV2|SEV3|SEV4]` を続けて付与する。
  - 例: `[#incident:IN-20260303-001][#sev:SEV2] ...`
- 恒久対策・暫定対策の区別が必要な場合は `[#action:permanent|mitigation]` を付与する。
- インシデント記録（`docs/IN-*.md`）へのリンクを同一エントリ内に必ず含める。
