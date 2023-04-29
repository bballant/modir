package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func decodeImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".png":
		return png.Decode(f)
	case ".jpg", ".jpeg":
		return jpeg.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run stitch_images.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	var imageFiles []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			imageFiles = append(imageFiles, filepath.Join(dir, file.Name()))
		}
	}

	count := len(imageFiles)
	if count == 0 {
		fmt.Println("No PNG or JPEG files found in the directory.")
		os.Exit(0)
	}

	gridSize := int(math.Sqrt(float64(count)))
	totalImages := gridSize * gridSize
	imageFiles = imageFiles[:totalImages]

	var images []image.Image
	for _, file := range imageFiles {
		img, err := decodeImage(file)
		if err != nil {
			fmt.Println("Error decoding image:", err)
			os.Exit(1)
		}
		images = append(images, img)
	}

	// Create an empty image to stitch the images together
	border := 5
	imgWidth := images[0].Bounds().Dx() + 2*border
	imgHeight := images[0].Bounds().Dy() + 2*border
	finalWidth := imgWidth * gridSize
	finalHeight := imgHeight * gridSize

	finalImage := image.NewRGBA(image.Rect(0, 0, finalWidth, finalHeight))
	draw.Draw(finalImage, finalImage.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			index := y*gridSize + x
			rect := image.Rect(x*imgWidth+border, y*imgHeight+border, (x+1)*imgWidth, (y+1)*imgHeight)
			draw.Draw(finalImage, rect, images[index], image.Point{}, draw.Src)
		}
	}

	outputFile := "stitched_images.png"
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating the final image:", err)
		os.Exit(1)
	}
	defer f.Close()

	err = png.Encode(f, finalImage)
	if err != nil {
		fmt.Println("Error saving the final image:", err)
		os.Exit(1)
	}

	fmt.Printf("Stitched images saved as %s\n", outputFile)
}
