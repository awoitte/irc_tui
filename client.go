package main

import (
	"fmt"

	irc "github.com/awoitte/irc"
)

func connect_to_irc(chat_messages, commands chan string) {

	tcp, err := irc.TCPConnect("irc.mozilla.org", "6667")
	if err != nil {
		fmt.Println("Error connecting via TCP: ", err)
		return
	}
	connection := irc.Connect("bamboo", tcp)
	go get_irc_messages(&connection, chat_messages)
	go dispatch_commands(&connection, commands, chat_messages)
}

func get_irc_messages(connection *irc.IRC, messages chan string) {
	connection.ReadLoop(func(message string) error {
		messages <- message
		return nil
	})
	close(messages)
}

func dispatch_commands(connection *irc.IRC, commands, chat_messages chan string) {
	for {
		command := <-commands
		translated := translate_command(command)
		chat_messages <- translated
		connection.SendRaw(translated)
	}
}
