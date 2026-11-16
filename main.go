package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/onewire/onewirereg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ds18b20"
	"periph.io/x/host/v3"
)

func main() {
	h, err := host.Init()
	if err != nil {
		panic(err)
	}
	fmt.Printf("loaded devices: %v\n", h.Loaded)
	fmt.Printf("skipped devices: %v\n", h.Skipped)
	fmt.Printf("failed devices: %v\n", h.Failed)

	fmt.Printf("one wire ifaces: %v\n", onewirereg.All())

	b, err := onewirereg.Open("") // open default 1-wire interface
	if err != nil {
		panic(err)
	}
	defer b.Close()

	addrs, err := b.Search(false) // search for all devices
	if err != nil {
		panic(err)
	}
	// TODO: assumes only one device on bus
	if len(addrs) < 1 {
		panic(fmt.Errorf("instruments: temperature sensor not found"))
	}
	addr := addrs[0]

	d, err := ds18b20.New(b, addr, 12) // connect using 12-bit resolution
	if err != nil {
		panic(err)
	}
	var temp physic.Env

	for {
		err = d.Sense(&temp) // blocking read temp
		if err != nil {
			panic(err)
		}
		fmt.Printf("temperature: %v C", temp.Temperature.Celsius())
		time.Sleep(1000)
	}

	// cmd := flag.Arg(0)
	// switch cmd {
	// case "serve":
	// 	serve()
	// case "cron":
	// 	cron()
	// default:
	// 	log.Fatalf("Command required. Supported commands: serve, cron")
	// }
}

// func serve() {
// 	log.Printf("starting server, env: %v\n", api.Env())
// 	db, err := database.Init()
// 	if err != nil {
// 		log.Fatalf("unable to start database: %v\n", err)
// 	}
// 	defer db.Close()

// 	s := server.NewServer(db)

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8000"
// 	}
// 	addr := "localhost:" + port
// 	fmt.Printf("Starting server on %v\n", addr)
// 	err = http.ListenAndServe(addr, s)
// 	if err != nil {
// 		if errors.Is(err, http.ErrServerClosed) {
// 			fmt.Printf("server shutdown\n")
// 		} else {
// 			log.Fatalf("unable to start server: %v\n", err)
// 		}
// 	}
// }

func cron() {
	// TODO
	// look at config for interval of different tasks
	// look for log for which tasks are overdue
	// execute all tasks
	// exit
}
