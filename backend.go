package main

import (
	"log"
	"net/http"
	"os"
)

var debug bool

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if os.Getenv("LISTEN_DEBUG") == "1" {
		debug = true
	}

	for i := 0; i < numWorkers; i++ {
		var s show
		var b band
		var v venue
		go recvAndWrite(&s)
		go recvAndWrite(&b)
		go recvAndWrite(&v)
	}

	go pullFromOMR()
	go agg()

	http.HandleFunc("/api/shows.json", func(w http.ResponseWriter, r *http.Request) {

		_, err := w.Write(cachedListings)
		if err != nil {
			log.Println("Error serving listings", err)
		}
	})

	// By default look for a static asset
	http.Handle("/", http.FileServer(http.Dir(os.Getenv("LISTEN_STATIC_DIR"))))

	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal(err)
	}
}
