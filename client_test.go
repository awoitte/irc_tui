package main

import (
	"testing"
	"time"
)

func Test_it_should_call_join_on_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	test_response_to_input("j test", input, output)
	if connection.join_called != "test" {
		t.Errorf("Join called with wrong value '%v'", connection.join_called)
	}
}

func Test_it_should_call_part_on_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	test_response_to_input("p", input, output)
	if !connection.part_called {
		t.Error("Part not called")
	}
}

func Test_it_should_call_quit_on_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	go func() { <-output }() //will wait to push echo message to screen
	test_channel_called_on_input("q", input, quit)
	if connection.raw_called != "QUIT" {
		t.Errorf("Raw called with wrong value '%v'", connection.raw_called)
	}
}
func Test_it_should_call_message_on_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	test_response_to_input("m test_message", input, output)
	if connection.message_called != "test_message" {
		t.Errorf("Message called with wrong value '%v'", connection.raw_called)
	}
}

func Test_it_should_call_whisper_on_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	test_response_to_input("w user test_message", input, output)
	if connection.raw_called != "PRIVMSG user test_message" {
		t.Errorf("Raw called with wrong value '%v'", connection.raw_called)
	}
}

func Test_it_should_call_raw_on_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	test_response_to_input("r user test_message", input, output)
	if connection.raw_called != "user test_message" {
		t.Errorf("Raw called with wrong value '%v'", connection.raw_called)
	}
}

func Test_it_should_give_help_short_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	var response string
	go func() { input <- "h" }()
	select {
	case response = <-output:
	case <-time.After(time.Millisecond * 100):
	}
	if response == "" {
		t.Errorf("No response recieved")
	}
}

func Test_it_should_send_message_on_invalid_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	test_response_to_input("test_message", input, output)
	if connection.message_called != "test_message" {
		t.Errorf("Message called with wrong value '%v'", connection.message_called)
	}
}

func Test_it_should_give_warning_on_invalid_input(t *testing.T) {
	connection := Mock_Connection{}
	client := IRC_client{&connection}
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)
	client.attach_listeners(output, input, quit)
	var response string
	go func() { input <- "j" }() // no channel name
	select {
	case response = <-output:
	case <-time.After(time.Millisecond * 100):
	}
	if response == "" {
		t.Errorf("No response recieved")
	}
}
