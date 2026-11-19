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
	"path"

	"github.com/rabidaudio/api/env"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Init() (*DB, error) {
	log.Printf("open db: %v\n", dbpath())
	if env.IsProd() {
		// make sure directory exists
		if err := os.MkdirAll(path.Dir(dbpath()), 0664); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", dbpath())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	dbi := &DB{db}
	err = RunMigrations(dbi)
	if err != nil {
		return nil, fmt.Errorf("migrate error: %v", err)
	}
	return dbi, nil
}

func Create() error {
	p := dbpath()
	_, err := os.Stat(p)
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

	return nil
}

func (db *DB) Wipe() {
	if !env.IsTest() {
		panic("Don't wipe real database")
	}
	err := os.Remove(dbpath())
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

func StorageRoot() string {
	if env.IsProd() {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		// overwrite to permanent filesystem
		return fmt.Sprintf("%v/storage", home)
	} else if env.IsTest() {
		return "../storage"
	}
	return "storage"
}

func dbpath() string {
	return fmt.Sprintf("%v/database_%v.db", StorageRoot(), env.Env())
}
