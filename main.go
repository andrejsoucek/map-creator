package main

import (
	"flag"
	"fmt"

	"github.com/andrejsoucek/map-creator/tiles"
	"github.com/andrejsoucek/map-creator/ui"
	progressbar "github.com/schollz/progressbar/v3"
)

var baseMapUrl string
var overlayUrl string
var minZoom int
var maxZoom int
var north float64
var west float64
var south float64
var east float64
var maxConcurrentDownloads int
var quality int

func parseFlags() {
	flag.StringVar(&baseMapUrl, "baseMapUrl", "https://nwy-tiles-api.prod.newaydata.com/tiles/{z}/{x}/{y}.png?path=2208/aero/latest", "map URL")
	flag.StringVar(&overlayUrl, "overlayUrl", "", "map URL")
	flag.IntVar(&minZoom, "minZoom", 8, "min zoom")
	flag.IntVar(&maxZoom, "maxZoom", 12, "max zoom")
	flag.Float64Var(&north, "n", 51.0835, "north bounding point")
	flag.Float64Var(&west, "w", 12.0475, "west bounding point")
	flag.Float64Var(&south, "s", 48.55, "south bounding point")
	flag.Float64Var(&east, "e", 18.9, "east bounding point")
	flag.IntVar(&maxConcurrentDownloads, "maxConcurrentDownloads", 4, "max concurrent downloads")
	flag.IntVar(&quality, "quality", 100, "quality")
	flag.Parse()
}

func main() {
	parseFlags()

	totalTiles := 0
	totalFilesToDownload := 0
	p := tiles.DownloadParams{
		BaseMapUrl:             baseMapUrl,
		OverlayUrl:             overlayUrl,
		MinZoom:                minZoom,
		MaxZoom:                maxZoom,
		North:                  north,
		West:                   west,
		South:                  south,
		East:                   east,
		MaxConcurrentDownloads: maxConcurrentDownloads,
	}

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
		downloadBar := progressbar.Default(int64(totalFilesToDownload), "Downloading tiles...")
		tiles.Download(p, downloadBar)
		// processBar := progressbar.Default(int64(totalFilesToDownload), "Processing tiles...")
		createBar := progressbar.Default(int64(totalTiles), "Creating map DB...")
		tiles.CreateAtlas(north, west, south, east, createBar)
	}
}
