package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	irc "github.com/awoitte/irc_client"
)

func join_command(command *Command, connection *irc.IRC, quit chan bool) error {
	arguments := command.arguments
	connection.SendRaw(fmt.Sprintf("JOIN #%v", arguments[0]))
	return nil
}

func part_command(command *Command, connection *irc.IRC, quit chan bool) error {
	arguments := command.arguments
	connection.SendRaw(fmt.Sprintf("PART #%v", arguments[0]))
	return nil
}

func quit_command(command *Command, connection *irc.IRC, quit chan bool) error {
	connection.SendRaw(fmt.Sprintf("QUIT"))
	quit <- true
	return nil
}

func message_command(command *Command, connection *irc.IRC, quit chan bool) error {
	arguments := command.arguments
	translated := fmt.Sprintf("PRIVMSG #%v %v", arguments[0], strings.Join(arguments[1:], " "))
	connection.SendRaw(translated)
	return nil
}

func whisper_command(command *Command, connection *irc.IRC, quit chan bool) error {
	arguments := command.arguments

	translated := fmt.Sprintf("PRIVMSG %v %v", arguments[0], strings.Join(arguments[1:], " "))
	connection.SendRaw(translated)
	return nil
}

func help_command(command *Command, connection *irc.IRC, quit chan bool) error {
	command_descriptions := &command_list
	for _, command_description := range *command_descriptions {
		log.Print(format_command_for_help(&command_description))
	}
	return nil
}

func raw_command(command *Command, connection *irc.IRC, quit chan bool) error {
	arguments := command.arguments
	connection.SendRaw(strings.Join(arguments, " "))
	return nil
}

func invalid_command(command *Command, connection *irc.IRC, quit chan bool) error {

	return errors.New("name didn't match any known commands")
}

func format_command_for_help(command *Command) string {
	arguments := strings.Join(command.arguments, " ")
	return fmt.Sprintln(command.name, ":", arguments, "\t-\t", command.description)
}
