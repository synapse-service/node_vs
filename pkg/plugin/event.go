package plugin

import "github.com/242617/other/dispatcher"

const (
	EventTypeEvent = "event"
	EventTypeExit  = "exit"
)

type Event struct {
	Type string
}

func (p *Plugin) OnEvent(callback func(event Event)) {
	f := func(event dispatcher.Event) { callback(event.Value.(Event)) }
	p.EventDispatcher.AddEventListener(EventTypeEvent, &f)
}

func (p *Plugin) OnExit(callback func()) {
	f := func(event dispatcher.Event) { callback() }
	p.EventDispatcher.AddEventListener(EventTypeExit, &f)
}
