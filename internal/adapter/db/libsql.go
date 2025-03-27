package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log"

	"github.com/pimg/certguard/internal/adapter/db/queries"
	"github.com/rubenv/sql-migrate"
	_ "github.com/tursodatabase/go-libsql"
)

//go:embed schema/*.sql
var dbMigrations embed.FS

type LibSqlStorage struct {
	DB      *sql.DB
	Queries *queries.Queries
}

func NewLibSqlStorage(db *sql.DB) *LibSqlStorage {
	return &LibSqlStorage{
		DB:      db,
		Queries: queries.New(db),
	}
}

func NewDBConnection(dbLocation string) (*sql.DB, error) {
	log.Println("Connecting to DB...")
	db, err := sql.Open("libsql", "file:"+dbLocation+"/certguard.db?_journal_mode=WAL&busy_timeout=5000_foreign_keys=on")
	if err != nil {
		log.Printf("Error connecting to DB: %v", err)
		return nil, err
	}

	return db, nil
}

func (s *LibSqlStorage) InitDB(ctx context.Context) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "schema",
	}

	n, err := migrate.ExecContext(ctx, s.DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return errors.Join(errors.New("error performing migrations"), err)
	}
	log.Printf("database initialized, applied %d migrations!", n)

	return nil
}

func (s *LibSqlStorage) CloseDB() error {
	if closeError := s.DB.Close(); closeError != nil {
		return errors.Join(errors.New("error closing database"), closeError)
	}
	log.Println("Database closed")

	return nil
}
