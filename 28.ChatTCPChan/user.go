package main

import (
	"net"
	"fmt"
	"github.com/golang/glog"
)

type user struct {
	Name  string
	r     *room
	conn  net.Conn
	MsgIn chan *message
	stop  chan struct{}
}

func (u *user) listenMsgIn() {
	for {
		select {
		case msg := <-u.MsgIn:
			_, err := fmt.Fprintf(u.conn, "%s: %s \n", msg.Sender.Name, msg.Msg)
			if err != nil {
				glog.Error(err)
				u.disconnect()
				return
			}
		case <-u.stop:
			return
		}
	}
}

func (u *user) listenMsgOut() {
	data := make([]byte, 1<<10)
	for {
		n, err := u.conn.Read(data)
		if err != nil {
			u.disconnect()
			return
		}
		msg := NewMessage(data[:n], u)
		u.r.Receive <- msg
	}
}

func NewUser(conn net.Conn, name string) *user {
	u := &user{
		Name:  name,
		conn:  conn,
		MsgIn: make(chan *message),
		stop:  make(chan struct{}),
		r:     Chat.lobby,
	}
	go u.listenMsgIn()
	go u.listenMsgOut()
	return u
}

func (u *user) changeName(s string) {
	if len(s) >= 1 {
		oldName := u.Name
		u.Name = s
		u.r.Send <- NewSystemMessage(fmt.Sprintf("User %s has changed name to %s", oldName, u.Name))
		return
	}
	u.MsgIn <- NewSystemMessage("Name is too short!")
}

func (u *user) disconnect() {
	u.r.removeUser(u)
	close(u.stop)
	u.conn.Close()
}

func NewSystemUser() *user {
	return &user{Name: "System"}
}
