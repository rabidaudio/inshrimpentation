package tasks

import (
	"context"
	"maps"
	"slices"

	"github.com/rabidaudio/inshrimpentation/database"
)

type Task string

const (
	MeasureWaterLevel  Task = Task(database.WaterLevel)
	MeasureTemperature Task = Task(database.Temperature)
	// TakePhoto
	// MeasureBrightness
	// MeasureAlgae
	//

	NeedsCleaning Task = "needs_cleaning"
)

type taskRunner func(ctx context.Context, opts any) ([]*database.Measurement, *database.Alert, error)

var taskRunners = map[Task]taskRunner{
	MeasureWaterLevel:  measureWaterLevel,
	MeasureTemperature: measureTemperature,
	NeedsCleaning:      needsCleaning,
}

func AllTasks() []Task {
	return slices.Collect(maps.Keys(taskRunners))
}
