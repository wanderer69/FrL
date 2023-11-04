package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wanderer69/debug"

	"github.com/wanderer69/FrL/internal/ws"
	exec "github.com/wanderer69/FrL/public/executor"
	print "github.com/wanderer69/tools/parser/print"
)

func main() {
	//fmt.Printf("file_in %v\r\n", fileIn)
	debug.NewDebug()

	wse := &ws.WSEnv{}

	// var websocketStreem *websocket.Conn
	//var callback func(wsEnv *ws.WSEnv, messageType int, mo ws.MessageOut) error

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
		//if callback != nil {
		str := fmt.Sprintf(frm, args...)
		var mo ws.MessageOut
		mo.Id = "0"
		mo.Cmd = "print"
		mo.Result = "Ok"
		mo.Answer = str
		toPrint <- mo
		/*
			err := callback(wse, 1, mo)
			if err != nil {
				fmt.Printf("error %v\r\n", err)
			}
		*/
		//}
	}

	go printByTime()

	output := print.NewOutput(printFunc)

	eb := exec.InitExecutorBase(0, output)
	e := exec.InitExecutor(eb, 0)

	// websocketStreem =
	wse.Server(eb, e)
	//callback = send

	http.ListenAndServe(":8083", nil)
}
