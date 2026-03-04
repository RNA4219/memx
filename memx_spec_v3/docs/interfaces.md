---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx インターフェース仕様（interfaces）

## 1. CLI I/O（v1 必須）

### 1.1 `mem in short`（IF-CLI-INGEST-REQ/RES）

#### Input（IF-CLI-INGEST-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `title` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `body` (`--stdin`/`--file`) | string | yes | none | non-empty UTF-8 text | 必須削除禁止・意味変更禁止 |
| `summary` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `source_type` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `origin` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `source_trust` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `sensitivity` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `tags` | array<string> | no | `[]` 相当（未指定可） | 各要素 string | 追加のみ許可（既存意味維持） |
| `--json` | boolean(flag) | no | false | flag | true 時は API response と同型を維持 |

#### Output（IF-CLI-INGEST-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `note` | object (`Note`) | yes | none | `Note` スキーマ準拠 | トップレベル構造変更禁止 |
| `note.id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `note.title` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.summary` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.body` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.created_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `note.updated_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `note.last_accessed_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `note.access_count` | integer | yes | none | integer | 必須削除禁止・意味変更禁止 |
| `note.source_type` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.origin` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.source_trust` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.sensitivity` | string | yes | none | string | 必須削除禁止・意味変更禁止 |

### 1.2 `mem out search`（IF-CLI-SEARCH-REQ/RES）

#### Input（IF-CLI-SEARCH-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `query` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `top_k` | integer | no | implementation default | integer | 追加のみ許可（既存意味維持） |
| `--json` | boolean(flag) | no | false | flag | true 時は API response と同型を維持 |

#### Output（IF-CLI-SEARCH-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `notes` | array<`Note`> | yes | none | 各要素 `Note` 準拠 | トップレベル構造変更禁止 |

### 1.3 `mem out show`（IF-CLI-SHOW-REQ/RES）

#### Input（IF-CLI-SHOW-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `--json` | boolean(flag) | no | false | flag | true 時は API response と同型を維持 |

#### Output（IF-CLI-SHOW-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `title` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `summary` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `body` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `created_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `updated_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `last_accessed_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `access_count` | integer | yes | none | integer | 必須削除禁止・意味変更禁止 |
| `source_type` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `origin` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `source_trust` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `sensitivity` | string | yes | none | string | 必須削除禁止・意味変更禁止 |

## 2. API I/O（v1 必須）

### 2.1 `POST /v1/notes:ingest`（IF-API-INGEST-REQ/RES）

#### Request（IF-API-INGEST-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `title` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `body` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `summary` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `source_type` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `origin` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `source_trust` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `sensitivity` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `tags` | array<string> | no | omitted | 各要素 string | 追加のみ許可（既存意味維持） |

#### Response 200（IF-API-INGEST-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `note` | object (`Note`) | yes | none | `Note` 準拠 | トップレベル構造変更禁止 |

### 2.2 `POST /v1/notes:search`（IF-API-SEARCH-REQ/RES）

#### Request（IF-API-SEARCH-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `query` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `top_k` | integer | no | omitted | integer | 追加のみ許可（既存意味維持） |

#### Response 200（IF-API-SEARCH-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `notes` | array<`Note`> | yes | none | 各要素 `Note` 準拠 | トップレベル構造変更禁止 |

### 2.3 `GET /v1/notes/{id}`（IF-API-GET-REQ/RES）

#### Request（IF-API-GET-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` (path) | string | yes | none | `minLength: 1` | 必須削除禁止・意味変更禁止 |

#### Response 200（IF-API-GET-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `title` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `summary` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `body` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `created_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `updated_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `last_accessed_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `access_count` | integer | yes | none | integer | 必須削除禁止・意味変更禁止 |
| `source_type` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `origin` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `source_trust` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `sensitivity` | string | yes | none | string | 必須削除禁止・意味変更禁止 |

## 3. 互換ルール
- 必須フィールド削除禁止。
- 既存フィールド意味変更禁止。
- 成功レスポンストップレベル構造変更禁止。
- 破壊変更は v2+ で段階移行（互換フラグまたは新バージョン導入）。

## 4. エラー面

### 4.1 ErrorCode × HTTP × retryable × クライアント動作（IF-ERR-MATRIX）
| ErrorCode | HTTP | retryable | クライアント動作 |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | 400 | false | 入力修正後に再実行（自動再試行しない） |
| `NOT_FOUND` | 404 | false | 対象 ID/条件を見直し、必要なら ingest 後に再実行 |
| `CONFLICT` | 409 | false | 競合解消（重複・状態整合）後に再実行 |
| `GATEKEEP_DENY` | 403 | false | ポリシー変更または権限見直しまで停止 |
| `FEATURE_DISABLED` | 409 | false | feature flag 有効化まで停止 |
| `INTERNAL` | 500 | conditional | 一時障害判定時のみ指数バックオフ再試行（最大2回） |

### 4.2 整合元
- エラーコード区分・運用差分は `error-contract.md` を正本運用要約として参照。
- API 契約正本は `docs/contracts/openapi.yaml`。

## 5. 契約変更手順（更新順序固定）
1. `memx_spec_v3/docs/interfaces.md` を更新し、I/F 項目ID（`IF-*`）を付与・更新する。
2. 同一 PR で `memx_spec_v3/docs/contracts/openapi.yaml` を確認・反映する。
3. 同一 PR で `memx_spec_v3/docs/contracts/cli-json.schema.json` を確認・反映する。
4. 同一 PR で `memx_spec_v3/docs/CONTRACTS.md` を確認・反映する。
5. `RUNBOOK.md` の `trace-req-*` と I/F 項目ID対応を更新し、受け入れ差分レビュー可能状態にする。


## 6. 付録: RUNBOOK連携 I/F ID（v1運用）
| I/F 項目ID | 対象 |
| --- | --- |
| `IF-GC-SHORT-REQ` | `mem gc short` の入力（`--dry-run` 等の実行パラメータ） |
| `IF-GC-SHORT-RES` | `mem gc short` の出力（dry-run 結果 JSON / 実行結果） |
