package main

import (
	"errors"
	"fmt"
	"strings"

	irc "github.com/awoitte/irc_client"
)

const read_chunk_size = 128

type IRC_client struct {
	connection *irc.IRC
}

func connect_to_irc(server, name string) (*IRC_client, error) {
	connection_details := strings.Split(server, ":")
	if len(connection_details) == 1 {
		connection_details = append(connection_details, "6667")
	}

	if len(connection_details) != 2 {
		return nil, errors.New("malformed server address")
	}

	tcp, err := irc.TCPConnect(connection_details[0], connection_details[1])
	if err != nil {
		return nil, fmt.Errorf("Error connecting via TCP: ", err)
	}
	connection := irc.Connect(name, tcp)
	return &IRC_client{&connection}, nil
}

func (client *IRC_client) attach_listeners(chat_messages, user_input chan string, quit chan bool) {
	go get_irc_messages(client.connection, chat_messages)
	go dispatch_commands(client.connection, user_input, chat_messages, quit)
}

func get_irc_messages(connection *irc.IRC, messages chan string) {
	connection.ReadLoop(read_chunk_size, func(message string) error {
		messages <- message
		if message == "QUIT" {
			close(messages)
			return errors.New("IRC Session quit")
		}

		return nil
	})
}
