package main

import (
	"fmt"
	"net"
)

func main() {
	cnn, err := net.Dial("tcp", ":10021")
	if err != nil {
		panic(err)
	}
	cnn.Write([]byte("hi"))

	buf := make([]byte, 1024)
	n, err := cnn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
	cnn.Close()
}
