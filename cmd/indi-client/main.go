package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type BaseAttrs struct {
	Name  string `xml:"name,attr"`
	Label string `xml:"label,attr"`
}

type VectorAttrs struct {
	Device    string `xml:"device,attr"`
	Group     string `xml:"group,attr"`
	State     string `xml:"state,attr"`
	Perm      string `xml:"perm,attr"`
	Timeout   int    `xml:"timeout,attr"`
	Timestamp string `xml:"timestamp,attr"`
	BaseAttrs
}

type DefText struct {
	XMLName xml.Name `xml:"defText"`
	Value   string   `xml:",chardata"`
	BaseAttrs
}

type DefTextVector struct {
	XMLName xml.Name  `xml:"defTextVector"`
	DefText []DefText `xml:"defText"`
	VectorAttrs
}

type DefSwitch struct {
	XMLName xml.Name `xml:"defSwitch"`
	Value   string   `xml:",chardata"`
	BaseAttrs
}

type DefSwitchVector struct {
	XMLName   xml.Name    `xml:"defSwitchVector"`
	Rule      string      `xml:"rule,attr"`
	DefSwitch []DefSwitch `xml:"defSwitch"`
	VectorAttrs
}

type DefNumber struct {
	XMLName xml.Name `xml:"defNumber"`
	Format  string   `xml:"format,attr"`
	Min     string   `xml:"min,attr"`
	Max     string   `xml:"max,attr"`
	Step    string   `xml:"step,attr"`
	Value   int      `xml:",chardata"`
	BaseAttrs
}

type DefNumberVector struct {
	XMLName   xml.Name    `xml:"defNumberVector"`
	DefNumber []DefNumber `xml:"defNumber"`
	VectorAttrs
}

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

// Trimmer is used to remove blank space from received XML
type Trimmer struct {
	dec *xml.Decoder
}

func (tr Trimmer) Token() (xml.Token, error) {
	t, err := tr.dec.Token()
	if cd, ok := t.(xml.CharData); ok {
		t = xml.CharData(bytes.TrimSpace(cd))
	}
	return t, err
}

func recv(c net.Conn) {
	raw := xml.NewDecoder(c)
	decoder := xml.NewTokenDecoder(Trimmer{raw})
	properties := make([]interface{}, 0)

	for {
		t, err := decoder.Token()
		if t == nil {
			if err == nil {
				continue
			}
			if err == io.EOF {
				log.Println("EOF")
				break
			}
		}

		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "defNumberVector":
				var defNumberVector DefNumberVector
				decoder.DecodeElement(&defNumberVector, &se)
				properties = append(properties, defNumberVector)
			case "defSwitchVector":
				var defSwitchVector DefSwitchVector
				decoder.DecodeElement(&defSwitchVector, &se)
				properties = append(properties, defSwitchVector)
			case "defTextVector":
				var defTextVector DefTextVector
				decoder.DecodeElement(&defTextVector, &se)
				properties = append(properties, defTextVector)
			case "defNumber":
				var defNumber DefNumber
				decoder.DecodeElement(&defNumber, &se)
				properties = append(properties, defNumber)
			case "defSwitch":
				var defSwitch DefNumber
				decoder.DecodeElement(&defSwitch, &se)
				properties = append(properties, defSwitch)
			case "defText":
				var defText DefNumber
				decoder.DecodeElement(&defText, &se)
				properties = append(properties, defText)
			default:
				log.Printf("Unhandled data type: %s\n", se.Name.Local)
			}
		default:
		}

		// log.Printf("Containing %d properties", len(properties))
	}
}
