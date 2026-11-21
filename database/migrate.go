package database

import (
	"database/sql"
	"fmt"
	"log"
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *DB) error {
	if db == nil {
		var err error
		db, err = Init()
		if err != nil {
			return err
		}
	}
	driver, err := sqlite.WithInstance(db.db.(*sql.DB), &sqlite.Config{
		MigrationsTable: "_migrations",
	})
	if err != nil {
		return err
	}
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	path := path.Join(root, "database/migrations")
	path = fmt.Sprintf("file://%v", path)
	m, err := migrate.NewWithDatabaseInstance(path, "sqlite", driver)
	if err != nil {
		return err
	}
	m.Log = migrateLogger{}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Printf("migrations: up-to-date")
		return nil
	}
	return err
}

type migrateLogger struct{}

func (_ migrateLogger) Printf(format string, v ...interface{}) {
	log.Printf("migrate:"+format, v...)
}

func (_ migrateLogger) Verbose() bool {
	return false
}
