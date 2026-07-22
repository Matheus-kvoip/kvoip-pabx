package events

// Event is an internal PBX domain event (call started, registered, etc).
type Event struct {
	Type string
	Data map[string]string
}

// Bus is a simple in-process pub/sub placeholder.
type Bus struct {
	subs []chan Event
}

func NewBus() *Bus {
	return &Bus{}
}

func (b *Bus) Publish(evt Event) {
	for _, ch := range b.subs {
		select {
		case ch <- evt:
		default:
		}
	}
}

func (b *Bus) Subscribe(buffer int) <-chan Event {
	ch := make(chan Event, buffer)
	b.subs = append(b.subs, ch)
	return ch
}
