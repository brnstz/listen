package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"time"

	"launchpad.net/goamz/s3"
)

// How many days to look for shows
const maxListingDays = 14

const listingDateFormat = "Mon Jan 02 2006 03PM"

// Saved listings in memory
var cachedListings []byte

// Given a bucket path, get JSON and unmarhal into object
func bucketJSON(path string, v interface{}) (err error) {
	bucket := getBucket()

	// Get json data for show
	b, err := bucket.Get(path)
	if err != nil {
		log.Println("Problem getting", path)
		return
	}

	// Unmarshal into v (should be a pointer)
	err = json.Unmarshal(b, v)
	if err != nil {
		log.Println("Problem reading", err)
		return
	}

	return
}

func showToListing(s show, bucket *s3.Bucket) (l listing, err error) {

	// Copy over the starting time
	l.Starts = s.Starts
	l.StartsFormatted = s.Starts.Format(listingDateFormat)

	// Get the venue object via slug
	var v venue
	err = bucketJSON(entityPath("/venue", s.Venue), &v)
	if err != nil {
		return
	}
	log.Printf("venue: %+v %v", v, showPath(s.Starts, s.Venue))

	// Copy over the full venue object
	l.Venue = v

	// Get each band via slug
	for _, slug := range s.Bands {
		var b band
		err = bucketJSON(entityPath("/band", slug), &b)
		if err != nil {
			return
		}

		l.Bands = append(l.Bands, b)
	}

	return
}

// Periodically aggregage current shows list for front page
func agg() {

	first := true
	for {

		// On first iteration, try to load cached listings from bucket
		if first {
			bucket := getBucket()
			b, err := bucket.Get(path.Join(rootPath, listingPath))
			if err == nil {
				cachedListings = b
			}
			first = false
		}

		func() {
			now := time.Now()
			bucket := getBucket()

			// Loop until passing max days
			days := 0

			// Save all the listings we find
			listings := []listing{}

			for days < maxListingDays {
				// Create the path for this day
				showDir := path.Join(
					rootPath, "/show", now.Format(datePath),
				)
				// Format for bucket.List, no preceeding slash, but yes to
				// ending slash
				showDir = fmt.Sprint(showDir[1:], "/")

				// List the bucket dir for this day
				resp, err := bucket.List(showDir, "/", "", 1000)
				log.Printf("%+v\n", resp)

				if err != nil {
					log.Println("Can't list dir", err)
					break
				}

				// For each show, try to add or skip it
				for _, key := range resp.Contents {

					var s show
					err = bucketJSON(key.Key, &s)
					if err != nil {
						continue
					}

					l, err := showToListing(s, bucket)
					if err != nil {
						continue
					}

					listings = append(listings, l)
				}

				// Set up next iteration
				now = now.AddDate(0, 0, 1)
				days++
			}

			b, err := json.Marshal(listings)
			if err != nil {
				log.Println(err)
				return
			}

			lpath := path.Join(rootPath, listingPath)
			err = bucket.Put(path.Join(rootPath, listingPath), b, "application/json", s3.Private)
			if err != nil {
				log.Println("Cannot store aggregate file", err)
				return
			} else {
				log.Println("Successfully wrote aggregate to", lpath, len(b))
				cachedListings = b
			}
		}()

		time.Sleep(time.Minute * 60)
	}
}
