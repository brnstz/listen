package main

import (
	"log"
	"os"
)

var debug bool

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if os.Getenv("LISTEN_DEBUG") == "1" {
		debug = true
	}

	for i := 0; i < numWorkers; i++ {
		var s show
		var b band
		var v venue
		go recvAndWrite(&s)
		go recvAndWrite(&b)
		go recvAndWrite(&v)
	}

	go pullFromOMR()

	forever := make(chan bool)
	<-forever
}
