package tiles

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func CreateAtlas(north float64, west float64, south float64, east float64, bar *progressbar.ProgressBar) {
	file, err := os.Open("./tmp")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	zooms, err := file.Readdirnames(0)
	for _, z := range zooms {
		files, err := ioutil.ReadDir("./tmp/" + z)
		if err != nil {
			log.Fatal(err)
		}
		for _, tile := range files {
			split := strings.Split(tile.Name(), "-")
			x, err := strconv.Atoi(split[0])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(split[1])
			if err != nil {
				log.Fatal(err)
			}
			z, err := strconv.Atoi(z)
			if err != nil {
				log.Fatal(err)
			}
			index := (((z << z) + x) << z) + y
			println(tile.Name(), index)
			bar.Add(1)
		}
	}
}
