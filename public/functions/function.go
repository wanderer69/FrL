package functions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	ops "github.com/wanderer69/FrL/public/operators"
	attr "github.com/wanderer69/tools/parser/attributes"
)

type Function struct {
	Name      string
	NumArgs   int
	Operators []*ops.Operator
	FileName  string
	Package   string
}

// внешняя функция
type ExternalFunction struct {
	Name    string
	Alias   string
	Func    func(args []interface{}) ([]interface{}, bool, error)
	NumArgs int
	Args    []string
}

type Method struct {
	Name      string
	Operators []*ops.Operator
}

func PrintFunction(r *Function) string {
	result := fmt.Sprintf("Function %v {\r\n", r.Name)
	for i := range r.Operators {
		c := r.Operators[i]
		result = result + fmt.Sprintf("%v\r\n", ops.PrintOperator(*c))
	}
	result = result + fmt.Sprintf("} # %v", r.FileName)
	return result
}

func PrintMethod(r *Method) string {
	result := fmt.Sprintf("Method %v {\r\n", r.Name)
	for i := range r.Operators {
		c := r.Operators[i]
		result = result + fmt.Sprintf("%v\r\n", ops.PrintOperator(*c))
	}
	result = result + "}"
	return result
}

func Func2Bin(f *Function) ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("function nil")
	}
	bb := []byte{}

	bb_ := attr.Save_lenght_value([]byte(f.Name))
	bb = append(bb, bb_...)

	vb := int32(len(f.Operators))
	b_in := make([]byte, 0, 4)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &vb); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	bb = append(bb, buf.Bytes()...)

	na := int32(f.NumArgs)
	b_in = make([]byte, 0, 4)
	buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &na); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	bb = append(bb, buf.Bytes()...)

	for i := range f.Operators {
		bb_, err := ops.Operator2Bin(f.Operators[i])
		if err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = append(bb, bb_...)
	}
	return bb, nil
}

func Bin2Func(bb []byte) (*Function, []byte, error) {
	res := Function{}

	var lenght_header attr.Lenght_header

	var buf = bytes.NewBuffer(make([]byte, 0, unsafe.Sizeof(lenght_header)))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	value := make([]byte, lenght_header.LenValue)

	if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	res.Name = string(value)

	var lenght int32

	if err := binary.Read(buf, binary.LittleEndian, &lenght); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	bb = buf.Bytes()

	var na int32

	if err := binary.Read(buf, binary.LittleEndian, &na); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	bb = buf.Bytes()

	res.NumArgs = int(na)

	for i := 0; i < int(lenght); i++ {
		a, bb_, err := ops.Bin2Operator(bb)
		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		bb = bb_
		res.Operators = append(res.Operators, a)
	}
	return &res, bb, nil
}
