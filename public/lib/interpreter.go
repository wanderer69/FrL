package frl

import (
	"fmt"
	"strconv"

	"github.com/wanderer69/FrL/public/script"
	attr "github.com/wanderer69/tools/parser/attributes"
	"github.com/wanderer69/tools/unique"

	fnc "github.com/wanderer69/FrL/public/functions"
	ops "github.com/wanderer69/FrL/public/operators"
	ns "github.com/wanderer69/tools/parser/new_strings"
	"github.com/wanderer69/tools/parser/parser"
	print "github.com/wanderer69/tools/parser/print"
)

type BreakPoint struct {
	FileName string
	LineNum  int
}

// Это интерпретатор кода который может транслировать исходный текст в массив операторов, интерпретировать его, сохранять операторы и загружать их.
type InterpreterEnv struct {
	Output *print.Output
	FE     *FrameEnvironment

	debug          int
	functions      []*fnc.Function
	functionsDict  map[string]*fnc.Function
	methods        []*fnc.Method
	methodsDict    map[string]*fnc.Method
	contextEnv     []*ContextEnv
	contextEnvDict map[string]*ContextEnv
	intFuncs       map[string]func(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error)
	intMethods     map[string]func(ie *InterpreterEnv, state int, if_ *InternalMethod) (*InternalMethod, bool, error)

	breakPointsList   map[string]*BreakPoint
	currentBreakPoint *BreakPoint
}

func (ie *InterpreterEnv) AddBreakPoints(breakPoints []*BreakPoint) {
	for i := range breakPoints {
		k := fmt.Sprintf("%v_%v", breakPoints[i].FileName, breakPoints[i].LineNum)
		ie.breakPointsList[k] = breakPoints[i]
	}
}

func (ie *InterpreterEnv) AddBreakPoint(breakPoint *BreakPoint) {
	k := fmt.Sprintf("%v_%v", breakPoint.FileName, breakPoint.LineNum)
	ie.breakPointsList[k] = breakPoint
}

func (ie *InterpreterEnv) GetCurrentBreakPoint() *BreakPoint {
	return ie.currentBreakPoint
}

func (ie *InterpreterEnv) ClearCurrentBreakPoint() {
	ie.currentBreakPoint = nil
}

type ContextEnv struct {
	id      string
	stack   []*ContextFunc
	current *ContextFunc
}

func (ce *ContextEnv) GetCurrentFunc() *ContextFunc {
	return ce.current
}

const (
	CFTypeFunction = 1
	CFTypeMethod   = 2
)

type ContextFunc struct {
	ce       *ContextEnv
	stack    []*Value
	args     []*Value
	varDict  map[string]*Value
	codeType byte // признак метод или функция
	method   *fnc.Method
	function *fnc.Function
	pos      int
	ops      []*ops.Operator
}

func (cf *ContextFunc) GetFunc() *fnc.Function {
	return cf.function
}

func (cf *ContextFunc) GetVarDict() map[string]*Value {
	return cf.varDict
}

func (ie *InterpreterEnv) SetFrameEnvironment(fe *FrameEnvironment) {
	ie.FE = fe
}

func (cf *ContextFunc) PrintContextFunc(o *print.Output) {
	tt := "--"

	switch cf.codeType {
	case CFTypeMethod:
		o.Print("%v%v", tt, cf.method.Name)
		for k, v := range cf.varDict {
			vv, ok := FromType(v)
			if ok {
				o.Print(" %v:%v", k, vv)

			} else {
				o.Print(" Error value %v", k)
			}
		}

	case CFTypeFunction:
		o.Print("%v%v", tt, cf.function.Name)
		for k, v := range cf.varDict {
			vv, ok := FromType(v)
			if ok {
				o.Print(" %v:%v", k, vv)

			} else {
				o.Print(" Error value %v", k)
			}
		}
	}

	o.Print("\r\nStack ")
	if len(cf.stack) > 0 {
		for i := range cf.stack {
			v := cf.stack[i]
			ss, ok := FromType(v)
			if !ok {
				o.Print("error %v", v)
			} else {
				o.Print("\r\n%v", ss)
			}
		}
		o.Print("\r\n")
	} else {
		o.Print(" empty\r\n")
	}

	o.Print("\r\n%vpos %v\r\n", tt, cf.pos)
}

// создает окружение интерпретатора
func NewInterpreterEnv() *InterpreterEnv {
	ie := InterpreterEnv{}
	ie.methodsDict = make(map[string]*fnc.Method)
	ie.intFuncs = make(map[string]func(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error))
	ie.functionsDict = make(map[string]*fnc.Function)
	ie.intMethods = make(map[string]func(ie *InterpreterEnv, state int, if_ *InternalMethod) (*InternalMethod, bool, error))
	ie.contextEnvDict = make(map[string]*ContextEnv)
	ie.breakPointsList = make(map[string]*BreakPoint)

	return &ie
}

func (ie *InterpreterEnv) SetDebug(d int) {
	ie.debug = d
}

func (ie *InterpreterEnv) CreateContextEnv() (*ContextEnv, error) {
	res := &ContextEnv{}
	res.id = unique.UniqueValue(8)
	ie.contextEnv = append(ie.contextEnv, res)
	ie.contextEnvDict[res.id] = res
	return res, nil
}

func (ie *InterpreterEnv) CreateContextFunc(ce *ContextEnv, name string, args []*Value) (*ContextFunc, error) {
	res := &ContextFunc{}
	res.ce = ce
	res.varDict = make(map[string]*Value)
	res.args = args

	f, ok := ie.functionsDict[name]
	if !ok {
		m, ok := ie.methodsDict[name]
		if !ok {
			return nil, fmt.Errorf("function or method %v not found", name)
		} else {
			res.method = m
			res.codeType = CFTypeMethod
			res.ops = m.Operators
		}
	} else {
		res.function = f
		res.codeType = CFTypeFunction
		res.ops = f.Operators
	}
	res.pos = 0

	if false {
		for i := range res.ops {
			s := ops.PrintOperator(*res.ops[i])
			ie.Output.Print("o: %v\r\n", s)
		}
	}

	return res, nil
}

func (ie *InterpreterEnv) AddMethod(method *fnc.Method) error {
	_, ok := ie.methodsDict[method.Name]
	if ok {
		return fmt.Errorf("method %v found", method.Name)
	}
	ie.methods = append(ie.methods, method)
	ie.methodsDict[method.Name] = method

	return nil
}

func (ie *InterpreterEnv) AddFunction(function *fnc.Function) error {
	_, ok := ie.functionsDict[function.Name]
	if ok {
		return fmt.Errorf("function %v found", function.Name)
	}
	ie.functions = append(ie.functions, function)
	ie.functionsDict[function.Name] = function

	return nil
}

func (ie *InterpreterEnv) BindFunction(fn func(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error)) error {
	// спрашиваем имя
	if_, _, _ := fn(ie, 0, nil)
	if if_ != nil {
		ie.intFuncs[if_.Name] = fn
		return nil
	}

	return fmt.Errorf("error bind function")
}

func (ie *InterpreterEnv) BindMethod(fn func(ie *InterpreterEnv, state int, if_ *InternalMethod) (*InternalMethod, bool, error)) error {
	// спрашиваем имя
	if_, _, _ := fn(ie, 0, nil)
	if if_ != nil {
		ie.intMethods[if_.Name] = fn
		return nil
	}

	return fmt.Errorf("error bind function")
}

// возвращаемый результат flag_mode bool - true если это встроенная и false если созданная,
//
//	result []*Value - результат если он есть
//	flag_quit - для встроенной функции признак завершения выполнения
//	err error - ошибка
func (ie *InterpreterEnv) CallFunction(name string, vl []*Value /*, cfo *ContextFunc */) (bool, []*Value, bool, error) {
	// проверяем внутреннюю функцию
	fn, ok := ie.intFuncs[name]
	if ok {
		// спрашиваем число аргументов
		if_, _, _ := fn(ie, 1, nil)
		if if_ != nil {
			if_.Args = vl
			if_, ok, err := fn(ie, 2, if_)
			if if_ != nil {
				return true, if_.Return, ok, err
			} else {
				return true, nil, ok, err
			}
		}
	} else {
		// проверяем созданную функцию
		_, ok := ie.functionsDict[name]
		if !ok {
			return false, nil, true, fmt.Errorf("function %v not found", name)
		}
		// формируем вызов
		ce := ie.contextEnv[len(ie.contextEnv)-1]
		cf, err1 := ie.CreateContextFunc(ce, name, vl)
		if err1 != nil {
			return false, nil, true, fmt.Errorf("create context func error %v", err1)
		}
		ce.stack = append(ce.stack, ce.current)
		ce.current = cf
	}
	return false, nil, true, nil
}

func (ie *InterpreterEnv) CallMethod(name string, vl []*Value, cf *ContextFunc) (bool, error) {
	fn := ie.intMethods[name]
	// спрашиваем число аргументов
	if_, _, _ := fn(ie, 1, nil)
	if if_ != nil {
		if_.Args = vl
		_, ok, err := fn(ie, 2, if_)
		return ok, err
	}
	return false, nil
}

// транслирует исходный код в последовательность операторов
func (ie *InterpreterEnv) TranslateText(name string, data string, debug int, o *print.Output) ([]*Value, error) {
	env := parser.NewEnv()
	env.Output = ie.Output

	script.MakeRules(env)

	rEnv := script.NewEnvironment()
	rp := script.FrameParser{}
	rp.Env = rEnv
	env.Struct = rp

	env.Debug = debug
	res, err := env.ParseString(data, o)
	if err != nil {
		return nil, err
	}
	if ie.debug > 3 {
		ie.Output.Print("ParseString: %v\r\n", res)
	}

	if ie.debug > 8 {
		for i := range rEnv.Relations {
			r := rEnv.Relations[i]
			s := script.PrintRelations(r)
			ie.Output.Print("%v\r\n", s)
		}
		for i := range rEnv.Frames {
			r := rEnv.Frames[i]
			s := script.PrintFrames(r)
			ie.Output.Print("%v\r\n", s)
		}
		for i := range rEnv.Methods {
			r := rEnv.Methods[i]
			s := fnc.PrintMethod(r)
			ie.Output.Print("%v\r\n", s)
		}
		for i := range rEnv.Functions {
			r := rEnv.Functions[i]
			s := fnc.PrintFunction(r)
			ie.Output.Print("%v\r\n", s)
		}
	}
	fvl := []*Value{}
	for i := range rEnv.Functions {
		rEnv.Functions[i].FileName = name
		ie.AddFunction(rEnv.Functions[i])
		ie.Output.Print("f %v\r\n", rEnv.Functions[i])
		fnc.PrintFunction(rEnv.Functions[i])

		fv := CreateValue(rEnv.Functions[i])
		fv.Print(ie.Output)
		fvl = append(fvl, fv)
	}
	for i := range rEnv.Methods {
		ie.AddMethod(rEnv.Methods[i])
	}

	return fvl, nil
}

// интерпретирует последовательность операторов функции
func (ie *InterpreterEnv) InterpreterFunc(ce *ContextEnv, name string, values []*Value) (*ContextFunc, error) {
	_, ok := ie.functionsDict[name]
	if !ok {
		return nil, fmt.Errorf("function %v not found", name)
	}
	cf, err1 := ie.CreateContextFunc(ce, name, values)
	if err1 != nil {
		return nil, fmt.Errorf("create context func error %v", err1)
	}
	// формируем вызов
	ce.current = cf
	return cf, nil
}

// интерпретирует последовательность операторов метода
func (ie *InterpreterEnv) InterpreterMethod(ce *ContextEnv, name string, values []*Value) (*ContextFunc, error) {
	_, ok := ie.methodsDict[name]
	if !ok {
		return nil, fmt.Errorf("method %v not found", name)
	}
	cf, err1 := ie.CreateContextFunc(ce, name, values)
	if err1 != nil {
		return nil, fmt.Errorf("create context func error %v", err1)
	}
	ce.current = cf

	return cf, nil
}

// выполняет один оператор
func (ie *InterpreterEnv) InterpreterFuncStep( /* cf *ContextFunc */ ) (bool, error) {
	/*
		var op Operator
		switch cf.code_type {
			case "method":
				op = cf.method.Operators[cf.pos]
			case "function":
		}
	*/
	if ie.currentBreakPoint != nil {
		return false, nil
	}
	ce := ie.contextEnv[len(ie.contextEnv)-1]
	cf := ce.current

	op := cf.ops[cf.pos]

	flagChangePos := true
	debug := ie.debug
	if (ie.debug & 4) > 0 {
		ie.Output.Print("op.Name %v\r\n", ops.OpCode2Name(op.Code))
	}
	if (ie.debug & 2) > 0 {
		cf.PrintContextFunc(ie.Output)
	}
	addQriaPointWoValue := func(name string, qria []QueryRelationItem) ([]QueryRelationItem, error) {
		sl := ns.ParseStringBySignList(name, []string{"."})
		if len(sl) > 3 {
			// ошибка!
			return nil, fmt.Errorf("too more points in %v", name)
		}
		slot_name := sl[0]
		qri := QueryRelationItem{ObjectType: "relation", Object: slot_name}
		if len(sl) == 3 {
			qri.Value = CreateValue(sl[2])
		}
		qria = append(qria, qri)
		return qria, nil
	}
	type TemplItem struct {
		Templ string
		Type  string
	}
	type TemplOut struct {
		Type  string
		Value string
	}
	templ := [][]TemplItem{
		{{Templ: "const", Type: "value"}},
		{{Templ: "const", Type: "slot"}, {Templ: ":", Type: "sign"}, {Templ: "const", Type: "value"}},
		{{Templ: "const", Type: "slot"}, {Templ: ":", Type: "sign"}, {Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}},
		{{Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}},
		{{Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}, {Templ: ":", Type: "sign"}, {Templ: "const", Type: "value"}},
		{{Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}, {Templ: ":", Type: "sign"}, {Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}},
		{{Templ: "const", Type: "slot"}, {Templ: ":", Type: "sign"}, {Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}, {Templ: "?", Type: "sign"}, {Templ: "const", Type: "var"}},
	}
	parseArg := func(lst []string) []TemplOut {
		res := []TemplOut{}
		for i := range templ {

			if len(templ[i]) == len(lst) {
				n := 0
				for j := range lst {
					if lst[j] == ":" {
						if lst[j] == templ[i][j].Templ {
							n = n + 1
						}
					} else if lst[j] == "?" {
						if lst[j] == templ[i][j].Templ {
							n = n + 1
						}
					} else {
						if templ[i][j].Templ == "const" {
							n = n + 1
						}
					}
				}
				if n == len(templ[i]) {
					for j := range lst {
						to := TemplOut{templ[i][j].Type, lst[j]}
						res = append(res, to)
					}
					return res
				}
			}
		}
		return nil
	}

	createQria := func(al []*attr.Attribute) ([]QueryRelationItem, bool, error) {
		qria := []QueryRelationItem{}
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
			t, name, array := attr.GetAttribute(op.Attributes[i])
			switch t {
			case attr.AttrTConst:
				ie.Output.Print("const %v\r\n", name)
			case attr.AttrTArray:
				if debug > 4 {
					ie.Output.Print("array %v\r\n", array)
				}
				res := parseArg(array)
				slotName := ""
				varName := ""
				varName2 := ""
				value := ""
				if res != nil {
					state := 0
					for i := range res {
						if debug > 5 {
							ie.Output.Print("state %v\r\n", state)
						}
						if res[i].Type == "slot" {
							if state == 0 {
								if (res[i].Value == "slot") || (res[i].Value == "слот") {
									state = 1
								} else {
									return nil, false, fmt.Errorf("must be slot. Now %v", res[i].Value)
								}
							}
						} else if res[i].Type == "sign" {
							if res[i].Value == ":" {
								if state == 1 {
									state = 2
								} else if state == 5 {
									state = 6
								}
							} else if res[i].Value == "?" {
								if state == 0 {
									state = 4
								} else if state == 2 {
									state = 11
								} else if state == 6 {
									state = 8
								} else if state == 12 {
									state = 14
								}
							}
						} else if res[i].Type == "var" {
							if state == 4 {
								varName = res[i].Value
								state = 5
							} else if state == 11 {
								varName = res[i].Value
								state = 12
							} else if state == 14 {
								varName2 = res[i].Value
								state = 15
							}
						} else if res[i].Type == "value" {
							if state == 2 {
								slotName = res[i].Value
								state = 3
							} else if state == 5 {
								value = res[i].Value
								state = 7
							} else if state == 8 {
								varName = res[i].Value
								state = 9
							} else if state == 0 {
								value = res[i].Value
								state = 10
							} else if state == 11 {
								value = res[i].Value
								state = 13
							}
						}
					}
					if debug > 6 {
						ie.Output.Print("slot_name %v var_name %v value %v\r\n", slotName, varName, value)
					}
					switch state {
					case 3:
						if debug > 4 {
							ie.Output.Print("slot: slot_name %v\r\n", slotName)
						}
						var err error
						qria, err = addQriaPointWoValue(slotName, qria)
						if err != nil {
							return nil, false, err
						}
					case 7:
						if debug > 4 {
							ie.Output.Print("slot_name %v value %v\r\n", slotName, value)
						}
						sl := ns.ParseStringBySignList(slotName, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						sl = ns.ParseStringBySignList(value, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						qri := QueryRelationItem{ObjectType: "relation", Object: slotName}
						qri.Value = CreateValue(value)
						qria = append(qria, qri)
					case 9:
						if debug > 4 {
							ie.Output.Print("slot_name %v var_name %v\r\n", slotName, varName)
						}
						sl := ns.ParseStringBySignList(slotName, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						sl = ns.ParseStringBySignList(value, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						v, ok := cf.varDict[varName]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varName)
						}
						qri := QueryRelationItem{ObjectType: "relation", Object: slotName}
						qri.Value = v
						qria = append(qria, qri)
					case 10:
						if debug > 4 {
							ie.Output.Print("value %v\r\n", value)
						}
						var err error
						qria, err = addQriaPointWoValue(value, qria)
						if err != nil {
							return nil, false, err
						}
					case 12:
						if debug > 4 {
							ie.Output.Print("slot: var_name %v\r\n", varName)
						}
						sl := ns.ParseStringBySignList(varName, []string{"."})
						if len(sl) > 3 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						varNameLocal := sl[0]
						v, ok := cf.varDict[varNameLocal]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varNameLocal)
						}
						vt := v.GetType()
						vv, ok := FromType(v)
						if debug > 4 {
							ie.Output.Print("vt %v, FromType(v) %v ok %v\r\n", vt, vv, ok)
						}
						if vt == VtString {
							varName, _ = FromType(v)
						} else {
							return nil, false, fmt.Errorf("slot name in variable %v must be string. Now %v", varNameLocal, vt)
						}

						qri := QueryRelationItem{ObjectType: "relation", Object: varName}
						if len(sl) == 3 {
							qri.Value = CreateValue(sl[2])
						}
						qria = append(qria, qri)
					case 13:
						if debug > 4 {
							ie.Output.Print("slot: value %v\r\n", value)
						}
						var err error
						qria, err = addQriaPointWoValue(value, qria)
						if err != nil {
							return nil, false, err
						}
					case 15:
						if debug > 4 {
							ie.Output.Print("slot: var_name %v var_name2 %v\r\n", varName, varName2)
						}
						sl := ns.ParseStringBySignList(varName, []string{"."})
						if len(sl) > 2 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						varNameLocal := sl[0]
						v, ok := cf.varDict[varNameLocal]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varNameLocal)
						}

						vt := v.GetType()
						vv, ok := FromType(v)
						if debug > 4 {
							ie.Output.Print("vt %v, FromType(v) %v ok %v\r\n", vt, vv, ok)
						}
						if vt == VtString {
							varName2, _ = FromType(v) // was varName
						} else {
							return nil, false, fmt.Errorf("slot name in variable %v must be string. Now %v", varNameLocal, vt)
						}
						sl = ns.ParseStringBySignList(varName2, []string{"."})
						if len(sl) > 1 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						varNamelocal2 := sl[0]
						v, ok = cf.varDict[varNamelocal2]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varNamelocal2)
						}
						qri := QueryRelationItem{ObjectType: "relation", Object: varNamelocal2}
						qri.Value = v
						qria = append(qria, qri)
					}
				}
			}
		}
		return qria, true, nil
	}
	createAria := func(al []*attr.Attribute) ([]AddRelationItem, bool, error) {
		aria := []AddRelationItem{}
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
			t, name, array := attr.GetAttribute(op.Attributes[i])
			switch t {
			case attr.AttrTConst:
				ie.Output.Print("const %v\r\n", name)
			case attr.AttrTArray:
				if debug > 2 {
					ie.Output.Print("array %v\r\n", array)
				}
				res := parseArg(array)
				slotName := ""
				varName := ""
				varName2 := ""
				value := ""
				if res != nil {
					state := 0
					if debug > 5 {
						ie.Output.Print("res %v\r\n", res)
					}
					for i := range res {
						if debug > 5 {
							ie.Output.Print("state %v\r\n", state)
						}
						if res[i].Type == "slot" {
							if state == 0 {
								if (res[i].Value == "slot") || (res[i].Value == "слот") {
									state = 1
								} else {
									return nil, false, fmt.Errorf("must be slot. Now %v", res[i].Value)
								}
							}
						} else if res[i].Type == "sign" {
							if res[i].Value == ":" {
								if state == 1 {
									state = 2
								} else if state == 5 {
									state = 6
								}
							} else if res[i].Value == "?" {
								if state == 0 {
									state = 4
								} else if state == 2 {
									state = 11
								} else if state == 6 {
									state = 8
								} else if state == 12 {
									state = 14
								}
							}
						} else if res[i].Type == "var" {
							if state == 4 {
								varName = res[i].Value
								state = 5
							} else if state == 11 {
								varName = res[i].Value
								state = 12
							} else if state == 14 {
								varName2 = res[i].Value
								state = 15
							}
						} else if res[i].Type == "value" {
							if state == 2 {
								slotName = res[i].Value
								state = 3
							} else if state == 5 {
								value = res[i].Value
								state = 7
							} else if state == 8 {
								varName = res[i].Value
								state = 9
							} else if state == 0 {
								value = res[i].Value
								state = 10
							} else if state == 11 {
								value = res[i].Value
								state = 13
							}
						}
					}
					if debug > 4 {
						ie.Output.Print("slot_name %v var_name %v value %v\r\n", slotName, varName, value)
					}
					switch state {
					case 3:
						if debug > 4 {
							ie.Output.Print("slot: slot_name %v\r\n", slotName)
						}
						sl := ns.ParseStringBySignList(slotName, []string{"."})
						if len(sl) > 3 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", slotName)
						}
						slotNameLocal := sl[0]
						ari := AddRelationItem{ObjectType: "relation", Object: slotNameLocal}
						if len(sl) == 3 {
							ari.Value = CreateValue(sl[2])
						}
						aria = append(aria, ari)
					case 5:
						if debug > 4 {
							ie.Output.Print("var_name: %v\r\n", varName)
						}
						sl := ns.ParseStringBySignList(varName, []string{"."})
						if len(sl) > 3 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						var_name_ := sl[0]
						v, ok := cf.varDict[var_name_]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", var_name_)
						}
						vt := v.GetType()
						vv, ok := FromType(v)
						if debug > 4 {
							ie.Output.Print("vt %v, FromType(v) %v ok %v\r\n", vt, vv, ok)
						}
						if vt == VtString {
							varName, _ = FromType(v)
						} else {
							return nil, false, fmt.Errorf("slot name in variable %v must be string. Now %v", var_name_, vt)
						}

						ari := AddRelationItem{ObjectType: "relation", Object: varName}
						if len(sl) == 3 {
							ari.Value = CreateValue(sl[2])
						}
						aria = append(aria, ari)
					case 7:
						if debug > 4 {
							ie.Output.Print("slot_name %v value %v\r\n", slotName, value)
						}
						sl := ns.ParseStringBySignList(slotName, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						sl = ns.ParseStringBySignList(value, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						ari := AddRelationItem{ObjectType: "relation", Object: slotName}
						ari.Value = CreateValue(value)
						aria = append(aria, ari)
					case 9:
						if debug > 4 {
							ie.Output.Print("slot_name %v var_name %v\r\n", slotName, varName)
						}
						sl := ns.ParseStringBySignList(slotName, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						sl = ns.ParseStringBySignList(value, []string{"."})
						if len(sl) > 0 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						v, ok := cf.varDict[varName]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varName)
						}
						ari := AddRelationItem{ObjectType: "relation", Object: slotName}
						ari.Value = v
						aria = append(aria, ari)
					case 10:
						if debug > 4 {
							ie.Output.Print("value %v\r\n", value)
						}
						sl := ns.ParseStringBySignList(value, []string{"."})
						if len(sl) > 3 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", value)
						}
						slotNameLocal := sl[0]
						ari := AddRelationItem{ObjectType: "relation", Object: slotNameLocal}
						if len(sl) == 3 {
							ari.Value = CreateValue(sl[2])
						}
						aria = append(aria, ari)
					case 12:
						if debug > 4 {
							ie.Output.Print("slot: var_name %v\r\n", varName)
						}
						sl := ns.ParseStringBySignList(varName, []string{"."})
						if len(sl) > 3 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						varNameLocal := sl[0]
						v, ok := cf.varDict[varNameLocal]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varNameLocal)
						}
						vt := v.GetType()
						vv, ok := FromType(v)
						if debug > 4 {
							ie.Output.Print("vt %v, FromType(v) %v ok %v\r\n", vt, vv, ok)
						}
						if vt == VtString {
							varName, _ = FromType(v)
						} else {
							return nil, false, fmt.Errorf("slot name in variable %v must be string. Now %v", varNameLocal, vt)
						}

						ari := AddRelationItem{ObjectType: "relation", Object: varName}
						if len(sl) == 3 {
							ari.Value = CreateValue(sl[2])
						}
						aria = append(aria, ari)
					case 13:
						if debug > 4 {
							ie.Output.Print("slot: value %v\r\n", value)
						}
						sl := ns.ParseStringBySignList(value, []string{"."})
						if len(sl) > 3 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", value)
						}
						slotNameLocal := sl[0]
						ari := AddRelationItem{ObjectType: "relation", Object: slotNameLocal}
						if len(sl) == 3 {
							ari.Value = CreateValue(sl[2])
						}
						aria = append(aria, ari)
					case 15:
						if debug > 4 {
							ie.Output.Print("slot: var_name %v var_name2 %v\r\n", varName, varName2)
						}
						sl := ns.ParseStringBySignList(varName, []string{"."})
						if len(sl) > 2 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						varNameLocal := sl[0]
						v, ok := cf.varDict[varNameLocal]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varNameLocal)
						}

						vt := v.GetType()
						vv, ok := FromType(v)
						if debug > 4 {
							ie.Output.Print("vt %v, FromType(v) %v ok %v\r\n", vt, vv, ok)
						}
						if vt == VtString {
							varName2, _ = FromType(v) // was varName
						} else {
							return nil, false, fmt.Errorf("slot name in variable %v must be string. Now %v", varNameLocal, vt)
						}

						sl = ns.ParseStringBySignList(varName2, []string{"."})
						if len(sl) > 1 {
							// ошибка!
							return nil, false, fmt.Errorf("too more points in %v", name)
						}
						varNameLocal2 := sl[0]
						v, ok = cf.varDict[varNameLocal2]
						if !ok {
							return nil, false, fmt.Errorf("variable %v not found", varNameLocal2)
						}
						ari := AddRelationItem{ObjectType: "relation", Object: varNameLocal2}
						ari.Value = v
						aria = append(aria, ari)
					}
				}
			}
		}
		return aria, true, nil
	}
	push := func(v *Value) {
		cf.stack = append(cf.stack, v)
	}
	pop := func() (*Value, error) {
		if len(cf.stack) > 0 {
			v := cf.stack[len(cf.stack)-1]
			cf.stack = cf.stack[:len(cf.stack)-1]
			return v, nil
		}
		return nil, fmt.Errorf("stack empty")
	}
	/*
		set := func(name string, v *Value) error {
			var_name := array[1]
			cf.var_dict[var_name] = v
			return nil
		}
		get := func(var_name string) (*Value, error) {
			v, ok := cf.var_dict[var_name]
			if !ok {
				return fmt.Errorf("variable %v not found", var_name)
			}
			return v, nil
		}
	*/
	setA := func(a *attr.Attribute, v *Value) error {
		t, name, array := attr.GetAttribute(a)
		switch t {
		case attr.AttrTConst:
			ie.Output.Print("const %v\r\n", name)
		case attr.AttrTArray:
			if len(array) == 2 {
				if array[0] == "?" {
					var_name := array[1]
					cf.varDict[var_name] = v
					return nil
				} else {
					return fmt.Errorf("attribute %v type %v not variable", name, t)
				}
			} else {
				return fmt.Errorf("array %v too long", array)
			}
		}
		return fmt.Errorf("attribute %v type %v not variable", name, t)
	}
	getA := func(a *attr.Attribute) (*Value, error) {
		t, name, array := attr.GetAttribute(a)
		switch t {
		case attr.AttrTConst:
			ie.Output.Print("const %v\r\n", name)
		case attr.AttrTArray:
			if len(array) == 2 {
				if array[0] == "?" {
					var_name := array[1]
					v, ok := cf.varDict[var_name]
					if !ok {
						return nil, fmt.Errorf("variable %v not found", var_name)
					}
					return v, nil
				}
			} else {
				return nil, fmt.Errorf("array %v too long", array)
			}
		}
		return nil, fmt.Errorf("attribute %v type %v not variable", name, t)
	}
	getNum := func(a *attr.Attribute) (int, error) {
		t, name, _ := attr.GetAttribute(a)
		if debug > 5 {
			ie.Output.Print("t %v name %v\r\n", t, name)
		}
		switch t {
		case attr.AttrTConst:
			ie.Output.Print("const %v\r\n", name)
		case attr.AttrTArray:
			ie.Output.Print("array %v\r\n", name)
		case attr.AttrTNumber:
			pos, _ := strconv.ParseInt(name, 10, 64)
			return int(pos), nil
		}
		return 0, fmt.Errorf("attribute %v type %v not variable", name, t)
	}

	// флаг возврата и значения для возврата
	flagReturn := false
	valueReturn := []*Value{}

	switch op.Code {
	case ops.OpCargs:
		// связывает аргументы и помещяет их в словарь переменных
		pos := 0
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
			if pos < len(cf.args) {
				setA(op.Attributes[i], cf.args[pos])
				pos = pos + 1
			} else {
				setA(op.Attributes[i], CreateValue(""))
				// ????
			}
		}
	case ops.OpCconst:
		if false {
			s := op.Attributes[0].Attribute2String()
			ie.Output.Print("a: %v\r\n", s)
		}
		// константу строки или символ в стек
		if len(op.Attributes) != 1 {
			return false, fmt.Errorf("too more attributes")
		}

		t, name, array := attr.GetAttribute(op.Attributes[0])
		switch t {
		case attr.AttrTConst:
			ie.Output.Print("const %v\r\n", name)
		case attr.AttrTArray:
			if len(array) == 1 {
				value := array[0]

				// проверяем тип
				vi, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					// ошибка, пробуем float
					vf, err := strconv.ParseFloat(value, 64)
					if err != nil {
						if value[0] == '"' {
							if value[len(value)-1] == '"' {
								value = value[1 : len(value)-1]
							}
						}
						v := CreateValue(value)
						push(v)
					} else {
						v := CreateValue(vf)
						push(v)
					}
				} else {
					v := CreateValue(int(vi))
					push(v)
				}
			} else {
				return false, fmt.Errorf("array %v too long", array)
			}
		}
	case ops.OpCfind_frame:
		qria, r, err := createQria(op.Attributes)
		if err != nil {
			return r, err
		}
		if ie.FE != nil {
			fl, err := ie.FE.QueryRelations(qria...)
			if err != nil {
				return false, err
			}
			if len(fl) > 0 {
				if len(fl) == 0 {
					v := CreateValue(fl)
					push(v)
				} else {
					vl := []*Value{}
					for i := range fl {
						f := fl[i]
						v := CreateValue(f)
						vl = append(vl, v)
						f.Print(ie.Output, true)
					}
					v := CreateValue(vl)
					push(v)
				}
			} else {
				v := CreateValue(nil)
				push(v)
			}
		} else {
			ie.Output.Print("qria %v\r\n", qria)
			v := CreateValue(nil)
			push(v)
		}

	case ops.OpCadd_slots:
		aria, r, err := createAria(op.Attributes)
		if err != nil {
			return r, err
		}
		if debug > 2 {
			ie.Output.Print("aria %v\r\n", aria)
		}

		if debug > 5 {
			ie.Output.Print("len(cf.stack) %v\r\n", len(cf.stack))
		}

		v, err := pop()
		if err != nil {
			return false, err
		}
		var add_value func(v *Value)
		add_value = func(v *Value) {
			vt := v.GetType()
			vv, ok := FromType(v)
			if debug > 6 {
				ie.Output.Print("vt %v, FromType(v) %v ok %v\r\n", vt, vv, ok)
			}
			switch vt {
			case VtInt:
			case VtFloat:
			case VtString:
			case VtFrame:
				f := v.GetValue().(*Frame)
				if ie.FE != nil {
					for i := range aria {
						if aria[i].ObjectType == "relation" {
							f.AddSlot(aria[i].Object)
							if aria[i].Value != nil {
								f.SetValue(aria[i].Object, aria[i].Value)
							}
						}
					}
					ie.FE.AddRelations(f, aria...)
				}
				//f.Print(true)
			case VtSlice:
				vl := v.GetValue().([]*Value)
				for i := range vl {
					add_value(vl[i])
				}
			}
		}
		add_value(v)
	case ops.OpCset:
		if false {
			s := op.Attributes[0].Attribute2String()
			ie.Output.Print("a: %v\r\n", s)
		}

		if len(op.Attributes) != 1 {
			return false, fmt.Errorf("too more attributes")
		}

		v, err := pop()
		if err != nil {
			return false, err
		}
		err = setA(op.Attributes[0], v)
		if err != nil {
			return false, err
		}

	case ops.OpCget:
		if false {
			s := op.Attributes[0].Attribute2String()
			ie.Output.Print("a: %v\r\n", s)
		}

		if len(op.Attributes) != 1 {
			return false, fmt.Errorf("too more attributes")
		}

		v, err := getA(op.Attributes[0])
		if err != nil {
			return false, err
		}
		push(v)

	case ops.OpCframe:
		aria, r, err := createAria(op.Attributes)
		if err != nil {
			return r, err
		}
		if debug > 2 {
			ie.Output.Print("aria %v\r\n", aria)
		}
		if ie.FE != nil {
			f := ie.FE.NewFrameWithRelation()
			for i := range aria {
				if aria[i].ObjectType == "relation" {
					f.AddSlot(aria[i].Object)
					if aria[i].Value != nil {
						f.SetValue(aria[i].Object, aria[i].Value)
					}
				}
			}
			ie.FE.AddRelations(f, aria...)
			v := CreateValue(f)
			push(v)
		}
	case ops.OpCunify:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}

	case ops.OpCcreate_iterator:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// считываем из стека значение
		v, err := pop()
		if err != nil {
			return false, err
		}
		// пробуем создать итератор
		iter_v, err := NewIterator(v)
		if err != nil {
			return false, fmt.Errorf("create_iterator: NewIterator: %v", err)
		}
		if false {
			err = iter_v.Print(ie.Output)
			if err != nil {
				return false, fmt.Errorf("create_iterator: Print: %v", err)
			}
		}
		// и кладем в стек
		push(iter_v)
	case ops.OpCiteration:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// берем из стека значение
		iter_v, err := pop()
		if err != nil {
			return false, err
		}
		// делаем итерацию
		v, err := iter_v.Iterate()
		if err != nil {

			if fmt.Sprintf("%v", err) == "empty" {
				v = CreateValue(false)
			} else {
				return false, fmt.Errorf("iteration: %v", err)
			}
		}
		// и кладем в стек
		push(v)
	case ops.OpCcheck_iteration:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// берем из стека значение
		iter_v, err := pop()
		if err != nil {
			return false, err
		}
		// получаем значение и кладем его в стек
		bb, err1 := iter_v.IsEnd()
		if err1 != nil {
			return false, err1
		}
		push(bb)
	case ops.OpCcall_function:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// должно быть два аргумента имя переменной и число аргументов в стеке
		if len(op.Attributes) == 2 {
			t, name, array := attr.GetAttribute(op.Attributes[0])
			if debug > 5 {
				ie.Output.Print("t %v name %v\r\n", t, name)
			}
			var_name := ""
			switch t {
			case attr.AttrTConst:
				ie.Output.Print("const %v\r\n", name)
			case attr.AttrTArray:
				if len(array) == 1 {
					var_name = array[0]
				} else {
					return false, fmt.Errorf("call: %v not variable", array)
				}
			case attr.AttrTNumber:
				ie.Output.Print("number %v\r\n", name)
			}

			pos, err := getNum(op.Attributes[1])
			if err != nil {
				return false, fmt.Errorf("call: %v", err)
			}
			// в стеке уже вычисленные значения надо найти в стеке
			vl := []*Value{}
			for i := 0; i < pos; i++ {
				v, err := pop()
				if err != nil {
					return false, fmt.Errorf("call: %v", err)
				}
				vl = append([]*Value{v}, vl...)
			}
			flag, res, ok, err := ie.CallFunction(var_name, vl /*, cf */)
			if err != nil {
				return ok, err
			} else {
				if flag {
					// вызов встроенного метода
					if res != nil {
						if len(res) > 0 {
							for i := range res {
								push(res[i])
							}
						}
					}
				} else {
					// вызов определяемого метода
					flagChangePos = false
				}
			}
		}
	case ops.OpCcall_method:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// должно быть два аргумента имя переменной и число аргументов в стеке
		if len(op.Attributes) == 2 {
			t, name, array := attr.GetAttribute(op.Attributes[0])
			if debug > 5 {
				ie.Output.Print("t %v name %v\r\n", t, name)
			}
			var_name := ""
			switch t {
			case attr.AttrTConst:
				ie.Output.Print("const %v\r\n", name)
			case attr.AttrTArray:
				if len(array) == 1 {
					var_name = array[0]
				} else {
					return false, fmt.Errorf("call: %v not variable", array)
				}
			case attr.AttrTNumber:
				ie.Output.Print("number %v\r\n", name)
			}

			pos, err := getNum(op.Attributes[1])
			if err != nil {
				return false, fmt.Errorf("call: %v", err)
			}
			// в стеке уже вычисленные значения надо найти в стеке
			vl := []*Value{}
			for i := 0; i < pos; i++ {
				v, err := pop()
				if err != nil {
					return false, fmt.Errorf("call: %v", err)
				}
				vl = append([]*Value{v}, vl...)
			}
			flag, res, ok, err := ie.CallFunction(var_name, vl /*, cf */)
			if err != nil {
				return ok, err
			} else {
				if flag {
					// вызов встроенного метода
					if res != nil {
						if len(res) > 0 {
							for i := range res {
								push(res[i])
							}
						}
					}
				} else {
					// вызов определяемого метода
					flagChangePos = false
				}
			}

		}
	case ops.OpCbranch:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		if len(op.Attributes) == 1 {
			pos, err := getNum(op.Attributes[0])
			if err != nil {
				return false, fmt.Errorf("branch: %v", err)
			}
			// переход
			flagChangePos = false
			cf.pos = cf.pos + int(pos)

		} else {
			return false, fmt.Errorf("branch: too more attributes %v", len(op.Attributes))
		}

	case ops.OpCdup:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// берем из стека значение
		v, err := pop()
		if err != nil {
			return false, err
		}
		// результат в стек
		push(v)
		push(v)

	case ops.OpCclear:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// берем из стека значение
		_, err := pop()
		if err != nil {
			return false, err
		}

	// предикаты
	case ops.OpCeq:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// д. б. два аргумента из стека
		// берем из стека левое значение
		l_v, err := pop()
		if err != nil {
			return false, err
		}
		// берем из стека правое значение
		r_v, err := pop()
		if err != nil {
			return false, err
		}
		// сравниваем два значения
		res := CompareValuesEq(l_v, r_v)
		// результат в стек
		push(CreateValue(res))

	case ops.OpClt:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// д. б. два аргумента из стека
		// берем из стека левое значение
		l_v, err := pop()
		if err != nil {
			return false, err
		}
		// берем из стека правое значение
		r_v, err := pop()
		if err != nil {
			return false, err
		}
		// сравниваем два значения
		res := CompareValuesLt(l_v, r_v)
		// результат в стек
		push(CreateValue(res))

	case ops.OpCgt:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// д. б. два аргумента из стека
		// берем из стека левое значение
		l_v, err := pop()
		if err != nil {
			return false, err
		}
		// берем из стека правое значение
		r_v, err := pop()
		if err != nil {
			return false, err
		}
		// сравниваем два значения
		res := CompareValuesGt(l_v, r_v)
		// результат в стек
		push(CreateValue(res))

	case ops.OpCempty:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// выбираем значение из стека
		v, err := pop()
		if err != nil {
			return false, err
		}

		vr := false
		// проверяем что значение переменной - непусто то есть отлично от nil
		vt := v.GetType()
		if vt == VtNil {
			vr = true
		}
		vv := CreateValue(vr)
		// и кладем его в стек
		push(vv)

	case ops.OpCbranch_if_false:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// получаем значение из стека
		v, err := pop()
		if err != nil {
			return false, err
		}
		if len(op.Attributes) == 1 {
			pos, err := getNum(op.Attributes[0])
			if err != nil {
				return false, fmt.Errorf("branch: %v", err)
			}
			vt := v.GetType()
			if vt == VtBool {
				vv := v.Bool()
				if vv {

				} else {
					// переход так как не истина
					flagChangePos = false
					cf.pos = cf.pos + int(pos)
				}
			} else {
				return false, fmt.Errorf("branch_if_false: bad type %v", vt)
			}
		} else {
			return false, fmt.Errorf("branch_if_false: too more attributes %v", len(op.Attributes))
		}
	case ops.OpCbranch_if_true:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// получаем значение из стека
		v, err := pop()
		if err != nil {
			return false, err
		}
		if len(op.Attributes) == 1 {
			pos, err := getNum(op.Attributes[0])
			if err != nil {
				return false, fmt.Errorf("branch: %v", err)
			}
			vt := v.GetType()
			if vt == VtBool {
				vv := v.Bool()
				if vv {
					// переход так как истина
					flagChangePos = false
					cf.pos = cf.pos + int(pos)
				}
			} else {
				return false, fmt.Errorf("branch_if_false: bad type %v", vt)
			}
		} else {
			return false, fmt.Errorf("branch_if_false: too more attributes %v", len(op.Attributes))
		}
	case ops.OpCbreak:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
	case ops.OpCreturn:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		// должен быть один аргумент - число в стеке
		if len(op.Attributes) == 1 {
			pos, err := getNum(op.Attributes[0])
			if err != nil {
				return false, fmt.Errorf("call: %v", err)
			}
			// в стеке уже вычисленные значения надо найти в стеке
			for i := 0; i < pos; i++ {
				v, err := pop()
				if err != nil {
					return false, fmt.Errorf("call: %v", err)
				}
				valueReturn = append([]*Value{v}, valueReturn...)
			}
			flagReturn = true
			flagChangePos = false
		}

	case ops.OpCcontinue:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}

	case ops.OpCdebug:
		if (ie.debug & 1) > 0 {
			for i := range op.Attributes {
				if true {
					s := op.Attributes[i].Attribute2String()
					ie.Output.Print("a: %v\r\n", s)
				}
				t, name, array := attr.GetAttribute(op.Attributes[0])
				if debug > 5 {
					ie.Output.Print("t %v name %v\r\n", t, name)
				}
				switch t {
				case attr.AttrTConst:
					ie.Output.Print("const %v\r\n", name)
				case attr.AttrTArray:
					ie.Output.Print("array %v\r\n", array)
				case attr.AttrTNumber:
					ie.Output.Print("number %v\r\n", name)
				}
			}
		}

	case ops.OpCslice:
		for i := range op.Attributes {
			if false {
				s := op.Attributes[i].Attribute2String()
				ie.Output.Print("a: %v\r\n", s)
			}
		}
		if len(op.Attributes) == 1 {
			pos, err := getNum(op.Attributes[0])
			if err != nil {
				return false, fmt.Errorf("slice: %v", err)
			}
			// в стеке уже вычисленные значения надо найти в стеке
			vl := []*Value{}
			for i := 0; i < pos; i++ {
				v, err := pop()
				if err != nil {
					return false, fmt.Errorf("slice: %v", err)
				}
				vl = append([]*Value{v}, vl...)
			}
			vs, _ := NewSlice(vl...)
			push(vs)
		} else {
			return false, fmt.Errorf("slice: too more attributes %v", len(op.Attributes))
		}

	case ops.OpCline:
		num := op.Attributes[0].Number + 1
		ie.Output.Print("file name %v function name %v line: %v\r\n", cf.function.FileName, cf.function.Name, num)
		k := fmt.Sprintf("%v_%v", cf.function.FileName, num)
		bp, ok := ie.breakPointsList[k]
		if ok {
			ie.currentBreakPoint = bp
		}
	default:
		return false, fmt.Errorf("bad operator code %v", op.Code)
	}
	// был возврат?
	if flagReturn {
		// надо проверить стек
		if len(ce.stack) > 0 {
			// в стеке есть
			v := ce.stack[len(ce.stack)-1]
			ce.stack = ce.stack[:len(ce.stack)-1]
			ce.current = v
			// и надо позаботиться о возврате - если был return, то выбираем то что он оставил и заносим в стек
			v.stack = append(v.stack, valueReturn...)
			v.pos = v.pos + 1
		}
		return false, nil
	}
	if len(cf.ops) == cf.pos+1 {
		// закончилось тело функции
		// надо проверить стек
		if len(ce.stack) > 0 {
			// в стеке есть
			v := ce.stack[len(ce.stack)-1]
			ce.stack = ce.stack[:len(ce.stack)-1]
			ce.current = v
			v.pos = v.pos + 1
			return false, nil
		} else {
			return true, nil
		}
	}
	if flagChangePos {
		cf.pos = cf.pos + 1
	}
	if len(cf.ops) <= cf.pos {
		return false, fmt.Errorf("pos %v out of range", cf.pos)
	}
	return false, nil
}
