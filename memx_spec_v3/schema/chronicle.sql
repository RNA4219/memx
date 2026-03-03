-- schema/chronicle.sql
-- short.sql との差分:
-- - 必須列: notes.working_scope を追加（Chronicle の時間軸グルーピング用）
-- - 任意列: notes.is_pinned を追加（重要ログ固定化）
-- - FTS: notes_fts を有効

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS notes (
  id                TEXT PRIMARY KEY,
  title             TEXT NOT NULL,
  summary           TEXT NOT NULL DEFAULT '',
  body              TEXT NOT NULL,
  created_at        TEXT NOT NULL,
  updated_at        TEXT NOT NULL,
  last_accessed_at  TEXT NOT NULL,
  access_count      INTEGER NOT NULL DEFAULT 0,
  source_type       TEXT NOT NULL,
  origin            TEXT NOT NULL DEFAULT '',
  source_trust      TEXT NOT NULL,
  sensitivity       TEXT NOT NULL,
  relevance         REAL,
  quality           REAL,
  novelty           REAL,
  importance_static REAL,
  route_override    TEXT,
  working_scope     TEXT NOT NULL,
  is_pinned         INTEGER NOT NULL DEFAULT 0
);

CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
  title,
  body,
  content='notes',
  content_rowid='rowid'
);

CREATE TABLE IF NOT EXISTS tags (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  route       TEXT NOT NULL,
  parent_id   INTEGER,
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL,
  usage_count INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(parent_id) REFERENCES tags(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS note_embeddings (
  note_id TEXT PRIMARY KEY,
  dim     INTEGER NOT NULL,
  vector  BLOB NOT NULL,
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);

PRAGMA user_version = 1;
