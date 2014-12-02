package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"time"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"

	"github.com/brnstz/ohmy"
	"github.com/streadway/amqp"
)

const (
	exchangeName = "listen"
	bucket       = "brnstz"
	path         = "/listen/shows"
)

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

		// FIXME: what are these?
		false,
		nil,
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

	// Success!
	return
}

func oneReader() {
	// Get connection to rabbit
	// FIXME: does it make sense to reconnect here?
	url := os.Getenv("AMQP_URL")
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Get a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()

	err = ensureExchange(ch)
	if err != nil {
		return
	}

	// Get current shows
	shows, err := ohmy.GetShows(ohmy.RegionNYC, 100)
	if err != nil {
		// Skip when there's an err, will be logged by the ohmy lib
		return
	}

	// Look at each show
	for _, show := range shows {
		log.Printf("%+v\n", show)
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

func s3Writer() {
	// Get connection to rabbit
	// FIXME: does it make sense to reconnect here?
	url := os.Getenv("AMQP_URL")
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Get a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()

	err = ensureExchange(ch)
	if err != nil {
		return
	}

	s3auth := aws.Auth{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
	s3conn := s3.New(s3auth, aws.Regions[os.Getenv("AWS_DEFAULT_REGION")])
	s3Bucket := s3conn.Bucket(bucket)
	log.Println(s3Bucket)

	msgs, err := receiveFromQueue(ch, "s3")
	for d := range msgs {
		log.Printf("Received %s\n", d.Body)
		log.Printf("%+v\n", d)
	}

}

func main() {
	go s3Writer()
	go reader()

	forever := make(chan bool)
	<-forever
}
