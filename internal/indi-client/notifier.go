package indiclient

type EventType uint8

const (
	Add EventType = iota
	Update
	Delete
)

type Event struct {
	EventType EventType
	Property  Property
}

type Observer interface {
	OnNotify(Event)
}

type Notifier interface {
	Register(Observer)
	Unregister(Observer)
	Notify(Event)
}
