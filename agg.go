package main

import (
	"encoding/json"
	"log"
	"path"
	"time"

	"launchpad.net/goamz/s3"
)

const maxListingDays = 14

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

	// Get the venue object via slug
	var v venue
	err = bucketJSON(entityPath("/show", s.Venue), &v)
	if err != nil {
		return
	}

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
}

// Periodically aggregage current shows list for front page
func agg() {
	for {

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

				// List the bucket dir for this day
				resp, err := bucket.List(showDir[1:], "/", "", 1000)
				if err != nil {
					log.Println("Can't list dir", err)
					break
				}

				// For each show
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
				now.AddDate(0, 0, 1)
				days++
			}

			b, err := json.Marshal(listings)
			if err != nil {
				log.Println(err)
				return
			}

			err = bucket.Put(path.Join(rootPath, listingPath), b, "application/json", s3.Private)
			if err != nil {
				log.Println("Cannot store aggregate file", err)
				continue
			} else {
				log.Println("Successfully wrote aggregate to", e.Path())
			}
		}()

		time.Sleep(time.Minute * 5)
	}
}
