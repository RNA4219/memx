-- schema/memopedia.sql
-- knowledge base store
-- short.sql との差分:
-- - notes.working_scope を追加（知識領域の分類キー）
-- - notes.is_pinned を追加（定義・方針の固定表示）
-- - memopedia_meta を追加（GC用メタ情報）

PRAGMA foreign_keys = ON;

-- =========================================
-- 1. notes（知識ベースノート）
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
  sensitivity       TEXT NOT NULL,         -- 'public' | 'internal' | 'confidential' | 'secret'

  relevance         REAL,                  -- 0〜1, null 可
  quality           REAL,
  novelty           REAL,
  importance_static REAL,

  -- タグ route を上書きするための昇格先指定
  route_override    TEXT,

  -- memopedia 固有列
  working_scope     TEXT NOT NULL,         -- 'self' | 'session:<id>' | 'project:<name>' | 'global'
  is_pinned         INTEGER NOT NULL DEFAULT 0  -- 1 なら Working Memory として固定
);

CREATE INDEX IF NOT EXISTS idx_notes_created_at
  ON notes(created_at);

CREATE INDEX IF NOT EXISTS idx_notes_last_accessed_at
  ON notes(last_accessed_at);

CREATE INDEX IF NOT EXISTS idx_notes_source_trust
  ON notes(source_trust);

CREATE INDEX IF NOT EXISTS idx_notes_sensitivity
  ON notes(sensitivity);

CREATE INDEX IF NOT EXISTS idx_notes_working_scope
  ON notes(working_scope);

CREATE INDEX IF NOT EXISTS idx_notes_is_pinned
  ON notes(is_pinned);


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
  route         TEXT NOT NULL,               -- 'chronicle' | 'memopedia' | 'both' | 'short_only'
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
-- 5. memopedia_meta（GC 用メタ情報）
-- =========================================

CREATE TABLE IF NOT EXISTS memopedia_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

-- 初期値

INSERT OR IGNORE INTO memopedia_meta(key, value)
  VALUES ('note_count', '0');

INSERT OR IGNORE INTO memopedia_meta(key, value)
  VALUES ('token_sum', '0');

INSERT OR IGNORE INTO memopedia_meta(key, value)
  VALUES ('last_gc_at', '1970-01-01T00:00:00Z');


-- =========================================
-- 6. lineage（蒸留・昇格・隔離の系譜）
-- =========================================

CREATE TABLE IF NOT EXISTS lineage (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  src_store      TEXT NOT NULL,      -- 'short' / 'chronicle' / 'memopedia'
  src_note_id    TEXT NOT NULL,
  dest_store     TEXT NOT NULL,      -- 'chronicle' / 'memopedia' / 'archive'
  dest_note_id   TEXT NOT NULL,
  relation       TEXT NOT NULL,      -- 'distilled_to' / 'merged_into' / 'observed' / 'reflected' / 'archived_from' 等
  created_at     TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lineage_src
  ON lineage(src_store, src_note_id);

CREATE INDEX IF NOT EXISTS idx_lineage_dest
  ON lineage(dest_store, dest_note_id);


-- =========================================
-- 7. スキーマバージョン
-- =========================================

PRAGMA user_version = 1;