package db

import (
	"database/sql"
	"fmt"
)

// shortSchemaSQL は schema/short.sql と同等の DDL を含む。
const shortSchemaSQL = `-- schema/short.sql

PRAGMA foreign_keys = ON;

-- =========================================
-- 1. notes（短期ノート本体）
-- =========================================

CREATE TABLE IF NOT EXISTS notes (
  id                TEXT PRIMARY KEY,      -- UUID 等

  title             TEXT NOT NULL,
  summary           TEXT NOT NULL DEFAULT '',
  body              TEXT NOT NULL,

  created_at        TEXT NOT NULL,         -- ISO8601
  updated_at        TEXT NOT NULL,
  last_accessed_at  TEXT NOT NULL,
  access_count      INTEGER NOT NULL DEFAULT 0,

  source_type       TEXT NOT NULL,         -- 'web' | 'file' | 'chat' | 'agent' | 'manual'
  origin            TEXT NOT NULL DEFAULT '', -- URL / path / agent 名など（無い場合は空文字）
  source_trust      TEXT NOT NULL,         -- 'trusted' | 'user_input' | 'untrusted'
  sensitivity       TEXT NOT NULL,         -- 'public' | 'internal' | 'secret'

  relevance         REAL,                  -- 0〜1, null 可
  quality           REAL,
  novelty           REAL,
  importance_static REAL,

  -- タグ route を上書きするための昇格先指定
  -- null or 'journal' | 'knowledge' | 'both' | 'archive_only'
  route_override    TEXT
);

CREATE INDEX IF NOT EXISTS idx_notes_created_at
  ON notes(created_at);

CREATE INDEX IF NOT EXISTS idx_notes_last_accessed_at
  ON notes(last_accessed_at);

CREATE INDEX IF NOT EXISTS idx_notes_source_trust
  ON notes(source_trust);

CREATE INDEX IF NOT EXISTS idx_notes_sensitivity
  ON notes(sensitivity);


-- =========================================
-- 2. notes_fts（全文検索用 FTS5）
--    content='notes' 方式
-- =========================================

CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
  title,
  body,
  content='notes',
  content_rowid='rowid'
);

-- notes → notes_fts 同期用トリガ
-- UPDATE 時は DELETE → INSERT で同期する（FTS5 の推奨パターン）。

CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
  INSERT INTO notes_fts(rowid, title, body)
  VALUES (new.rowid, new.title, new.body);
END;

CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
  DELETE FROM notes_fts WHERE rowid = old.rowid;
  INSERT INTO notes_fts(rowid, title, body)
  VALUES (new.rowid, new.title, new.body);
END;

CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
  DELETE FROM notes_fts WHERE rowid = old.rowid;
END;


-- =========================================
-- 3. tags / note_tags（タグと対応付け）
-- =========================================

CREATE TABLE IF NOT EXISTS tags (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT NOT NULL UNIQUE,        -- 正規化タグ名
  route         TEXT NOT NULL,               -- 'journal' | 'knowledge' | 'both' | 'short_only'
  parent_id     INTEGER,                     -- 代表タグ or 親タグ（NULL 許容）
  created_at    TEXT NOT NULL,
  updated_at    TEXT NOT NULL,
  usage_count   INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(parent_id) REFERENCES tags(id)
    ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_tags_name
  ON tags(name);

CREATE INDEX IF NOT EXISTS idx_tags_parent
  ON tags(parent_id);


CREATE TABLE IF NOT EXISTS note_tags (
  note_id       TEXT NOT NULL,
  tag_id        INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id)
    ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id
  ON note_tags(tag_id);

CREATE INDEX IF NOT EXISTS idx_note_tags_note_id
  ON note_tags(note_id);


-- =========================================
-- 4. note_embeddings（埋め込みベクトル）
-- =========================================

CREATE TABLE IF NOT EXISTS note_embeddings (
  note_id   TEXT PRIMARY KEY,
  dim       INTEGER NOT NULL,
  vector    BLOB NOT NULL,     -- float32 配列等（長さ dim * 4 バイト）
  FOREIGN KEY (note_id) REFERENCES notes(id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_note_embeddings_dim
  ON note_embeddings(dim);


-- =========================================
-- 5. short_meta（GC 用メタ情報）
-- =========================================

CREATE TABLE IF NOT EXISTS short_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

-- 初期値（既存行があれば無視される）

INSERT OR IGNORE INTO short_meta(key, value)
  VALUES ('note_count', '0');

INSERT OR IGNORE INTO short_meta(key, value)
  VALUES ('token_sum', '0');

INSERT OR IGNORE INTO short_meta(key, value)
  VALUES ('last_gc_at', '1970-01-01T00:00:00Z');


-- =========================================
-- 6. lineage（蒸留・昇格・隔離の系譜）
-- =========================================

CREATE TABLE IF NOT EXISTS lineage (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  src_store      TEXT NOT NULL,      -- 'short' / 'journal' / 'knowledge'
  src_note_id    TEXT NOT NULL,
  dest_store     TEXT NOT NULL,      -- 'journal' / 'knowledge' / 'archive'
  dest_note_id   TEXT NOT NULL,
  relation       TEXT NOT NULL,      -- 'distilled_to' / 'merged_into' / 'observed' / 'reflected' / 'archived_from' 等
  created_at     TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lineage_src
  ON lineage(src_store, src_note_id);

CREATE INDEX IF NOT EXISTS idx_lineage_dest
  ON lineage(dest_store, dest_note_id);


-- =========================================
-- 7. resolver_documents / resolver_chunks
-- =========================================

CREATE TABLE IF NOT EXISTS resolver_documents (
  doc_id TEXT PRIMARY KEY,
  doc_type TEXT NOT NULL,
  title TEXT NOT NULL,
  source_path TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  summary TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL,
  tags_json TEXT NOT NULL DEFAULT '[]',
  feature_keys_json TEXT NOT NULL DEFAULT '[]',
  task_ids_json TEXT NOT NULL DEFAULT '[]',
  acceptance_criteria_json TEXT NOT NULL DEFAULT '[]',
  forbidden_patterns_json TEXT NOT NULL DEFAULT '[]',
  definition_of_done_json TEXT NOT NULL DEFAULT '[]',
  dependencies_json TEXT NOT NULL DEFAULT '[]',
  importance TEXT NOT NULL DEFAULT 'reference'
);

CREATE INDEX IF NOT EXISTS idx_resolver_documents_type
  ON resolver_documents(doc_type);

CREATE INDEX IF NOT EXISTS idx_resolver_documents_version
  ON resolver_documents(version);

CREATE TABLE IF NOT EXISTS resolver_chunks (
  chunk_id TEXT PRIMARY KEY,
  doc_id TEXT NOT NULL,
  heading TEXT NOT NULL DEFAULT '',
  heading_path_json TEXT NOT NULL DEFAULT '[]',
  ordinal INTEGER NOT NULL,
  body TEXT NOT NULL,
  token_estimate INTEGER NOT NULL DEFAULT 0,
  importance TEXT NOT NULL DEFAULT 'reference',
  FOREIGN KEY (doc_id) REFERENCES resolver_documents(doc_id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_resolver_chunks_doc
  ON resolver_chunks(doc_id, ordinal);

CREATE TABLE IF NOT EXISTS resolver_document_links (
  src_doc_id TEXT NOT NULL,
  dst_doc_id TEXT NOT NULL,
  link_type TEXT NOT NULL,
  PRIMARY KEY (src_doc_id, dst_doc_id, link_type)
);

CREATE TABLE IF NOT EXISTS resolver_read_receipts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id TEXT NOT NULL,
  doc_id TEXT NOT NULL,
  version TEXT NOT NULL,
  chunk_ids_json TEXT NOT NULL DEFAULT '[]',
  reader TEXT NOT NULL,
  read_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_resolver_read_receipts_task
  ON resolver_read_receipts(task_id, read_at);

-- =========================================
-- 8. スキーマバージョン
-- =========================================

PRAGMA user_version = 1;
`

// migrateShort は short.db に対してスキーマを適用する。
func migrateShort(db *sql.DB) error {
	if _, err := db.Exec(shortSchemaSQL); err != nil {
		return fmt.Errorf("apply short schema: %w", err)
	}
	return nil
}
