package main

import (
	"testing"
	"net"
	"strings"
	"time"
)

var (
	host      string = "127.127.127.127:65000"
	userName  string = "TestName"
	userName2 string = "TestName_2"
	roomName  string = "TestRoom"
	command   string = "/command"
	text      string = "Test text"
	suffix    string = "\r\n"
	data      []byte = make([]byte, 1<<10)
)

func TestRunApp(t *testing.T) {
	stop := make(chan empty)
	ready := make(chan empty)

	go runApp(&host, stop, ready)
	<-ready

	if len(Chat.lobby.users) != 0 {
		t.Error("Created Chat lobby is not empty!")
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Fatal(err)
	}

	n, err := conn.Read(data)

	if err != nil {
		t.Fatal(err, data[:n])
	} else if n == 0 {
		t.Error("We did not receive any data from application after connection cover!")
	}

	n, err = conn.Write([]byte(userName))
	if err != nil {
		t.Error(err)
	} else if n == 0 {
		t.Error("We could not send any data to connection!")
	}

	n, err = conn.Read(data)
	if err != nil {
		t.Fatal(err, data[:n])
	} else if n == 0 {
		t.Error("We did not receive any data from application after registration!")
	}

	if len(Chat.lobby.users) != 1 {
		t.Error("New TestUser was not added to lobby after registration!")
	}
	user, ok := Chat.lobby.users[userName]
	if !ok {
		t.Error("User was not added to Chat!")
	} else if user.Name != userName {
		t.Error("Username was corrupted!", user.Name, "!=", userName)
	}
	stop <- empty{}
	<- stop
}

func TestNewMessage(t *testing.T) {
	user := NewSystemUser()

	msg := NewMessage([]byte(command+" "+text+" "+suffix), user)
	if *msg.Sender != *user {
		t.Error("User in new message was corrupted! ", msg.Sender, "!=", user)
	} else if msg.Command != Command(command) {
		t.Error("Command in new message was corrupted! ", msg.Command, "!=", command)
	} else if msg.Msg != text || strings.HasSuffix(msg.Msg, suffix) {
		t.Error("Text in new message was corrupted! ", msg.Msg, "!=", text)
	}

	msg = NewMessage([]byte(command+" "+suffix), user)
	if msg.Command != Command(command) {
		t.Error("Command in new message was corrupted! ", msg.Command, "!=", command)
	} else if msg.Msg != "" || strings.HasSuffix(msg.Msg, suffix) {
		t.Error("Text in new message was corrupted!", msg.Msg, "!= \"\"", )
	}
}

func connectNewUser(name string, conn net.Conn) error {
	_, err := conn.Write([]byte(name))
	if err != nil {
		return err
	}

	_, err = conn.Read(data)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond)
	return nil
}

func TestRoom(t *testing.T) {
	stop := make(chan empty)
	ready := make(chan empty)

	go runApp(&host, stop, ready)
	<-ready

	_, err := Chat.NewRoom("")
	if err == nil {
		t.Error("Room name must have at least one char!")
	}

	conn1, err := net.Dial("tcp", host)
	if err != nil {
		t.Error(err)
	}
	err = connectNewUser(userName, conn1)
	if err != nil {
		t.Fatal(err)
	}
	u1, ok := Chat.lobby.users[userName]
	if !ok {
		t.Fatal("First user was not added!")
	}

	conn2, err := net.Dial("tcp", host)
	if err != nil {
		t.Error(err)
	}
	err = connectNewUser(userName2, conn2)
	if err != nil {
		t.Fatal(err)
	}
	u2, ok := Chat.lobby.users[userName2]
	if !ok {
		t.Fatal("Second user was not added!")
	}

	// Tests start
	if len(Chat.lobby.users) == 0 || u1.r.name != Chat.lobby.name || u2.r.name != Chat.lobby.name {
		t.Fatal("Some users were not added to room!")
	}

	// Create test room and add first user to it
	Chat.lobby.Receive <- &message{u1, roomName, create_room}
	time.Sleep(time.Millisecond)

	_, err = conn1.Read(data)
	if err != nil {
		t.Error(err)
	}

	if len(Chat.rooms) != 1 {
		t.Fatal("Room was not created! len(Chat.rooms) = ", len(Chat.rooms))
	}

	// Add second user to test room
	Chat.lobby.Receive <- &message{u2, roomName, change_room}
	time.Sleep(time.Millisecond)

	_, err = conn2.Read(data)
	if err != nil {
		t.Error(err)
	}

	if len(Chat.rooms[roomName].users) != 2 {
		t.Fatal("Second user was not added! len(Chat.rooms) = ", len(Chat.rooms))
	}
}
