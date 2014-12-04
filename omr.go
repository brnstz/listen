package main

import (
	"log"
	"time"

	"github.com/brnstz/ohmy"
)

// Periodically pull updates from ohmy and publish it to our exchange
func pullFromOMR() {
	for {

		func() {
			ch, conn, err := connect()
			if err != nil {
				log.Fatal("Cannot connect", err)
			}
			defer conn.Close()
			defer ch.Close()

			// Get current shows
			shows, err := ohmy.GetShows(ohmy.RegionNYC, numShows)
			if err != nil {
				// Skip when there's an err, will be logged by the ohmy lib
				return
			}

			// Iterate over shows to find individual items to process, i.e.,
			// shows and bands and venues
			for _, show := range shows {
				publishAsGob(show, ch, routeShow)

				publishAsGob(show.Venue, ch, routeVenue)

				for _, band := range show.Bands {
					publishAsGob(band, ch, routeBand)
				}
			}

		}()

		time.Sleep(time.Minute * 60)
	}
}
