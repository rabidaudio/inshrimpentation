package instruments

import (
	"periph.io/x/devices/v3/ads1x15"
)

// WaterLevel in inches
func WaterLevel() (float64, error) {
	raw, err := readAdc(ads1x15.Channel0, 16)
	if err != nil {
		return 0, err
	}
	// TODO: convert resistance to inches
	return float64(raw), nil
}
