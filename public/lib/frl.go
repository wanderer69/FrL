package frl

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	fnc "github.com/wanderer69/FrL/public/functions"
	uqe "github.com/wanderer69/tools/unique"
)

func IsFrame(value interface{}) bool {
	// проверяем тип
	aa := reflect.ValueOf(value).Kind().String()
	if string(aa) == "ptr" {
		ss := reflect.TypeOf(value)
		vv := string(ss.String())[1:]
		sl := strings.Split(vv, ".")
		if len(sl) == 2 {
			if sl[1] == "Frame" {
				return true
			}
		}
	}
	return false
}

// тип итератор. замыкает значение списка либо фрейма и дальше позволяет вызовом функции next получать очередное значение
type Iterator struct {
	typev        ValueType
	value        *Value
	pos          int
	flag         bool
	func_iterate func() (*Value, error)
}

func (iter *Iterator) CreateIterate() func() (*Value, error) {
	switch iter.typev {
	case VtString:
		str := iter.value.GetValue().(string)
		iterate := func() (*Value, error) {
			ss := ""
			if iter.pos < len(str) {
				runeValue, width := utf8.DecodeRuneInString(str[iter.pos:])
				ss = string(runeValue)
				iter.pos = iter.pos + width
				if iter.pos < len(str) {
					iter.flag = false
				} else {
					iter.flag = true
				}
			} else {
				return nil, errors.New("empty")
			}
			vo := CreateValue(ss)
			return vo, nil
		}
		return iterate
	case VtFrame:
		f := iter.value.GetValue().(*Frame)
		ff := f.ht.Iterate()
		iterate := func() (*Value, error) {
			slot, flag, err := ff()
			iter.flag = flag
			if err != nil {
				return nil, errors.New("empty")
			}
			iter.flag = flag
			iter.pos = iter.pos + 1
			vo := CreateValue(slot)
			return vo, nil
		}
		return iterate
	case VtSlice:
		list := iter.value.GetValue().([]*Value)
		iterate := func() (*Value, error) {
			var vv *Value
			if iter.pos < len(list) {
				vv = list[iter.pos]
				iter.pos = iter.pos + 1
				if iter.pos < len(list) {
					iter.flag = false
				} else {
					iter.flag = true
				}
			} else {
				return nil, errors.New("empty")
			}
			return vv, nil
		}
		return iterate
	}
	return nil
}

func (iter *Iterator) Iterate() (*Value, error) {
	if iter.flag {
		return nil, errors.New("iterate end")
	}
	return iter.func_iterate()
}

func NewIterator(v *Value) (*Value, error) {
	// на входе счислимый тип строка, фрейм , список
	iter := &Iterator{}
	vt := v.GetType()
	switch vt {
	case VtString:
	case VtFrame:
	case VtSlice:
	default:
		return nil, fmt.Errorf("type %v not iterated", vt)
	}
	iter.typev = vt
	iter.value = v
	iter.pos = 0
	iter.flag = false

	iter.func_iterate = iter.CreateIterate()

	iv := &Value{}
	iv.value = iter
	tt, ok := ToType(iter)
	if ok {
		iv.typev = tt
	}
	return iv, nil
}

func (v *Value) Iterate() (vr *Value, err error) {
	// на входе счислимый тип строка, фрейм , список
	vt := v.GetType()
	switch vt {
	case VtIterator:
		iter := v.Iterator()
		vr, err = iter.Iterate()
	default:
		return nil, fmt.Errorf("type %v is not iterator", vt)
	}
	return
}

func (v *Value) IsEnd() (vr *Value, err error) {
	// на входе счислимый тип строка, фрейм , список
	vt := v.GetType()
	switch vt {
	case VtIterator:
		iter := v.Iterator()
		vr = CreateValue(iter.flag)
	default:
		return nil, fmt.Errorf("type %v is not iterator", vt)
	}
	return
}

func ToType(value interface{}) (ValueType, bool) {
	typev := ValueType(-1)
	res := false
	// проверяем тип
	if value == nil {
		typev = VtNil
		res = true
	} else {
		switch reflect.TypeOf(value).String() {
		case "bool":
			typev = VtBool
			res = true
		case "int":
			typev = VtInt
			res = true
		case "float64":
			typev = VtFloat
			res = true
		case "string":
			typev = VtString
			res = true
		default:
			aa := reflect.ValueOf(value).Kind().String()
			if string(aa) == "ptr" {
				ss := reflect.TypeOf(value)
				vv := string(ss.String())[1:]
				sl := strings.Split(vv, ".")
				if len(sl) == 2 {
					if sl[1] == "Frame" {
						typev = VtFrame
						res = true
					} else if sl[1] == "Iterator" {
						typev = VtIterator
						res = true
					} else if sl[1] == "Slot" {
						typev = VtSlot
						res = true
					} else if sl[1] == "Stream" {
						typev = VtStream
						res = true
					} else if sl[1] == "Function" {
						typev = VtFunction
						res = true
					} else if sl[1] == "Value" {
						typev = VtIterator
						res = true
					}
				}
			} else if string(aa) == "slice" {
				ss := reflect.TypeOf(value)
				vv := string(ss.String())
				sl := strings.Split(vv, ".")
				if len(sl) == 2 {
					if sl[1] == "Value" {
						typev = VtSlice
						res = true
					}
				}
			} else if string(aa) == "struct" {
				ss := reflect.TypeOf(value)
				vv := string(ss.String())[1:]
				sl := strings.Split(vv, ".")
				if len(sl) == 2 {
					if sl[1] == "Iterator" {
						typev = VtIterator
						res = true
					} else if sl[1] == "Slot" {
						typev = VtSlot
						res = true
					}
				}
			}
		}
	}
	return typev, res
}

func CompareValuesEq(value1 *Value, value2 *Value) bool {
	if value1.typev != value2.typev {
		return false
	}
	res := false

	// проверяем тип
	switch value1.typev {
	case VtNil:
		res = true
	case VtBool:
		if value1.value == value2.value {
			res = true
		}
	case VtInt:
		if value1.value == value2.value {
			res = true
		}
	case VtFloat:
		if value1.value == value2.value {
			res = true
		}
	case VtString:
		if value1.value == value2.value {
			res = true
		}
	case VtFrame:
	case VtSlice:
		v_1 := value1.value.([]*Value)
		v_2 := value2.value.([]*Value)
		if len(v_1) != len(v_2) {
			return false
		}
		flag := true
		for i := range v_1 {
			f := CompareValuesEq(v_1[i], v_2[i])
			if !f {
				break
			}
		}
		res = flag
	case VtIterator:
	case VtSlot:
		slot_1 := value1.value.(*Slot)
		slot_2 := value2.value.(*Slot)
		if slot_1.GetSlotName() != slot_2.GetSlotName() {
			return false
		}
		if slot_1.GetSlotProperty() != slot_2.GetSlotProperty() {
			return false
		}
		vl_1 := slot_1.GetSlotValue()
		vl_2 := slot_2.GetSlotValue()
		if len(vl_1) != len(vl_2) {
			return false
		}
		flag := true
		for i := range vl_1 {
			f := CompareValuesEq(vl_1[i], vl_2[i])
			if !f {
				break
			}
		}
		res = flag
	}
	return res
}

func CompareValuesLt(value1 *Value, value2 *Value) bool {
	if value1.typev != value2.typev {
		return false
	}
	res := false
	// проверяем тип
	switch value1.typev {
	case VtInt:
		if value1.value.(int) < value2.value.(int) {
			res = true
		}
	case VtFloat:
		if value1.value.(float64) < value2.value.(float64) {
			res = true
		}
	case VtString:
		if len(value1.value.(string)) < len(value2.value.(string)) {
			res = true
		}
	case VtFrame:
	case VtSlice:
		v_1 := value1.value.([]*Value)
		v_2 := value2.value.([]*Value)
		if len(v_1) < len(v_2) {
			res = true
		}
	case VtIterator:
	case VtSlot:
	}
	return res
}

func CompareValuesGt(value1 *Value, value2 *Value) bool {
	if value1.typev != value2.typev {
		return false
	}
	res := false
	// проверяем тип
	switch value1.typev {
	case VtInt:
		if value1.value.(int) > value2.value.(int) {
			res = true
		}
	case VtFloat:
		if value1.value.(float64) > value2.value.(float64) {
			res = true
		}
	case VtString:
		if len(value1.value.(string)) > len(value2.value.(string)) {
			res = true
		}
	case VtFrame:
	case VtSlice:
		v_1 := value1.value.([]*Value)
		v_2 := value2.value.([]*Value)
		if len(v_1) > len(v_2) {
			res = true
		}
	case VtIterator:
	case VtSlot:
	}
	return res
}

func FromType(value *Value) (string, bool) {
	str := ""
	res := false
	// проверяем тип
	switch value.typev {
	case VtNil:
		str = fmt.Sprintf("%v", value.value)
		res = true
	case VtBool:
		str = fmt.Sprintf("%v", value.value)
		res = true
	case VtInt:
		str = fmt.Sprintf("%v", value.value)
		res = true
	case VtFloat:
		str = fmt.Sprintf("%v", value.value)
		res = true
	case VtString:
		str = fmt.Sprintf("%v", value.value)
		res = true
	case VtFrame:
		str = value.value.(*Frame).ToString()
		res = true
	case VtSlice:
		str = "["
		vl := value.value.([]*Value)
		for i := range vl {
			ss, ok := FromType(vl[i])
			if ok {
				if i == 0 {
					str = str + ss
				} else {
					str = str + ", " + ss
				}
			}
		}
		str = str + "]"
		res = true
	case VtIterator:
		iter := value.value.(*Iterator)
		str = fmt.Sprintf("iterator type %v pos %v", iter.typev, iter.pos)
		res = true
	case VtSlot:
		slot := value.value.(*Slot)
		str = fmt.Sprintf("slot %v property %v ", slot.GetSlotName(), slot.GetSlotProperty())
		vl := slot.GetSlotValue()
		for i := range vl {
			ss, ok := FromType(vl[i])
			if i == 0 {
				str = str + ss
			} else {
				if ok {
					str = str + ", " + ss
				}
			}
		}
		res = true
	case VtStream:
		stream := value.value.(*Stream)
		str = fmt.Sprintf("stream %v", stream.SourceType)
		res = true
	case VtFunction:
		fn := value.value.(*fnc.Function)
		str = fmt.Sprintf("function %v args %v", fn.Name, fn.NumArgs)
		res = true
	}
	return str, res
}

type FrameEnvironment struct {
	FrameDict map[string][]*Frame // словарь фреймов в разных разрезах
	Frames    []*Frame            // список всех вреймов
}

func NewFrameEnvironment() *FrameEnvironment {
	fe := FrameEnvironment{}
	fe.FrameDict = make(map[string][]*Frame)
	// надо добавить фрейм с определением отношения
	f := NewFrame()
	// добавляем поле уникального идентификатора
	f.AddSlot("ID")
	id := uqe.UniqueValue(10)
	v, _ := f.Set("ID", id)
	fe.AddRelations(f, AddRelationItem{"frame", "", v})
	relation := "отношение"
	v, _ = f.Set("отношение", relation)
	fe.AddRelations(f, AddRelationItem{"relation", "отношение", v})

	return &fe
}

/*
	func isNil(i interface{}) bool {
		if i == nil {
			return true
		}
		switch reflect.TypeOf(i).Kind() {
		case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
			return reflect.ValueOf(i).IsNil()
		}
		return false
	}
*/
func (fe *FrameEnvironment) QueryRelations(qria ...QueryRelationItem) ([]*Frame, error) {
	// может быть три стратегии
	// 1 - когда есть запрос конкретного фрейма - тогда его обрабатываем в первую очередь
	// 2 - когда просто запросы отношений.
	// 3 - если попалось несколько запросов конкретного фрейма. - это скорее всего ошибка

	n := 0
	pos := -1
	// смотрим какие есть запросы
	for i := range qria {
		switch qria[i].ObjectType {
		case "frame":
			pos = i
			n = n + 1
		case "relation":
			// pos = i
		default:
			return nil, fmt.Errorf("bad common.QueryRelationItem ObjectType %v", qria[i].ObjectType)
		}
	}
	if n > 1 {
		return nil, fmt.Errorf("ObjectType frame more 1")
	}
	lst := []*Frame{}
	ok := false
	num := 0
	if n == 1 {
		// выполняем запрос на глобальном словаре
		mask := fmt.Sprintf("frame_%v", qria[pos].Value)
		lst, ok = fe.FrameDict[mask]
		if !ok {
			return lst, nil
		}
	} else {
		mask := fmt.Sprintf("relation_%v", qria[num].Object)
		lst_, ok_ := fe.FrameDict[mask]
		if !ok_ {
			return lst, nil
		}
		if qria[num].Value != nil {
			// проверяем
			for i := range lst_ {
				vala, err := lst_[i].Get(qria[num].Object)
				if err == nil {
					for {
						val, err := vala()
						if err != nil {
							break
						}
						// вот здесь надо правильно сравнивать
						if qria[num].Value.typev == val.typev {
							if qria[num].Value.value == val.value {
								lst = append(lst, lst_[i])
							}
						}
					}
				}
			}
		} else {
			lst = lst_
		}
		num = num + 1
	}
	// вот теперь можно оставшееся
	if num > pos {
		if (num + 1) > len(qria)-1 {
			return lst, nil
		}
	}

	lst_ := lst
	lst = []*Frame{}

	for {
		if num == pos {
			num = num + 1
		}

		if num > len(qria)-1 {
			break
		}

		for i := range lst_ {
			vala, err := lst_[i].Get(qria[num].Object)
			if err == nil {
				for {
					val, err := vala()
					if err != nil {
						break
					}
					if qria[num].Value != nil {
						// вот здесь надо правильно сравнивать
						if qria[num].Value.typev == val.typev {
							if qria[num].Value.value == val.value {
								lst = append(lst, lst_[i])
							}
						}
					} else {
						lst = append(lst, lst_[i])
					}
				}
			}
		}
		num = num + 1
	}
	return lst, nil
}

func (fe *FrameEnvironment) AddRelations(Frame *Frame, qri ...AddRelationItem) {
	maska := []string{}
	for i := range qri {
		mask := ""
		switch qri[i].ObjectType {
		case "frame":
			mask = fmt.Sprintf("frame_%v", qri[i].Value)
		case "relation":
			mask = fmt.Sprintf("relation_%v", qri[i].Object)
		default:
			panic("bad common.QueryRelationItem ObjectType")
		}
		maska = append(maska, mask)
	}

	for i := 0; i < len(maska); i++ {
		s := maska[i]
		ff := fe.FrameDict[s]
		fe.FrameDict[s] = append(ff, Frame)
	}
}

func (fe *FrameEnvironment) DeleteRelations(relation string) {
	mask := fmt.Sprintf("relation_%v", relation)
	delete(fe.FrameDict, mask)
}

func (fe *FrameEnvironment) DeleteFrameRelations(f *Frame) {
}

func (fe *FrameEnvironment) NewFrameWithRelation() *Frame {
	f := NewFrame()
	// добавляем поле уникального идентификатора
	f.AddSlot("ID")
	id := uqe.UniqueValue(10)
	v, _ := f.Set("ID", id)
	fe.AddRelations(f, AddRelationItem{"frame", "", v})
	fe.Frames = append(fe.Frames, f)
	return f
}

func (fe *FrameEnvironment) AddRelation(f1 *Frame, relation string) error {
	// проверяем что это отношение
	pv := &Value{}
	pv.value = relation
	tt, ok := ToType(relation)
	if ok {
		pv.typev = tt
	}
	_, err := fe.QueryRelations(QueryRelationItem{"relation", "отношение", pv})
	if err != nil {
		return fmt.Errorf("relation %v found in frame", relation)
	}
	// добавляем поле уникального идентификатора
	fe.AddRelations(f1, AddRelationItem{"relation", relation, nil})
	return nil
}

func (fe *FrameEnvironment) AddRelationWithoutCheck(f1 *Frame, relation string) error {
	// проверяем что это отношение
	pv := &Value{}
	pv.value = relation
	tt, ok := ToType(relation)
	if ok {
		pv.typev = tt
	}
	_, err := fe.QueryRelations(QueryRelationItem{"relation", "отношение", pv})
	if err != nil {
		return fmt.Errorf("relation %v found in frame", relation)
	}
	// добавляем поле уникального идентификатора
	fe.AddRelations(f1, AddRelationItem{"relation", relation, nil})
	return nil
}

func (fe *FrameEnvironment) DeleteRelation(f1 *Frame, relation string) error {
	// проверяем что этот слот есть
	_, err := f1.GetValue(relation)
	if err != nil {
		return fmt.Errorf("relation %v not found in frame", relation)
	}
	f1.DeleteSlot(relation)
	fe.DeleteRelations(relation)
	return nil
}
