package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	rootURL := flag.String("root", "https://www.calhoun.io/", "root URL to create sitemap from")
	flag.Parse()

	// Build a string of the sitemap of domain from a URL.
	sMap, err := SiteMap(*rootURL)
	if err != nil {
		log.Fatal(err)
	}

	// Print the map.
	fmt.Println(sMap)
}
