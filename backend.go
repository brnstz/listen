package main

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	qDial = "amqp://guest:guest@192.168.59.103:5672"
)

func setup() *amqp.Channel {
	conn, err := amqp.Dial(qDial)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic(err)
	}

	return ch
}

func send() {
	ch := setup()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello"),
	}

	err := ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		msg,    //
	)

	if err != nil {
		panic(err)
	}

	log.Println("sent it")
}

func receive() {
	ch := setup()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when used
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, //queue name
		"",     // consumer
		true,   // auto-ack FIXME
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received %s\n", d.Body)
			log.Printf("%+v\n", d)
		}
	}()

}

func main() {
	receive()
	send()
	log.Println("after send")

	forever := make(chan bool)
	<-forever
}
