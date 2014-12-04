package main

const (
	// FIXME
	staticDir = "/Users/bseitz/go/src/github.com/brnstz/sandbox/listen/html"
)

/*
func getListings(listReq chan bool, resp chan []byte) {

}

func getShows(w http.ResponseWriter, r *http.Request) {
	// Get
	bucket := getBucket()
	lpath := path.Join(rootPath, listingPath)

}

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
*/
