package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler ...
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	mapper := func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path // aliased path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, req, dest, http.StatusFound) // num vs statusfound
			return                                        // make sure to return so fallback doesn't run
		}
		fallback.ServeHTTP(w, req)
	}
	return mapper
}

// YAMLHandler ...
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	rs, err := parseYAML(yml)
	if err != nil {
		log.Fatal(err)
	}

	rsMap := buildMap(rs)

	return MapHandler(rsMap, fallback), nil
}

// JSONHandler ...
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	rs, err := parseJSON(json)
	if err != nil {
		log.Fatal(err)
	}

	rsMap := buildMapFromJSON(rs)

	return MapHandler(rsMap, fallback), nil
}

// type for parsed yaml data
type yamlRedirects []struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// parseYAML returns a yamlRedirects type
func parseYAML(yml []byte) (yamlRedirects, error) {
	var rs yamlRedirects
	err := yaml.Unmarshal(yml, &rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// buildMap converts a yamlRedirects type into a strings map
func buildMap(rs yamlRedirects) map[string]string {
	rsMap := make(map[string]string)
	for _, r := range rs {
		rsMap[r.Path] = r.URL
	}
	return rsMap
}

// type for parsed json data
type jsonRedirects []struct {
	Path string `json:"Path"`
	URL  string `json:"URL"`
}

// parseJSON returns a redirects type
func parseJSON(js []byte) (jsonRedirects, error) {
	var rs jsonRedirects
	err := json.Unmarshal(js, &rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// buildMap converts a redirects type into a strings map
func buildMapFromJSON(rs jsonRedirects) map[string]string {
	rsMap := make(map[string]string)
	for _, r := range rs {
		rsMap[r.Path] = r.URL
	}
	return rsMap
}
