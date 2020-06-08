package main

import (
	"flag"
	"log"
)

func main() {
	vPtr := flag.String("v", "html", "either html or cmd")
	flag.Parse()

	var err error
	switch *vPtr {
	case "cmd":
		err = ServeCmdLnStory()
	default:
		err = ServeHTMLStory()
	}
	if err != nil {
		log.Fatal(err)
	}
}
