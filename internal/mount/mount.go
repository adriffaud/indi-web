package mount

import indiclient "github.com/adriffaud/indi-web/internal/indi-client"

type Mount struct {
	client    *indiclient.Client
	Connected bool
	Driver    string
	Parked    bool
	Tracking  bool
	RA        string
	DEC       string
}

func NewMount() Mount {
	return Mount{
		RA:        "00:00:00",
		DEC:       "00:00:00",
		Parked:    true,
		Tracking:  false,
		Connected: false,
	}
}

func (m Mount) Connect() {}

func (m Mount) Disconnect() {}

func (m Mount) Park() {}

func (m Mount) Unpark() {}

func (m Mount) StartTracking() {}

func (m Mount) StopTracking() {}
