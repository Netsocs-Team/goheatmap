package main

import (
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"

	"github.com/Netsocs-Team/goheatmap"
	"github.com/Netsocs-Team/goheatmap/schemes"
)

func main() {
	points := []goheatmap.DataPoint{}
	for range 350 {
		points = append(points,
			goheatmap.P(rand.Float64(), rand.Float64()))
	}

	// scheme, _ := schemes.FromImage("../schemes/fire.png")
	scheme := schemes.AlphaFire

	img := goheatmap.Heatmap(image.Rect(0, 0, 1024, 1024),
		points, 150, 128, scheme)

	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}
