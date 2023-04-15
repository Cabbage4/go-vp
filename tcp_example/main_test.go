package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	cnn, err := net.Dial("tcp", ":10021")
	if err != nil {
		panic(err)
	}
	cnn.Write([]byte(strings.Repeat("hgggg", 10000)))

	buf := make([]byte, 1024)
	n, err := cnn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
	cnn.Close()
}
