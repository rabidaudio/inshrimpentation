package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/rabidaudio/inshrimpentation/api"
)

func main() {
	log.Printf("starting server, env: %v\n", api.Env())

	var err error

	ctx, cancel := context.WithCancel(context.Background())

	p, err := net.ListenUnixgram("unixgram", &net.UnixAddr{Name: "test.sock", Net: "unixgram"})
	if err != nil {
		panic(err)
	}
	defer p.Close()
	defer os.Remove("test.sock")

	// writer
	go func() {
		i := 0
		mu := sync.Mutex{}
		clients := map[string]*net.UnixAddr{}

		go func() {
			buf := make([]byte, 1)
			for {
				_, addr, err := p.ReadFromUnix(buf)
				if errors.Is(err, net.ErrClosed) {
					return
				}
				if err != nil {
					panic(err)
				}
				mu.Lock()
				if buf[0] == 1 {
					// connect
					clients[addr.String()] = addr
				} else {
					// disconnect
					delete(clients, addr.String())
				}
				mu.Unlock()
			}
		}()

		for {
			mu.Lock()
			for _, c := range clients {
				_, err := p.WriteToUnix(fmt.Appendf(nil, "hello %v\n", i), c)
				if err != nil {
					delete(clients, c.String())
				}
			}
			mu.Unlock()
			i += 1

			timeout := time.After(500 * time.Millisecond)
			select {
			case <-timeout:
				continue
			case <-ctx.Done():
				fmt.Println("context closed")
				return
			}
		}
	}()

	for i := range 3 {
		go func(i int) {
			name := fmt.Sprintf("test_%v.sock", i)
			p, err := net.DialUnix("unixgram",
				&net.UnixAddr{Name: name, Net: "unixgram"},
				&net.UnixAddr{Name: "test.sock", Net: "unixgram"})
			if err != nil {
				panic(err)
			}
			defer p.Close()
			defer os.Remove(name)

			// send connect
			_, err = p.Write([]byte{1})
			if err != nil {
				panic(err)
			}
			// defer p.Write([]byte{0})

			fmt.Printf("start [%v]\n", i)

			scanner := bufio.NewScanner(p)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Printf("Recv [%v]: %v\n", i, line)
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("error [%v]: %v\n", i, err)
			}
		}(i)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)

	// db, err := database.Init()
	// if err != nil {
	// 	log.Fatalf("unable to start database: %v\n", err)
	// }
	// defer db.Close()

	// s := server.NewServer(db)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8000"
	// }
	// addr := "localhost:" + port
	// fmt.Printf("Starting server on %v\n", addr)
	// err = http.ListenAndServe(addr, s)
	// if err != nil {
	// 	if errors.Is(err, http.ErrServerClosed) {
	// 		fmt.Printf("server shutdown\n")
	// 	} else {
	// 		log.Fatalf("unable to start server: %v\n", err)
	// 	}
	// }
}
