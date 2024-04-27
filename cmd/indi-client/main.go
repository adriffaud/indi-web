package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	log.Println("Connecting to INDI server")

	conn, err := net.Dial("tcp", "localhost:7624")
	if err != nil {
		log.Fatalf("could not create INDI client: %q", err)
	}
	defer conn.Close()

	log.Println("connected")

	exit := make(chan string)

	go recv(conn)

	_, err = fmt.Fprint(conn, "<getProperties version=\"1.7\"/>")
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	// Wait forever until user kills the process
	for {
		select {
		case <-exit:
			os.Exit(0)
		}
	}
}

func recv(c net.Conn) {
	decoder := xml.NewDecoder(c)

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}

		switch t.(type) {
		case xml.StartElement:
			log.Println("============================")
			log.Printf("%+v\n", t)
		case xml.EndElement:
		}
	}
}
