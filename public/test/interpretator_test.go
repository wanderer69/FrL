package frl_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/debug"

	"github.com/wanderer69/FrL/internal/addons/convertor"
	exec "github.com/wanderer69/FrL/public/executor"
	fnc "github.com/wanderer69/FrL/public/functions"
	frl "github.com/wanderer69/FrL/public/lib"

	//	uqe "github.com/wanderer69/FrL/src/lib/unique"
	uqe "github.com/wanderer69/tools/unique"

	print "github.com/wanderer69/tools/parser/print"
)

func TestFrame(t *testing.T) {
	debug.NewDebug()

	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}

	output := print.NewOutput(printFunc)

	t.Run("frame_simple", func(t *testing.T) {
		buffer := ""

		compareResult := func(state int) {
			samples := ""
			switch state {
			case 1:
				samples = "ID () 1\r\n"
			case 2:
				samples = "ID () 1, отношение () value2\r\n"
			case 3:
				samples = "ID () 1, отношение () value2, наименование () value1\r\n"
			case 4:
				samples = "ID () 1, отношение () value2, наименование () value1, имя () value1\r\n"
			case 5:
				samples = "ID () 1, отношение () value2, наименование () value1, имя () value1\r\n"
			}
			fmt.Printf("%v", buffer)
			require.Equal(t, samples, buffer)
			buffer = ""
		}
		printFuncSelect := func(frm string, args ...any) {
			buffer = buffer + fmt.Sprintf(frm, args...)
		}

		outputSelect := print.NewOutput(printFuncSelect)

		f := frl.NewFrame()
		err := f.AddSlot("ID")
		require.NoError(t, err)
		_, err = f.Set("ID", 1)
		require.NoError(t, err)

		fmt.Printf("--- ID\r\n")
		f.Print(outputSelect, true)
		compareResult(1)
		fmt.Printf("-\r\n")

		err = f.AddSlot("отношение")
		require.NoError(t, err)
		_, err = f.Set("отношение", "value2")
		require.NoError(t, err)
		fmt.Printf("--- отношение\r\n")
		f.Print(outputSelect, true)
		compareResult(2)
		fmt.Printf("-\r\n")

		err = f.AddSlot("наименование")
		require.NoError(t, err)
		_, err = f.Set("наименование", "value1")
		require.NoError(t, err)
		fmt.Printf("--- наименование\r\n")
		f.Print(outputSelect, true)
		compareResult(3)
		fmt.Printf("-_\r\n")

		err = f.AddSlot("имя")
		require.NoError(t, err)
		_, err = f.Set("имя", "value1")
		require.NoError(t, err)
		fmt.Printf("--- имя\r\n")
		f.Print(outputSelect, true)
		compareResult(4)
		fmt.Printf("-\r\n")

		f1 := frl.NewFrame()
		err = f1.AddSlot("slot4")
		require.NoError(t, err)
		_, err = f1.Set("slot4", []int{1, 2, 3})
		require.NoError(t, err)
		f.Set("slot2", f1)
		require.NoError(t, err)
		fmt.Printf("---\r\n")
		//state = 5
		f.Print(outputSelect, true)
		compareResult(5)
		fmt.Printf("-\r\n")

		samples := [][]string{{"ID", "1"}, {"отношение", "value2"}, {"наименование", "value1"}, {"имя", "value1"}}
		pos := 0
		ff := f.Iterate()
		for {
			s, ok, err := ff()
			if err != nil {
				break
			}
			ssl := s.GetSlotValue()
			for i := range ssl {
				ss, _ := frl.FromType(ssl[i])
				fmt.Printf("k %v v %v\r\n", s.GetSlotName(), ss)
				require.Equal(t, samples[pos][0], s.GetSlotName())
				require.Equal(t, samples[pos][1], ss)
				pos++
			}
			if ok {
				break
			}
		}
		require.True(t, false)
	})

	t.Run("value_store", func(t *testing.T) {
		buffer := ""
		mask := "ID () %v, отношение () %v, наименование () %v, тип_отношения () %v\r\n"
		var esimatedDataPtr *string
		compareResult := func(state int) {
			samples := ""
			switch state {
			case 1:
				samples = *esimatedDataPtr
				//"ID () PQ6U85RO24, отношение () 60RR0I51R9, наименование () отношение_60RR0I51R9, тип_отношения () CHA32\r\n"
			case 2:
				samples = "1, 2.5, 3, " + *esimatedDataPtr
			}
			fmt.Printf("%v", buffer)
			require.Equal(t, samples, buffer)
			buffer = ""
		}
		printFuncSelect := func(frm string, args ...any) {
			buffer = buffer + fmt.Sprintf(frm, args...)
		}

		outputSelect := print.NewOutput(printFuncSelect)

		// настраиваем окружение
		fe := frl.NewFrameEnvironment()

		f := frl.NewFrame()
		// добавляем поле уникального идентификатора
		f.AddSlot("ID")
		id := uqe.UniqueValue(10)
		relation := uqe.UniqueValue(10)
		relationName := "отношение_" + relation
		relationType := uqe.UniqueValue(5)

		esimatedData := fmt.Sprintf(mask, id, relation, relationName, relationType)
		esimatedDataPtr = &esimatedData

		v, err := f.Set("ID", id)
		require.NoError(t, err)
		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})

		err = f.AddSlot("отношение")
		require.NoError(t, err)
		v, err = f.Set("отношение", relation)
		require.NoError(t, err)
		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "отношение", Value: v})

		err = f.AddSlot("наименование")
		require.NoError(t, err)
		v, err = f.Set("наименование", relationName)
		require.NoError(t, err)
		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "наименование", Value: v})

		err = f.AddSlot("тип_отношения")
		require.NoError(t, err)
		_, err = f.Set("тип_отношения", relationType)
		require.NoError(t, err)

		f_v := frl.CreateValue(f)
		err = f_v.Print(outputSelect)
		compareResult(1)
		require.NoError(t, err)

		bb, err := frl.SaveValueStore(f_v)
		require.NoError(t, err)
		vv, bbLast, err := frl.LoadValueStore(bb)
		require.NoError(t, err)
		require.Len(t, bbLast, 0)

		err = vv.Print(outputSelect)
		require.NoError(t, err)
		compareResult(1)

		slice_v := frl.CreateValue([]*frl.Value{frl.CreateValue(1), frl.CreateValue(2.5), frl.CreateValue("3"), frl.CreateValue(f)})
		err = slice_v.Print(outputSelect)
		require.NoError(t, err)
		compareResult(2)

		bb, err = frl.SaveValueStore(slice_v)
		require.NoError(t, err)

		vv, bbLast, err = frl.LoadValueStore(bb)
		require.NoError(t, err)
		require.Len(t, bbLast, 0)

		err = vv.Print(outputSelect)
		require.NoError(t, err)
		compareResult(2)

		require.True(t, false)
	})

	t.Run("relations_store", func(t *testing.T) {
		os.Remove("./Frames")
		ns, err := frl.NewStore("./Frames", output)
		require.NoError(t, err)

		// настраиваем окружение
		fe := frl.NewFrameEnvironment()
		//fe.FrameDict = make(map[string][]*frl.Frame)

		// добавляем фрейм с отношением "имя"
		// надо добавить фрейм с определением отношения
		f := frl.NewFrame()
		// добавляем поле уникального идентификатора
		err = f.AddSlot("ID")
		require.NoError(t, err)

		id := uqe.UniqueValue(10)
		v, err := f.Set("ID", id)
		require.NoError(t, err)

		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
		relation := "наименование"
		err = f.AddSlot("отношение")
		require.NoError(t, err)

		v, err = f.Set("отношение", relation)
		require.NoError(t, err)

		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "отношение", Value: v})

		r, err := convertor.LoadRelation("/home/user/Go_projects/SemanticNet/data/relation.txt")
		require.NoError(t, err)

		if false {
			for i := range r {
				fmt.Printf("%v\r\n", r[i])
				for j := range r[i].RelationItem {
					fmt.Printf("\t%v\r\n", r[i].RelationItem[j])
				}
			}
		}

		for i := range r {
			for j := range r[i].RelationItem {
				// надо добавить фрейм с определением отношения
				f := frl.NewFrame()
				// добавляем поле уникального идентификатора
				err = f.AddSlot("ID")
				require.NoError(t, err)
				id := uqe.UniqueValue(10)
				v, err := f.Set("ID", id)
				require.NoError(t, err)
				fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
				relation := r[i].RelationItem[j].Relation

				err = f.AddSlot("отношение")
				require.NoError(t, err)
				v, err = f.Set("отношение", relation)
				require.NoError(t, err)
				fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "отношение", Value: v})

				err = f.AddSlot("наименование")
				require.NoError(t, err)
				v, err = f.Set("наименование", "отношение_"+relation)
				require.NoError(t, err)
				fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "наименование", Value: v})

				err = f.AddSlot("тип_отношения")
				require.NoError(t, err)
				_, err = f.Set("тип_отношения", r[i].RelationType)
				require.NoError(t, err)

				if len(r[i].RelationItem[j].RelationExample) > 0 {
					isSlotExist := false
					for k := range r[i].RelationItem[j].RelationExample {
						if !isSlotExist {
							err = f.AddSlot("пример")
							require.NoError(t, err)
							isSlotExist = true
						}
						_, err = f.Set("пример", r[i].RelationItem[j].RelationExample[k])
						require.NoError(t, err)
					}
				}

				/*
					for k, _ := range r[i].relationItem[j].relation_example {
						f.Set("пример", r[i].relationItem[j].relation_example[k])
					}
				*/
				fe.Frames = append(fe.Frames, f)
				//f.Print(true)
			}
		}

		if false {
			for i := range fe.Frames {
				f := fe.Frames[i]
				f.Print(output, true)
			}
		}

		for i := range fe.Frames {
			f := fe.Frames[i]
			ff := f.Iterate()
			frame_ids, err := f.GetValue("ID")
			require.NoError(t, err)

			frame_id := frame_ids[0]
			for {
				s, ok, err := ff()
				if err != nil {
					break
				}
				ssl := s.GetSlotValue()
				slot_name := s.GetSlotName()
				slot_property := s.GetSlotProperty()
				if slot_name != "ID" {
					for j := range ssl {
						err := ns.SaveFrameRecord(frame_id, slot_name, slot_property, ssl[j], 0)
						require.NoError(t, err)
					}
				}
				if ok {
					break
				}
			}
		}

		require.True(t, false)
	})

	t.Run("test_to_type", func(t *testing.T) {
		v, err := frl.ToType(true)
		require.True(t, err)
		require.Equal(t, 1, v)

		//		fmt.Printf("v %v, err %v\r\n", v, err)
		v, err = frl.ToType(1)
		require.True(t, err)
		require.Equal(t, 2, v)

		//		fmt.Printf("v %v, err %v\r\n", v, err)
		v, err = frl.ToType(1.1)
		require.True(t, err)
		require.Equal(t, 3, v)

		//		fmt.Printf("v %v, err %v\r\n", v, err)
		v, err = frl.ToType("qwert")
		require.True(t, err)
		require.Equal(t, 4, v)

		//		fmt.Printf("v %v, err %v\r\n", v, err)
		v, err = frl.ToType(frl.NewFrame())
		require.True(t, err)
		require.Equal(t, 5, v)

		//		fmt.Printf("v %v, err %v\r\n", v, err)

		v, err = frl.ToType([]*frl.Value{frl.CreateValue("1"), frl.CreateValue("2"), frl.CreateValue("3")})
		require.True(t, err)
		require.Equal(t, 6, v)

		//		fmt.Printf("v %v, err %v\r\n", v, err)

		vv := frl.CreateValue([]*frl.Value{frl.CreateValue("1"), frl.CreateValue("2"), frl.CreateValue("3")})
		v1, err1 := frl.NewIterator(vv)
		require.NoError(t, err1)
		//		fmt.Printf("vl %v, err %v\r\n", v1, err1)

		v, err = frl.ToType(v1)
		require.True(t, err)
		require.Equal(t, 7, v)

		v, err = frl.ToType(nil)
		//		fmt.Printf("v %v, err %v\r\n", v, err)
		require.True(t, err)
		require.Equal(t, 0, v)

		require.True(t, false)
	})

	t.Run("uri_parse", func(t *testing.T) {
		uri := "file://"
		dct, err := frl.ParseURI(uri)
		require.NoError(t, err)
		fmt.Printf("dct %v\r\n", dct)

		dctExpected := map[string]string{"fragment": "", "path": "", "query": "", "schema": "file", "source": ""}
		require.Equal(t, dctExpected, dct)

		require.True(t, false)
	})

	t.Run("types_test", func(t *testing.T) {
		buffer := ""
		mask := "ID () %v, наименование () фрейм1\r\n"
		mask3 := "slot ID property  , %v\r\n"
		var esimatedDataPtr *string
		prev := ""
		compareResult := func(state int) {
			samples := ""
			switch state {
			case 1:
				samples = *esimatedDataPtr
				prev = samples
			case 2:
				samples = "iterator type 5 pos 0\r\n"
			case 3:
				samples = *esimatedDataPtr
			case 4:
				samples = "slot наименование property  , фрейм1\r\n"
			case 5:
				samples = "1, 2, 3, " + prev
			case 6:
				samples = "iterator type 6 pos 0\r\n"
			case 7:
				samples = "1\r\n"
			case 8:
				samples = "2\r\n"
			case 9:
				samples = "3\r\n"
			case 10:
				samples = prev
			case 11:
				samples = "iterator type 4 pos 0\r\n"
			case 12:
				samples = "1\r\n"
			case 13:
				samples = "2\r\n"
			case 14:
				samples = "3\r\n"
			case 15:
				samples = "б\r\n"
			case 16:
				samples = "а\r\n"
			case 17:
				samples = "р\r\n"
			case 18:
				samples = "б\r\n"
			case 19:
				samples = "а\r\n"
			case 20:
				samples = "р\r\n"
			}
			fmt.Printf("%v", buffer)
			require.Equal(t, samples, buffer)
			buffer = ""
		}
		printFuncSelect := func(frm string, args ...any) {
			buffer = buffer + fmt.Sprintf(frm, args...)
		}

		outputSelect := print.NewOutput(printFuncSelect)

		id := uqe.UniqueValue(7)
		estimatedData := fmt.Sprintf(mask, id)
		esimatedDataPtr = &estimatedData

		fe := frl.NewFrameEnvironment()
		//fe.FrameDict = make(map[string][]*frl.Frame)

		f := frl.NewFrame()
		// добавляем поле уникального идентификатора
		err := f.AddSlot("ID")
		require.NoError(t, err)

		v, err := f.Set("ID", id)
		require.NoError(t, err)

		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
		fe.Frames = append(fe.Frames, f)
		err = f.AddSlot("наименование")
		require.NoError(t, err)
		v, err = f.Set("наименование", "фрейм1")
		require.NoError(t, err)
		fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "наименование", Value: v})

		frame_v := frl.CreateValue(f)
		err = frame_v.Print(outputSelect)
		require.NoError(t, err)
		compareResult(1)

		iter_v, err := frl.NewIterator(frame_v)
		require.NoError(t, err)

		err = iter_v.Print(outputSelect)
		require.NoError(t, err)
		compareResult(2)

		estimatedData = fmt.Sprintf(mask3, id)
		esimatedDataPtr = &estimatedData

		pos := 3
		for {
			v, err := iter_v.Iterate()
			if err != nil {
				break
			}
			if v != nil {
				err = v.Print(outputSelect)
				require.NoError(t, err)
				compareResult(pos)
				pos++
			} else {
				break
			}
		}

		slice_v := frl.CreateValue([]*frl.Value{frl.CreateValue(1), frl.CreateValue("2"), frl.CreateValue("3"), frame_v})
		err = slice_v.Print(outputSelect)
		require.NoError(t, err)
		compareResult(5)

		iter_v1, err1 := frl.NewIterator(slice_v)
		require.NoError(t, err1)

		err = iter_v1.Print(outputSelect)
		require.NoError(t, err)
		compareResult(6)

		pos = 7
		for {
			v, err := iter_v1.Iterate()
			if err != nil {
				break
			}
			if v != nil {
				err = v.Print(outputSelect)
				require.NoError(t, err)
				compareResult(pos)
				pos++
			} else {
				break
			}
		}

		str_v := frl.CreateValue("123барбар")
		iter_v2, err2 := frl.NewIterator(str_v)
		require.NoError(t, err2)

		err = iter_v2.Print(outputSelect)
		require.NoError(t, err)
		compareResult(11)

		pos = 12
		for {
			v, err := iter_v2.Iterate()
			if err != nil {
				break
			}
			if v != nil {
				err = v.Print(outputSelect)
				require.NoError(t, err)
				compareResult(pos)
				pos++
			} else {
				break
			}
		}

		require.True(t, false)
	})

	t.Run("relations_load", func(t *testing.T) {
		// настраиваем окружение
		fe := frl.NewFrameEnvironment()
		//fe.FrameDict = make(map[string][]*frl.Frame)

		ns, err := frl.NewStore("./Frames", output)
		require.NoError(t, err)

		loader, err1 := ns.LoadFrameRecord(0)
		require.NoError(t, err1)

		fm := make(map[string]*frl.Frame)
		for {
			frame_id, slot_name, slot_property, slot_value, err := loader()
			if err != nil {
				break
			}
			//fmt.Printf("fi %v\r\n", fi )
			// find frame by id
			var f *frl.Frame
			/*
				f, ok := fm[frame_id]
				require.False(t, ok)
				require.Nil(t, f)
			*/
			fl, err := fe.QueryRelations(frl.QueryRelationItem{ObjectType: "frame", Object: "", Value: frl.CreateValue(frame_id)})
			require.NoError(t, err)

			// fmt.Printf("fl %v\r\n", fl)
			if len(fl) == 0 {
				f = frl.NewFrame()
				// добавляем поле уникального идентификатора
				err = f.AddSlot("ID")
				require.NoError(t, err)
				v, err := f.Set("ID", frame_id)
				require.NoError(t, err)
				fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
				fe.Frames = append(fe.Frames, f)
				fm[frame_id] = f
			} else {
				f = fl[0]
			}
			err = f.AddSlot(slot_name)
			require.NoError(t, err)
			err = f.SetSlotProperty(slot_name, slot_property)
			require.NoError(t, err)
			_, err = f.SetValue(slot_name, slot_value)
			require.NoError(t, err)
			fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: slot_name, Value: slot_value})
		}
		if false {
			for i := range fe.Frames {
				f := fe.Frames[i]
				f.Print(output, true)
			}
		}
		if false {
			fd := fe.GetFrameDict()
			for k, v := range fd {
				fmt.Printf("key %v len v %v\r\n", k, len(v))
				/*
					for i, _ := range v {
						f := v[i]
						f.Print(true)
					}
				*/
			}
		}
		// тесты для запросов
		fmt.Printf("тесты для запросов\r\n")
		//ll := fe.FrameDict["relation_наименование"]
		//fmt.Printf("%v \r\n", ll)
		qris := []frl.QueryRelationItem{
			{
				ObjectType: "relation",
				Object:     "наименование",
				Value:      frl.CreateValue("отношение_предикат сравнения"),
			},
			{
				ObjectType: "relation",
				Object:     "тип_отношения",
				Value:      frl.CreateValue("Предикаты отношения, связи (действия или состояния):"),
			},
		}
		fl, err := fe.QueryRelations(qris...)
		require.NoError(t, err)
		require.Len(t, fl, 1)

		if len(fl) > 0 {
			for i := range fl {
				f := fl[i]
				f.Print(output, true)
			}
		}

		// поиск в базе
		qri1 := frl.QueryRelationItem{
			ObjectType: "relation",
			Object:     "наименование",
			Value:      frl.CreateValue("отношение_предикат сравнения"),
		}
		ff, err := ns.Find(&qri1)
		require.NoError(t, err)
		for {
			frameId, slotName, slotProperty, slotValue, err := ff()
			if err != nil {
				break
			}
			fmt.Printf("frame_id %v, slot_name %v, slot_property %v, slot_value %v\r\n", frameId, slotName, slotProperty, slotValue)
			require.Equal(t, "наименование", slotName)
			require.Equal(t, "отношение_предикат сравнения", slotValue)
		}
		require.True(t, false)
	})
}

func TestTranslatorExec(t *testing.T) {
	debug.NewDebug()
	path := "../../data/scripts/lang/"

	files := []string{
		"test_вложенный_для_каждого.frm",
		"test_встроенных_функций.frm",
		"test_вызов_функции.frm",
		"test_вызов_функции_с_возвратом.frm",
		"test_для_каждого.frm",
		"test_если.frm",
		//		"test_нагрузочный.frm",
		//		"test_нагрузочный_памяти.frm",
		"test_пока.frm",
		"test_пока_вложенный.frm",
		"test_потока.frm",
		"test_потока_full.frm",
		"test_присваивание_константы_в_переменную.frm",
		"test_присваивание_константы_поиск_фрейма_в_переменную.frm",
		"test_присваивание_списка_в_переменную.frm",
		"test_форматировать.frm",
		"test_фрейм.frm",
	}
	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}
	output := print.NewOutput(printFunc)
	for _, fileIn := range files {
		fmt.Printf("file_in %v\r\n", path+fileIn)
		t.Run("exec "+fileIn, func(t *testing.T) {
			eb := exec.InitExecutorBase(0, output)
			e := exec.InitExecutor(eb, output, 0)
			err := e.Exec(path+fileIn, "пример1", "1", "2")
			require.NoError(t, err)
		})
	}
}

func TestTranslatorExecBad(t *testing.T) {
	debug.NewDebug()
	path := "../../data/scripts/lang/"

	files := []string{
		//		"test_плохой_файл.frm",
		"test_плохой_файл_простой.frm",
	}
	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}
	output := print.NewOutput(printFunc)
	for _, fileIn := range files {
		fmt.Printf("file_in %v\r\n", path+fileIn)
		t.Run("exec "+fileIn, func(t *testing.T) {
			eb := exec.InitExecutorBase(0, output)
			e := exec.InitExecutor(eb, output, 0)
			err := e.Exec(path+fileIn, "пример1", "1", "2")
			//require.NoError(t, err)
			require.ErrorContains(t, err, "translate error")
		})
	}
}

func TestTranslatorExecLineNum(t *testing.T) {
	debug.NewDebug()
	path := "../../data/scripts/lang/"

	files := []string{
		"test_пока.frm",
		//"test_фрейм.frm",
		//"test_вызов_функции_с_возвратом.frm",
	}
	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}
	output := print.NewOutput(printFunc)
	for _, fileIn := range files {
		fmt.Printf("file_in %v\r\n", path+fileIn)
		t.Run("exec "+fileIn, func(t *testing.T) {
			eb := exec.InitExecutorBase(0xff, output)
			e := exec.InitExecutor(eb, output, 0)
			err := e.Exec(path+fileIn, "пример1", "1", "2")
			require.NoError(t, err)
		})
	}
}

func TestStore(t *testing.T) {
	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}

	output := print.NewOutput(printFunc)

	// настраиваем окружение
	fe := frl.NewFrameEnvironment()
	//fe.FrameDict = make(map[string][]*frl.Frame)

	// заполняем
	ns, err := frl.NewStore("./Frames", output)
	require.NoError(t, err)

	// добавляем фрейм с отношением "имя"
	// надо добавить фрейм с определением отношения
	f := frl.NewFrame()
	// добавляем поле уникального идентификатора
	f.AddSlot("ID")
	id := uqe.UniqueValue(10)
	v, _ := f.Set("ID", id)
	fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
	relation := "наименование"
	f.AddSlot("отношение")
	v, _ = f.Set("отношение", relation)
	fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "отношение", Value: v})

	r, err := convertor.LoadRelation("relation.txt")
	require.NoError(t, err)

	if false {
		for i := range r {
			fmt.Printf("%v\r\n", r[i])
			for j := range r[i].RelationItem {
				fmt.Printf("\t%v\r\n", r[i].RelationItem[j])

			}
		}
	}

	for i := range r {
		for j := range r[i].RelationItem {
			// надо добавить фрейм с определением отношения
			f := frl.NewFrame()
			// добавляем поле уникального идентификатора
			f.AddSlot("ID")
			id := uqe.UniqueValue(10)
			v, _ := f.Set("ID", id)
			fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
			relation := r[i].RelationItem[j].Relation

			f.AddSlot("отношение")
			v, _ = f.Set("отношение", relation)
			fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "отношение", Value: v})

			f.AddSlot("наименование")
			v, _ = f.Set("наименование", "отношение_"+relation)
			fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: "наименование", Value: v})

			f.AddSlot("тип_отношения")
			f.Set("тип_отношения", r[i].RelationType)

			if len(r[i].RelationItem[j].RelationExample) > 0 {
				for k := range r[i].RelationItem[j].RelationExample {
					f.AddSlot("пример")
					f.Set("пример", r[i].RelationItem[j].RelationExample[k])
				}
			}

			/*
				for k, _ := range r[i].relationItem[j].relation_example {
					f.Set("пример", r[i].relationItem[j].relation_example[k])
				}
			*/
			fe.Frames = append(fe.Frames, f)
			//f.Print(true)
		}
	}

	if false {
		for i := range fe.Frames {
			f := fe.Frames[i]
			f.Print(output, true)
		}
	}
	for i := range fe.Frames {
		f := fe.Frames[i]
		ff := f.Iterate()
		frame_ids, err := f.GetValue("ID")
		if err != nil {
			fmt.Printf("get value %v\r\n", err)
			return
		}
		frame_id := frame_ids[0]
		for {
			s, ok, err := ff()
			if err != nil {
				break
			}
			ssl := s.GetSlotValue()
			slot_name := s.GetSlotName()
			slot_property := s.GetSlotProperty()
			// fmt.Printf("sn %v len ssl %v\r\n", sn, len(ssl))
			if slot_name == "ID" {
				// frame_id = ssl[0]
			} else {
				for j := range ssl {
					err := ns.SaveFrameRecord(frame_id, slot_name, slot_property, ssl[j], 0)
					if err != nil {
						fmt.Printf("SaveFrameRecord: %v\r\n", err)
						return
					}
				}
			}
			if ok {
				break
			}
		}
	}

	// проверяем

	if true {
		ns, err := frl.NewStore("./Frames", output)
		require.NoError(t, err)

		loader, err := ns.LoadFrameRecord(0)
		require.NoError(t, err)

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
				err := f.AddSlot("ID")
				require.NoError(t, err)

				v, err := f.Set("ID", frame_id)
				require.NoError(t, err)

				fe.AddRelations(f, frl.AddRelationItem{ObjectType: "frame", Object: "", Value: v})
				fe.Frames = append(fe.Frames, f)
				fm[frame_id] = f
			} else {
				f = fl[0]
			}
			err = f.AddSlot(slot_name)
			require.NoError(t, err)

			err = f.SetSlotProperty(slot_name, slot_property)
			require.NoError(t, err)

			_, err = f.SetValue(slot_name, slot_value)
			require.NoError(t, err)
			fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: slot_name, Value: slot_value})
		}
	}
	if false {
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

		fileIn := ""
		bb, err := os.ReadFile(fileIn)
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

		data, err := os.ReadFile(fileIn)
		if err != nil {
			fmt.Print(err)
			return
		}

		_, err = ie.TranslateText(fileIn, string(data), 0, ie.Output)
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
	require.True(t, false)
}
