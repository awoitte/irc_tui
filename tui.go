package main

import (
	"time"

	ui "github.com/gizak/termui"
)

const text_view_limit = 4096
const refresh_rate = time.Millisecond * 100

type State struct {
	current_text string
}

func show_messages(messages chan string, quit chan bool, input chan string) {
	state := State{"Loading..."}
	go runTUI(&state, quit, input)
	go append_messages_to_text_view(messages, &state)
}

func append_messages_to_text_view(messages chan string, state *State) {
	var current_text string
	for {
		message := <-messages
		current_text = current_text + message
		if len(current_text) > text_view_limit {
			current_text = get_last_n_characters_in_string(current_text, text_view_limit)
		}
		state.current_text = current_text
	}
}

func runTUI(state *State, quit chan bool, sent_lines chan string) {
	err := ui.Init()
	if err != nil {
		panic("crashed at init")
	}
	defer ui.Close()

	chat_list := generate_chat_list()

	key_inputs := make(chan string)
	input_box := generate_input_box(key_inputs, sent_lines)

	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, chat_list)),
		ui.NewRow(ui.NewCol(12, 0, input_box)))
	ui.Body.Align()

	shutdown_and_push_quit_on_key("C-q", quit)
	go func() {
		for {
			time.Sleep(refresh_rate)
			update_chat_list(state.current_text, chat_list)
			ui.Render(ui.Body)
		}
	}()
	refresh_on_resize()
	pass_keyboard_events_to_channel(key_inputs)
	ui.Loop()
}

func refresh_ui() {
	ui.Body.Width = ui.TermWidth()
	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

func pass_keyboard_events_to_channel(key_inputs chan string) {
	ui.Handle("/sys", func(e ui.Event) {
		k, ok := e.Data.(ui.EvtKbd)
		if ok {
			key_inputs <- k.KeyStr
		}
	})
}

func refresh_on_resize() {
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		refresh_ui()
	})

}

func shutdown_and_push_quit_on_key(key string, quit chan bool) {
	ui.Handle("/sys/kbd/"+key, func(ui.Event) {
		ui.StopLoop()
		quit <- true
	})
}
