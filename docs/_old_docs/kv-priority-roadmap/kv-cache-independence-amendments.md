# KV Cache Independence Requirement 追記案

作成日: 2026-03-08
目的: 現行 `kv-cache-independence.md` に対して、設計ブレを減らすための追記候補を 6 点整理する。
方針: 元文書は直接上書きせず、追記案として別紙化する。

---

## 追記案 1: typed_ref canonical format の固定

### 追加理由

`agent-taskstate` / `memx-core` / `tracker-bridge` の参照形式が揃っていないと、`source refs`、`entity_link`、`lineage`、`context rebuild` の接続点が壊れる。

### 追記先候補

- 第 4 章 解決方針
- 第 6 章 機能要求
- 第 9 章 受入条件

### 追記文案

#### 4. 解決方針 への追記

- repo / DB / service をまたぐ参照は `typed_ref` に統一する
- `typed_ref` の canonical format は 4 セグメント文字列とする
- migration 期間中は旧形式の読込を許可してもよいが、新規出力は canonical format に統一する

#### 6. 機能要求 への追加

### FR-008 typed_ref 正規化

システムは、外部状態・bundle・根拠参照で用いる `typed_ref` を canonical format で保存・返却しなければならない。

canonical format は以下とする。

```txt
<domain>:<entity_type>:<provider>:<entity_id>
```

最低限、以下を満たさなければならない。

- 新規生成される ref は canonical format であること
- migration 期間中は旧 ref を canonical format へ正規化可能であること
- 実在性確認と形式妥当性確認を分離すること

#### 9. 受入条件 への追加

### AC-006 typed_ref 一貫性

`agent-taskstate` / `memx-core` / `tracker-bridge` をまたぐ主要参照が、同一 canonical `typed_ref` 形式で保存・再利用・追跡可能であること。

---

## 追記案 2: current dashboard と state history の分離

### 追加理由

`current state` と `state history` を同じ概念で扱うと、ダッシュボード更新と履歴付き遷移の責務が混ざる。KV 独立性では「現在地」だけでなく「どうそこに至ったか」が必要である。

### 追記先候補

- 第 6 章 機能要求 FR-001 / FR-007
- 第 8 章 失敗条件
- 第 10 章 ゴールデンパス受入シナリオ

### 追記文案

#### FR-001 状態外部化 への追記

外部保持する状態は、少なくとも以下の 2 系統に分離されなければならない。

- current dashboard
- append-only history

current dashboard は最新の作業地点を示すために用い、history は状態遷移と判断経路の追跡に用いる。

#### FR-007 状態遷移明示化 の置換追記

システムは、task の進行状態を以下の 2 層で保持しなければならない。

- 最新状態を示す materialized current state
- append-only の state transition history

状態変更は暗黙更新ではなく、履歴付き遷移として記録されなければならない。
また、`task.status` の変更は専用の状態遷移処理を経由し、直接更新してはならない。

最低限、状態遷移履歴は以下を持たなければならない。

- from state
- to state
- reason
- actor
- related run ref
- changed timestamp

#### 8. 失敗条件 への追記

- current state は復元できるが state history が失われている
- task.status が直接更新され、履歴付き遷移として辿れない

---

## 追記案 3: context bundle の必須監査項目の明確化

### 追加理由

現行文書でも `source refs` と `raw included flag` はあるが、再現性と監査性のためには `purpose`, `rebuild level`, `generator version`, `diagnostics` まで固定した方がよい。

### 追記先候補

- 第 6 章 FR-006
- 第 7 章 NFR-001 / NFR-002
- 第 9 章 受入条件

### 追記文案

#### FR-006 継続用 bundle 保存 の置換追記

bundle は少なくとも以下を持たなければならない。

- purpose
- rebuild level
- summary
- state snapshot
- decision digest
- open question digest
- source refs
- raw included flag
- generator version
- generated timestamp
- diagnostics

`diagnostics` は最低限、以下を表現可能でなければならない。

- missing refs
- unsupported refs
- partial bundle flag
- resolver warning

#### NFR-001 再現性 への追記

同一状態集合と同一 generator version から再構成される bundle は、意味的に再現可能でなければならない。

#### NFR-002 可監査性 への追記

bundle 本体だけでなく、bundle 生成時に使用された `source refs`、`generator version`、`diagnostics` を後から監査可能でなければならない。

#### 9. 受入条件 への追加

### AC-007 Bundle 監査性

任意の bundle について、少なくとも `purpose`、`source refs`、`raw included flag`、`generator version`、`diagnostics` を確認できること。

---

## 追記案 4: memx-core 依存の段階化

### 追加理由

`memx-core` は必要だが、現時点の実装成熟度を考えると、最初から lineage / journal / knowledge 前提で受入条件を書くと、要件と実装が乖離しやすい。

### 追記先候補

- 第 4 章 解決方針
- 第 5 章 目標
- 第 9 章 受入条件

### 追記文案

#### 4. 解決方針 への追記

記憶基盤連携は summary-first を原則とし、段階的に導入する。

- Phase 1: task 継続に必要な構造状態を外部化する
- Phase 2: summary-first で evidence / knowledge / artifact を再構成に利用する
- Phase 3: lineage / journal / distilled knowledge により根拠追跡を強化する

#### 5. 目標 への追記

本要求の達成は段階的に評価してよい。ただし、最終的には raw evidence と distilled knowledge の双方へ辿れる構造を持たなければならない。

#### 9. 受入条件 への追加

### AC-008 段階的導入整合

初期段階では `work state` のみで再開可能であってもよいが、後続段階では summary-first の memory integration を通じて evidence / knowledge / artifact を bundle に取り込めること。

---

## 追記案 5: tracker 情報は optional input であり必須依存ではないことの明記

### 追加理由

外部 tracker は有用だが、resume の必須条件にすると KV 独立性が壊れる。tracker 側障害や未同期状態でも task は継続できる必要がある。

### 追記先候補

- 第 3 章 問題定義
- 第 6 章 FR-003
- 第 8 章 失敗条件
- 第 11 章 非対象

### 追記文案

#### FR-003 Context Rebuild への追記

tracker issue snapshot や外部 issue 情報は optional input として利用してよいが、task 再開の必須条件としてはならない。

#### 8. 失敗条件 への追記

- 外部 tracker 情報がないと task を再開できない
- 外部 issue snapshot が欠けると current state を復元できない

#### 11. 非対象 への追記

外部 tracker を内部作業状態の正本として扱うことは本要求の対象外とする。

---

## 追記案 6: stale state / stale bundle への競合制御

### 追加理由

複数エージェントやセッション断絶後の再開では、古い bundle や古い current state に基づく更新が起きやすい。これを放置すると、再開可能性より先に状態破壊が起きる。

### 追記先候補

- 第 6 章 機能要求
- 第 7 章 非機能要求
- 第 8 章 失敗条件

### 追記文案

#### 6. 機能要求 への追加

### FR-009 競合検出

システムは、stale state または stale bundle に基づく更新を検出可能でなければならない。

最低限、以下のいずれかを持たなければならない。

- state revision
- task version
- expected current state
- bundle generated_at / source snapshot version

不一致時は暗黙 merge を行わず、競合として扱わなければならない。

#### NFR-004 劣化耐性 への追記

再開時に古い state や bundle が使われた場合でも、誤更新より競合検出を優先しなければならない。

#### 8. 失敗条件 への追記

- stale bundle からの更新を検出できない
- 古い current state に基づいて superseded decision や古い next action を採用してしまう

---

## 補足

今止めて修正する優先度は以下。

1. typed_ref canonical format
2. current dashboard と state history の分離
3. context bundle の監査項目固定
4. memx-core 依存の段階化
5. tracker optional input 明記
6. stale state / stale bundle 競合制御

最初の 4 点は、今の段階で要件に入れておかないと後工程の差し戻しコストが高い。
