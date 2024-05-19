package indiclient

import "encoding/xml"

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

type DelProperty struct {
	XMLName xml.Name `xml:"delProperty"`
	Device  string   `xml:"device,attr"`
	Name    string   `xml:"name,attr"`
}
