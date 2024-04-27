package main

import (
	"fmt"
	"log"
	"os"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

const HOST = "localhost:7624"

func main() {
	log.Println("Connecting to INDI server")

	client, err := indiclient.New(HOST)
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
