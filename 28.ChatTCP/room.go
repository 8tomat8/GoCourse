package main

import (
	"fmt"
	"errors"
	"sync"
	"github.com/golang/glog"
)

type room struct {
	muRW      sync.RWMutex
	name      string
	users     []*user
	usersChan chan *user
}

type rooms map[string]*room

var activeRooms rooms = make(map[string]*room)

func NewRoom(name string) (*room, error) {
	if len(name) <= 0 {
		return nil, errors.New("Name should contain at least 1 symbol!\n")
	}

	if _, ok := activeRooms[name]; ok {
		return nil, fmt.Errorf("Room with name \"%v\" already exist! Simply send \"ADD_USER %v\" to connect.\n", name, name)
	}

	newRoom := &room{
		name:      name,
		usersChan: make(chan *user),
	}
	go newRoom.listenForNewUsers()

	activeRooms[name] = newRoom

	return newRoom, nil
}

func (r *room) listenForNewUsers() {
	for {
		u := <-r.usersChan

		r.muRW.Lock()
		r.users = append(r.users, u)
		r.muRW.Unlock()

		u.r = r
		err := u.receiveMsg(fmt.Sprintf("Hi! You have been added to room: %s. Users in room: %d\n", r.name, len(r.users)))
		if err != nil {
			glog.Error(err)
		}
		r.announceNewUser(u)
	}
}

func (r *room) announceNewUser(newUser *user) {
	r.sendMsg(fmt.Sprintf("User %s has connected to room\n", newUser.name), newUser)
}

func (r *room) announceUserChangedName(u *user, oldName string) {
	r.sendMsg(fmt.Sprintf("User %s has has changed username to %s\n", oldName, u.name), u)
}

func (r *room) sendMsg(msg string, sender *user) {
	for _, u := range r.users {
		if *u == *sender {
			continue
		}
		go u.receiveMsg(sender.name + ": \033[1m" + msg + "\033[0m" + "\n")
	}
	return
}

func (r *room) removeUser(removeUser *user) error {
	r.muRW.Lock()
	defer r.muRW.Unlock()

	index := -1

	for i, u := range r.users {
		if *u == *removeUser {
			index = i
		}
	}
	if index == -1 {
		err := errors.New("User is already removed from this room!")
		glog.Error(err)
		return err
	}
	r.users[len(r.users) - 1], r.users[index] = r.users[index], r.users[len(r.users) - 1]
	r.users = r.users[:len(r.users) - 1]
	r.sendMsg(fmt.Sprintf("User %s has left this room =(", removeUser.name), removeUser)
	return nil
}

func (r *room) changeRoom(newRoom *room, u *user) error {
	err := r.removeUser(u)
	if err != nil {
		return err
	}
	newRoom.usersChan <- u
	return nil
}
