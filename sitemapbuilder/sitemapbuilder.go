// Package sitemapbuilder builds an XML sitemap of a website.
package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Page represents one page in the domain.
type Page struct {
	XMLName xml.Name `xml:"url"`
	URL     string   `xml:"loc"`
}

// XMLMap represents all pages in the domain.
type XMLMap struct {
	XMLName   xml.Name `xml:"urlset"`
	Namespace string   `xml:"xmlns,attr"`
	Pages     []Page
}

type empty struct{}

// SiteMap takes in a root URL and returns an XML sitemap.
func SiteMap(rootURL string) (string, error) {
	clean(&rootURL)
	fmt.Printf("Building sitemap for %s...\n\n", rootURL)

	unvisited := map[string]empty{rootURL: empty{}}
	found := map[string]empty{rootURL: empty{}}

	sitePages, err := crawlSite(rootURL, unvisited, found)
	if err != nil {
		return "", err
	}
	siteXML, err := buildXML(sitePages)
	if err != nil {
		return "", err
	}

	return siteXML, nil
}

func buildXML(sitePages []Page) (string, error) {

	doc := XMLMap{Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9", Pages: sitePages}
	var xmlData bytes.Buffer
	w := bufio.NewWriter(&xmlData)

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")

	if err := enc.Encode(doc); err != nil {
		return "", err
	}
	siteXML := xmlData.String()

	return siteXML, nil
}

// crawlSite crawls over pages in the domain, returning a slice of Pages.
func crawlSite(
	rootURL string,
	unvisited map[string]empty,
	found map[string]empty) ([]Page, error) {
	var SitePages []Page
	for len(unvisited) > 0 {
		for lnk := range unvisited {
			// Remove the link from the list of unvisiteds.
			delete(unvisited, lnk)

			// Get the page's links, create the page, and add it to SitePages.
			uid, err := pageURLsInDomain(lnk, rootURL)
			if err != nil {
				return nil, err
			}
			thisPage := Page{URL: lnk}
			SitePages = append(SitePages, thisPage)

			// Add unfound pages to found and unvisted lists.
			for _, dl := range uid {
				if _, ok := found[dl]; !ok {
					found[dl] = empty{}
					unvisited[dl] = empty{}
				}
			}
		}
	}
	return SitePages, nil
}

// pageURLsInDomain filters URLs from a given page for the ones in the domain.
func pageURLsInDomain(lnk string, rootURL string) ([]string, error) {
	urls, err := getPageLinks(lnk)
	if err != nil {
		return nil, err
	}
	uid, err := filterLinks(rootURL, urls)
	if err != nil {
		return nil, err
	}
	// Use keyList to return a slice of the URLs.
	return keyList(uid), nil
}

// filterLinks filters a map to contain only URLs in the domain of the root URL.
func filterLinks(rootURL string, urls []string) (map[string]empty, error) {
	pr, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	rootHost := pr.Host

	uid := make(map[string]empty)
	for _, u := range urls {
		pu, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		var fileExt string
		if len(u) > 4 {
			fileExt = u[len(u)-4:]
		}
		// Check if the link is in the domain of the root URL.
		if (pu.Host == rootHost || pu.Host == "") && (fileExt != ".jpg" && fileExt != ".png") {
			clean(&pu.Path)
			if _, ok := uid[rootURL+pu.Path]; !ok {
				uid[rootURL+pu.Path] = empty{}
			}
		}
	}
	return uid, nil
}

// getPageLinks returns a slice of links from a page.
func getPageLinks(url string) ([]string, error) {
	s, err := pageToString(url)
	if err != nil {
		return nil, err
	}
	r := strings.NewReader(string(s))
	l, err := ParseLinks(r)
	if err != nil {
		return nil, err
	}
	urls := GetURLs(l)

	return urls, nil
}

// pageToString takes in a url and returns a string of the HTML.
func pageToString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	str := fmt.Sprintf("%s\n", html)

	return str, nil
}

// keyList is a utility that takes a map of links and creates a list of its keys.
func keyList(mp map[string]empty) []string {
	var ret []string
	for key := range mp {
		ret = append(ret, key)
	}
	return ret
}

// clean is a utility to remove a trailing slash on a string.
func clean(s *string) {
	*s = strings.TrimRight(*s, "/")
}

// HACK: adding  contents here because I'm not yet sure how to use
// go 1.13 with private and local repos.

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
