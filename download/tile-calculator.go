package download

import (
	"math"
)

type TileCalculator struct {
	Converter *Converter
}

func (tc *TileCalculator) CalculateTiles(zoom int, north float64, west float64, south float64, east float64) int {
	top := tc.Converter.Lat2tile(north, zoom)
	left := tc.Converter.Lon2tile(west, zoom)
	bottom := tc.Converter.Lat2tile(south, zoom)
	right := tc.Converter.Lon2tile(east, zoom)
	width := math.Abs(float64(left-right)) + 1
	height := math.Abs(float64(top-bottom)) + 1

	return int(width * height)
}

func NewTileCalculator(converter *Converter) *TileCalculator {
	return &TileCalculator{Converter: converter}
}
