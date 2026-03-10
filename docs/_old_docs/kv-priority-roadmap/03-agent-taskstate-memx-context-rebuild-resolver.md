# Priority 4: agent-taskstate x memx-core の context rebuild resolver

## 目的

`agent-taskstate` の bundle build を「単なる ref の寄せ集め」から、「再開に必要な文脈を再構成する処理」へ進化させる。

そのために、`agent-taskstate` から `memx-core` の evidence / knowledge / artifact を引く resolver を定義する。

## 結論

再開時の中心は `agent-taskstate` の `ContextRebuildService` とし、`memx-core` はその下で使う memory substrate として扱う。

依存方向:

```txt
agent-taskstate -> memx-core
agent-taskstate -> tracker snapshot provider
runtime -> agent-taskstate
```

`memx-core` 自体に task 概念を入れない。

## 現状の問題

- `agent-taskstate` の build は task/state/decision/question を集めるところまで
- `memx-core` の思想は合っているが、公開面はまだ short 中心
- lineage / journal / knowledge を前提にしすぎると、今すぐ動かない

## 進め方

### 方針

- interface は先に固める
- 実装は段階導入する
- 最初は `summary-first`, `selected-raw`, `diagnostics` の 3 点を成立させる

## Resolver の責務

### 1. typed_ref 解決

入力:

- `memx:evidence:local:*`
- `memx:artifact:local:*`
- `memx:knowledge:local:*`

出力:

- `resolved`
- `unresolved`
- `unsupported`

### 2. summary-first retrieval

通常は以下を優先する。

- knowledge
- artifact metadata
- evidence summary

raw は必要時のみ selected inclusion とする。

### 3. raw descent

以下のとき raw を取る。

- review 前
- conflicting summaries
- high priority open question
- low confidence decision
- investigation / verification ステップ
- operator 明示要求

### 4. diagnostics

bundle build が不完全でも落としきらない。

最低限返すもの:

- `missing_refs`
- `unsupported_refs`
- `resolver_warnings`
- `partial_bundle`

## 実装フェーズ

### Phase 3A: interface 固定

`agent-taskstate` 側に以下の抽象 I/F を置く。

```txt
ResolveRef(ref) -> ResolvedRef
ResolveMany(refs) -> ResolveReport
LoadSummary(ref) -> SummaryPayload
LoadSelectedRaw(ref, selector) -> RawPayload
```

### Phase 3B: current memx-core adapter

現時点の `memx-core` 実装に合わせて、まずは以下だけ使う。

- ingest 済み short note の参照
- search / show による evidence 相当の取得
- artifact metadata が未薄い場合は unresolved を返す

この段階では「理想的な memx 連携」ではなく、「今ある surface で bundle を再構成できること」を優先する。

### Phase 3C: knowledge / lineage 拡張

`memx-core` 側で以下が増えたら resolver を拡張する。

- knowledge search
- artifact registry
- lineage trace
- journal / knowledge

## Context Bundle に追加すべき内容

- `purpose`
- `summary`
- `decision_digest`
- `open_question_digest`
- `artifact_refs`
- `evidence_refs`
- `source_refs`
- `raw_included`
- `resolver_diagnostics`

## 受入条件

- 会話履歴ゼロで `agent-taskstate context build` を再実行できる
- `memx` ref があれば summary-first で bundle に反映される
- raw が必要な局面だけ selected raw を含められる
- ref 解決不能でも bundle build 自体は継続できる
- bundle から source refs を辿って根拠に戻れる

## 完了の定義

- `agent-taskstate` の rebuild が本当に `context rebuild` になる
- `memx-core` が「保存箱」ではなく「再開用 memory substrate」として効き始める

## 依存

- P1（`kv-cache-independence-amendments.md`）完了が必須
- P2（`01-typed-ref-unification.md`）完了が必須
- P3（`02-agent-taskstate-state-history-and-bundle-audit.md`）の bundle/source 構造が先に必要
