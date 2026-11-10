package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rabidaudio/inshrimpentation/api"
	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/server"
)

func main() {
	cmd := flag.Arg(0)
	switch cmd {
	case "serve":
		serve()
	case "cron":
		cron()
	default:
		log.Fatalf("Command required. Supported commands: serve, cron")
	}
}

func serve() {
	log.Printf("starting server, env: %v\n", api.Env())
	db, err := database.Init()
	if err != nil {
		log.Fatalf("unable to start database: %v\n", err)
	}
	defer db.Close()

	s := server.NewServer(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := "localhost:" + port
	fmt.Printf("Starting server on %v\n", addr)
	err = http.ListenAndServe(addr, s)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server shutdown\n")
		} else {
			log.Fatalf("unable to start server: %v\n", err)
		}
	}
}

func cron() {
	// TODO
	// look at config for interval of different tasks
	// look for log for which tasks are overdue
	// execute all tasks
	// exit
}
