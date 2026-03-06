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

		// FTS5 仮想テーブル
		`CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
  title,
  body,
  content='notes',
  content_rowid='rowid'
);`,

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

		// note_tags テーブル
		`CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`,

		// note_embeddings テーブル
		`CREATE TABLE IF NOT EXISTS note_embeddings (
  note_id TEXT PRIMARY KEY,
  dim     INTEGER NOT NULL,
  vector  BLOB NOT NULL,
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);`,
	}
}

// getMemopediaDDL は memopedia 用の DDL 文を返す。
func getMemopediaDDL() []string {
	return []string{
		`PRAGMA foreign_keys = ON;`,

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

		`CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
  title,
  body,
  content='notes',
  content_rowid='rowid'
);`,

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

		`CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`,

		`CREATE TABLE IF NOT EXISTS note_embeddings (
  note_id TEXT PRIMARY KEY,
  dim     INTEGER NOT NULL,
  vector  BLOB NOT NULL,
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);`,
	}
}

// getArchiveDDL は archive 用の DDL 文を返す。
func getArchiveDDL() []string {
	return []string{
		`PRAGMA foreign_keys = ON;`,

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

		`CREATE TABLE IF NOT EXISTS note_tags (
  note_id TEXT NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);`,
	}
}

// truncate は文字列を指定長で切り詰める。
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}