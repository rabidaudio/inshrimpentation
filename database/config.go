package database

import (
	"context"
	"fmt"
	"time"
)

const createConfigSQL = `
create table configs (
	id integer not null primary key,
	type string not null,
	interval integer not null,
	last_completed timestamp
);
CREATE UNIQUE INDEX idx_configs_type ON configs(type);`

type Config struct {
	MeasurementType MeasurementType
	Interval        time.Duration
	LastCompleted   time.Time
}

func (db *DB) GetConfigs(ctx context.Context) ([]Config, error) {
	rows, err := db.QueryContext(ctx, "SELECT type, interval, last_completed FROM configs ORDER BY last_completed ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Config
	for rows.Next() {
		var c Config
		if err := rows.Scan(&c.MeasurementType, &c.Interval, &c.LastCompleted); err != nil {
			return nil, fmt.Errorf("db: query configs: %w", err)
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: query configs: %w", err)
	}
	return results, nil
}

func (db *DB) SetInterval(ctx context.Context, mt MeasurementType, i time.Duration) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO configs (type, interval) VALUES (?, ?)
		ON CONFLICT(type) DO UPDATE SET interval = ?`, string(mt), int64(i), int64(i))
	if err != nil {
		return fmt.Errorf("db: upsert configs: %w", err)
	}
	return nil
}

func (db *DB) TrackRun(ctx context.Context, mt MeasurementType, runat time.Time) error {
	_, err := db.ExecContext(ctx, `
		UPDATE configs SET last_completed = ?
		WHERE type = ?`, runat, string(mt))
	if err != nil {
		return fmt.Errorf("db: update configs: %w", err)
	}
	return nil
}
