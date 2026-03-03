````markdown
# memx

Local-first personal memory & knowledge store for LLM agents.

`memx` は、ローカルLLM／エージェントから使うことを前提にした  
**4層構成（short / chronicle / memopedia / archive）** のメモリ／知識ストアと CLI です。

ログをただ溜めるのではなく、

- 短期メモ（短期記憶）
- 時系列の出来事（エピソード記憶）
- 抽象化された知識（意味記憶）
- 退避された履歴（長期保管）

に分けて管理し、「蒸留（distill）」と GC（Observer / Reflector）を通じて  
ログを「参照しやすい知識」に育てていくことを目標としています。

> ※まだ **pre-alpha** 段階です。仕様が固まりつつあるところで、実装はこれから。

---

## Governance Docs

- [BLUEPRINT.md](./BLUEPRINT.md) - 方針（目的・設計原則・互換性方針）
- [RUNBOOK.md](./RUNBOOK.md) - 手順（主要運用フロー・実行手順）
- [GUARDRAILS.md](./GUARDRAILS.md) - ガードレール（安全性・互換性・移行制約）
- [EVALUATION.md](./EVALUATION.md) - 受け入れ基準（評価条件・判定基準）
- [CHECKLISTS.md](./CHECKLISTS.md) - チェックリスト（実装・運用確認項目）

---

## Goals

- ローカルで完結する「個人用・長期メモリ＆知識ストア」を提供する
- LLM／エージェントから使いやすい **シンプルな CLI と API** を用意する
- 「短期メモ → 観測ノート → 知識ページ」への**一方向のフロー**で、記憶の整理を自動化する
- Semantic Recall（意味ベース検索）＋タグ＋FTS で **検索性能を重視** する
- Web UI や常時稼働エージェントではなく、**バッチ／都度呼び出し前提の設計**にする

---

## Architecture Overview

### 4 Stores

物理的に 4つの SQLite DB を持ちます（ATTACH 前提）。

- `short.db`
  - すべてのノートが最初に入る「短期メモ」ストア
  - 断片的なログ・生のメモなど
- `chronicle.db`
  - 日記・旅程・プロジェクト進捗など、「時間軸で意味を持つログ」
- `memopedia.db`
  - 用語定義・設計・ポリシーなど、「時間軸から独立した知識ベース」
- `archive.db`
  - 古い／優先度の低いノートを退避するストア
  - 通常検索からは外すが、バックトラック用に保持

### Storage / Indexing

各ストアは基本的に同じ構造を持ちます（`short.db` は superset）:

- `notes` … ノート本体
- `tags` / `note_tags` … タグとノートの多対多
- `note_embeddings` … 意味検索用のベクター（埋め込み）
- `notes_fts` … FTS5 による全文検索インデックス（archive は optional）

`short.db` にのみ、追加で:

- `short_meta` … GC 用メタ情報（note_count / token_sum / last_gc_at など）
- `lineage` … 蒸留・昇格・退避の系譜（どのノートがどこへ統合されたか）

### LLM Roles

LLM は役割ごとに分けて扱います：

- **EmbeddingClient**
  - テキスト → ベクター（埋め込み）
  - Semantic Recall で使用
- **MiniLLMClient**
  - タグ生成
  - `relevance / quality / novelty / importance` の初期スコア推定
  - 機密度（`sensitivity`）推定
  - 軽量モデル（1B〜3B）想定
- **ReflectLLMClient**
  - 観測ノートクラスタの要約（Observer）
  - Memopedia ページの更新（Reflector）
  - 7B〜30B クラス想定

これらは Go の interface（`go/llm_client.go`）として定義され、  
`db.Conn` に注入して使う構成になっています。

---

## CLI Design (draft)

コマンド名は `mem` を想定しています。

### 基本コマンド

- `mem in short`
  - 短期ストア（short.db）へのノート投入
- `mem out search`
  - FTS5 によるキーワード検索
- `mem out recall`
  - ベクター検索＋前後文脈の Semantic Recall
- `mem gc short`
  - short.db の GC（Observer / Reflector を含む）

今後追加予定のもの：

- `mem distill`
  - 手動での蒸留／統合
- `mem working`
  - Memopedia の Working Memory（常時ピン留めノート）の操作
- `mem lineage`
  - あるノートがどのログから来たかを辿る

### 入力例（mem in short）

```bash
mem in short \
  --title "Qwen3.5-27B ローカルメモ" \
  --file ./note.txt \
  --source-type web \
  --origin "https://example.com/article"
````

オプション例:

* `--no-llm`

  * タグ付け・スコアリング・埋め込み生成を行わず、「生ノート」として保存
  * 後からバッチで LLM を流す運用も想定

### 検索例（mem out recall）

```bash
mem out recall "Qwen3.5-27B ベンチマーク" \
  --scope self \
  --stores short,chronicle,memopedia \
  --top-k 8 \
  --range 3
```

* クエリを embed → 類似ノート上位 `top-k` を anchor 取得
* 各 anchor の前後 `range` 件を同ストアから連結して返す
* 将来的には Working Memory（pinned ノート）を常に先頭に含める予定

---

## Status

* [x] 要件定義・アーキテクチャ設計（v1.2）
* [x] `short.db` 用のスキーマ (`schema/short.sql`)
* [x] Go 側の骨組み (`db.Conn`, `MustOpenAll`, migrate, LLM/Gatekeeper interface)
* [ ] CLI 実装（`mem in/out/gc`）
* [ ] Semantic Recall 実装（Recall API）
* [ ] GC（Observer / Reflector）の実装
* [ ] chronicle / memopedia / archive 向けスキーマ・実装

まだ「動くソフトウェア」というより、
**仕様とスケルトンが整った状態** です。

---

## Tech Stack

* Language: Go
* Storage: SQLite3

  * WAL モード／foreign_keys ON
  * FTS5（content table モード）
* Interface: CLI（単一バイナリ）
* 外部ベクターDBには依存せず、SQLite＋BLOB＋将来の拡張（`sqlite-vec` 等）で済ませる方針。

---

## Inspirations / Acknowledgements

このプロジェクトは、いくつかの既存アイデア・OSSに強く影響を受けています。
コードはゼロから書いていますが、コンセプト面で多くのヒントをもらっているので、ここで明示的に感謝します。

* SAIVerse

  * エージェント指向の設計や、長期メモリの扱い方に関する議論・実装から、「人間の記憶構造を意識したエージェント設計」という発想をかなり借りています。
* Mastra

  * Semantic Recall / Working Memory / Observational Memory といった概念と、その実装方針から強い影響を受けています。
    `memx` の設計は、これらのアイデアをローカル・CLI・Go/SQLite 文脈に落とし込んだ「別物」ですが、発想の出発点として大きく参考にしています。

また、

* 各種 LLM 長期記憶論文・ブログポスト
* 人間の記憶モデル（短期記憶／ワーキングメモリ／長期記憶など）に関する一般的な知見

からも、間接的にインスピレーションを受けています。

※もし将来的に、具体的な論文名・記事名・リポジトリ名を明示したくなったら、ここに追記していく予定です。

---

## License

Apache-2.0

---

## Roadmap (rough)

1. `mem in short --no-llm` / `mem out search` の実装
2. EmbeddingClient を繋いで `mem out recall` を実装
3. GC の Phase 0（トリガ判定）だけ実装
4. Observer → Reflector の順に GC を拡張
5. chronicle / memopedia / archive のスキーマ・マイグレーション実装
6. `mem lineage` / `mem working` など、日常利用で欲しくなったものから順次追加

---

## Contributing

現時点では個人用プロジェクトとしてスタートしていますが、
アイデア・Issue・PR などのフィードバックは歓迎です。(PR連打は困るので一度フォークしてからにしてください。)

* 「こういう検索フローが欲しい」
* 「このタグ設計だとこういうケースで詰まる」
* 「この論文・OSSのアイデアも組み込めるのでは？」

といった提案があれば、Issue で教えてもらえると助かります。
