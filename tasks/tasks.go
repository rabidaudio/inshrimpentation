package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/instruments"
)

func measureWaterLevel(ctx context.Context, opts any) ([]*database.Measurement, *database.Alert, error) {
	t := time.Now()
	v, err := instruments.WaterLevel()
	if err != nil {
		return nil, nil, err
	}
	m := database.Measurement{
		Time:  t,
		Type:  database.WaterLevel,
		Value: v,
	}
	var a *database.Alert
	limit := opts.(database.WaterLevelOptions).LowLevel
	if limit != 0 && v < float64(limit) {
		a = &database.Alert{
			Type:       string(MeasureWaterLevel),
			Time:       t,
			Suppressed: false,
			Message:    fmt.Sprintf("Water level %v below limit %v", v, limit),
		}
	}
	return []*database.Measurement{&m}, a, nil
}

func measureTemperature(ctx context.Context, opts any) ([]*database.Measurement, *database.Alert, error) {
	t := time.Now()
	v, err := instruments.Temperature()
	if err != nil {
		return nil, nil, err
	}
	m := database.Measurement{
		Time:  t,
		Type:  database.Temperature,
		Value: v,
	}
	return []*database.Measurement{&m}, nil, nil
}

func needsCleaning(ctx context.Context, opts any) ([]*database.Measurement, *database.Alert, error) {
	return nil, &database.Alert{
		Type:       string(NeedsCleaning),
		Time:       time.Now(),
		Suppressed: false,
		Message:    "It's time to clean the tank",
	}, nil
}
