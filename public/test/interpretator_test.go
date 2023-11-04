package frl_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	exec "github.com/wanderer69/FrL/internal/lib/executor"
	"github.com/wanderer69/debug"
	print "github.com/wanderer69/tools/parser/print"
)

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
			e := exec.InitExecutor(eb, 0)
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
			e := exec.InitExecutor(eb, 0)
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
			e := exec.InitExecutor(eb, 0)
			err := e.Exec(path+fileIn, "пример1", "1", "2")
			require.NoError(t, err)
		})
	}
}
