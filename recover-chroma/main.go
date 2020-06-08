package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()

	var dbg bool
	flag.BoolVar(&dbg, "debug", false, "toggle debug mode")
	flag.Parse()

	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/source/", sourceFileHandler)
	mux.HandleFunc("/", hello)

	fmt.Println("running at http://localhost:4444")
	log.Fatal(http.ListenAndServe(":4444", recoverHandler(mux, dbg)))
}

// panicDemo panics immediately.
func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

// panicAfterDemo writes to the page then panics.
func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

// funcThatPanics panics to enable panic handling.
func funcThatPanics() {
	panic("Oh no!")
}

// hello prints a greeting.
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

// recoverHandler is a wrapper handler that recovers from panics.
func recoverHandler(next http.Handler, dbg bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Recover and print a message to the screen.
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				debug.PrintStack()

				if dbg == false {
					fmt.Fprintf(w, "<p>Something went wrong</p>")
					return
				}
				panicHTML := parsePanic(string(debug.Stack()))
				fmt.Fprint(w, panicHTML)
			}
		}()

		// Buffer and then send the response.
		rw := responseWriter{ResponseWriter: w}
		next.ServeHTTP(&rw, r)
		rw.flush()
	})
}

// sourceFileHandler renders highlighted go files to the browser.
func sourceFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file contents.
	path := r.FormValue("path")
	baseFileURL, err := url.QueryUnescape(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	fb, err := ioutil.ReadFile(baseFileURL)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// Get the line number.
	lNum := r.FormValue("line")
	lNumInt, err := strconv.Atoi(lNum)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// Format the file.
	buf, err := formatGoSource(fb, lNumInt)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	fmt.Fprintf(w, buf.String())
}

// formatGoSource uses chroma to add source code highlighting
// to go files being handled by the server.
func formatGoSource(fb []byte, lNumInt int) (bytes.Buffer, error) {
	var buf bytes.Buffer

	// All files will be go files.
	lexer := lexers.Get("go")

	// Build a custom line height style alongside paraiso-dark.
	builder := styles.Get("paraiso-dark").Builder()
	builder.Add(chroma.LineHighlight, "#433442 bg:#433442")
	style, err := builder.Build()
	if err != nil {
		return buf, err
	}

	// Create the formatter and iterator.
	formatter := html.New(
		html.Standalone(true),
		html.WithLineNumbers(true),
		html.HighlightLines([][2]int{{lNumInt, lNumInt}}))

	iterator, err := lexer.Tokenise(nil, string(fb))
	if err != nil {
		return buf, err
	}

	// Format the file.
	if err = formatter.Format(&buf, style, iterator); err != nil {
		return buf, err
	}

	return buf, nil
}

// parsePanic returns the panic stack trace after finding
// links to files and wrapping them in anchor tags.
func parsePanic(s string) string {
	lines := strings.Split(s, "\n")

	// Rebuild the panic stack trace line by line.
	for i, ln := range lines {
		re := regexp.MustCompile(`\/(.)*.go`)
		found := re.Find([]byte(ln))

		// Only lines with go files will be changed.
		if found == nil {
			continue
		}
		pathURL := url.QueryEscape(string(found))

		lnInfo := strings.Split(ln, ":")[1]
		lnNum := strings.Fields(lnInfo)[0]

		// Add the URL value for the line number.
		v := url.Values{}
		v.Add("path", pathURL)
		v.Add("line", lnNum)
		srcURL := "/source" + "?" + v.Encode()

		// Build the enw line and replace the old one..
		newLine := fmt.Sprintf("\t<a href='%s'>%s</a>:%s", srcURL, found, lnInfo)
		lines[i] = newLine
	}
	parsedStr := strings.Join(lines, "\n")

	return "<h1>Panic!</h1><pre>" + parsedStr + "</pre>"
}

// responseWriter implements ResponseWriter to
// buffer http responses until a final write.
type responseWriter struct {
	http.ResponseWriter
	writes     [][]byte
	statusCode int
}

// Write buffers writes to the responseWriter.
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

// WriteHeader buffers status code assignments
// to the responseWriter
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

// flush writes the buffered data to the ResponseWriter.
func (rw *responseWriter) flush() {
	if rw.statusCode == 0 {
		rw.statusCode = 200
	}
	rw.ResponseWriter.WriteHeader(rw.statusCode)
	for _, write := range rw.writes {
		rw.ResponseWriter.Write(write)
	}
}
