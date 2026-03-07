package db

import (
	"database/sql"
	"fmt"
)

// migrateChronicle は chronicle.db に対してスキーマを適用する。
func migrateChronicle(db *sql.DB) error {
	// user_version をチェックして、既にマイグレーション済みならスキップ
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		return fmt.Errorf("check user_version: %w", err)
	}
	if version >= 1 {
		return nil
	}

	// DDL を実行
	ddls := getChronicleDDL()
	for _, ddl := range ddls {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("apply chronicle schema: %w (ddl: %s)", err, truncate(ddl, 50))
		}
	}

	// user_version を設定
	if _, err := db.Exec("PRAGMA user_version = 1;"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}

// migrateMemopedia は memopedia.db に対してスキーマを適用する。
func migrateMemopedia(db *sql.DB) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		return fmt.Errorf("check user_version: %w", err)
	}
	if version >= 1 {
		return nil
	}

	ddls := getMemopediaDDL()
	for _, ddl := range ddls {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("apply memopedia schema: %w (ddl: %s)", err, truncate(ddl, 50))
		}
	}

	if _, err := db.Exec("PRAGMA user_version = 1;"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}

// migrateArchive は archive.db に対してスキーマを適用する。
func migrateArchive(db *sql.DB) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		return fmt.Errorf("check user_version: %w", err)
	}
	if version >= 1 {
		return nil
	}

	ddls := getArchiveDDL()
	for _, ddl := range ddls {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("apply archive schema: %w (ddl: %s)", err, truncate(ddl, 50))
		}
	}

	if _, err := db.Exec("PRAGMA user_version = 1;"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}

// getChronicleDDL は chronicle 用の DDL 文を返す。
func getChronicleDDL() []string {
	return []string{
		`PRAGMA foreign_keys = ON;`,

		// notes テーブル（working_scope, is_pinned を追加）
		`CREATE TABLE IF NOT EXISTS notes (
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
);`,

		// インデックス
		`CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_last_accessed_at ON notes(last_accessed_at);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_source_trust ON notes(source_trust);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_sensitivity ON notes(sensitivity);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_working_scope ON notes(working_scope);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_is_pinned ON notes(is_pinned);`,

		// FTS5 仮想テーブル
		`CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
  title,
  body,
  content='notes',
  content_rowid='rowid'
);`,

		// FTS 同期トリガー
		`CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
  INSERT INTO notes_fts(rowid, title, body) VALUES (new.rowid, new.title, new.body);
END;`,
		`CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
  DELETE FROM notes_fts WHERE rowid = old.rowid;
  INSERT INTO notes_fts(rowid, title, body) VALUES (new.rowid, new.title, new.body);
END;`,
		`CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
  DELETE FROM notes_fts WHERE rowid = old.rowid;
END;`,

		// tags テーブル
		`CREATE TABLE IF NOT EXISTS tags (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  route       TEXT NOT NULL,
  parent_id   INTEGER,
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL,
  usage_count INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(parent_id) REFERENCES tags(id) ON DELETE SET NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_parent ON tags(parent_id);`,

		// note_tags テーブル
		`CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id ON note_tags(tag_id);`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_note_id ON note_tags(note_id);`,

		// note_embeddings テーブル
		`CREATE TABLE IF NOT EXISTS note_embeddings (
  note_id TEXT PRIMARY KEY,
  dim     INTEGER NOT NULL,
  vector  BLOB NOT NULL,
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);`,
		`CREATE INDEX IF NOT EXISTS idx_note_embeddings_dim ON note_embeddings(dim);`,

		// chronicle_meta テーブル
		`CREATE TABLE IF NOT EXISTS chronicle_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);`,
		`INSERT OR IGNORE INTO chronicle_meta(key, value) VALUES ('note_count', '0');`,
		`INSERT OR IGNORE INTO chronicle_meta(key, value) VALUES ('token_sum', '0');`,
		`INSERT OR IGNORE INTO chronicle_meta(key, value) VALUES ('last_gc_at', '1970-01-01T00:00:00Z');`,

		// lineage テーブル
		`CREATE TABLE IF NOT EXISTS lineage (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  src_store      TEXT NOT NULL,
  src_note_id    TEXT NOT NULL,
  dest_store     TEXT NOT NULL,
  dest_note_id   TEXT NOT NULL,
  relation       TEXT NOT NULL,
  created_at     TEXT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_src ON lineage(src_store, src_note_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_dest ON lineage(dest_store, dest_note_id);`,
	}
}

// getMemopediaDDL は memopedia 用の DDL 文を返す。
func getMemopediaDDL() []string {
	return []string{
		`PRAGMA foreign_keys = ON;`,

		// notes テーブル（working_scope, is_pinned を追加）
		`CREATE TABLE IF NOT EXISTS notes (
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
);`,

		// インデックス
		`CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_last_accessed_at ON notes(last_accessed_at);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_source_trust ON notes(source_trust);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_sensitivity ON notes(sensitivity);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_working_scope ON notes(working_scope);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_is_pinned ON notes(is_pinned);`,

		// FTS5 仮想テーブル
		`CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
  title,
  body,
  content='notes',
  content_rowid='rowid'
);`,

		// FTS 同期トリガー
		`CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
  INSERT INTO notes_fts(rowid, title, body) VALUES (new.rowid, new.title, new.body);
END;`,
		`CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
  DELETE FROM notes_fts WHERE rowid = old.rowid;
  INSERT INTO notes_fts(rowid, title, body) VALUES (new.rowid, new.title, new.body);
END;`,
		`CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
  DELETE FROM notes_fts WHERE rowid = old.rowid;
END;`,

		// tags テーブル
		`CREATE TABLE IF NOT EXISTS tags (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  route       TEXT NOT NULL,
  parent_id   INTEGER,
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL,
  usage_count INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(parent_id) REFERENCES tags(id) ON DELETE SET NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_parent ON tags(parent_id);`,

		// note_tags テーブル
		`CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id ON note_tags(tag_id);`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_note_id ON note_tags(note_id);`,

		// note_embeddings テーブル
		`CREATE TABLE IF NOT EXISTS note_embeddings (
  note_id TEXT PRIMARY KEY,
  dim     INTEGER NOT NULL,
  vector  BLOB NOT NULL,
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);`,
		`CREATE INDEX IF NOT EXISTS idx_note_embeddings_dim ON note_embeddings(dim);`,

		// memopedia_meta テーブル
		`CREATE TABLE IF NOT EXISTS memopedia_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);`,
		`INSERT OR IGNORE INTO memopedia_meta(key, value) VALUES ('note_count', '0');`,
		`INSERT OR IGNORE INTO memopedia_meta(key, value) VALUES ('token_sum', '0');`,
		`INSERT OR IGNORE INTO memopedia_meta(key, value) VALUES ('last_gc_at', '1970-01-01T00:00:00Z');`,

		// lineage テーブル
		`CREATE TABLE IF NOT EXISTS lineage (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  src_store      TEXT NOT NULL,
  src_note_id    TEXT NOT NULL,
  dest_store     TEXT NOT NULL,
  dest_note_id   TEXT NOT NULL,
  relation       TEXT NOT NULL,
  created_at     TEXT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_src ON lineage(src_store, src_note_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_dest ON lineage(dest_store, dest_note_id);`,
	}
}

// getArchiveDDL は archive 用の DDL 文を返す。
func getArchiveDDL() []string {
	return []string{
		`PRAGMA foreign_keys = ON;`,

		// notes テーブル（working_scope, is_pinned なし）
		`CREATE TABLE IF NOT EXISTS notes (
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
  route_override    TEXT
);`,

		// インデックス
		`CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_last_accessed_at ON notes(last_accessed_at);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_source_trust ON notes(source_trust);`,
		`CREATE INDEX IF NOT EXISTS idx_notes_sensitivity ON notes(sensitivity);`,

		// tags テーブル
		`CREATE TABLE IF NOT EXISTS tags (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  route       TEXT NOT NULL,
  parent_id   INTEGER,
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL,
  usage_count INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(parent_id) REFERENCES tags(id) ON DELETE SET NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_parent ON tags(parent_id);`,

		// note_tags テーブル
		`CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id ON note_tags(tag_id);`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_note_id ON note_tags(note_id);`,

		// archive_meta テーブル
		`CREATE TABLE IF NOT EXISTS archive_meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);`,
		`INSERT OR IGNORE INTO archive_meta(key, value) VALUES ('note_count', '0');`,
		`INSERT OR IGNORE INTO archive_meta(key, value) VALUES ('token_sum', '0');`,
		`INSERT OR IGNORE INTO archive_meta(key, value) VALUES ('last_gc_at', '1970-01-01T00:00:00Z');`,

		// lineage テーブル
		`CREATE TABLE IF NOT EXISTS lineage (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  src_store      TEXT NOT NULL,
  src_note_id    TEXT NOT NULL,
  dest_store     TEXT NOT NULL,
  dest_note_id   TEXT NOT NULL,
  relation       TEXT NOT NULL,
  created_at     TEXT NOT NULL
);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_src ON lineage(src_store, src_note_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_dest ON lineage(dest_store, dest_note_id);`,
	}
}

// truncate は文字列を指定長で切り詰める。
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}