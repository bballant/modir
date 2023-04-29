package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	gameTime := flag.Int("t", 52, "Length of time in minutes for the game")
	formation := flag.Int("f", 322, "Formation of the game")
	flag.Parse()

	width, height := 400, 300
	fieldColor := color.White
	lineColor := color.Black
	playerColor := color.Gray{Y: 128}
	lineThickness := 3
	changesTextOffsetX := 300

	r := csv.NewReader(os.Stdin)
	rows, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	maxImages := 8
	imagesPerCol := 4
	cols := 2

	imgWidth := width*cols + changesTextOffsetX*cols
	imgHeight := height * imagesPerCol
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	draw.Draw(img, img.Bounds(), &image.Uniform{fieldColor}, image.ZP, draw.Src)

	for i, row := range rows {
		if i == 0 || i > maxImages {
			continue // skip header and limit to a maximum of 8 images
		}

		colIndex := (i - 1) / imagesPerCol
		rowIndex := (i - 1) % imagesPerCol

		offsetX := colIndex * (width + changesTextOffsetX)
		offsetY := rowIndex * height

		// Draw field, center circle, and goal boxes
		drawField(img, offsetX, offsetY, width, height, lineColor, lineThickness)

		playerRadius := 10
		playerPositions := getPositions(offsetX, offsetY, width, height, *formation)

		playerNames := make([]string, len(row))
		copy(playerNames, row)

		drawPlayers(img, playerColor, playerRadius, playerPositions[:], playerNames[:])

		addLabel(img, strconv.Itoa(i), offsetX+10, offsetY+height-10)

		drawChanges(img, offsetX+width+10, offsetY+height-10, []string{timeInGame(i, len(rows)-1, *gameTime)})

		var subs []string

		if i == 1 && len(rows) > 2 {
			nextRow := rows[2]
			for idx, name := range row {
				if nextRow[idx] != name {
					subs = append(subs, name)
				}
			}
		} else if i > 1 {
			prevRow := rows[i-1]
			maxNameLen := 0
			for idx, name := range row {
				if prevRow[idx] != name {
					if len(name) > maxNameLen {
						maxNameLen = len(name)
					}
				}
			}
			for idx, name := range row {
				if prevRow[idx] != name {
					pos := playerPositions[idx].symbol
					name = fmt.Sprintf("%s %-*s", pos, maxNameLen, name)
					subs = append(subs, name+" for "+prevRow[idx])
				}
			}
		}

		if len(subs) > 0 {
			drawChanges(img, offsetX+width+10, offsetY+20, subs)
		}
	}

	fileName := "soccer_fields.png"
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

type Position struct {
	symbol string
	x, y   int
}

func getPositions(offsetX, offsetY, width, height, formation int) []Position {
	switch formation {
	case 322:
		return []Position{
			{"GK", width/20 + offsetX, height/2 + offsetY},
			{"LB", width/4 + offsetX, height/4 + offsetY},
			{"CB", width/4 + offsetX, height/2 + offsetY},
			{"RB", width/4 + offsetX, height*3/4 + offsetY},
			{"LM", width*4/8 + offsetX, height*3/8 + offsetY},
			{"RM", width*4/8 + offsetX, height*5/8 + offsetY},
			{"LF", width*6/8 + offsetX, height*3/8 + offsetY},
			{"RF", width*6/8 + offsetX, height*5/8 + offsetY},
		}
	case 331:
		return []Position{
			{"GK", width/20 + offsetX, height/2 + offsetY},
			{"LB", width/4 + offsetX, height/2 + offsetY},
			{"CB", width/4 + offsetX, height/4 + offsetY},
			{"RB", width/4 + offsetX, height*3/4 + offsetY},
			{"LM", width*4/8 + offsetX, height/2 + offsetY},
			{"CM", width*4/8 + offsetX, height/4 + offsetY},
			{"RM", width*4/8 + offsetX, height*3/4 + offsetY},
			{"ST", width*6/8 + offsetX, height/2 + offsetY},
		}
	default:
		return []Position{}
	}
}

func timeInGame(period int, totalPeriods int, totalTime int) string {
	timeNum := 0.0
	if period > 0 {
		timeNum = float64(period) / float64(totalPeriods)
	}
	return decimalToTimeString(timeNum * float64(totalTime))
}

func decimalToTimeString(decimal float64) string {
	hours := int(decimal)
	minutes := int(math.Round((decimal - float64(hours)) * 60))
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func drawField(img *image.RGBA, offsetX, offsetY, width, height int, lineColor color.Color, lineThickness int) {
	// Draw field outline
	drawThickLine(img, lineColor, lineThickness, offsetX+0, offsetY+0, offsetX+width-1, offsetY+0)
	drawThickLine(img, lineColor, lineThickness, offsetX+width-1, offsetY+0, offsetX+width-1, offsetY+height-1)
	drawThickLine(img, lineColor, lineThickness, offsetX+width-1, offsetY+height-1, offsetX+0, offsetY+height-1)
	drawThickLine(img, lineColor, lineThickness, offsetX+0, offsetY+height-1, offsetX+0, offsetY+0)

	// Draw center circle
	drawCircle(img, lineColor, width/2+offsetX, height/2+offsetY, height/5, false)

	// Draw center line
	drawLine(img, lineColor, width/2+offsetX, offsetY, width/2+offsetX, height+offsetY)

	// left goal box
	drawLine(img, lineColor, offsetX, offsetY+60, offsetX+60, offsetY+60)
	drawLine(img, lineColor, offsetX+60, offsetY+60, offsetX+60, offsetY+height-60)
	drawLine(img, lineColor, offsetX, offsetY+height-60, offsetX+60, offsetY+height-60)

	// right goal box
	drawLine(img, lineColor, offsetX+width-60, offsetY+60, offsetX+width, offsetY+60)
	drawLine(img, lineColor, offsetX+width-60, offsetY+60, offsetX+width-60, offsetY+height-60)
	drawLine(img, lineColor, offsetX+width-60, offsetY+height-60, offsetX+width, offsetY+height-60)
}

func drawChanges(img *image.RGBA, startX, startY int, changes []string) {
	// Load font
	fontBytes, err := ioutil.ReadFile("/usr/share/fonts/truetype/cousine/Cousine Bold Italic Nerd Font Complete.ttf")
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	fontFace := truetype.NewFace(f, &truetype.Options{Size: 18})

	// Draw changes text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: fontFace,
	}

	for i, change := range changes {
		changeText := fmt.Sprintf("%s", change)
		d.Dot = fixed.Point26_6{
			X: fixed.I(startX),
			Y: fixed.I(startY + i*18),
		}
		d.DrawString(changeText)
	}
}

func drawPlayers(img *image.RGBA, c color.Color, r int, pos []Position, names []string) {
	for i, p := range pos {
		drawCircle(img, c, p.x, p.y, r, true)
		addLabel(img, names[i], p.x, p.y)
	}
}

func drawLabel(img *image.RGBA, label string, x, y int) {
	col := color.Black
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
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

func drawThickLine(img *image.RGBA, c color.Color, t, x1, y1, x2, y2 int) {
	for i := 0; i < t; i++ {
		drawLine(img, c, x1+i, y1+i, x2+i, y2+i)
	}
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

func drawCircle(img *image.RGBA, c color.Color, x, y, r int, filled bool) {
	for i := x - r; i <= x+r; i++ {
		for j := y - r; j <= y+r; j++ {
			if (i-x)*(i-x)+(j-y)*(j-y) <= r*r {
				if filled || (i-x)*(i-x)+(j-y)*(j-y) >= (r-1)*(r-1) {
					img.Set(i, j, c)
				}
			}
		}
	}
}
