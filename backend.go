package main

func main() {

	var s show
	var b band
	var v venue
	go recvAndWrite(&s)
	go recvAndWrite(&b)
	go recvAndWrite(&v)

	go pullFromOMR()

	forever := make(chan bool)
	<-forever
}
