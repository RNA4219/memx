---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - API

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 6. API 要件（v1.3 追加）

- Requirement ID: `REQ-API-001`

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`
- `memx_spec_v3/docs/quickstart.md`

### 6-1. 目的

- ツール/AI から呼びやすい **安定 JSON I/F** を提供する。
- CLI は API の薄いラッパとして実装し、ビジネスロジックを持たない。

### 6-2. 提供形態

- **HTTP（ローカル）**：`mem api serve` で起動。
  - 例：`http://127.0.0.1:7766`
  - 将来的に unix socket 対応も想定。
- **in-proc**：CLI や別バイナリが同一プロセスで呼ぶ。

### 6-3. エンドポイント（v1）
- ADR: [ADR-0002: v1必須3エンドポイント固定方針](../../docs/ADR/ADR-0002-v1-required-endpoints.md)

- `GET /healthz` → `ok`
- `POST /v1/notes:ingest`
  - request: `{title, body, summary?, source_type?, origin?, source_trust?, sensitivity?, tags?}`
  - response: `{note: Note}`
- `POST /v1/notes:search`
  - request: `{query, top_k?}`
  - response: `{notes: Note[]}`
- `GET /v1/notes/{id}` → `Note`
- `POST /v1/gc:run`（SHOULD (v1.x): `mem.features.gc_short=true` 時のみ有効な実験機能）
  - request: `{target, options?}`
  - response: `{status}`

`POST /v1/gc:run` は v1 MUST には含めない。SHOULD としては、公開可否（route mount）と実行可否（flag 判定）を分離し、無効時は `NOT_FOUND`（公開 OFF）または `INTERNAL`（公開 ON かつ `mem.features.gc_short=false` 時）を返す。

クライアントは `NOT_FOUND` を「未公開（または未提供）」、`INTERNAL` を「公開中だが実行不可（flag OFF を含む）のサーバー側失敗」の運用解釈で扱う。`FAILED_PRECONDITION` / `FEATURE_DISABLED` 相当コードは v1 契約外とし、導入時は `REQ-ERR-001` の拡張手順に従う。

### 6-3-1. v1必須3エンドポイント契約（`requirements.md` × `go/api/types.go` 照合）

#### POST `/v1/notes:ingest`

**request (`NotesIngestRequest`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `title` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。v1 では削除・意味変更禁止。 |
| `body` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。v1 では削除・意味変更禁止。 |
| `summary` | `string` | 任意 | `""` | 追加バリデーションなし | 任意フィールドのため、未指定互換を維持。 |
| `source_type` | `string` | 任意 | `"manual"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定（クライアント未送信を維持）。 |
| `origin` | `string` | 任意 | `""` | trim のみ | 任意フィールドとして後方互換追加可。 |
| `source_trust` | `string` | 任意 | `"user_input"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定。 |
| `sensitivity` | `string` | 任意 | `"internal"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定。 |
| `tags` | `[]string` | 任意 | `[]` 扱い（未指定時は処理なし） | 各要素 trim、空要素は無視 | 任意配列として未指定・空配列を同等扱い。 |

**response (`NotesIngestResponse`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `note` | `Note` | 必須 | なし | 保存成功時に返却 | v1 では wrapper 構造（`{note: ...}`）を維持。 |
| `note.id` | `string` | 必須 | なし | 32 hex 形式の ID を生成 | 型・キー名固定。 |
| `note.title` | `string` | 必須 | なし | request 値（trim 後） | 型固定。 |
| `note.summary` | `string` | 必須 | `""` 許容 | request 値 | 必須キーだが空文字を許容。 |
| `note.body` | `string` | 必須 | なし | request 値（trim 後） | 型固定。 |
| `note.created_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.updated_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.last_accessed_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.access_count` | `int64` | 必須 | `0` | ingest 直後は `0` | 数値型固定。 |
| `note.source_type` | `string` | 必須 | `"manual"` 補完あり | request/補完値 | 補完込みで常に返却。 |
| `note.origin` | `string` | 必須 | `""` | request 値 | 必須キーとして維持。 |
| `note.source_trust` | `string` | 必須 | `"user_input"` 補完あり | request/補完値 | 補完込みで常に返却。 |
| `note.sensitivity` | `string` | 必須 | `"internal"` 補完あり | request/補完値 | 補完込みで常に返却。 |

#### POST `/v1/notes:search`

**request (`NotesSearchRequest`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `query` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。 |
| `top_k` | `int` | 任意 | `20`（`<=0` 時） | `<=0` は service で `20` に補正 | 任意数値として未指定互換を維持。 |

**response (`NotesSearchResponse`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `notes` | `[]Note` | 必須 | `[]` | 一致結果を最大 `top_k` 件返却 | v1 では配列 wrapper 構造を維持。 |
| `notes[].*` | `Note` 各フィールド | 必須 | `Note` 契約準拠 | `Note` と同一 | `Note` フィールドの型・キー名を維持。 |

#### GET `/v1/notes/{id}`

**request（path parameter）**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `id` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | path パラメータ必須を維持。 |

**response (`Note`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `id` | `string` | 必須 | なし | ノート存在時に返却、非存在は `NOT_FOUND` | 型・キー名固定。 |
| `title` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `summary` | `string` | 必須 | `""` 許容 | 保存済み値 | 必須キーとして維持。 |
| `body` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `created_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `updated_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `last_accessed_at` | `string` | 必須 | なし | 取得時刻で更新 | 取得時更新の挙動を維持。 |
| `access_count` | `int64` | 必須 | なし | 取得ごとに +1 | カウンタ型・加算挙動を維持。 |
| `source_type` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `origin` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `source_trust` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `sensitivity` | `string` | 必須 | なし | 保存済み値 | 型固定。 |

### 6-4. エラーモデル
- ADR: [ADR-0003: ErrorCodeとretryable設計（再試行可/不可の境界）](../../docs/ADR/ADR-0003-errorcode-retryable-boundary.md)

- Requirement ID: `REQ-ERR-001`

本節を ErrorCode 契約の正本とし、`error-contract.md` は本節の運用向け要約として同期する。

共通：

```json
{"code":"INVALID_ARGUMENT","message":"...","details":{}}
```

- `INVALID_ARGUMENT` → 400
- `NOT_FOUND` → 404
- `CONFLICT` → 409
- `GATEKEEP_DENY` → 403
- `INTERNAL` → 500

`go/service/errors` と `go/api/errors.go` の現行方針（`ErrInvalidArgument` / `ErrNotFound` を明示マッピングし、それ以外は `INTERNAL` へフォールバック）を前提に、ErrorCode を 2 段で再定義する。

#### ErrorCode 区分（v1）

| 区分 | ErrorCode | HTTP | 契約レベル | 条件 |
| --- | --- | --- | --- | --- |
| v1必須保証 | `INVALID_ARGUMENT` | 400 | MUST | 常時有効。 |
| v1必須保証 | `NOT_FOUND` | 404 | MUST | 常時有効。 |
| v1必須保証 | `INTERNAL` | 500 | MUST | 常時有効。未分類エラーのフォールバック先。 |
| v1.x拡張（feature/sentinel依存） | `CONFLICT` | 409 | SHOULD | service sentinel（例: `ErrConflict`）を実装し `mapError` へ明示追加した場合のみ返却。未実装時は `INTERNAL`。 |
| v1.x拡張（feature/sentinel依存） | `GATEKEEP_DENY` | 403 | SHOULD | gatekeeper deny 系 sentinel（例: `ErrGatekeepDeny`）を実装し `mapError` へ明示追加した場合のみ返却。未実装時は `INTERNAL`。 |

#### 現行実装との差分注記（`go/api/types.go` / `go/api/http_server.go`）

- `go/api/types.go` には `CONFLICT` / `GATEKEEP_DENY` 定数が定義済みだが、返却は service sentinel と `mapError` 実装に依存する。
- `go/api/http_server.go` は `writeErr` で `CONFLICT=409` / `GATEKEEP_DENY=403` を処理可能だが、上流が当該コードを返さない限り到達しない。
- sentinel 未実装時は `mapError` 方針に従って `INTERNAL`（500）へフォールバックする。

| 代表事象 | service 層の分類 | API `code` | HTTP | 再試行可否 | 備考 |
| --- | --- | --- | --- | --- | --- |
| 入力不備（必須欠落・形式不正・空文字） | `ErrInvalidArgument` | `INVALID_ARGUMENT` | 400 | 不可 | クライアント入力修正が必要。 |
| DB ロック（`database is locked` / busy timeout 超過） | 現行は汎用エラー（将来 `ErrConflict` 候補） | `INTERNAL`（将来 `CONFLICT` 検討可） | 500（将来 409 検討可） | 条件付き可 | 短時間で解消しうるため指数バックオフ再試行対象。長時間継続時は運用アラート。 |
| LLM タイムアウト（上流 timeout / 502/503/504） | 汎用エラー | `INTERNAL` | 500 | 条件付き可 | 一時障害として再試行対象。最大回数超過で失敗扱い。 |
| Gatekeeper deny / needs_human(v1.3 では deny 扱い) | 将来 sentinel 化（`ErrGatekeepDeny`） | `GATEKEEP_DENY`（未実装時は `INTERNAL`） | 403（未実装時は 500） | 不可 | ポリシー判定のため再試行では解消しない。入力/運用判断が必要。 |

#### エラーコード別 再試行可否表

| API `code` | HTTP | 再試行可否 | ルール |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | 400 | 不可 | 入力を修正して再実行。 |
| `NOT_FOUND` | 404 | 不可 | 対象 ID・クエリを見直す。 |
| `CONFLICT` | 409 | 条件付き可 | DB ロック等の一時競合のみ再試行。恒久競合は不可。 |
| `GATEKEEP_DENY` | 403 | 不可 | ポリシー deny は即時失敗。再試行禁止。 |
| `INTERNAL` | 500 | 条件付き可 | 原因が一時障害（DB lock / LLM timeout / upstream 502/503/504）の場合のみ指数バックオフ。恒久障害は不可。 |

### 6-5. 互換性ポリシー（v1）

- v1 系（`/v1/*`）では、既存クライアントを壊さない後方互換を維持する。

許容する変更（v1 内）：

- request/response への**任意フィールド追加**（既存必須フィールドは維持）。
- 新規エンドポイント追加（既存エンドポイントの契約は維持）。
- 既存フィールドのバリデーション強化のうち、既存の正常系入力を拒否しない変更。

禁止する変更（v1 内）：

- 既存の必須フィールド削除。
- 既存フィールドの意味変更（同じキー名で別意味にする変更）。
- 既存レスポンスの型変更（例：`string` → `object`）や、既存成功レスポンスの構造破壊。

破壊変更が必要な場合の移行手順：

1. まずは **新エンドポイント**（例：`/v1/notes:search2`）で並行提供し、既存挙動を維持する。
2. 並行提供で吸収できない場合は **`/v2` を新設**し、v1 を非推奨化する。
3. `CHANGES.md` に移行期限・差分・移行例を記載し、少なくとも 1 リリースは移行猶予を置く。

CLI 出力との互換責任範囲：

- 人間向けの通常表示（非 `--json`）は、可読性改善のための文言・並び変更を許容する。
- `--json` 出力は API の互換ポリシーと同等に扱い、v1 では破壊的変更を禁止する。
- 機械連携は API または CLI `--json` を正とし、互換責任はこの 2 系統に対して負う。