# OSMDroid SQLite Map Creator

## Running
### Flags
```
  -atlasName string
        name of the generated atlas
  -baseMapUrl string
        map URL
  -baseMapType string (jpeg or png)
        image type of the tiles for the base map
  -country string
        country code, supported: CZ, SK, overrides manual coordinates (default "CZ")
  -maxConcurrentDownloads int
        max concurrent downloads (default 4)
  -maxZoom int
        max zoom (default 12)
  -minZoom int
        min zoom (default 8)
  -overlayUrl string
        map URL
  -overlayType string (jpeg or png)
        image type of the tiles for the overlay map
  -quality int
        quality (default 100)
  -n float
        north bounding point
  -s float
        south bounding point
  -w float
        west bounding point
  -e float
        east bounding point
```
### Example
```
go run main.go --maxZoom=8 --baseMapUrl='https://nwy-tiles-api.prod.newaydata.com/tiles/{z}/{x}/{y}.jpg?path=2208/base/latest' --overlayUrl='https://nwy-tiles-api.prod.newaydata.com/tiles/{z}/{x}/{y}.png?path=2208/aero/latest' --quality 70 --country CZ
```

## TODO
- creating atlas in a goroutine
- architecture refactoring

## Links
### Inspired by
 https://sourceforge.net/p/mobac/code/HEAD/tree/trunk/MOBAC/mobac/src/main/java/mobac/program/atlascreators/OsmdroidSQLite.java#l128
 https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#ECMAScript_.28JavaScript.2FActionScript.2C_etc..29

### Coordinates picker
https://www.latlong.net/

