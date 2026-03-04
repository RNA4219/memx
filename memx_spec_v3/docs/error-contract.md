# Error Contract (v1)

`memx_spec_v3/docs/contracts/openapi.yaml` を API/Error 契約の正本とし、本書は `requirements.md` に従属する運用向け要約とする。

- 本書の位置づけと参照導線は `spec.md` を起点に確認する。
- request/response のフィールド契約は `memx_spec_v3/docs/contracts/openapi.yaml`（API）および `memx_spec_v3/docs/contracts/cli-json.schema.json`（CLI `--json`）を参照する。

## 正本運用ルール（重複定義解消）

- エラーコード・HTTP ステータス・エラーボディ構造の正本は `memx_spec_v3/docs/contracts/openapi.yaml` の `components.responses` / `components.schemas.Error*` とする。
- 本書に同一内容を再定義しない。必要な記載は「実装差分」「運用メモ」「移行手順」に限定する。
- 契約変更時は schema を先に更新し、その後に本書へ差分要約を反映する。

## ErrorCode追加時の必須更新対象

新規 `ErrorCode` を追加する場合は、下表の更新を 1 つでも欠いた時点でレビュー不合格とする。

| 更新対象 | 必須更新内容 | 目的 |
| --- | --- | --- |
| `go/api/errors.go` | service sentinel から API `Error.code` へのマッピングを追加し、既存コードの意味を変更しない。 | transport 応答で新規 code を返せる状態にする。 |
| `go/service/errors.go` | service 層の sentinel/error 定義を追加し、retryable 判定の根拠を service 層で確定する。 | ドメイン由来の再試行可否を transport から分離する。 |
| `docs/error-contract.md` | ErrorCode×HTTP×retryable マトリクスと差分注記を更新する。 | 運用解釈と実装の同期を維持する。 |
| `docs/interfaces.md` 4章 | 外部インタフェースのエラー契約（code / status / retryable）を更新する。 | 利用者向け契約を最新版へ同期する。 |

## retryable判定の責務境界

- service 層（`go/service/errors.go`）
  - `retryable=true/false` の一次判定責務を持つ。
  - 判定基準はドメイン要因（恒久失敗/一時失敗）とし、HTTP 都合で上書きしない。
- transport 層（`go/api/errors.go`, `go/api/http_server.go`）
  - service 判定結果を保持したまま、`ErrorCode -> HTTP status` の写像のみを担う。
  - transport 層で retryable の意味を再解釈・反転しない。

## 破壊的変更の禁止事項（要件ID）

| Requirement ID | 禁止事項 | 判定基準 |
| --- | --- | --- |
| `REQ-ERR-001` | 既存 `ErrorCode` の意味変更を禁止する。 | 同一 code の利用者影響（原因分類/対処方針）が変わる変更は不可。必要な場合は新規 code を追加する。 |
| `REQ-ERR-002` | 既存 `ErrorCode` の HTTP ステータス変更を禁止する。 | 既存 code の status を別値に変更する差分は不可。移行が必要な場合は別 code を導入し段階移行する。 |

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
| v1.x運用解釈（`/v1/gc:run`） | `NOT_FOUND` / `INTERNAL` | `404` / `500` | SHOULD | route 非公開時は `NOT_FOUND`、route 公開かつ flag OFF の現行実装は `INTERNAL`（FAILED_PRECONDITION 相当の暫定フォールバック）として扱う。 |


## 現行実装との差分注記

- `go/api/types.go` には `CONFLICT` / `GATEKEEP_DENY` が定義済み。
  ただし `go/api/errors.go` の現行マップは `ErrInvalidArgument` / `ErrNotFound` 優先で、それ以外を `INTERNAL` にフォールバックするため、sentinel 未実装時に 409/403 は実運用で返らない。
- `go/api/http_server.go` の `writeErr` は `CONFLICT=409` / `GATEKEEP_DENY=403` を処理可能。
  ただし上流から当該 `Error.code` が渡された場合に限り有効で、未実装時のフォールバックは `INTERNAL=500`。

## 追加ケース: GC flag OFF 時の運用解釈（`POST /v1/gc:run`）

`POST /v1/gc:run` は SHOULD 実験機能のため、運用では「公開可否」と「実行可否」を分けて解釈する。

- route 非公開（サーバーが `/v1/gc:run` をマウントしない）
  - HTTP status: `404 Not Found`
  - body: 標準 `NOT_FOUND` エラー
- route 公開かつ flag OFF（`mem.features.gc_short=false`）
  - HTTP status: `500 Internal Server Error`
  - body: 標準 `INTERNAL` エラー
  - 解釈: `FAILED_PRECONDITION` 相当の「実行前提未充足」を、現行 API 実装の都合で `INTERNAL` にフォールバックしている状態

## 実装整合メモ（`go/api/http_server.go`）

- `writeErr` は `INVALID_ARGUMENT=400` / `NOT_FOUND=404` / `CONFLICT=409` / `GATEKEEP_DENY=403` / その他 `500` を返す。
- `FEATURE_DISABLED` は現行 `ErrorCode` に未定義のため、flag OFF を専用コードで返さない。
- したがって `/v1/gc:run` の flag OFF は、上流が専用 sentinel を返さない限り `INTERNAL=500` と解釈する。
