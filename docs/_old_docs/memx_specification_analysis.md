# memx Specification Analysis

## 1. Executive Summary

memx は、LLM エージェント向けのローカルファーストな知識/記憶基盤であり、v1 では「投入（ingest）・検索（search）・参照（show）」を API 正本契約に基づき提供することが中核責務である。4層ストア（short/journal/knowledge/archive）を前提に、GC（Observer/Reflector）と archive 補償フローは v1.x の実験機能として段階導入する。

設計文書はレイヤ責務・4DB責務・NFR運用導線まで明確化されており、主要REQのトレーサビリティ表も存在するため、基盤設計に着手可能な状態に近い。一方で、契約・トレーサビリティ記述に「将来拡張を前提とした先行参照（未定義フィールド/未導入I/F）」が混在し、設計初期での境界誤解リスクが残る。

## 2. Specification Inventory

### requirements
- `memx_spec_v3/docs/requirements.md`（要件正本）
- `workflow-cookbook/docs/requirements.md`（別系統テンプレ/サンプル系）
- `workflow-cookbook/docs/requirements.template.md`

### design
- `memx_spec_v3/docs/design.md`（設計正本）
- `workflow-cookbook/docs/design.md`（別系統）
- `workflow-cookbook/docs/design.template.md`

### interface
- `memx_spec_v3/docs/interfaces.md`

### contract
- `memx_spec_v3/docs/contracts/openapi.yaml`
- `memx_spec_v3/docs/contracts/cli-json.schema.json`
- `memx_spec_v3/docs/CONTRACTS.md`
- `memx_spec_v3/docs/error-contract.md`

### review
- `memx_spec_v3/docs/reviews/DESIGN-REVIEW-20260304-001.md`
- `memx_spec_v3/docs/reviews/DESIGN-ACCEPTANCE-20260304.md`
- `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304.md`
- `memx_spec_v3/docs/reviews/DESIGN-CHAPTER-VALIDATION-20260304-002.md`
- `memx_spec_v3/docs/reviews/README.md`
- `memx_spec_v3/docs/reviews/TEMPLATE.md`

### task
- ルート `TASK.*.md` 群（ブートストラップ、GC、章検証、優先度見直し等）
- `docs/TASKS.md`
- `workflow-cookbook/TASK.codex.md`

### other（設計前提として重要）
- `memx_spec_v3/docs/traceability.md`
- `memx_spec_v3/docs/spec.md`
- `memx_spec_v3/docs/operations-spec.md`
- `docs/ADR/ADR-0001-4db-boundary.md`
- `docs/ADR/ADR-0002-v1-required-endpoints.md`
- `docs/ADR/ADR-0003-errorcode-retryable-boundary.md`
- `memx_spec_v3/memory_policy.yaml`
- `memx_spec_v3/schema/*.sql`

## 3. System Responsibilities

### primary system goal
- ローカル単体運用で、CLI/API からノート投入・検索・参照を後方互換を維持して提供する。
- 将来の蒸留/昇格/退避（Observer/Reflector/GC）を見据えた4層記憶アーキテクチャを保持する。

### non-goals
- Web UI。
- マルチユーザー向け認証認可/監査基盤の本格提供。
- クラウド前提の分散構成や外部ベクタDB必須化。
- v1中の破壊的 API/CLI 変更。

### constraints
- CLI は API の薄いラッパであること（1:1 マッピング原則）。
- API/CLI JSON 同型維持。
- fail-closed（機密入力拒否）を維持。
- feature flag 既定 OFF で SHOULD 機能を隔離。
- 4DB 境界（short/journal/knowledge/archive）と lineage 記録前提。
- NFR（性能・RTO/RPO・補償収束・監査証跡）を運用証跡で判定。

## 4. Architecture Candidate Components

1. **CLI Adapter**
   - 責務: 入力整形、API 呼び出し、表示整形。
   - 入力: ファイル/stdin/CLI flags。
   - 出力: 人間向け表示 + `--json` 同型出力。
   - 依存: API Contract, CLI JSON Schema。

2. **HTTP API Layer**
   - 責務: 契約検証、HTTP ステータス/エラー整形、安定エンドポイント提供。
   - 入力: `/v1/notes:ingest`, `/v1/notes:search`, `/v1/notes/{id}`, `/v1/gc:run`。
   - 出力: `Note`, `Notes*Response`, `Error`。
   - 依存: OpenAPI 契約。

3. **Service / Usecase Layer**
   - 責務: ingest/search/show/gc 判定ロジックの唯一入口。
   - 入力: API DTO。
   - 出力: ドメイン結果、再試行可否付きエラー分類。
   - 依存: Gatekeeper, DB, LLM Clients, memory policy。

4. **Gatekeeper / Policy Engine**
   - 責務: `memory_store`, `memory_output`, `archive_move`, `archive_purge` 判定。
   - 入力: sensitivity, actor, approval context。
   - 出力: allow/deny（fail-closed）。
   - 依存: `GUARDRAILS.md`, セキュリティ運用規約。

5. **Storage Orchestrator (4DB)**
   - 責務: store別保存/検索/退避、lineage整合、FTS更新、short_meta更新。
   - 入力: note payload, GC candidate set, retention condition。
   - 出力: 永続化データ、監査対象イベント。
   - 依存: `schema/*.sql`, `memory_policy.yaml`。

6. **Distillation Pipeline (Observer/Reflector)**
   - 責務: shortクラスタ観測→journal化→knowledge反映。
   - 入力: short notes, tags, embeddings, policy threshold。
   - 出力: observation notes, reflected pages, lineage(`observed`/`reflected`)。
   - 依存: LLM（MiniLLM/ReflectLLM）, embedding。

7. **Archive & Compensation Manager**
   - 責務: short→archive 退避、補償再実行、削除判定の安全化。
   - 入力: GC計画、lineage状態、archive実在確認。
   - 出力: archived data, compensation log, delete-ready判定。
   - 依存: operations-spec, incident artifacts。

8. **Operations/Evidence Layer**
   - 責務: RTO/RPO 判定、waiver管理、INチケット証跡統合。
   - 入力: incident timeline, recovery logs。
   - 出力: 合否判定用証跡（`IN-*`, `incident-summary`, `recovery-log`）。
   - 依存: RUNBOOK, EVALUATION, operations-spec。

## 5. Major Data Flows

### A. memory write pipeline（ingestion）
CLI入力 → API `notes:ingest` → Service 検証 → Gatekeeper（保存可否）→ short保存（notes/tags/embeddings/fts）→ short_meta更新 → Note応答。

### B. memory distillation pipeline（observer/reflector）
GCトリガ判定（short_meta + memory_policy）→ Observerで短期ノート抽出・クラスタ化・観測ノート生成（journal）→ lineage `observed` 記録 → Reflectorでテーマ統合（knowledge更新/作成）→ lineage `reflected` 記録。

### C. memory retrieval pipeline（search/show）
- search: CLI/APIクエリ → Service検索戦略決定 → Gatekeeper検索可否 → FTS/メタ検索（必要時LLM再ランキング）→ notes配列応答。
- show: note id → 主キー参照 + Gatekeeper閲覧判定 → 単一 Note 応答。
- recall（将来）: embedding 類似検索 + 前後文脈連結 + working memory マージ。

### D. archive lifecycle（GC/補償/削除）
候補選定 → archive insert → lineage `archived_from` 記録 → short delete（条件成立時のみ）→ 失敗時は重複許容でデータ喪失回避 → 再実行で lineage/実体突合し収束。

## 6. Interface Boundaries

### CLI
- v1必須: `mem in short`, `mem out search`, `mem out show`。
- v1.x SHOULD: `mem gc short`（feature flag有効時のみ）。
- `--json` は API response 同型。

### API
- 必須: `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}`。
- 実験: `POST /v1/gc:run`（公開可否と実行可否を分離）。
- エラー: `INVALID_ARGUMENT / NOT_FOUND / CONFLICT / GATEKEEP_DENY / FEATURE_DISABLED / INTERNAL`。

### internal services
- Usecase service（業務ロジック）
- Gatekeeper（policy）
- GC orchestration（observer/reflector/archive compensation）
- Evidence collector（運用NFR判定）

### storage layers
- `short.db`, `journal.db`, `knowledge.db`, `archive.db`。
- `lineage` による遷移証跡。
- `notes_fts` 同期規則と short_meta 閾値判定。

## 7. Traceability Status

### 判定
- **主要REQは形式上トレース可能**（要件→設計→I/F→評価→契約の表が整備済み）。
- **ただし整合性リスクあり**（参照先に未定義要素を含む行が存在）。

### unmapped / weakly-mapped requirements（実質）
1. distill/recall/working/tag/meta/lineage の FUTURE 節は、設計・契約で正式導線が未固定。
2. observer/reflector の詳細は requirements に濃く、openapi 契約は `gc:run` 最小形のみで情報欠落。
3. トレーサビリティ表に、契約上未定義のフィールド参照（例: `working_scope`, `lineage`）が含まれ、機械検証しにくい。

### design sections without clear requirement linkage（実質）
- 設計章内の運用手順/テンプレ移行チェックリストは REQ直結性が弱く、実装設計境界としては補助情報。

## 8. Identified Risks

1. **責務あいまい化リスク**
   - v1最小I/Fと将来機能仕様（recall/distill）が同一文書内で混在し、設計対象範囲の誤読を誘発。
2. **契約先行/不足リスク**
   - traceability が OpenAPI 未定義フィールドを参照しており、契約主導設計時に齟齬が出る。
3. **ライフサイクル競合リスク**
   - archive の「重複許容→後収束」方針は妥当だが、収束判定・再実行回数・削除許可条件の実装一貫性が崩れるとデータ整合性事故を起こす。
4. **運用依存リスク**
   - NFR達成判定が運用証跡（INチケット、opsログ）依存のため、設計段階で observability スキーマ未固定だと評価不能。
5. **feature flag 境界リスク**
   - `/v1/gc:run` の route公開可否と実行可否の二重制御が誤実装されると、404/500 契約逸脱が起きやすい。

## 9. Missing Information

1. `gc:run` の詳細レスポンス契約（phase別結果、planned_ops、lineage結果）の正本不足。
2. distillation の入出力DTO・失敗時補償契約（observer/reflector個別）のI/F定義不足。
3. 4DB横断トランザクション不在時の一貫性モデル（idempotency key、再実行識別子）定義不足。
4. retrieval policy（FTS vs semantic recall の切替条件、ranking統合規則）の最終正本不足。
5. 運用証跡アーティファクト（JSON/NDJSON）のスキーマ定義と生成責務の固定不足。

## 10. Design Readiness Verdict

**NEEDS_SPEC_REFINEMENT**

理由：
- v1中核（ingest/search/show）の設計着手には十分な材料がある。
- ただし、memory lifecycle 全体（distillation/retrieval高度化/archive補償）の契約粒度と traceability 一貫性に不足があり、先行して「契約の空欄埋め」と「FUTURE節の設計対象切り分け」を実施すべき。
