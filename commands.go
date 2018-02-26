package main

import (
	"log"
	"strings"

	irc "github.com/awoitte/irc_client"
)

type Command struct {
	name            string
	description     string
	arguments       []string
	execute_command func(*Command, *irc.IRC, chan bool) error
}

type ParsedInput struct {
	name      string
	arguments []string
}

var (
	command_list []Command
)

func dispatch_commands(
	connection *irc.IRC,
	user_input,
	chat_messages chan string,
	quit chan bool) {

	init_commands()

	for {
		command_text := <-user_input
		command := convert_into_command(command_text)
		err := command.execute_command(&command, connection, quit)
		if err != nil {
			log.Print("couldn't match command: ", command.name)
		}
	}
}

func init_commands() {
	command_list = []Command{
		Command{
			"JOIN",
			"join a channel",
			[]string{"<channel name>"},
			join_command},
		Command{
			"PART",
			"leave a channel",
			[]string{"<channel name>"},
			part_command},
		Command{
			"MSG",
			"send message to a channel",
			[]string{"<channel name>", "<message>"},
			message_command},
		Command{
			"WHISPER",
			"send message to a user",
			[]string{"<user>", "<message>"},
			whisper_command},
		Command{
			"QUIT",
			"leave the server",
			[]string{},
			quit_command},
		Command{
			"HELP",
			"show this message",
			[]string{},
			help_command},
		Command{
			"RAW",
			"send unaltered command to the server",
			[]string{"<command>"},
			raw_command}}
}

func get_command_list() []Command {
	return command_list
}

func convert_into_command(text string) Command {
	parsed := parse_input(text)
	command_description := get_command_named(parsed.name)
	return convert_parsed_to_command(parsed, command_description)
}

func get_command_named(name string) (command *Command) {
	for _, description := range command_list {
		if strings.HasPrefix(description.name, name) {
			return &description
		}
	}
	return nil
}

func parse_input(text string) ParsedInput {
	text_parts := strings.Split(text, " ")
	name := strings.Replace(text_parts[0], ":", "", 1)
	NAME := strings.ToUpper(name)

	arguments := []string{}
	if len(text_parts) > 1 {
		arguments = text_parts[1:]
	}

	return ParsedInput{
		NAME,
		arguments}
}

func convert_parsed_to_command(parsed ParsedInput, command_description *Command) Command {

	if command_description != nil {
		is_valid := validate_command_arguments(command_description, parsed.arguments)

		if is_valid {
			command := *command_description
			command.arguments = parsed.arguments
			return command
		}
	}

	return Command{
		"INVALID",
		"incorrect name or number of arguments",
		[]string{},
		invalid_command}

}

func validate_command_arguments(command_description *Command, arguments []string) bool {
	if command_description != nil && len(arguments) >= len(command_description.arguments) {
		return true
	}

	return false
}
