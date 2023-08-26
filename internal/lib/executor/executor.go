package executor

import (
	"fmt"
	"os"

	frl "github.com/wanderer69/FrL/internal/lib"
	print "github.com/wanderer69/tools/parser/print"
)

type ExecutorBase struct {
	fe     *frl.FrameEnvironment
	ie     *frl.InterpreterEnv
	output *print.Output
}

type Executor struct {
	eb    *ExecutorBase
	debug int
}

func InitExecutorBase(debug int, output *print.Output) *ExecutorBase {
	eb := &ExecutorBase{}
	// настраиваем окружение
	fe := frl.NewFrameEnvironment()
	fe.FrameDict = make(map[string][]*frl.Frame)

	ie := frl.NewInterpreterEnv()
	ie.SetDebug(debug) //xfd xff xff
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

	ie.SetFrameEnvironment(fe)

	eb.fe = fe
	eb.ie = ie
	eb.output = output
	eb.ie.Output = output
	eb.ie.FE = fe
	return eb
}

func InitExecutor(eb *ExecutorBase, debug int) *Executor {
	return &Executor{
		eb:    eb,
		debug: debug,
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
	return e.ExecString(fileIn, string(data), funcStartName, args)
}

func (e *Executor) ExecString(fileName string, data string, funcStartName string, args ...interface{}) error {
	_, err := e.eb.ie.TranslateText(fileName, data, e.debug, e.eb.ie.Output)
	if err != nil {
		return fmt.Errorf("translate error: %w", err)
	}

	ce, err := e.eb.ie.CreateContextEnv()
	if err != nil {
		return fmt.Errorf("create context error %w", err)
	}

	if len(funcStartName) == 0 {
		return nil
	}
	values := []*frl.Value{}
	for _, arg := range args {
		values = append(values, frl.CreateValue(arg))
	}
	_, err = e.eb.ie.InterpreterFunc(ce, funcStartName, values)
	if err != nil {
		return fmt.Errorf("intrepreter function error %w", err)
	}
	for {
		flag, err := e.eb.ie.InterpreterFuncStep()
		if err != nil {
			return fmt.Errorf("interpreter  function step %w", err)
		}
		if flag {
			break
		}
	}
	return nil
}

type SourceItem struct {
	Name        string
	SourceCode  string
	Breakpoints []int
}

func (e *Executor) ExecuteFuncWithManyFiles(
	sourceItems []SourceItem,
	callback func(string, int, [][]string),
	funcStartName string,
	args ...interface{},
) error {
	for _, sourceItem := range sourceItems {
		breakPoints := []*frl.BreakPoint{}
		for _, breakpoint := range sourceItem.Breakpoints {
			breakPoint := frl.BreakPoint{FileName: sourceItem.Name, LineNum: breakpoint}
			breakPoints = append(breakPoints, &breakPoint)
		}
		_, err := e.eb.ie.TranslateText(sourceItem.Name, sourceItem.SourceCode, e.debug, e.eb.ie.Output)
		if err != nil {
			return fmt.Errorf("translate error: %w", err)
		}
		if len(breakPoints) > 0 {
			e.eb.ie.AddBreakPoints(breakPoints)
		}
	}

	ce, err := e.eb.ie.CreateContextEnv()
	if err != nil {
		return fmt.Errorf("create context error %w", err)
	}

	if len(funcStartName) == 0 {
		return nil
	}
	values := []*frl.Value{}
	for _, arg := range args {
		values = append(values, frl.CreateValue(arg))
	}
	_, err = e.eb.ie.InterpreterFunc(ce, funcStartName, values)
	if err != nil {
		return fmt.Errorf("intrepreter function error %w", err)
	}
	for {
		flag, err := e.eb.ie.InterpreterFuncStep()
		if err != nil {
			return fmt.Errorf("interpreter  function step %w", err)
		}
		if flag {
			break
		}
		bp := e.eb.ie.GetCurrentBreakPoint()
		if bp != nil {
			cf := ce.GetCurrentFunc()
			data := [][]string{}
			for k, v := range cf.GetVarDict() {
				data = append(data, []string{cf.GetFunc().Name, k, fmt.Sprintf("%v", v.GetType()), v.String()})
			}

			if callback != nil {
				callback(bp.FileName, bp.LineNum, data)
			}
			e.eb.ie.ClearCurrentBreakPoint()
		}
	}
	return nil
}

func (e *Executor) TranslateManyFiles(
	sourceItems []SourceItem,
) error {
	for _, sourceItem := range sourceItems {
		_, err := e.eb.ie.TranslateText(sourceItem.Name, sourceItem.SourceCode, e.debug, e.eb.ie.Output)
		if err != nil {
			return fmt.Errorf("translate error: %w", err)
		}
	}
	return nil
}
