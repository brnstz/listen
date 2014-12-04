package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/brnstz/ohmy"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

const (
	exchangeName = "listen"
	bucketName   = "brnstz"

	routeShow  = "show"
	routeVenue = "venue"
	routeBand  = "band"

	rootPath = "/listen"

	listingPath = "listings.json"

	numShows = 100

	datePath = "/2006/01/02"
)

// A listing is aggregated data of a show with all details
type listing struct {
	Starts *time.Time `json:"starts"`
	Venue  venue      `json:"venue"`
	Bands  []band     `json:"bands"`
}

// A show with only references to venues/bands, not full details
type show struct {
	Starts *time.Time `json:"starts"`
	Venue  string     `json:"venue"`
	Bands  []string   `json:"bands"`
}

func (s *show) QueueName() string {
	return "show"
}

func (s *show) RouteName() string {
	return routeShow
}

func (s *show) Process(b []byte) (err error) {
	var incoming ohmy.Show

	err = decode(b, &incoming)
	if err != nil {
		return
	}

	s.Starts = incoming.Starts
	s.Venue = incoming.Venue.Slug

	temp := make([]string, len(incoming.Bands))

	for i, band := range incoming.Bands {
		temp[i] = band.Slug
	}

	s.Bands = temp

	return
}

func (s *show) Path() string {
	return path.Join(
		rootPath, "/show", s.Starts.Format(datePath),
		fmt.Sprint(s.Starts.Unix()),
		fmt.Sprintf("%d-%s.json", s.Starts.Unix(), s.Venue),
	)
}

// The band that is playing with track listings
type band struct {
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Tracks      []track    `json:"tracks"`
	LastUpdated *time.Time `json:"last_updated"`
}

func (b *band) QueueName() string {
	return "band"
}

func (b *band) RouteName() string {
	return routeBand
}

func (b *band) Process(by []byte) (err error) {
	var incoming ohmy.Band

	err = decode(by, &incoming)
	if err != nil {
		return
	}

	b.Name = incoming.Name
	b.Slug = incoming.Slug

	// FIXME: get tracks from spotify, etc.
	//b.Tracks =
	//b.LastUpdated

	return
}

func (b *band) Path() string {
	return path.Join(
		rootPath, "/band", fmt.Sprintf("%s.json", b.Slug),
	)
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

func (v *venue) QueueName() string {
	return "venue"
}

func (v *venue) RouteName() string {
	return routeVenue
}

func (v *venue) Process(b []byte) (err error) {
	var incoming ohmy.Venue

	err = decode(b, &incoming)
	if err != nil {
		return
	}

	now := time.Now()

	v.Address = incoming.Address
	v.Latitude = incoming.Latitude
	v.Longitude = incoming.Longitude
	v.Name = incoming.Name
	v.Slug = incoming.Slug
	v.LastUpdated = &now

	return
}

func (v *venue) Path() string {
	return path.Join(
		rootPath, "/venue", fmt.Sprintf("%s.json", v.Slug),
	)
}

type track struct {
	// Source of the track to hint the display (e.g., spotify)
	Source string `json:"source"`

	// The HTML to display to link to the track. Usually embedded player.
	HTML string `json:"html"`
}

func getBucket() *s3.Bucket {
	s3auth := aws.Auth{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
	s3conn := s3.New(s3auth, aws.Regions[os.Getenv("AWS_DEFAULT_REGION")])

	return s3conn.Bucket(bucketName)
}
