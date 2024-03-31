package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/wanderer69/debug"

	"github.com/wanderer69/FrL/internal/ws"
	exec "github.com/wanderer69/FrL/public/executor"
	frl "github.com/wanderer69/FrL/public/lib"
	print "github.com/wanderer69/tools/parser/print"
)

func main() {
	debug.NewDebug()

	var port int
	flag.IntVar(&port, "port", 8083, "server port")

	wse := &ws.WSEnv{}

	toPrint := make(chan ws.MessageOut)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	printByTime := func() {
		// flag
		n := 0
		queueSize := 10
		queue := make([]ws.MessageOut, 0, queueSize)
		for {
			select {
			case <-ticker.C:
				n = n + 1
				if n == 200 {
					//v := t.String()
					//log.Println("write:", v)
					n = 0
					if len(queue) == queueSize {
						for _, msg := range queue {
							err := wse.Send(1, msg)
							if err != nil {
								fmt.Printf("error %v\r\n", err)
							}
						}
						queue = make([]ws.MessageOut, 0, queueSize)
					}
				}
			case msg := <-toPrint:
				if len(msg.Id) == 0 {
					break
				}
				if len(queue) == queueSize {
					for _, msgo := range queue {
						err := wse.Send(1, msgo)
						if err != nil {
							fmt.Printf("error %v\r\n", err)
						}
					}
					err := wse.Send(1, msg) // callback(wse, 1, msg)
					if err != nil {
						fmt.Printf("error %v\r\n", err)
					}
					queue = make([]ws.MessageOut, 0, queueSize)
				} else {
					queue = append(queue, msg)
				}
			}
		}
	}
	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
		str := fmt.Sprintf(frm, args...)
		var mo ws.MessageOut
		mo.Id = "0"
		mo.Cmd = "print"
		mo.Result = "Ok"
		mo.Answer = str
		toPrint <- mo
	}

	go printByTime()

	output := print.NewOutput(printFunc)

	translatePrintFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}
	outputTranslate := print.NewOutput(translatePrintFunc)

	basePrintFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
	}
	outputBaseTranslate := print.NewOutput(basePrintFunc)

	eb := exec.InitExecutorBase(0, outputBaseTranslate)
	extFunctions := make(map[string]func(args []*frl.Value) ([]*frl.Value, bool, error))
	e := exec.InitExecutor(eb, extFunctions, output, outputTranslate, 0)

	wse.Server(eb, e)

	path := fmt.Sprintf(":%v", port)
	http.ListenAndServe(path, nil)
}
