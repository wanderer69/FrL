package frl

import "sync"

type Map struct {
	mu            sync.Mutex
	valueByString map[string]*Value
}

func NewMap() *Map {
	return &Map{
		valueByString: make(map[string]*Value),
	}
}
