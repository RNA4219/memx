---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
priority: high
---

# memx 要件トレーサビリティ（traceability）

本書は `memx_spec_v3/docs/requirements.md#task-seed-source-fixed` を起点に、REQ-ID の設計・I/F・評価・契約の対応を 1 行 1 要件で固定化する。

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
| `REQ-NFR-002` | `memx_spec_v3/docs/requirements.md#5-2-可用性復旧整合性回復運用nfr` | `memx_spec_v3/docs/design.md#6-1-nfr設計性能--復旧--整合性回復`, `memx_spec_v3/docs/operations-spec.md#3-rtorpo-判定req-nfr-002` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#運用nfr可用性復旧整合性回復合否基準` | `artifacts/ops/incident-summary.json`, `docs/IN-<実日付>-<連番>.md` |
| `REQ-NFR-003` | `memx_spec_v3/docs/requirements.md#5-2-可用性復旧整合性回復運用nfr` | `memx_spec_v3/docs/design.md#6-1-nfr設計性能--復旧--整合性回復`, `memx_spec_v3/docs/operations-spec.md#5-時系列フロー検知一次切り分け緩和復旧事後レビュー` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#運用nfr可用性復旧整合性回復合否基準` | `artifacts/ops/incident-summary.json`, `artifacts/ops/recovery-log.ndjson` |
| `REQ-NFR-004` | `memx_spec_v3/docs/requirements.md#5-2-可用性復旧整合性回復運用nfr`, `memx_spec_v3/docs/requirements.md#5-2-2-再試行方針運用固定` | `memx_spec_v3/docs/design.md#6-1-nfr設計性能--復旧--整合性回復`, `memx_spec_v3/docs/operations-spec.md#5-時系列フロー検知一次切り分け緩和復旧事後レビュー` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#運用nfr可用性復旧整合性回復合否基準` | `artifacts/ops/incident-summary.json`, `artifacts/ops/recovery-log.ndjson`, `docs/IN-<実日付>-<連番>.md` |
| `REQ-NFR-005` | `memx_spec_v3/docs/requirements.md#5-3-整合性回復要件archive-補償フロー` | `memx_spec_v3/docs/design.md#6-1-nfr設計性能--復旧--整合性回復`, `memx_spec_v3/docs/operations-spec.md#4-補償フロー収束条件req-nfr-005` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#運用nfr可用性復旧整合性回復合否基準` | `artifacts/ops/recovery-log.ndjson`, `artifacts/ops/incident-summary.json`, `docs/IN-<実日付>-<連番>.md` |
| `REQ-NFR-006` | `memx_spec_v3/docs/requirements.md#5-4-インシデント記録docsin-md最小監査項目`, `memx_spec_v3/docs/requirements.md#5-4-1-waiver-時の必須記録docsin-md-運用連動` | `memx_spec_v3/docs/design.md#6-2-req-nfr-006監査記録--waiver-責務境界`, `memx_spec_v3/docs/operations-spec.md#2-waiver-記録必須項目req-nfr-006`, `memx_spec_v3/docs/operations-spec.md#6-必須証跡ファイル一覧とキー定義` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#運用nfr可用性復旧整合性回復合否基準` | `docs/IN-<実日付>-<連番>.md`, `artifacts/ops/incident-summary.json`, `artifacts/ops/recovery-log.ndjson` |

## 2. ストア系 REQ-ID トレーサビリティ表（主要REQ以外）

### 2-1. short

| Requirement ID | Source | Design Mapping | Interface Mapping | Evaluation Mapping | Contract Mapping |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-SHORT-001` | `memx_spec_v3/docs/requirements.md#short-store-要求` | `memx_spec_v3/docs/design.md#2-1-store別設計詳細shortchroniclememopediaarchive` | `memx_spec_v3/docs/interfaces.md#1-1-mem-in-shortif-cli-ingest-reqres` | `EVALUATION.md#req-api-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:ingest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesIngestRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesIngestResponse` |
| `REQ-STORE-SHORT-002` | `memx_spec_v3/docs/requirements.md#short-store-要求` | `memx_spec_v3/docs/design.md#4-4-gc-dry-run` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#req-gc-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1gc:run`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunResponse` |

### 2-2. chronicle

| Requirement ID | Source | Design Mapping | Interface Mapping | Evaluation Mapping | Contract Mapping |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-CHR-001` | `memx_spec_v3/docs/requirements.md#chronicle-store-要求` | `memx_spec_v3/docs/design.md#2-1-store別設計詳細shortchroniclememopediaarchive` | `memx_spec_v3/docs/interfaces.md#2-api-iov1-必須` | `EVALUATION.md#req-api-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesIngestRequest/properties/dest_scope`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/Note/properties/working_scope` |
| `REQ-STORE-CHR-002` | `memx_spec_v3/docs/requirements.md#chronicle-store-要求` | `memx_spec_v3/docs/design.md#4-2-search` | `memx_spec_v3/docs/interfaces.md#2-api-iov1-必須` | `EVALUATION.md#req-api-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:search`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesSearchRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesSearchResponse` |

### 2-3. memopedia

| Requirement ID | Source | Design Mapping | Interface Mapping | Evaluation Mapping | Contract Mapping |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-MP-001` | `memx_spec_v3/docs/requirements.md#memopedia-store-要求` | `memx_spec_v3/docs/design.md#2-1-store別設計詳細shortchroniclememopediaarchive` | `memx_spec_v3/docs/interfaces.md#2-api-iov1-必須` | `EVALUATION.md#req-api-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/NotesIngestRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/Note` |
| `REQ-STORE-MP-002` | `memx_spec_v3/docs/requirements.md#memopedia-store-要求` | `memx_spec_v3/docs/design.md#2-1-store別設計詳細shortchroniclememopediaarchive` | `memx_spec_v3/docs/interfaces.md#2-api-iov1-必須` | `EVALUATION.md#req-err-001-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/responses/NotFoundError`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/responses/InternalError`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/ErrorCode` |

### 2-4. archive

| Requirement ID | Source | Design Mapping | Interface Mapping | Evaluation Mapping | Contract Mapping |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-ARC-001` | `memx_spec_v3/docs/requirements.md#archive-store-要求` | `memx_spec_v3/docs/design.md#2-1-store別設計詳細shortchroniclememopediaarchive` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#req-nfr-005-passfail` | `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunResponse`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/Note/properties/lineage` |
| `REQ-STORE-ARC-002` | `memx_spec_v3/docs/requirements.md#archive-store-要求` | `memx_spec_v3/docs/design.md#2-2-securityretention-設計` | `memx_spec_v3/docs/interfaces.md#6-付録-runbook連携-if-idv1運用` | `EVALUATION.md#req-ret-001-passfail-waiver` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1gc:run`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunRequest`, `memx_spec_v3/docs/contracts/openapi.yaml#/components/schemas/GCRunResponse` |

## 3. 運用ルール

- マッピング記法は `path#Section`（ドキュメント）または `path#/json-pointer`（契約）に統一する。
- `requirements.md` の REQ-ID（主要REQ/主要REQ以外を含む）を追加・変更した場合は、本書を同一 PR で必ず更新する。
- `spec.md` の更新順序に従い、`requirements.md` 更新直後に本書を更新する。
