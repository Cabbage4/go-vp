package main

import (
	"flag"
	"fmt"
	"go-vp/butin"
	"io"
	"net"
)

var (
	outPort    int
	serverPort int
)

func main() {
	flag.IntVar(&outPort, "outPort", 10021, "")
	flag.IntVar(&serverPort, "serverPort", 10022, "")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", outPort))
	if err != nil {
		panic(err)
	}

	sLn, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		panic(err)
	}

	for {
		cnn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		sCnn, err := sLn.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			if _, err := io.Copy(cnn, sCnn); err != nil {
				if !butin.IsSkipError(err) {
					fmt.Println(err)
				}
			}

			cnn.Close()
			sCnn.Close()
		}()

		go func() {
			if _, err := io.Copy(sCnn, cnn); err != nil {
				if !butin.IsSkipError(err) {
					fmt.Println(err)
				}
			}

			cnn.Close()
			sCnn.Close()
		}()
	}
}
