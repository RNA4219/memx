---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# CONTRACTS

API/CLI の機械可読契約（フィールド単位）の正本サマリ。

- API 正本: `memx_spec_v3/docs/contracts/openapi.yaml`
- CLI `--json` 正本: `memx_spec_v3/docs/contracts/cli-json.schema.json`

## API Contracts

### POST /v1/notes:ingest
#### Request
| field | type | required | constraints |
| --- | --- | --- | --- |
| title | string | yes | non-empty |
| body | string | yes | non-empty |
| summary | string | no | string |
| source_type | string | no | string |
| origin | string | no | string |
| source_trust | string | no | string |
| sensitivity | string | no | string |
| tags | array<string> | no | each item is string |

#### Response 200
| field | type | required | constraints |
| --- | --- | --- | --- |
| note | object | yes | `Note` schema |

### POST /v1/notes:search
#### Request
| field | type | required | constraints |
| --- | --- | --- | --- |
| query | string | yes | non-empty |
| top_k | integer | no | integer |

#### Response 200
| field | type | required | constraints |
| --- | --- | --- | --- |
| notes | array<object> | yes | each item is `Note` |

### GET /v1/notes/{id}
#### Request
| field | type | required | constraints |
| --- | --- | --- | --- |
| id (path) | string | yes | minLength: 1 |

#### Response 200 (`Note`)
| field | type | required | constraints |
| --- | --- | --- | --- |
| id | string | yes | non-empty |
| title | string | yes | string |
| summary | string | yes | string |
| body | string | yes | string |
| created_at | string | yes | timestamp string |
| updated_at | string | yes | timestamp string |
| last_accessed_at | string | yes | timestamp string |
| access_count | integer | yes | int64 |
| source_type | string | yes | string |
| origin | string | yes | string |
| source_trust | string | yes | string |
| sensitivity | string | yes | string |

## CLI Contracts（`--json`）
- `mem in short --json` は `NotesIngestResponse` と同型。
- `mem out search --json` は `NotesSearchResponse` と同型。
- `mem out show --json` は `Note` と同型。

## Error Contracts
| http_status | code | retryable | description |
| --- | --- | --- | --- |
| 400 | INVALID_ARGUMENT | false | 入力検証エラー |
| 404 | NOT_FOUND | false | ノート未検出 |
| 409 | CONFLICT | false | 競合（状態不整合/重複） |
| 403 | GATEKEEP_DENY | false | ポリシー拒否（fail-closed） |
| 409 | FEATURE_DISABLED | false | 機能フラグ無効 |
| 500 | INTERNAL | conditional | 一時障害時のみ再試行 |
