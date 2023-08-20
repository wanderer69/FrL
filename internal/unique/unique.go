package unique

import (
	"math/rand"
	"time"
)

var _unique map[string]string

func init() {
	InitUniqueValue()
}

func InitUniqueValue() {
	rand.Seed(time.Now().UTC().UnixNano())
	_unique = make(map[string]string)
}

func UniqueValue(len_n int) string {
	var str string
	for {
		var bytes_array []byte
		for i := 0; i < len_n; i++ {
			bytes := rand.Intn(35)
			if bytes > 9 {
				bytes = bytes + 7
			}
			bytes_array = append(bytes_array, byte(bytes+16*3))
		}
		str = string(bytes_array)
		if _, ok := _unique[str]; !ok {
			_unique[str] = str
			break
		}
	}
	return str
}
