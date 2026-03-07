package service

import (
	"context"
	"database/sql"
)

func queryRowsWithFallback(ctx context.Context, conn *sql.DB, ftsSQL string, ftsArgs []interface{}, likeSQL string, likeArgs []interface{}) (*sql.Rows, error) {
	rows, err := conn.QueryContext(ctx, ftsSQL, ftsArgs...)
	if err == nil {
		return rows, nil
	}
	return conn.QueryContext(ctx, likeSQL, likeArgs...)
}

func likePattern(query string) string {
	return "%" + query + "%"
}
