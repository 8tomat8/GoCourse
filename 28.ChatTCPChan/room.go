package main

import "fmt"

type room struct {
	c            *chat
	name         string
	users        map[string]*user
	Receive      chan *message
	Send         chan *message
	usersChan    chan *user
	sentMessages uint
}

func (r *room) sendMsg() {
	for {
		msg := <-r.Send
		for _, u := range r.users {
			if *u == *msg.Sender {
				continue
			}
			go func(u *user) { u.MsgIn <- msg }(u)
		}
	}
}

func (r *room) listenMsg() {
	for {
		msg := <-r.Receive
		switch msg.Command {
		case send_msg:
			r.Send <- msg
			r.sentMessages += 1
		case lobby:
			r.removeUser(msg.Sender)
			Chat.lobby.usersChan <- msg.Sender
		case help:
			fallthrough
		default:
			msg.Sender.MsgIn <- NewSystemMessage(`
	List of avalible commands:
	/send_msg [text]
	/lobby - Leave room and go to lobby

	Lobby commands:
	/change_name [name]
	/change_room [room_name]
	/create_room  [room_name]
	/list_rooms`)
		}

	}
}

func (r *room) listenLobbyMsg() {
	for {
		msg := <-r.Receive
		switch msg.Command {
		case change_name:
			msg.Sender.changeName(msg.Msg)
		case create_room:
			newRoom, err := r.c.NewRoom(msg.Msg)
			if err != nil {
				msg.Sender.MsgIn <- NewSystemMessage(err.Error())
				continue
			}
			msg.Sender.MsgIn <- NewSystemMessage("New room was successfully created.")

			msg.Sender.r.removeUser(msg.Sender)
			newRoom.usersChan <- msg.Sender
		case list_rooms:
			var list string = "\nList of avalible rooms:\n"
			for _, elem := range r.c.rooms {
				list = list + fmt.Sprintf("%s - %d users online\n", elem.name, len(elem.users))
			}
			msg.Sender.MsgIn <- NewSystemMessage(list)
		case change_room:
			targetRoom, ok := r.c.rooms[msg.Msg]
			if !ok {
				msg.Sender.MsgIn <- NewSystemMessage(fmt.Sprintf("Room with name %s does not exist.", msg.Msg))
				continue
			}
			r.removeUser(msg.Sender)
			targetRoom.usersChan <- msg.Sender
		case help:
			fallthrough
		default:
			msg.Sender.MsgIn <- NewSystemMessage(`
	List of avalible commands:
	/change_name [name]
	/change_room [room_name]
	/create_room  [room_name]
	/list_rooms

	Room commands:
	/send_msg [text]
	/lobby - Leave room and go to lobby`)
		}
	}
}

func (r *room) listenUsers() {
	for {
		newUser := <-r.usersChan
		newUser.r = r

		r.users[newUser.Name] = newUser

		r.Send <- NewSystemMessage(fmt.Sprintf("User %s connected to this room (%s)", newUser.Name, r.name))
	}
}

func (r *room) removeUser(u *user) {
	delete(r.users, u.Name)
	r.Send <- NewSystemMessage(fmt.Sprintf("User %s has left this room.", u.Name))
}
