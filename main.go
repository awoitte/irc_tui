package main

import (
	"fmt"
	"os"
	"time"
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

	irc.attach_listeners(chat_messages, user_input, quit)
	start_tui(chat_messages, quit, user_input)

	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
