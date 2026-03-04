---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# CONTRACTS

本書は**正本スキーマへの索引・抜粋**を提供する。フィールド定義そのものの正本は保持しない。

- API 正本: `memx_spec_v3/docs/contracts/openapi.yaml`
- CLI `--json` 正本: `memx_spec_v3/docs/contracts/cli-json.schema.json`

## 1. 役割定義

- `interfaces.md`
  - 人間可読の I/O 説明と互換方針を示す補助仕様。
- `CONTRACTS.md`（本書）
  - 正本スキーマへの索引・抜粋のみを扱う。
  - 重複定義（フィールド型・required・制約の再定義）を持たない。

## 2. API 契約索引（openapi.yaml）

- `POST /v1/notes:ingest`
  - Request: `#/components/schemas/NotesIngestRequest`
  - Response 200: `#/components/schemas/NotesIngestResponse`
- `POST /v1/notes:search`
  - Request: `#/components/schemas/NotesSearchRequest`
  - Response 200: `#/components/schemas/NotesSearchResponse`
- `GET /v1/notes/{id}`
  - Response 200: `#/components/schemas/Note`

## 3. CLI JSON 契約索引（cli-json.schema.json）

- `mem in short --json`
  - 正本: `#/definitions/NotesIngestResponse`
- `mem out search --json`
  - 正本: `#/definitions/NotesSearchResponse`
- `mem out show --json`
  - 正本: `#/definitions/Note`

## 4. エラー契約索引

- HTTP/API エラー正本: `contracts/openapi.yaml` の `components.responses` / `components.schemas.Error*`
- 運用向け要約: `error-contract.md`

## 5. レビュー時の必須確認

- 契約差分がある PR は、`operations-spec.md` の「契約差分チェック手順」を実施し、
  `openapi.yaml` / `cli-json.schema.json` を起点にレビューする。
