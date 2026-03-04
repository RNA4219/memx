---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
priority: high
---

# memx 要件トレーサビリティ（traceability）

本書は `memx_spec_v3/docs/requirements.md#task-seed-source-fixed` を起点に、主要 REQ-ID の設計・I/F・評価・契約の対応を 1 行 1 要件で固定化する。

## 1. 主要 REQ-ID トレーサビリティ表

| Requirement ID | Source | Design Mapping | Interface Mapping | Evaluation Mapping | Contract Mapping |
| --- | --- | --- | --- | --- | --- |
| `REQ-CLI-001` | `memx_spec_v3/docs/requirements.md#3-cli-要件` | `memx_spec_v3/docs/design.md#4-1-ingest` | `memx_spec_v3/docs/interfaces.md#1-cli-iov1-必須` | `EVALUATION.md#req-cli-001-passfail` | `memx_spec_v3/docs/contracts/cli-json.schema.json#/x-requirement-id`, `memx_spec_v3/docs/contracts/cli-json.schema.json#/definitions/NotesIngestResponse`, `memx_spec_v3/docs/contracts/cli-json.schema.json#/definitions/NotesSearchResponse` |
| `REQ-API-001` | `memx_spec_v3/docs/requirements.md#6-api-要件v13-追加` | `memx_spec_v3/docs/design.md#4-1-ingest` | `memx_spec_v3/docs/interfaces.md#2-api-iov1-必須` | `EVALUATION.md#req-api-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:ingest`, `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:search`, `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes~1{id}`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/Note` |
| `REQ-GC-001` | `memx_spec_v3/docs/requirements.md#3-5-mem-gc-shortobserver--reflector` | `memx_spec_v3/docs/design.md#4-4-gc-dry-run` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#req-gc-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1gc:run`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunResponse` |
| `REQ-SEC-001` | `memx_spec_v3/docs/requirements.md#2-7-security--retention-requirements` | `memx_spec_v3/docs/design.md#2-2-securityretention-設計` | `memx_spec_v3/docs/interfaces.md#1-1-mem-in-shortif-cli-ingest-reqres` | `EVALUATION.md#req-sec-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesIngestRequest/properties/sensitivity`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/Note/properties/sensitivity`, `memx_spec_v3/docs/contracts/cli-json.schema.json#/definitions/Note/properties/sensitivity` |
| `REQ-RET-001` | `memx_spec_v3/docs/requirements.md#2-7-security--retention-requirements` | `memx_spec_v3/docs/design.md#2-2-securityretention-設計` | `memx_spec_v3/docs/interfaces.md#5-契約変更手順更新順序固定` | `EVALUATION.md#req-ret-001-passfail-waiver` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1gc:run` |
| `REQ-SEC-AUD-001` | `memx_spec_v3/docs/requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5` | `memx_spec_v3/docs/design.md#2-2-securityretention-設計` | `memx_spec_v3/docs/interfaces.md#5-契約変更手順更新順序固定` | `EVALUATION.md#req-ret-001-passfail-waiver` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunResponse` |
| `REQ-SEC-AUD-002` | `memx_spec_v3/docs/requirements.md#2-7-2-actor--approval--audit-責任分界表2-7-12-7-5` | `memx_spec_v3/docs/design.md#2-2-securityretention-設計` | `memx_spec_v3/docs/interfaces.md#5-契約変更手順更新順序固定` | `EVALUATION.md#req-ret-001-passfail-waiver` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunResponse` |
| `REQ-SEC-GRD-001` | `memx_spec_v3/docs/requirements.md#2-7-5-guardrails-fail-closed-との整合チェック要件` | `memx_spec_v3/docs/design.md#2-2-securityretention-設計` | `memx_spec_v3/docs/interfaces.md#4-エラー面` | `EVALUATION.md#v1-受け入れ基準release-scope-matrix-準拠` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/responses/GatekeepDenyError`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/ErrorCode` |
| `REQ-ERR-001` | `memx_spec_v3/docs/requirements.md#6-4-エラーモデル` | `memx_spec_v3/docs/design.md#4-4-gc-dry-run` | `memx_spec_v3/docs/interfaces.md#4-1-errorcode--http--retryable--クライアント動作if-err-matrix` | `EVALUATION.md#req-err-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/ErrorCode`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/Error`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/responses/InvalidArgumentError`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/responses/NotFoundError`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/responses/InternalError` |
| `REQ-NFR-001` | `memx_spec_v3/docs/requirements.md#5-1-性能目標v1必須3エンドポイント` | `memx_spec_v3/docs/design.md#6-1-nfr設計性能--復旧--整合性回復` | `memx_spec_v3/docs/interfaces.md#5-契約変更手順更新順序固定` | `EVALUATION.md#req-nfr-001-passfail-waiver` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:ingest`, `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:search`, `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes~1{id}` |

## 2. 運用ルール

- マッピング記法は `path#Section`（ドキュメント）または `path#/json-pointer`（契約）に統一する。
- `requirements.md` の主要 REQ-ID 追加/更新時は、本表を同一 PR で更新する。
- `spec.md` の更新順序に従い、`requirements.md` 更新直後に本書を更新する。
