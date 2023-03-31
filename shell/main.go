package main

import (
	"flag"
	"fmt"
	"go-vp/butin"
	"io"
	"net"
	"runtime"
	"time"
)

var (
	addr   string
	target int
)

func main() {
	flag.StringVar(&addr, "addr", ":10022", "")
	flag.IntVar(&target, "target", 8081, "")
	flag.Parse()

	worker := func() {
		for {
			cnn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var tCnn net.Conn
			buf := make([]byte, 1024)
			for {
				n, err := cnn.Read(buf)
				if err != nil {
					if !butin.IsSkipError(err) {
						fmt.Println(err)
					}
					break
				}

				if n == 0 {
					continue
				}

				if tCnn == nil {
					tCnn, err = net.Dial("tcp", fmt.Sprintf(":%d", target))
					if err != nil {
						fmt.Println(err)
						break
					}

					go io.Copy(cnn, tCnn)
				}

				if _, err := tCnn.Write(buf[:n]); err != nil {
					fmt.Println(err)
					break
				}
			}

			cnn.Close()
			if tCnn != nil {
				tCnn.Close()
			}
		}
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go worker()
	}

	for {
		<-time.After(1 * time.Minute)
	}
}
