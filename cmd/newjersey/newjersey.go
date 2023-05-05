package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func drawLine(img *image.RGBA, c color.Color, x1, y1, x2, y2 int) {
	// Bresenham's line algorithm
	dx := x2 - x1
	dy := y2 - y1
	sx := 1
	sy := 1

	if dx < 0 {
		sx = -1
		dx = -dx
	}

	if dy < 0 {
		sy = -1
		dy = -dy
	}

	err := dx - dy

	for {
		img.Set(x1, y1, c)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err = err - dy
			x1 = x1 + sx
		}
		if e2 < dx {
			err = err + dx
			y1 = y1 + sy
		}
	}
}

func main() {

	njBorderCoordinates := []Coordinate{
		{41.3571, -74.6957},
		{41.1851, -74.7918},
		{40.9975, -75.1419},
		{40.8391, -75.0800},
		{40.6255, -75.2047},
		{40.4637, -74.9855},
		{40.2223, -74.7733},
		{39.9639, -75.1348},
		{39.8465, -75.1405},
		{39.7211, -75.3866},
		{39.4022, -75.5641},
		{38.9398, -74.9067},
		{39.2276, -74.5861},
		{39.4501, -74.3176},
		{39.6394, -74.0454},
		{40.0363, -74.1182},
		{40.2231, -73.9884},
		{40.4698, -73.9951},
		{40.5438, -74.1454},
		{40.8351, -74.3802},
	}

	// Find minimum latitude and longitude
	minLatitude := math.MaxFloat64
	minLongitude := math.MaxFloat64

	for _, coordinate := range njBorderCoordinates {
		if coordinate.Latitude < minLatitude {
			minLatitude = coordinate.Latitude
		}
		if coordinate.Longitude < minLongitude {
			minLongitude = coordinate.Longitude
		}
	}

	// Map coordinates to the new origin (0, 0)
	newCoordinates := make([]Coordinate, len(njBorderCoordinates))
	for i, coordinate := range njBorderCoordinates {
		newCoordinates[i] = Coordinate{
			Latitude:  coordinate.Latitude - minLatitude,
			Longitude: coordinate.Longitude - minLongitude,
		}
	}

	// Create an empty image with a white background
	width := 500
	height := 500
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}

	// Find the scaling factors
	scaleLat := float64(height) / (njBorderCoordinates[0].Latitude - njBorderCoordinates[10].Latitude)
	scaleLong := float64(width) / (njBorderCoordinates[19].Longitude - njBorderCoordinates[12].Longitude)

	fmt.Println("Scale ", scaleLat, scaleLong)

	scaleLat = 200
	scaleLong = 200

	// Draw lines connecting the coordinates
	for i := 0; i < len(newCoordinates)-1; i++ {
		fmt.Println(newCoordinates[i])
		x1 := int(newCoordinates[i].Longitude * scaleLong)
		y1 := int(newCoordinates[i].Latitude * scaleLat)
		x2 := int(newCoordinates[i+1].Longitude * scaleLong)
		y2 := int(newCoordinates[i+1].Latitude * scaleLat)
		fmt.Println(height-y1, x2, height-y2)
		drawLine(img, color.Black, x1, height-y1, x2, height-y2)
	}
	// Connect the last point to the first one
	x1 := int(newCoordinates[len(newCoordinates)-1].Longitude * scaleLong)
	y1 := int(newCoordinates[len(newCoordinates)-1].Latitude * scaleLat)
	x2 := int(newCoordinates[0].Longitude * scaleLong)
	y2 := int(newCoordinates[0].Latitude * scaleLat)
	drawLine(img, color.Black, x1, height-y1, x2, height-y2)

	// Save the image as a PNG file
	file, err := os.Create("new_jersey_outline.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
