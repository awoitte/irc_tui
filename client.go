package main

import (
	"fmt"

	irc "github.com/awoitte/irc"
)

func connect_to_irc(messages chan string, input chan string) {

	tcp, err := irc.TCPConnect("irc.mozilla.org", "6667")
	if err != nil {
		fmt.Println("Error connecting via TCP: ", err)
		return
	}
	connection := irc.Connect("bamboo", tcp)
	go get_irc_messages(&connection, messages)
	go dispatch_commands(&connection, input)

}

func get_irc_messages(connection *irc.IRC, messages chan string) {
	connection.ReadLoop(func(message string) error {
		//fmt.Println(message)
		messages <- message
		return nil
	})
	close(messages)
}

func dispatch_commands(connection *irc.IRC, input chan string) {
	for {
		command := <-input
		connection.SendRaw(command)
	}
}
