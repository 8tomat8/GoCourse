package main

import (
	"flag"
	"net"
	"github.com/golang/glog"
	"fmt"
	"strings"
)

var Chat *chat

func main() {
	host := flag.String("host", "127.0.0.1:12345", "Host")
	flag.Parse()

	Chat = NewChat()
	listener, err := net.Listen("tcp", *host)
	if err != nil {
		glog.Fatal(err)
	}

	for {
		newConn, err := listener.Accept()
		if err != nil {
			glog.Error(err)
			continue
		}
		go handleNewConn(newConn)
	}
}

func handleNewConn(newConn net.Conn) {
	data := make([]byte, 1 << 6)
	_, err := fmt.Fprint(newConn, "Please choose nickname: ")
	if err != nil {
		newConn.Close()
		return
	}
	n, err := newConn.Read(data)
	if err != nil {
		newConn.Close()
		return
	}
	Chat.users <- NewUser(newConn, strings.TrimSuffix(string(data[:n]), "\r\n"))
}
