package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func InitDB(dbPath string) (*sql.DB, error) {
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_TOKEN")

	var dsn string
	if tursoURL != "" && tursoToken != "" {
		dsn = fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)
	} else if dbPath != "" {
		// Local dev fallback: use a local file via libsql
		dsn = fmt.Sprintf("file:%s", dbPath)
	} else {
		dsn = "file:./bulletin.db"
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(schema)); err != nil {
		return nil, err
	}
	return db, nil
}
