package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rabidaudio/inshrimpentation/api"
	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/server"
)

func main() {
	log.Printf("starting server, env: %v\n", api.Env())

	db, err := database.Init()
	if err != nil {
		log.Fatalf("unable to start database: %v\n", err)
	}
	defer db.Close()

	s := server.NewServer(db)
	log.Printf("starting server")
	err = http.ListenAndServe(":8000", s)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server shutdown\n")
		} else {
			log.Fatalf("unable to start server: %v\n", err)
		}
	}
}
