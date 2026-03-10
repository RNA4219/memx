# Priority 2: typed_ref 統一

## 目的

`agent-taskstate` / `memx-core` / `tracker-bridge` の参照表現を 1 つに固定する。

KV Cache Independence の観点では、`source_refs`、`entity_link`、`lineage`、`context_bundle_source` が同じ文字列表現で扱えないと、再開・監査・根拠追跡が途中で壊れる。

## 結論

canonical format は以下で固定する。

```txt
<domain>:<entity_type>:<provider>:<entity_id>
```

例:

- `agent-taskstate:task:local:task_01J...`
- `agent-taskstate:decision:local:dec_01J...`
- `memx:evidence:local:ev_01J...`
- `memx:artifact:local:art_01J...`
- `tracker:issue:jira:PROJ-123`
- `tracker:issue:github:owner/repo#123`

## この形式を採る理由

- `tracker` は provider を外せない
- `agent-taskstate` 側の契約文書は既に 4 セグメント前提
- `local` と `external provider` を同じ規則で扱える
- 将来 `resolver` や `bundle source` の集約時に条件分岐を減らせる

## 現状の問題

- `agent-taskstate` の契約文書は 4 セグメント
- `agent-taskstate` の現実装は 3 セグメント
- `memx-core` の実装は `memx:<type>:<id>` 固定
- `tracker-bridge` は 3 セグメント生成かつ `len(parts) >= 3` の緩い検証

この状態では、同じ ref に見えても repo ごとに別物になる。

## スコープ

- canonical format の決定
- 3 repo の parser / formatter / validator の統一
- 互換期間の migration ルール整備
- テストケースの共通化

## 非対象

- 各 provider API への解決処理
- ref 先エンティティの存在確認
- URL や URI 形式の標準化

## 実施内容

### 1. 契約を 4 セグメントへ統一

- `agent-taskstate` の docs を正本として残す
- `memx-core` docs / code を 4 セグメントへ更新する
- `tracker-bridge` docs / code を 4 セグメントへ更新する

### 2. 移行期は read-both / write-one

- parser は一時的に 3 セグメントも受理してよい
- formatter / emit は必ず 4 セグメントを出す
- DB に古い 3 セグメントが残っていても直ちに破壊しない

### 3. validator を repo ごとに独自実装しない

- 共通ルール:
  - 4 セグメントに split 可能
  - `domain`, `entity_type`, `provider`, `entity_id` が空でない
  - `domain` は既知 namespace
  - 実在性確認は別責務

### 4. canonicalization を追加する

- `memx:evidence:01H...` を入力したら、一時互換として `memx:evidence:local:01H...` へ正規化可能にする
- `agent-taskstate:task:01H...` も同様に `agent-taskstate:task:local:01H...` へ正規化可能にする
- この正規化は migration 用であり、最終的な出力は常に canonical

## 変更対象

- `agent-taskstate/docs/contracts/typed-ref.md`
- `agent-taskstate/docs/agent-taskstate_cli.py`
- `memx-core/memx_spec_v3/go/api/typed_ref.go`
- `memx-core` の interface / requirements 周辺 docs
- `tracker-bridge-materials/tracker_bridge/refs.py`
- `tracker-bridge-materials` の requirements / sqlite spec

## 受入条件

- 3 repo が同じ canonical format を出力する
- 3 repo が同じ validation rule を持つ
- `agent-taskstate` bundle source / `memx` lineage / `tracker-bridge` entity_link で同じ ref を共有できる
- 既存 3 セグメント ref を migration 期間中に読み込める

## 完了の定義

- typed_ref に関する仕様差分がなくなる
- 4 セグメント canonical format が docs と code の両方で一致する
- 以後の P3-P5 がこの契約の上で進められる

## 依存

- P1（`kv-cache-independence-amendments.md`）完了後に着手する
