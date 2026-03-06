package db

import (
	"database/sql"
	"fmt"
	"log"

	// Pure-Go SQLite driver. If you prefer CGO-based sqlite3, you can swap this
	// import and the driver name in sql.Open.
	_ "modernc.org/sqlite"
)

// MustOpenAll は short.db をメインとして開き、
// chronicle / memopedia / archive を ATTACH + migrate した Conn を返す。
func MustOpenAll(paths Paths) *Conn {
	c, err := OpenAll(paths)
	if err != nil {
		log.Fatalf("failed to open db(s): %v", err)
	}
	return c
}

// OpenAll は short.db をメインとして開き、
// chronicle / memopedia / archive を ATTACH + migrate した Conn を返す。
//
// v1.3 以降では、この関数を Service / API 層から呼び出し、
// CLI は直接 DB を触らない。
func OpenAll(paths Paths) (*Conn, error) {
	dsn := fmt.Sprintf("file:%s", paths.Short)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open short.db: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign_keys: %w", err)
	}
	if _, err := db.Exec(`PRAGMA journal_mode = WAL;`); err != nil {
		// WAL が使えない環境もあるため warning 扱い。
		log.Printf("warning: failed to set WAL mode: %v", err)
	}

	if err := migrateShort(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate short.db: %w", err)
	}

	// 他ストアを個別に開いてマイグレーション後、ATTACH する
	if paths.Chronicle != "" {
		if err := migrateAndAttach(db, paths.Chronicle, "chronicle", migrateChronicle); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("migrate/attach chronicle.db: %w", err)
		}
	}
	if paths.Memopedia != "" {
		if err := migrateAndAttach(db, paths.Memopedia, "memopedia", migrateMemopedia); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("migrate/attach memopedia.db: %w", err)
		}
	}
	if paths.Archive != "" {
		if err := migrateAndAttach(db, paths.Archive, "archive", migrateArchive); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("migrate/attach archive.db: %w", err)
		}
	}

	return &Conn{DB: db}, nil
}

// Close は基盤となる *sql.DB をクローズする。
func (c *Conn) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// migrateAndAttach は DB ファイルを個別に開いてマイグレーションし、その後 ATTACH する。
// これにより、各ストアのスキーマを独立して管理できる。
func migrateAndAttach(mainDB *sql.DB, path string, schemaName string, migrateFunc func(*sql.DB) error) error {
	// 1. 個別に開く
	dsn := fmt.Sprintf("file:%s", path)
	storeDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("open %s: %w", schemaName, err)
	}
	defer storeDB.Close()

	// 2. マイグレーション実行
	if err := migrateFunc(storeDB); err != nil {
		return fmt.Errorf("migrate %s: %w", schemaName, err)
	}

	// 3. 閉じる（WALのコミット等）
	if err := storeDB.Close(); err != nil {
		return fmt.Errorf("close %s after migration: %w", schemaName, err)
	}

	// 4. ATTACH
	attachSQL := fmt.Sprintf("ATTACH DATABASE '%s' AS %s;", path, schemaName)
	if _, err := mainDB.Exec(attachSQL); err != nil {
		return fmt.Errorf("attach %s: %w", schemaName, err)
	}

	return nil
}