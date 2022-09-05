package tiles

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
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
	Quality                int
}

type xyz struct {
	x int
	y int
	z int
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
	xyzs := createXyzs(dp.MinZoom, dp.MaxZoom, dp.North, dp.West, dp.South, dp.East)
	var wg sync.WaitGroup
	// limit to four downloads at a time, this is called a semaphore
	limiter := make(chan struct{}, dp.MaxConcurrentDownloads)
	for _, xyz := range xyzs {
		wg.Add(1)
		baseUrl := formatUrl(dp.BaseMapUrl, xyz)
		overlayUrl := formatUrl(dp.OverlayUrl, xyz)
		go get(&wg, limiter, baseUrl, overlayUrl, xyz, dp.Quality, bar)
	}
	wg.Wait()
}

func formatUrl(url string, xyz xyz) string {
	url = strings.Replace(url, "{x}", strconv.Itoa(xyz.x), 1)
	url = strings.Replace(url, "{y}", strconv.Itoa(xyz.y), 1)
	url = strings.Replace(url, "{z}", strconv.Itoa(xyz.z), 1)

	return url
}

func createXyzs(minZoom int, maxZoom int, north float64, west float64, south float64, east float64) []xyz {
	xyzs := []xyz{}
	for z := minZoom; z <= maxZoom; z++ {
		left := convert.Lon2tile(west, z)
		right := convert.Lon2tile(east, z)
		for x := left; x <= right; x++ {
			top := convert.Lat2tile(north, z)
			bottom := convert.Lat2tile(south, z)
			for y := top; y <= bottom; y++ {
				tile := xyz{
					x: x,
					y: y,
					z: z,
				}
				xyzs = append(xyzs, tile)
			}
		}
	}
	return xyzs
}

func get(wg *sync.WaitGroup, sema chan struct{}, baseUrl string, overlayUrl string, xyz xyz, quality int, bar *progressbar.ProgressBar) {
	sema <- struct{}{}
	defer func() {
		<-sema
		wg.Done()
	}()

	client := &http.Client{Timeout: 5 * time.Second}

	resBase, err := client.Get(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resBase.Body.Close()

	base, err := jpeg.Decode(resBase.Body)
	if err != nil {
		log.Fatalf("failed to decode base image: %s", err)
	}
	b := base.Bounds()
	output := image.NewRGBA(b)
	draw.Draw(output, b, base, image.ZP, draw.Src)

	if overlayUrl != "" {
		resOverlay, err := client.Get(overlayUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer resOverlay.Body.Close()

		overlay, err := png.Decode(resOverlay.Body)
		if err != nil {
			log.Fatalf("failed to decode overlay image: %s", err)
		}

		draw.Draw(output, b, overlay, image.ZP, draw.Over)
	}
	saveFile(output, xyz.z, fmt.Sprintf("%d-%d", xyz.x, xyz.y), quality)
	bar.Add(1)
	return
}

func saveFile(img *image.RGBA, zoom int, fileName string, quality int) {
	cwd, _ := os.Getwd()
	folderName := "tmp"
	os.MkdirAll(folderName+"/"+"/"+strconv.Itoa(zoom), os.ModePerm)

	path := filepath.Join(cwd, folderName, strconv.Itoa(zoom), fileName)
	newFilePath := filepath.FromSlash(path)
	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
}
