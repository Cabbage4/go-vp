package main

import (
	"flag"
	"fmt"
	"go-vp/butin"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/websocket"
)

var (
	outPort    int
	serverPort int
	secretKey  string

	cnnChan = make(chan net.Conn, runtime.NumCPU())
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

	go func() {
		http.Handle("/", websocket.Handler(worker))
		http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
	}()

	for {
		cnn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		cnnChan <- cnn
	}
}

func worker(w *websocket.Conn) {
	cnn := <-cnnChan

	defer func() {
		w.Close()
		cnn.Close()
	}()

	var signature string
	if err := websocket.Message.Receive(w, &signature); err != nil {
		log.Println(err)
		return
	}
	if len(signature) != (64+128)/8 {
		log.Println("signature len error")
		return
	}
	if err := butin.CheckSignature([]byte(signature), secretKey); err != nil {
		log.Println(err)
		return
	}

	ch := make(chan bool, 2)
	startTime := time.Now()

	go func() {
		defer func() { ch <- true }()

		var buf string
		for {
			err := websocket.Message.Receive(w, &buf)
			if err != nil {
				log.Println(err)
				return
			}
			cnn.Write([]byte(buf))
		}
	}()

	go func() {
		defer func() { ch <- true }()

		buf := make([]byte, 1024)
		for {
			n, err := cnn.Read(buf)
			if err != nil {
				log.Println(err)
				break
			}

			websocket.Message.Send(w, buf[:n])
		}
	}()

	<-ch
	log.Printf("cost time = %s\n", time.Since(startTime))
}
