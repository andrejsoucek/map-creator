package tiles

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

func CreateAtlas(atlasName string, north float64, west float64, south float64, east float64, bar *progressbar.ProgressBar) {
	file, err := os.Open("./tmp")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s.sqlite", atlasName))
	checkErr(err)

	defer file.Close()
	defer finalizeDb(db)

	configureDb(db)
	createTable(db)

	zooms, err := file.Readdirnames(0)
	for _, z := range zooms {
		dir := "./tmp/" + z
		fileInfos, err := ioutil.ReadDir(dir)
		checkErr(err)

		for _, fileInfo := range fileInfos {
			split := strings.Split(fileInfo.Name(), "-")
			x, err := strconv.Atoi(split[0])
			checkErr(err)

			y, err := strconv.Atoi(split[1])
			checkErr(err)

			z, err := strconv.Atoi(z)
			checkErr(err)

			tile, err := os.ReadFile(dir + "/" + fileInfo.Name())
			checkErr(err)

			index := (((z << z) + x) << z) + y
			insert(db, index, tile, atlasName)
			bar.Add(1)
		}
	}
}

func finalizeDb(db *sql.DB) {
	db.Exec("PRAGMA journal_mode=DELETE")
	db.Close()
}

func configureDb(db *sql.DB) {
	db.Exec("PRAGMA journal_mode=WAL")
	db.SetMaxOpenConns(1)
}

func createTable(db *sql.DB) {
	st, err := db.Prepare("CREATE TABLE IF NOT EXISTS tiles (key INTEGER PRIMARY KEY, provider TEXT, tile BLOB)")
	checkErr(err)

	_, err = st.Exec()
	checkErr(err)
}

func insert(db *sql.DB, index int, tile []byte, provider string) {
	st, err := db.Prepare("INSERT OR REPLACE INTO tiles (key, provider, tile) VALUES (?, ?, ?)")
	checkErr(err)

	_, err = st.Exec(index, provider, tile)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
