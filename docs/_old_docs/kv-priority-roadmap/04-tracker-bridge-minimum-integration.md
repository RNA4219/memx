# Priority 5: tracker-bridge の最小統合

## 目的

外部 issue を「外界の窓口」として `agent-taskstate` / `memx-core` に接続する。

ただし tracker を正本にしない。`tracker-bridge` は最後まで補助層に留める。

## 結論

P4 のゴールは full sync ではない。

最小で成立すべきなのは以下。

1. issue を取得できる
2. `agent-taskstate task` とリンクできる
3. 重要更新を `sync_event` に残せる
4. 必要な snapshot だけを `agent-taskstate context build` に渡せる
5. `done` / `review` などを外部へ短文投影できる

## なぜ優先度 4 なのか

- KV 問題の本丸は `agent-taskstate` と `memx-core`
- tracker は再開の補助情報であって、内部状態の正本ではない
- ここを先に厚くすると、外部 issue が内部思考を侵食する

## 守るべき境界

### tracker-bridge が持つもの

- `tracker_connection`
- `issue_cache`
- `entity_link`
- `sync_event`

### tracker-bridge が持たないもの

- task state の正本
- decision の正本
- evidence / knowledge の正本
- context build policy

## MVP スコープ

### 1. inbound

- issue fetch
- issue normalize
- issue_cache 更新
- sync_event 記録

### 2. linking

- `tracker:issue:*:*` と `agent-taskstate:task:local:*` の link
- role:
  - `primary`
  - `related`
  - `duplicate`
  - `blocks`

### 3. outbound

- status 反映
- short comment 反映

### 4. snapshot export

`agent-taskstate context build` が必要時に使える最小 snapshot を返す。

例:

- remote key
- title
- status
- assignee
- updated_at
- last sync result

## 先にやるべきこと

### 1. adapter を最低限完成させる

現状の Jira adapter は fetch 未実装なので、まずは 1 provider を通す。

推奨順:

1. Jira
2. GitHub Issues
3. Linear

### 2. typed_ref を P1 の canonical に合わせる

例:

- `tracker:issue:jira:PROJ-123`
- `tracker:issue:github:owner/repo#123`

### 3. agent-taskstate 連携は suggestion ベースに留める

tracker 更新を見て、いきなり `agent-taskstate task.status` を書き換えない。

最初は以下に留める。

- 更新候補の生成
- operator 確認
- run / sync_event 記録

## 受入条件

- 外部 issue を 1 件 fetch して `issue_cache` に保存できる
- issue と `agent-taskstate task` を entity_link できる
- inbound / outbound の sync_event を追跡できる
- `agent-taskstate context build` が tracker snapshot を optional input として使える
- tracker を見失っても `agent-taskstate` と `memx-core` だけで task 継続は可能

## 完了の定義

- tracker が内部状態の正本にならずに済む
- 外界同期を足しても KV Cache Independence の原則を壊さない

## 依存

- P1（`kv-cache-independence-amendments.md`）完了が必須
- P2（`01-typed-ref-unification.md`）完了が必須
- P3（`02-agent-taskstate-state-history-and-bundle-audit.md`）完了後に進めるのが安全
- P4（`03-agent-taskstate-memx-context-rebuild-resolver.md`）完了後に進めるのが安全
