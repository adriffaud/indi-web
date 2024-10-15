package config

type Mount struct {
	Connected bool
	Driver    string
	Parked    bool
	Tracking  bool
	RA        string
	DEC       string
}

func NewMount() Mount {
	return Mount{RA: "00:00:00", DEC: "00:00:00", Parked: true}
}
