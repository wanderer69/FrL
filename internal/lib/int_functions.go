package frl

//	"github.com/wanderer69/FrL/src/lib/common"

// функции интерпретатора
// встроенные
type InternalFunction struct {
	Name    string
	NumArgs int
	Args    []*Value
	Return  []*Value
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

func Print_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "печатать"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: -1} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			for i := range if_.Args {
				ss, ok := FromType(if_.Args[i])
				if ok {
					ie.Output.Print("%v\t", ss)
				}
			}
			ie.Output.Print("\r\n")
			return nil, true, nil
		}
	}
	return nil, false, nil
}

// встроенные методы
// int - методы + - / * string
// float - методы
// string - методы + slice trim split integer float
// slot - методы value_get property_get

func AddNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "сложить"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Add(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func SubNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "вычесть"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Sub(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func MulNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "умножить"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Mul(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func DivNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "делить"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Div(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func FromStringNumber_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "из_строки"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			err := if_.Args[0].FromString(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, if_.Args[0])
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func ConcatString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "склеить"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Concat(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func SliceString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "срез"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SliceString(if_.Args[1], if_.Args[2])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func TrimString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "обрезать"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Trim(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func SplitString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "отрезать"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].Split(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func FromNumberString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "из_числа"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].FromNumber(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func GetNameSlot_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_имя_слота"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SlotGetName()
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func GetValueSlot_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_значение_слота"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SlotGetValue()
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func GetPropertySlot_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_свойство_слота"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SlotGetProperty()
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func ItemSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "элемент"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номер элемента
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SliceItem(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func SliceSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "слайс"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SliceSlice(if_.Args[1], if_.Args[2])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func InsertSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "вставить"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SliceInsert(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func AppendSlice_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "добавить"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].SliceAppend(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func CreateStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "поток"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := StreamCreate(if_.Args[0])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func OpenStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "открыть_поток"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			err := if_.Args[0].StreamOpen()
			if err != nil {
				return nil, false, err
			}
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func ReadStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "читать_поток"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			cnt, v, err := if_.Args[0].StreamRead()
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			if_.Return = append(if_.Return, cnt)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func WriteStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "записать_поток"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			err := if_.Args[0].StreamWrite(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func CloseStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "закрыть_поток"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			err := if_.Args[0].StreamClose()
			if err != nil {
				return nil, false, err
			}
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func ControlSetStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "установить_настройки_потока"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			err := if_.Args[0].StreamControlSet(if_.Args[1], if_.Args[2])
			if err != nil {
				return nil, false, err
			}
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func ControlGetStream_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "получить_настройки_потока"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход слайс и номера начала и конца
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v, err := if_.Args[0].StreamControlGet(if_.Args[1])
			if err != nil {
				return nil, false, err
			}
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func SprintfString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "форматировать"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: -1} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			var fmt *Value
			fmt = nil
			args := []*Value{}
			for i := range if_.Args {
				if i == 0 {
					fmt = if_.Args[i]
				} else {
					args = append(args, if_.Args[i])
				}
			}
			if fmt != nil {
				v, err := fmt.SprintfString(args...)
				if err != nil {
					return nil, false, err
				}
				if_.Return = append(if_.Return, v)
				return if_, true, nil
			}
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func IsType_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "тип"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 1} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v := if_.Args[0].IsType()
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func UUID_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "уникальный_идентификатор"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 0} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			v := UUIDString()
			if_.Return = append(if_.Return, v)
			return if_, true, nil
		}
	}
	return nil, false, nil
}

func AddSlotFrame_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "добавить_слот"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			if_.Args[0].FrameAddSlot(if_.Args[1])
			// if_.Return = append(if_.Return, v)
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func SetSlotFrame_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "добавить_значение_в_слот"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 3} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			if_.Args[0].FrameSetSlot(if_.Args[1], if_.Args[2])
			// if_.Return = append(if_.Return, v)
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func DeleteSlotFrame_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "удалить_слот"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			if_.Args[0].FrameDeleteSlot(if_.Args[1])
			// if_.Return = append(if_.Return, v)
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func EvalString_internal(ie *InterpreterEnv, state int, if_ *InternalFunction) (*InternalFunction, bool, error) {
	// принцип аналогичен команде однако есть отличие так как вычисление идет в две итерации
	// 0. регистрация
	// 1. оценка и связывание аргументов
	// 2. собственно вычисление
	switch state {
	case 0:
		if_n := &InternalFunction{Name: "оценить_выражение"} // имя
		return if_n, false, nil
	case 1:
		if_n := &InternalFunction{NumArgs: 2} // принимает на вход список
		return if_n, false, nil
	case 2:
		if if_ != nil {
			_, err := if_.Args[0].EvalString(ie)
			// if_.Return = append(if_.Return, v)
			return nil, false, err
		}
	}
	return nil, false, nil
}
