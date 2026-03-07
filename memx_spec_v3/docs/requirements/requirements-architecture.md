---
owner: memx-core
status: active
last_reviewed_at: 2026-03-06
next_review_due: 2026-06-06
---

# memx 要求事項 - 全体アーキテクチャ

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

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
- ADR: [ADR-0001: 4DB分割と責務境界](../../docs/ADR/ADR-0001-4db-boundary.md)

物理的に 4 つの SQLite DB を用意する：

- `short.db` ✅
  - 短期メモ／ログの一次格納場所。
  - すべてのノートはまずここに入る。
- `chronicle.db` ✅（スキーマ完了、CRUD未実装）
  - 日記・旅程・プロジェクト進捗など、「時間軸で意味を持つログ」。
- `memopedia.db` ✅（スキーマ完了、CRUD未実装）
  - 用語定義・設計・方針など、「時間軸から独立した知識ベース」。
- `archive.db` ✅（スキーマ完了、CRUD未実装）
  - GC によって退避された、古い or 低優先度のノート。
  - 通常検索からは外すが、バックトラックのために保持する。

### 1-2. ストレージ層

- すべての DB に共通して、以下を持つ：
  - `notes` … ノート本体
  - `tags` / `note_tags` … タグとその対応
  - `*_meta` … GC トリガ用メタ情報（short_meta / chronicle_meta / memopedia_meta / archive_meta）
  - `lineage` … 「どのストアのどのノートが、どこに蒸留／昇格されたか」の系譜情報
- 検索対象ストア（short / chronicle / memopedia）は追加で：
  - `note_embeddings` … ベクター検索用埋め込み
  - `notes_fts` … FTS5 を使った全文検索（FTS同期トリガー付き）
- archive は検索対象外のため：
  - `note_embeddings` / `notes_fts` を持たない

### 1-2-1. `migrate_other.go` 実装状況（2026-03-06 更新）

**実装完了**:
- `go/db/migrate_other.go` の `migrateChronicle` / `migrateMemopedia` / `migrateArchive` を実装済み。
- 各ストアを個別に開いてマイグレーション後、ATTACH する方式を採用。
- `PRAGMA user_version` による冪等性保証（version >= 1 でスキップ）。
- 適用順序: `PRAGMA foreign_keys` → `notes` → インデックス → `notes_fts`（採用時）→ FTSトリガー → `tags`/`note_tags` → `note_embeddings`（採用時）→ `*_meta` → `lineage`。

**ストア別スキーマ構成**:
| Store | notes | tags | note_tags | note_embeddings | notes_fts | FTSトリガー | `*_meta` | lineage | `working_scope` | `is_pinned` |
| --- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| short | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | - | - |
| chronicle | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ (必須) | ✓ (任意) |
| memopedia | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ (必須) | ✓ (任意) |
| archive | ✓ | ✓ | ✓ | - | - | - | ✓ | ✓ | - | - |

**インデックス（全ストア共通）**:
- `idx_notes_created_at`, `idx_notes_last_accessed_at`
- `idx_notes_source_trust`, `idx_notes_sensitivity`
- `idx_tags_name`, `idx_tags_parent`
- `idx_note_tags_tag_id`, `idx_note_tags_note_id`
- `idx_lineage_src`, `idx_lineage_dest`
- chronicle/memopedia 追加: `idx_notes_working_scope`, `idx_notes_is_pinned`

### 1-2-2. short.sql との差分方針（store 別）

| Store | notes | tags | note_tags | note_embeddings | notes_fts | `working_scope` | `is_pinned` | 初期 `user_version` |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| short | 採用 | 採用 | 採用 | 採用 | 採用 | 未採用 | 未採用 | 1 |
| chronicle | 採用 | 採用 | 採用 | 採用 | 採用 | v1 で必須列として採用 | v1 で任意列として採用（`0/1`） | 1 |
| memopedia | 採用 | 採用 | 採用 | 採用 | 採用 | v1 で必須列として採用 | v1 で任意列として採用（`0/1`） | 1 |
| archive | 採用 | 採用 | 採用 | 不採用 | 不採用 | 未採用 | 未採用 | 1 |

- 差分 DDL は `schema/chronicle.sql` / `schema/memopedia.sql` / `schema/archive.sql` に明記する。
- `working_scope` / `is_pinned` の導入対象は chronicle・memopedia のみとし、archive には導入しない。

### 1-2-3. store 別要求（short / chronicle / memopedia / archive）

本節の要件IDは、Task Seed 作成時に **1要件=1タスク（0.5d 粒度）** で起票可能な最小単位として固定する。

#### short store 要求

| 区分 | MUST (v1) | SHOULD (v1.x) | FUTURE (v2+) | 許可変更 | 禁止変更 |
| --- | --- | --- | --- | --- | --- |
| short | `REQ-STORE-SHORT-001`: `mem in short` で `short.notes` へ保存できること。 | `REQ-STORE-SHORT-002`: `mem gc short --dry-run` で削除対象候補のみ返せること（feature flag 有効時）。 | `REQ-STORE-SHORT-9XX`: decay 学習ベースの自動優先度調整。 | 任意列追加、任意CLIオプション追加、dry-run の診断項目追加（後方互換維持）。 | 既存保存項目の必須化変更、既定挙動での自動削除、`--json` 出力キー破壊。 |

| Requirement ID | 入力条件 | 出力条件 | エラー条件 | 非機能条件 | Done条件 |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-SHORT-001` | `mem in short --title <text> [--body <text> | --stdin]` または同等 API 入力（`title`/`body`）で `body` 非空。 | `short.notes` に1件挿入され、`id`/`created_at` を返す。 | `body` 空文字は `INVALID_ARGUMENT`。保存失敗は `INTERNAL`。 | ingest p95 が `REQ-NFR-001` 閾値以内。 | `trace-req-api-001` 相当の ingest 検証 + `short.notes` 実在確認。 |
| `REQ-STORE-SHORT-002` | `mem.features.gc_short=true` かつ `mem gc short --dry-run` 実行。 | 削除候補件数と対象IDを返し、DB更新は0件。 | feature flag 無効時（route 公開時）は `INTERNAL`（500）を返す。 | dry-run 実行で `short` ロック時間が運用閾値（5秒）以内。 | dry-run 前後の `short.notes` 件数不変 + レポート出力確認。 |

#### chronicle store 要求

| 区分 | MUST (v1) | SHOULD (v1.x) | FUTURE (v2+) | 許可変更 | 禁止変更 |
| --- | --- | --- | --- | --- | --- |
| chronicle | `REQ-STORE-CHR-001`: `working_scope` 必須で `chronicle.notes` 保存できること。 | `REQ-STORE-CHR-002`: タグ/期間による抽出ビューを追加可能。 | `REQ-STORE-CHR-9XX`: 時系列クラスタ自動要約。 | 非破壊インデックス追加、任意フィルタ追加、任意メタ列追加。 | `working_scope` 必須性の撤回、時系列ソート既定の反転、既存ID体系変更。 |

| Requirement ID | 入力条件 | 出力条件 | エラー条件 | 非機能条件 | Done条件 |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-CHR-001` | `dest_scope=chronicle` かつ `working_scope` 指定で昇格/保存。 | `chronicle.notes` に保存され、`working_scope` が欠落しない。 | `working_scope` 未指定は `INVALID_ARGUMENT`。 | 保存 + 検索の往復 p95 が `REQ-NFR-001` 閾値以内。 | 保存後 `mem out show`（または API）で `working_scope` が取得できる。 |
| `REQ-STORE-CHR-002` | `from/to` 期間またはタグ指定で検索実行。 | 条件一致ノートのみ時系列降順で返す。 | 不正期間（`from > to`）は `INVALID_ARGUMENT`。 | 10,000件条件で検索 p95 が 1.5 秒以内。 | 期間フィルタ有/無の差分テストで期待件数一致。 |

#### memopedia store 要求

| 区分 | MUST (v1) | SHOULD (v1.x) | FUTURE (v2+) | 許可変更 | 禁止変更 |
| --- | --- | --- | --- | --- | --- |
| memopedia | `REQ-STORE-MP-001`: 用語/方針ノートを `memopedia.notes` に保存できること。 | `REQ-STORE-MP-002`: 既存ページへの追記（reflect）を追加可能。 | `REQ-STORE-MP-9XX`: 知識グラフ自動リンク生成。 | 任意セクション追加、任意参照メタ追加、後方互換な検索キー追加。 | 既存ノート本文の自動上書き、同一ID再利用による履歴破壊、必須列削除。 |

| Requirement ID | 入力条件 | 出力条件 | エラー条件 | 非機能条件 | Done条件 |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-MP-001` | `dest_scope=memopedia` で `title` と `body` を指定。 | `memopedia.notes` に保存され、検索対象に含まれる。 | `title` / `body` 欠落は `INVALID_ARGUMENT`。保存失敗は `INTERNAL`。 | 保存後の全文検索 p95 が 1.0 秒以内。 | 保存→検索で新規ノートIDがヒットすることを確認。 |

> 互換性注記: 旧用語 `content` は廃止表現とし、本仕様では `body` へ統一する。
| `REQ-STORE-MP-002` | 既存 `memopedia` ノートID + 追記本文で reflect 実行。 | 追記結果が新規ノートまたは版管理付きで保持される。 | 対象ID未存在は `NOT_FOUND`。競合は `CONFLICT`（未実装時 `INTERNAL`）。 | 追記処理でデータ喪失ゼロ（元本文復元可能）。 | reflect 後に旧版参照可 + 新版検索可を確認。 |

#### archive store 要求

| 区分 | MUST (v1) | SHOULD (v1.x) | FUTURE (v2+) | 許可変更 | 禁止変更 |
| --- | --- | --- | --- | --- | --- |
| archive | `REQ-STORE-ARC-001`: short からの退避ノートを `archive.notes` へ保持できること。 | `REQ-STORE-ARC-002`: retention 期限切れ候補の監査付き purge 提示。 | `REQ-STORE-ARC-9XX`: 暗号化アーカイブ/復号キー分離。 | 保持メタ列追加、監査ログ項目追加、dry-run 出力列追加。 | 監査ログなし purge 実行、retention 無視削除、lineage 未記録での short 削除。 |

| Requirement ID | 入力条件 | 出力条件 | エラー条件 | 非機能条件 | Done条件 |
| --- | --- | --- | --- | --- | --- |
| `REQ-STORE-ARC-001` | GC 判定で `archive` 退避対象となった `short` ノート。 | `archive.notes` 保存 + `lineage(archived_from)` 記録後にのみ short 削除可能。 | `archive` 保存失敗時は short 削除禁止、`INTERNAL` 返却。 | 補償処理込みで `REQ-NFR-005`（30分以内収束）を満たす。 | `archive 実在 + lineage 実在` を満たすケースでのみ削除判定成立。 |
| `REQ-STORE-ARC-002` | retention 超過候補を `--dry-run` で照会。 | 候補ID/期限/想定削除件数を返し、監査ログ下書きを生成。 | 監査コンテキスト不足は `INVALID_ARGUMENT`。 | purge dry-run p95 が 2.0 秒以内。 | dry-run 結果と監査ログ下書きが同一件数で一致。 |

#### 1-2-3-a. Task Seed 分解ルール（0.5d 固定）

- `REQ-STORE-*` は 1 ID を 1 Task Seed として起票する。
- 1 Task Seed の想定工数は `0.5d`（実装 + 検証 + 文書更新）を上限とする。
- 依存がある場合は `Depends on: REQ-...` で要件IDを明記し、同一 Task Seed 内で複数要件を混在させない。
- 2 つ以上の要件を同時に変更する場合は、Task Seed を要件ID単位に分割して並行可否を明示する。

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

**実装状況（2026-03-06 更新）**:
- `go/db/gatekeeper.go`: インターフェース定義（`Gatekeeper`, `GatekeeperCheckRequest`, `GatekeeperDecision` 等）
- `go/db/gatekeeper_impl.go`: `DefaultGatekeeper` 実装完了
  - プロファイル: `STRICT` / `NORMAL` / `DEV`
  - 判定結果: `allow` / `deny` / `needs_human`
  - fail-closed: `sensitivity=secret` は常に `deny`
  - テスト用ヘルパー: `AllowAllGatekeeper`, `DenyAllGatekeeper`

**プロファイル別動作**:
| Profile | `sensitivity=secret` | `source_trust=untrusted` | `source_trust=user_input` | `source_trust=trusted` |
| --- | --- | --- | --- | --- |
| `DEV` | deny | allow | allow | allow |
| `NORMAL` | deny | needs_human | allow | allow |
| `STRICT` | deny | needs_human | needs_human | allow |

**Service層統合**:
- `service/service.go`: `IngestShort` で `Gatekeeper.Check` を呼び出し
- `service/errors.go`: `ErrPolicyDenied`, `ErrNeedsHuman` 追加
- `api/errors.go`: `GATEKEEP_DENY` エラーコードへマッピング

**入力バリデーション（実装済み）**:
- `title`: 最大500文字
- `body`: 最大100,000文字
- `source_type`: `web` | `file` | `chat` | `agent` | `manual`
- `source_trust`: `trusted` | `user_input` | `untrusted`
- `sensitivity`: `public` | `internal` | `secret`
- `query` (search): 最大1,000文字
- `top_k` (search): 最大100