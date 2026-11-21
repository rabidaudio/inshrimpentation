package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigs(t *testing.T) {
	// assert.NoError(t, Create())

	db, err := Init()
	assert.NoError(t, err)

	err = db.SetInterval(t.Context(), string(WaterLevel), 1*time.Hour)
	assert.NoError(t, err)

	res1, err := db.GetConfigs(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res1))
	c1 := res1[0]
	assert.Equal(t, string(WaterLevel), c1.Type)
	assert.Equal(t, 1*time.Hour, c1.Interval)
	assert.Equal(t, time.UnixMilli(0), c1.LastCompleted)
	assert.Equal(t, WaterLevelOptions{}, c1.Options, "should the zero value")
	assert.True(t, c1.Enabled)

	err = db.SetConfigOptions(t.Context(), string(WaterLevel), WaterLevelOptions{
		LowLevel: 2.5,
	})
	assert.NoError(t, err)

	res2, err := db.GetConfigs(t.Context())
	assert.Equal(t, 1, len(res2))
	c2 := res2[0]
	assert.Equal(t, string(WaterLevel), c2.Type)
	assert.Equal(t, float32(2.5), c2.Options.(WaterLevelOptions).LowLevel)

	tt := time.Unix(1795107960, 0)
	err = db.TrackRun(t.Context(), string(WaterLevel), tt)
	assert.NoError(t, err)

	res3, err := db.GetConfigs(t.Context())
	assert.Equal(t, 1, len(res3))
	c3 := res3[0]
	assert.Equal(t, string(WaterLevel), c3.Type)
	assert.Equal(t, tt, c3.LastCompleted)

	err = db.SetConfigEnabled(t.Context(), string(WaterLevel), false)
	assert.NoError(t, err)

	res4, err := db.GetConfigs(t.Context())
	assert.Equal(t, 1, len(res4))
	c4 := res4[0]
	assert.Equal(t, string(WaterLevel), c4.Type)
	assert.False(t, c4.Enabled)
}
