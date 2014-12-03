package main

import "time"

// The current list of shows
type shows struct {
	Starts *time.Time
	Venue  venue
	Bands  []band
}

// The band that is playing with track listings
type band struct {
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Tracks      []track    `json:"tracks"`
	LastUpdated *time.Time `json:"last_updated"`
}

// The venue where the show is happening
type venue struct {
	Address     string     `json:"address"`
	Latitude    string     `json:"string"`
	Longitude   string     `json:"string"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	LastUpdated *time.Time `json:"last_updated"`
}

type track struct {
	// Source of the track to hint the display (e.g., spotify)
	Source string `json:"source"`

	// The HTML to display to link to the track. Usually embedded player.
	HTML string `json:"html"`
}
