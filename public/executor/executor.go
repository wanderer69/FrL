package executor

import (
	"fmt"
	"os"
	"time"

	"github.com/wanderer69/FrL/public/entity"
	frl "github.com/wanderer69/FrL/public/lib"
	print "github.com/wanderer69/tools/parser/print"
)

type ExecutorBase struct {
	fe     *frl.FrameEnvironment
	output *print.Output
}

type Executor struct {
	eb    *ExecutorBase
	ie    *frl.InterpreterEnv
	debug int
	// output *print.Output
	//	tickerPool   []*time.Ticker
	eventManager *frl.EventManager
}

func NewExecutorBase(debug int, output *print.Output) *ExecutorBase {
	return &ExecutorBase{
		// настраиваем окружение
		fe:     frl.NewFrameEnvironment(),
		output: output,
	}
}

func InitExecutor(
	eb *ExecutorBase,
	extFunctions map[string]func(args []*frl.Value) ([]*frl.Value, bool, error),
	output *print.Output,
	outputTranslate *print.Output,
	debug int,
	eventManager *frl.EventManager,
) *Executor {
	ie := frl.NewInterpreterEnv()
	ie.SetDebug(debug)
	ie.BindFunction(frl.Print_internal)
	ie.BindFunction(frl.AddNumber_internal)
	ie.BindFunction(frl.SubNumber_internal)
	ie.BindFunction(frl.MulNumber_internal)
	ie.BindFunction(frl.DivNumber_internal)
	ie.BindFunction(frl.FromStringNumber_internal)
	ie.BindFunction(frl.ConcatString_internal)
	ie.BindFunction(frl.SliceString_internal)
	ie.BindFunction(frl.TrimString_internal)
	ie.BindFunction(frl.SplitString_internal)
	ie.BindFunction(frl.FromNumberString_internal)
	ie.BindFunction(frl.GetNameSlot_internal)
	ie.BindFunction(frl.GetValueSlot_internal)
	ie.BindFunction(frl.GetPropertySlot_internal)
	ie.BindFunction(frl.ItemSlice_internal)
	ie.BindFunction(frl.SliceSlice_internal)
	ie.BindFunction(frl.InsertSlice_internal)
	ie.BindFunction(frl.AppendSlice_internal)
	ie.BindFunction(frl.CreateStream_internal)
	ie.BindFunction(frl.OpenStream_internal)
	ie.BindFunction(frl.ReadStream_internal)
	ie.BindFunction(frl.WriteStream_internal)
	ie.BindFunction(frl.CloseStream_internal)
	ie.BindFunction(frl.ControlSetStream_internal)
	ie.BindFunction(frl.ControlGetStream_internal)
	ie.BindFunction(frl.SprintfString_internal)
	ie.BindFunction(frl.IsType_internal)
	ie.BindFunction(frl.UUID_internal)
	ie.BindFunction(frl.AddSlotFrame_internal)
	ie.BindFunction(frl.SetSlotFrame_internal)
	ie.BindFunction(frl.DeleteSlotFrame_internal)
	ie.BindFunction(frl.EvalString_internal)
	ie.BindFunction(frl.OpenDataBase_internal)
	ie.BindFunction(frl.FindInDataBase_internal)
	ie.BindFunction(frl.CloseDataBase_internal)
	ie.BindFunction(frl.StoreInDataBase_internal)
	ie.BindFunction(frl.SetTimerEvent_internal)
	ie.BindFunction(frl.SetChannelEvent_internal)
	ie.BindFunction(frl.FireEvent_internal)
	ie.BindFunction(frl.GetSlot_internal)
	ie.BindFunction(frl.DoneEvent_internal)
	ie.ExternalFunctions = extFunctions

	ie.SetFrameEnvironment(eb.fe)
	ie.FE = eb.fe

	ie.Output = output
	ie.OutputTranslate = outputTranslate

	return &Executor{
		eb:           eb,
		debug:        debug,
		ie:           ie,
		eventManager: eventManager,
	}
}

func (e *Executor) Exec(fileIn string, funcStartName string, args ...interface{}) error {
	if len(fileIn) == 0 {
		return fmt.Errorf("empty file name")
	}

	data, err := os.ReadFile(fileIn)
	if err != nil {
		return fmt.Errorf("exec load file: %w", err)
	}
	return e.ExecString(&entity.SourceItem{
		Name:       fileIn,
		SourceCode: string(data),
	}, nil, funcStartName, args)
}

func (e *Executor) ExecString(
	sourceItems *entity.SourceItem,
	callback func(string, int, [][]string, []*entity.Variable),
	funcStartName string,
	args ...interface{},
) error {
	initFuncName, _, err := e.ie.TranslateText(sourceItems.Name, sourceItems.SourceCode, e.debug, e.ie.OutputTranslate)
	if err != nil {
		return fmt.Errorf("translate error: %w", err)
	}

	ce, err := e.ie.CreateContextEnv()
	if err != nil {
		return fmt.Errorf("create context error %w", err)
	}

	// всегда вызываем функцию инициализации
	if len(initFuncName) > 0 {
		_, err = e.ie.InterpreterFunc(ce, initFuncName, []*frl.Value{})
		if err != nil {
			return fmt.Errorf("interepreter function %v error %w", initFuncName, err)
		}
		for {
			flag, err := e.ie.InterpreterFuncStep()
			if err != nil {
				return fmt.Errorf("interpreter %v function step %w", initFuncName, err)
			}
			if flag {
				break
			}
		}
	}

	if len(funcStartName) == 0 {
		return nil
	}

	values := []*frl.Value{}
	for _, arg := range args {
		values = append(values, frl.CreateValue(arg))
	}
	_, err = e.ie.InterpreterFunc(ce, funcStartName, values)
	if err != nil {
		return fmt.Errorf("intrepreter %v function error %w", funcStartName, err)
	}

	for {
		flag, err := e.interpreterStep(ce, callback)
		if err != nil {
			return fmt.Errorf("interpreter  function step %w", err)
		}
		if flag {
			break
		}
	}

	return e.setEvents(ce, callback)
}

func (e *Executor) ExecuteFuncWithManyFiles(
	sourceItems []entity.SourceItem,
	callback func(string, int, [][]string, []*entity.Variable),
	funcStartName string,
	args ...interface{},
) error {
	initFuncList := []string{}
	for _, sourceItem := range sourceItems {
		breakPoints := []*frl.BreakPoint{}
		for _, breakpoint := range sourceItem.Breakpoints {
			breakPoint := frl.BreakPoint{FileName: sourceItem.Name, LineNum: breakpoint}
			breakPoints = append(breakPoints, &breakPoint)
		}
		initFuncName, _, err := e.ie.TranslateText(sourceItem.Name, sourceItem.SourceCode, e.debug, e.ie.OutputTranslate)
		if err != nil {
			return fmt.Errorf("translate error: %w", err)
		}
		if len(breakPoints) > 0 {
			e.ie.AddBreakPoints(breakPoints)
		}
		initFuncList = append(initFuncList, initFuncName)
	}

	ce, err := e.ie.CreateContextEnv()
	if err != nil {
		return fmt.Errorf("create context error %w", err)
	}

	/*
		interpreterStep := func() (bool, error) {
			flag, err := e.ie.InterpreterFuncStep()
			if err != nil {
				return false, err
			}
			if flag {
				return true, nil
			}
			bp := e.ie.GetCurrentBreakPoint()
			if bp != nil {
				cf := ce.GetCurrentFunc()
				fn := cf.GetFunc()
				fnName := fn.Name
				data := [][]string{}
				variables := []*entity.Variable{}
				for k, v := range cf.GetVarDict() {
					data = append(data, []string{fnName, k, fmt.Sprintf("%v", v.GetType()), v.String()})
					variable := entity.Variable{
						FuncName: fnName,
						Name:     k,
						Type:     v.GetType().String(),
						Value:    v.String(),
					}
					variables = append(variables, &variable)
				}

				if callback != nil {
					callback(bp.FileName, bp.LineNum, data, variables)
				}
				e.ie.ClearCurrentBreakPoint()
			}
			return false, nil
		}
	*/

	for i := range initFuncList {
		if len(initFuncList[i]) > 0 {
			_, err = e.ie.InterpreterFunc(ce, initFuncList[i], []*frl.Value{})
			if err != nil {
				return fmt.Errorf("intrepreter function error %w", err)
			}
			for {
				flag, err := e.interpreterStep(ce, callback)
				if err != nil {
					return fmt.Errorf("interpreter  function step %w", err)
				}
				if flag {
					break
				}
				/*
					flag, err := e.ie.InterpreterFuncStep()
					if err != nil {
						return fmt.Errorf("interpreter  function step %w", err)
					}
					if flag {
						break
					}
					bp := e.ie.GetCurrentBreakPoint()
					if bp != nil {
						cf := ce.GetCurrentFunc()
						fn := cf.GetFunc()
						fnName := fn.Name
						data := [][]string{}
						variables := []*entity.Variable{}
						for k, v := range cf.GetVarDict() {
							data = append(data, []string{fnName, k, fmt.Sprintf("%v", v.GetType()), v.String()})
							variable := entity.Variable{
								FuncName: fnName,
								Name:     k,
								Type:     v.GetType().String(),
								Value:    v.String(),
							}
							variables = append(variables, &variable)
						}

						if callback != nil {
							callback(bp.FileName, bp.LineNum, data, variables)
						}
						e.ie.ClearCurrentBreakPoint()
					}
				*/
			}
		}
	}
	if len(funcStartName) == 0 {
		return nil
	}
	values := []*frl.Value{}
	for _, arg := range args {
		values = append(values, frl.CreateValue(arg))
	}
	_, err = e.ie.InterpreterFunc(ce, funcStartName, values)
	if err != nil {
		return fmt.Errorf("intrepreter function error %w", err)
	}

	for {
		flag, err := e.interpreterStep(ce, callback)
		if err != nil {
			return fmt.Errorf("interpreter  function step %w", err)
		}
		if flag {
			break
		}
	}

	return e.setEvents(ce, callback)
	/*
		if len(e.ie.Events) == 0 {
			return nil
		}

		type eventTimer struct {
			id string
			tm *frl.Value
		}

		type eventChannel struct {
			id   string
			data *frl.Value
		}

		e.ie.SetDone(make(chan struct{}))
		events := e.ie.Events // e.eventManager.GetEvents()

		et := make(chan *eventTimer, 3)
		ec := make(chan *eventChannel, 3)

		type eventData struct {
			id    string
			event *frl.Event
		}
		eventByID := make(map[string]*eventData)

		for i := range events {
			switch events[i].Type {
			case "duration":
				tt := time.NewTicker(events[i].Duration)
				id := uuid.NewString()
				eventByID[id] = &eventData{
					id:    id,
					event: events[i],
				}
				go func() {
					for {
						select {
						case tm := <-tt.C:
							tmValue := frl.NewValue(int(frl.VtString), tm.Format(time.RFC3339))
							et <- &eventTimer{
								id: id,
								tm: tmValue,
							}
						case <-e.ie.GetDone():
							return
						}
					}
				}()

			case "channel":
				id := uuid.NewString()
				eventByID[id] = &eventData{
					id:    id,
					event: events[i],
				}
				cs, ok := e.ie.Channels[events[i].Channel]
				if !ok {
					continue
				}
				cv := make(chan *frl.Value)
				cs.Value = cv
				go func() {
					for {
						select {
						case data := <-cv:
							ec <- &eventChannel{
								id:   id,
								data: data,
							}
						case <-e.ie.GetDone():
							return
						}
					}
				}()
			}
		}

		for {
			select {
			case ee := <-et:
				event, ok := e.ie.EventsByID[ee.id]
				if ok {
					if len(funcStartName) == 0 {
						return nil
					}
					values := []*frl.Value{ee.tm}
					_, err = e.ie.InterpreterFunc(ce, event.Fn, values)
					if err != nil {
						return fmt.Errorf("interpreter function error %w", err)
					}
					for {
						flag, err := e.interpreterStep(ce, callback)
						if err != nil {
							return fmt.Errorf("interpreter  function step %w", err)
						}
						if flag {
							break
						}
					}
				}
			case ee := <-ec:
				event, ok := e.ie.EventsByID[ee.id]
				if ok {
					if len(funcStartName) == 0 {
						return nil
					}
					values := []*frl.Value{ee.data}
					_, err = e.ie.InterpreterFunc(ce, event.Fn, values)
					if err != nil {
						return fmt.Errorf("intrepreter function error %w", err)
					}
					for {
						flag, err := e.interpreterStep(ce, callback)
						if err != nil {
							return fmt.Errorf("interpreter  function step %w", err)
						}
						if flag {
							break
						}
					}
				}
			case <-e.ie.GetDone():
				return nil
			}
		}
	*/
}

func (e *Executor) interpreterStep(
	ce *frl.ContextEnv,
	callback func(string, int, [][]string, []*entity.Variable),
) (bool, error) {
	flag, err := e.ie.InterpreterFuncStep()
	if err != nil {
		return false, err
	}
	if flag {
		return true, nil
	}
	bp := e.ie.GetCurrentBreakPoint()
	if bp != nil {
		cf := ce.GetCurrentFunc()
		fn := cf.GetFunc()
		fnName := fn.Name
		data := [][]string{}
		variables := []*entity.Variable{}
		for k, v := range cf.GetVarDict() {
			data = append(data, []string{fnName, k, fmt.Sprintf("%v", v.GetType()), v.String()})
			variable := entity.Variable{
				FuncName: fnName,
				Name:     k,
				Type:     v.GetType().String(),
				Value:    v.String(),
			}
			variables = append(variables, &variable)
		}

		if callback != nil {
			callback(bp.FileName, bp.LineNum, data, variables)
		}
		e.ie.ClearCurrentBreakPoint()
	}
	return false, nil
}

func (e *Executor) setEvents(
	ce *frl.ContextEnv,
	callback func(string, int, [][]string, []*entity.Variable),
) error {
	if len(e.ie.Events) == 0 {
		return nil
	}

	type eventTimer struct {
		id string
		tm *frl.Value
	}

	type eventChannel struct {
		id   string
		data *frl.Value
	}

	e.ie.SetDone(make(chan struct{}))
	events := e.ie.Events

	et := make(chan *eventTimer, 3)
	ec := make(chan *eventChannel, 3)

	done := make(chan struct{}, 2)
	isDone := false
	for i := range events {
		switch events[i].Type {
		case "duration":
			tt := time.NewTicker(events[i].Duration)
			go func(id string) {
				for {
					select {
					case tm := <-tt.C:
						tmValue := frl.NewValue(frl.VtString, tm.Format(time.RFC3339))
						et <- &eventTimer{
							id: id,
							tm: tmValue,
						}
					case <-e.ie.GetDone():
						done <- struct{}{}
						isDone = true
						return
					}
				}
			}(events[i].ID)

		case "channel":
			cs, ok := e.ie.Channels[events[i].Channel]
			if !ok {
				continue
			}
			cv := make(chan *frl.Value)
			cs.Value = cv
			go func(id string) {
				for {
					select {
					case data := <-cv:
						ec <- &eventChannel{
							id:   id,
							data: data,
						}
					case <-e.ie.GetDone():
						done <- struct{}{}
						isDone = true
						return
					}
				}
			}(events[i].ID)
		}
	}

	for {
		select {
		case ee := <-et:
			event, ok := e.ie.EventsByID[ee.id]
			if ok {
				if len(event.Fn) == 0 {
					return nil
				}
				values := []*frl.Value{ee.tm}
				_, err := e.ie.InterpreterFunc(ce, event.Fn, values)
				if err != nil {
					return fmt.Errorf("intrepreter function error %w", err)
				}
				for {
					flag, err := e.interpreterStep(ce, callback)
					if err != nil {
						return fmt.Errorf("interpreter  function step %w", err)
					}
					if flag {
						break
					}
					if isDone {
						return nil
					}
				}
			}
		case ee := <-ec:
			event, ok := e.ie.EventsByID[ee.id]
			if ok {
				if len(event.Fn) == 0 {
					return nil
				}
				values := []*frl.Value{ee.data}
				_, err := e.ie.InterpreterFunc(ce, event.Fn, values)
				if err != nil {
					return fmt.Errorf("intrepreter function error %w", err)
				}
				for {
					flag, err := e.interpreterStep(ce, callback)
					if err != nil {
						return fmt.Errorf("interpreter  function step %w", err)
					}
					if flag {
						break
					}
					if isDone {
						return nil
					}
				}
			}
			//		case <-done:
			//			return nil
		}
		if isDone {
			return nil
		}

	}
}

func (e *Executor) TranslateManyFiles(
	sourceItems []entity.SourceItem,
) error {
	for _, sourceItem := range sourceItems {
		_, _, err := e.ie.TranslateText(sourceItem.Name, sourceItem.SourceCode, e.debug, e.ie.OutputTranslate)
		if err != nil {
			return fmt.Errorf("translate error: %w", err)
		}
	}
	return nil
}
