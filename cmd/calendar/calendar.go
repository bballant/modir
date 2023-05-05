package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	// Accept command-line arguments for month and year
	monthFlag := flag.Int("month", int(time.Now().Month()), "Month (1-12)")
	yearFlag := flag.Int("year", time.Now().Year(), "Year (e.g., 2023)")
	flag.Parse()

	month := time.Month(*monthFlag)
	year := *yearFlag

	if month < 1 || month > 12 {
		fmt.Println("Invalid month. Please provide a month between 1 and 12.")
		os.Exit(1)
	}

	if year < 0 {
		fmt.Println("Invalid year. Please provide a positive year value.")
		os.Exit(1)
	}

	date := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	// Create a new image with size 400 x 300
	margin := 40
	topOffset := 10
	calWidth := 400
	calHeight := 300
	imageWidth := calWidth + 2*margin
	imageHeight := calHeight + 2*margin
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Fill the image with white color
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			img.Set(x, y, color.White)
		}
	}

	monthYearLabel := fmt.Sprintf("%s %d", month, year)
	addLabel(img, monthYearLabel, margin+5, 20)

	offsetX := calWidth / 7
	offsetY := calHeight / 5

	// draw vertical grid lines
	for x := 0; x < 8; x++ {
		drawLine(img, color.Black, x*offsetX+margin, margin+topOffset, x*offsetX+margin, calHeight+margin+topOffset)
	}

	// draw horizontal grid lines
	for y := 0; y < 6; y++ {
		drawLine(img, color.Black, margin, y*offsetY+margin+topOffset, calWidth+margin, y*offsetY+margin+topOffset)
	}

	// Draw the day labels
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for i, day := range days {
		addLabel(img, day, i*offsetX+margin+5, margin-5+topOffset)
	}

	// Get the date for the first day of the month
	weekday := int(date.Weekday())
	nextMonth := date.AddDate(0, 1, 0)
	numDays := int(nextMonth.Sub(date).Hours() / 24)

	// Draw the dates for the month
	for i := 1; i <= numDays; i++ {
		x := (weekday + i - 1) % 7
		y := (weekday + i - 1) / 7
		if y >= 5 {
			break // stop drawing if we reached the end of the calendar grid
		}
		addLabel(img, fmt.Sprint(i), x*offsetX+margin+5, y*offsetY+margin+15+topOffset)
	}

	// Save the image as a PNG file
	file, err := os.Create("calendar.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func getMonthInfo(month time.Month, year int) (startWeekday int, numDays int) {
	// Get the date for the first day of the month
	date := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	// Get the weekday (0 = Sunday, 1 = Monday, etc.) for the first day of the month
	startWeekday = int(date.Weekday())

	// Calculate the number of days in the month
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	numDays = int(nextMonth.Sub(date).Hours() / 24)

	return startWeekday, numDays
}

func addLabel(img *image.RGBA, label string, x, y int) {
	// Load font
	fontBytes, err := ioutil.ReadFile("/usr/share/fonts/truetype/cousine/Cousine Bold Italic Nerd Font Complete.ttf")
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	fontFace := truetype.NewFace(f, &truetype.Options{Size: 14})

	// Draw label text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: fontFace,
	}

	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d.DrawString(label)
}

func drawLine(img *image.RGBA, c color.Color, x1, y1, x2, y2 int) {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	distance := math.Sqrt(dx*dx + dy*dy)

	for t := 0.0; t <= 1.0; t += 1.0 / distance {
		x := float64(x1) + dx*t
		y := float64(y1) + dy*t
		img.Set(int(x), int(y), c)
	}
}
