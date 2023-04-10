package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"time"

	"golang.org/x/net/websocket"
)

var (
	addr      string
	target    int
	secretKey string

	ch = make(chan struct{}, runtime.NumCPU())
)

func main() {
	flag.StringVar(&addr, "addr", "localhost:10022", "")
	flag.IntVar(&target, "target", 8081, "")
	flag.StringVar(&secretKey, "secretKey", "jerryzhuo@abcd", "")
	flag.Parse()

	for {
		go worker()
		<-ch
	}
}

func worker() {
	defer func() {
		<-time.After(2 * time.Second)
		ch <- struct{}{}
	}()

	w, err := websocket.Dial(fmt.Sprintf("ws://%s", addr), "", fmt.Sprintf("http://%s", addr))
	if err != nil {
		log.Println(err)
		return
	}

	cnn, err := net.Dial("tcp", fmt.Sprintf(":%d", target))
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		w.Close()
		cnn.Close()
	}()

	ch := make(chan bool, 3)
	startTime := time.Now()

	go func() {
		// 解决NAT问题
		<-time.After(2 * time.Hour)
		ch <- true
	}()

	go func() {
		defer func() { ch <- true }()

		if _, err := io.Copy(w, cnn); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer func() { ch <- true }()

		if _, err := io.Copy(cnn, w); err != nil {
			log.Println(err)
		}
	}()

	<-ch
	log.Printf("cost time = %s\n", time.Since(startTime))
}
