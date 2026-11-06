package level

type LevelMeter interface {
	// return the current water level in inches
	WaterLevel() (float64, error)
}

func New() (LevelMeter, error) {
	// TODO: return real hardware backed version instead
	return mock{}, nil
}
