package database

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"time"
)

var optionsTypes = map[string]reflect.Type{
	string(WaterLevel): reflect.TypeOf(WaterLevelOptions{}),
}

var defaultTime = time.UnixMicro(0)
var defaultInterval = time.Duration(math.MaxInt64)

type Config struct {
	Type          string
	Interval      time.Duration
	LastCompleted time.Time
	Enabled       bool
	Options       any
}

type WaterLevelOptions struct {
	LowLevel float32 `json:"low_level"`
}

func (c *Config) DueAt() time.Time {
	return c.LastCompleted.Add(c.Interval)
}

func (c *Config) DueIn() time.Duration {
	return time.Until(c.DueAt())
}

func (c *Config) IsDue() bool {
	return c.DueIn() <= 0
}

func (db *DB) GetConfigs(ctx context.Context) ([]*Config, error) {
	rows, err := db.db.QueryContext(ctx,
		`SELECT type, interval, last_completed, enabled, options FROM configs ORDER BY last_completed ASC, type ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Config
	for rows.Next() {
		var c Config
		var optionsJson string
		if err := rows.Scan(&c.Type, &c.Interval, &c.LastCompleted, &c.Enabled, &optionsJson); err != nil {
			return nil, fmt.Errorf("db: query configs: %w", err)
		}
		if ot, ok := optionsTypes[string(c.Type)]; ok {
			v := reflect.New(ot)
			if err = json.Unmarshal([]byte(optionsJson), v.Interface()); err != nil {
				return nil, fmt.Errorf("db: query configs: parse json [%v] %w", optionsJson, err)
			}
			c.Options = reflect.Indirect(v).Interface()
		}
		results = append(results, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: query configs: %w", err)
	}
	return results, nil
}

func (db *DB) SetInterval(ctx context.Context, t string, i time.Duration) error {
	_, err := db.db.ExecContext(ctx,
		`INSERT INTO configs (type, interval, last_completed) VALUES (?, ?, ?)
		ON CONFLICT(type) DO UPDATE SET interval = ?`,
		t, int64(i), defaultTime, int64(i))
	if err != nil {
		return fmt.Errorf("db: upsert configs: %w", err)
	}
	return nil
}

func (db *DB) SetConfigEnabled(ctx context.Context, t string, enabled bool) error {
	_, err := db.db.ExecContext(ctx,
		`INSERT INTO configs (type, enabled, interval, last_completed) VALUES (?, ?, ?, ?)
		ON CONFLICT(type) DO UPDATE SET enabled = ?`,
		t, enabled, defaultInterval, defaultTime, enabled)
	if err != nil {
		return fmt.Errorf("db: upsert configs: %w", err)
	}
	return nil
}

func (db *DB) SetConfigOptions(ctx context.Context, t string, opts any) error {
	ot := reflect.TypeOf(opts)
	if ot == reflect.MapOf(reflect.TypeFor[string](), reflect.TypeFor[any]()) {
		// trust the format and write it directly
	} else if ot != optionsTypes[t] {
		return fmt.Errorf("db: upsert configs: invalid option type for %v: %v", t, opts)
	}
	optsJson, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("db: upsert configs: marshal opts: %v", err)
	}

	_, err = db.db.ExecContext(ctx,
		`INSERT INTO configs (type, options, interval, last_completed) VALUES (?, ?, ?, ?)
		ON CONFLICT(type) DO UPDATE SET options = ?`,
		t, string(optsJson), defaultInterval, defaultTime, string(optsJson))
	if err != nil {
		return fmt.Errorf("db: upsert configs: %w", err)
	}
	return nil
}

func (db *DB) TrackRun(ctx context.Context, t string, runat time.Time) error {
	_, err := db.db.ExecContext(ctx,
		`UPDATE configs SET last_completed = ? WHERE type = ?`, runat, t)
	if err != nil {
		return fmt.Errorf("db: update configs: %w", err)
	}
	return nil
}
