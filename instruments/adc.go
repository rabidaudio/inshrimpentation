package instruments

import (
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

const vref = 3300 * physic.MilliVolt

func readAdc(channel ads1x15.Channel, avgsamples int) (int32, error) {
	_, err := host.Init()
	if err != nil {
		return 0, err
	}

	bus, err := i2creg.Open("") // open default i2c device (dev 1)
	if err != nil {
		return 0, err
	}
	defer bus.Close()

	adc, err := ads1x15.NewADS1015(bus, &ads1x15.DefaultOpts)
	if err != nil {
		return 0, err
	}

	pin, err := adc.PinForChannel(channel, vref, 10*physic.KiloHertz, ads1x15.BestQuality)
	if err != nil {
		return 0, err
	}
	defer pin.Halt()

	samples := make([]int32, avgsamples)
	for i := range avgsamples {
		reading, err := pin.Read()
		if err != nil {
			return 0, err
		}
		samples[i] = reading.Raw
	}
	var sum int64 = 0
	for i := range avgsamples {
		sum += int64(i)
	}

	return int32(sum / int64(avgsamples)), nil
}
