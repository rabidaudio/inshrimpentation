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
	database.WaterLevel: WaterLevel,
	// database.Temperature
	// database.Brightness
	// database.Algae
	database.Snapshot: Snapshot,
}

func Measure(ctx context.Context) error {
	t := time.Now()

	db := database.GetDB(ctx)
	conf, err := db.GetConfigs(ctx)
	if err != nil {
		return err
	}

	for _, c := range conf {
		if time.Since(c.LastCompleted) < c.Interval {
			continue // not time to run yet
		}
		err = runTest(ctx, db, c, t)
		if err != nil {
			// log warning and continue
			// TODO: alarm
			log.Printf("instruments: error running %v: %v", c.MeasurementType, err)
		}
	}

	// TODO: trigger alarms
	// TODO: alarm for time since filter/water clean

	return nil
}

func runTest(ctx context.Context, db *database.DB, c database.Config, t time.Time) error {
	i, ok := instruments[c.MeasurementType]
	if !ok {
		return fmt.Errorf("todo handle %v", c.MeasurementType)
	}
	v, err := i()
	if err != nil {
		return err
	}
	return db.Record(ctx, database.Measurement{
		Time:            time.Now(),
		Value:           v,
		MeasurementType: c.MeasurementType,
	})
}
