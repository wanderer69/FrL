package frl

import (
	"fmt"
	//	"github.com/wanderer69/FrL/src/lib/common"
)

// метод это выполнимый тип данных
// скрипт - это нечто выполнимое - то есть по сути это объект который может проивзодить действия над фреймами

type Method struct {
	Name string   // имя метода
	Args []*Value // список имен переменных аргументов
	Body string   // тело скрипта на языке фреймов
}

func NewMethod(name string, b string) (*Method, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("empty body")
	}
	return &Method{Name: name, Body: b}, nil
}
