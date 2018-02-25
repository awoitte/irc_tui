package main

import "time"

func main() {

	messages := make(chan string)
	quit := make(chan bool)
	commands := make(chan string)
	go connect_to_irc(messages, commands)
	go show_messages(messages, quit, commands)

	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
