package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/rabidaudio/inshrimpentation/api"
)

func main() {
	log.Printf("starting server, env: %v\n", api.Env())

	var err error
	if err = syscall.Mkfifo("foo_pipe", 0666); err != nil {
		panic(err)
	}
	defer os.Remove("foo_pipe")

	ctx, cancel := context.WithCancel(context.Background())

	for i := range 3 {
		go func() {
			f, err := os.OpenFile("foo_pipe", os.O_RDONLY, os.ModeNamedPipe)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			fmt.Printf("start [%v]\n", i)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Printf("Recv [%v]: %v\n", i, line)
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("error [%v]: %v\n", i, err)
			}
		}()
	}

	// writer
	go func(ctx context.Context) {
		f, err := os.OpenFile("foo_pipe", os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		i := 0
		for {
			_, err := f.WriteString(fmt.Sprintf("hello %v\n", i))
			i += 1
			if err != nil {
				panic(err)
			}
			f.Sync()

			timeout := time.After(100 * time.Millisecond)
			select {
			case <-timeout:
				continue
			case <-ctx.Done():
				fmt.Println("context closed")
				return
			}
		}
	}(ctx)

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
