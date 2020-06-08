// A minimal program to write SVGs to a web browser.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	svg "github.com/ajstarks/svgo"
)

func main() {
	http.Handle("/circle", http.HandlerFunc(circleHandler))
	http.Handle("/super-ball", http.HandlerFunc(superBallHandler))
	http.Handle("/graph", http.HandlerFunc(graphHandler))

	fmt.Println("server is running at http://localhost:2003")

	log.Fatal(http.ListenAndServe(":2003", nil))
}

func graphHandler(w http.ResponseWriter, req *http.Request) {
	fb, err := ioutil.ReadFile("image.svg")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "image/svg+xml")

	bb := bytes.NewBuffer(fb)
	io.Copy(w, bb)
}

func circleHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(500, 500)
	s.Circle(250, 250, 250, "fill:none;stroke:black")
	s.End()
}

func superBallHandler(w http.ResponseWriter, req *http.Request) {
	f, err := os.Create("image.svg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	width := 500
	height := 500

	lg := []svg.Offcolor{
		{Offset: 0, Color: "red", Opacity: .5},
		{Offset: 100, Color: "purple", Opacity: .5},
	}

	g := svg.New(f)
	g.Start(width, height)
	g.Title("Gradients")
	g.LinearGradient("h", 0, 0, 100, 0, lg)
	g.LinearGradient("v", 0, 0, 0, 100, lg)

	// g.Rect(0, 0, width-10, height-10, "fill:url(#h)")
	// g.Rect(10, 10, width, height, "fill:url(#v)")

	g.Circle(width/2, height/2, width/2, "fill:url(#v)")
	g.Circle(width/2, height/2, width/2, "fill:url(#h)")

	g.End()

	// Use the rsvg convert cmd to create a png image from the svg.
	convertCmd := "rsvg-convert image.svg > image-from-svg.png"
	cmd := exec.Command("bash", "-c", convertCmd)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}
