package download

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/andrejsoucek/map-creator/bounds"
	"github.com/andrejsoucek/map-creator/byte2image"
	"github.com/andrejsoucek/map-creator/flags"
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

type Downloader struct {
	DownloadParams DownloadParams
	Converter      *Converter
	Bar            *progressbar.ProgressBar
}

func (d *Downloader) Download() {
	xyzs := createXyzs(
		d.Converter,
		d.DownloadParams.MinZoom,
		d.DownloadParams.MaxZoom,
		d.DownloadParams.North,
		d.DownloadParams.West,
		d.DownloadParams.South,
		d.DownloadParams.East,
	)
	baseDecoder := byte2image.NewDecoder(d.DownloadParams.BaseMapUrl)
	overlayDecoder := byte2image.NewDecoder(d.DownloadParams.OverlayUrl)
	var wg sync.WaitGroup
	for _, xyz := range xyzs {
		wg.Add(1)
		baseUrl := formatUrl(d.DownloadParams.BaseMapUrl, xyz)
		overlayUrl := formatUrl(d.DownloadParams.OverlayUrl, xyz)
		go get(&wg, baseUrl, baseDecoder, overlayUrl, overlayDecoder, xyz, d.DownloadParams.Quality, d.Bar)
	}
	wg.Wait()
}

func CreateDownloadParams(f flags.Flags) DownloadParams {
	var c bounds.CountryBounds
	north := f.North
	west := f.West
	south := f.South
	east := f.East
	if f.Country != "" {
		switch f.Country {
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
	return DownloadParams{
		BaseMapUrl: f.BaseMapUrl,
		OverlayUrl: f.OverlayUrl,
		MinZoom:    f.MinZoom,
		MaxZoom:    f.MaxZoom,
		North:      north,
		West:       west,
		South:      south,
		East:       east,
		Quality:    f.Quality,
	}
}

func formatUrl(url string, xyz xyz) string {
	url = strings.Replace(url, "{x}", strconv.Itoa(xyz.x), 1)
	url = strings.Replace(url, "{y}", strconv.Itoa(xyz.y), 1)
	url = strings.Replace(url, "{z}", strconv.Itoa(xyz.z), 1)

	return url
}

func createXyzs(c *Converter, minZoom int, maxZoom int, north float64, west float64, south float64, east float64) []xyz {
	xyzs := []xyz{}
	for z := minZoom; z <= maxZoom; z++ {
		left := c.Lon2tile(west, z)
		right := c.Lon2tile(east, z)
		for x := left; x <= right; x++ {
			top := c.Lat2tile(north, z)
			bottom := c.Lat2tile(south, z)
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
	client := &http.Client{}

	resBase, err := client.Get(baseUrl)
	if resBase != nil {
		defer resBase.Body.Close()
	}
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
		if resOverlay != nil {
			defer resOverlay.Body.Close()
		}
		if err != nil {
			log.Fatal(err)
		}

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

	jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
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
