package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"

	print "github.com/wanderer69/FrL/internal/lib/print"
)

const (
	AttrTConst  = 1
	AttrTArray  = 2
	AttrTNumber = 3
)

type Attribute struct {
	Type   byte // константа знак строка массив число
	Const  string
	Array  []string
	Number int
}

func NewAttribute(type_ byte, value string, array []string) *Attribute {
	a := Attribute{}
	a.Type = type_
	switch a.Type {
	case AttrTConst:
		a.Const = value
	case AttrTArray:
		a.Array = array
	case AttrTNumber:
		v, _ := strconv.ParseInt(value, 10, 64)
		a.Number = int(v)
	default:
		panic(fmt.Sprintf("Bad type %v. Must be const, array, string, sign, number\r\n", a.Type))
	}
	return &a
}

func (a *Attribute) PrintAttribute(o *print.Output) {
	switch a.Type {
	case AttrTConst:
		o.Print("const %v\r\n", a.Const)
	case AttrTArray:
		o.Print("array %v\r\n", a.Array)
	case AttrTNumber:
		o.Print("number %v\r\n", a.Number)
	default:
		panic(fmt.Sprintf("Bad type %v. Must be const, array, string, sign, number\r\n", a.Type))
	}
}

func (a *Attribute) Attribute2String() string {
	res := ""
	switch a.Type {
	case AttrTConst:
		res = fmt.Sprintf(" const %v", a.Const)
	case AttrTArray:
		res = fmt.Sprintf(" array %v", a.Array)
	case AttrTNumber:
		res = fmt.Sprintf(" number %v", a.Number)
	default:
		panic(fmt.Sprintf("Bad type %v. Must be const, array, string, sign, number\r\n", a.Type))
	}
	return res
}

func GetAttribute(a *Attribute) (byte, string, []string) {
	if a != nil {
		switch a.Type {
		case AttrTConst:
			return a.Type, a.Const, nil
		case AttrTArray:
			return a.Type, "", a.Array
		case AttrTNumber:
			return a.Type, fmt.Sprintf("%v", a.Number), nil
		default:
			panic(fmt.Sprintf("Bad type %v. Must be const, array, string, sign, number\r\n", a.Type))
		}
	}
	return 0, "", nil
}

type Lenght_header struct {
	LenValue int32
}

func Save_lenght_value(bb []byte) []byte {
	lenght_header := Lenght_header{LenValue: int32(len(bb))}

	len_all := (int)(unsafe.Sizeof(lenght_header)) + len(bb)
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &lenght_header); err != nil {
		fmt.Println(err)
	}
	if lenght_header.LenValue > 0 {
		if err := binary.Write(buf, binary.LittleEndian, bb); err != nil {
			fmt.Println(err)
		}
	}
	return buf.Bytes()
}

func (a *Attribute) Attribute2Bin() ([]byte, error) {
	bb := []byte{}
	vb := byte(a.Type)
	bb = append(bb, vb)

	switch a.Type {
	case AttrTConst:
		bb_ := Save_lenght_value([]byte(a.Const))
		bb = append(bb, bb_...)
	case AttrTArray:
		vb := int32(len(a.Array))
		b_in := make([]byte, 0, 4)
		var buf = bytes.NewBuffer(b_in)
		if err := binary.Write(buf, binary.LittleEndian, &vb); err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = append(bb, buf.Bytes()...)
		for i := range a.Array {
			bb_ := Save_lenght_value([]byte(a.Array[i]))
			bb = append(bb, bb_...)
		}
	case AttrTNumber:
		vb := int32(len(a.Array))
		b_in := make([]byte, 0, 4)
		var buf = bytes.NewBuffer(b_in)
		if err := binary.Write(buf, binary.LittleEndian, &vb); err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = append(bb, buf.Bytes()...)
	default:
		panic(fmt.Sprintf("Bad type %v/ Must be const, array, string, sign, number\r\n", a.Type))
	}
	return bb, nil
}

func Bin2Attribute(bb []byte) (*Attribute, []byte, error) {
	a := Attribute{}

	var t byte
	var buf = bytes.NewBuffer(make([]byte, 0, 1))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &t); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	a.Type = t
	switch a.Type {
	case AttrTConst:
		var lenght_header Lenght_header

		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		value := make([]byte, lenght_header.LenValue)
		if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
			fmt.Println(err)
			return nil, nil, err
		}

		a.Const = string(value)
	case AttrTArray:
		var lenght int32

		if err := binary.Read(buf, binary.LittleEndian, &lenght); err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		for i := 0; i < int(lenght); i++ {
			var lenght_header Lenght_header

			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, nil, err
			}
			value := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				fmt.Println(err)
				return nil, nil, err
			}
			a.Array = append(a.Array, string(value))
		}
	case AttrTNumber:
		var v int32

		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			fmt.Println(err)
			return nil, nil, err
		}

		a.Number = int(v)
	default:
		panic(fmt.Sprintf("Bad type %v/ Must be const, array, string, sign, number\r\n", a.Type))
	}
	return &a, buf.Bytes(), nil
}
