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
	// Cluster a few points randomly around Lawrence station.
	lawrence := goheatmap.P(-121.996158, 37.370713)

	points := []goheatmap.DataPoint{lawrence}
	for range 35 {
		points = append(points,
			goheatmap.P(lawrence.X()+(rand.Float64()/100.0)-0.005,
				lawrence.Y()+(rand.Float64()/100.0)-0.005))
	}

	kmlout, err := os.Create("test.kml")
	if err != nil {
		log.Fatalf("Error creating kml file:  %v", err)
	}
	defer func() {
		err := kmlout.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	imgout, err := os.Create("test.png")
	if err != nil {
		log.Fatalf("Error creating image file:  %v", err)
	}
	defer func() {
		err := imgout.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	img, err := goheatmap.KML(image.Rect(0, 0, 1024, 1024),
		points, 200, 128, schemes.AlphaFire, "test.png", kmlout)
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(imgout, img)
	if err != nil {
		log.Fatal(err)
	}
}
