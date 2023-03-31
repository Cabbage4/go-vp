package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	for {
		cnn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		cnn.Write([]byte("hello world"))
		cnn.Close()
	}
}
