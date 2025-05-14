package main

import (
	"bufio"
	"encoding/json"
	"image"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Netsocs-Team/goheatmap"
	"github.com/Netsocs-Team/goheatmap/schemes"
)

const url = "http://single.couchbase.net/wikipedia2012/_design/places/_spatial/points?bbox=-180,-90,180,90"

type row struct {
	ID       string `json:"id"`
	Value    int    `json:"value"`
	Geometry struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

func main() {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("error opening strem from %q: %v", url, err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalf("error when closing the body: %v", err)
		}
	}()

	r := bufio.NewReader(res.Body)
	_, _, err = r.ReadLine()
	if err != nil {
		log.Fatalf("error reading first line: %v", err)
	}

	locations := make([]goheatmap.DataPoint, 0, 200000)

	for {
		bytes, isPrefix, err := r.ReadLine()
		if err != nil {
			log.Fatalf("error reading line: %v", err)
		}
		if isPrefix {
			log.Fatalf("crap, that was a prefix...")
		}
		if bytes[0] != '{' {
			break
		}
		r := row{}
		if bytes[len(bytes)-1] == ',' {
			bytes = bytes[:len(bytes)-1]
		}
		err = json.Unmarshal(bytes, &r)
		if err != nil {
			log.Printf("couldn't parse %v: %v", string(bytes), err)
			break
		}

		locations = append(locations, goheatmap.P(r.Geometry.Coordinates[0],
			r.Geometry.Coordinates[1]))
	}
	end := time.Now()

	log.Printf("parsed %d items in %s",
		int64(len(locations)), start)

	out, err := os.Create("wikipedia.kmz")
	if err != nil {
		log.Fatalf("error making output file:  %v", err)
	}

	defer func() {
		err := out.Close()
		if err != nil {
			log.Fatalf("error when closing 'wikipedia.kmz' doc: %v", err)
		}
	}()

	err = goheatmap.KMZ(image.Rect(0, 0, 8192, 4096), locations, 50, 96,
		schemes.AlphaFire, out)
	if err != nil {
		log.Fatalf("error generating thingy: %v", err)
	}

	log.Printf("completed heatmap generation in %s", time.Since(end))
}
