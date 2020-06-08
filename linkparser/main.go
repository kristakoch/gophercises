package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {

	fPtr := flag.String("f", "./ex/ex1.html", "name of file to parse")
	flag.Parse()

	c, err := ioutil.ReadFile(*fPtr)
	if err != nil {
		log.Fatal(err)
	}
	r := strings.NewReader(string(c))

	l, err := ParseLinks(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(l)

	// lp.PrintLinks(l)
	// fmt.Println(lp.GetURLs(l))
	// fmt.Println(lp.GetLinkText(l))
}
