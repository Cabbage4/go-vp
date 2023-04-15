package main

import (
	"fmt"
	"net"
	"strings"
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

		cnn.Write([]byte(strings.Repeat("hello-", 100)))
		cnn.Close()
	}
}
