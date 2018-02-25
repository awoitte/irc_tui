package main

import (
	"fmt"
	"strings"
)

type Command struct {
	name      string
	arguments []string
}

var command_list = map[string][]string{
	"JOIN":    {"join a channel", "channel name"},
	"PART":    {"leave a channel", "channel name"},
	"MSG":     {"send message to a channel", "channel name", "message"},
	"WHISPER": {"send message to a user", "user", "message"},
	"QUIT":    {"leave the server"},
}

func translate_command(command_string string) string {

	if strings.Index(command_string, ":") == 0 {
		command := convert_into_command(command_string)
		return convert_command_to_raw_message(command)

	}
	fmt.Println("not a command, treating as raw", command_string)

	return command_string
}

func convert_into_command(text string) Command {
	arguments := strings.Split(text, " ")
	name := strings.Replace(arguments[0], ":", "", 1)
	NAME := strings.ToUpper(name)
	command := Command{NAME, arguments[1:]}

	return validate_command(command)
}

func validate_command(command Command) Command {
	command_description, ok := match_command_description(&command)

	if ok != true || len(command.arguments) < len(command_description)-1 {
		fmt.Println("incorrect name or number of arguments. requires %v", len(command_description))
		return Command{"INVALID", []string{}}
	}

	return command
}

func match_command_description(command *Command) ([]string, bool) {
	command_description, ok := command_list[command.name]
	if ok != true {
		return check_if_command_starts_with(command)
	}
	return command_description, true
}

func check_if_command_starts_with(command *Command) ([]string, bool) {
	for full_name, description := range command_list {
		if strings.Index(full_name, command.name) == 0 {
			command.name = full_name
			return description, true
		}
	}
	return nil, false
}

//TODO: this should be stored in one struct
func convert_command_to_raw_message(command Command) string {
	arguments := command.arguments

	switch command.name {
	case "JOIN":
		return fmt.Sprintf("JOIN #%v", arguments[0])
	case "PART":
		return fmt.Sprintf("PART #%v", arguments[0])
	case "QUIT":
		return fmt.Sprintf("QUIT", arguments[0])
	case "MSG":
		return fmt.Sprintf("PRIVMSG #%v %v", arguments[0], strings.Join(arguments[1:], " "))
	case "WHISPER":
		return fmt.Sprintf("PRIVMSG %v %v", arguments[0], strings.Join(arguments[1:], " "))
	}
	fmt.Println("name didn't match %v", command.name)
	return ""
}
