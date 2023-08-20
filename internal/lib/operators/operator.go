package operators

import (
	"bytes"
	"encoding/binary"
	"fmt"

	attr "github.com/wanderer69/FrL/internal/lib/attributes"
)

type Operator struct {
	Code       byte
	Attributes []*attr.Attribute // перечень атрибутов факта
}

const (
	OpCargs            = 1
	OpCconst           = 2
	OpCfind_frame      = 3
	OpCadd_slots       = 4
	OpCset             = 5
	OpCget             = 6
	OpCframe           = 7
	OpCunify           = 8
	OpCcreate_iterator = 9
	OpCiteration       = 10
	OpCcheck_iteration = 11
	OpCcall_function   = 12
	OpCcall_method     = 13
	OpCbranch          = 14
	OpCdup             = 15
	OpCclear           = 16
	OpCeq              = 17
	OpClt              = 18
	OpCgt              = 19
	OpCempty           = 20
	OpCbranch_if_false = 21
	OpCbranch_if_true  = 22
	OpCbreak           = 23
	OpCreturn          = 24
	OpCcontinue        = 25
	OpCdebug           = 26
	OpCslice           = 27
	OpCline            = 28
)

func OpName2Code(name string) byte {
	var res byte = 0
	switch name {
	case "args":
		res = OpCargs
	case "const":
		res = OpCconst
	case "find_frame":
		res = OpCfind_frame
	case "add_slots":
		res = OpCadd_slots
	case "set":
		res = OpCset
	case "get":
		res = OpCget
	case "frame":
		res = OpCframe
	case "unify":
		res = OpCunify
	case "create_iterator":
		res = OpCcreate_iterator
	case "iteration":
		res = OpCiteration
	case "check_iteration":
		res = OpCcheck_iteration
	case "call_function":
		res = OpCcall_function
	case "call_method":
		res = OpCcall_method
	case "branch":
		res = OpCbranch
	case "dup":
		res = OpCdup
	case "clear":
		res = OpCclear
	case "eq":
		res = OpCeq
	case "lt":
		res = OpClt
	case "gt":
		res = OpCgt
	case "empty":
		res = OpCempty
	case "branch_if_false":
		res = OpCbranch_if_false
	case "branch_if_true":
		res = OpCbranch_if_true
	case "break":
		res = OpCbreak
	case "return":
		res = OpCreturn
	case "continue":
		res = OpCcontinue
	case "debug":
		res = OpCdebug
	case "slice":
		res = OpCslice
	case "line":
		res = OpCline
	default:
		panic(fmt.Errorf("bad operator name %v", name))
	}
	return res
}

func OpCode2Name(c byte) string {
	res := ""
	switch c {
	case OpCargs:
		res = "args"
	case OpCconst:
		res = "const"
	case OpCfind_frame:
		res = "find_frame"
	case OpCadd_slots:
		res = "add_slots"
	case OpCset:
		res = "set"
	case OpCget:
		res = "get"
	case OpCframe:
		res = "frame"
	case OpCunify:
		res = "unify"
	case OpCcreate_iterator:
		res = "create_iterator"
	case OpCiteration:
		res = "iteration"
	case OpCcheck_iteration:
		res = "check_iteration"
	case OpCcall_function:
		res = "call_function"
	case OpCcall_method:
		res = "call_method"
	case OpCbranch:
		res = "branch"
	case OpCdup:
		res = "dup"
	case OpCclear:
		res = "clear"
	case OpCeq:
		res = "eq"
	case OpClt:
		res = "lt"
	case OpCgt:
		res = "gt"
	case OpCempty:
		res = "empty"
	case OpCbranch_if_false:
		res = "branch_if_false"
	case OpCbranch_if_true:
		res = "branch_if_true"
	case OpCbreak:
		res = "break"
	case OpCreturn:
		res = "return"
	case OpCcontinue:
		res = "continue"
	case OpCdebug:
		res = "debug"
	case OpCslice:
		res = "slice"
	case OpCline:
		res = "line"
	}
	return res
}

func PrintOperator(o Operator) string {
	result := fmt.Sprintf("%v (", OpCode2Name(o.Code))
	for i := range o.Attributes {
		result = result + fmt.Sprintf("%v ", o.Attributes[i].Attribute2String())
	}
	result = result + fmt.Sprintf(")")
	return result
}

func Operator2Bin(o *Operator) ([]byte, error) {
	bb := []byte{}
	bb = append(bb, o.Code)

	vb := int32(len(o.Attributes))
	b_in := make([]byte, 0, 4)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &vb); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	bb = append(bb, buf.Bytes()...)
	for i := range o.Attributes {
		bb_, err := o.Attributes[i].Attribute2Bin()
		if err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
		bb = append(bb, bb_...)
	}
	return bb, nil
}

func Bin2Operator(bb []byte) (*Operator, []byte, error) {
	res := Operator{}
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
	res.Code = t
	var lenght int32

	if err := binary.Read(buf, binary.LittleEndian, &lenght); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	bb = buf.Bytes()
	for i := 0; i < int(lenght); i++ {
		a, bb_, err := attr.Bin2Attribute(bb)
		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		bb = bb_
		res.Attributes = append(res.Attributes, a)
	}
	return &res, bb, nil
}
