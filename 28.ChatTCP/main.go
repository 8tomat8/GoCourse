package main

import (
	"net"
	"flag"
	"github.com/golang/glog"
	"fmt"
	"strings"
	"errors"
)

var maxTries int

func main() {
	host := flag.String("host", "127.0.0.1:12345", "Host")
	maxTries = int(*flag.Uint("maxTries", 3, "Max tries"))
	flag.Parse()

	listener, err := net.Listen("tcp", *host)
	if err != nil {
		glog.Fatal(err)
	}

	// DELETE THIS SHIT!!!!!
	for _, e := range []string{"test"} {
		_, err := NewRoom(e)
		if err != nil {
			panic(err)
		}
	}

	for {
		newConn, err := listener.Accept()

		if err != nil {
			glog.Error(err)
		}
		go handleNewConn(newConn)
	}
}

func handleNewConn(conn net.Conn) {

	newUser, err := NewUser(conn)
	if err != nil {
		glog.Error(err)
		fmt.Fprint(conn, err)
		conn.Close()
		return
	}

	data := make([]byte, 256)
	err = newUser.changeName(data)
	if err != nil {
		glog.Error(err)
		return
	}

	r, err := selectRoom(data, newUser)
	if err != nil {
		glog.Error(err)
		return
	}
	r.usersChan <- newUser
}

func parseCommand(s string) []string {
	s = strings.Trim(strings.TrimSuffix(s, "\r\n"), " ")
	params := strings.SplitN(s, " ", 2)
	for i, p := range params {
		params[i] = strings.Trim(p, " ")
	}
	return params
}

func askUserName(u *user, data []byte) (userName string, err error) {
	i := 1
	var n int

	for ; i <= maxTries; i++ {
		u.receiveMsg("Please set your name: ")
		n, err = u.conn.Read(data)
		if err != nil {
			glog.Error(err)
			u.conn.Close()
			return
		}
		params := parseCommand(string(data[:n]))

		if len(params) >= 1 {
			userName = params[0]
		}
		break
	}

	if i >= maxTries {
		err = errors.New("Sorry, but you are an idiot... Name isn't changed.")
	}
	return
}

func selectRoom(data []byte, u *user) (targetRoom *room, err error) {
	i := 1
	var n int

ConnLoop:
	for ; i <= maxTries; i++ {
		err = u.receiveMsg("\nTo choose room select room_name from list and send ADD_USER [room_name]\n\n")
		if err != nil {
			glog.Error(err)
			continue
		}

		// Print list of actual rooms
		for name := range activeRooms {
			err = u.receiveMsg(fmt.Sprintf("Room: %v \n", name))
			if err != nil {
				glog.Error(err)
				continue
			}
		}

		err = u.receiveMsg("\nOr, send CREATE_ROOM [new_room_name] to create new one\n")
		if err != nil {
			glog.Error(err)
		}

		// Reading response from new client
		n, err = u.conn.Read(data)
		if err != nil {
			glog.Error(err)
			u.conn.Close()
			return
		}
		params := parseCommand(string(data[:n]))
		if len(params) != 2 {
			err = u.receiveMsg("Command should have exactly 1 parametem\n")
			continue
		}

		switch params[0] {
		case "CREATE_ROOM":
			targetRoom, err = NewRoom(params[1])
			if err != nil {
				glog.Error(err)
				_ = u.receiveMsg("\n" + err.Error() + "\n")
				continue ConnLoop
			}
		case "ADD_USER":
			r, ok := activeRooms[params[1]]
			if !ok {
				err = u.receiveMsg(fmt.Sprintf("\nRoom with name \"%v\" was not found!\n\n", params[1]))
				break
			}
			targetRoom = r
		default:
			continue
		}
		//targetRoom.usersChan <- u
		break
	}

	if i >= maxTries {
		u.conn.Close()
	}
	return
}
