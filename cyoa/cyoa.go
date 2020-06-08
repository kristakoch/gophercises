// Package cyoa creates a command line or HTML choose your own
// adventure story from a JSON file.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// ServeCmdLnStory gets the story data and then
// serves the story through the command line.
func ServeCmdLnStory() error {
	s, err := buildStoryMap()
	if err != nil {
		return err
	}
	arc := "intro"
	for {
		fmt.Printf("\n--------------------------\n%s\n--------------------------\n", s[arc].Title)
		for _, p := range s[arc].Story {
			fmt.Printf("\n%s\n\n", p)
		}
		numOptions := len(s[arc].Options)
		var choice int
		if numOptions < 1 {
			break
		}
		fmt.Println("Here are your options...")
		for i, o := range s[arc].Options {
			fmt.Println(o.TextOption)
			fmt.Printf("Press %v to venture.\n\n", i)
		}
		fmt.Scanln(&choice)
		for choice < 0 || choice > (numOptions-1) { // if alpha chars are entered, defaults to 0
			fmt.Printf("Must choose a number between 0 and %v.", numOptions-1)
			fmt.Scanln(&choice)
		}
		arc = s[arc].Options[choice].ArcOption
	}
	return nil
}

// ServeHTMLStory gets the story data, creats the template,
// creates the pages, and then serves the story on localhost.
func ServeHTMLStory() error {
	s, err := buildStoryMap()
	if err != nil {
		return err
	}

	tmpl, err := template.New("gopher-story-template").Parse(storyTemplate)
	if err != nil {
		return err
	}

	handler := NewStoryHandler(s, tmpl)

	fmt.Println("Story is running at http://localhost:1313 ...")
	http.ListenAndServe(":1313", handler)

	return nil
}

// StoryHandler is type for handling story requests
type StoryHandler struct {
	s map[string]arc
	t *template.Template
}

// NewStoryHandler creates a story handler
func NewStoryHandler(story map[string]arc, tmpl *template.Template) StoryHandler {
	return StoryHandler{story, tmpl}
}

// implements Handler interface through adding
// ServeHTTP method to StoryHandler
func (s StoryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if path == "" {
		path = "intro"
	}
	err := s.t.Execute(w, s.s[path])
	if err != nil {
		log.Fatal(err)
	}
}

// buildStoryMap reads from the json file and creats a map
// of the story arcs
func buildStoryMap() (map[string]arc, error) {
	f, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		return nil, err
	}
	fileStr := string(f)

	var s map[string]arc
	err = json.Unmarshal([]byte(fileStr), &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// arc is the type of object unmarshalled from json data
type arc struct {
	Title   string     `json:"title"`
	Story   []string   `json:"story"`
	Options []struct { // may be more readable to create Option type instead of leaving it as an undefined struct
		TextOption string `json:"text"`
		ArcOption  string `json:"arc"`
	} `json:"options"`
}

// storyTemplate is the template used for HTML story arc files
const storyTemplate = `
<div style="max-width: 800px; margin:auto; padding: 45px 3%; font-family: Courier New">
	<a href="/" style="font-size: 40px; text-decoration: none">
		🕳
	</a>
	<h1>{{.Title}}</h1>
	<div>
		{{range .Story}}
			<p>{{.}}</p>
		{{end}}
	</div>
	{{range .Options}}
		<div style="padding: 20px 0px">
			<a href="/{{.ArcOption}}">{{.TextOption}}</a>
		</div>
	{{end}}
</div>
`
