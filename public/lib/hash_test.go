package frl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	uqe "github.com/wanderer69/tools/unique"
)

func TestHash(t *testing.T) {
	itemsCount := 16
	ht := NewHashTable(8, 16)
	slots := make([]*Slot, itemsCount)
	for i := 0; i < itemsCount; i++ {
		k := fmt.Sprintf("slotName%v", i)
		slots[i] = NewSlot(k)
	}
	for i := 0; i < itemsCount; i++ {
		k := uqe.UniqueValue(1)
		v := NewValue(-1, k)
		s := slots[i]
		s.Set(v)
		ht.Add(s)
	}
	expectedHT := &HashTable{
		size: 0x8,
		mul:  0x10,
		table: []HTNode{
			{
				collision: 0x2,
				enable:    true,
				slot:      slots[0],
			},
			{
				collision: 0x2,
				enable:    true,
				slot:      slots[1],
			},
			{
				collision: 0x1,
				enable:    true,
				slot:      slots[2],
			},
			{
				collision: 0x1,
				enable:    true,
				slot:      slots[3],
			},
			{
				collision: 0x1,
				enable:    true,
				slot:      slots[4],
			},
			{
				collision: 0x1,
				enable:    true,
				slot:      slots[5],
			},
			{
				collision: 0x0,
				enable:    true,
				slot:      slots[6],
			},
			{
				collision: 0x0,
				enable:    true,
				slot:      slots[7],
			},
		},
	}
	require.Equal(t, expectedHT, ht)
}
