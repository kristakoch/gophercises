package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()

	var dbg bool
	flag.BoolVar(&dbg, "debug", false, "toggle debug mode")
	flag.Parse()

	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)

	fmt.Println("running at http://localhost:4444")
	log.Fatal(http.ListenAndServe(":4444", recoverHandler(mux, dbg)))
}

func recoverHandler(next http.Handler, dbg bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				debug.PrintStack()

				if dbg == false {
					fmt.Fprintf(w, "<p>Something went wrong</p>")
					return
				}
				fmt.Fprintf(w, "<h1>Panic!</h1><pre>"+string(debug.Stack())+"</pre>")

			}
		}()

		rw := responseWriter{ResponseWriter: w}

		next.ServeHTTP(&rw, r)

		rw.flush()
	})
}

type responseWriter struct {
	http.ResponseWriter
	writes     [][]byte
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *responseWriter) flush() {
	if rw.statusCode == 0 {
		rw.statusCode = 200
	}
	rw.ResponseWriter.WriteHeader(rw.statusCode)
	for _, write := range rw.writes {
		rw.ResponseWriter.Write(write)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
