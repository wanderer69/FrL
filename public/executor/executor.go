package executor

import (
	"fmt"
	"os"

	frl "github.com/wanderer69/FrL/public/lib"
	print "github.com/wanderer69/tools/parser/print"
)

type ExecutorBase struct {
	fe     *frl.FrameEnvironment
	output *print.Output
}

type Executor struct {
	eb     *ExecutorBase
	ie     *frl.InterpreterEnv
	debug  int
	output *print.Output
}

func InitExecutorBase(debug int, output *print.Output) *ExecutorBase {
	return &ExecutorBase{
		// настраиваем окружение
		fe:     frl.NewFrameEnvironment(),
		output: output,
	}
}

func InitExecutor(eb *ExecutorBase, output *print.Output, debug int) *Executor {
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

	ie.SetFrameEnvironment(eb.fe)
	ie.FE = eb.fe

	ie.Output = output

	return &Executor{
		eb:    eb,
		debug: debug,
		ie:    ie,
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
	_, err := e.ie.TranslateText(fileName, data, e.debug, e.ie.Output)
	if err != nil {
		return fmt.Errorf("translate error: %w", err)
	}

	ce, err := e.ie.CreateContextEnv()
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
	_, err = e.ie.InterpreterFunc(ce, funcStartName, values)
	if err != nil {
		return fmt.Errorf("intrepreter function error %w", err)
	}
	for {
		flag, err := e.ie.InterpreterFuncStep()
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

type Variable struct {
	FuncName string
	Name     string
	Type     string
	Value    string
}

func (e *Executor) ExecuteFuncWithManyFiles(
	sourceItems []SourceItem,
	callback func(string, int, [][]string, []*Variable),
	funcStartName string,
	args ...interface{},
) error {
	for _, sourceItem := range sourceItems {
		breakPoints := []*frl.BreakPoint{}
		for _, breakpoint := range sourceItem.Breakpoints {
			breakPoint := frl.BreakPoint{FileName: sourceItem.Name, LineNum: breakpoint}
			breakPoints = append(breakPoints, &breakPoint)
		}
		_, err := e.ie.TranslateText(sourceItem.Name, sourceItem.SourceCode, e.debug, e.ie.Output)
		if err != nil {
			return fmt.Errorf("translate error: %w", err)
		}
		if len(breakPoints) > 0 {
			e.ie.AddBreakPoints(breakPoints)
		}
	}

	ce, err := e.ie.CreateContextEnv()
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
	_, err = e.ie.InterpreterFunc(ce, funcStartName, values)
	if err != nil {
		return fmt.Errorf("intrepreter function error %w", err)
	}
	for {
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
			variables := []*Variable{}
			for k, v := range cf.GetVarDict() {
				data = append(data, []string{fnName, k, fmt.Sprintf("%v", v.GetType()), v.String()})
				variable := Variable{
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
	}
	return nil
}

func (e *Executor) TranslateManyFiles(
	sourceItems []SourceItem,
) error {
	for _, sourceItem := range sourceItems {
		_, err := e.ie.TranslateText(sourceItem.Name, sourceItem.SourceCode, e.debug, e.ie.Output)
		if err != nil {
			return fmt.Errorf("translate error: %w", err)
		}
	}
	return nil
}
