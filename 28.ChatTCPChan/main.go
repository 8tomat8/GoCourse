package main

import (
	"flag"
	"net"
	"github.com/golang/glog"
	"fmt"
	"strings"
	"os/signal"
	"os"
)

var Chat *chat

func init() {
	Chat = NewChat()
}

type empty struct{}

func main() {

	host := flag.String("host", "127.0.0.1:12345", "Host")
	flag.Parse()

	ready := make(chan empty)
	stop := make(chan empty)
	go runApp(host, stop, ready)
	<-ready
	glog.Info("Chat server started on ", host)

	sig := make(chan<- os.Signal)
	signal.Stop(sig)

	stop <- empty{}
	glog.Info("Chat server has been interrupted...")
}

func runApp(host *string, stop chan empty, ready chan empty) {
	listener, err := net.Listen("tcp", *host)
	if err != nil {
		glog.Fatal(err)
	}
	defer listener.Close()

	conns := make(chan net.Conn)

	go func() {
		var errorsCount uint = 0
		for {
			newConn, err := listener.Accept()
			if err != nil {
				glog.Error(err)
				errorsCount += 1
				if errorsCount >= 2 {
					return
				} else {
					continue
				}
			}
			errorsCount = 0
			conns <- newConn
		}
	}()

	close(ready)

	for {
		select {
		case <-stop:
			stop <- empty{}
			return
		case conn := <- conns:
			go handleNewConn(conn)
		}
	}
}

func handleNewConn(newConn net.Conn) {
	data := make([]byte, 1<<6)
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
