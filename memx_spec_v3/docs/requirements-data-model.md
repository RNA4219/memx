---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - データモデル

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

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

本節を `sensitivity`・保持・削除・監査証跡の**正本要件**とする。運用手順は [付録 G: Security & Privacy](../../docs/addenda/G_Security_Privacy.md) と [memx最小版: セキュリティ運用補足](../../docs/security/minimal_operations.md#3-gatekeeper-deny--needs_human-発生時の運用手順) を参照し、重複定義は行わない。
- Requirement ID: `REQ-SEC-001`（Security）
- Requirement ID: `REQ-RET-001`（Retention）

#### 2-7-1. `sensitivity` 別の保存可否・マスキング・保持期間

| sensitivity | 永続保存可否 | マスキング要件 | 保持期間要件 |
| --- | --- | --- | --- |
| `public` | `short/chronicle/memopedia/archive` へ保存可 | 原則マスク不要。ただし認証情報パターン検知時は `***REDACTED***` へ置換して保存/出力 | `short` 退避後 `archive` で最大 365 日保持し、期限到達後は物理削除対象 |
| `internal` | ローカル環境の `short/chronicle/memopedia/archive` のみ保存可（外部共有禁止） | CLI/API 出力および運用ログで識別子を部分マスク（例: `a***@example.com`） | `short` 退避後 `archive` で最大 180 日保持し、期限到達後は物理削除対象 |
| `secret` | 永続保存禁止（`memory_store` で `deny`） | 保存前・出力前とも全量マスク（`***REDACTED***`）し、原文は監査ログにも残さない | 保持対象外。検知時はイベントメタデータのみ最小 90 日保持 |

- `needs_human` は v1.3 では `deny` 同等で扱う（fail-closed）。
- 保持期間の起算日は `archive` へ退避した `created_at` とし、期限判定は日次バッチまたは `mem gc --purge` 実行時に行う。
- `internal` / `secret` の運用判断ログは、最小保持期間を 90 日とし、運用手順への固定リンクは [memx最小版: セキュリティ運用補足 §3](../../docs/security/minimal_operations.md#3-gatekeeper-deny--needs_human-発生時の運用手順) とする。

#### 2-7-2. actor / approval / audit 責任分界表（2-7-1〜2-7-5）

| 対象節 | operation | actor（実行主体） | approval（承認要否） | audit（監査ログ必須項目） |
| --- | --- | --- | --- | --- |
| 2-7-1 | `memory_store` | API/Service（Gatekeeper フック必須） | 不要（判定は Gatekeeper が自動実施） | `op`, `decision`, `sensitivity`, `executor`, `executed_at`, `result` |
| 2-7-1 | `memory_output` | API/CLI 出力処理（Gatekeeper フック必須） | 不要（判定は Gatekeeper が自動実施） | `op`, `decision`, `sensitivity`, `executor`, `executed_at`, `result` |
| 2-7-3 | `archive_move` | 自動: `mem gc` 定期ジョブ / 手動: Repo Maintainer または委任運用担当 | 手動時は Repo Maintainer の実行責任を必須 | `op`, `src_note_id`, `dest_note_id`, `sensitivity`, `executor`, `executed_at`, `result`, `reason`, `retryable`, `owner`, `next_attempt_at` |
| 2-7-4 | `archive_purge` | 自動: Maintainer 設定の定期ジョブ / 手動: Repo Maintainer（必要時 Security Champion へ事後報告） | 手動実行時のみ必須（Repo Maintainer） | `op`, `note_id`, `sensitivity`, `retention_expired_at`, `executor`, `executed_at`, `result`, `reason`, `retryable`, `owner`, `next_attempt_at` |
| 2-7-5 | fail-closed 整合チェック | 要件更新者（PR Author）+ Repo Maintainer（レビュー責任） | マージ前レビュー承認を必須 | `check_id`, `guardrail_ref`, `requirement_ref`, `status`, `owner`, `checked_at`, `evidence` |

`archive_move` / `archive_purge` の監査ログ必須項目は、成功/失敗を問わず `result`, `reason`, `retryable`, `owner`, `next_attempt_at` を固定（追加任意・削除禁止）とする。

`REQ-SEC-AUD-*` の固定フィールド定義（削除/名称変更禁止）:
- `evidence_file_path`: 証跡ファイルパス。`artifacts/ops/` 配下の実在パスを記録する。
- `required_keys`: 当該operationの監査ログに必須なキー配列。最低限 `result`, `reason`, `retryable`, `owner`, `next_attempt_at` を含む。
- `retention_days`: 監査証跡の保持期間（日）。`archive_move` / `archive_purge` とも `90` 以上を必須とする。

- `REQ-SEC-AUD-001` (`archive_move`): 上記固定項目と固定フィールドを満たす。
- `REQ-SEC-AUD-002` (`archive_purge`): 上記固定項目と固定フィールドを満たす。

#### 2-7-3. `archive` 退避条件（Short → Archive）

| actor | approval | audit |
| --- | --- | --- |
| 自動: `mem gc` 定期ジョブ / 手動: Repo Maintainer または委任運用担当 | 手動時は Repo Maintainer の実行責任を必須 | `op`, `src_note_id`, `dest_note_id`, `sensitivity`, `executor`, `executed_at`, `result`, `reason`, `retryable`, `owner`, `next_attempt_at`（`result/reason/retryable/owner/next_attempt_at` は固定） |

`phase3` の退避実行は、次の条件を**すべて満たす**場合のみ許可する。

1. GC 候補選定済みである（`memory_policy.yaml.gc.short` の閾値判定を通過）。
2. `route_override` やタグルーティングで `chronicle/memopedia/both` への強制昇格が指定されていない。
3. `sensitivity` が `public` または `internal` である（`secret` は対象外）。
4. `archive` 側 Insert と `lineage(relation='archived_from')` の記録に成功している。
5. 失敗時は `REQ-SEC-AUD-001` を満たし、次回再試行計画を監査ログへ残す。
6. `result`, `reason`, `retryable`, `owner`, `next_attempt_at` は成功/失敗を問わず必須（固定）で記録する。

#### 2-7-4. 物理削除条件（Archive Purge）

| actor | approval | audit |
| --- | --- | --- |
| 自動: Maintainer 設定の定期ジョブ / 手動: Repo Maintainer（必要時 Security Champion へ事後報告） | 手動実行時のみ必須（Repo Maintainer） | `op`, `note_id`, `sensitivity`, `retention_expired_at`, `executor`, `executed_at`, `result`, `reason`, `retryable`, `owner`, `next_attempt_at`（`result/reason/retryable/owner/next_attempt_at` は固定） |

`archive` からの物理削除は、次の条件を**すべて満たす**場合のみ許可する。

1. 2-7-1 の保持期間を超過している。
2. 対象ノートにインシデント調査中フラグ（legal/investigation hold）が無い。
3. 削除直前に `lineage` と `archive.notes` の対応整合を確認済みである。
4. 削除不能時は `REQ-SEC-AUD-002` を満たし、次回削除計画を監査ログへ残す。
5. `result`, `reason`, `retryable`, `owner`, `next_attempt_at` は成功/失敗を問わず必須（固定）で記録する。

#### 2-7-5. GUARDRAILS fail-closed との整合チェック要件

| actor | approval | audit |
| --- | --- | --- |
| 要件更新者（PR Author）+ Repo Maintainer レビュー | マージ前レビュー承認を必須 | `check_id`, `guardrail_ref`, `requirement_ref`, `status`, `owner`, `checked_at`, `evidence` |

整合チェックは次の `REQ-SEC-GRD-001` を必須とする。

- `REQ-SEC-GRD-001-1`: [GUARDRAILS](../../GUARDRAILS.md) の fail-closed 方針（`needs_human` は `deny` 相当）と 2-7-1 の判定規則が一致している。
- `REQ-SEC-GRD-001-2`: `memory_store` / `memory_output` / `archive_move` / `archive_purge` の 4 operation が、2-7-2 の actor/approval/audit と矛盾しない。
- `REQ-SEC-GRD-001-3`: 差分が生じた場合、requirements を正本として同一PRで GUARDRAILS と運用手順書（`docs/security/minimal_operations.md`）を追従更新する。