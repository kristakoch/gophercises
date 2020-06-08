package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gophercises/quiet_hn/hn"
	"github.com/patrickmn/go-cache"
)

func main() {
	// Parse flags.
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	// Make the template.
	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	// Create a new cache.
	c := cache.New(5*time.Minute, 10*time.Minute)

	http.HandleFunc("/", handler(numStories, tpl, c))

	// Start the server.
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(
	numStories int,
	tpl *template.Template,
	c *cache.Cache) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var stories []item
		var err error
		st, found := c.Get("stories")
		start := time.Now()
		if found {
			stBytes := st.([]byte)
			err = json.Unmarshal(stBytes, &stories)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			stories, err = getTopStories(numStories)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			storyBytes, err := json.Marshal(stories)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			c.Set("stories", storyBytes, cache.DefaultExpiration)

		}

		// Create the page using the story data.
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getTopStories(numStories int) ([]item, error) {
	// Get the IDs of the top stories from the HN API.
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}

	// Get each story and sort the resulting slice.
	stories := getStoryData(10, numStories, ids)
	sort.Slice(stories, func(i, j int) bool {
		return stories[i].Rank < stories[j].Rank
	})

	return stories, nil
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type item struct {
	hn.Item
	Host string
	Rank int
}

// getStoryData creates a worker pool to distribute
// getting story data from the HN API.
func getStoryData(numWorkers int, numStories int, ids []int) []item {
	var stories []item
	jobs := make(chan map[string]int)
	results := make(chan item)

	// Create the workers.
	for w := 0; w <= numWorkers; w++ {
		go work(w, jobs, results)
	}

	// Send the jobs into the worker pool.
	go func() {
		for i, sid := range ids {
			jobs <- map[string]int{"sid": sid, "rank": i}
		}
		close(jobs)
	}()

	// Collect the story data.
	for a := 0; a < numStories; a++ {
		res := <-results
		stories = append(stories, res)
	}

	return stories
}

// Work initializes a worker that listens for jobs
// on the jobs channel.
func work(wid int, jobs chan map[string]int, results chan item) {
	// fmt.Printf("worker %v initialized\n", wid)
	var client hn.Client
	for s := range jobs {
		// fmt.Printf("worker %v processing story id %v...\n", wid, id)
		// Get the story data.
		id := s["sid"]
		hnItem, err := client.GetItem(id)
		if err != nil {
			return
		}
		item := parseHNItem(hnItem)

		// Set the rank.
		r := s["rank"]
		item.Rank = r

		// Send the data.
		if isStoryLink(item) {
			results <- item
		}
	}
	// fmt.Printf("worker %v exiting\n", wid)
	return
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}
