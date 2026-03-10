# cookbook-resolver 要件定義

## 1. 文書情報

- 文書名: cookbook-resolver 要件定義
- 文書種別: requirements
- 版: v0.1
- 作成日: 2026-03-10
- 状態: Draft

## 2. 背景

現行の `workflow-cookbook` は、運用ガイド・設計原則・評価観点・ガードレールを Markdown 群として整理した、非常に有用なプロトタイプである。

一方で、生成AIエージェントの実運用においては、以下の問題がある。

- 参照すべきファイルを読まずに進行する
- 読んでも更新差分を追従できない
- task や Birdseye が更新対象であるにもかかわらず、契約として保持できていない
- 会話コンテキストや都度の md 参照に依存し、継続運用時に破綻しやすい
- 長文資料を毎回食わせる運用は、トークンコストと認識ブレの両面で不利である

このため、単なる Markdown 参照ではなく、機能や task から読むべき文書と必要 chunk を解決し、その参照結果を契約として残せる仕組みが必要である。

## 3. 目的

本システムの目的は、`workflow-cookbook`、要件定義、仕様書、Birdseye 系資料を、人間向け Markdown のまま保持しつつ、エージェントに対しては以下を可能にすることである。

- 読むべき文書を機能名または task から解決する
- 文書全文ではなく必要な chunk を返す
- 読んだ文書と版を task と結びつけて記録する
- 文書更新時に stale を検知する
- 既存リポジトリ群の責務を崩さずに追加可能である

## 4. 解決したい問題

### 4.1 現在の問題

- エージェントは参照対象を保証付きで読まない
- 更新に対する expected version 契約がない
- 既読記録が保持されず、前回何を見たかが残らない
- Birdseye と task が静的文書のように扱われ、状態オブジェクトとして機能していない
- `workflow-cookbook` は便利だが、読むことが努力目標であり、拘束力がない

### 4.2 本質的な課題

本質的な課題は「長文をどう分割するか」ではない。  
「何を読ませるべきかを契約として解決し、その結果を継続可能な状態として保持できるか」である。

## 5. 対象読者

- エージェント基盤を設計・実装する開発者
- `workflow-cookbook` を参照して実装を行うAIエージェント
- task / spec / Birdseye の追従性を管理したい運用者

## 6. スコープ

### 6.1 本要件で実装対象とするもの

- Markdown 文書を登録する仕組み
- 文書を chunk 化して保持する仕組み
- feature 名、task 名、topic から文書を解決する API
- 必要 chunk を取得する API
- 読了記録を task 側に残すための連携面
- 文書 version の差分に基づく stale 判定
- Skill から呼びやすい薄い API surface

### 6.2 本要件で実装対象外とするもの

- 文書本文の高度な編集 UI
- 汎用 CMS 化
- 完全自動の意味差分判定
- 外部 SaaS への本格同期ロジック
- 全文書のリアルタイム自動整形
- エージェント本体の orchestration 全体

## 7. 前提条件

本要件は、以下の既存資産を前提とする。

- `memx-core`
  - 検索・記憶・参照補助を担う
  - API と Skill の追加先候補
- `agent-taskstate`
  - task / state / decision / question / run の正本
  - read receipt や stale 判定結果の格納先
- `tracker-bridge-materials`
  - 外部 tracker や将来の Birdseye view 連携の境界
- `workflow-cookbook`
  - 人間向けの原本資料群

## 8. 配置方針

本機能は、公開済みの `memx-core` 本体へ直接混入させるのではなく、以下のいずれかで実装する。

- `memx-core` の検証用ブランチ
- `memx-core` の private fork
- `memx-*` 系の別リポジトリ
- `cookbook-resolver` のような隣接リポジトリ

※ 本家 `memx-core` の責務を不必要に肥大化させないことを優先する。

## 9. システム概念

本システムは、文書世界を契約世界へ接続する薄い層である。

役割分担は以下とする。

- `workflow-cookbook`
  - 人間向け原本
- resolver
  - feature / task から参照文書と chunk を解決
- `memx-core`
  - 文書・chunk の検索と取得
- `agent-taskstate`
  - read receipt、dependency、expected version、stale の正本
- `tracker-bridge-materials`
  - 将来の外部連携

## 10. 用語定義

### 10.1 Doc

管理対象となる1つの文書。  
例: requirement、spec、cookbook、Birdseye capsule

### 10.2 Chunk

Doc をエージェント参照用に分割した最小単位。  
単純分割ではなく、見出し構造や意味単位を優先する。

### 10.3 Resolve

feature 名や task 名などを入力に、読むべき doc と chunk を導出すること。

### 10.4 Read Receipt

エージェントが、どの task に対して、どの doc のどの version を読み、どの chunk を参照したかを記録する情報。

### 10.5 Stale

過去に参照した doc version と現在 version が一致せず、task の再確認が必要な状態。

## 11. ユースケース

### 11.1 実装前に必要資料を解決する

開発エージェントは feature 名を指定し、必読 doc と推奨 doc を受け取る。  
その後、必要 chunk だけ取得して読み、task に読了記録を残す。

### 11.2 task 進行中に仕様の更新を検知する

spec の version が更新された際、過去 version を参照していた task を stale として扱い、再読を要求する。

### 11.3 Birdseye と spec の関連を追跡する

Birdseye 項目と spec doc を論理参照で結び、どの feature がどの観点に依存しているかを追跡する。

## 12. 機能要件

### 12.1 文書登録

システムは Markdown 文書を doc として登録できなければならない。

最低限保持する項目は以下とする。

- doc_id
- doc_type
- title
- source_path
- version
- updated_at
- tags
- feature_keys
- summary
- body

### 12.2 chunk 生成

システムは doc を chunk に分割して保持できなければならない。

chunk 生成は以下を満たすこと。

- 見出し単位を優先する
- 長すぎる節のみ再分割する
- chunk ごとに ordinal を持つ
- chunk に importance を付与できる

### 12.3 文書解決

システムは feature 名、task_id、topic のいずれかから、読むべき doc を解決できなければならない。

解決結果には最低限以下を含むこと。

- required docs
- recommended docs
- current version
- reason
- top chunks

### 12.4 chunk 取得

システムは doc_id または chunk_id を指定して、必要部分を返せなければならない。

取得方式として最低限以下を提供する。

- doc 全文取得
- heading 指定取得
- query 指定取得
- top-k chunk 取得

### 12.5 読了記録

システムは、doc と chunk の読了結果を task に紐づけて記録できなければならない。

記録対象は以下とする。

- task_id
- doc_id
- version
- chunk_ids
- read_at
- reader

### 12.6 stale 判定

システムは、task に紐づく read receipt と現在の doc version を比較し、stale を判定できなければならない。

最低限のルールは以下とする。

- `doc_id` が同一
- `version` が不一致
- 不一致時は stale とする

### 12.7 契約解決

システムは feature または task_id から、最低限の契約情報を返せなければならない。

対象は以下とする。

- required docs
- acceptance criteria
- forbidden patterns
- definition of done
- 関連 dependency

## 13. 非機能要件

### 13.1 可搬性

- ローカル実行可能であること
- SQLite を利用できること
- 既存 CLI / API / Skill との統合が容易であること

### 13.2 拡張性

- doc_type を追加しやすいこと
- chunking 戦略を差し替え可能であること
- agent-taskstate 連携を強くしすぎず、疎結合に保つこと

### 13.3 保守性

- `memx-core` 本体の汎用責務を崩さないこと
- resolver 機能は分離可能な構成であること
- schema 変更の影響範囲が明確であること

### 13.4 コスト効率

- 毎回全文を返す設計にしないこと
- 解決結果は chunk 単位で返し、コンテキスト消費を抑えること
- 検索と参照解決は長文読み込みより安価であること

## 14. データ要件

### 14.1 documents

保持項目の例

- doc_id
- doc_type
- title
- source_path
- version
- updated_at
- summary
- tags_json
- feature_keys_json
- body

### 14.2 document_chunks

保持項目の例

- chunk_id
- doc_id
- heading_path
- ordinal
- body
- token_estimate
- importance

### 14.3 document_links

保持項目の例

- src_doc_id
- dst_doc_id
- link_type

link_type 例

- depends_on
- implements
- derived_from
- supersedes
- related_to

### 14.4 task 側追加情報

`agent-taskstate` に保持する候補

- read_receipt
- task_dependency
- stale_reason

## 15. インターフェース要件

### 15.1 docs ingest

文書を登録し、必要に応じて chunk 化する API を提供すること。

想定入力

- doc_type
- title
- source_path
- version
- tags
- feature_keys
- body
- chunking

### 15.2 docs resolve

feature または task を入力に、必要 doc を返す API を提供すること。

想定出力

- required
- recommended
- reason
- version
- top_chunks

### 15.3 chunks get

doc_id または query から chunk を返す API を提供すること。

### 15.4 reads ack

read receipt を記録する API を提供すること。

### 15.5 stale check

task_id を入力に stale 状態を判定する API を提供すること。

## 16. Skill 要件

本システムは、エージェントから使いやすい Skill 面を持つこと。

最低限必要な Skill は以下とする。

- `/resolve-docs`
- `/read-chunks`
- `/ack-docs`
- `/stale-check`
- `/resolve-contract`

※ ただし Skill は見た目の入口であり、本体は API 契約と状態保持層であること。

## 17. chunk 方針

chunk の目的は長文分割そのものではない。  
契約可能な最小参照単位を作ることである。

そのため、chunk 方針は以下とする。

- 基本は見出し単位
- 長すぎる節のみ再分割
- summary chunk を別途用意可能
- chunk ごとに version の母体である doc を必ず参照可能
- chunk 単独で意味が崩れにくいこと

## 18. 制約

- 本家 `memx-core` に即時混入しない
- `agent-taskstate` を task 正本として尊重する
- `tracker-bridge-materials` を外部連携境界として扱う
- 文書原本は当面 Markdown のまま維持する
- 全文編集可能な CMS 化は行わない

## 19. 受け入れ条件

以下を満たした時点で MVP 完了とする。

- `workflow-cookbook` の主要文書を doc として登録できる
- feature 名から required / recommended docs を返せる
- doc から必要 chunk を取得できる
- 読了した doc version と chunk_ids を task に紐づけて記録できる
- doc version 更新時に stale 判定できる

## 20. リスク

- resolver の責務が肥大化すると `memx-core` 本体を侵食する
- chunking の粒度設計を誤ると、逆に参照効率が悪化する
- taskstate 連携が曖昧だと、読了記録だけ存在して拘束力がない
- Birdseye 側を静的文書のまま扱うと stale 管理が弱い
- 本家リポジトリへ直接混ぜると OSS としての見通しが悪くなる

## 21. 実装方針

推奨する進め方は以下とする。

- 第1段階
  - branch または fork 上で PoC
  - documents / document_chunks を追加
  - resolve / get を先に作る

- 第2段階
  - read receipt を `agent-taskstate` と連携
  - stale 判定を実装

- 第3段階
  - `workflow-cookbook` を ingest する
  - Birdseye 系資料を紐づける
  - resolve 精度を改善する

## 22. 将来拡張

将来的には以下を検討可能とする。

- query 解決結果に基づく required section 提示
- forbidden pattern 自動検出
- definition of done の自動展開
- 外部 tracker / Birdseye view 同期
- 変更影響範囲の自動推定
- semantic diff に基づく stale 強度判定

## 23. 結論

作るべきものは、単なる Markdown 返却ツールではない。  
機能や task を入口として、読むべき資料と必要 chunk を解決し、その版を契約として残せる resolver surface である。

本要件は、その resolver surface を、既存の `memx-core`、`agent-taskstate`、`tracker-bridge-materials` を前提として最小構成で成立させるためのものである。