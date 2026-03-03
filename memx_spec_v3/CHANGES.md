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
