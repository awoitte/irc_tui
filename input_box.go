package main

import ui "github.com/gizak/termui"

func generate_input_box(key_inputs chan string, sent_lines chan string) *ui.Par {
	input_box := generate_input_box_ui_par()

	go func() {
		for {
			pressed_key := <-key_inputs
			handle_key_input(pressed_key, sent_lines, input_box)
		}
	}()

	return input_box
}

func generate_input_box_ui_par() *ui.Par {
	input_box := ui.NewPar("")
	input_box.TextFgColor = ui.ColorWhite
	input_box.Border = false
	input_box.Height = 2
	return input_box
}

func handle_key_input(key string, sent_lines chan string, input_box *ui.Par) {
	switch key {
	case "<space>":
		input_box.Text += " "
	case "C-8": //backspace
		input_box.Text = remove_last_letter(input_box.Text)
	case "<delete>":
		input_box.Text = remove_first_letter(input_box.Text)
	case "C-u": //erase line
		input_box.Text = ""
	case "<enter>":
		sent_lines <- input_box.Text
		input_box.Text = ""
	case "<left>":
	case "<right>":
	case "<up>":
	case "<down>":
	case "<escape>":
	default:
		if len(key) == 1 {
			input_box.Text += key
		}
	}
}
