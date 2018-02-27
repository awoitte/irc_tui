package main

type IRC_Connection interface {
	ReadLoop(int, func(string) error) error
	SendRaw(string)
	JoinChannel(string) error
	PartChannel() error
	SendMessage(string) error
}
