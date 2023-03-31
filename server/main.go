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
	secretKey  string
)

func main() {
	flag.IntVar(&outPort, "outPort", 10021, "")
	flag.IntVar(&serverPort, "serverPort", 10022, "")
	flag.StringVar(&secretKey, "secretKey", "jerryzhuo@abcd", "")
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
			cnn.Close()

			fmt.Println(err)
			continue
		}

		signatureBuf := make([]byte, (64+128)/8)
		n, err := sCnn.Read(signatureBuf)
		if err != nil {
			cnn.Close()
			sCnn.Close()

			fmt.Println(err)
			continue
		}

		if n != (64+128)/8 {
			cnn.Close()
			sCnn.Close()

			fmt.Println("signature len error")
			continue
		}

		if err := butin.CheckSignature(signatureBuf, secretKey); err != nil {
			cnn.Close()
			sCnn.Close()

			fmt.Println(err)
			continue
		}

		fmt.Printf("[%s->%s]link\n", cnn.RemoteAddr().String(), sCnn.RemoteAddr().String())

		go func() {
			if _, err := io.Copy(cnn, sCnn); err != nil {
				if !butin.IsSkipError(err) {
					fmt.Printf("[%s->%s]error|%s\n", cnn.RemoteAddr().String(), sCnn.RemoteAddr().String(), err)
				}
			}

			cnn.Close()
			sCnn.Close()
		}()

		go func() {
			if _, err := io.Copy(sCnn, cnn); err != nil {
				if !butin.IsSkipError(err) {
					fmt.Printf("[%s->%s]error|%s\n", cnn.RemoteAddr().String(), sCnn.RemoteAddr().String(), err)
				}
			}

			cnn.Close()
			sCnn.Close()
			fmt.Printf("[%s->%s]done\n", cnn.RemoteAddr().String(), sCnn.RemoteAddr().String())
		}()
	}
}
