package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rabidaudio/api/env"
	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/instruments"
	"github.com/rabidaudio/inshrimpentation/server"
)

func main() {
	var cmd string
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "serve":
		serve()
	case "measure":
		measure()
	case "cron":
		cron()
	default:
		log.Fatalf("Command required. Supported commands: serve, cron")
	}
}

func serve() {
	log.Printf("starting server, env: %v\n", env.Env())
	db, err := database.Init()
	if err != nil {
		log.Fatalf("unable to start database: %v\n", err)
	}
	defer db.Close()

	s := server.NewServer(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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

func measure() {
	ctx := context.Background()

	mm, err := instruments.Measure(ctx)
	if err != nil {
		log.Fatalf("measure: %v\n", err)
	}
	for _, m := range mm {
		log.Printf("[%v] %v: %v\n", m.Time, m.MeasurementType, m.Formatted())
	}
}

func cron() {
	// TODO
	// look at config for interval of different tasks
	// look for log for which tasks are overdue
	// execute all tasks
	// exit
}
