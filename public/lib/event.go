package frl

import (
	"fmt"
)

// событие - это механизм который позволяет
type Event struct {
	Type string // тип события
}

func NewEvent(t string) (*Event, error) {
	switch t {
	case "timer":
	default:
		return nil, fmt.Errorf("type %v not valid", t)
	}
	return &Event{Type: t}, nil
}
