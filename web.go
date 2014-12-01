package main

import (
	"io"
	"log"
	"net/http"
)

const (
	// FIXME
	staticDir = "/Users/bseitz/go/src/github.com/brnstz/sandbox/listen/html"
)

func testHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("hello")
	io.WriteString(w, "hello there")
}

func main() {
	// Test API stuff
	http.HandleFunc("/api/test.json", testHandle)

	// By default look for a static asset
	http.Handle("/", http.FileServer(http.Dir(staticDir)))

	// FIXME
	err := http.ListenAndServe(":8003", nil)
	if err != nil {
		panic(err)
	}
}
