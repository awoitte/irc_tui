package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterh/liner"
)

var (
	history_file_path = filepath.Join(os.TempDir(), ".irc_tui_history")
)

const text_view_limit = 4096
const refresh_rate = time.Millisecond * 100

func start_tui(chat_messages chan string, quit chan bool, input chan string) {
	log.Print("Loading...")
	go run_tui(quit, input)
	go append_messages_to_text_view(chat_messages)
}

func append_messages_to_text_view(messages chan string) {
	for {
		message, ok := <-messages
		if ok == true {
			log.Print(message)
		} else {
			break
		}
	}
}

func run_tui(quit chan bool, user_input chan string) {
	line := initialize_liner()
	defer line.Close()

	read_history(line)
	defer close_history(line)

	for {
		time.Sleep(refresh_rate)
		if command, err := line.Prompt(">"); err == nil {
			if command != "" {
				user_input <- command
				line.AppendHistory(command)
			}
		} else if err == liner.ErrPromptAborted {
			log.Print("Aborted")
			quit <- true
			break
		} else {
			log.Print("Error reading line: ", err)
		}
	}
	close(user_input)
}

func initialize_liner() *liner.State {
	line := liner.NewLiner()

	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)

	line.SetCompleter(get_omnicomplete_list)
	return line
}

func read_history(line *liner.State) {
	if f, err := os.Open(history_file_path); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func close_history(line *liner.State) {
	if f, err := os.Create(history_file_path); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func get_omnicomplete_list(input string) []string {
	possible_completions := []string{}
	for _, command := range command_list {
		if strings.HasPrefix(command.name, strings.ToUpper(input)) {
			possible_completions = append(possible_completions, command.name)
		}
	}
	return possible_completions
}
