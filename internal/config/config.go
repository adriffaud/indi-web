package config

type Mount struct {
	Connected bool
	Driver    string
	Parked    bool
	Tracking  bool
	RA        string
	DEC       string
}
