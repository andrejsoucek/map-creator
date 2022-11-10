package tiles

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andrejsoucek/map-creator/byte2image"
	"github.com/andrejsoucek/map-creator/convert"
	"github.com/schollz/progressbar/v3"
)

type DownloadParams struct {
	BaseMapUrl string
	OverlayUrl string
	MinZoom    int
	MaxZoom    int
	North      float64
	West       float64
	South      float64
	East       float64
	Quality    int
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
	baseDecoder := byte2image.NewDecoder(dp.BaseMapUrl)
	overlayDecoder := byte2image.NewDecoder(dp.OverlayUrl)
	var wg sync.WaitGroup
	for _, xyz := range xyzs {
		wg.Add(1)
		baseUrl := formatUrl(dp.BaseMapUrl, xyz)
		overlayUrl := formatUrl(dp.OverlayUrl, xyz)
		go get(&wg, baseUrl, baseDecoder, overlayUrl, overlayDecoder, xyz, dp.Quality, bar)
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

func get(
	wg *sync.WaitGroup,
	baseUrl string,
	baseDecoder byte2image.Decoder,
	overlayUrl string,
	overlayDecoder byte2image.Decoder,
	xyz xyz,
	quality int,
	bar *progressbar.ProgressBar,
) {
	defer func() {
		wg.Done()
	}()

	fileName := fmt.Sprintf("%d-%d", xyz.x, xyz.y)
	exists := fileExists("output/tiles" + "/" + strconv.Itoa(xyz.z) + "/" + fileName)
	if exists {
		log.Println("File already exists: " + fileName)
		bar.Add(1)
		if overlayUrl != "" {
			bar.Add(1)
		}
		return
	}

	log.Println("Downloading: " + fileName)
	client := &http.Client{Timeout: 5 * time.Second}

	resBase, err := client.Get(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resBase.Body.Close()

	base, err := baseDecoder.Decode(resBase.Body)
	if err != nil {
		log.Fatalf("failed to decode base image, error: %s, url: %s", err, baseUrl)
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

		overlay, err := overlayDecoder.Decode(resOverlay.Body)
		if err != nil {
			log.Fatalf("failed to decode overlay image, error: %s, url: %s", err, overlayUrl)
		}

		draw.Draw(output, b, overlay, image.ZP, draw.Over)
		bar.Add(1)
	}
	saveFile(output, xyz.z, fileName, quality)
	bar.Add(1)
	return
}

func saveFile(img *image.RGBA, zoom int, fileName string, quality int) {
	cwd, _ := os.Getwd()
	folderName := "output/tiles"
	os.MkdirAll(folderName+"/"+"/"+strconv.Itoa(zoom), os.ModePerm)

	path := filepath.Join(cwd, folderName, strconv.Itoa(zoom), fileName)
	newFilePath := filepath.FromSlash(path)
	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
