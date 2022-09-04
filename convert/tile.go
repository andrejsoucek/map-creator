package convert

import (
	"math"
)

func Lon2tile(lon float64, zoom int) int {
	n := math.Exp2(float64(zoom))
	return int(math.Floor((lon + 180.0) / 360.0 * n))
}

func Lat2tile(lat float64, zoom int) int {
	n := math.Exp2(float64(zoom))
	return int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n))
}

func Deg2num(lat float64, lon float64, zoom int) (x int, y int) {
	n := math.Exp2(float64(zoom))
	x = Lon2tile(lon, zoom)
	if float64(x) >= n {
		x = int(n - 1)
	}
	y = Lat2tile(lat, zoom)
	return
}
