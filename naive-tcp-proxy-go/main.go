package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <listen-addr> <upstream-addr>", os.Args[0])
	}

	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var downstreamConns, upstreamConns uint64

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for range t.C {
			log.Printf("downstream = %d, upstream = %d",
				atomic.LoadUint64(&downstreamConns), atomic.LoadUint64(&upstreamConns))
		}
	}()

	for {
		downstream, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		start := time.Now()

		go func() {
			atomic.AddUint64(&downstreamConns, 1)
			defer atomic.AddUint64(&downstreamConns, ^uint64(0))

			upstream, err := net.DialTimeout("tcp", os.Args[2], 100*time.Millisecond)
			if err != nil {
				log.Println(err)
				return
			}
			atomic.AddUint64(&upstreamConns, 1)
			defer atomic.AddUint64(&upstreamConns, ^uint64(0))

			var once sync.Once
			cleanup := func() {
				downstream.Close()
				upstream.Close()
				log.Printf("duration = %s", time.Since(start))
			}

			go func() {
				defer once.Do(cleanup)
				io.CopyBuffer(downstream, upstream, nil)
			}()

			defer once.Do(cleanup)
			io.CopyBuffer(upstream, downstream, nil)
		}()
	}
}
