package main

import (
	"image"
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

	kmzout, err := os.Create("test.kmz")
	if err != nil {
		log.Fatalf("Error creating kml file:  %v", err)
	}
	defer func() {
		err := kmzout.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = goheatmap.KMZ(image.Rect(0, 0, 1024, 1024),
		points, 200, 128, schemes.AlphaFire, kmzout)
	if err != nil {
		log.Fatalf("Error creating heatmap: %v", err)
	}
}
