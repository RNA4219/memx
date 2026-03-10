## `scope.md`

````md
# memx-core Scope

## 1. Scope Statement
memx-core は、LLM / Agent 向けの **外部記憶基盤** である。  
本システムは、Evidence・Knowledge・Artifact metadata・Lineage を保存 / 取得 / 検索 / 要約 / 退避する責務を持つ。

memx-core は **案件状態管理**、**外部トラッカー同期**、**実行計画生成** を責務に含めない。

---

## 2. In Scope

### 2.1 Evidence Storage
memx-core は、生データまたは証拠断片を保存できる。

対象例:
- ログ
- transcript
- 文書断片
- 外部引用
- テスト結果
- diff
- 実行結果

含まれる機能:
- evidence 登録
- evidence chunk 登録
- evidence 参照
- chunk 単位参照
- source metadata 保持

---

### 2.2 Knowledge Storage
memx-core は、Evidence から抽出・圧縮された知識を保存できる。

対象例:
- summary
- fact
- policy
- profile
- failure pattern
- domain knowledge

含まれる機能:
- knowledge 登録
- knowledge 更新
- knowledge 検索
- confidence / validity 管理

---

### 2.3 Artifact Registry
memx-core は、成果物のメタデータを保持できる。

対象例:
- 仕様書
- コード
- テスト
- レポート
- bundle
- prompt set

含まれる機能:
- artifact metadata 登録
- version / hash / uri 保持
- artifact 参照

---

### 2.4 Lineage Tracking
memx-core は、Evidence / Knowledge / Artifact 間の由来関係を保持できる。

対象例:
- derived_from
- supports
- contradicts
- summarizes
- cites

含まれる機能:
- lineage edge 登録
- lineage edge 検索
- 上流 / 下流トレース

---

### 2.5 Summarize / Distill
memx-core は、Raw Evidence から Summary / Distilled Knowledge を生成する処理を支援する。

含まれる機能:
- summarize request の受付
- distill request の受付
- 元 evidence と生成結果の lineage 保持
- 蒸留結果の登録

---

### 2.6 Archive / GC
memx-core は、利用頻度やライフサイクルに応じて記憶を退避 / 整理できる。

含まれる機能:
- archive 移送
- GC 候補判定
- archive 参照
- retention policy 適用

---

### 2.7 Retrieval
memx-core は、上位システムからの参照要求に応じて relevant な記憶を返せる。

含まれる機能:
- id 参照
- store 指定参照
- kind 指定検索
- typed ref 解決
- summary 優先取得
- raw evidence 参照

---

## 3. Out of Scope

### 3.1 Task / Work Management
memx-core は以下を扱わない。

- task
- task_state
- decision
- open_question
- run
- checkpoint
- context_bundle
- current_step
- done_when
- workflow status

これらは work system の責務とする。

---

### 3.2 External Tracker Synchronization
memx-core は以下を扱わない。

- Jira / GitHub Issues / Backlog / Redmine 接続
- issue status 同期
- comment 同期
- assignee 同期
- sprint / board 反映
- remote issue cache

これらは tracker bridge の責務とする。

---

### 3.3 Orchestration / Planning
memx-core は以下を扱わない。

- execution planning
- agent orchestration
- verifier / planner
- runtime scheduling
- context build policy
- workflow rule execution

これらは runtime / orchestration 層の責務とする。

---

## 4. Ownership Boundaries

### memx-core が正本
- evidence
- evidence_chunk
- knowledge_card
- artifact metadata
- lineage_edge

### memx-core が参照するが正本ではないもの
- work task state
- decisions
- open questions
- tracker issue / comment / status

---

## 5. Dependency Direction

memx-core は下位基盤として振る舞う。

依存方向:

```text
work-system      -> memx-core
tracker-bridge   -> work-system
runtime          -> work-system + memx-core + tracker-bridge
````

制約:

* memx-core は work-system を import しない
* memx-core は tracker-bridge を import しない
* 参照は typed ref を使う
* DB 横断 FK は持たない

---

## 6. Storage Layers

### short

一時投入・未整理記憶を格納する。

### journal

時系列イベント・履歴を格納する。

### knowledge

抽象化済み知識を格納する。

### archive

低頻度アクセスの退避記憶を格納する。

---

## 7. Success Criteria

memx-core は以下を満たすこと。

1. Evidence を登録し chunk 単位で参照できる
2. Knowledge を保存し検索できる
3. Artifact metadata を登録できる
4. Lineage を辿れる
5. Summary / Distill 結果が元 evidence に遡れる
6. Task state と tracker sync を責務に含めていない
7. 上位システムから typed ref で安全に参照できる