package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlerts(t *testing.T) {
	db, err := Init()
	assert.NoError(t, err)

	// create an alert
	a1 := Alert{
		Type:       "foo",
		Time:       time.Unix(1795105087, 0),
		Suppressed: false,
		Message:    "hello world",
	}
	err = db.AddAlert(t.Context(), &a1)
	assert.NoError(t, err, "should insert")

	assert.NotEqual(t, 0, a1.id, "should have set an id")

	// query for unsuppressed alerts
	res1, err := db.GetAlerts(t.Context(), false)
	assert.NoError(t, err, "should query")
	assert.Equal(t, 1, len(res1), "should return 1 result")

	a2 := res1[0]
	assert.EqualValues(t, a1, *a2, "should have the created alert")

	// suppress
	err = db.SupressAlert(t.Context(), a2)
	assert.NoError(t, err, "should allow suppressing alerts")

	// should no longer return
	res2, err := db.GetAlerts(t.Context(), false)
	assert.NoError(t, err, "should query")
	assert.Equal(t, 0, len(res2), "should return no results")

	// ...unless all alerts are included
	res3, err := db.GetAlerts(t.Context(), true)
	assert.NoError(t, err, "should query")
	assert.Equal(t, 1, len(res3), "should return 1 result")
	a3 := res3[0]
	assert.Equal(t, a1.id, a3.id, "should be the same alert")
	assert.Equal(t, true, a3.Suppressed, "should be suppressed")
}
