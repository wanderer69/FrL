package frl

import "fmt"

//	"github.com/wanderer69/FrL/src/lib/common"

// функции интерпретатора
// встроенные
type InternalFunction struct {
	Name    string
	NumArgs int
	Args    []*Value
	Return  []*Value
}

// внешние
type ExternalFunction struct {
	Name    string
	NumArgs int
	Func    func(args []interface{}) ([]interface{}, bool, error)

	Args   []*Value
	Return []*Value
}

// методы объектов интерпретатора
// встроенные
type InternalMethod struct {
	Name    string
	Type    string // применимость метода к объектам
	NumArgs int
	Args    []*Value
	Return  []*Value
}

func Print_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "печатать"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: -1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			for i := range args {
				ss, ok := FromType(args[i])
				if ok {
					ie.Output.Print("%v\t", ss)
				}
			}
			ie.Output.Print("\r\n")
			return nil, nil, true, nil
		}
	}
	return nil, nil, false, nil
}

// встроенные методы
// int - методы + - / * string
// float - методы
// string - методы + slice trim split integer float
// slot - методы value_get property_get

func AddNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "сложить"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Add(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			//result = append(result, v)
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func SubNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "вычесть"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Sub(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func MulNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "умножить"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Mul(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func DivNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "делить"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Div(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func FromStringNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "из_строки"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			err := args[0].FromString(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, args[0])
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func ConcatString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "склеить"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Concat(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func SliceString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "срез"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SliceString(args[1], args[2])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func TrimString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "обрезать"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Trim(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func SplitString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "отрезать"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].Split(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func FromNumberString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "из_числа"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].FromNumber(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func GetNameSlot_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_имя_слота"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SlotGetName()
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func GetValueSlot_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_значение_слота"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SlotGetValue()
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func GetPropertySlot_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_свойство_слота"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SlotGetProperty()
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func ItemSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "элемент"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номер элемента
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SliceItem(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func SliceSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "слайс"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SliceSlice(args[1], args[2])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func InsertSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "вставить"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SliceInsert(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func AppendSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "добавить"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].SliceAppend(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func CreateStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "поток"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := StreamCreate(args[0])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func OpenStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "открыть_поток"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			//result := []*Value{}
			err := args[0].StreamOpen()
			if err != nil {
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func ReadStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "читать_поток"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			cnt, v, err := args[0].StreamRead()
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			result = append(result, cnt)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func WriteStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "записать_поток"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			err := args[0].StreamWrite(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func CloseStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "закрыть_поток"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			err := args[0].StreamClose()
			if err != nil {
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func ControlSetStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "установить_настройки_потока"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			err := args[0].StreamControlSet(args[1], args[2])
			if err != nil {
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func ControlGetStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_настройки_потока"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v, err := args[0].StreamControlGet(args[1])
			if err != nil {
				return nil, nil, false, err
			}
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func SprintfString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "форматировать"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: -1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			var fmt *Value
			fmt = nil
			args_lst := []*Value{}
			for i := range args {
				if i == 0 {
					fmt = args[i]
				} else {
					args_lst = append(args_lst, args[i])
				}
			}
			if fmt != nil {
				v, err := fmt.SprintfString(args_lst...)
				if err != nil {
					return nil, nil, false, err
				}
				result = append(result, v)
				return if_, result, true, nil
			}
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func IsType_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "тип"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v := args[0].IsType()
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func UUID_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "уникальный_идентификатор"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 0} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			v := UUIDString()
			result = append(result, v)
			return if_, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func AddSlotFrame_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "добавить_слот"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			args[0].FrameAddSlot(args[1])
			// result = append(result, v)
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func SetSlotFrame_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "добавить_значение_в_слот"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			args[0].FrameSetSlot(args[1], args[2])
			// result = append(result, v)
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func DeleteSlotFrame_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "удалить_слот"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			args[0].FrameDeleteSlot(args[1])
			// result = append(result, v)
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func EvalString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "оценить_выражение"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			_, err := args[0].EvalString(ie)
			// result = append(result, v)
			return nil, nil, false, err
		}
	}
	return nil, nil, false, nil
}

func OpenDataBase_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "открыть_базу_данных"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			if args[0].GetType() != VtString {
				return nil, nil, false, fmt.Errorf("bad type database name %v", args[0].GetType())
			}
			pathToDB := args[0].String()

			db := NewDataBase()
			err := db.Connect(DataBaseTypeSimple, pathToDB, ie.Output)
			if err != nil {
				return nil, nil, false, err
			}
			v := CreateValue(db)
			result = append(result, v)
			return nil, result, true, nil
		}
	}
	return nil, nil, false, nil
}

func CloseDataBase_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "закрыть_базу_данных"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			//result := []*Value{}
			if args[0].GetType() != VtDataBase {
				return nil, nil, false, fmt.Errorf("must be type database, has %v", args[0].GetType())
			}
			db := args[0].DataBase()
			db.Close()

			//result = append(result, v)
			return nil, nil, false, nil
		}
	}
	return nil, nil, false, nil
}

func FindInDataBase_internal(ie *InterpreterEnv, state int, if_ *InternalFunction, args []*Value) (*InternalFunction, []*Value, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "найти_в_базе_данных"} // имя
		return if_n, nil, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, nil, false, nil
	case 2:
		if if_ != nil {
			result := []*Value{}
			if args[0].GetType() != VtString {
				return nil, nil, false, fmt.Errorf("bad type database name %v", args[0].GetType())
			}
			pathToDB := args[0].String()

			db := NewDataBase()
			err := db.Connect(DataBaseTypeSimple, pathToDB, ie.Output)
			if err != nil {
				return nil, nil, false, err
			}
			v := CreateValue(db)
			result = append(result, v)
			return nil, result, true, nil
		}
	}
	return nil, nil, false, nil
}
