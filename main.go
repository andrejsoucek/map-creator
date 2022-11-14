package main

import (
	"fmt"

	"github.com/andrejsoucek/map-creator/download"
	"github.com/andrejsoucek/map-creator/flags"
	"github.com/andrejsoucek/map-creator/tiles"
	"github.com/andrejsoucek/map-creator/ui"
	progressbar "github.com/schollz/progressbar/v3"
)

func main() {
	f := flags.ParseFlags()

	totalTiles := 0
	totalFilesToDownload := 0

	dp := download.CreateDownloadParams(f)
	c := download.Converter{}
	calc := download.NewTileCalculator(&c)

	for i := f.MinZoom; i <= f.MaxZoom; i++ {
		count := calc.CalculateTiles(i, dp.North, dp.West, dp.South, dp.East)
		totalTiles = totalTiles + count
		if dp.OverlayUrl != "" {
			count = count * 2
		}
		totalFilesToDownload += count
		fmt.Printf("Calculated %d tiles to download and process for zoom level %d\n", count, i)
	}

	ok := ui.YesNoPrompt("Would you like to proceed?", true)
	if ok {
		d := download.Downloader{
			DownloadParams: dp,
			Converter:      download.Converter{},
			Bar:            progressbar.Default(int64(totalFilesToDownload), "Downloading and processing tiles..."),
		}
		d.Download()
		createBar := progressbar.Default(int64(totalTiles), "Creating map DB...")
		tiles.CreateAtlas(f.AtlasName, f.North, f.West, f.South, f.East, createBar)
	}
}
