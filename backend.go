package main

func main() {

	var s show
	var b band
	go recvAndWrite(&s)
	go recvAndWrite(&b)
	go pullFromOMR()

	forever := make(chan bool)
	<-forever
}
