package main

import "strings"

type Command string

const (
	send_msg    Command = "/send_msg"
	lobby       Command = "/lobby"
	help        Command = "/help"
	change_name Command = "/change_name"
	create_room Command = "/create_room"
	list_rooms  Command = "/list_rooms"
	change_room Command = "/change_room"
)

type message struct {
	Sender  *user
	Msg     string
	Command Command
}

func NewMessage(data []byte, sender *user) *message {
	s := string(data)
	s = strings.Trim(strings.TrimSuffix(s, "\r\n"), " ")

	params := strings.SplitN(s, " ", 2)
	for i, p := range params {
		params[i] = strings.Trim(p, " ")
	}
	if len(params) == 1 {
		params = append(params, "")
	}

	return &message{Sender: sender, Command: Command(params[0]), Msg: params[1]}
}

func NewSystemMessage(s string) *message {
	return &message{Sender: Chat.systemUser, Msg: "\033[1m" + s + "\033[0m"}
}
