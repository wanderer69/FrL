package frl

import (
	"errors"
	"fmt"

	print "github.com/wanderer69/FrL/internal/lib/print"
)

// тип слот
// фраме слот
type Slot struct {
	name     string
	value    []*Value
	property string
	next     *Slot
}

func NewSlot(name string) *Slot {
	return &Slot{name: name}
}

func (slot *Slot) Set(v *Value) {
	slot.value = append(slot.value, v)
}

func (slot *Slot) GetSlotName() string {
	return slot.name
}

func (slot *Slot) GetSlotValue() []*Value {
	return slot.value
}

func (slot *Slot) GetSlotProperty() string {
	return slot.property
}

// тип фрейм
// фраме
type Frame struct {
	ht *HashTable
}

func NewFrame() *Frame {
	return &Frame{NewHashTable(8, 16)}
}

func (f *Frame) AddSlot(name string) error {
	slot := Slot{name: name}
	return f.ht.Add(&slot)
}

func (f *Frame) AddSlotRaw(name string, prop string, value []*Value) {
	slot := Slot{name: name, property: prop, value: value}
	f.ht.Add(&slot)
}

func (f *Frame) DeleteSlot(name string) error {
	err := f.ht.Delete(name)
	if err != nil {
		return fmt.Errorf("undefined slot name %v", name)
	}
	return nil
}

func (f *Frame) Set(name string, value interface{}) (*Value, error) {
	s_, ok := f.ht.Lookup(name)
	if ok {
		slot := s_
		pv := &Value{}
		pv.Value(value)
		tt, ok := ToType(value)
		if ok {
			pv.Typev(tt)
		}
		slot.value = append(slot.value, pv)
		return pv, nil
	} else {
		return nil, fmt.Errorf("undefined slot name %v", name)
	}
}

func (f *Frame) SetValue(name string, value *Value) (*Value, error) {
	s_, ok := f.ht.Lookup(name)
	if ok {
		slot := s_
		slot.value = append(slot.value, value)
		return value, nil
	} else {
		return nil, fmt.Errorf("undefined slot name %v", name)
	}
}

func (f *Frame) SetValues(name string, values []*Value) error {
	s_, ok := f.ht.Lookup(name)
	if ok {
		slot := s_
		slot.value = values
		return nil
	} else {
		return fmt.Errorf("undefined slot name %v", name)
	}
}

func (f *Frame) SetSlotProperty(name string, property string) error {
	s_, ok := f.ht.Lookup(name)
	if ok {
		slot := s_
		slot.property = property
	} else {
		return fmt.Errorf("undefined slot name %v", name)
	}
	return nil
}

func (f *Frame) Append(name string, value interface{}) error {

	return nil
}

// Get return iterator by values array
func (f *Frame) Get(name string) (func() (*Value, error), error) {
	slot, ok := f.ht.Lookup(name)
	if ok {
		pos := 0
		r_func := func() (*Value, error) {
			ptr := pos
			if pos < len(slot.value) {
				pos = pos + 1
				return slot.value[ptr], nil
			}
			return nil, errors.New("empty")
		}
		return r_func, nil
	}
	return nil, fmt.Errorf("undefined slot name %v", name)
}

func (f *Frame) GetValue(name string) ([]*Value, error) {
	s_, ok := f.ht.Lookup(name)
	if ok {
		slot := s_
		return slot.value, nil
	} else {
		return nil, fmt.Errorf("undefined slot name %v", name)
	}
}

func (f *Frame) ToString() string {
	res := ""
	ff := f.ht.Iterate()
	flag := false
	for {
		slot, tt, err := ff()
		if err != nil {
			break
		}
		if flag {
			res = res + ", "
		} else {
			flag = true
		}
		res = res + fmt.Sprintf("%v (%v) ", slot.name, slot.property)
		for i := range slot.value {
			switch slot.value[i].GetType() {
			case VtBool:
				res = res + fmt.Sprintf("%v", slot.value[i].value)
			case VtInt:
				res = res + fmt.Sprintf("%v", slot.value[i].value)
			case VtFloat:
				res = res + fmt.Sprintf("%v", slot.value[i].value)
			case VtString:
				res = res + fmt.Sprintf("%v", slot.value[i].value)
			case VtFrame:
				ff := slot.value[i].value.(*Frame)
				res = res + fmt.Sprintf("%v {", slot.name)
				res = res + ff.ToString()
				res = res + "}"
			case VtSlice:
				res = res + "["
				vl := slot.value[i].value.([]*Value)
				for i := range vl {
					ss, ok := FromType(vl[i])
					if ok {
						if i == 0 {
							res = res + ss
						} else {
							res = res + ", " + ss
						}
					}
				}
				res = res + "]"
			case VtIterator:
				iter := slot.value[i].value.(*Iterator)
				res = res + fmt.Sprintf("iterator type %v pos %v", iter.typev, iter.pos)
			case VtSlot:
				slot := slot.value[i].value.(*Slot)
				res = res + fmt.Sprintf("slot %v property %v [", slot.GetSlotName(), slot.GetSlotProperty())
				vl := slot.GetSlotValue()
				for i := range vl {
					ss, ok := FromType(vl[i])
					if ok {
						if i == 0 {
							res = res + ss
						} else {
							res = res + ", " + ss
						}
					}
				}
				res = res + "]"
			default:
				res = res + fmt.Sprintf("%v", slot.value[i].value)
			}
		}
		if tt {
			break
		}
	}
	return res
}

func (f *Frame) Print(o *print.Output, flag_n bool) {
	ff := f.ht.Iterate()
	flag := false
	for {
		slot, tt, err := ff()
		if err != nil {
			break
		}
		if flag {
			o.Print(", ")
		} else {
			flag = true
		}

		o.Print("%v (%v) ", slot.name, slot.property)

		for i := range slot.value {
			switch slot.value[i].typev {
			case VtBool:
				o.Print("%v", slot.value[i].value)
			case VtInt:
				o.Print("%v", slot.value[i].value)
			case VtFloat:
				o.Print("%v", slot.value[i].value)
			case VtString:
				o.Print("%v", slot.value[i].value)
			case VtFrame:
				ff := slot.value[i].value.(*Frame)
				o.Print("%v {", slot.name)
				ff.Print(o, false)
				o.Print("}")
			case VtSlice:
				o.Print("[")
				vl := slot.value[i].value.([]*Value)
				for j := range vl {
					ss, ok := FromType(vl[j])
					if ok {
						if i == 0 {
							o.Print("%v", ss)
						} else {
							o.Print(", %v", ss)
						}
					}
				}
				o.Print("]")
			case VtIterator:
				iter := slot.value[i].value.(*Iterator)
				o.Print("iterator type %v pos %v", iter.typev, iter.pos)
			case VtSlot:
				slot := slot.value[i].value.(*Slot)
				o.Print("slot %v property %v [", slot.GetSlotName(), slot.GetSlotProperty())
				vl := slot.GetSlotValue()
				for j := range vl {
					ss, ok := FromType(vl[j])
					if ok {
						if i == 0 {
							o.Print("%v", ss)
						} else {
							o.Print("%v, ", ss)
						}
					}
				}
				o.Print("]")
			default:
				o.Print("%v", slot.value[i].value)
			}
		}
		if tt {
			break
		}
	}
	if flag_n {
		o.Print("\r\n")
	}
}

func (f *Frame) Iterate() func() (*Slot, bool, error) {
	ff := f.ht.Iterate()
	iterate := func() (*Slot, bool, error) {
		slot, flag, err := ff()
		if err != nil {
			return nil, true, errors.New("empty")
		}
		return slot, flag, nil
	}
	return iterate
}

// сопоставление фреймов f2 имеет незаполненные слоты (не имеющие значение)
// f2.Unify(f1) -> на выходе фрейм из незаполенных слотов
func (f1 *Frame) Unify(f2 *Frame) (*Frame, error) {
	f := NewFrame()

	return f, nil
}

// сопоставление группы фреймов
// UnifyAll(f1, f2, ...)
func UnifyAll(f1 *Frame, f2 ...*Frame) (*Frame, error) {
	f := NewFrame()

	return f, nil
}

func (f1 *Frame) Add(f2 *Frame) (*Frame, error) {
	f := NewFrame()

	return f, nil
}

func (f1 *Frame) Sub(f2 *Frame) (*Frame, error) {
	f := NewFrame()

	return f, nil
}

func (f1 *Frame) Compare(f2 *Frame) (*Frame, error) {
	f := NewFrame()

	return f, nil
}
