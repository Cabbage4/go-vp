package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
)

const count = 100

func TestGet(t *testing.T) {
	g := new(sync.WaitGroup)
	g.Add(count)

	f := func() {
		defer g.Done()

		rsp, err := http.Get("http://127.0.0.1:10021")
		if err != nil {
			log.Println(err)
			return
		}

		b, err := io.ReadAll(rsp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(b))
		rsp.Body.Close()
	}

	for i := 0; i < count; i++ {
		go f()
	}

	g.Wait()
}
