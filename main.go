package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {

	ctx, _ := context.WithCancel(context.Background())

	path := "test.mp4"

	// f, err := os.Create(path)
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	if err := syscall.Mkfifo(path, 0666); err != nil && !os.IsExist(err) {
		panic(err)
	}
	defer os.Remove(path)

	args := []string{
		"-f", "avfoundation",
		"-framerate", "30",
		"-video_size", "640x480",
		"-i", "0",
		"-pix_fmt", "uyvy422",
		"-b:v", "800k",
		"-t", "10",
		"-an", "-y",
		path,
	}

	log.Printf("ffmpeg %v\n", strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	done := make(chan bool)

	go func() {
		time.Sleep(5 * time.Second)
		cmd.Cancel()
		done <- true
	}()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	if err != nil {
		panic(err)
	}

	<-done

	f2, err := os.Create("test2.mp4")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	io.Copy(f2, &buf)

	// log.Printf("starting server, env: %v\n", api.Env())
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
