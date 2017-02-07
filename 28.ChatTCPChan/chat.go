package main

import (
	"errors"
)

type chat struct {
	users      chan *user
	systemUser *user
	lobby      *room
	rooms      map[string]*room
}

func (c *chat) NewRoom(name string) (*room, error) {
	if len(name) < 1 {
		return nil, errors.New("Room mane is too short!")
	}
	r := &room{
		name:      name,
		usersChan: make(chan *user),
		Receive:   make(chan *message, 5),
		Send:      make(chan *message, 5),
		users:     make(map[string]*user),
	}
	go r.listenMsg()
	go r.sendMsg()
	go r.listenUsers()

	c.rooms[r.name] = r
	return r, nil
}

func (c *chat) listenUsers() {
	for {
		newUser := <-c.users
		newUser.r = c.lobby

		c.lobby.usersChan <- newUser
	}
}

func NewChat() *chat {
	c := &chat{
		users:      make(chan *user),
		systemUser: NewSystemUser(),
		rooms:      make(map[string]*room),
	}
	r := &room{
		c:            c,
		name:         "Lobby",
		Receive:      make(chan *message, 5),
		Send:         make(chan *message, 5),
		usersChan:    make(chan *user),
		users:        make(map[string]*user),
	}
	go r.listenLobbyMsg()
	go r.sendMsg()
	go r.listenUsers()
	c.lobby = r

	go c.listenUsers()

	return c
}
