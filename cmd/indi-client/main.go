package main

import (
	"flag"
	"fmt"
	"log"

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

	err = client.GetProperties()
	if err != nil {
		log.Fatalf("could not get INDI properties: %q", err)
	}

	properties := make([]interface{}, 0)

	// Wait forever until user kills the process
	for {
		select {
		case <-client.Channel:
			v := <-client.Channel
			fmt.Printf("%T\n", v)
			properties = append(properties, v)
		}
	}
}
