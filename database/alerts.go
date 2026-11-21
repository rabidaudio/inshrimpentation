package database

import (
	"context"
	"fmt"
	"time"
)

type Alert struct {
	id         int64
	Type       string
	Time       time.Time
	Suppressed bool
	Message    string
}

func (db *DB) AddAlert(ctx context.Context, a *Alert) error {
	return db.InTx(ctx, nil, func(db *DB) error {
		existing, err := db.db.QueryContext(ctx,
			`SELECT 1 FROM alerts WHERE type = ? and suppressed = ? LIMIT 1`, a.Type, false)
		if err != nil {
			return fmt.Errorf("db: insert alert: %w", err)
		}
		defer existing.Close()

		if existing.Next() {
			return nil // already exists, don't make a new one
		}

		res, err := db.db.ExecContext(ctx,
			`INSERT INTO alerts (type, time, suppressed, message) VALUES (?, ?, ?, ?) RETURNING id`,
			a.Type, a.Time, a.Suppressed, a.Message)
		if err != nil {
			return fmt.Errorf("db: insert alert: %w", err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("db: insert alert: %w", err)
		}
		a.id = id
		return nil
	})
}

func (db *DB) SupressAlert(ctx context.Context, a *Alert) error {
	_, err := db.db.ExecContext(ctx,
		`UPDATE alerts SET suppressed = TRUE WHERE id = ?`, a.id)
	if err != nil {
		return fmt.Errorf("db: update alert: %w", err)
	}
	return nil
}

func (db *DB) GetAlerts(ctx context.Context, all bool) ([]*Alert, error) {
	query := `SELECT id, type, time, suppressed, message FROM alerts`
	if !all {
		query += ` WHERE suppressed = FALSE`
	}
	query += ` ORDER BY time DESC`
	rows, err := db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Alert
	for rows.Next() {
		var a Alert
		if err := rows.Scan(&a.id, &a.Type, &a.Time, &a.Suppressed, &a.Message); err != nil {
			return nil, fmt.Errorf("db: query alerts: %w", err)
		}
		results = append(results, &a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: query alerts: %w", err)
	}
	return results, nil
}
