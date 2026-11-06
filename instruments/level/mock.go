package level

type mock struct{}

func (m mock) WaterLevel() (float64, error) {
	return 8.5, nil
}
