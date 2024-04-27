package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

func main() {
	var host string

	flag.StringVar(&host, "host", "localhost:7624", "INDI server address")
	flag.Parse()

	log.Printf("Connecting to INDI server at %s\n", host)

	client, err := indiclient.New(host)
	if err != nil {
		log.Fatalf("could not create INDI client: %q", err)
	}

	log.Println("connected")

	exit := make(chan string)
	err = client.GetProperties()
	if err != nil {
		log.Fatalf("could not get INDI properties: %q", err)
	}

	for {
		select {
		// Wait forever until user kills the process
		case <-exit:
			os.Exit(0)
		case <-client.Data:
			for v := range client.Data {
				fmt.Printf("%+v\n", v)
			}
		}
	}
}
