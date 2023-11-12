package main

import (
	"flag"
	"fmt"

	"os"

	fnc "github.com/wanderer69/FrL/public/functions"
	frl "github.com/wanderer69/FrL/public/lib"
	print "github.com/wanderer69/tools/parser/print"
)

func main() {
	var file_in string
	flag.StringVar(&file_in, "file_in", "", "input frm file")

	var debug_file string
	flag.StringVar(&debug_file, "debug_file", "", "file debug configuration")

	//	var mode_list string
	//	flag.StringVar(&mode_list, "mode_list", "", "list of modes hash,frame_simple,test_to_type,relations_store,relations_load,script_load,types_test")

	flag.Parse()

	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}

	output := print.NewOutput(printFunc)

	if len(file_in) == 0 {
		flag.PrintDefaults()
	}

	// настраиваем окружение
	fe := frl.NewFrameEnvironment()
	fe.FrameDict = make(map[string][]*frl.Frame)

	if false {
		ns, err := frl.NewStore("./Frames", output)
		if err != nil {
			fmt.Printf("err %v\r\n", err)
			return
		}

		loader, err1 := ns.LoadFrameRecord(0)
		if err1 != nil {
			fmt.Printf("err %v\r\n", err1)
			return
		}
		fm := make(map[string]*frl.Frame)
		for {
			frame_id, slot_name, slot_property, slot_value, err := loader()
			if err != nil {
				break
			}
			//fmt.Printf("fi %v\r\n", fi )
			// find frame by id
			var f *frl.Frame

			fl, err := fe.QueryRelations(frl.QueryRelationItem{ObjectType: "frame", Object: "", Value: frl.CreateValue(frame_id)})
			if err != nil {
				fmt.Printf("err %v\r\n", err)
			}
			// fmt.Printf("fl %v\r\n", fl)
			if len(fl) == 0 {
				f = frl.NewFrame()
				// добавляем поле уникального идентификатора
				f.AddSlot("ID")
				v, _ := f.Set("ID", frame_id)
				fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
				fe.Frames = append(fe.Frames, f)
				fm[frame_id] = f
			} else {
				f = fl[0]
			}
			f.AddSlot(slot_name)
			f.SetSlotProperty(slot_name, slot_property)
			f.SetValue(slot_name, slot_value)
			fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: slot_name, Value: slot_value})
		}
	}
	ie := frl.NewInterpreterEnv()
	ie.SetDebug(0) //xfd xff xff
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

	if false {
		bb, err := os.ReadFile(file_in)
		if err != nil {
			fmt.Print(err)
			return
		}

		for {
			fp, bb_, err1 := fnc.Bin2Func(bb)
			if err1 != nil {
				fmt.Print(err1)
				return
			}
			fmt.Printf("len bb_ %v\r\n", len(bb_))
			s := fnc.PrintFunction(fp)
			fmt.Printf("%v\r\n", s)
			if len(bb_) > 0 {
				bb = bb_
			} else {
				break
			}
		}

		data, err := os.ReadFile(file_in)
		if err != nil {
			fmt.Print(err)
			return
		}

		_, err = ie.TranslateText(file_in, string(data), 0, ie.Output)
		if err != nil {
			fmt.Print(err)
			return
		}

		ce, err := ie.CreateContextEnv()
		if err != nil {
			fmt.Printf("create context error %v", err)
			return
		}

		values := []*frl.Value{frl.CreateValue("1"), frl.CreateValue("2")}
		_, err1 := ie.InterpreterFunc(ce, "пример1", values)
		if err1 != nil {
			fmt.Print(err1)
			return
		}
		for {
			flag, err1 := ie.InterpreterFuncStep( /*cf*/ )
			if err1 != nil {
				fmt.Print(err1)
				return
			}
			if flag {
				break
			}
		}
	}
}
