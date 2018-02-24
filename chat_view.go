package main

import (
	"strings"
	"time"

	ui "github.com/gizak/termui"
)

func runTUI(state *State, quit chan bool, input chan string) {
	err := ui.Init()
	if err != nil {
		panic("crashed at init")
	}
	defer ui.Close()

	chat_list := ui.NewList()
	chat_list.ItemFgColor = ui.ColorWhite
	chat_list.BorderLabel = "Lines"
	chat_list.BorderFg = ui.ColorCyan
	chat_list.Height = ui.TermHeight() - 2
	chat_list.Items = []string{}

	input_box := ui.NewPar("")
	input_box.TextFgColor = ui.ColorWhite
	input_box.Border = false
	input_box.Height = 2

	ui.Body.AddRows(ui.NewRow(ui.NewCol(12, 0, chat_list)),
		ui.NewRow(ui.NewCol(12, 0, input_box)))
	ui.Body.Align()

	ui.Handle("/sys/kbd/C-q", func(ui.Event) {
		ui.StopLoop()
		quit <- true
	})

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			only_newline := strings.Replace(state.current_text, "\r", "", -1)
			lines := strings.Split(only_newline, "\n")
			overflow := len(lines) - chat_list.Height + 1
			overflow_end := int_max(overflow, 0)
			chat_list.Items = lines[overflow_end:]

			ui.Render(ui.Body)
		}
	}()

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Clear()
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys", func(e ui.Event) {
		k, ok := e.Data.(ui.EvtKbd)
		if ok {
			switch k.KeyStr {
			case "<space>":
				input_box.Text += " "
			case "C-8": //backspace
				new_length := int_max(0, len(input_box.Text)-1)
				input_box.Text = get_first_n_characters_in_string(input_box.Text, new_length)

			case "<delete>":
				new_length := int_max(0, len(input_box.Text)-1)
				input_box.Text = get_last_n_characters_in_string(input_box.Text, new_length)
			case "C-u": //erase line
				input_box.Text = ""
			case "<enter>":
				input <- input_box.Text
				input_box.Text = ""
			case "<left>":
			case "<right>":
			case "<up>":
			case "<down>":
			case "<escape>":
			default:
				input_box.Text += k.KeyStr
			}
		}
	})
	ui.Loop()
}
