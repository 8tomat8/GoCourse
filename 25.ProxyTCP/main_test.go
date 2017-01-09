package main

import (
	"testing"
	"net"
)

func createListener(addr string, t *testing.T) (net.Conn, net.Conn) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	connIn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	connOut, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	return connIn, connOut
}

func TestHandleConnect(t *testing.T) {
	var n int
	var err error
	response := make([]byte, 1024)

	testReq := []byte("foo")
	testRes := []byte("bar")

	rConnIn, rConnOut := createListener(":65321", t)
	lConnIn, lConnOut := createListener(":65322", t)

	_, err = lConnIn.Write(testReq)
	if err != nil {
		t.Fatal(err)
	}

	go handleConnect(lConnOut, rConnIn)

	n, err = rConnOut.Read(response)
	if err != nil {
		t.Fatal(err)
	}
	if string(response[:n]) != string(testReq) {
		t.Fatal("Proxy corrupts your request!")
	}

	_, err = rConnOut.Write(testRes)
	if err != nil {
		t.Fatal(err)
	}
	n, err = lConnIn.Read(response)
	if err != nil {
		t.Fatal(err)
	}

	if string(response[:n]) != string(testRes) {
		t.Error("Proxy corrupts your response!")
	}
}
