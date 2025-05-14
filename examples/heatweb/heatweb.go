package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Netsocs-Team/goheatmap"
	"github.com/Netsocs-Team/goheatmap/schemes"
)

const maxInputLength = 10000

type csvpoint []string

func (c csvpoint) X() float64 {
	x, _ := strconv.ParseFloat(c[0], 64)
	return x
}

func (c csvpoint) Y() float64 {
	x, _ := strconv.ParseFloat(c[1], 64)
	return x
}

func parseInt(vals url.Values, v string, def, min, max int) int {
	rv, err := strconv.ParseInt(vals.Get(v), 10, 32)
	if err != nil || int(rv) < min || int(rv) > max {
		return def
	}
	return int(rv)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {

	vals := req.URL.Query()

	width := parseInt(vals, "w", 1024, 100, 4096)
	height := parseInt(vals, "h", 768, 100, 4096)
	dotsize := parseInt(vals, "d", 200, 1, 256)
	opacity := uint8(parseInt(vals, "o", 128, 1, 255))

	defer func() {
		err := req.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	lr := io.LimitReader(req.Body, maxInputLength)
	cr := csv.NewReader(lr)

	data := []goheatmap.DataPoint{}
	reading := true
	for reading {
		rec, err := cr.Read()
		switch err {
		case io.EOF:
			reading = false
		case nil:
			data = append(data, csvpoint(rec))
		default:
			log.Printf("Other error:  %#v", err)
			w.WriteHeader(400)
			_, err := fmt.Fprintf(w, "error reading data: %v", err)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	w.Header().Set("Content-type", "application/vnd.google-earth.kmz")
	w.WriteHeader(200)

	err := goheatmap.KMZ(image.Rect(0, 0, width, height),
		data, dotsize, opacity, schemes.AlphaFire, w)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", rootHandler)

	srv := &http.Server{
		Addr:         ":1756",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
