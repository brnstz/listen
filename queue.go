package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"

	"github.com/streadway/amqp"
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
func receiveFromQueue(ch *amqp.Channel, name, route string) (msgs <-chan amqp.Delivery, err error) {

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
		route,
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

func publishAsGob(value interface{}, ch *amqp.Channel, route string) (err error) {

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

	log.Printf("Publishing %+v to %v %v\n", value, exchangeName, route)
	// Publish our gob
	err = ch.Publish(
		// Send to our exchange
		exchangeName,
		route,

		// not mandatory, not immediate
		false, false,

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
