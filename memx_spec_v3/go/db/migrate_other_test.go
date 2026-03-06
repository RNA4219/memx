package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrateChronicle(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "chronicle.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	if err := migrateChronicle(db); err != nil {
		t.Fatalf("migrateChronicle failed: %v", err)
	}

	// user_version が設定されているか確認
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}

	// notes テーブルが存在するか確認（sqlite_master で確認）
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='notes';").Scan(&tableName)
	if err != nil {
		t.Errorf("notes table should exist: %v", err)
	}
	if tableName != "notes" {
		t.Errorf("expected table name 'notes', got: %s", tableName)
	}

	// tags テーブルが存在するか確認
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tags';").Scan(&tableName)
	if err != nil {
		t.Errorf("tags table should exist: %v", err)
	}
	if tableName != "tags" {
		t.Errorf("expected table name 'tags', got: %s", tableName)
	}
}

func TestMigrateChronicle_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "chronicle.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// 1回目のマイグレーション
	if err := migrateChronicle(db); err != nil {
		t.Fatalf("first migrateChronicle failed: %v", err)
	}

	// 2回目のマイグレーション（再実行安全性の確認）
	if err := migrateChronicle(db); err != nil {
		t.Fatalf("second migrateChronicle failed: %v", err)
	}

	// user_version が 1 のままであることを確認
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestMigrateMemopedia(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "memopedia.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateMemopedia(db); err != nil {
		t.Fatalf("migrateMemopedia failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestMigrateMemopedia_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "memopedia.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateMemopedia(db); err != nil {
		t.Fatalf("first migrateMemopedia failed: %v", err)
	}

	if err := migrateMemopedia(db); err != nil {
		t.Fatalf("second migrateMemopedia failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestMigrateArchive(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "archive.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateArchive(db); err != nil {
		t.Fatalf("migrateArchive failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}

	// archive には FTS がないことを確認
	// notes_fts テーブルは存在しないはず
	var exists int
	err = db.QueryRow("SELECT 1 FROM notes_fts LIMIT 1;").Scan(&exists)
	if err == nil {
		t.Error("archive should not have notes_fts table")
	}
}

func TestMigrateArchive_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "archive.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateArchive(db); err != nil {
		t.Fatalf("first migrateArchive failed: %v", err)
	}

	if err := migrateArchive(db); err != nil {
		t.Fatalf("second migrateArchive failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestOpenAll(t *testing.T) {
	tmpDir := t.TempDir()

	paths := Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Chronicle: filepath.Join(tmpDir, "chronicle.db"),
		Memopedia: filepath.Join(tmpDir, "memopedia.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	conn, err := OpenAll(paths)
	if err != nil {
		t.Fatalf("OpenAll failed: %v", err)
	}
	defer conn.Close()

	// 各ストアの user_version を確認
	stores := []struct {
		name  string
		query string
	}{
		{"main", "PRAGMA main.user_version;"},
		{"chronicle", "PRAGMA chronicle.user_version;"},
		{"memopedia", "PRAGMA memopedia.user_version;"},
		{"archive", "PRAGMA archive.user_version;"},
	}

	for _, store := range stores {
		var version int
		if err := conn.DB.QueryRow(store.query).Scan(&version); err != nil {
			t.Errorf("failed to check %s.user_version: %v", store.name, err)
		} else if version != 1 {
			t.Errorf("expected %s.user_version=1, got: %d", store.name, version)
		}
	}
}

func openDB(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite", dsn)
}