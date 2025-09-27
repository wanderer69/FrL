package frl

import (
	"fmt"
	"time"
)

// событие - это механизм который позволяет обрабатывать события
type Event struct {
	Type     string // тип события
	Duration time.Duration
	Channel  string
	Fn       string
	ID       string
}

type EventManager struct {
	events []*Event
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (em *EventManager) AddEvent(t string, d time.Duration, channel string, fn string) error {
	switch t {
	case "timer":

	case "channel":

	default:
		return fmt.Errorf("type %v not valid", t)
	}
	em.events = append(em.events, &Event{
		Type:     t,
		Duration: d,
		Channel:  channel,
		Fn:       fn,
	})
	return nil
}

func (em *EventManager) GetEvents() []*Event {
	return em.events
}
