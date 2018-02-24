package main

import "time"

func main() {

	messages := make(chan string)
	quit := make(chan bool)
	input := make(chan string)
	go connect_to_irc(messages, input)
	go show_messages(messages, quit, input)

	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
