package instruments

import (
	"context"
	"time"

	"github.com/rabidaudio/inshrimpentation/database"
	"golang.org/x/sync/errgroup"
)

type Instrument func() (float64, error)

var instruments = map[database.MeasurementType]Instrument{
	database.WaterLevel:  WaterLevel,
	database.Temperature: Temperature,
}

// Measure returns a collection of one-off measurements, without writing
// anything to the database.
func Measure(ctx context.Context) ([]*database.Measurement, error) {
	vv := make(chan *database.Measurement)
	eg, _ := errgroup.WithContext(ctx) // TODO: Instruments should take context
	for t, i := range instruments {
		eg.Go(func() error {
			v, err := i()
			now := time.Now()
			if err != nil {
				return err
			}
			vv <- &database.Measurement{
				Type:  t,
				Value: v,
				Time:  now,
			}
			return nil
		})
	}
	go func() {
		eg.Wait()
		close(vv)
	}()

	var mm []*database.Measurement
	for m := range vv {
		mm = append(mm, m)
	}
	return mm, nil
}
