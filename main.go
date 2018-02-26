package main

import (
	"fmt"
	"os"

	tui "github.com/awoitte/input_output_tui"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("usage: irc_tui server[:port] name")
		return
	}
	chat_messages := make(chan string)
	quit := make(chan bool)
	user_input := make(chan string)

	irc, err := connect_to_irc(args[0], args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	irc.attach_listeners(user_input, chat_messages, quit)

	tui.Start(chat_messages, user_input, quit)
}
