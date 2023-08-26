package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	exec "github.com/wanderer69/FrL/internal/lib/executor"
	"github.com/wanderer69/debug"
	print "github.com/wanderer69/tools/parser/print"
)

type wsEnv struct {
	websocketStreem *websocket.Conn
}

func server(wsEnv *wsEnv, eb *exec.ExecutorBase, e *exec.Executor) *websocket.Conn {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var websocketStreem *websocket.Conn
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		websocketStreem, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Websocket Connected!")
		wsEnv.websocketStreem = websocketStreem
		listen(wsEnv, eb, e)
	})
	return websocketStreem
}

type Call struct {
	CallName string
	Args     []string
}

type Source struct {
	Name   string
	Source string
}

type Breakpoint struct {
	Name   string
	Number int
}

type Variable struct {
	Function string
	Variable string
	Type     string
	Value    string
}

type Message struct {
	Id          string
	Cmd         string
	Call        Call
	Sources     []Source
	Breakpoints []Breakpoint
	Result      string
	Answer      string
}

type MessageOut struct {
	Id        string
	Cmd       string
	Result    string
	Answer    string
	Variables []Variable
}

func send(wsEnv *wsEnv, messageType int, mo MessageOut) error {
	messageResponse, err := json.MarshalIndent(mo, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	if err := wsEnv.websocketStreem.WriteMessage(messageType, messageResponse); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func listen(wsEnv *wsEnv, eb *exec.ExecutorBase, e *exec.Executor) {
	command := make(chan string)
	defer close(command)
	messageOut := make(chan MessageOut)
	defer close(messageOut)
	var messageType int

	go func() {
		for {
			msgOut := <-messageOut
			/*
				messageResponse, err := json.MarshalIndent(msgOut, "", "  ")
				if err != nil {
					log.Println(err)
					return
				}

				if err := wsEnv.websocketStreem.WriteMessage(messageType, messageResponse); err != nil {
					log.Println(err)
					return
				}
			*/
			err := send(wsEnv, messageType, msgOut)
			if err != nil {
				break
			}
		}
	}()

	for {
		mType, messageContent, err := wsEnv.websocketStreem.ReadMessage()
		messageType = mType
		if err != nil {
			log.Println(err)
			return
		}

		var mi Message
		err = json.Unmarshal(messageContent, &mi)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("mi %#v\r\n", mi)

		switch mi.Cmd {
		case "execute":
			var mo MessageOut
			mo.Id = mi.Id
			mo.Cmd = mi.Cmd
			mo.Result = "Ok"
			mo.Answer = "answer test"

			sourceItems := []exec.SourceItem{}
			for i := range mi.Sources {
				breakpoints := []int{}
				for i := range mi.Breakpoints {
					if mi.Breakpoints[i].Name == mi.Sources[i].Name {
						breakpoints = append(breakpoints, mi.Breakpoints[i].Number)
					}
				}
				sourceItems = append(sourceItems, exec.SourceItem{
					Name:        mi.Sources[i].Name,
					SourceCode:  mi.Sources[i].Source,
					Breakpoints: breakpoints})
			}
			callback := func(name string, number int, data [][]string) {
				for {
					var moCB MessageOut
					moCB.Id = mi.Id
					moCB.Cmd = "breakpoint"
					moCB.Result = "Ok"
					moCB.Answer = fmt.Sprintf("breakpoint: %v %v", name, number)
					for i := range data {
						v := Variable{
							Function: data[i][0],
							Variable: data[i][1],
							Type:     data[i][2],
							Value:    data[i][3],
						}
						moCB.Variables = append(moCB.Variables, v)
					}
					/*
						messageResponse, err := json.MarshalIndent(mo, "", "  ")
						if err != nil {
							log.Println(err)
							return
						}

						if err := wsEnv.websocketStreem.WriteMessage(messageType, messageResponse); err != nil {
							log.Println(err)
							return
						}
					*/
					messageOut <- moCB

					c := <-command
					if c == "run" {
						break
					}
				}
			}
			go func() {
				var moProc MessageOut
				moProc.Id = mi.Id
				moProc.Cmd = "stop"
				moProc.Result = "Ok"
				err := e.ExecuteFuncWithManyFiles(sourceItems, callback, mi.Call.CallName, mi.Call.Args)
				if err != nil {
					log.Println(err)
					moProc.Result = "Error"
					moProc.Answer = fmt.Sprintf("execute: %v", err)
				}
				messageOut <- moProc
				/*
					messageResponse, err := json.MarshalIndent(mo, "", "  ")
					if err != nil {
						log.Println(err)
						return
					}

					if err := wsEnv.websocketStreem.WriteMessage(messageType, messageResponse); err != nil {
						log.Println(err)
						return
					}
				*/
			}()
			send(wsEnv, messageType, mo)

		case "run":
			command <- "run"

		case "check":
			var mo MessageOut
			mo.Id = mi.Id
			mo.Cmd = mi.Cmd
			mo.Result = "Ok"
			mo.Answer = "answer test"

			sourceItems := []exec.SourceItem{}
			for i := range mi.Sources {
				sourceItems = append(sourceItems, exec.SourceItem{
					Name:       mi.Sources[i].Name,
					SourceCode: mi.Sources[i].Source,
				})
			}
			err := e.TranslateManyFiles(sourceItems)
			if err != nil {
				log.Println(err)
				mo.Result = "Error"
				mo.Answer = fmt.Sprintf("check: %v", err)
			}
			send(wsEnv, messageType, mo)
		}
		/*
			messageResponse, err := json.MarshalIndent(mo, "", "  ")
			if err != nil {
				log.Println(err)
				return
			}

			if err := wsEnv.websocketStreem.WriteMessage(messageType, messageResponse); err != nil {
				log.Println(err)
				return
			}
		*/
	}
}

func main() {
	//fmt.Printf("file_in %v\r\n", fileIn)
	debug.NewDebug()

	wse := &wsEnv{}

	// var websocketStreem *websocket.Conn
	var callback func(wsEnv *wsEnv, messageType int, mo MessageOut) error

	toPrint := make(chan MessageOut)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	printByTime := func() {
		// flag
		n := 0
		queueSize := 10
		queue := make([]MessageOut, 0, queueSize)
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
							err := callback(wse, 1, msg)
							if err != nil {
								fmt.Printf("error %v\r\n", err)
							}
						}
						queue = make([]MessageOut, 0, queueSize)
					}
				}
			case msg := <-toPrint:
				if len(msg.Id) == 0 {
					break
				}
				if len(queue) == queueSize {
					for _, msgo := range queue {
						err := callback(wse, 1, msgo)
						if err != nil {
							fmt.Printf("error %v\r\n", err)
						}
					}
					err := callback(wse, 1, msg)
					if err != nil {
						fmt.Printf("error %v\r\n", err)
					}
					queue = make([]MessageOut, 0, queueSize)
				} else {
					queue = append(queue, msg)
				}
			}
		}
	}
	printFunc := func(frm string, args ...any) {
		fmt.Printf(frm, args...)
		if callback != nil {
			str := fmt.Sprintf(frm, args...)
			var mo MessageOut
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
		}
	}

	go printByTime()

	output := print.NewOutput(printFunc)

	eb := exec.InitExecutorBase(0, output)
	e := exec.InitExecutor(eb, 0)

	// websocketStreem =
	server(wse, eb, e)
	callback = send

	http.ListenAndServe(":8083", nil)
}
