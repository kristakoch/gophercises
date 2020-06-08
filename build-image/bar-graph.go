package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"

	svg "github.com/ajstarks/svgo"
	"github.com/pbnjay/pixfont"
)

func main() {

	data := []monthlyData{
		{"Jan", 100},
		{"Feb", 330},
		{"Mar", 230},
		{"Apr", 140},
		{"May", 300},
		{"Jun", 200},
		{"Jul", 80},
		{"Aug", 110},
		{"Sep", 40},
		{"Oct", 310},
		{"Nov", 230},
		{"Dec", 130},
	}

	// log.Printf("creating png bar graph")
	// makePngBarGraph(data)

	// log.Printf("drawing png bar graph")
	// drawPngBarGraph(data)

	log.Printf("creating svg bar graph")
	makeSvgBarGraph(data)

}

type monthlyData struct {
	month  string
	amount int
}

var (
	blue   = color.NRGBA{R: 0, G: 166, B: 221, A: 255}
	white  = color.NRGBA{R: 255, G: 255, B: 255, A: 230}
	orange = color.NRGBA{R: 255, G: 150, B: 20, A: 200}
	black  = color.Black
)
var padding = 10

func makeSvgBarGraph(data []monthlyData) error {
	f, err := os.Create("image.svg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	width, height := 800, 400

	// Create the background.
	g := svg.New(f)
	g.Start(width, height)
	g.Rect(0, 0, width, height, "fill:white")

	// Create the bars
	barWidth := width / len(data)
	for i, md := range data {

		// Make the bar.
		g.Rect(
			barWidth*i+padding, height-md.amount, barWidth-padding, md.amount-padding, "stroke:none;fill:#2FCEDC",
		)

		// Add the value label above the bar.
		g.Text(
			barWidth*i+padding, height-md.amount-padding, fmt.Sprintf("%d", md.amount), "stroke:black; fill:black",
		)

		// Add the month label within the bar.
		g.Text(
			barWidth*i+padding*2, height-padding*2, md.month, "stroke:white; fill:white",
		)

	}

	avg := averageMonthlyData(data)
	g.Rect(
		0, height-avg, width, 2, "fill:orange",
	)

	g.End()

	convertCmd := "rsvg-convert image.svg > graph-png-from-svg.png"
	cmd := exec.Command("bash", "-c", convertCmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func drawPngBarGraph(data []monthlyData) error {

	// Create a background with these dimensions.
	var width, height = 800, 400

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	barWidth := width / len(data)
	for i, md := range data {
		// Draw the bar.
		draw.Draw(
			m, image.Rect(i*barWidth+padding, height-padding, (i*barWidth)+barWidth,
				height-md.amount), &image.Uniform{blue}, image.ZP, draw.Src,
		)

		// Add the value label above the bar.
		pixfont.DrawString(
			m, i*barWidth+padding, height-md.amount-padding, fmt.Sprintf("%d", md.amount), black,
		)

		// Add the month label within the bar.
		pixfont.DrawString(
			m, i*barWidth+(padding*2), height-23, md.month, color.White,
		)

	}

	avg := averageMonthlyData(data)
	draw.Draw(m, image.Rect(0, height-avg-1, width, height-avg+1), &image.Uniform{orange}, image.ZP, draw.Src)

	createPngImage(m, "image.jpg")

	return nil
}

func makePngBarGraph(data []monthlyData) error {

	// Create a background with these dimensions.
	var width, height = 800, 400

	rec := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rec)
	fillRectangleByPixel(img, 0, 0, width, height, white)

	// Create the bar graph bars.
	barWidth := width / len(data)

	for i, md := range data {

		// Add the bar.
		fillRectangleByPixel(
			img, i*barWidth+padding, height-md.amount, (i*barWidth)+barWidth, height-padding, blue,
		)

		// Add the value label above the bar.
		pixfont.DrawString(
			img, i*barWidth+padding, height-md.amount-10, fmt.Sprintf("%d", md.amount), black,
		)

		// Add the month label within the bar.
		pixfont.DrawString(
			img, i*barWidth+(padding*2), height-23, md.month, color.White,
		)
	}

	// Draw the line for the average of the data points.
	avg := averageMonthlyData(data)
	fillRectangleByPixel(img, 0, height-avg-1, width, height-avg+1, orange)

	// Use the generated image to create a png.
	if err := createPngImage(img, "image.png"); err != nil {
		return err
	}

	return nil
}

func fillRectangleByPixel(
	img *image.RGBA, xStart, yStart, xEnd, yEnd int, rgbColor color.Color,
) {
	for x := xStart; x < xEnd; x++ {
		for y := yStart; y < yEnd; y++ {
			img.Set(x, y, rgbColor)
		}
	}
}

func averageMonthlyData(data []monthlyData) int {
	var sum int
	for _, md := range data {
		sum += md.amount
	}

	return sum / len(data)
}

func createPngImage(img image.Image, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
