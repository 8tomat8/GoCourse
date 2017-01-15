package main

import (
	"net"
	"fmt"
	"github.com/golang/glog"
	"io"
)

type user struct {
	r    *room
	name string
	conn net.Conn
}

func NewUser(conn net.Conn) (*user, error) {
	newUser := &user{conn:conn}
	go newUser.listener()
	return newUser, nil
}

func (u *user) changeName(data []byte) error {
	i := 1
	for ; i < maxTries; i++ {
		userName, err := askUserName(u, data)
		if err != nil {
			u.receiveMsg(err.Error())
			continue
		}

		if len(userName) == 0 {
			u.receiveMsg("Name should contain at least 1 symbol!\n")
			continue
		}
		u.name = userName
		u.receiveMsg(fmt.Sprintf("Your username was succsessfully changed to %s\n", u.name))
		break
	}
	if i >= maxTries {
		return fmt.Errorf("User %s can't get message from %d tries. User was disconnected!", u.name, maxTries)
	}
	return nil
}

func (u *user) receiveMsg(msg string) error {
	i := 1
	for ; i < maxTries; i++ {
		_, err := fmt.Fprint(u.conn, msg)
		if err != nil {
			glog.Error(err)
			continue
		}
		break
	}
	if i >= maxTries {
		err := u.r.removeUser(u)
		if err != nil {
			return err
		}
		u.conn.Close()
		return fmt.Errorf("User %s can't get message from %d tries. User was disconnected!", u.name, maxTries)
	}

	return nil
}

func (u *user) listener() {
	data := make([]byte, 4096)
	defer u.conn.Close()

userListenerLoop:
	for {
		n, err := u.conn.Read(data)
		if err == io.EOF {
			u.r.removeUser(u)
			return
		} else if err != nil {
			glog.Error(err)
			return
		}
		params := parseCommand(string(data[:n]))
		if len(params) < 1 {
			continue
		}

		switch params[0] {
		case "SEND_MSG":
			if len(params) > 1 {
				u.r.sendMsg(params[1], u)
			}
		case "CHANGE_ROOM":
			data := make([]byte, 256)
			newRoom, err := selectRoom(data, u)
			if err != nil {
				glog.Error(err)
				u.receiveMsg(fmt.Sprintf(err.Error()))
				continue userListenerLoop
			}
			u.r.changeRoom(newRoom, u)
		case "CHANGE_NAME":
			data := make([]byte, 256)
			err := u.changeName(data)
			if err != nil {

			}
			oldName := u.name

			u.r.announceUserChangedName(u, oldName)
		}

	}

}
