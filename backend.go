package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"

	"github.com/brnstz/ohmy"
	"github.com/streadway/amqp"
)

const (
	exchangeName       = "listen"
	bucketName         = "brnstz"
	rootPath           = "/listen"
	showPath           = "/shows"
	showPathTimeFormat = "2006-01-02_15"
)

// Try to connect, returning either both channel and connection, or an error.
// Also ensure the exchange exists while we're at it. Caller is responsible
// for closing the connection and channel.
func connect() (ch *amqp.Channel, conn *amqp.Connection, err error) {

	// Get connection to rabbit
	url := os.Getenv("AMQP_URL")
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Println(err)
		return
	}

	// Get a channel
	ch, err = conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}

	// Make sure the exchange exists
	err = ensureExchange(ch)
	if err != nil {
		return
	}

	// Success
	return
}

// Declare our exchange
func ensureExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		// Fanout queue called listen
		exchangeName, "fanout",

		// yes durable
		true,

		// no auto-delete, no internal, no noWait
		false, false, false,

		// no extra arguments
		nil,
	)
}

// Bind to an exchange and get a Go channel of messages
func receiveFromQueue(ch *amqp.Channel, name string) (msgs <-chan amqp.Delivery, err error) {

	// Create the queue
	q, err := ch.QueueDeclare(
		name,

		// false to durable, delete when used
		false, false,

		// true to exclusive
		true,

		// no no wait and nil arguments
		false, nil,
	)
	if err != nil {
		log.Println(err)
		return
	}

	// Bind queue to our exchange
	err = ch.QueueBind(
		q.Name,

		// blank routing key
		"",

		exchangeName,

		// no no wait and nil extra args
		false, nil,
	)
	if err != nil {
		log.Println(err)
		return
	}

	// Start receiving messages
	msgs, err = ch.Consume(
		q.Name,

		// consumer
		"",

		// auto-ack
		true,

		//  false to exclusive, no local, no wait
		false, false, false,

		nil, // args
	)
	if err != nil {
		log.Println(err)
	}

	return
}

func publishAsGob(value interface{}, ch *amqp.Channel) (err error) {

	// Setup stuff for gob
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	// Encode into a gob
	err = enc.Encode(value)
	if err != nil {
		log.Println("Error encoding value %#v %v", value, err)
		return
	}

	// Create a message with the gob
	msg := amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        buff.Bytes(),
	}

	// Publish our gob
	err = ch.Publish(
		// Send to our exchange
		exchangeName,

		// routing key, not mandatory, not immediate
		"", false, false,

		msg,
	)
	if err != nil {
		log.Printf("Error publishing message: %+v %v\n", msg, err)
		return
	}

	// Success
	return
}

// Given data in src, decode it as a gob into dst
func decode(src []byte, dst interface{}) (err error) {
	buf := bytes.NewBuffer(src)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(dst)

	return
}

func oneReader() {
	ch, conn, err := connect()
	if err != nil {
		log.Fatal("Cannot connect", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Get current shows
	shows, err := ohmy.GetShows(ohmy.RegionNYC, 100)
	if err != nil {
		// Skip when there's an err, will be logged by the ohmy lib
		return
	}

	// Look at each show
	for _, show := range shows {
		publishAsGob(show, ch)
	}
}

// Periodically pull updates from ohmy and publish it to our exchange
func reader() {
	for {
		oneReader()
		time.Sleep(time.Minute * 5)
	}
}

/*
func showPath(s *ohmy.Show) {
	fullPath := path.Join(rootPath, showPath)
}
*/

func s3Writer() {
	ch, conn, err := connect()
	if err != nil {
		log.Fatal("Cannot connect", err)
	}
	defer conn.Close()
	defer ch.Close()

	s3auth := aws.Auth{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
	s3conn := s3.New(s3auth, aws.Regions[os.Getenv("AWS_DEFAULT_REGION")])
	bucket := s3conn.Bucket(bucketName)

	msgs, err := receiveFromQueue(ch, "s3")
	for d := range msgs {

		// Decode the rabbit message into a Show object
		var show ohmy.Show
		err = decode(d.Body, &show)
		if err != nil {
			continue
		}

		// Create the path to store this show under
		fullPath := path.Join(
			rootPath, showPath,
			show.Starts.Format(showPathTimeFormat),
			fmt.Sprint(show.Venue.Slug, ".json"),
		)

		b, err := json.Marshal(show)
		if err != nil {
			log.Println("Cannot encode show", err)
			continue
		}

		err = bucket.Put(fullPath, b, "application/json", s3.Private)
		if err != nil {
			log.Println("Cannot store show", err)
			continue
		}

	}

}

func main() {
	go s3Writer()
	go reader()

	forever := make(chan bool)
	<-forever
}
