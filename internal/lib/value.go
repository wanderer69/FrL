package frl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"

	"github.com/wanderer69/debug"

	uuid "github.com/satori/go.uuid"
	fnc "github.com/wanderer69/FrL/internal/lib/functions"
	print "github.com/wanderer69/FrL/internal/lib/print"
)

const (
	VtNil      = 0
	VtBool     = 1
	VtInt      = 2
	VtFloat    = 3
	VtString   = 4
	VtFrame    = 5
	VtSlice    = 6
	VtIterator = 7
	VtSlot     = 8
	VtStream   = 9
	VtFunction = 10
)

// значение слота и базовый тип интерпретатора
// значение слота
type Value struct {
	typev int
	value interface{}
}

func (v *Value) Value(value interface{}) {
	v.value = value
}

func (v *Value) Typev(t int) {
	v.typev = t
}

func (value *Value) GetType() int {
	return value.typev
}

func NewValue(vtype int, value interface{}) *Value {
	return &Value{vtype, value}
}

func CreateValue(value interface{}) *Value {
	pv := &Value{}
	pv.value = value
	tt, ok := ToType(value)
	if ok {
		pv.typev = tt
	}
	return pv
}

func RestoreValue(vtype int, value string) *Value {
	result := &Value{}
	result.typev = vtype
	switch vtype {
	case VtBool:
		i, _ := strconv.ParseInt(value, 10, 64)
		if value == "true" {
			result.value = true
		} else {
			result.value = false
		}
		result.value = int(i)
	case VtInt:
		i, _ := strconv.ParseInt(value, 10, 64)
		result.value = int(i)
	case VtFloat:
		f, _ := strconv.ParseFloat(value, 64)
		result.value = f
	case VtString:
		result.value = value
	case VtFrame:
		result.value = value
	case VtSlice:
		result.value = value
	case VtIterator:
		result.value = value
	case VtSlot:
		result.value = value
	case VtStream:
		result.value = value
	case VtFunction:
		result.value = value
	}

	return result
}

func (v *Value) GetValue() interface{} {
	return v.value
}

func (value *Value) Nil() *bool {
	switch value.typev {
	case VtNil:
		return nil
	}
	a := new(bool)
	return a
}

func (value *Value) Bool() bool {
	switch value.typev {
	case VtBool:
		return value.value.(bool)
	}
	return false
}

func (value *Value) Int() int {
	switch value.typev {
	case VtInt:
		return value.value.(int)
	}
	return 0
}

func (value *Value) Float() float64 {
	switch value.typev {
	case VtFloat:
		return value.value.(float64)
	}
	return math.NaN()
}

func (value *Value) String() string {
	switch value.typev {
	case VtString:
		return value.value.(string)
	}
	return ""
}

func (value *Value) Frame() *Frame {
	switch value.typev {
	case VtFrame:
		return value.value.(*Frame)
	}
	return nil
}

func (value *Value) Slice() []*Value {
	switch value.typev {
	case VtSlice:
		return value.value.([]*Value)
	}
	return nil
}

func (value *Value) Iterator() *Iterator {
	switch value.typev {
	case VtIterator:
		return value.value.(*Iterator)
	}
	return nil
}

func (value *Value) Slot() *Slot {
	switch value.typev {
	case VtSlot:
		return value.value.(*Slot)
	}
	return nil
}

func (value *Value) Stream() *Stream {
	switch value.typev {
	case VtStream:
		return value.value.(*Stream)
	}
	return nil
}

func (value *Value) Function() *fnc.Function {
	switch value.typev {
	case VtFunction:
		return value.value.(*fnc.Function)
	}
	return nil
}

func (value *Value) IsType() *Value {
	v := CreateValue(value.typev)
	return v
}

/*
func (v *Value) GetValue() interface{} {
	switch value.typev {
	case VtInt:
		return value.value.(int)
	case VtFloat:
		return value.value.(int)
	case VtString:
		return value.value.(int)
	case VtFrame:
		return value.value.(*Frame)
	case VtSlice:
		return value.value.([]*Value)
	case VtIterator:
		return value.value.(*Iterator)
	case VtSlot:
		return value.value.(*Slot)
	default:
		return nil
	}
}
*/

func (value *Value) Print(o *print.Output) error {
	// проверяем тип
	switch value.typev {
	case VtBool:
		o.Print("%v\r\n", value.value)
	case VtInt:
		o.Print("%v\r\n", value.value)
	case VtFloat:
		o.Print("%v\r\n", value.value)
	case VtString:
		o.Print("%v\r\n", value.value)
	case VtFrame:
		str := value.value.(*Frame).ToString()
		o.Print("%v\r\n", str)
	case VtSlice:
		str := ""
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
		o.Print("%v\r\n", str)
	case VtIterator:
		iter := value.value.(*Iterator)
		str := fmt.Sprintf("iterator type %v pos %v", iter.typev, iter.pos)
		o.Print("%v\r\n", str)
	case VtSlot:
		slot := value.value.(*Slot)
		str := fmt.Sprintf("slot %v property %v ", slot.GetSlotName(), slot.GetSlotProperty())
		vl := slot.GetSlotValue()
		for i := range vl {
			ss, ok := FromType(vl[i])
			if ok {
				str = str + ", " + ss
			}
		}
		o.Print("%v\r\n", str)
	case VtStream:
	case VtFunction:
		fnc := value.value.(*fnc.Function)
		str := fmt.Sprintf("function %v num args %v ", fnc.Name, fnc.NumArgs)
		o.Print("%v\r\n", str)
	default:
		return fmt.Errorf("error type %v", value.typev)
	}
	return nil
}

func (value1 *Value) Add(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtInt:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(int) + value2.value.(int))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(int) + int(value2.value.(float64)))
			return v, nil
		}
	case VtFloat:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(float64) + float64(value2.value.(int)))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(float64) + value2.value.(float64))
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) Sub(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtInt:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(int) - value2.value.(int))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(int) - int(value2.value.(float64)))
			return v, nil
		}
	case VtFloat:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(float64) - float64(value2.value.(int)))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(float64) - value2.value.(float64))
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) Mul(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtInt:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(int) * value2.value.(int))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(int) * int(value2.value.(float64)))
			return v, nil
		}
	case VtFloat:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(float64) * float64(value2.value.(int)))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(float64) * value2.value.(float64))
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) Div(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtInt:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(int) / value2.value.(int))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(int) / int(value2.value.(float64)))
			return v, nil
		}
	case VtFloat:
		switch value2.typev {
		case VtInt:
			v := CreateValue(value1.value.(float64) / float64(value2.value.(int)))
			return v, nil
		case VtFloat:
			v := CreateValue(value1.value.(float64) / value2.value.(float64))
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) FromString(value2 *Value) error {
	switch value1.typev {
	case VtInt:
		switch value2.typev {
		case VtString:
			v, err := strconv.ParseInt(value2.value.(string), 10, 64)
			if err != nil {
				return err
			}
			value1.value = int(v)
			return nil
		}
	case VtFloat:
		switch value2.typev {
		case VtString:
			v, err := strconv.ParseFloat(value2.value.(string), 64)
			if err != nil {
				return err
			}
			value1.value = v
			return nil
		}
	}
	return nil
}

func (value1 *Value) Concat(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		switch value2.typev {
		case VtString:
			v := CreateValue(value1.value.(string) + value2.value.(string))
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) SliceString(value2 *Value, value3 *Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		switch value2.typev {
		case VtInt:
			switch value3.typev {
			case VtInt:
				s := value1.value.(string)
				l := value2.value.(int)
				r := value3.value.(int)
				v := CreateValue(s[l:r])
				return v, nil
			}
		}
	}
	return nil, nil
}

func (value1 *Value) SprintfString(values ...*Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		s := value1.value.(string)
		sl := strings.Split(s, "%")
		cnt := 0
		ss := ""
		for i := range sl {
			if sl[i][0] == '?' {
				s1 := sl[i][1:]
				if cnt < len(values) {
					v := values[cnt]
					str, _ := FromType(v)
					ss = ss + str
					cnt = cnt + 1
				}
				ss = ss + s1
			} else if sl[i][0] == '%' {
				ss = ss + sl[i]
			} else {
				ss = ss + sl[i]
			}
		}
		ssl := strings.Split(ss, "\\")
		sss := ""
		for i := range ssl {
			if ssl[i][0] == 'r' {
				sss = sss + string([]byte{0x0d}) + ssl[i][1:]
			} else if ssl[i][0] == 'n' {
				sss = sss + string([]byte{0x0a}) + ssl[i][1:]
			} else if ssl[i][0] == 't' {
				sss = sss + string([]byte{0x09}) + ssl[i][1:]
			} else if ssl[i][0] == '\\' {
				sss = sss + string([]byte{0x5c}) + ssl[i][1:]
			} else {
				sss = sss + ssl[i]
			}
		}
		v := CreateValue(sss)
		return v, nil
	}
	return nil, nil
}

func (value1 *Value) Trim(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		switch value2.typev {
		case VtString:
			s := strings.Trim(value1.value.(string), value2.value.(string))
			v := CreateValue(s)
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) Split(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		switch value2.typev {
		case VtString:
			s := strings.Split(value1.value.(string), value2.value.(string))
			va := []*Value{}
			for i := range s {
				v := CreateValue(s[i])
				va = append(va, v)
			}
			v := CreateValue(va)
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) FromNumber(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		switch value2.typev {
		case VtInt:
			s := fmt.Sprintf("%v", value2.value.(int))
			v := CreateValue(s)
			return v, nil
		case VtFloat:
			s := fmt.Sprintf("%v", value2.value.(float64))
			v := CreateValue(s)
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) SlotGetName() (*Value, error) {
	switch value1.typev {
	case VtSlot:
		s := value1.value.(*Slot).GetSlotName()
		v := CreateValue(s)
		return v, nil
	}
	return nil, nil
}

func (value1 *Value) SlotGetValue() (*Value, error) {
	switch value1.typev {
	case VtSlot:
		s := value1.value.(*Slot).GetSlotValue()
		v := CreateValue(s)
		return v, nil
	}
	return nil, nil
}

func (value1 *Value) SlotGetProperty() (*Value, error) {
	switch value1.typev {
	case VtSlot:
		s := value1.value.(*Slot).GetSlotProperty()
		v := CreateValue(s)
		return v, nil
	}
	return nil, nil
}

func (value1 *Value) FrameAddSlot(value2 *Value) error {
	switch value1.typev {
	case VtFrame:
		switch value2.typev {
		case VtString:
			f := value1.value.(*Frame)
			n := value2.value.(string)
			f.AddSlot(n)
			return nil
		}
	}
	return nil
}

func (value1 *Value) FrameSetSlot(value2 *Value, value3 *Value) error {
	switch value1.typev {
	case VtFrame:
		switch value1.typev {
		case VtString:
			f := value1.value.(*Frame)
			n := value1.value.(string)
			f.SetValue(n, value3)
			return nil
		}
	}
	return nil
}

func (value1 *Value) FrameDeleteSlot(value2 *Value) error {
	switch value1.typev {
	case VtFrame:
		switch value2.typev {
		case VtString:
			f := value1.value.(*Frame)
			n := value2.value.(string)
			f.DeleteSlot(n)
			return nil
		}
	}
	return nil
}

func (value1 *Value) SliceItem(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtSlice:
		switch value2.typev {
		case VtInt:
			s := value1.value.([]*Value)
			ii := value2.value.(int)
			v := s[ii]
			return v, nil
		}
	}
	return nil, nil
}

func (value1 *Value) SliceSlice(value2 *Value, value3 *Value) (*Value, error) {
	switch value1.typev {
	case VtSlice:
		switch value2.typev {
		case VtInt:
			switch value3.typev {
			case VtInt:
				s := value1.value.([]*Value)
				i := value2.value.(int)
				j := value3.value.(int)
				vv := s[i:j]
				v := CreateValue(vv)
				return v, nil
			}
		}
	}
	return nil, nil
}

func (value1 *Value) SliceInsert(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtSlice:
		s := value1.value.([]*Value)
		value1.value = append([]*Value{value2}, s...)
		return value1, nil
	}
	return nil, nil
}

func (value1 *Value) SliceAppend(value2 *Value) (*Value, error) {
	switch value1.typev {
	case VtSlice:
		s := value1.value.([]*Value)
		value1.value = append(s, value2)
		return value1, nil
	}
	return nil, nil
}

func NewSlice(values ...*Value) (*Value, error) {
	v := CreateValue(values)
	return v, nil
}

func StreamCreate(value1 *Value) (*Value, error) {
	switch value1.typev {
	case VtString:
		uri := value1.value.(string)
		s, err := NewStream(uri)
		if err != nil {
			return nil, err
		}
		v := CreateValue(s)
		return v, nil
	}
	return nil, nil
}

func (value1 *Value) StreamOpen() error {
	switch value1.typev {
	case VtStream:
		s := value1.value.(*Stream)
		err := s.Open("")
		return err
	}
	return nil
}

func (value1 *Value) StreamRead() (*Value, *Value, error) {
	switch value1.typev {
	case VtStream:
		s := value1.value.(*Stream)
		cnt, data, err := s.Read()
		if err != nil {
			return nil, nil, err
		}
		v_cnt := CreateValue(cnt)
		v_data := CreateValue(string(data))
		return v_cnt, v_data, err
	}
	return nil, nil, nil
}

func (value1 *Value) StreamWrite(v_data *Value) error {
	switch value1.typev {
	case VtStream:
		s := value1.value.(*Stream)
		switch v_data.typev {
		case VtString:
			data := v_data.value.(string)
			err := s.Write([]byte(data))
			return err
		}
	}
	return nil
}

func (value1 *Value) StreamClose() error {
	switch value1.typev {
	case VtStream:
		s := value1.value.(*Stream)
		err := s.Close()
		return err
	}
	return nil
}

func (value1 *Value) StreamControlSet(v_cmd *Value, v_data *Value) error {
	switch value1.typev {
	case VtStream:
		s := value1.value.(*Stream)
		switch v_cmd.typev {
		case VtString:
			cmd := v_cmd.value.(string)
			switch v_data.typev {
			case VtString:
				data_str := v_data.value.(string)
				data := make(map[string]string)
				data_str_l := strings.Split(data_str, ";")
				for i := range data_str_l {
					ss := strings.Trim(data_str_l[i], " \r\n\t;")
					if len(ss) > 0 {
						ssl := strings.Split(ss, "=")
						if len(ssl) != 2 {
							return fmt.Errorf("separator not found in %v", ss)
						}
						data[ssl[0]] = ssl[1]
					}
				}
				err := s.ControlSet(cmd, data)
				return err
			}
		}
	}
	return nil
}

func (value1 *Value) StreamControlGet(v_cmd *Value) (*Value, error) {
	switch value1.typev {
	case VtStream:
		s := value1.value.(*Stream)
		switch v_cmd.typev {
		case VtString:
			cmd := v_cmd.value.(string)
			data, err := s.ControlGet(cmd)
			if err != nil {
				return nil, err
			}
			v_data := CreateValue(data)
			return v_data, nil
		}
	}
	return nil, nil
}

func (value1 *Value) EvalString(ie *InterpreterEnv) (*Value, error) {
	switch value1.typev {
	case VtString:
		s := value1.value.(string)
		ce := ie.contextEnv[len(ie.contextEnv)-1]
		cf := ce.current

		fvl, err := ie.TranslateText(cf.function.Name, s, 0, ie.Output)
		if err != nil {
			return nil, err
		}
		v_data := CreateValue(fvl)
		return v_data, nil
	}
	return nil, nil
}

func UUIDString() *Value {
	data := uuid.NewV4().String()
	v_data := CreateValue(data)
	return v_data
}

type Lenght_header struct {
	LenValue int32
}

func SaveLenghtValue(bb []byte) []byte {
	lenghtHeader := Lenght_header{LenValue: int32(len(bb))}

	len_all := (int)(unsafe.Sizeof(lenghtHeader)) + len(bb)
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &lenghtHeader); err != nil {
		fmt.Println(err)
	}
	if lenghtHeader.LenValue > 0 {
		if err := binary.Write(buf, binary.LittleEndian, bb); err != nil {
			fmt.Println(err)
		}
	}
	return buf.Bytes()
}

func SaveLenght(n int) []byte {
	lenghtHeader := Lenght_header{LenValue: int32(n)}

	len_all := (int)(unsafe.Sizeof(lenghtHeader))
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &lenghtHeader); err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

type ValueStoreHeader struct {
	Len  int32
	Type int32
	LenD int32
}

func SaveValueStore(v *Value) ([]byte, error) {
	// для длинных полей отдельная TLV с указанием поля
	debug.Alias("value.SaveValueStore").Printf("Value_store\r\n")
	debug.Alias("value.SaveValueStore").Printf("%#v\r\n", *v)
	valueStoreHeader := ValueStoreHeader{Type: int32(v.typev)}

	var bb []byte
	len_all := (int)(unsafe.Sizeof(valueStoreHeader))
	switch v.typev {
	case VtBool:
		vb := v.value.(bool)
		b_in := make([]byte, 0, 1)
		var buf = bytes.NewBuffer(b_in)
		if err := binary.Write(buf, binary.LittleEndian, &vb); err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = buf.Bytes()
	case VtInt:
		// value_store_header.
		value_int := int64(v.value.(int))
		b_in := make([]byte, 0, 8)
		var buf = bytes.NewBuffer(b_in)
		if err := binary.Write(buf, binary.LittleEndian, &value_int); err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = buf.Bytes()
	case VtFloat:
		value_float := float64(v.value.(float64))
		b_in := make([]byte, 0, 8)
		var buf = bytes.NewBuffer(b_in)
		if err := binary.Write(buf, binary.LittleEndian, &value_float); err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = buf.Bytes()
	case VtString:
		bb = SaveLenghtValue([]byte(v.value.(string)))
		len_all = len_all + len(bb)
	case VtFrame:
		f := v.value.(*Frame)
		ff := f.ht.Iterate()
		l := 0
		for {
			slot, tt, err := ff()
			if err != nil {
				break
			}
			bb_ := SaveLenghtValue([]byte(slot.GetSlotName()))
			bb = append(bb, bb_...)
			bb_ = SaveLenghtValue([]byte(slot.GetSlotProperty()))
			bb = append(bb, bb_...)

			vl := slot.GetSlotValue()
			bb_ = SaveLenght(len(vl))
			bb = append(bb, bb_...)
			for i := range vl {
				h, err := SaveValueStore(vl[i])
				if err != nil {
					return []byte{}, err
				}
				bb_ := SaveLenghtValue([]byte(h))
				bb = append(bb, bb_...)
			}
			l = l + 1
			if tt {
				break
			}
		}
		len_all = len_all + len(bb)
		valueStoreHeader.LenD = int32(l)
	case VtSlice:
		sl := v.value.([]*Value)
		for i := range sl {
			csh := sl[i]
			h, err := SaveValueStore(csh)
			if err != nil {
				return []byte{}, err
			}
			bb_ := SaveLenghtValue([]byte(h))
			bb = append(bb, bb_...)
		}
		len_all = len_all + len(bb)
		valueStoreHeader.LenD = int32(len(v.value.([]*Value)))
	case VtFunction:
		fn := v.value.(*fnc.Function)
		bb_, err := fnc.Func2Bin(fn)
		if err != nil {
			return []byte{}, err
		}
		bb = append(bb, bb_...)
		len_all = len_all + len(bb)
	}
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	valueStoreHeader.Len = int32(len_all)
	if err := binary.Write(buf, binary.LittleEndian, &valueStoreHeader); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &bb); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func LoadValueStore(bb []byte) (*Value, []byte, error) {
	debug.Alias("value.LoadValueStore").Printf("Load_value_store\r\n")
	v := &Value{}
	var value_store_header ValueStoreHeader

	var buf = bytes.NewBuffer(make([]byte, 0, len(bb)))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &value_store_header); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}
	v.typev = int(value_store_header.Type)

	switch v.typev {
	case VtBool:
	case VtInt:
		var v_int int64
		if err := binary.Read(buf, binary.LittleEndian, &v_int); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		v.value = v_int
	case VtFloat:
		var v_float float64
		if err := binary.Read(buf, binary.LittleEndian, &v_float); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		v.value = v_float
	case VtString:
		var lenght_header Lenght_header

		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		ss := make([]byte, lenght_header.LenValue)
		if err := binary.Read(buf, binary.LittleEndian, &ss); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		v.value = string(ss)
	case VtFrame:
		f := NewFrame()
		for i := 0; i < int(value_store_header.LenD); i++ {
			var lenght_header Lenght_header

			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			slot_name := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &slot_name); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}

			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			slot_prop := make([]byte, lenght_header.LenValue)
			if lenght_header.LenValue > 0 {
				if err := binary.Read(buf, binary.LittleEndian, &slot_prop); err != nil {
					fmt.Println(err)
					return nil, []byte{}, err
				}
			}
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			f.AddSlot(string(slot_name))
			//slice := []*Value{}
			for i := 0; i < int(lenght_header.LenValue); i++ {
				var lenght_header1 Lenght_header

				if err := binary.Read(buf, binary.LittleEndian, &lenght_header1); err != nil {
					fmt.Println(err)
					return nil, []byte{}, err
				}
				ssh := make([]byte, lenght_header1.LenValue)
				if lenght_header1.LenValue > 0 {
					if err := binary.Read(buf, binary.LittleEndian, &ssh); err != nil {
						fmt.Println(err)
						return nil, []byte{}, err
					}
				}

				csh, _, err := LoadValueStore(ssh)
				if err != nil {
					fmt.Println(err)
					return nil, []byte{}, err
				}
				f.SetValue(string(slot_name), csh)
				//slice = append(slice, csh)
			}
		}

		v.value = f
	case VtSlice:
		slice := []*Value{}
		for i := 0; i < int(value_store_header.LenD); i++ {
			var lenght_header Lenght_header

			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			ssh := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &ssh); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}

			csh, _, err := LoadValueStore(ssh)
			if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			slice = append(slice, csh)
		}
		v.value = slice
	case VtFunction:
		fp, _, err1 := fnc.Bin2Func(buf.Bytes())
		if err1 != nil {
			fmt.Print(err1)
			return nil, []byte{}, err1
		}
		v.value = fp
	}
	bb_r := buf.Next(buf.Len())
	return v, bb_r, nil
}
