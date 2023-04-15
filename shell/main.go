package main

import (
	"flag"
	"fmt"
	"go-vp/butin"
	"io"
	"log"
	"net"
	"runtime"
	"time"
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
		ch <- struct{}{}

		go func() {
			if err := worker(); err != nil {
				log.Println(err)
			}
		}()
	}
}

func worker() error {
	defer func() {
		<-ch
	}()

	cnn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer cnn.Close()

	signature, err := butin.GenSignature(time.Now().Unix(), secretKey)
	if err != nil {
		return err
	}
	if _, err := cnn.Write(signature); err != nil {
		return err
	}

	// NAT timeout
	cnn.SetReadDeadline(time.Now().Add(2 * time.Hour))

	tranDataKeyBuf := make([]byte, len(butin.TranDataKey))
	n, err := cnn.Read(tranDataKeyBuf)
	if err != nil {
		return err
	}
	if string(tranDataKeyBuf[:n]) != butin.TranDataKey {
		return fmt.Errorf("tranDataKey error:%s", tranDataKeyBuf[:n])
	}

	tnn, err := net.Dial("tcp", fmt.Sprintf(":%d", target))
	if err != nil {
		return err
	}
	defer tnn.Close()

	ch := make(chan struct{}, 3)
	go func() {
		if _, err := io.Copy(cnn, tnn); err != nil {
			log.Println(err)
		}
		ch <- struct{}{}
	}()

	go func() {
		if _, err := io.Copy(tnn, cnn); err != nil {
			log.Println(err)
		}
		ch <- struct{}{}
	}()

	<-ch

	return nil
}
