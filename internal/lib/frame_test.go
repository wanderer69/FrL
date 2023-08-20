package frl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	print "github.com/wanderer69/FrL/internal/lib/print"
)

func TestFrame(t *testing.T) {
	buffer := ""
	clearBuffer := func() {
		buffer = ""
	}
	printFunc := func(frm string, args ...any) {
		str := fmt.Sprintf(frm, args...)
		buffer = buffer + str
	}
	output := print.NewOutput(printFunc)

	f := NewFrame()
	f.AddSlot("ID")
	f.Set("ID", 1)
	f.Print(output, true)
	require.Contains(t, buffer, "ID () 1")

	clearBuffer()
	f.AddSlot("отношение")
	f.Set("отношение", "value2")
	f.Print(output, true)
	require.Contains(t, buffer, "отношение () value2")

	clearBuffer()
	f.AddSlot("наименование")
	f.Set("наименование", "value1")
	f.Print(output, true)
	require.Contains(t, buffer, "наименование () value1")

	clearBuffer()
	f.AddSlot("имя")
	f.Set("имя", "value3")
	f.Print(output, true)
	require.Contains(t, buffer, "имя () value3")
}

func TestFrameGetSlotValue(t *testing.T) {
	buffer := ""
	clearBuffer := func() {
		buffer = ""
	}
	printFunc := func(frm string, args ...any) {
		str := fmt.Sprintf(frm, args...)
		buffer = buffer + str
	}
	output := print.NewOutput(printFunc)

	f := NewFrame()
	f.AddSlot("ID")
	f.Set("ID", 1)
	f.AddSlot("отношение")
	f.Set("отношение", "value2")
	f.AddSlot("наименование")
	f.Set("наименование", "value1")
	f.AddSlot("имя")
	f.Set("имя", "value3")

	clearBuffer()
	f1 := NewFrame()
	f1.AddSlot("ID")
	f1.Set("ID", 2)
	f1.AddSlot("slot4")
	f1.Set("slot4", []int{1, 2, 3})
	f.AddSlot("slot2")
	f.Set("slot2", f1)
	fmt.Printf("---\r\n")
	f.Print(output, true)
	fmt.Printf("-\r\n")
	require.Contains(t, buffer, "slot2 () slot2 {ID () 2, slot4 () [1 2 3]}")

	expectedNames := []string{"slot2", "ID", "отношение", "наименование", "имя"}
	expectedValues := []string{"ID () 2, slot4 () [1 2 3]", "1", "value2", "value1", "value3"}
	counter := 0
	ff := f.Iterate()
	for {
		s, ok, err := ff()
		if err != nil {
			break
		}
		ssl := s.GetSlotValue()
		for i := range ssl {
			ss, _ := FromType(ssl[i])
			fmt.Printf("k %v v %v\r\n", s.GetSlotName(), ss)
			require.Equal(t, expectedNames[counter], s.GetSlotName())
			require.Equal(t, expectedValues[counter], ss)
		}
		counter++
		if ok {
			break
		}
	}
}
