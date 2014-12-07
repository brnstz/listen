package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	spotifySearch    = "https://api.spotify.com/v1/search?"
	spotifyTopTracks = "https://api.spotify.com/v1/artists/%s/top-tracks?country=us"
	spotifyEmbed     = `<iframe src="https://embed.spotify.com/?uri=%v" frameborder="0" allowtransparency="true"></iframe>`
)

type spotifySearchResults struct {
	Artists struct {
		Items []struct {
			Id   string
			URI  string
			Name string
		}
	}
}

type spotifyTopTracksResult struct {
	Tracks []struct {
		Id   string
		URI  string
		Name string
	}
}

// Get JSON via HTTP. Pass in a pointer in v to unmarshal response into that
// var.
func getJSON(getURL string, v interface{}) (err error) {

	// Get the response
	resp, err := http.Get(getURL)
	if err != nil {
		log.Println("Can't GET URL", getURL, err)
		return
	}
	defer resp.Body.Close()

	// Read the body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read response from URL", getURL, err)
		return
	}

	// Marshal into v
	err = json.Unmarshal(b, v)
	if err != nil {
		log.Println("Can't read JSON from URL", getURL, err)
		return
	}

	return
}

func searchSpotify(band string) (tracks []track) {

	// Construct search URL
	v := url.Values{}
	v.Set("q", band)
	v.Set("type", "artist")
	searchURL := fmt.Sprint(spotifySearch, v.Encode())

	// Search for the artist
	var sr spotifySearchResults
	err := getJSON(searchURL, &sr)
	if err != nil {
		return
	}

	// Need at least one artist
	if len(sr.Artists.Items) < 1 {
		log.Println("No artists in response from", searchURL)
		return
	}

	// Save the top artist
	id := sr.Artists.Items[0].Id
	name := sr.Artists.Items[0].Name

	log.Printf("Found %v %v\n", id, name)

	// Construct tracks URL
	tracksURL := fmt.Sprintf(spotifyTopTracks, id)

	// Get tracks
	var tr spotifyTopTracksResult
	err = getJSON(tracksURL, &tr)
	if err != nil {
		return
	}

	tracks = make([]track, len(tr.Tracks))
	for i, spotifyTrack := range tr.Tracks {
		tracks[i] = track{
			Source: "spotify",
			HTML:   fmt.Sprintf(spotifyEmbed, spotifyTrack.URI),
		}
	}

	return tracks
}
