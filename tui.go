package main

const text_view_limit = 4096

type State struct {
	current_text string
}

//TODO: this will end up in a controller file
func show_messages(messages chan string, quit chan bool, input chan string) {
	//setup
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
