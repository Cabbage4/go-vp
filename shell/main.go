package main

import (
	"flag"
	"fmt"
	"go-vp/butin"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

var (
	addr      string
	target    int
	secretKey string

	errorCount int
)

func main() {
	flag.StringVar(&addr, "addr", ":10022", "")
	flag.IntVar(&target, "target", 8081, "")
	flag.StringVar(&secretKey, "secretKey", "jerryzhuo@abcd", "")
	flag.Parse()

	g := new(sync.WaitGroup)

	worker := func() {
		defer g.Done()

		for {
			if errorCount > 10 {
				fmt.Println("errorCount > 10")
				return
			}

			cnn, err := net.Dial("tcp", addr)
			if err != nil {
				errorCount++
				fmt.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}

			signature, err := butin.GenSignature(time.Now().Unix(), secretKey)
			if err != nil {
				errorCount++
				fmt.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}

			cnn.Write(signature)
			errorCount = 0

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

					fmt.Printf("[%s->%s]link\n", cnn.RemoteAddr().String(), tCnn.RemoteAddr().String())
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
				fmt.Printf("[%s->%s]done\n", cnn.RemoteAddr().String(), tCnn.RemoteAddr().String())
				tCnn = nil
			}
		}
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		g.Add(1)
		go worker()
	}

	g.Wait()
}
