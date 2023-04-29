package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func drawPetal(img draw.Image, x, y, width, height, angle float64) {
	var points []image.Point
	angleRad := angle * math.Pi / 180

	for i := 0.0; i <= 1.0; i += 0.01 {
		angle := 2 * math.Pi * i
		r := width * height / (math.Sqrt(math.Pow(height*math.Cos(angle), 2) + math.Pow(width*math.Sin(angle), 2)))
		points = append(points, image.Pt(
			int(x+r*math.Cos(angle+angleRad)),
			int(y+r*math.Sin(angle+angleRad)),
		))
	}

	for i := range points {
		if i > 0 {
			x1, y1 := points[i-1].X, points[i-1].Y
			x2, y2 := points[i].X, points[i].Y
			drawLine(img, x1, y1, x2, y2, color.Black)
		}
	}
}

func drawLine(img draw.Image, x1, y1, x2, y2 int, col color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := float64(x1)*(1-t) + float64(x2)*t
		y := float64(y1)*(1-t) + float64(y2)*t
		img.Set(int(x), int(y), col)
	}
}

func main() {
	width := 500
	height := 500

	// Create a blank canvas with the specified dimensions
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// Set the background color to white
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	// Draw the flower's center
	centerX, centerY := width/2, height/2
	for y := centerY - 20; y <= centerY+20; y++ {
		for x := centerX - 20; x <= centerX+20; x++ {
			if (x-centerX)*(x-centerX)+(y-centerY)*(y-centerY) <= 20*20 {
				canvas.Set(x, y, color.Black)
			}
		}
	}

	// Draw the petals
	petalWidth := 30.0
	petalHeight := 50.0
	for i := 0; i < 8; i++ {
		angle := float64(i) * 45
		drawPetal(canvas, float64(centerX), float64(centerY)-60, petalWidth, petalHeight, angle)
	}

	// Save the canvas as a PNG file
	file, err := os.Create("simple_flower_no_lib.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, canvas)
	if err != nil {
		panic(err)
	}
}
