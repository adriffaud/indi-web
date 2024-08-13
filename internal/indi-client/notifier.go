package indiclient

type Event struct {
	Message string
}

type Observer interface {
	OnNotify(Event)
}

type Notifier interface {
	Register(Observer)
	Unregister(Observer)
	Notify(Event)
}
