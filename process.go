package main

import (
	"encoding/json"
	"log"

	"launchpad.net/goamz/s3"
)

// An entity from queue and writes to s3 after some possible conversion
type entity interface {
	QueueName() string
	RouteName() string
	Path() string
	Process(b []byte) error
}

func recvAndWrite(e entity) {
	ch, conn, err := connect()
	if err != nil {
		log.Fatal("Cannot connect", err)
	}
	defer conn.Close()
	defer ch.Close()

	bucket := getBucket()

	msgs, err := receiveFromQueue(ch, e.QueueName(), e.RouteName())
	for d := range msgs {

		err = e.Process(d.Body)
		if err != nil {
			continue
		}

		b, err := json.Marshal(e)
		if err != nil {
			log.Println("Cannot encode entity", err)
			continue
		}

		if debug {
			log.Printf("Received from queue: %s %v %v\n", b, e.QueueName(), e.RouteName())
		}

		err = bucket.Put(e.Path(), b, "application/json", s3.Private)
		if err != nil {
			log.Println("Cannot store entity", err)
			continue
		}
	}
}
