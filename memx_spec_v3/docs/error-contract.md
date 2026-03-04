# Error Contract (v1)

`memx_spec_v3/docs/contracts/openapi.yaml` を API/Error 契約の正本とし、本書は `requirements.md` に従属する運用向け要約とする。

- 本書の位置づけと参照導線は `spec.md` を起点に確認する。
- request/response のフィールド契約は `memx_spec_v3/docs/contracts/openapi.yaml`（API）および `memx_spec_v3/docs/contracts/cli-json.schema.json`（CLI `--json`）を参照する。

## 正本運用ルール（重複定義解消）

- エラーコード・HTTP ステータス・エラーボディ構造の正本は `memx_spec_v3/docs/contracts/openapi.yaml` の `components.responses` / `components.schemas.Error*` とする。
- 本書に同一内容を再定義しない。必要な記載は「実装差分」「運用メモ」「移行手順」に限定する。
- 契約変更時は schema を先に更新し、その後に本書へ差分要約を反映する。

## ErrorCode × HTTP × retryable × クライアント動作（運用マトリクス）

| ErrorCode | HTTP | retryable | クライアント動作 |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | `400 Bad Request` | false | 入力修正後に再実行（自動再試行しない） |
| `NOT_FOUND` | `404 Not Found` | false | 対象ID/検索条件を見直し、必要なら再投入後に再実行 |
| `CONFLICT` | `409 Conflict` | false | 競合解消後に再実行 |
| `GATEKEEP_DENY` | `403 Forbidden` | false | ポリシー・権限変更まで停止 |
| `FEATURE_DISABLED` | `409 Conflict` | false | 機能フラグ有効化まで停止 |
| `INTERNAL` | `500 Internal Server Error` | conditional | 一時障害時のみ指数バックオフ再試行（最大2回） |

## ErrorCode 区分（schema 同期）

| 区分 | Error code | HTTP status | 契約レベル | 実装メモ |
| --- | --- | --- | --- | --- |
| v1必須保証 | `INVALID_ARGUMENT` | `400 Bad Request` | MUST | `service.ErrInvalidArgument` または HTTP 層の入力不備（`invalid json` / `id is required`）で返却。 |
| v1必須保証 | `NOT_FOUND` | `404 Not Found` | MUST | `service.ErrNotFound` で返却。 |
| v1必須保証 | `INTERNAL` | `500 Internal Server Error` | MUST | 未分類エラーのフォールバック。 |
| v1.x拡張（feature/sentinel依存） | `CONFLICT` | `409 Conflict` | SHOULD | service sentinel（例: `ErrConflict`）+ `mapError` 明示マップ時のみ返却。未実装時は `INTERNAL` へフォールバック。 |
| v1.x拡張（feature/sentinel依存） | `GATEKEEP_DENY` | `403 Forbidden` | SHOULD | gatekeeper deny sentinel（例: `ErrGatekeepDeny`）+ `mapError` 明示マップ時のみ返却。未実装時は `INTERNAL` へフォールバック。 |
| v1.x拡張（feature flag依存） | `FEATURE_DISABLED` | `409 Conflict` | SHOULD | feature flag 無効時に返却（例: `gc_short` 無効）。 |

## 現行実装との差分注記

- `go/api/types.go` には `CONFLICT` / `GATEKEEP_DENY` が定義済み。
  ただし `go/api/errors.go` の現行マップは `ErrInvalidArgument` / `ErrNotFound` 優先で、それ以外を `INTERNAL` にフォールバックするため、sentinel 未実装時に 409/403 は実運用で返らない。
- `go/api/http_server.go` の `writeErr` は `CONFLICT=409` / `GATEKEEP_DENY=403` を処理可能。
  ただし上流から当該 `Error.code` が渡された場合に限り有効で、未実装時のフォールバックは `INTERNAL=500`。

## 追加ケース: GC feature disabled（`POST /v1/gc:run`）

`mem.features.gc_short=false` の場合は全環境共通で以下を返す。

- HTTP status: `409 Conflict`
- body:

```json
{
  "code": "FEATURE_DISABLED",
  "message": "gc_short feature is disabled"
}
```
