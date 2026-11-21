package instruments

import (
	"fmt"

	"github.com/rabidaudio/api/env"
	"periph.io/x/conn/v3/onewire"
	"periph.io/x/conn/v3/onewire/onewirereg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ds18b20"
	"periph.io/x/host/v3"
)

// Temperature connects to the DS18B20 via 1-wire
// protocol and returns a temperature in Celsius
// with a resolution of 12-bits (0.0625C steps).
func Temperature() (float64, error) {
	if !env.IsProd() {
		return 20.1, nil // mock data
	}

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
	addr, ok := findTempSensor(addrs)
	if !ok {
		return 0, fmt.Errorf("instruments: temperature sensor not found")
	}

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

func findTempSensor(addrs []onewire.Address) (onewire.Address, bool) {
	for _, addr := range addrs {
		fam := ds18b20.Family(addr)
		if fam == ds18b20.DS18B20 {
			return addr, true
		}
	}
	return 0, false
}
