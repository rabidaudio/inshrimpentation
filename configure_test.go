package main

import (
	"testing"
	"time"

	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/tasks"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	db, err := database.Init()
	assert.NoError(t, err)

	ctx := database.WithDB(t.Context(), db)

	err = configure(ctx, "default_config.yml")
	assert.NoError(t, err)

	confs, err := db.GetConfigs(ctx)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(confs))

	assert.Equal(t, string(tasks.NeedsCleaning), confs[0].Type)
	assert.Equal(t, 4*7*24*time.Hour, confs[0].Interval)
	assert.True(t, confs[0].Enabled)
	assert.True(t, confs[0].IsDue())

	assert.Equal(t, string(tasks.MeasureTemperature), confs[1].Type)
	assert.Equal(t, 15*time.Minute, confs[1].Interval)
	assert.False(t, confs[1].Enabled)
	assert.True(t, confs[1].IsDue())

	assert.Equal(t, string(tasks.MeasureWaterLevel), confs[2].Type)
	assert.Equal(t, 1*time.Hour, confs[2].Interval)
	assert.Equal(t, float32(-1.5), confs[2].Options.(database.WaterLevelOptions).LowLevel)
	assert.True(t, confs[2].Enabled)
	assert.True(t, confs[2].IsDue())
}

func TestCheck(t *testing.T) {
	db, err := database.Init()
	assert.NoError(t, err)

	ctx := database.WithDB(t.Context(), db)

	err = configure(ctx, "default_config.yml")
	assert.NoError(t, err)

	err = check(ctx)
	assert.NoError(t, err)
}
