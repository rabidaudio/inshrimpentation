package instruments

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabidaudio/inshrimpentation/database"
)

type Instrument func() (float64, error)

var instruments = map[database.MeasurementType]Instrument{
	database.WaterLevel:  WaterLevel,
	database.Temperature: Temperature,
	// database.Brightness
	// database.Algae
	database.Snapshot: Snapshot,
}

func Measure(ctx context.Context) ([]*database.Measurement, error) {
	t := time.Now()

	db := database.GetDB(ctx)
	conf, err := db.GetConfigs(ctx)
	if err != nil {
		return nil, err
	}

	var measurements []*database.Measurement

	for _, c := range conf {
		if time.Since(c.LastCompleted) < c.Interval {
			continue // not time to run yet
		}
		m, err := runTest(ctx, db, c, t)
		if err != nil {
			// log warning and continue
			// TODO: alarm
			log.Printf("instruments: error running %v: %v", c.MeasurementType, err)
		}
		measurements = append(measurements, m)
	}

	// TODO: trigger alarms
	// TODO: alarm for time since filter/water clean

	return measurements, nil
}

func runTest(ctx context.Context, db *database.DB, c database.Config, t time.Time) (m *database.Measurement, err error) {
	i, ok := instruments[c.MeasurementType]
	if !ok {
		return m, fmt.Errorf("todo handle %v", c.MeasurementType)
	}
	m.Time = time.Now()
	m.MeasurementType = c.MeasurementType
	v, err := i()
	if err != nil {
		return
	}

	m.Value = v
	err = db.Record(ctx, m)
	return
}
