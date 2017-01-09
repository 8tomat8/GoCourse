package main

import (
	"net"
	"log"
	"flag"
)

func main() {
	var listen = flag.String("l", "127.0.0.1:11112", "host:port")
	var remote = flag.String("r", "127.0.0.1:11111", "host:port")
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	lConn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer lConn.Close()

	rConn, err := net.Dial("tcp", *remote)
	if err != nil {
		panic(err)
	}
	defer rConn.Close()

	handleConnect(lConn, rConn)
}

func handleConnect(in net.Conn, out net.Conn) error {
	data := make([]byte, 1024)
	for {
		read, err := in.Read(data)
		if err != nil {
			log.Fatalln(err)
			continue
		} else if read == 0 {
			continue
		}

		written, err := out.Write(data[:read])
		if err != nil {
			log.Fatalln(err)
			continue
		} else if read != written {
			log.Fatalln("Request was corrupted.")
			continue
		}

		read, err = out.Read(data)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		written, err = in.Write(data[:read])
		if err != nil {
			log.Fatalln(err)
			continue
		} else if read != written {
			log.Fatalln("Response was corrupted.")
			continue
		}
	}
}
