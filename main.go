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
	messages := make(chan string)
	quit := make(chan bool)
	commands := make(chan string)

	irc, err := connect_to_irc(args[0], args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	irc.attach_listeners(messages, commands, quit)
	show_messages(messages, quit, commands)

	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
