---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - CLI

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

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

- Requirement ID: `REQ-GC-001`

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
  4. 観測ノートは `chronicle.db` に `notes` として Insert。
  5. `lineage` に `relation='observed'` を記録。

- Phase 2: Reflector
  1. `chronicle` 側で、同一テーマ（タグ／トピック）に属する観測ノート群を抽出。
  2. `memopedia` に既存ページがあれば：
     - 既存本文 + 観測ノート群をコンテキストとして、"統合された最新版ページ" を生成（Update）。
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