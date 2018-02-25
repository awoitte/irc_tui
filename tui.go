package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterh/liner"
)

//TODO: command list should be consolidated
var (
	history_fn = filepath.Join(os.TempDir(), ".irc_tui_history")
	commands   = []string{":join", ":part", ":quit", ":msg", ":whisper"}
)

const text_view_limit = 4096
const refresh_rate = time.Millisecond * 100

func show_messages(messages chan string, quit chan bool, input chan string) {
	log.Print("Loading...")
	go runTUI(quit, input)
	go append_messages_to_text_view(messages)
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

func runTUI(quit chan bool, sent_lines chan string) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	for {
		time.Sleep(refresh_rate)
		if command, err := line.Prompt(">"); err == nil {
			if command != "" {
				sent_lines <- command
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

	if f, err := os.Create(history_fn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
