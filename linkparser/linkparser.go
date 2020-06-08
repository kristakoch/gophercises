// Package linkparser contains functions to parse links from HTML files.
package main

// todo:
// refactor to pass in reader
//  instead of file name

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// ParseLinks reads in a file, parses it into a node tree,
// and returns an array of Links from the HTML document.
func ParseLinks(r *strings.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	Links := getLinks(doc)

	return Links, nil
}

// Link is a type that holds parsed link information.
type Link struct {
	Href string
	Text string
}

// PrintLinks prints a formatted list of Links.
func PrintLinks(ls []Link) {
	for i, l := range ls {
		fmt.Printf("%d. Link text: %s, URL: %s\n", i+1, l.Href, l.Text)
	}
}

// GetURLs returns a slice of only the urls
// in a list of Links.
func GetURLs(ls []Link) []string {
	var links []string
	for _, l := range ls {
		links = append(links, l.Href)
	}
	return links
}

// GetLinkText returns a slice of only the link text
// in a list of Links.
func GetLinkText(ls []Link) []string {
	var txt []string
	for _, l := range ls {
		txt = append(txt, l.Text)
	}
	return txt
}

// getLinks is a helper to get the href attribute
// values of link tags in a given node tree.
func getLinks(n *html.Node) []Link {
	var Links []Link
	isLink := false
	if n.Type == html.ElementNode && n.Data == "a" {
		isLink = true
		for _, a := range n.Attr {
			if a.Key == "href" {
				s := txtNodes(n)
				s = strings.Trim(s, " ")
				Links = append(Links, Link{Href: a.Val, Text: s})
				break
			}
		}
	}
	// Range over the node's children and recursively collect
	// a href node values.
	if !isLink {
		c := n.FirstChild
		for c != nil {
			Links = append(Links, getLinks(c)...)
			c = c.NextSibling
		}
	}
	return Links
}

// txtNodes is a helper to get the text node
// values in a given node tree.
func txtNodes(n *html.Node) string {
	s := ""
	if n.Type == html.TextNode {
		s = fmt.Sprintf("%v ", strings.Trim(n.Data, " \n\t"))
	}

	// Range over the node's children and recursively collect
	// text node values.
	c := n.FirstChild
	for c != nil {
		s += txtNodes(c)
		c = c.NextSibling
	}

	return s
}
