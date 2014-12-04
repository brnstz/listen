package main

func main() {

	var s show
	go recvAndWrite(&s)
	go pullFromOMR()

	forever := make(chan bool)
	<-forever
}
