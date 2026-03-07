-- schema/archive.sql
-- archive store
-- short.sql との差分:
-- - working_scope/is_pinned は採用しない（アーカイブは検索対象外）
-- - notes_fts は無効（アーカイブの通常検索対象外）
-- - note_embeddings は無効（アーカイブはベクトル検索対象外）
-- - archive_meta を追加

PRAGMA foreign_keys = ON;

-- =========================================
-- 1. notes（アーカイブノート）
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
-- 2. tags / note_tags（タグと対応付け）
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
-- 3. archive_meta（管理用メタ情報）
-- =========================================

CREATE TABLE IF NOT EXISTS archive_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

-- 初期値

INSERT OR IGNORE INTO archive_meta(key, value)
  VALUES ('note_count', '0');

INSERT OR IGNORE INTO archive_meta(key, value)
  VALUES ('token_sum', '0');

INSERT OR IGNORE INTO archive_meta(key, value)
  VALUES ('last_gc_at', '1970-01-01T00:00:00Z');


-- =========================================
-- 4. lineage（退避・削除の系譜）
-- =========================================

CREATE TABLE IF NOT EXISTS lineage (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  src_store      TEXT NOT NULL,      -- 'short' / 'chronicle' / 'memopedia'
  src_note_id    TEXT NOT NULL,
  dest_store     TEXT NOT NULL,      -- 'archive'
  dest_note_id   TEXT NOT NULL,
  relation       TEXT NOT NULL,      -- 'archived_from' / 'purged_from' 等
  created_at     TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lineage_src
  ON lineage(src_store, src_note_id);

CREATE INDEX IF NOT EXISTS idx_lineage_dest
  ON lineage(dest_store, dest_note_id);


-- =========================================
-- 5. スキーマバージョン
-- =========================================

PRAGMA user_version = 1;