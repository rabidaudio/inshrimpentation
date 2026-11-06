package db

import (
	"context"
	"fmt"
	"time"
)

const createMeasurementSQL = `
create table measurements (
	id integer not null primary key,
	time timestamp not null,
	type varchar not null,
	value numeric not null
);`

type MeasurementType string

const (
	WaterLevel  MeasurementType = "water_level"
	Temperature MeasurementType = "temperature"
)

type Measurement struct {
	time            time.Time
	measurementType MeasurementType
	value           float64
}

func (db *DB) Record(ctx context.Context, m Measurement) error {
	_, err := db.ExecContext(ctx, "INSERT INTO measurements (time, type, value) VALUES (?, ?, ?)", m.time, string(m.measurementType), m.value)
	if err != nil {
		return fmt.Errorf("db: create measurement: %w", err)
	}
	return nil
}

func (db *DB) GetMeasurements(ctx context.Context, mt MeasurementType, timeRange [2]time.Time) ([]Measurement, error) {
	rows, err := db.QueryContext(ctx, "SELECT time, type, value FROM measurements WHERE type = ? AND time >= ? and time < ? ORDER BY TIME ASC", mt, timeRange[0], timeRange[1])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Measurement
	for rows.Next() {
		var m Measurement
		if err := rows.Scan(&m.time, &m.measurementType, &m.value); err != nil {
			return nil, fmt.Errorf("db: query measurements: %w", err)
		}
		results = append(results, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: query measurements: %w", err)
	}
	return results, nil
}
