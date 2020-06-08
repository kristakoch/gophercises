package main

// todo: take some time to think about
// the possible options people can enter
// don't let it break based on that.
// return useful error messages or let it run
// ex --- used to break when entering "db" instead
// of "usedb" for database name. should always either
// return an error or use a default

// think about: data omissions, trash data, bad combos
// flags are dataType, filePath, and dbName

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dataType := flag.String("dt", "yaml", "enter yaml, json, or usedb for data type")
	filePath := flag.String("fp", "none", "location of json or yaml file with redirect rules")
	dbName := flag.String("dn", "test", "enter name of database, table is assumed to be url_short, columns path and url")
	flag.Parse()

	// set the data
	var fData string
	var err error
	switch *dataType {
	case "yaml":
		fData, err = setYAML(*filePath)
	case "json":
		fData, err = setJSON(*filePath)
	case "usedb":
		fData = *dbName
		err = nil
	}
	if err != nil {
		log.Fatal(err)
	}

	// build the main (map) handler
	mux := defaultMux()
	pathsToUrls := getDefaultURLs()
	mapHandler := MapHandler(pathsToUrls, mux)

	// build the handler based on flags
	// should be refactored so the system never breaks
	// for any possible flag input
	var handler http.HandlerFunc
	switch *dataType {
	case "yaml":
		handler, err = YAMLHandler([]byte(fData), mapHandler)
	case "json":
		handler, err = JSONHandler([]byte(fData), mapHandler)
	case "usedb":
		var dbData map[string]string
		dbData, err = getDBDataMap(fData)
		handler = MapHandler(dbData, mux)
	default:
		handler = MapHandler(pathsToUrls, mux)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

// default mux
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultMsg)
	return mux
}

// default landing page for url not in shortener data
func defaultMsg(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Oí! Tudo bem?")
}

// getFileContents returns a string of the file contents
func getFileContents(fp string) (string, error) {
	dat, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// getDefaultURLs returns a default string map to be consumed by the map handler
func getDefaultURLs() map[string]string {
	return map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
}

// setYAML gets the contents of a yaml file and returns the string
// or a default yaml string
func setYAML(fp string) (string, error) {
	var yaml string
	if fp != "none" {
		fileYaml, err := getFileContents(fp)
		if err != nil {
			return "", err
		}
		yaml = fileYaml
	} else {
		yaml = `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution`
	}
	return yaml, nil
}

// setJSON gets the contents of a json file and returns the string
// or a default json string
func setJSON(fp string) (string, error) {
	var json string
	if fp != "none" {
		fileJSON, err := getFileContents(fp)
		if err != nil {
			return "", err
		}
		json = fileJSON
	} else {
		json = `[{"Path":"/urlshort","URL":"https://github.com/gophercises/urlshort"},{"Path":"/urlshort-final","URL":"https://github.com/gophercises/urlshort/tree/solution"}]`
	}
	return json, nil
}

// getDBDataMap gets data from db with a given name and returns a map
func getDBDataMap(dbt string) (map[string]string, error) {
	dbInfo := fmt.Sprintf("root@tcp(127.0.0.1)/%s", dbt)
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT path, url FROM url_short`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbData := make(map[string]string)
	for rows.Next() {
		var path, url string
		if err := rows.Scan(&path, &url); err != nil {
			return nil, err
		}
		dbData[path] = url
	}
	return dbData, nil
}
