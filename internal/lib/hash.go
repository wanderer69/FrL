package frl

import (
	"fmt"

	print "github.com/wanderer69/tools/parser/print"
)

type HTNode struct {
	collision uint16
	enable    bool
	slot      *Slot
}

const HASHTAB_SIZE = 71
const HASHTAB_MUL = 31

type HashTable struct {
	size  uint32
	mul   uint32
	table []HTNode
}

func NewHashTable(size uint32, mul uint32) *HashTable {
	ht := HashTable{size: size, mul: mul}
	ht.table = make([]HTNode, ht.size)
	for i := 0; i < (int)(ht.size); i++ {
		ht.table[i].enable = false
		ht.table[i].collision = 0
	}
	return &ht
}

func (ht *HashTable) Hash(key string) uint32 {
	var h uint32 = 0
	for _, p := range []byte(key) {
		h = h*ht.mul + (uint32)(p)
	}
	return h % ht.size
}

func (ht *HashTable) Add(slot *Slot) error {
	index := ht.Hash(slot.name)
	if ht.table[index].enable {
		ikv := ht.table[index].slot
		for {
			if ikv.next == nil {
				if ikv.name == slot.name {
					// slot exist
					return fmt.Errorf("slot %v exist", slot.name)
				} else {
					ikv.next = slot
				}
				break
			} else {
				if ikv.name == slot.name {
					return fmt.Errorf("slot %v exist", slot.name)
				} else {
					ikv = ikv.next
				}
			}
		}
		ht.table[index].collision = ht.table[index].collision + 1
	} else {
		ht.table[index].enable = true
		ht.table[index].slot = slot
	}
	return nil
}

func (ht *HashTable) Delete(name string) error {
	index := ht.Hash(name)
	if ht.table[index].enable {
		ikv := ht.table[index].slot
		var pkv *Slot = nil
		flag := false
		for {
			if ikv.next == nil {
				if ikv.name == name {
					if pkv != nil {
						pkv.next = nil
					}
					flag = true
				}
				break
			} else {
				if ikv.name == name {
					pkv.next = ikv.next
					flag = true
					break
				} else {
					pkv = ikv
					ikv = ikv.next
				}
			}
		}
		if !flag {
			return fmt.Errorf("unknown name %v", name)
		}
		return nil
	}
	return fmt.Errorf("unknown name %v", name)
}

func (ht *HashTable) Print(o *print.Output) {
	o.Print("size %v mul %v\r\n", ht.size, ht.mul)
	for i := 0; i < (int)(ht.size); i++ {
		o.Print("enable %v collision %v ", ht.table[i].enable, ht.table[i].collision)
		if ht.table[i].enable {
			ikv := ht.table[i].slot
			for {
				if ikv.next == nil {
					break
				} else {
					o.Print("%v %v ", ikv.name, ikv.value)
					ikv = ikv.next
				}
			}
		}
		o.Print("\r\n")
	}
}

func (ht *HashTable) Iterate() func() (*Slot, bool, error) {
	pos := 0
	state := 0
	var ikv *Slot
	iterate := func() (*Slot, bool, error) {
		var res *Slot
		next_state := 0
		pos_ := 0
		if pos < (int)(ht.size) {
			for {
				switch state {
				case 0:
					if pos < (int)(ht.size) {
						if ht.table[pos].enable {
							ikv = ht.table[pos].slot
							state = 1
						} else {
							pos = pos + 1
						}
					} else {
						return nil, true, fmt.Errorf("hash empty")
					}
				case 1:
					if ikv.next == nil {
						pos = pos + 1
						next_state = 0
						state = 2
						res = ikv
						pos_ = pos
					} else {
						ikv_ := ikv
						ikv = ikv.next
						res = ikv_
						if ikv != nil {
							state = 1
							return res, false, nil
						} else {
							next_state = 1
							state = 3
						}
					}
				case 2:
					if pos_ < (int)(ht.size) {
						if ht.table[pos_].enable {
							state = next_state
							return res, false, nil
						} else {
							pos_ = pos_ + 1
						}
					} else {
						return res, true, nil
					}
				case 3:
					if ikv.next == nil {
						state = 4
						pos_ = pos + 1
					} else {
						state = next_state
						return res, false, nil
					}
				case 4:
					if pos_ < (int)(ht.size) {
						if ht.table[pos_].enable {
							state = next_state
							return res, false, nil
						} else {
							pos_ = pos_ + 1
						}
					} else {
						return res, true, nil
					}
				}
			}
		}
		return nil, true, fmt.Errorf("hash empty")
	}
	return iterate
}

func (ht *HashTable) Lookup(name string) (*Slot, bool) {
	index := ht.Hash(name)
	if ht.table[index].enable {
		ikv := ht.table[index].slot
		for {
			if ikv.next == nil {
				if ikv.name == name {
					return ikv, true
				}
				break
			} else {
				if ikv.name == name {
					return ikv, true
				}
				ikv = ikv.next
			}
		}
	}
	return nil, false
}
