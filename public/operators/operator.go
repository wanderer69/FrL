package operators

import (
	"bytes"
	"encoding/binary"
	"fmt"

	attr "github.com/wanderer69/tools/parser/attributes"
)

type Operator struct {
	Code       byte
	Attributes []*attr.Attribute // перечень атрибутов факта
}

const (
	OpCargs           = 1
	OpCconst          = 2
	OpCfindFrame      = 3
	OpCaddSlots       = 4
	OpCset            = 5
	OpCget            = 6
	OpCframe          = 7
	OpCunify          = 8
	OpCcreateIterator = 9
	OpCiteration      = 10
	OpCcheckIteration = 11
	OpCcallFunction   = 12
	OpCcallMethod     = 13
	OpCbranch         = 14
	OpCdup            = 15
	OpCclear          = 16
	OpCeq             = 17
	OpClt             = 18
	OpCgt             = 19
	OpCempty          = 20
	OpCbranchIfFalse  = 21
	OpCbranchIfTrue   = 22
	OpCbreak          = 23
	OpCreturn         = 24
	OpCcontinue       = 25
	OpCdebug          = 26
	OpCslice          = 27
	OpCline           = 28
	OpCfindSlot       = 29
	OpCtemplate       = 30
)

func OpName2Code(name string) byte {
	var res byte = 0
	switch name {
	case "args":
		res = OpCargs
	case "const":
		res = OpCconst
	case "findFrame":
		res = OpCfindFrame
	case "addSlots":
		res = OpCaddSlots
	case "set":
		res = OpCset
	case "get":
		res = OpCget
	case "frame":
		res = OpCframe
	case "template":
		res = OpCtemplate
	case "unify":
		res = OpCunify
	case "createIterator":
		res = OpCcreateIterator
	case "iteration":
		res = OpCiteration
	case "checkIteration":
		res = OpCcheckIteration
	case "callFunction":
		res = OpCcallFunction
	case "callMethod":
		res = OpCcallMethod
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
	case "branchIfFalse":
		res = OpCbranchIfFalse
	case "branchIfTrue":
		res = OpCbranchIfTrue
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
	case "findSlot":
		res = OpCfindSlot
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
	case OpCfindFrame:
		res = "findFrame"
	case OpCaddSlots:
		res = "addSlots"
	case OpCset:
		res = "set"
	case OpCget:
		res = "get"
	case OpCframe:
		res = "frame"
	case OpCtemplate:
		res = "template"
	case OpCunify:
		res = "unify"
	case OpCcreateIterator:
		res = "createIterator"
	case OpCiteration:
		res = "iteration"
	case OpCcheckIteration:
		res = "checkIteration"
	case OpCcallFunction:
		res = "callFunction"
	case OpCcallMethod:
		res = "callMethod"
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
	case OpCbranchIfFalse:
		res = "branchIfFalse"
	case OpCbranchIfTrue:
		res = "branchIfTrue"
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
	case OpCfindSlot:
		res = "findSlot"
	}
	return res
}

func PrintOperator(o Operator) string {
	result := fmt.Sprintf("%v (", OpCode2Name(o.Code))
	for i := range o.Attributes {
		result = result + o.Attributes[i].Attribute2String()
	}
	result = result + ")"
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
