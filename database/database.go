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

type dbLike interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB struct {
	db   dbLike
	isTx bool
}

func (db *DB) InTx(ctx context.Context, opts *sql.TxOptions, fn func(db *DB) error) error {
	if db.isTx {
		return fn(db)
	}
	tx, err := db.db.(*sql.DB).BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	err = fn(&DB{db: tx, isTx: true})
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) Close() error {
	if db.isTx {
		return fmt.Errorf("db: cannot close: inside transaction")
	}
	return db.db.(*sql.DB).Close()
}

func Init() (*DB, error) {
	log.Printf("open db: %v\n", dbpath())
	if env.IsProd() {
		// make sure directory exists
		if err := os.MkdirAll(path.Dir(dbpath()), 0664); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbpath()))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	dbi := &DB{db: db, isTx: false}
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
		return fmt.Sprintf("%v/storage", home)
	}
	root, err := findProjectRoot()
	if err != nil {
		panic(err)
	}
	return path.Join(root, "storage")
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		_, err := os.Stat(path.Join(dir, "go.mod"))
		if err == nil {
			return dir, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		if dir == "/" {
			return "", fmt.Errorf("unable to find project root")
		}
		dir = path.Dir(dir)
	}
}

func dbpath() string {
	if env.IsTest() {
		return ":memory:"
	}
	return fmt.Sprintf("%v/database_%v.db", StorageRoot(), env.Env())
}
