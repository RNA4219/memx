## 1. Purpose
memx-core は、LLM / Agent 向けの **外部記憶基盤** である。  
本システムは、証拠・圧縮知識・成果物参照・由来関係を保存 / 検索 / 要約 / 退避することを目的とする。

memx-core は **案件の進行管理** を行わない。  
また、**外部トラッカー同期** や **オーケストレーション** も責務に含めない。

---

## 2. Product Definition
memx-core は以下を提供する。

- Evidence の保存と参照
- Raw データの chunk 管理
- Distilled Knowledge の保存と参照
- Artifact のメタデータ登録
- Lineage の保存と参照
- Summary / Distill
- Archive / GC
- Local-first な運用前提の安全な記憶基盤

memx-core は、上位システムから参照される **Memory Substrate** として振る舞う。

---

## 3. In Scope

### 3.1 Evidence Management
memx-core は以下のような生データまたは証拠断片を保持できなければならない。

- ログ
- transcript
- 文書断片
- 外部引用
- テスト結果
- diff
- 実行結果

#### Requirements
- Evidence を一意な ID で登録できること
- source URI, source hash, recorded_at などの由来情報を保持できること
- Raw 全文を chunk 単位で保持できること
- 必要時に特定 chunk を参照できること

---

### 3.2 Knowledge Management
memx-core は、Evidence から抽出・圧縮された知識を保持できなければならない。

対象例:

- summary
- fact
- policy
- profile
- failure pattern
- domain knowledge

#### Requirements
- Knowledge を一意な ID で保存できること
- kind / scope / confidence を持てること
- valid_from / valid_to により有効期間を表現できること
- Raw Evidence と分離して扱えること

---

### 3.3 Artifact Registry
memx-core は、成果物そのものではなく **成果物参照の正本** を保持しなければならない。

対象例:

- 仕様書
- コード
- テスト
- レポート
- プロンプト束
- 出力 bundle

#### Requirements
- Artifact を一意な ID で登録できること
- URI, version, content hash を保持できること
- 実体ファイルが外部にあっても参照可能であること

---

### 3.4 Lineage Tracking
memx-core は、Evidence / Knowledge / Artifact 間の由来関係を保持できなければならない。

対象例:

- derived_from
- supports
- contradicts
- summarizes
- cites

#### Requirements
- typed ref ベースで関係を保存できること
- from_ref / to_ref / edge_type を保持できること
- 参照元 / 参照先の追跡が可能であること

---

### 3.5 Summarize / Distill
memx-core は、Raw Evidence から Distilled Knowledge を生成するための要約 / 蒸留処理を支援しなければならない。

#### Requirements
- Raw から Summary を作成できること
- Summary から Pattern / Fact を作成できること
- 蒸留結果が元 Evidence を辿れること
- Summary / Distill は上書きではなく新規生成として扱えること

---

### 3.6 Archive / GC
memx-core は、利用頻度やライフサイクルに応じて記憶を退避 / 整理できなければならない。

#### Requirements
- 非アクティブデータを archive に移せること
- GC 対象を判定できること
- 必要時に archive から再参照できること

---

## 4. Out of Scope

memx-core は以下を責務に含めない。

### 4.1 Task / Work State Management
以下は work 系システムの責務とする。

- task
- task_state
- decision
- open_question
- run
- context_bundle
- current_step
- done_when
- internal execution status

### 4.2 External Tracker Synchronization
以下は tracker bridge 系システムの責務とする。

- Jira / BTS / GitHub Issues 接続
- issue status 同期
- comment 同期
- assignee 同期
- sprint / board 情報
- remote issue cache

### 4.3 Orchestration / Planning
以下は runtime / orchestration / workflow 系の責務とする。

- 実行計画生成
- agent orchestration
- verifier / planner
- context build policy
- workflow rule execution

---

## 5. Ownership Boundaries

### memx-core が正本となるもの
- evidence
- evidence_chunk
- knowledge_card
- artifact metadata
- lineage_edge

### memx-core が参照するが正本ではないもの
- task / task_state
- decision / open_question
- tracker issue / comment / status

---

## 6. Dependency Rules

memx-core は下位基盤として振る舞わなければならない。

#### Rules
- memx-core は work system を知らないこと
- memx-core は tracker bridge を知らないこと
- repo 間参照は typed ref のみを使うこと
- DB 横断 FK を持たないこと
- 上位システムが memx-core を参照する一方向依存であること

依存方向は以下とする。

```text
work-system      -> memx-core
tracker-bridge   -> work-system
runtime          -> work-system + memx-core + tracker-bridge
````

---

## 7. Storage Model

memx-core は以下の論理ストアを持つ。

### short

一時投入・未整理記憶を格納する。

### journal

時系列イベント・履歴を格納する。

### knowledge

抽象化済み知識を格納する。

### archive

低頻度アクセスの退避記憶を格納する。

#### Requirements

* 各ストアは責務が重複しすぎないこと
* Raw / Journal / Distilled / Archived を論理的に分離できること
* 上位層が利用目的に応じてストアを選べること

---

## 8. Functional Requirements

### FR-01 Evidence Ingest

Evidence を登録できること。

### FR-02 Evidence Retrieve

Evidence とその chunk を参照できること。

### FR-03 Knowledge Ingest

Knowledge を登録できること。

### FR-04 Knowledge Search

Knowledge を検索できること。

### FR-05 Artifact Register

Artifact metadata を登録できること。

### FR-06 Lineage Register

Lineage を登録できること。

### FR-07 Lineage Trace

任意の entity から lineage を辿れること。

### FR-08 Summarize

Evidence から summary を生成できること。

### FR-09 Distill

summary / evidence から distilled knowledge を生成できること。

### FR-10 Archive / GC

archive と GC を実行できること。

---

## 9. Non-Functional Requirements

### NFR-01 Local-first

ローカル環境で完結可能であること。

### NFR-02 Deterministic Identity

Evidence / Knowledge / Artifact / Lineage は安定した ID を持つこと。

### NFR-03 Referential Safety

typed ref により参照整合性を保ちやすいこと。

### NFR-04 Extensibility

上位の work system / tracker bridge / runtime と疎結合に接続できること。

### NFR-05 Auditability

知識や要約が元証拠に遡れること。

### NFR-06 Storage Efficiency

Raw を常時展開せず、chunk / archive / distill を前提に保存効率を保てること。

---

## 10. Acceptance Criteria

memx-core は以下を満たしたとき受け入れ可能とする。

1. Evidence を保存し chunk 参照できる
2. Distilled Knowledge を保存し検索できる
3. Artifact metadata を登録できる
4. Lineage を登録し追跡できる
5. Summary / Distill により Raw から圧縮知識を生成できる
6. Task state や external tracker sync を責務に含めていない
7. 上位システムが typed ref 経由で利用できる

---

## 11. Design Principles

* memx-core は **記憶基盤** に徹する
* memx-core は **案件状態管理** を持たない
* memx-core は **外部チケット同期** を持たない
* memx-core は **上位層から読まれる下位基盤** である
* memx-core は **普段は圧縮知識を返し、必要時だけ raw に降りられる構造** を支える