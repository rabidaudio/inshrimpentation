package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rabidaudio/api/env"
	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/instruments"
	"github.com/rabidaudio/inshrimpentation/server"
	"github.com/rabidaudio/inshrimpentation/tasks"
)

func main() {
	var cmd string
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db, err := database.Init()
	if err != nil {
		log.Fatalf("unable to start database: %v\n", err)
	}
	defer db.Close()
	dbctx := database.WithDB(ctx, db)

	switch cmd {
	case "measure":
		measure(ctx)
	case "configure":
		if len(os.Args) < 3 {
			log.Fatalf("path required")
		}
		path := os.Args[2]
		err = configure(dbctx, path)
		if err != nil {
			log.Fatalf("configure: %v\n", err)
		}
	case "check":
		err := check(dbctx)
		if err != nil {
			log.Fatalf("configure: %v\n", err)
		}

	default:
		cctx, cancel := context.WithCancelCause(dbctx)

		go func() {
			err := monitor(cctx)
			if err != nil {
				log.Printf("monitor: problem running tasks: %v\n", err)
			}
			cancel(err) // interrupt the server
		}()
		serve(cctx)
	}
}

func serve(ctx context.Context) {
	log.Printf("server: starting, env: %v\n", env.Env())
	db := database.GetDB(ctx)

	app := server.NewServer(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	s := http.Server{Addr: addr, Handler: app}

	esrv := make(chan error)
	go func() {
		fmt.Printf("server: starting on %v\n", addr)
		err := s.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("server: shutdown\n")
			} else {
				esrv <- err
			}
		}
		close(esrv)
	}()

	select {
	case err := <-esrv:
		if err != nil {
			log.Fatalf("server: error %v", err)
		}
	case <-ctx.Done():
		// Give in-flight requests a bounded window to finish.
		sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Shutdown(sctx); err != nil {
			fmt.Println("server: graceful shutdown failed, forcing close:", err)
			_ = s.Close()
			os.Exit(1)
		}
	}
}

// periodically check the status of the tank
func monitor(ctx context.Context) error {
	log.Printf("monitor: starting\n")
	for {
		check(ctx)

		time.Sleep(5 * time.Minute) // TODO: what's the lowest interval?
		// TODO: should we read the database and schedule? what if the configurations change?
	}
}

func check(ctx context.Context) error {
	log.Printf("monitor: running tasks\n")
	_, alerts, err := tasks.RunAllTasks(ctx, false)
	if err != nil {
		return err
	}
	for _, a := range alerts {
		log.Printf("alert: %v: %v\n", a.Type, a.Message)
	}
	return nil
}

func measure(ctx context.Context) {
	mm, err := instruments.Measure(ctx)
	if err != nil {
		log.Fatalf("measure: %v\n", err)
	}
	for _, m := range mm {
		log.Printf("[%v] %v: %v\n", m.Time, m.Type, m.Formatted())
	}
}
