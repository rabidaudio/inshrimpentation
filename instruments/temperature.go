package instruments

import (
	"fmt"

	"periph.io/x/conn/v3/onewire/onewirereg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ds18b20"
	"periph.io/x/host/v3"
)

// Temperature connects to the DS18B20 via 1-wire
// protocol and returns a temperature in Celsius
// with a resolution of 12-bits (0.0625C steps).
func Temperature() (float64, error) {
	_, err := host.Init()
	if err != nil {
		return 0, err
	}

	b, err := onewirereg.Open("") // open default 1-wire interface
	if err != nil {
		return 0, err
	}
	defer b.Close()

	addrs, err := b.Search(false) // search for all devices
	if err != nil {
		return 0, err
	}
	// TODO: assumes only one device on bus
	if len(addrs) < 1 {
		return 0, fmt.Errorf("instruments: temperature sensor not found")
	}
	addr := addrs[0]

	d, err := ds18b20.New(b, addr, 12) // connect using 12-bit resolution
	if err != nil {
		return 0, err
	}

	var temp physic.Env
	err = d.Sense(&temp) // blocking read temp
	if err != nil {
		return 0, err
	}
	return temp.Temperature.Celsius(), nil
}
