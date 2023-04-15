package main

import (
	"flag"
	"fmt"
	"go-vp/butin"
	"io"
	"log"
	"net"
	"runtime"
)

var (
	cport     int
	sport     int
	secretKey string

	sch = make(chan net.Conn, runtime.NumCPU())
	cch = make(chan net.Conn, runtime.NumCPU())
)

func main() {
	flag.IntVar(&cport, "outPort", 10021, "")
	flag.IntVar(&sport, "serverPort", 10022, "")
	flag.StringVar(&secretKey, "secretKey", "jerryzhuo@abcd", "")
	flag.Parse()

	cln, err := net.Listen("tcp", fmt.Sprintf(":%d", cport))
	if err != nil {
		panic(err)
	}

	sln, err := net.Listen("tcp", fmt.Sprintf(":%d", sport))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			ccn, err := cln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			cch <- ccn
		}
	}()

	go func() {
		for {
			scn, err := sln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			signatureBuf := make([]byte, 1024)
			n, err := scn.Read(signatureBuf)
			if err != nil {
				log.Println(err)
				continue
			}
			if n != (64+128)/8 {
				log.Println("signature len error")
				return
			}
			if err := butin.CheckSignature(signatureBuf[:n], secretKey); err != nil {
				log.Println(err)
				continue
			}

			sch <- scn
		}
	}()

	for {
		ccn := <-cch
		scn := <-sch

		if _, err := scn.Write([]byte(butin.TranDataKey)); err != nil {
			log.Println(err)
			continue
		}

		go func() {
			defer func() {
				ccn.Close()
				scn.Close()
			}()

			ch := make(chan struct{}, 2)
			go func() {
				if _, err := io.Copy(ccn, scn); err != nil {
					log.Println(err)
				}
				ch <- struct{}{}
			}()

			go func() {
				if _, err := io.Copy(scn, ccn); err != nil {
					log.Println(err)
				}
				ch <- struct{}{}
			}()

			<-ch
		}()
	}
}
