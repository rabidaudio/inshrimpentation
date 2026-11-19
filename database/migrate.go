package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/rabidaudio/api/env"

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
	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{
		MigrationsTable: "_migrations",
	})
	if err != nil {
		return err
	}
	path := "database/migrations"
	if env.IsTest() {
		path = "../" + path
	}
	path = fmt.Sprintf("file://%v", path)
	m, err := migrate.NewWithDatabaseInstance(path, "sqlite", driver)
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
