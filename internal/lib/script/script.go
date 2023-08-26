package script

import (
	"fmt"
	"strings"

	fnc "github.com/wanderer69/FrL/internal/lib/functions"
	ops "github.com/wanderer69/FrL/internal/lib/operators"
	attr "github.com/wanderer69/tools/parser/attributes"
	ns "github.com/wanderer69/tools/parser/new_strings"
	uqe "github.com/wanderer69/tools/unique"

	parser "github.com/wanderer69/tools/parser/parser"
)

// грамматика !!!
/* верхнеуровневые элементы

<symbols, == отношения> <{, > - добавление отношений во фреймы
<symbols, > <symbols, > <symbols, > - отношение

<symbols, == фреймы> <{, > - определение фреймов
<symbols, == фрейм> <(, > - определение фреймА, после ключевого слова идет список

<(, > <symbols, == => > <?, > <symbols,> - присваивание результатат поиска
<[, > <symbols, == => > <?, > <symbols,> - присваивание результатат поиска созданного из списка

<(, > <symbols, == %> <symbols, == => > <symbols, == ?> <symbols,> - присваивание списка слотов найденного фрейма

<(, > <symbols, == ?> <symbols,> <symbols, == => > <symbols, == ?> <symbols,> - унификация (сущность.объект)?элемент-класс

<symbols, == если> <(, >  <{, > - если
<symbols, == если> <(, >  <{, > <symbols, == иначе> <{, > - если - иначе

<symbols, == для> <symbols, == каждого> <(, >  <symbols, == => > <symbols, == ?> <symbols,> <{, > - цикл для каждого фрейма
<symbols, == для> <symbols, == каждого> <(, >  <symbols, == %> <symbols, == => > <symbols, == ?> <symbols,> <{, > - цикл для каждого слота фрейма

<symbols, == @> <symbols, > <(, > <{, > -  определение метода
<symbols, > <(, > - вызов метода

<symbols, [0] == "?"> <symbols, == => > <symbols,> - переприсваивание значения переменой другой переменной

<string, > <symbols, == => > <symbols,> - константу строку в переменную
<(, > <symbols, == => > <symbols,> - константу список в переменную
<[, > <symbols, == => > <symbols,> - константу массив в переменную
<{, > <symbols, == => > <symbols,> - константу словарь в переменную
<symbols, == Факт> <(, >  <symbols, == => > <symbols,> - определение факта
<symbols, == Шаблон> <(, >  <symbols, == => > <symbols,> - определение шаблона
<symbols, == Лисп> <(, > - вставка на чистом Лиспе
<symbols, > <(, > - вызов функции без возврата значения
<symbols, > <symbols, == => > <symbols,>  - вызов функции с возвратом значения
среднеуровневые элементы
операторы
<symbols, == Если> <(, >  <{, > - если
<symbols, == Если> <(, >  <{, > <symbols, == Иначе> <{, > - если иначе
<symbols, == Цикл> <symbols, == по> <symbols, [0] == "?">  <symbols, == => > <symbols,> <{, > - Цикл по
<symbols, == Вернуть> <symbols, > - вернуть
список в обпределении тринара или шаблона
<symbols, > <symbols, > <symbols, > - список тринаров

*/

type FrameParserStackItem struct {
	ConditionOps []*ops.Operator
	ExecOps      []*ops.Operator
}

type FrameParser struct {
	CurOp     ops.Operator
	Operators []*ops.Operator // список операторов

	Stack    []*FrameParserStackItem
	StackPos int

	Env *Environment
}

type Relations struct {
	Relations []*ops.Operator
}

type Frames struct {
	Frames []*ops.Operator
}

func PrintRelations(r *Relations) string {
	result := "Relations {\r\n"
	for i := range r.Relations {
		c := r.Relations[i]
		result = result + fmt.Sprintf("%v\r\n", ops.PrintOperator(*c))
	}
	result = result + "}"
	return result
}

func PrintFrames(r *Frames) string {
	result := "Frames {\r\n"
	for i := range r.Frames {
		c := r.Frames[i]
		result = result + fmt.Sprintf("%v\r\n", ops.PrintOperator(*c))
	}
	result = result + "}"
	return result
}

type Environment struct {
	Relations []*Relations
	Frames    []*Frames
	Functions []*fnc.Function
	Methods   []*fnc.Method
}

func NewEnvironment() *Environment {
	env := Environment{}
	return &env
}

func (env *Environment) AddRelations(r Relations) {
	env.Relations = append(env.Relations, &r)
}

func (env *Environment) AddFrames(r Frames) {
	env.Frames = append(env.Frames, &r)
}

func (env *Environment) AddMethod(r fnc.Method) {
	env.Methods = append(env.Methods, &r)
}

func (env *Environment) AddFunction(r fnc.Function) {
	env.Functions = append(env.Functions, &r)
}

func ParseArg(val string) (*attr.Attribute, error) {
	ssl := ns.ParseStringBySignList(val, []string{"?", ":"})
	a := attr.NewAttribute(attr.AttrTArray, "", ssl)
	return a, nil
}

func fRelations(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		// список операторов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)
		rp.Operators = []*ops.Operator{} // список операторов

		rp.Stack = []*FrameParserStackItem{}
		rp.StackPos = -1
		env.Struct = rp

	case 1:
		body := env.CE.ResultGenerate
		result = fmt.Sprintf("(relation %v)", body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		r := Relations{rp.Operators}
		rp.Env.AddRelations(r)
		env.Struct = rp
	}
	return result, nil
}

func fFrames(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		// список операторов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)

		rp.Operators = []*ops.Operator{} // список операторов

		rp.Stack = []*FrameParserStackItem{}
		rp.StackPos = -1
		env.Struct = rp

	case 1:
		body := env.CE.ResultGenerate
		result = fmt.Sprintf("(frames %v)", body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		r := Frames{rp.Operators}
		rp.Env.AddFrames(r)
		env.Struct = rp
	}
	return result, nil
}

func fFunction(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		function_name := pi.Items[1].Data
		env.CE.StringVars["function_name"] = function_name
		// список операторов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)

		rp.Operators = []*ops.Operator{} // список операторов

		rp.Stack = []*FrameParserStackItem{}
		rp.StackPos = -1
		env.Struct = rp
	case 1:
		body := env.CE.ResultGenerate
		function_name := env.CE.StringVars["function_name"]
		args_list := pi.Items[2].Data

		result = fmt.Sprintf("(function %v (%v) %v)", function_name, args_list, body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op}, rp.Operators...)
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("args")

		b := strings.Trim(args_list, " ")
		args := strings.Split(b, ",")
		for i := range args {
			arg := args[i]
			arg = strings.Trim(arg, " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}
		na := len(args)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op}, rp.Operators...)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}

		r := fnc.Function{Name: function_name, NumArgs: na, Operators: rp.Operators}
		rp.Env.AddFunction(r)
		env.Struct = rp
	}
	return result, nil
}

func fMethod(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		method_name := pi.Items[1].Data
		env.CE.StringVars["method_name"] = method_name
		// список операторов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)

		rp.Operators = []*ops.Operator{} // список операторов

		rp.Stack = []*FrameParserStackItem{}
		rp.StackPos = -1
		env.Struct = rp
	case 1:
		body := env.CE.ResultGenerate
		method_name := env.CE.StringVars["method_name"]
		args_list := pi.Items[2].Data

		result = fmt.Sprintf("(method %v (%v) %v)", method_name, args_list, body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op}, rp.Operators...)
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("args")

		b := strings.Trim(args_list, " ")
		args := strings.Split(b, ",")
		for i := range args {
			arg := args[i]
			arg = strings.Trim(arg, " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op}, rp.Operators...)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}

		r := fnc.Method{Name: method_name, Operators: rp.Operators}
		rp.Env.AddMethod(r)
		env.Struct = rp
	}
	return result, nil
}

func fReturn(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// значение - константа или значение переменной
		value := pi.Items[1].Data
		result = fmt.Sprintf("(return %v)", value)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		// вначале надо вычислить аргументы
		b1 := strings.Trim(value, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op := ops.Operator{}
		op.Code = ops.OpName2Code("return")
		a1 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op.Attributes = append(op.Attributes, a1)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func fCallFunction(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		function_name := pi.Items[0].Data
		args_list := pi.Items[1].Data
		result = fmt.Sprintf("(call function %v %v)", function_name, args_list)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)
		/*
			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
			} else {
				rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
			}
		*/
		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		// вначале надо вычислить аргументы
		b1 := strings.Trim(args_list, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op := ops.Operator{}
		op.Code = ops.OpName2Code("call_function")

		b0 := strings.Trim(function_name, " ")
		a, err := ParseArg(b0)
		if err != nil {
			return "", err
		}
		op.Attributes = append(op.Attributes, a)
		a1 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op.Attributes = append(op.Attributes, a1)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

/*
вызов метода состоит из имени объекта либо переменной и через точку само имя метода
метод может быть у любого встроенного типа кроме nil
*/
func fCallMethod(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		method_name := pi.Items[0].Data
		args_list := pi.Items[1].Data
		result = fmt.Sprintf("(call method %v %v)", method_name, args_list)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		// разделяем объект и метод
		b0 := strings.Trim(args_list, " ")
		msl := ns.ParseStringBySignList(b0, []string{"."})
		state := 0
		var obj *attr.Attribute
		var m_name *attr.Attribute
		for i := range msl {
			arg := strings.Trim(msl[i], " ")
			if len(arg) > 0 {
				if arg == "." {
					state = 1
				} else {
					a, err := ParseArg(arg)
					if err != nil {
						return "", err
					}
					if state == 0 {
						// ожидаем объект
						t, _, array := attr.GetAttribute(a)
						switch t {
						case attr.AttrTConst:
							obj = a
						case attr.AttrTArray:
							if len(array) == 2 {
								if array[0] == "?" {
									obj = a
								} else {
									return "", fmt.Errorf("not variable %v", array)
								}
							} else {
								return "", fmt.Errorf("too many items %v", array)
							}
						}
					} else {
						if state == 1 {
							// ожидаем объект
							t, _, array := attr.GetAttribute(a)
							switch t {
							case attr.AttrTConst:
								m_name = a
							case attr.AttrTArray:
								if len(array) == 2 {
									if array[0] == "?" {
										m_name = a
									} else {
										return "", fmt.Errorf("not variable %v", array)
									}
								} else {
									return "", fmt.Errorf("too many items %v", array)
								}
							}
						} else {
							return "", fmt.Errorf("bad type %v", arg)
						}
					}
				}
			}
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}
		// вначале надо вычислить аргументы
		b1 := strings.Trim(args_list, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("call_method")

		op.Attributes = append(op.Attributes, obj)
		op.Attributes = append(op.Attributes, m_name)
		a1 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op.Attributes = append(op.Attributes, a1)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func fCallFunctionWithAssignment(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		function_name := pi.Items[0].Data
		args_list := pi.Items[1].Data
		var_name := pi.Items[3].Data
		result = fmt.Sprintf("%v = (call function %v %v)", var_name, function_name, args_list)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		/*
			op := ops.Operator{}
			op.Code = ops.OpName2Code("line")
			a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
			env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
			op.Attributes = append(op.Attributes, a)

			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
			} else {
				rp.Operators = append(rp.Operators, &op)
			}
		*/
		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)
		/*
			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
			} else {
				rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
			}
		*/
		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		// вначале надо вычислить аргументы
		b1 := strings.Trim(args_list, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op := ops.Operator{}
		op.Code = ops.OpName2Code("call_function")

		b0 := strings.Trim(function_name, " ")
		a, err := ParseArg(b0)
		if err != nil {
			return "", err
		}
		op.Attributes = append(op.Attributes, a)
		a1 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op.Attributes = append(op.Attributes, a1)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op12 := ops.Operator{}
		op12.Code = ops.OpName2Code("set")
		// переменная итератора
		a, err = ParseArg(var_name)
		if err != nil {
			return "", err
		}
		op12.Attributes = append(op12.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op12)
		} else {
			rp.Operators = append(rp.Operators, &op12)
		}

		env.Struct = rp
	}
	return result, nil
}

/*
вызов метода состоит из имени объекта либо переменной и через точку само имя метода
метод может быть у любого встроенного типа кроме nil
*/
func fCallMethodWithAssignment(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		method_name := pi.Items[0].Data
		args_list := pi.Items[1].Data
		result = fmt.Sprintf("(call method %v %v)", method_name, args_list)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		/*
			op := ops.Operator{}
			op.Code = ops.OpName2Code("line")
			a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
			env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
			op.Attributes = append(op.Attributes, a)

			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
			} else {
				rp.Operators = append(rp.Operators, &op)
			}
		*/
		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		// разделяем объект и метод
		b0 := strings.Trim(args_list, " ")
		msl := ns.ParseStringBySignList(b0, []string{"."})
		state := 0
		var obj *attr.Attribute
		var m_name *attr.Attribute
		for i := range msl {
			arg := strings.Trim(msl[i], " ")
			if len(arg) > 0 {
				if arg == "." {
					state = 1
				} else {
					a, err := ParseArg(arg)
					if err != nil {
						return "", err
					}
					if state == 0 {
						// ожидаем объект
						t, _, array := attr.GetAttribute(a)
						switch t {
						case attr.AttrTConst:
							obj = a
						case attr.AttrTArray:
							if len(array) == 2 {
								if array[0] == "?" {
									obj = a
								} else {
									return "", fmt.Errorf("not variable %v", array)
								}
							} else {
								return "", fmt.Errorf("too many items %v", array)
							}
						}
					} else {
						if state == 1 {
							// ожидаем объект
							t, _, array := attr.GetAttribute(a)
							switch t {
							case attr.AttrTConst:
								m_name = a
							case attr.AttrTArray:
								if len(array) == 2 {
									if array[0] == "?" {
										m_name = a
									} else {
										return "", fmt.Errorf("not variable %v", array)
									}
								} else {
									return "", fmt.Errorf("too many items %v", array)
								}
							}
						} else {
							return "", fmt.Errorf("bad type %v", arg)
						}
					}
				}
			}
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}
		// вначале надо вычислить аргументы
		b1 := strings.Trim(args_list, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op := ops.Operator{}
		op.Code = ops.OpName2Code("call_method")

		op.Attributes = append(op.Attributes, obj)
		op.Attributes = append(op.Attributes, m_name)
		a1 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op.Attributes = append(op.Attributes, a1)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func fCallFunctionWithAssignmentMany(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		function_name := pi.Items[0].Data
		args_list := pi.Items[1].Data
		var_names := pi.Items[3].Data
		result = fmt.Sprintf("%v = (call function %v %v)", var_names, function_name, args_list)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		/*
			op := ops.Operator{}
			op.Code = ops.OpName2Code("line")
			a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
			env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
			op.Attributes = append(op.Attributes, a)

			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
			} else {
				rp.Operators = append(rp.Operators, &op)
			}
		*/
		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}
		// вначале надо вычислить аргументы
		b1 := strings.Trim(args_list, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op := ops.Operator{}
		op.Code = ops.OpName2Code("call_function")

		b0 := strings.Trim(function_name, " ")
		a, err := ParseArg(b0)
		if err != nil {
			return "", err
		}
		op.Attributes = append(op.Attributes, a)
		a1 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op.Attributes = append(op.Attributes, a1)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		// вначале надо вычислить аргументы
		b1 = strings.Trim(var_names, " ")
		sl = ns.ParseStringBySignList(b1, []string{","})
		ll = 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("clear")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("set")
						} else {
							op.Code = ops.OpName2Code("clear")
						}
					} else {
						op.Code = ops.OpName2Code("clear")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		env.Struct = rp
	}
	return result, nil
}

func fIf(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список аргументов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)
		si := FrameParserStackItem{}

		rp.Stack = append(rp.Stack, &si)
		rp.StackPos = rp.StackPos + 1

		env.Struct = rp
	case 1:
		env.CE.StringVars["condition"] = env.CE.ResultGenerate
		env.CE.PiCnt = 1
		env.CE.NextState = 2
		env.CE.State = 100
	case 2:
		body := env.CE.ResultGenerate
		cond := env.CE.StringVars["condition"]
		result = fmt.Sprintf("(if %v %v)", cond, body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		/*
			op := ops.Operator{}
			op.Code = ops.OpName2Code("line")
			a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
			env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
			op.Attributes = append(op.Attributes, a)

			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
			} else {
				rp.Operators = append(rp.Operators, &op)
			}
		*/
		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		/*
			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
			} else {
				rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
			}
		*/

		// добавляем оператор условия
		cops := rp.Stack[rp.StackPos].ConditionOps

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, cops...)
		} else {
			rp.Operators = append(rp.Operators, cops...)
		}
		// добавляем переход
		eops := rp.Stack[rp.StackPos].ExecOps
		l := len(eops)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("branch_if_false")
		a := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", l), nil)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}
		// добавляем выполняемые в случае исполнения условия команды
		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, eops...)
		} else {
			rp.Operators = append(rp.Operators, eops...)
		}

		rp.StackPos = rp.StackPos - 1
		if len(rp.Stack) > 0 {
			rp.Stack = rp.Stack[:len(rp.Stack)-1]
		} else {
			rp.Stack = []*FrameParserStackItem{}
		}
		env.Struct = rp
	}
	return result, nil
}

func fCondition1(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// левая часть условия
		l := pi.Items[0].Data
		sign := pi.Items[1].Data
		r := pi.Items[2].Data
		result = fmt.Sprintf("(%v %v %v)", l, sign, r)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		opln := ops.Operator{}
		opln.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		opln.Attributes = append(opln.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &opln)
		} else {
			rp.Operators = append(rp.Operators, &opln)
		}

		// надо понимать, что передалось - переменная или константа.
		al, err := ParseArg(l)
		if err != nil {
			return "", err
		}
		opl := ops.Operator{}

		t, _, array := attr.GetAttribute(al)
		switch t {
		case attr.AttrTConst:
			opl.Code = ops.OpName2Code("const")
		case attr.AttrTArray:
			if len(array) == 2 {
				if array[0] == "?" {
					opl.Code = ops.OpName2Code("get")
				} else {
					opl.Code = ops.OpName2Code("const")
				}
			} else {
				opl.Code = ops.OpName2Code("const")
			}
		}
		opl.Attributes = append(opl.Attributes, al)

		ar, err := ParseArg(r)
		if err != nil {
			return "", err
		}
		opr := ops.Operator{}

		t, _, array = attr.GetAttribute(ar)
		switch t {
		case attr.AttrTConst:
			opr.Code = ops.OpName2Code("const")
		case attr.AttrTArray:
			if len(array) == 2 {
				if array[0] == "?" {
					opr.Code = ops.OpName2Code("get")
				} else {
					opr.Code = ops.OpName2Code("const")
				}
			} else {
				opr.Code = ops.OpName2Code("const")
			}
		}
		opr.Attributes = append(opr.Attributes, ar)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &opr)
		} else {
			rp.Operators = append(rp.Operators, &opr)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &opl)
		} else {
			rp.Operators = append(rp.Operators, &opl)
		}

		op := ops.Operator{}
		switch sign {
		case "==":
			op.Code = ops.OpName2Code("eq")
		case "<":
			op.Code = ops.OpName2Code("lt")
		case ">":
			op.Code = ops.OpName2Code("gt")
		case "=<":
			op.Code = ops.OpName2Code("le")
		case "=>":
			op.Code = ops.OpName2Code("ge")
		}

		rule_name := l + " " + r
		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func fCondition2(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// левая часть условия
		l := pi.Items[0].Data
		r := pi.Items[2].Data
		result = fmt.Sprintf("(gt %v %v)", l, r)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		op := ops.Operator{}
		op.Code = ops.OpName2Code("gt")

		rule_name := l + " " + r
		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func fCondition3(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// левая часть условия
		l := pi.Items[0].Data
		r := pi.Items[2].Data
		result = fmt.Sprintf("(lt %v %v)", l, r)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		op := ops.Operator{}
		op.Code = ops.OpName2Code("lt")

		rule_name := l + " " + r
		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func fForEach(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список аргументов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)
		si := FrameParserStackItem{}

		rp.Stack = append(rp.Stack, &si)
		rp.StackPos = rp.StackPos + 1

		env.Struct = rp

	case 1:
		// поиск
		s1 := pi.Items[3].Data
		s2 := pi.Items[5].Data
		body := env.CE.ResultGenerate

		result = fmt.Sprintf("(for_each (%v = %v) %v)", s2, s1, body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)
		/*
			op := ops.Operator{}
			op.Code = ops.OpName2Code("line")
			a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
			env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
			op.Attributes = append(op.Attributes, a)

			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
			} else {
				rp.Operators = append(rp.Operators, &op)
			}
		*/
		bb_s := []*ops.Operator{}

		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)
		bb_s = append(bb_s, &op_l)
		/*
			if rp.StackPos >= 0 {
				rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
			} else {
				rp.Operators = append(rp.Operators, &op_l)
			}
		*/
		eops := rp.Stack[rp.StackPos].ExecOps
		eops_l := len(eops)

		// find_frame (......) -> stack
		// stack -> create_iterator -> stack
		// stack -> set var iteration_xxx
		// iteration (iteration_xxx) -> stack
		// stack -> set var <var>

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)
		bb_s = append(bb_s, &op_d)

		expr_parser := func(s string) error {
			ssl := parser.ParseArgListFull(s1, env.Output)
			state := 0
			symbol := ""
			args := ""
			for i := range ssl {
				switch ssl[i][0] {
				case "string":
					ss := strings.Trim(ssl[i][1], " \r\n")
					if len(ss) > 0 {
						// Это строка
						if state == 0 {
							state = 3
							symbol = ss
						}
					}
				case "symbols":
					ss := strings.Trim(ssl[i][1], " \r\n")
					if len(ss) > 0 {
						// Это просто символ
						if state == 0 {
							state = 1
							symbol = ss
						}
					}
				case "(":
					ss := strings.Trim(ssl[i][1], " \r\n")
					if len(ss) > 0 {
						// это список
						if state == 0 {
							// список ?
							args = ss
						} else if state == 1 {
							state = 2
							args = ss
						}
					}
				}
			}
			switch state {
			case 0:
				// заменяем на вызов поиска фрейма
				op1 := ops.Operator{}
				op1.Code = ops.OpName2Code("find_frame")

				sl := strings.Split(args, ",")
				for i := range sl {
					arg := strings.Trim(sl[i], " ")
					a, err := ParseArg(arg)
					if err != nil {
						return err
					}
					op1.Attributes = append(op1.Attributes, a)
				}
				bb_s = append(bb_s, &op1)

			case 1, 3:
				// проверяем что это переменная
				if symbol[0] == '?' {
					// строим вызов получения переменной
					op14 := ops.Operator{}
					op14.Code = ops.OpName2Code("get")
					a12, err1 := ParseArg(symbol)
					if err1 != nil {
						return err1
					}
					op14.Attributes = append(op14.Attributes, a12)

					bb_s = append(bb_s, &op14)
				} else {
					// нет это константа
					op1 := ops.Operator{}
					op1.Code = ops.OpName2Code("const")

					a, err := ParseArg(symbol)
					if err != nil {
						return err
					}
					op1.Attributes = append(op1.Attributes, a)

					bb_s = append(bb_s, &op1)
				}
			case 2:
				switch symbol {
				case "найти":
					{
						// заменяем на вызов поиска фрейма
						op1 := ops.Operator{}
						op1.Code = ops.OpName2Code("find_frame")

						sl := strings.Split(args, ",")
						for i := range sl {
							arg := strings.Trim(sl[i], " ")
							a, err := ParseArg(arg)
							if err != nil {
								return err
							}
							op1.Attributes = append(op1.Attributes, a)
						}
						bb_s = append(bb_s, &op1)
					}
				}
			}
			return nil
		}

		err = expr_parser(s1)
		if err != nil {
			return "", err
		}
		// дублируем стек
		op21 := ops.Operator{}
		op21.Code = ops.OpName2Code("dup")

		bb_s = append(bb_s, &op21)

		// проверяем, что на выходе что то есть, а не пустота
		op22 := ops.Operator{}
		op22.Code = ops.OpName2Code("empty")

		bb_s = append(bb_s, &op22)

		// делаем переход на один оператор дальше если не истина (то есть не пусто)
		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("branch_if_false")
		a3 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", 3), nil)
		op3.Attributes = append(op3.Attributes, a3)

		bb_s = append(bb_s, &op3)

		// выбираем результат из стека и переходим на конец цикла
		op31 := ops.Operator{}
		op31.Code = ops.OpName2Code("clear")

		bb_s = append(bb_s, &op31)

		op32 := ops.Operator{}
		op32.Code = ops.OpName2Code("branch")
		a32 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", 1+3+eops_l+3), nil)
		op32.Attributes = append(op32.Attributes, a32)

		bb_s = append(bb_s, &op32)

		// так как результат не пуст то создаем итератор
		op11 := ops.Operator{}
		op11.Code = ops.OpName2Code("create_iterator")

		bb_s = append(bb_s, &op11)

		op12 := ops.Operator{}
		op12.Code = ops.OpName2Code("set")
		// переменная итератора
		sss := fmt.Sprintf("?variable_%v", uqe.UniqueValue(8))
		a, err := ParseArg(sss)
		if err != nil {
			return "", err
		}
		op12.Attributes = append(op12.Attributes, a)

		bb_s = append(bb_s, &op12)

		bb := []*ops.Operator{}

		// эта точка откуда начинется итерация - загружаем в стек из переменной
		op14 := ops.Operator{}
		op14.Code = ops.OpName2Code("get")
		a12, err1 := ParseArg(sss)
		if err1 != nil {
			return "", err1
		}
		op14.Attributes = append(op14.Attributes, a12)

		bb = append(bb, &op14)

		// делаем итерацию
		op13 := ops.Operator{}
		op13.Code = ops.OpName2Code("iteration")

		bb = append(bb, &op13)

		// сохраняем результат в стек
		op2 := ops.Operator{}
		op2.Code = ops.OpName2Code("set")
		a2, err2 := ParseArg(s2)
		if err2 != nil {
			return "", err2
		}
		op2.Attributes = append(op2.Attributes, a2)
		bb = append(bb, &op2)

		// добавляем тело цикла
		bb = append(bb, eops...)

		op15 := ops.Operator{}
		op15.Code = ops.OpName2Code("get")
		a13, err1 := ParseArg(sss)
		if err1 != nil {
			return "", err1
		}
		op15.Attributes = append(op15.Attributes, a13)

		bb = append(bb, &op15)

		// проверяем что итератор не кончился - в стеке true если все кончилось
		op131 := ops.Operator{}
		op131.Code = ops.OpName2Code("check_iteration")

		bb = append(bb, &op131)

		op33 := ops.Operator{}
		op33.Code = ops.OpName2Code("branch_if_false")
		a31 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", -(len(bb))), nil)
		op33.Attributes = append(op33.Attributes, a31)
		bb = append(bb, &op33)

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, bb_s...)
		} else {
			rp.Operators = append(rp.Operators, bb_s...)
		}

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, bb...)
		} else {
			rp.Operators = append(rp.Operators, bb...)
		}

		rp.StackPos = rp.StackPos - 1
		if len(rp.Stack) > 0 {
			rp.Stack = rp.Stack[:len(rp.Stack)-1]
		} else {
			rp.Stack = []*FrameParserStackItem{}
		}

		env.Struct = rp
	}
	return result, nil
}

func fWhile(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список аргументов
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(FrameParser)
		si := FrameParserStackItem{}

		rp.Stack = append(rp.Stack, &si)
		rp.StackPos = rp.StackPos + 1

		env.Struct = rp

	case 1:
		env.CE.StringVars["condition"] = env.CE.ResultGenerate
		env.CE.PiCnt = 1
		env.CE.NextState = 2
		env.CE.State = 100

	case 2:
		body := env.CE.ResultGenerate
		cond := env.CE.StringVars["condition"]
		result = fmt.Sprintf("(while (%v) %v)", cond, body)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		bb_s := []*ops.Operator{}

		opl := ops.Operator{}
		opl.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		opl.Attributes = append(opl.Attributes, a)
		bb_s = append(bb_s, &opl)

		eops := rp.Stack[rp.StackPos].ExecOps
		eops_l := len(eops)

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)
		bb_s = append(bb_s, &op_d)

		// !!!! начало циклически исполняемого блока
		bb := []*ops.Operator{}

		// добавляем оператор условия
		cops := rp.Stack[rp.StackPos].ConditionOps
		cops_l := len(cops)

		bb = append(bb, cops...)

		// делаем переход в конец если не истина (то есть не пусто)
		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("branch_if_false")
		a3 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", 1+eops_l+1), nil)
		op3.Attributes = append(op3.Attributes, a3)
		bb = append(bb, &op3)

		// добавляем тело цикла
		bb = append(bb, eops...)

		op32 := ops.Operator{}
		op32.Code = ops.OpName2Code("branch")
		a32 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", -(1+eops_l+cops_l)), nil)
		op32.Attributes = append(op32.Attributes, a32)
		bb = append(bb, &op32)

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, bb_s...)
		} else {
			rp.Operators = append(rp.Operators, bb_s...)
		}

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, bb...)
		} else {
			rp.Operators = append(rp.Operators, bb...)
		}

		rp.StackPos = rp.StackPos - 1
		if len(rp.Stack) > 0 {
			rp.Stack = rp.Stack[:len(rp.Stack)-1]
		} else {
			rp.Stack = []*FrameParserStackItem{}
		}

		env.Struct = rp
	}
	return result, nil
}

func fRelation(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// Отношение
		s1 := pi.Items[0].Data
		s2 := pi.Items[1].Data
		s3 := pi.Items[2].Data
		result = fmt.Sprintf(" %v %v %v", s1, s2, s3)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("relation")

		for i := 0; i < 3; i++ {
			arg := pi.Items[i].Data
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func fFrame(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// Отношение
		s := pi.Items[1].Data
		result = fmt.Sprintf("(frame %v)", s)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("frame")

		sl := strings.Split(s, ",")
		for i := range sl {
			arg := strings.Trim(sl[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func fFrameWithAssignment(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// Отношение
		s := pi.Items[1].Data
		s2 := pi.Items[3].Data
		result = fmt.Sprintf("%v = frame(%v)", s2, s)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append([]*ops.Operator{&op_d}, rp.Stack[rp.StackPos].ExecOps...)
		} else {
			rp.Operators = append([]*ops.Operator{&op_d}, rp.Operators...)
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("frame")

		sl := strings.Split(s, ",")
		for i := range sl {
			arg := strings.Trim(sl[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op2 := ops.Operator{}
		op2.Code = ops.OpName2Code("set")
		sa := strings.Split(s2, ",")
		for i := range sa {
			arg := strings.Trim(sa[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op2.Attributes = append(op2.Attributes, a)
		}
		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op2)
		} else {
			rp.Operators = append(rp.Operators, &op2)
		}

		env.Struct = rp
	}
	return result, nil
}

func fFindWithAssignment(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// поиск
		s1 := pi.Items[0].Data
		s2 := pi.Items[2].Data
		result = fmt.Sprintf("(%v = find_frame_with_assignment %v)", s2, s1)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		op1 := ops.Operator{}
		op1.Code = ops.OpName2Code("find_frame")

		sl := ns.ParseStringBySignList(s1, []string{","})
		for i := range sl {
			arg := strings.Trim(sl[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op1.Attributes = append(op1.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op1)
		} else {
			rp.Operators = append(rp.Operators, &op1)
		}

		op2 := ops.Operator{}
		op2.Code = ops.OpName2Code("set")
		sa := strings.Split(s2, ",")
		for i := range sa {
			arg := strings.Trim(sa[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op2.Attributes = append(op2.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op2)
		} else {
			rp.Operators = append(rp.Operators, &op2)
		}

		env.Struct = rp
	}
	return result, nil
}

func fFindAndAdd(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// поиск
		s1 := pi.Items[0].Data
		s2 := pi.Items[2].Data
		result = fmt.Sprintf("(%v = f_find_and_add %v)", s2, s1)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		op1 := ops.Operator{}
		op1.Code = ops.OpName2Code("find_frame")

		sa := ns.ParseStringBySignList(s1, []string{","})
		for i := range sa {
			arg := strings.Trim(sa[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op1.Attributes = append(op1.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op1)
		} else {
			rp.Operators = append(rp.Operators, &op1)
		}

		// дублируем стек
		op21 := ops.Operator{}
		op21.Code = ops.OpName2Code("dup")

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op21)
		} else {
			rp.Operators = append(rp.Operators, &op21)
		}

		// проверяем, что на выходе что то есть, а не пустота
		op22 := ops.Operator{}
		op22.Code = ops.OpName2Code("empty")

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op22)
		} else {
			rp.Operators = append(rp.Operators, &op22)
		}

		// делаем переход на один оператор дальше если не истина (то есть не пусто)
		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("branch_if_false")
		a3 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", 3), nil)
		op3.Attributes = append(op3.Attributes, a3)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op3)
		} else {
			rp.Operators = append(rp.Operators, &op3)
		}

		// выбираем результат из стека и переходим на конец цикла
		op31 := ops.Operator{}
		op31.Code = ops.OpName2Code("clear")

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op31)
		} else {
			rp.Operators = append(rp.Operators, &op31)
		}

		op32 := ops.Operator{}
		op32.Code = ops.OpName2Code("branch")
		a32 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", 1+1), nil)
		op32.Attributes = append(op32.Attributes, a32)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op32)
		} else {
			rp.Operators = append(rp.Operators, &op32)
		}

		// в стеке то к чему надо добавить
		op2 := ops.Operator{}
		op2.Code = ops.OpName2Code("add_slots")
		sa = strings.Split(s2, ",")
		for i := range sa {
			arg := strings.Trim(sa[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op2.Attributes = append(op2.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op2)
		} else {
			rp.Operators = append(rp.Operators, &op2)
		}

		env.Struct = rp
	}
	return result, nil
}

func fUnify(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// переменная
		s1 := pi.Items[0].Data
		s2 := pi.Items[2].Data
		s3 := pi.Items[4].Data
		result = fmt.Sprintf("(%v = unify %v by %v)", s3, s2, s1)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		op1 := ops.Operator{}
		op1.Code = ops.OpName2Code("find_frame")

		sl := strings.Split(s1, ",")
		for i := range sl {
			arg := strings.Trim(sl[i], " ")
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op1.Attributes = append(op1.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op1)
		} else {
			rp.Operators = append(rp.Operators, &op1)
		}

		op2 := ops.Operator{}
		op2.Code = ops.OpName2Code("unify")
		a, err = ParseArg(s2)
		if err != nil {
			return "", err
		}
		op2.Attributes = append(op2.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op2)
		} else {
			rp.Operators = append(rp.Operators, &op2)
		}

		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("set")
		a2, err2 := ParseArg(s3)
		if err2 != nil {
			return "", err2
		}
		op3.Attributes = append(op3.Attributes, a2)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op3)
		} else {
			rp.Operators = append(rp.Operators, &op3)
		}
		env.Struct = rp

	}
	return result, nil
}

func fBreak(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		result = "(break)"
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("break")

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp

	}
	return result, nil
}

func fContinue(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		result = "(continue)"
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op := ops.Operator{}
		op.Code = ops.OpName2Code("line")
		a := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		op = ops.Operator{}
		op.Code = ops.OpName2Code("continue")

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Operators = append(rp.Operators, &op)
		}

		env.Struct = rp

	}
	return result, nil
}

func fSet(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// переменная
		s1 := pi.Items[0].Data
		s2 := pi.Items[2].Data
		result = fmt.Sprintf("(%v = %v)", s1, s2)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		op1 := ops.Operator{}
		op1.Code = ops.OpName2Code("const")

		a, err := ParseArg(s1)
		if err != nil {
			return "", err
		}
		op1.Attributes = append(op1.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op1)
		} else {
			rp.Operators = append(rp.Operators, &op1)
		}

		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("set")
		a2, err2 := ParseArg(s2)
		if err2 != nil {
			return "", err2
		}
		op3.Attributes = append(op3.Attributes, a2)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op3)
		} else {
			rp.Operators = append(rp.Operators, &op3)
		}
		env.Struct = rp

	}
	return result, nil
}

func fSetList(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// переменная
		s1 := pi.Items[0].Data
		s2 := pi.Items[2].Data
		result = fmt.Sprintf("([%v] = %v)", s1, s2)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op_l := ops.Operator{}
		op_l.Code = ops.OpName2Code("line")
		a_l := &attr.Attribute{Type: attr.AttrTNumber, Number: pi.Items[0].LineNumBegin}
		//env.Output.Print("pi.Items[0].LineNumBegin %v", pi.Items[0].LineNumBegin)
		op_l.Attributes = append(op_l.Attributes, a_l)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_l)
		} else {
			rp.Operators = append(rp.Operators, &op_l)
		}

		op_d := ops.Operator{}
		op_d.Code = ops.OpName2Code("debug")
		a1_d, err := ParseArg("text")
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a1_d)
		a2_d, err := ParseArg(result)
		if err != nil {
			return "", err
		}
		op_d.Attributes = append(op_d.Attributes, a2_d)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op_d)
		} else {
			rp.Operators = append(rp.Operators, &op_d)
		}

		// вначале надо вычислить список
		b1 := strings.Trim(s1, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		ll := 0
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op := ops.Operator{}

				t, _, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					op.Code = ops.OpName2Code("const")
					op.Attributes = append(op.Attributes, a)
				case attr.AttrTArray:
					if len(array) == 2 {
						if array[0] == "?" {
							op.Code = ops.OpName2Code("get")
						} else {
							op.Code = ops.OpName2Code("const")
						}
					} else {
						op.Code = ops.OpName2Code("const")
					}
					op.Attributes = append(op.Attributes, a)
				}
				ll = ll + 1
				if rp.StackPos >= 0 {
					rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
				} else {
					rp.Operators = append(rp.Operators, &op)
				}
			}
		}

		op1 := ops.Operator{}
		op1.Code = ops.OpName2Code("slice")
		a32 := attr.NewAttribute(attr.AttrTNumber, fmt.Sprintf("%v", ll), nil)
		op1.Attributes = append(op1.Attributes, a32)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op1)
		} else {
			rp.Operators = append(rp.Operators, &op1)
		}

		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("set")
		a2, err2 := ParseArg(s2)
		if err2 != nil {
			return "", err2
		}
		op3.Attributes = append(op3.Attributes, a2)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op3)
		} else {
			rp.Operators = append(rp.Operators, &op3)
		}
		env.Struct = rp

	}
	return result, nil
}

func fSymbol(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// символ
		s := pi.Items[0].Data
		result = fmt.Sprintf(" %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func fString(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// строка
		s := pi.Items[0].Data
		result = fmt.Sprintf(" %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func fVariable(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// переменная
		s := pi.Items[1].Data
		result = fmt.Sprintf("?%v", s)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op3 := ops.Operator{}
		op3.Code = ops.OpName2Code("get")
		a2, err2 := ParseArg(s)
		if err2 != nil {
			return "", err2
		}
		op3.Attributes = append(op3.Attributes, a2)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &op3)
		} else {
			rp.Operators = append(rp.Operators, &op3)
		}

	}
	return result, nil
}

func fConst(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// символ
		s := pi.Items[0].Data
		result = fmt.Sprintf("%v", s)
		env.CE.State = 1000

		rp := env.Struct.(FrameParser)

		op1 := ops.Operator{}
		op1.Code = ops.OpName2Code("const")

		a, err := ParseArg(s)
		if err != nil {
			return "", err
		}
		op1.Attributes = append(op1.Attributes, a)

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op1)
		} else {
			rp.Operators = append(rp.Operators, &op1)
		}

	}
	return result, nil
}

func MakeRules(env *parser.Env) {
	if true {
		defer func() {
			r := recover()
			if r != nil {
				env.Output.Print("%v\r\n", r)
				return
			}
		}()
	}
	items := []string{"определение фреймА", "присваивание результатат поиска",
		"добавление слотов в результататы поиска", "присваивание результатат поиска созданного из списка",
		"присваивание списка слотов найденного фрейма", "вернуть", "унификация",
		"если", "если - иначе", "вызов функции", "для каждого элемента", "вызов метода", "прервать",
		"продолжить", "присвоить", "фрейм присвоить переменной", "присвоить_список",
		"вызов функции с присваиванием", "вызов метода с присваиванием", "пока",
		"вызов функции с присваиванием нескольких значений",
	}

	//<symbols, == отношения> <{, > - добавление отношений во фреймы
	gr := parser.MakeRule("добавление отношений во фреймы", env)
	gr.AddItemToRule("symbols", "", 1, "отношения", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"отношение"}, env)
	gr.AddRuleHandler(fRelations, env)

	//<symbols, > <symbols, > <symbols, > - отношение
	gr = parser.MakeRule("отношение", env)
	gr.AddItemToRule("symbols|string", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols|string", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols|string", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fRelation, env)

	//<symbols, == фреймы> <{, > - определение фреймов
	gr = parser.MakeRule("определение фреймов", env)
	gr.AddItemToRule("symbols", "", 1, "фреймы", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"определение фреймА"}, env)
	gr.AddRuleHandler(fFrames, env)

	//<symbols, == фрейм> <(, > - определение фреймА, после ключевого слова идет список
	gr = parser.MakeRule("определение фреймА", env)
	gr.AddItemToRule("symbols", "", 1, "фрейм", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", ";", []string{"список_аргументов"}, env) // , "список"
	gr.AddRuleHandler(fFrame, env)

	//<symbols, == фрейм> <(, > - определение фреймА, после ключевого слова идет список
	gr = parser.MakeRule("фрейм присвоить переменной", env)
	gr.AddItemToRule("symbols", "", 1, "фрейм", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env) // , "список"
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "[0]", 1, "?", ";", []string{}, env)
	gr.AddRuleHandler(fFrameWithAssignment, env)

	//<(, > <symbols, == : > <?, > <(, > - поиск фрейма либо фреймов и добавление слотов и значений
	gr = parser.MakeRule("добавление слотов в результататы поиска", env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env)
	gr.AddItemToRule("symbols", "", 1, ":", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", ";", []string{"список_аргументов"}, env)
	gr.AddRuleHandler(fFindAndAdd, env)

	//<(, > <symbols, == => > <?, > <symbols,> - присваивание результатат поиска критерии поиска задаются либо списком слотов и значений либо списком объявлений слотов и значений
	gr = parser.MakeRule("присваивание результатат поиска", env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "[0]", 1, "?", ";", []string{}, env)
	gr.AddRuleHandler(fFindWithAssignment, env)

	//<(, > <symbols, == ?> <symbols,> <symbols, == => > <symbols, == ?> <symbols,> - унификация (сущность.объект)?элемент-класс
	gr = parser.MakeRule("унификация", env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env)
	gr.AddItemToRule("symbols", "", 1, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fUnify, env)

	//<symbols, == если> <(, >  <{, > - если
	gr = parser.MakeRule("если", env)
	gr.AddItemToRule("symbols", "", 1, "если", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"условие1", "условие2", "условие3"}, env)
	gr.AddItemToRule("{", "", 0, "", ";", items, env)
	gr.AddRuleHandler(fIf, env)

	//<symbols, == для> <symbols, == каждого> <(, >  <symbols, == => > <symbols, == ?> <symbols,> <{, > - цикл для каждого фрейма
	gr = parser.MakeRule("для каждого элемента", env)
	gr.AddItemToRule("symbols", "", 1, "для", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "каждого", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "элемента", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env) // "переменная", "вызов функции", "константа" // "список_аргументов"
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", items, env)
	gr.AddRuleHandler(fForEach, env)

	//<symbols, == пока> <(, >  <symbols, == => > <symbols, == ?> <symbols,> <{, > - цикл для каждого фрейма
	gr = parser.MakeRule("пока", env)
	gr.AddItemToRule("symbols", "", 1, "пока", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"условие1", "условие2", "условие3"}, env) // "переменная", "вызов функции", "константа" // "список_аргументов"
	gr.AddItemToRule("{", "", 0, "", ";", items, env)
	gr.AddRuleHandler(fWhile, env)

	//<symbols, == @> <symbols, > <(, > <{, > -  определение метода
	gr = parser.MakeRule("определение функции", env)
	gr.AddItemToRule("symbols", "", 1, "функция", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", items, env)
	gr.AddRuleHandler(fFunction, env)

	//<symbols, == @> <symbols, > <(, > <{, > -  определение метода
	gr = parser.MakeRule("определение метода", env)
	gr.AddItemToRule("symbols", "", 1, "метод", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", items, env)
	gr.AddRuleHandler(fMethod, env)

	//<symbols, > <(, > - вернуть
	gr = parser.MakeRule("вернуть", env)
	gr.AddItemToRule("symbols", "", 1, "вернуть", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fReturn, env)

	//<symbols, > <(, > - вызов функции
	gr = parser.MakeRule("вызов функции", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fCallFunction, env)

	//<symbols, > <(, > - вызов метода
	gr = parser.MakeRule("вызов метода", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fCallMethod, env)

	//<symbols, > <(, > - вызов функции с присваиванием
	gr = parser.MakeRule("вызов функции с присваиванием", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fCallFunctionWithAssignment, env)

	//<symbols, > <(, > - вызов метода с присваиванием
	gr = parser.MakeRule("вызов метода с присваиванием", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fCallMethodWithAssignment, env)

	//<symbols, > <(, > - вызов функции с присваиванием нескольких значений
	gr = parser.MakeRule("вызов функции с присваиванием нескольких значений", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fCallFunctionWithAssignmentMany, env)

	//<symbols, == прервать> - прервать выполнение
	gr = parser.MakeRule("прервать", env)
	gr.AddItemToRule("symbols", "", 1, "прервать", "", []string{}, env)
	gr.AddRuleHandler(fBreak, env)

	//<symbols, == продолжить> - продолжить выполнение пропустив часть блока
	gr = parser.MakeRule("продолжить", env)
	gr.AddItemToRule("symbols", "", 1, "продолжить", "", []string{}, env)
	gr.AddRuleHandler(fContinue, env)

	//<string, > <symbols, == => > <symbols,> - константу строку в переменную
	gr = parser.MakeRule("присвоить", env)
	gr.AddItemToRule("symbols|string", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fSet, env)

	//<[, > <symbols, == => > <symbols,> - константу строку в переменную
	gr = parser.MakeRule("присвоить_список", env)
	gr.AddItemToRule("[", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "=>", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", ";", []string{}, env)
	gr.AddRuleHandler(fSetList, env)

	//<symbols, > <symbols, == == >  <строка> - условие 1 if
	gr = parser.MakeRule("условие1", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "==", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fCondition1, env)

	//<symbols, > <symbols, == > >  <symbols, > - условие 2 if
	gr = parser.MakeRule("условие2", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, ">", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fCondition2, env)

	//<symbols, > <symbols, == < >  <symbols, > - условие 3 if
	gr = parser.MakeRule("условие3", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 1, "<", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fCondition3, env)

	// среднеуровневые элементы
	// список в определении тринара или шаблона
	// <symbols, > - просто символ
	gr = parser.MakeRule("символ", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fSymbol, env)

	// <string, > - просто строка
	gr = parser.MakeRule("строка", env)
	gr.AddItemToRule("string", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fString, env)

	// ?<variable name>
	// <symbols, == ?> - переменная
	gr = parser.MakeRule("переменная", env)
	gr.AddItemToRule("symbols", "", 0, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fVariable, env)
	/*
		//<symbols, > <(, > - вызов функции
		gr = parser.MakeRule("вызов функции", env)
		gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
		gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
		gr.AddRuleHandler(f_call_function, env)
	*/
	// <string, > - просто строка
	gr = parser.MakeRule("константа", env)
	gr.AddItemToRule("symbols|string", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fConst, env)

	// ?<variable name>:<attribute name>
	// <symbols, == ?> <symbols, > <symbols, == :> <symbols, >- атрибут переменной
	gr = parser.MakeRule("атрибут переменной", env)
	gr.AddItemToRule("symbols", "", 0, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, ":", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)

	// ?<variable name>:<attribute name>
	// <symbols, > <symbols, == :> <symbols, >- атрибут
	gr = parser.MakeRule("атрибут", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, ":", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)

	high_level_array := []string{"добавление отношений во фреймы", "определение фреймов", "определение функции", "определение метода"}

	expr_array := []string{"атрибут переменной", "строка", "символ"}

	env.SetHLAEnv(high_level_array)
	env.SetEAEnv(expr_array)
	env.SetBGRAEnv()
}
