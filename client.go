package main

import (
	"errors"
	"fmt"
	"strings"

	irc "github.com/awoitte/irc_client"
)

const read_chunk_size = 128

type IRC_client struct {
	connection IRC_Connection
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
		return nil, fmt.Errorf("Error connecting via TCP: %v", err)
	}
	connection := irc.Connect(name, tcp)
	return &IRC_client{&connection}, nil
}

func (client *IRC_client) attach_listeners(chat_messages, user_input chan string, quit chan bool) {
	go get_irc_messages(client.connection, chat_messages)
	go dispatch_commands(client.connection, user_input, chat_messages, quit)
}

func get_irc_messages(connection IRC_Connection, messages chan string) {
	connection.ReadLoop(read_chunk_size, func(message string) error {
		cleaned_message := strings.Replace(message, "\r", "", -1)
		lines := strings.Split(cleaned_message, "\n")

		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				messages <- line
			}
		}

		return nil
	})
}
