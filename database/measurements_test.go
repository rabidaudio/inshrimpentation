package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMeasurements(t *testing.T) {
	assert.NoError(t, Create())

	db, err := Init()
	assert.NoError(t, err)
	defer db.Wipe()

	m := Measurement{
		time:            time.Date(2026, 11, 1, 15, 00, 00, 00, time.UTC),
		measurementType: WaterLevel,
		value:           8.5,
	}

	err = db.Record(t.Context(), m)
	assert.NoError(t, err)

	// read the measurements in that year
	mm, err := db.GetMeasurements(t.Context(), WaterLevel, [2]time.Time{
		time.Date(2026, 1, 1, 00, 00, 00, 00, time.UTC),
		time.Date(2027, 1, 1, 00, 00, 00, 00, time.UTC),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mm))
	assert.EqualValues(t, m, mm[0])

	// read measurements from the previous year and expect none
	mm2, err := db.GetMeasurements(t.Context(), WaterLevel, [2]time.Time{
		time.Date(2025, 1, 1, 00, 00, 00, 00, time.UTC),
		time.Date(2026, 1, 1, 00, 00, 00, 00, time.UTC),
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(mm2))
}
