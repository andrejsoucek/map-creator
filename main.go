package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/andrejsoucek/map-creator/bounds"
	"github.com/andrejsoucek/map-creator/tiles"
	"github.com/andrejsoucek/map-creator/ui"
	progressbar "github.com/schollz/progressbar/v3"
)

var baseMapUrl string
var overlayUrl string
var country string
var atlasName string
var minZoom int
var maxZoom int
var north float64
var west float64
var south float64
var east float64
var maxConcurrentDownloads int
var quality int

func parseFlags() {
	flag.StringVar(&baseMapUrl, "baseMapUrl", "", "map URL")
	flag.StringVar(&overlayUrl, "overlayUrl", "", "map URL")
	flag.StringVar(&country, "country", "CZ", "country code, supported: CZ, SK, overrides manual coordinates")
	flag.StringVar(&atlasName, "atlasName", "", "name of the generated atlas")
	flag.IntVar(&minZoom, "minZoom", 8, "min zoom")
	flag.IntVar(&maxZoom, "maxZoom", 12, "max zoom")
	flag.Float64Var(&north, "n", 0, "north bounding point")
	flag.Float64Var(&west, "w", 0, "west bounding point")
	flag.Float64Var(&south, "s", 0, "south bounding point")
	flag.Float64Var(&east, "e", 0, "east bounding point")
	flag.IntVar(&maxConcurrentDownloads, "maxConcurrentDownloads", 4, "max concurrent downloads")
	flag.IntVar(&quality, "quality", 100, "quality")
	flag.Parse()
}

func createDownloadParams() tiles.DownloadParams {
	var c bounds.CountryBounds
	if country != "" {
		switch country {
		case "CZ":
			c = bounds.CzechRepublic()
		case "SK":
			c = bounds.Slovakia()
		default:
			log.Fatal("Unknown country, use n, w, s, e flags instead.")
		}
		north = c.North
		west = c.West
		south = c.South
		east = c.East
	}
	return tiles.DownloadParams{
		BaseMapUrl:       baseMapUrl,
		OverlayUrl:       overlayUrl,
		MinZoom:          minZoom,
		MaxZoom:          maxZoom,
		North:            north,
		West:             west,
		South:            south,
		East:             east,
		Quality:          quality,
	}
}

func main() {
	parseFlags()

	totalTiles := 0
	totalFilesToDownload := 0

	p := createDownloadParams()

	for i := minZoom; i <= maxZoom; i++ {
		count := tiles.CalculateTiles(i, north, west, south, east)
		totalTiles = totalTiles + count
		if overlayUrl != "" {
			count = count * 2
		}
		totalFilesToDownload += count
		fmt.Printf("Calculated %d tiles to download and process for zoom level %d\n", count, i)
	}

	ok := ui.YesNoPrompt("Would you like to proceed?", true)
	if ok {
		downloadBar := progressbar.Default(int64(totalFilesToDownload), "Downloading and processing tiles...")
		tiles.Download(p, downloadBar)
		createBar := progressbar.Default(int64(totalTiles), "Creating map DB...")
		tiles.CreateAtlas(atlasName, north, west, south, east, createBar)
	}
}
