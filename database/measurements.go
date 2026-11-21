package database

import (
	"context"
	"fmt"
	"time"
)

type MeasurementType string

const (
	WaterLevel  MeasurementType = "water_level"
	Temperature MeasurementType = "temperature"

	Snapshot MeasurementType = "snapshot"
)

type Measurement struct {
	Time  time.Time
	Type  MeasurementType
	Value float64
}

func (m *Measurement) Formatted() string {
	switch m.Type {
	case WaterLevel:
		return fmt.Sprintf("%v inches", m.Value)
	case Temperature:
		return fmt.Sprintf("%v °C", m.Value)
	case Snapshot:
		return fmt.Sprintf("snapshots/%v.jpg", int(m.Value))
	default:
		panic(fmt.Errorf("unsupported type: %v", m.Type))
	}
}

func (db *DB) Record(ctx context.Context, m *Measurement) error {
	_, err := db.db.ExecContext(ctx,
		`INSERT INTO measurements (time, type, value)
		VALUES (?, ?, ?)`,
		m.Time, string(m.Type), m.Value)
	if err != nil {
		return fmt.Errorf("db: create measurement: %w", err)
	}
	return nil
}

func (db *DB) GetMeasurements(ctx context.Context, mt MeasurementType, timeRange [2]time.Time) ([]*Measurement, error) {
	rows, err := db.db.QueryContext(ctx,
		`SELECT time, type, value FROM measurements
		WHERE type = ? AND time >= ? and time < ? ORDER BY TIME ASC`,
		mt, timeRange[0], timeRange[1])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Measurement
	for rows.Next() {
		var m Measurement
		if err := rows.Scan(&m.Time, &m.Type, &m.Value); err != nil {
			return nil, fmt.Errorf("db: query measurements: %w", err)
		}
		results = append(results, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: query measurements: %w", err)
	}
	return results, nil
}
