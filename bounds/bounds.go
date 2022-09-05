package bounds

type CountryBounds struct {
	North float64
	West  float64
	South float64
	East  float64
}

func CzechRepublic() CountryBounds {
	return CountryBounds{
		North: 51.0835,
		West:  12.0475,
		South: 48.55,
		East:  18.9,
	}
}

func Slovakia() CountryBounds {
	return CountryBounds{
		North: 49.624,
		West:  16.777,
		South: 47.7,
		East:  22.62,
	}
}
