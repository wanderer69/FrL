package frl

import (
	"fmt"
)

// скрипт - это нечто выполнимое - то есть по сути это объект который может проивзодить действия над фреймами

type Script struct {
	Condition string   // условие срабатывания
	Args      []string // список имен переменных аргументов
	Body      string   // тело скрипта на языке фреймов
}

func NewScript(c string, b string) (*Script, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("empty body")
	}
	return &Script{Condition: c, Body: b}, nil
}
