package flags

import "flag"

type Flags struct {
	BaseMapUrl string
	OverlayUrl string
	MinZoom    int
	MaxZoom    int
	North      float64
	West       float64
	South      float64
	East       float64
	Quality    int
	Country    string
	AtlasName  string
}

func ParseFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.BaseMapUrl, "baseMapUrl", "", "map URL")
	flag.StringVar(&f.OverlayUrl, "overlayUrl", "", "map URL")
	flag.StringVar(&f.Country, "country", "", "country code, supported: CZ, SK, overrides manual coordinates")
	flag.StringVar(&f.AtlasName, "atlasName", "", "name of the generated atlas")
	flag.IntVar(&f.MinZoom, "minZoom", 8, "min zoom")
	flag.IntVar(&f.MaxZoom, "maxZoom", 12, "max zoom")
	flag.Float64Var(&f.North, "n", 0, "north bounding point")
	flag.Float64Var(&f.West, "w", 0, "west bounding point")
	flag.Float64Var(&f.South, "s", 0, "south bounding point")
	flag.Float64Var(&f.East, "e", 0, "east bounding point")
	flag.IntVar(&f.Quality, "quality", 100, "quality")
	flag.Parse()

	return f
}
