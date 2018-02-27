package main

import (
	"time"
)

type Mock_Connection struct {
	join_called    string
	raw_called     string
	part_called    bool
	message_called string
}

func (mock *Mock_Connection) ReadLoop(_ int, _ func(string) error) error {
	return nil
}
func (mock *Mock_Connection) SendRaw(command string) {
	mock.raw_called = command
}
func (mock *Mock_Connection) JoinChannel(channel string) error {
	mock.join_called = channel
	return nil
}
func (mock *Mock_Connection) PartChannel() error {
	mock.part_called = true
	return nil
}
func (mock *Mock_Connection) SendMessage(message string) error {
	mock.message_called = message
	return nil
}

func test_response_to_input(
	input_message string,
	input, output chan string) (response string) {

	go func() {
		input <- input_message
	}()

	select {
	case response = <-output:
	case <-time.After(time.Millisecond * 100):
		panic("Timed out waiting for response")
	}
	<-time.After(time.Millisecond * 10) //Mock_Connection recieves value immediately AFTER message is sent
	return
}

func test_channel_called_on_input(
	input_message string,
	input chan string,
	channel chan bool) {

	go func() {
		input <- input_message
	}()

	select {
	case <-channel:
		return
	case <-time.After(time.Millisecond * 100):
		panic("Timed out waiting for response")
	}
}
