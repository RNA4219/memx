# Error Contract (v1)

`memx_spec_v3/docs/requirements.md` の「6-4. エラーモデル」を正本とし、本書は運用向け要約とする。

> request/response のフィールド契約は `memx_spec_v3/docs/requirements.md` の「6-3-1. v1必須3エンドポイント契約（`requirements.md` × `go/api/types.go` 照合）」を参照。

## ErrorCode 区分（requirements 6-4 同期）

| 区分 | Error code | HTTP status | 契約レベル | 実装メモ |
| --- | --- | --- | --- | --- |
| v1必須保証 | `INVALID_ARGUMENT` | `400 Bad Request` | MUST | `service.ErrInvalidArgument` または HTTP 層の入力不備（`invalid json` / `id is required`）で返却。 |
| v1必須保証 | `NOT_FOUND` | `404 Not Found` | MUST | `service.ErrNotFound` で返却。 |
| v1必須保証 | `INTERNAL` | `500 Internal Server Error` | MUST | 未分類エラーのフォールバック。 |
| v1.x拡張（feature/sentinel依存） | `CONFLICT` | `409 Conflict` | SHOULD | service sentinel（例: `ErrConflict`）+ `mapError` 明示マップ時のみ返却。未実装時は `INTERNAL` へフォールバック。 |
| v1.x拡張（feature/sentinel依存） | `GATEKEEP_DENY` | `403 Forbidden` | SHOULD | gatekeeper deny sentinel（例: `ErrGatekeepDeny`）+ `mapError` 明示マップ時のみ返却。未実装時は `INTERNAL` へフォールバック。 |

## 現行実装との差分注記

- `go/api/types.go` には `CONFLICT` / `GATEKEEP_DENY` が定義済み。
  ただし `go/api/errors.go` の現行マップは `ErrInvalidArgument` / `ErrNotFound` 優先で、それ以外を `INTERNAL` にフォールバックするため、sentinel 未実装時に 409/403 は実運用で返らない。
- `go/api/http_server.go` の `writeErr` は `CONFLICT=409` / `GATEKEEP_DENY=403` を処理可能。
  ただし上流から当該 `Error.code` が渡された場合に限り有効で、未実装時のフォールバックは `INTERNAL=500`。

## 代表 JSON スキーマ例

- `INVALID_ARGUMENT`: `{ "code": "INVALID_ARGUMENT", "message": "invalid argument", "details": {"field": "query"} }`（`details` は任意）
- `NOT_FOUND`: `{ "code": "NOT_FOUND", "message": "not found" }`
- `INTERNAL`: `{ "code": "INTERNAL", "message": "<internal error message>" }`
- `mapError` の分岐順は `ErrNotFound` → `ErrInvalidArgument` → 文字列救済 (`contains("invalid argument")`) → `INTERNAL`。
- HTTP ステータス変換は `writeErr` にて `INVALID_ARGUMENT=400`, `NOT_FOUND=404`, それ以外未定義コードは `500`。
- JSON のエラーレスポンスは `api.Error` 型（`code`, `message`, `details?`）。


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
