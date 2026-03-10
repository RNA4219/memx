# Priority 3: agent-taskstate の state history と bundle 監査面の実装

## 目的

`agent-taskstate` を「現在地の正本」から「再開の正本」へ引き上げる。

KV Cache Independence の要求に対して、`agent-taskstate` が最低限持つべきなのは以下。

- current state
- state history
- decision
- open question
- run
- context bundle
- source refs
- raw included flag

現状の `agent-taskstate` は、latest dashboard と軽量 bundle まではあるが、履歴付き遷移と bundle 監査の面が不足している。

## 結論

`agent-taskstate` は以下の 2 系統を明確に分ける。

1. 現在地ダッシュボード
2. append-only の履歴と再開証跡

## 推奨データモデル

### A. current dashboard

- `task`
- `task_state_current`

役割:

- 現在の status
- current_step
- constraints
- done_when
- current_summary
- artifact_refs
- evidence_refs
- context_policy

### B. append-only history

- `state_transition`
- `decision`
- `open_question`
- `run`
- `context_bundle`
- `context_bundle_source`

役割:

- 何がどう変わったか
- 何を根拠に bundle を作ったか
- raw を含めたか
- どの generator version で生成したか

## 重要な設計判断

### 1. state history は別テーブルに分ける

現在の `task_state` を状態遷移履歴に流用すると、dashboard と history の責務が混ざる。

推奨:

- `task_state_current`: 最新ダッシュボード
- `state_transition`: append-only 履歴

最低カラム:

- `id`
- `task_id`
- `from_status`
- `to_status`
- `reason`
- `actor_type`
- `actor_id`
- `run_id`
- `changed_at`

### 2. status 更新は 1 箇所に集約する

- 直接 `task.status` を更新しない
- すべて `StateTransitionService` を通す
- `task.status` は materialized current state として扱う

### 3. context bundle に監査情報を持たせる

`context_bundle` に最低限追加するもの:

- `purpose`
- `summary`
- `rebuild_level`
- `raw_included`
- `generator_version`
- `generated_at`

`context_bundle_source` に最低限追加するもの:

- `context_bundle_id`
- `typed_ref`
- `source_kind`
- `selected_raw`
- `metadata_json`

### 4. bundle 本体と source refs を分離する

bundle body に ref を埋めるだけでは監査しづらい。

推奨:

- bundle は digest と summary を持つ
- source refs は別テーブルで append する

## 実施内容

### 1. docs と implementation のズレを解消する

- `task_state` が dashboard なのか history なのかを固定する
- `docs/schema/agent-taskstate.sql` と CLI 実装のモデル差分を解消する

### 2. 状態遷移履歴を実装する

- `task set-status` 時に `state_transition` を必ず追加
- 例外遷移は reason 必須
- `run_id` がある場合は遷移と紐づける

### 3. bundle 監査面を実装する

- `purpose`
- `summary`
- `rebuild_level`
- `raw_included`
- `generator_version`
- `source_refs`

を保存する

### 4. `context build` の入出力を強化する

今の `included_*_refs` だけで終わらせず、以下を持つ。

- `decision_digest`
- `open_question_digest`
- `tracker_refs`
- `source_refs`
- `diagnostics`

## 受入条件

- `agent-taskstate` 単体で current state と state history を両方取得できる
- どの status 変更も履歴として辿れる
- 最新 bundle に対して source refs を列挙できる
- bundle ごとに raw included の有無を確認できる
- generator version を保存して比較できる

## 完了の定義

- `agent-taskstate` が「現在地だけを持つツール」ではなく「再開証跡を持つツール」になる
- `kv-cache-independence.md` の FR-006 / FR-007 に直接対応できる

## 依存

- P1（`kv-cache-independence-amendments.md`）完了後に着手する
- P2（`01-typed-ref-unification.md`）完了後に着手する
