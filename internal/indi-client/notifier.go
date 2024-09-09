package indiclient

type EventType uint8

const (
	Add EventType = iota
	Update
	Delete
	Message
	Timeout
)

type Event struct {
	EventType EventType
	Property  Property
	Message   string
}

type Observer interface {
	OnNotify(Event)
}

type Notifier interface {
	Register(Observer)
	Unregister(Observer)
	Notify(Event)
}
