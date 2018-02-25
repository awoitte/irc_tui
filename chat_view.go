package main

import (
	"strings"

	ui "github.com/gizak/termui"
)

func generate_chat_list() *ui.List {

	chat_list := ui.NewList()
	chat_list.ItemFgColor = ui.ColorWhite
	chat_list.BorderLabel = "Lines"
	chat_list.BorderFg = ui.ColorCyan
	chat_list.Height = ui.TermHeight() - 2
	chat_list.Items = []string{}

	return chat_list
}

func update_chat_list(text string, chat_list *ui.List) {
	only_newline := strings.Replace(text, "\r", "", -1)
	lines := strings.Split(only_newline, "\n")
	overflow := len(lines) - chat_list.Height + 1
	overflow_end := int_max(overflow, 0)
	chat_list.Items = lines[overflow_end:]

}
