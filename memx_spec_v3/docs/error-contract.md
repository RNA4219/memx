# Error Contract (v1)

`memx_spec_v3/go/api/errors.go` と `memx_spec_v3/go/api/http_server.go` の現行実装に合わせた、v1 のエラー契約。

> request/response のフィールド契約は `memx_spec_v3/docs/requirements.md` の「6-3-1. v1必須3エンドポイント契約（`requirements.md` × `go/api/types.go` 照合）」を参照。

## 対応表

| Error code | HTTP status | JSON スキーマ例 | 代表原因（実装準拠） |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | `400 Bad Request` | `{ "code": "INVALID_ARGUMENT", "message": "invalid argument", "details": {"field": "query"} }` (`details` は任意) | `service.ErrInvalidArgument` が返る。`errors.go` の救済ロジックで、`err.Error()` に `"invalid argument"` を含む文字列ラップ例外も同コードに変換される。HTTP 層では JSON decode 失敗（`invalid json`）や `id` 欠落（`id is required`）も同コード。 |
| `NOT_FOUND` | `404 Not Found` | `{ "code": "NOT_FOUND", "message": "not found" }` | `service.ErrNotFound` が返る。 |
| `INTERNAL` | `500 Internal Server Error` | `{ "code": "INTERNAL", "message": "<internal error message>" }` | 上記以外のすべてのエラー。`errors.go` では `err.Error()` を `message` に格納。 |

## 実装照合メモ

- `mapError` の分岐順は `ErrNotFound` → `ErrInvalidArgument` → 文字列救済 (`contains("invalid argument")`) → `INTERNAL`。
- HTTP ステータス変換は `writeErr` にて `INVALID_ARGUMENT=400`, `NOT_FOUND=404`, それ以外未定義コードは `500`。
- JSON のエラーレスポンスは `api.Error` 型（`code`, `message`, `details?`）。
