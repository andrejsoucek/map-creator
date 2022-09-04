package tiles

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andrejsoucek/map-creator/convert"
	"github.com/schollz/progressbar/v3"
)

type DownloadParams struct {
	BaseMapUrl             string
	OverlayUrl             string
	MinZoom                int
	MaxZoom                int
	North                  float64
	West                   float64
	South                  float64
	East                   float64
	MaxConcurrentDownloads int
}

type tileParam struct {
	url      string
	tileType TileType
	zoom     int
	fileName string
}

func CalculateTiles(zoom int, north float64, west float64, south float64, east float64) int {
	top := convert.Lat2tile(north, zoom)
	left := convert.Lon2tile(west, zoom)
	bottom := convert.Lat2tile(south, zoom)
	right := convert.Lon2tile(east, zoom)
	width := math.Abs(float64(left-right)) + 1
	height := math.Abs(float64(top-bottom)) + 1

	return int(width * height)
}

func Download(
	dp DownloadParams,
	bar *progressbar.ProgressBar,
) {
	tileParams := createTileParams(Base, dp.BaseMapUrl, dp.MinZoom, dp.MaxZoom, dp.North, dp.West, dp.South, dp.East)
	if dp.OverlayUrl != "" {
		tileParams = append(
			tileParams,
			createTileParams(Overlay, dp.OverlayUrl, dp.MinZoom, dp.MaxZoom, dp.North, dp.West, dp.South, dp.East)...,
		)
	}
	var wg sync.WaitGroup
	// limit to four downloads at a time, this is called a semaphore
	limiter := make(chan struct{}, dp.MaxConcurrentDownloads)
	for _, tP := range tileParams {
		wg.Add(1)
		go get(&wg, limiter, tP, bar)
	}
	wg.Wait()
}

func createTileParams(tileType TileType, mapUrl string, minZoom int, maxZoom int, north float64, west float64, south float64, east float64) []tileParam {
	tiles := []tileParam{}
	for z := minZoom; z <= maxZoom; z++ {
		left := convert.Lon2tile(west, z)
		right := convert.Lon2tile(east, z)
		for x := left; x <= right; x++ {
			top := convert.Lat2tile(north, z)
			bottom := convert.Lat2tile(south, z)
			for y := top; y <= bottom; y++ {
				url := strings.Replace(mapUrl, "{x}", strconv.Itoa(x), 1)
				url = strings.Replace(url, "{y}", strconv.Itoa(y), 1)
				url = strings.Replace(url, "{z}", strconv.Itoa(z), 1)
				tile := tileParam{
					url:      url,
					tileType: tileType,
					zoom:     z,
					fileName: fmt.Sprintf("%d-%d", x, y),
				}
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles
}

func get(wg *sync.WaitGroup, sema chan struct{}, tP tileParam, bar *progressbar.ProgressBar) {
	sema <- struct{}{}
	defer func() {
		<-sema
		wg.Done()
	}()

	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(tP.url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var buf bytes.Buffer
	// I'm copying to a buffer before writing it to file
	// I could also just use IO copy to write it to the file
	// directly and save memory by dumping to the disk directly.
	io.Copy(&buf, res.Body)
	// write the bytes to file
	saveFile(&buf, tP.zoom, tP.fileName, tP.tileType)
	bar.Add(1)
	return
}

func saveFile(buf *bytes.Buffer, zoom int, fileName string, tileType TileType) {
	cwd, _ := os.Getwd()
	folderName := "tmp"
	os.MkdirAll(folderName+"/"+tileType.toString()+"/"+strconv.Itoa(zoom), os.ModePerm)

	path := filepath.Join(cwd, folderName, tileType.toString(), strconv.Itoa(zoom), fileName)
	newFilePath := filepath.FromSlash(path)
	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	file.Write(buf.Bytes())
}
