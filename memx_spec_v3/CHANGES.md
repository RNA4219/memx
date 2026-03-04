# Changes

## v3 (requirements v1.3)


- [#incident:IN-20260303-002][#sev:SEV2][#action:mitigation] `POST /v1/gc:run` の flag OFF 契約を v1 正本に合わせて `INTERNAL/500` フォールバック継続へ統一（[記録](docs/IN-20260303-002.md)）。
- 移行条件: `service.ErrFeatureDisabled` 実装と `go/api/errors.go` の `FEATURE_DISABLED -> 409` 明示マッピングが同一リリースで有効化された時点で 409 へ移行する。
- CLI と API を分離：CLI は API の薄いラッパ。
- API（HTTP + in-proc）を追加：`/v1/notes:ingest`, `/v1/notes:search`, `/v1/notes/{id}` など。
- Service(usecase) 層を追加：短期ストアの ingest/search/get を最小実装。
- DB 層を `go/db` に整理（OpenAll を追加、MustOpenAll は互換）。


## 破壊変更時必須チェックリスト

- 適用条件: 既存 CLI/API の必須フィールド削除・型変更・意味変更、既存コマンド/エンドポイント/エラーコード削除、既定 `--json` 同型性の破壊。
- 記載先: `CHANGELOG.md` と `memx_spec_v3/CHANGES.md` の双方（同一日付で記載）。
- 必須項目:
  - [ ] 対象 I/F（API/CLI `--json`/エンドポイント名）
  - [ ] 変更種別（削除/型変更/意味変更）
  - [ ] 影響範囲（影響するクライアント・コマンド）
  - [ ] 移行先（新エンドポイント or `/v2`）
  - [ ] 移行期限（YYYY-MM-DD または次メジャー）
  - [ ] 移行手順（2ステップ以上）
  - [ ] 互換期間中の挙動（並行提供/警告表示/明示フラグ時のみ有効化）

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
  - 例: `[#incident:IN-20260303-002] search フィルタ条件の修正`
- 重大度を追記する場合は `[#sev:SEV1|SEV2|SEV3|SEV4]` を続けて付与する。
  - 例: `[#incident:IN-20260303-002][#sev:SEV2] ...`
- 恒久対策・暫定対策の区別が必要な場合は `[#action:permanent|mitigation]` を付与する。
- インシデント記録（`docs/IN-*.md`）へのリンクを同一エントリ内に必ず含める。
- 記録形式は次の順序で固定する: `[#incident:...][#sev:...][#action:...] <変更要約>（[記録](...)）`

### 運用例（固定フォーマット）

- `[#incident:IN-20260303-002][#sev:SEV2][#action:permanent] 検索条件の境界値判定を修正（[記録](../docs/IN-20260303-002.md)）`
