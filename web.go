package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/brnstz/ohmy"
)

const (
	// FIXME
	staticDir = "/Users/bseitz/go/src/github.com/brnstz/sandbox/listen/html"
	numShows  = 100
)

func getShows(w http.ResponseWriter, r *http.Request) {
	shows, err := ohmy.GetShows(ohmy.RegionNYC, numShows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(shows)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(j)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

/*
func main() {
	// Test API stuff
	http.HandleFunc("/api/shows.json", getShows)

	// By default look for a static asset
	http.Handle("/", http.FileServer(http.Dir(staticDir)))

	// FIXME
	err := http.ListenAndServe(":8003", nil)
	if err != nil {
		panic(err)
	}
}*/
