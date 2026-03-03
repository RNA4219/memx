---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
priority: high
---

# memx 要件定義（v1.3）

## 0. 目的とスコープ

### 目的

- 個人用・ローカル運用の「知識＋記憶」管理システムを提供する。
- LLM／エージェントが利用しやすい形で、
  - 短期メモ（ログ）
  - 時系列の出来事（Chronicle）
  - 抽象化された知識ベース（Memopedia）
  - 歴史的ログ（Archive）
  を分離・格納し、検索・蒸留・昇格できるようにする。
- CLI + API + SQLite ベースで、言語や実行環境に強く依存しない「知識 OS の下層」を構成する。

### 非ゴール

- Web UI の提供。
- マルチユーザー運用やリモート公開を前提にしたサーバー運用（認証・権限・監査などを含む）。
- 「常駐が必須」なプロセス設計（API サーバーは任意で起動できればよい）。
- 完全自動運転のエージェントフレームワーク。あくまで「記憶と知識の基盤」。

## 0-1. Release Scope Matrix

### CLI

| MUST (v1) | SHOULD (v1.x) | FUTURE (v1.1+) |
| --- | --- | --- |
| `mem in short`, `mem out search`, `mem out show` | `mem gc short`（`mem.features.gc_short=true` 時のみ有効な実験機能） | `mem out recall`, `mem working`, `mem tag`, `mem meta`, `mem lineage`, `mem distill`, `mem out context` |

### API

| MUST (v1) | SHOULD (v1.x) | FUTURE (v1.1+) |
| --- | --- | --- |
| `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}` | `POST /v1/gc:run`（`mem.features.gc_short=true` 時のみ有効な実験機能） | Recall/Working/Tag/Meta/Lineage 系 API |

### 受け入れ条件

| 区分 | 条件 |
| --- | --- |
| v1必須 | 入出力互換（CLI→API の入出力マッピングが保持されること） |
| v1必須 | エラーコード（入力不備: 400系 / 内部障害: 500系 を返すこと） |
| v1必須 | 最小性能目標（`ingest`/`search`/`show` がローカル単体で実用応答時間を維持すること） |

---

## 1. 全体アーキテクチャ

### 1-0. レイヤリング（v1.3 変更点）

**目的**：
人間向けの CLI と、ツール/AI 向けの API を分離し、
UI/呼び出し手段が変わっても中核ロジックを固定できるようにする。

基本フロー：

```
[Human] CLI  →  [Tool/AI] API  →  Service(Usecase)  →  DB/LLM/Gatekeeper
```

- CLI は「入力の整形」と「表示」だけを担当し、DB へ直接アクセスしない。
- API は JSON の安定 I/F（HTTP または in-proc）として提供する。
- Service はビジネスロジックの唯一の入口。
- DB/LLM/Gatekeeper はインフラ層。

### 1-1. ストア構成（DBファイル）

物理的に 4 つの SQLite DB を用意する：

- `short.db`
  - 短期メモ／ログの一次格納場所。
  - すべてのノートはまずここに入る。
- `chronicle.db`
  - 日記・旅程・プロジェクト進捗など、「時間軸で意味を持つログ」。
- `memopedia.db`
  - 用語定義・設計・方針など、「時間軸から独立した知識ベース」。
- `archive.db`
  - GC によって退避された、古い or 低優先度のノート。
  - 通常検索からは外すが、バックトラックのために保持する。

### 1-2. ストレージ層

- すべての DB に共通して、以下を持つ想定：
  - `notes` … ノート本体
  - `tags` / `note_tags` … タグとその対応
  - `note_embeddings` … ベクター検索用埋め込み
  - `notes_fts` … FTS5 を使った全文検索（※ archive は省略可能）
- `short.db` のみ、追加で：
  - `short_meta` … GC トリガ用メタ情報
  - `lineage` … 「どのストアのどのノートが、どこに蒸留／昇格されたか」の系譜情報

### 1-2-1. `migrate_other.go` TODO 対応方針（先行固定）

- `go/db/migrate_other.go` の `migrateChronicle` / `migrateMemopedia` / `migrateArchive` は、`memx_spec_v3/schema/*.sql` の DDL ファイル適用に統一する。
- 適用順序は `notes` → `notes_fts`（採用時のみ）→ `tags`/`note_tags` → `note_embeddings`（採用時のみ）→ `PRAGMA user_version`。
- 実装では `schemaName` を SQL ファイル選択に使用し、`Exec` の多段実行失敗時は即時にエラーを返す（部分適用はトランザクションでロールバック）。
- `PRAGMA user_version` は「初期作成時に 1」「破壊的/非互換 DDL のみ +1」を共通ルールとし、データ移行を伴わないインデックス追加等は据え置く。

### 1-2-2. short.sql との差分方針（store 別）

| Store | notes | tags | note_tags | note_embeddings | notes_fts | `working_scope` | `is_pinned` | 初期 `user_version` |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| short | 採用 | 採用 | 採用 | 採用 | 採用 | 未採用 | 未採用 | 1 |
| chronicle | 採用 | 採用 | 採用 | 採用 | 採用 | v1 で必須列として採用 | v1 で任意列として採用（`0/1`） | 1 |
| memopedia | 採用 | 採用 | 採用 | 採用 | 採用 | v1 で必須列として採用 | v1 で任意列として採用（`0/1`） | 1 |
| archive | 採用 | 採用 | 採用 | 不採用 | 不採用 | 未採用 | 未採用 | 1 |

- 差分 DDL は `schema/chronicle.sql` / `schema/memopedia.sql` / `schema/archive.sql` に明記する。
- `working_scope` / `is_pinned` の導入対象は chronicle・memopedia のみとし、archive には導入しない。

### 1-3. 検索層

- FTS5 によるキーワード検索（content table モード + DELETE/INSERT トリガ）。
- `note_embeddings` によるベクター検索（Semantic Recall）。
- 全ストアを横断して検索できる `mem out recall` を提供。

#### ベクター検索の実装方針

- v1 では、ノート数が数千程度までは Go 側での愚直な cosine 類似度計算を許容する。
- 将来的に以下のような SQLite 拡張による高速化を想定し、API からは隠蔽する：
  - `sqlite-vec` / `sqlite-vss` などのベクター検索拡張
- 要件レベルでは「EmbeddingClient + note_embeddings テーブルを前提としたベクター検索インターフェース」を固定し、内部実装は後から差し替え可能とする。

### 1-4. 評価・ポリシー層（Judgement & Policy）

- 各ノートは評価軸を持つ：
  - `relevance`（関連度）
  - `quality`（質）
  - `novelty`（新規性）
  - `importance_static`（静的な重要度）
- 評価は、ルール＋軽量 LLM（MiniLLM）で行う。
- `memory_policy.yaml` で閾値や禁止パターン、decay ポリシーを管理する。

### 1-5. セーフティ層（Gatekeeper）

- Gatekeeper 1B モデル＋ルール群を用意し、
  - `memory_store`（保存前）
  - `memory_output`（出力前）
  のタイミングでチェックできるフックを持つ。
- 判断：`allow` / `deny` / `needs_human`。
- Go 側では `go/db/gatekeeper.go` のインターフェースで表現し、`Conn` 経由で利用する。
- `needs_human` は v1.3 では **保留ではなく deny 相当** として扱い、保存/出力を中断する（後で人手承認フローを導入するまでは fail-closed）。
- API 返却は `go/api/errors.go` の現行マッピング方針に合わせ、`service.ErrInvalidArgument` / `service.ErrNotFound` 以外は `INTERNAL` にフォールバックする。将来 `GATEKEEP_DENY` 等を返す場合も、まず service 層で sentinel error を定義し `mapError` へ明示追加する。

---

## 2. データモデル（short.db）

### 2-1. Note（短期ノート）

`schema/short.sql` に定義される `notes` テーブル：

- 主キー
  - `id: TEXT` … UUID 等
- 本文まわり
  - `title: TEXT NOT NULL` … 一行要約
  - `summary: TEXT NOT NULL DEFAULT ''` … 3〜5 行の要約（LLM失敗時や --no-llm 時は空文字でよい）
  - `body: TEXT NOT NULL` … 元テキスト or 要約本文
- 日付・アクセス
  - `created_at: TEXT NOT NULL` … ISO8601
  - `updated_at: TEXT NOT NULL`
  - `last_accessed_at: TEXT NOT NULL`
  - `access_count: INTEGER NOT NULL DEFAULT 0`
- ソース情報
  - `source_type: TEXT NOT NULL` … `'web' | 'file' | 'chat' | 'agent' | 'manual'` 等
  - `origin: TEXT NOT NULL DEFAULT ''` … URL / ファイルパス / エージェント名など。manual などで不在のときは空文字。
  - `source_trust: TEXT NOT NULL` … `'trusted' | 'user_input' | 'untrusted'`
  - `sensitivity: TEXT NOT NULL` … `'public' | 'internal' | 'secret'`
- 評価スコア
  - `relevance: REAL`（0〜1, null 可）
  - `quality: REAL`
  - `novelty: REAL`
  - `importance_static: REAL`
- ルーティング
  - `route_override: TEXT`
    - `NULL` or `'chronicle' | 'memopedia' | 'both' | 'archive_only'`

### 2-2. タグ

- `tags` / `note_tags` で多対多を管理。
- `tags.route` で昇格先のデフォルトを決める：
  - `'chronicle' | 'memopedia' | 'both' | 'short_only'`
- `tags.parent_id` によりタグ統合・エイリアス管理を可能にする。
- 将来、タグ蒸留時に類似タグの canonicalization を行うことを想定。

### 2-3. 埋め込み

- `note_embeddings`：
  - `note_id: TEXT`
  - `dim: INTEGER`
  - `vector: BLOB`（float32 配列）

### 2-4. メタ情報（short_meta）

- `short_meta`：
  - `key: TEXT PRIMARY KEY`
  - `value: TEXT`

用途：

- `note_count` … ノート数の概算（ヒューリスティック）
- `token_sum` … 全ノートのトークン数概算
- `last_gc_at` … 最後に GC を回した時間

運用方針：

- INSERT/DELETE 毎の厳密なトリガ管理は行わず、「近似的なヒント」として扱う。
- `mem in` / `mem gc` などで opportunistic に更新しつつ、GC 実行直前には `SELECT COUNT(*)` 等で正確値を再計算して閾値判定に使う。
- これにより、途中失敗による meta のズレを許容しつつ、安全側（GC をサボらない方向）に寄せる。

### 2-5. 系譜（lineage）

- `lineage`：
  - `id: INTEGER PRIMARY KEY AUTOINCREMENT`
  - `src_store: TEXT` … `'short' / 'chronicle' / 'memopedia'`
  - `src_note_id: TEXT`
  - `dest_store: TEXT` … `'chronicle' / 'memopedia' / 'archive'`
  - `dest_note_id: TEXT`
  - `relation: TEXT` … `'distilled_to' | 'merged_into' | 'observed' | 'reflected' | 'archived_from'` 等
  - `created_at: TEXT`

将来：

- `mem lineage show <note-id>` のような CLI を追加し、あるノートの出自・派生元を可視化することを想定。

### 2-6. スキーマバージョン管理

- `short.sql` の末尾で `PRAGMA user_version = 1;` を設定する。
- 将来のスキーマ変更は、`user_version` をインクリメントし、`migrateShort` 内で `PRAGMA user_version` の値に応じて ALTER を実行する方針とする。
- `CREATE TABLE IF NOT EXISTS` 方式と組み合わせることで、初回作成／既存DBの変更を両立させる。

### 2-7. Security & Retention Requirements

本節を `sensitivity`・保持・削除の**正本要件**とする。運用手順は [付録 G: Security & Privacy](../../docs/addenda/G_Security_Privacy.md) と [memx最小版: セキュリティ運用補足](../../docs/security/minimal_operations.md) を参照し、重複定義は行わない。

#### 2-7-1. `sensitivity` 別の保存可否・マスキング・保持期間

| sensitivity | 永続保存可否 | マスキング要件 | 保持期間要件 |
| --- | --- | --- | --- |
| `public` | `short/chronicle/memopedia/archive` へ保存可 | 原則マスク不要。ただし認証情報パターン検知時は `***REDACTED***` へ置換して保存/出力 | `short` 退避後 `archive` で最大 365 日保持し、期限到達後は物理削除対象 |
| `internal` | ローカル環境の `short/chronicle/memopedia/archive` のみ保存可（外部共有禁止） | CLI/API 出力および運用ログで識別子を部分マスク（例: `a***@example.com`） | `short` 退避後 `archive` で最大 180 日保持し、期限到達後は物理削除対象 |
| `secret` | 永続保存禁止（`memory_store` で `deny`） | 保存前・出力前とも全量マスク（`***REDACTED***`）し、原文は監査ログにも残さない | 保持対象外。検知時は処理を中断し、イベントメタデータのみ 90 日保持 |

- `needs_human` は v1.3 では `deny` 同等で扱う（fail-closed）。
- 保持期間の起算日は `archive` へ退避した `created_at` とし、期限判定は日次バッチまたは `mem gc --purge` 実行時に行う。

#### 2-7-2. `archive` 退避条件（Short → Archive）

`phase3` の退避実行は、次の条件を**すべて満たす**場合のみ許可する。

1. GC 候補選定済みである（`memory_policy.yaml.gc.short` の閾値判定を通過）。
2. `route_override` やタグルーティングで `chronicle/memopedia/both` への強制昇格が指定されていない。
3. `sensitivity` が `public` または `internal` である（`secret` は対象外）。
4. `archive` 側 Insert と `lineage(relation='archived_from')` の記録に成功している。

実行主体（actor）:
- 自動実行: `mem gc` 定期ジョブ（Maintainer 管理下の実行ユーザー）。
- 手動実行: Repo Maintainer（または委任された運用担当者）。

監査ログ要件（必須）:
- `op=archive_move`, `src_note_id`, `dest_note_id`, `sensitivity`, `executor`, `executed_at`, `result` を記録する。
- 失敗時は `error_code` と再試行可否（`retryable/non_retryable`）を追加する。

#### 2-7-3. 物理削除条件（Archive Purge）

`archive` からの物理削除は、次の条件を**すべて満たす**場合のみ許可する。

1. 2-7-1 の保持期間を超過している。
2. 対象ノートにインシデント調査中フラグ（legal/investigation hold）が無い。
3. 削除直前に `lineage` と `archive.notes` の対応整合を確認済みである。

実行主体（actor）:
- 自動実行: Maintainer が設定した定期ジョブのみ。
- 手動実行: Repo Maintainer が `mem gc --purge` を実行し、必要時は Security Champion へ事後報告。

監査ログ要件（必須）:
- `op=archive_purge`, `note_id`, `sensitivity`, `retention_expired_at`, `executor`, `executed_at`, `result` を記録する。
- 削除不能時は `reason`, `next_attempt_at`, `owner` を残し、次回削除計画を明示する。


---

## 3. CLI 要件

- Requirement ID: `REQ-CLI-001`

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`

### 3-1. 全体

- コマンド名：`mem`
- v1.3 以降、CLI は **API の薄いラッパ** として実装する。
  - 例：`mem in short ...` は `POST /v1/notes:ingest` に対応
  - 例：`mem out search ...` は `POST /v1/notes:search` に対応
  - CLI のオプションは、原則 API の request フィールドへ 1:1 でマップ

- サブコマンド構成（v1.3 時点）：
  - `mem in` … ノートの投入
    - `mem in short` … 短期ストアへの投入
  - `mem out` … ノートの取得
    - `mem out search` … FTS ベース検索
    - `mem out recall` … Semantic Recall
    - `mem out show` … 単一ノート表示
    - `mem out context` … LLMコンテキスト向け出力（将来）
  - `mem api` … API 操作
    - `mem api serve` … ローカル API サーバーを起動（任意）
  - `mem gc` … GC／蒸留
    - `mem gc short`
  - `mem distill` … 手動蒸留（将来）
  - `mem working` … Working Memory 操作（将来、memopedia に対して）
  - `mem tag` … タグ操作（将来）
  - `mem meta` … メタ情報表示（将来）
  - `mem lineage` … 系譜の可視化（将来）

### 3-2. `mem in short`

役割：生テキストから short ノートを作成し、API 経由で `short.db` に保存する。

例：

```bash
mem in short   --title "Qwen3.5-27B ローカルメモ"   --file ./note.txt   --source-type web   --origin "https://example.com/article"
```

オプション案：

- `--no-llm` … MiniLLM/Embedding を使わず、生テキスト＋最低限のメタだけ保存する。`summary` は空文字、タグは空のまま。
- `--tags` … 手動タグを付与する（カンマ区切り）。

処理フロー（v1.3）：

1. CLI は `file` / stdin から本文を読み込み、request を組み立てる。
2. CLI は API（in-proc もしくは HTTP）へ request を送る。
3. API/Service 側で以下を実行：
   - Gatekeeper（kind=`memory_store`）で保存可否を確認（必要なら）
   - `--no-llm` 相当の分岐（v1.3 ではフックのみ）
   - `tags` / `note_tags` / `notes` / `note_embeddings` / `notes_fts` を更新
   - `short_meta` を近似的に更新
4. CLI はレスポンス（note id など）を人間向けに整形して表示する。

### 3-3. `mem out search`（FTS）

役割：キーワード検索（FTS5）。

例：

```bash
mem out search "Qwen3.5 ベンチ"   --store short   --limit 10
```

- `notes_fts` は content table モードとし、UPDATE 時は `DELETE → INSERT` で同期する（`schema/short.sql` 参照）。

### 3-4. `mem out recall`（Semantic Recall）

Mastra の Semantic Recall 相当。

例：

```bash
mem out recall "Qwen3.5-27B ベンチマーク結果"   --scope self   --stores short,chronicle,memopedia   --top-k 8   --range 3
```

パラメータ：

- `--scope`：
  - `self`（デフォルト）
  - `session`（将来拡張）
  - `project:<name>`（将来拡張）
- `--stores`：検索対象ストア（カンマ区切り）
- `--top-k`：ベクター検索の anchor 数
- `--range`：anchor 前後何件を同ストアから連結するか

内部仕様（疑似仕様 / 実装可能レベル）：

1. クエリを EmbeddingClient で embed → ベクター取得。
2. 指定ストアの `note_embeddings` を横断し cosine 類似度を計算（将来 sqlite-vec 等に差し替え）。
   - 類似度式：`score = dot(q, v) / (||q|| * ||v||)`（q: クエリ埋め込み, v: ノート埋め込み）
   - `||q|| == 0` または `||v|| == 0` の場合は `score = 0` とみなす。
   - 取得対象は `score >= 0.20` のノートのみ（閾値）。
3. スコア上位 `top-k` 件を anchor とする。
   - `top-k` の有効範囲は `1..50`。未指定時は `8`。
   - `top-k > 50` 指定時は `50` に丸める。
   - 同点時タイブレークは `created_at DESC` → `id ASC`。
4. anchor ごとに `created_at` ベースで `range` 件の Before/After ノートを取得。
   - `range` は `0..10` の整数（未指定時 `3`）。
   - 先頭ノートでは `Before` は存在分のみ（0 件を許容）。
   - 末尾ノートでは `After` は存在分のみ（0 件を許容）。
5. `--stores` は以下で正規化する。
   - 入力文字列を `,` で分割し、前後空白を trim、空要素を除去。
   - 小文字化して `short|chronicle|memopedia|archive` に解決する。
   - 重複は先に出現した順で一意化する。
   - 未指定時は `short,chronicle,memopedia`。
   - 不正値を含む場合は 400 系入力エラーとして失敗させる。
6. `Conn.Embed == nil` の場合は実行モードで分岐する。
   - デフォルトはエラー（`semantic recall requires embedding client`）。
   - 明示フラグ（例: `--allow-fts-fallback`）指定時のみ FTS 限定検索へフォールバックする。
7. Working Memory（memopedia の pinned ノート）がある場合は、結果の先頭にマージする。

### 3-5. `mem gc short`（Observer / Reflector）

スコープ区分：**SHOULD (v1.x)**。v1 では `mem.features.gc_short=true` を明示した場合のみ有効化する実験機能とし、デフォルトでは無効。

Mastra の Observational Memory を参考にした GC。

例：

```bash
mem gc short          # 通常実行
mem gc short --dry-run
```

オプション：

- `--dry-run` … 実際には DB を変更せず、予定されている操作だけ表示。

フロー：

- Phase 0: トリガ判定
  - `short_meta` から note_count / token_sum / last_gc_at を参照し、
    `memory_policy.yaml.gc.short` の閾値を使って判定する。
  - 判定に使うキー（`short_meta` 由来）は以下とする：
    - `soft_limit_notes`: `1200`
      - `note_count >= 1200` かつ `last_gc_at` から `min_interval_minutes` 以上経過で GC 実行対象。
    - `hard_limit_notes`: `2000`
      - `note_count >= 2000` なら `min_interval_minutes` を無視して強制実行。
    - `min_interval_minutes`: `180`
      - 直近 GC 実行から 180 分未満の場合、soft limit 到達のみではスキップ。
  - 実際に GC を行う場合は、`SELECT COUNT(*)` 等で正確値を取得してから閾値を確認。
  - 設定参照元は `memory_policy.yaml.gc.short` のみとし、`go/db/gc.go` 実装時に定数の重複定義を禁止する。

- `--dry-run` の予定操作フォーマット（JSON）

```json
{
  "target": "short",
  "phase": "phase0|phase1|phase2|phase3",
  "decision": {
    "should_run": true,
    "reason": "soft_limit_reached|hard_limit_reached|interval_not_elapsed",
    "metrics": {
      "note_count": 1324,
      "soft_limit_notes": 1200,
      "hard_limit_notes": 2000,
      "minutes_since_last_gc": 241,
      "min_interval_minutes": 180
    }
  },
  "planned_ops": [
    {
      "op": "observe_cluster",
      "src_note_ids": ["n1", "n2"],
      "dest_store": "chronicle"
    },
    {
      "op": "archive_move",
      "src_note_id": "n3",
      "dest_store": "archive",
      "lineage_relation": "archived_from"
    }
  ]
}
```

- Phase 1: Observer
  1. 古い／アクセスが少ない short ノートを対象集合として抽出。
  2. タグ＋embedding 類似度でクラスタリング。
  3. 各クラスタを MiniLLM/ReflectLLM に渡し、「観測ノート（Observation）」を生成。
  4. 관찰 노트는 `chronicle.db` 에 `notes` 로 Insert。（※実装時に日本語へ）
  5. `lineage` に `relation='observed'` を記録。

- Phase 2: Reflector
  1. `chronicle` 側で、同一テーマ（タグ／トピック）に属する観測ノート群を抽出。
  2. `memopedia` に既存ページがあれば：
     - 既存本文 + 관찰 노트群을 컨텍스트로, "統合された最新版ページ" を生成（Update）。
  3. なければ：新規ページとして Insert。
  4. `lineage` に `relation='reflected'` を記録。

- Phase 3: Short → Archive（補償設計）
  - short → archive の退避は、SQLite の ATTACH の制約により「完全な原子的操作」にはできない。
  - したがって、次のポリシーを取る：
    - 先に `archive` 側へ Insert → `lineage` に `archived_from` を記録。
    - 最後に `short` 側から Delete。
    - 途中で失敗した場合でも、short に元データが残る or archive に複製が残る形を優先し、「データ喪失より重複を許容」する。
  - `mem gc` 再実行時に、lineage と実データを突き合わせて「重複を整理する」処理を追加可能とする。
  - 再実行時の整合ルール（重複許容後の収束条件）：
    1. `lineage(src_store='short', src_note_id, dest_store='archive', relation='archived_from')` が存在し、かつ `archive.notes.id=dest_note_id` が存在する場合、同一 `src_note_id` の archive 追加 Insert は行わない。
    2. `lineage` があるのに `archive.notes.id=dest_note_id` が欠損している場合、同一 `src_note_id` の再退避を 1 回だけ許可し、新しい `dest_note_id` で lineage を追記する（過去 lineage は監査用に保持）。
    3. archive への複製が存在し lineage が欠損している場合、`src_note_id + dest_note_id + relation='archived_from'` で lineage を補完してから short 側 Delete 判定を行う。
    4. short 側 Delete は「archive 実在 + 対応 lineage 実在」を満たす場合のみ実行する。

### 3-6. `mem working`（Working Memory）※将来

- `memopedia.db` の `notes` に以下の列を追加予定：
  - `working_scope: TEXT` … `NULL` or `'global'` or `'session:<id>'` or `'project:<name>'`
  - `is_pinned: INTEGER` … 1 なら Working Memory として常時読み出し

CLI 想定：

```bash
mem working pin <note-id> --scope global
mem working list --scope global
mem working unpin <note-id>
```

検索系（`mem out recall` / `mem out context`）は、該当する `working_scope` の `is_pinned=1` ノートを必ず先頭に含める。

---

## 4. LLM 戦略

### 4-1. 役割分離

LLM は少なくとも 3 役割に分離する：

1. EmbeddingClient
   - テキスト → ベクター 変換。
   - Semantic Recall で使用。
2. MiniLLMClient
   - タグ生成・スコアリング（`relevance / quality / novelty / importance`）・`sensitivity` 推定。
   - 軽量モデル（1B〜3B）を想定。
3. ReflectLLMClient
   - クラスタ要約（Observer）・Memopedia ページ更新（Reflector）。
   - 7B〜27B クラスのモデルを想定。

Go 側では `go/db/llm_client.go` に interface を定義し、`db.Conn` にこれらをフィールドとして注入して使う。

### 4-2. 同期／非同期の扱い

- `mem in` 実行時に全ての LLM を同期で呼ぶとレイテンシが伸びる。
- v1 では：
  - `mem in` では最低限のフィールド（title/body/source_*）だけ即時保存し、
  - タグ付け・スコアリング・埋め込み生成はオプションで非同期キューに積む実装も許容範囲とする。
- CLI としては：
  - `--no-llm` で完全に LLM を使わない形
  - デフォルトでは同期処理（ただし後でオプションで非同期化も検討）
- 要件レベルでは、「LLM を使うか／どのタイミングで使うか」を `mem in` のフラグと設定ファイルで切り替え可能にする。

### 4-3. 設定例

`config.yaml` のイメージ：

```yaml
llm:
  embed:
    provider: local
    endpoint: "http://localhost:8000/embed"
  mini:
    provider: local
    endpoint: "http://localhost:8001/generate"
  reflect:
    provider: local
    endpoint: "http://localhost:8002/generate"
```

`memory_policy.yaml`（GC 関連キー雛形）：

```yaml
version: 1

gc:
  short:
    soft_limit_notes: 1200
    hard_limit_notes: 2000
    min_interval_minutes: 180
    target_delete_batch_size: 200
    max_archive_retries: 1
```
### 4-4. 各クライアント共通の呼び出し契約

対象：`EmbeddingClient` / `MiniLLMClient` / `ReflectLLMClient` / Gatekeeper 呼び出し。

- タイムアウト：**1 リクエスト 15 秒**（`context.WithTimeout` 等で必須化）。
- 最大リトライ回数：**2 回**（初回 + リトライ 2 = 最大 3 試行、指数バックオフ）。
- 再試行可エラー：
  - ネットワーク断、接続リセット、タイムアウト
  - HTTP 429 / 502 / 503 / 504
- 再試行不可エラー：
  - 入力不正（HTTP 400 相当）
  - 認証/認可失敗（HTTP 401/403 相当）
  - モデル仕様不一致（JSON スキーマ不整合、必須フィールド欠落）
- ingest 時の部分失敗ポリシー：
  - `notes` 保存成功をコミット境界の最小単位とし、ノート本体保存は継続。
  - タグ生成・`note_tags`・埋め込み生成の失敗は後追い再実行対象として記録し、ingest 全体は成功扱いにできる。
  - ただし Gatekeeper が `deny` / `needs_human` を返した場合は fail-closed で ingest 全体を失敗にする。
- エラーコードのマッピング方針（`go/api/errors.go` 整合）：
  - API 返却コードは `INVALID_ARGUMENT` / `NOT_FOUND` / `INTERNAL` を最小集合として保証する。
  - クライアント個別エラーは service 層で sentinel error に正規化し、`go/api/errors.go` の `mapError` に 1 箇所で集約マップする。
  - 未分類エラーは互換性維持のため常に `INTERNAL` にフォールバックする。

---

## 5. 非機能要件

### 5-1. 性能目標（v1必須3エンドポイント）

計測条件は以下で固定する。

- データセット条件: short ストア 10,000 件、本文 1 ノートあたり約 500 文字（UTF-8 プレーンテキスト）
- 実行環境: ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）
- ウォームアップ: 各エンドポイント 20 回
- 本計測: 各エンドポイント 200 回（ウォームアップ除外）

| 操作 | API | p50 目標 | p95 目標 |
| --- | --- | --- | --- |
| ingest | `POST /v1/notes:ingest` | `<= 120ms` | `<= 250ms` |
| search | `POST /v1/notes:search` | `<= 80ms` | `<= 180ms` |
| show | `GET /v1/notes/{id}` | `<= 40ms` | `<= 90ms` |
- Requirement ID: `REQ-NFR-001`

#### 計測プロトコル

- データ生成条件（固定）:
  - 対象ストアは `short` 固定。
  - 事前投入データは 10,000 件固定。
  - 各ノート本文は約 500 文字（UTF-8 プレーンテキスト）固定。
- 計測回数（固定）:
  - 各エンドポイントごとにウォームアップ 20 回を先に実行。
  - 本計測は各エンドポイント 200 回で固定し、ウォームアップを集計から除外する。
- 除外条件（固定）:
  - ウォームアップ 20 回は母集団に含めない。
  - `REQ-NFR-001` 判定では `ingest` / `search` / `show` の 3 エンドポイント以外を除外する。
  - ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）以外の環境で得た計測値は正式判定から除外する。
- 集計方法（固定）:
  - 本計測 200 回のレイテンシ分布から `p50` と `p95` を算出する。
  - 合否判定値は `p50_ms` / `p95_ms` のみとし、平均値や最大値は参考情報として扱う。
  - 出力フォーマットは `RUNBOOK.md` の `artifacts/perf/perf-result.json` に定義されたキーを正本とする。

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`
- `RUNBOOK.md`

- OS：ローカル（Linux / macOS / Windows）で動作する CLI を想定。
- DB：SQLite3（WAL モード、foreign_keys ON）。
- 言語：Go（単一バイナリビルドを前提）。
- 依存：
  - SQLite ドライバ（例：`modernc.org/sqlite` / `github.com/mattn/go-sqlite3`）。
  - CLI は標準 `flag` でもよい（薄いラッパが前提）。
  - HTTP API は標準 `net/http` でよい。
- セキュリティ：
  - APIキーや秘密情報は `memory_policy.yaml` + Gatekeeper により保存前にブロック。
- 拡張性：
  - `chronicle.db` / `memopedia.db` / `archive.db` のスキーマは、`short.db` と同様の `notes` / `tags` / `note_tags` / `note_embeddings` / `notes_fts` 構造をベースとする。
  - 将来、Working Memory／プロジェクトスコープ／セッションスコープを追加しても、既存の CLI と DB 構造を壊さない。
- トランザクションと一貫性：
  - ATTACH を跨いだ完全な原子的トランザクションは SQLite の仕様上保証できない。
  - 設計上「データ喪失より重複を許容する」ポリシーとし、lineage による追跡・再蒸留で整合性を取り戻せるようにする。

---

## 6. API 要件（v1.3 追加）

- Requirement ID: `REQ-API-001`

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`
- `memx_spec_v3/docs/quickstart.md`

### 6-1. 目的

- ツール/AI から呼びやすい **安定 JSON I/F** を提供する。
- CLI は API の薄いラッパとして実装し、ビジネスロジックを持たない。

### 6-2. 提供形態

- **HTTP（ローカル）**：`mem api serve` で起動。
  - 例：`http://127.0.0.1:7766`
  - 将来的に unix socket 対応も想定。
- **in-proc**：CLI や別バイナリが同一プロセスで呼ぶ。

### 6-3. エンドポイント（v1）

- `GET /healthz` → `ok`
- `POST /v1/notes:ingest`
  - request: `{title, body, summary?, source_type?, origin?, source_trust?, sensitivity?, tags?}`
  - response: `{note: Note}`
- `POST /v1/notes:search`
  - request: `{query, top_k?}`
  - response: `{notes: Note[]}`
- `GET /v1/notes/{id}` → `Note`
- `POST /v1/gc:run`（SHOULD (v1.x): `mem.features.gc_short=true` 時のみ有効な実験機能）
  - request: `{target, options?}`
  - response: `{status}`

`POST /v1/gc:run` は v1 MUST には含めない。`mem.features.gc_short=false` の場合、全環境で **HTTP 409 + `{ "code": "FEATURE_DISABLED", "message": "gc_short feature is disabled" }`** を固定で返す（404 へのフォールバック禁止）。

この無効時挙動は「デプロイ単位で選択」ではなく **全環境共通固定** とし、クライアントは `FEATURE_DISABLED` を恒久エラーとして扱う（同一条件での自動リトライ不可、運用者による feature flag 変更後のみ再試行可）。

### 6-3-1. v1必須3エンドポイント契約（`requirements.md` × `go/api/types.go` 照合）

#### POST `/v1/notes:ingest`

**request (`NotesIngestRequest`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `title` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。v1 では削除・意味変更禁止。 |
| `body` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。v1 では削除・意味変更禁止。 |
| `summary` | `string` | 任意 | `""` | 追加バリデーションなし | 任意フィールドのため、未指定互換を維持。 |
| `source_type` | `string` | 任意 | `"manual"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定（クライアント未送信を維持）。 |
| `origin` | `string` | 任意 | `""` | trim のみ | 任意フィールドとして後方互換追加可。 |
| `source_trust` | `string` | 任意 | `"user_input"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定。 |
| `sensitivity` | `string` | 任意 | `"internal"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定。 |
| `tags` | `[]string` | 任意 | `[]` 扱い（未指定時は処理なし） | 各要素 trim、空要素は無視 | 任意配列として未指定・空配列を同等扱い。 |

**response (`NotesIngestResponse`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `note` | `Note` | 必須 | なし | 保存成功時に返却 | v1 では wrapper 構造（`{note: ...}`）を維持。 |
| `note.id` | `string` | 必須 | なし | 32 hex 形式の ID を生成 | 型・キー名固定。 |
| `note.title` | `string` | 必須 | なし | request 値（trim 後） | 型固定。 |
| `note.summary` | `string` | 必須 | `""` 許容 | request 値 | 必須キーだが空文字を許容。 |
| `note.body` | `string` | 必須 | なし | request 値（trim 後） | 型固定。 |
| `note.created_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.updated_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.last_accessed_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.access_count` | `int64` | 必須 | `0` | ingest 直後は `0` | 数値型固定。 |
| `note.source_type` | `string` | 必須 | `"manual"` 補完あり | request/補完値 | 補完込みで常に返却。 |
| `note.origin` | `string` | 必須 | `""` | request 値 | 必須キーとして維持。 |
| `note.source_trust` | `string` | 必須 | `"user_input"` 補完あり | request/補完値 | 補完込みで常に返却。 |
| `note.sensitivity` | `string` | 必須 | `"internal"` 補完あり | request/補完値 | 補完込みで常に返却。 |

#### POST `/v1/notes:search`

**request (`NotesSearchRequest`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `query` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。 |
| `top_k` | `int` | 任意 | `20`（`<=0` 時） | `<=0` は service で `20` に補正 | 任意数値として未指定互換を維持。 |

**response (`NotesSearchResponse`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `notes` | `[]Note` | 必須 | `[]` | 一致結果を最大 `top_k` 件返却 | v1 では配列 wrapper 構造を維持。 |
| `notes[].*` | `Note` 各フィールド | 必須 | `Note` 契約準拠 | `Note` と同一 | `Note` フィールドの型・キー名を維持。 |

#### GET `/v1/notes/{id}`

**request（path parameter）**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `id` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | path パラメータ必須を維持。 |

**response (`Note`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `id` | `string` | 必須 | なし | ノート存在時に返却、非存在は `NOT_FOUND` | 型・キー名固定。 |
| `title` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `summary` | `string` | 必須 | `""` 許容 | 保存済み値 | 必須キーとして維持。 |
| `body` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `created_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `updated_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `last_accessed_at` | `string` | 必須 | なし | 取得時刻で更新 | 取得時更新の挙動を維持。 |
| `access_count` | `int64` | 必須 | なし | 取得ごとに +1 | カウンタ型・加算挙動を維持。 |
| `source_type` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `origin` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `source_trust` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `sensitivity` | `string` | 必須 | なし | 保存済み値 | 型固定。 |

### 6-4. エラーモデル

本節を ErrorCode 契約の正本とし、`error-contract.md` は本節の運用向け要約として同期する。

共通：

```json
{"code":"INVALID_ARGUMENT","message":"...","details":{}}
```

- `INVALID_ARGUMENT` → 400
- `NOT_FOUND` → 404
- `CONFLICT` → 409
- `GATEKEEP_DENY` → 403
- `INTERNAL` → 500

`go/service/errors` と `go/api/errors.go` の現行方針（`ErrInvalidArgument` / `ErrNotFound` を明示マッピングし、それ以外は `INTERNAL` へフォールバック）を前提に、ErrorCode を 2 段で再定義する。

#### ErrorCode 区分（v1）

| 区分 | ErrorCode | HTTP | 契約レベル | 条件 |
| --- | --- | --- | --- | --- |
| v1必須保証 | `INVALID_ARGUMENT` | 400 | MUST | 常時有効。 |
| v1必須保証 | `NOT_FOUND` | 404 | MUST | 常時有効。 |
| v1必須保証 | `INTERNAL` | 500 | MUST | 常時有効。未分類エラーのフォールバック先。 |
| v1.x拡張（feature/sentinel依存） | `CONFLICT` | 409 | SHOULD | service sentinel（例: `ErrConflict`）を実装し `mapError` へ明示追加した場合のみ返却。未実装時は `INTERNAL`。 |
| v1.x拡張（feature/sentinel依存） | `GATEKEEP_DENY` | 403 | SHOULD | gatekeeper deny 系 sentinel（例: `ErrGatekeepDeny`）を実装し `mapError` へ明示追加した場合のみ返却。未実装時は `INTERNAL`。 |

#### 現行実装との差分注記（`go/api/types.go` / `go/api/http_server.go`）

- `go/api/types.go` には `CONFLICT` / `GATEKEEP_DENY` 定数が定義済みだが、返却は service sentinel と `mapError` 実装に依存する。
- `go/api/http_server.go` は `writeErr` で `CONFLICT=409` / `GATEKEEP_DENY=403` を処理可能だが、上流が当該コードを返さない限り到達しない。
- sentinel 未実装時は `mapError` 方針に従って `INTERNAL`（500）へフォールバックする。

| 代表事象 | service 層の分類 | API `code` | HTTP | 再試行可否 | 備考 |
| --- | --- | --- | --- | --- | --- |
| 入力不備（必須欠落・形式不正・空文字） | `ErrInvalidArgument` | `INVALID_ARGUMENT` | 400 | 不可 | クライアント入力修正が必要。 |
| DB ロック（`database is locked` / busy timeout 超過） | 現行は汎用エラー（将来 `ErrConflict` 候補） | `INTERNAL`（将来 `CONFLICT` 検討可） | 500（将来 409 検討可） | 条件付き可 | 短時間で解消しうるため指数バックオフ再試行対象。長時間継続時は運用アラート。 |
| LLM タイムアウト（上流 timeout / 502/503/504） | 汎用エラー | `INTERNAL` | 500 | 条件付き可 | 一時障害として再試行対象。最大回数超過で失敗扱い。 |
| Gatekeeper deny / needs_human(v1.3 では deny 扱い) | 将来 sentinel 化（`ErrGatekeepDeny`） | `GATEKEEP_DENY`（未実装時は `INTERNAL`） | 403（未実装時は 500） | 不可 | ポリシー判定のため再試行では解消しない。入力/運用判断が必要。 |

#### エラーコード別 再試行可否表

| API `code` | HTTP | 再試行可否 | ルール |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | 400 | 不可 | 入力を修正して再実行。 |
| `NOT_FOUND` | 404 | 不可 | 対象 ID・クエリを見直す。 |
| `CONFLICT` | 409 | 条件付き可 | DB ロック等の一時競合のみ再試行。恒久競合は不可。 |
| `GATEKEEP_DENY` | 403 | 不可 | ポリシー deny は即時失敗。再試行禁止。 |
| `INTERNAL` | 500 | 条件付き可 | 原因が一時障害（DB lock / LLM timeout / upstream 502/503/504）の場合のみ指数バックオフ。恒久障害は不可。 |

### 6-5. 互換性ポリシー（v1）

- v1 系（`/v1/*`）では、既存クライアントを壊さない後方互換を維持する。

許容する変更（v1 内）：

- request/response への**任意フィールド追加**（既存必須フィールドは維持）。
- 新規エンドポイント追加（既存エンドポイントの契約は維持）。
- 既存フィールドのバリデーション強化のうち、既存の正常系入力を拒否しない変更。

禁止する変更（v1 内）：

- 既存の必須フィールド削除。
- 既存フィールドの意味変更（同じキー名で別意味にする変更）。
- 既存レスポンスの型変更（例：`string` → `object`）や、既存成功レスポンスの構造破壊。

破壊変更が必要な場合の移行手順：

1. まずは **新エンドポイント**（例：`/v1/notes:search2`）で並行提供し、既存挙動を維持する。
2. 並行提供で吸収できない場合は **`/v2` を新設**し、v1 を非推奨化する。
3. `CHANGES.md` に移行期限・差分・移行例を記載し、少なくとも 1 リリースは移行猶予を置く。

CLI 出力との互換責任範囲：

- 人間向けの通常表示（非 `--json`）は、可読性改善のための文言・並び変更を許容する。
- `--json` 出力は API の互換ポリシーと同等に扱い、v1 では破壊的変更を禁止する。
- 機械連携は API または CLI `--json` を正とし、互換責任はこの 2 系統に対して負う。

---

## 11. インシデント対応要件（運用）

- セキュリティ/品質インシデントは `docs/INCIDENT_TEMPLATE.md` に従って記録する。
- すべてのインシデント記録は `IN-YYYYMMDD-XXX` 形式のIDを持つ。
- 初動時点で「検知」「影響」「5 Whys」「再発防止」「タイムライン」を最低限記載する。
- サンプル記録: [`docs/IN-YYYYMMDD-001.md`](../../docs/IN-YYYYMMDD-001.md)
