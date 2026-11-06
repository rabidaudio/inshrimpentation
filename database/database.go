// Package db exposes a sqlite3 database for recording
// measurements
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rabidaudio/inshrimpentation/api"
)

type DB struct {
	*sql.DB
}

func path() string {
	env := api.Env()
	dir := ""
	if env == api.Test {
		dir = "../"
	}
	return fmt.Sprintf("%vstorage/database_%v.db", dir, env)
}

func Init() (*DB, error) {
	log.Printf("open db: %v\n", path())
	db, err := sql.Open("sqlite3", path())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func Create() error {
	_, err := os.Stat(path())
	if err == nil {
		// file already exists
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// create
	db, err := Init()
	if err != nil {
		return err
	}
	defer db.Close()

	// create tables
	_, err = db.Exec(createMeasurementSQL)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Wipe() {
	if !api.IsEnv(api.Test) {
		panic("Don't wipe real database")
	}
	err := os.Remove(path())
	if err != nil {
		panic(err)
	}
}

func WithDB(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, "DB", db)
}

func GetDB(ctx context.Context) *DB {
	db, ok := ctx.Value("DB").(*DB)
	if !ok {
		panic(fmt.Errorf("No DB found in context"))
	}
	return db
}
