# OSMDroid SQLite Map Creator

## Running
```
go run main.go --maxZoom=8 --baseMapUrl='https://nwy-tiles-api.prod.newaydata.com/tiles/{z}/{x}/{y}.jpg?path=2208/base/latest' --overlayUrl='https://nwy-tiles-api.prod.newaydata.com/tiles/{z}/{x}/{y}.png?path=2208/aero/latest' --quality 70 --country CZ
```

## Links
https://sourceforge.net/p/mobac/code/HEAD/tree/trunk/MOBAC/mobac/src/main/java/mobac/program/atlascreators/OsmdroidSQLite.java#l128

https://www.latlong.net/

https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#ECMAScript_.28JavaScript.2FActionScript.2C_etc..29

https://pkg.go.dev/github.com/mattn/go-sqlite3

https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance