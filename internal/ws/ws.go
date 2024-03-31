package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/wanderer69/FrL/public/entity"
	exec "github.com/wanderer69/FrL/public/executor"
)

type WSEnv struct {
	websocketStreem *websocket.Conn
}

func (wsEnv *WSEnv) Server(eb *exec.ExecutorBase, e *exec.Executor) *websocket.Conn {
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
		wsEnv.Listen(eb, e)
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

func (wsEnv *WSEnv) Send(messageType int, mo MessageOut) error {
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

func (wsEnv *WSEnv) Listen(eb *exec.ExecutorBase, e *exec.Executor) {
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
			err := wsEnv.Send(messageType, msgOut)
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

			sourceItems := []entity.SourceItem{}
			for i := range mi.Sources {
				breakpoints := []int{}
				for i := range mi.Breakpoints {
					if mi.Breakpoints[i].Name == mi.Sources[i].Name {
						breakpoints = append(breakpoints, mi.Breakpoints[i].Number)
					}
				}
				sourceItems = append(sourceItems, entity.SourceItem{
					Name:        mi.Sources[i].Name,
					SourceCode:  mi.Sources[i].Source,
					Breakpoints: breakpoints})
			}
			callback := func(name string, number int, data [][]string, variables []*entity.Variable) {
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
			wsEnv.Send(messageType, mo)

		case "run":
			command <- "run"

		case "check":
			var mo MessageOut
			mo.Id = mi.Id
			mo.Cmd = mi.Cmd
			mo.Result = "Ok"
			mo.Answer = "answer test"

			sourceItems := []entity.SourceItem{}
			for i := range mi.Sources {
				sourceItems = append(sourceItems, entity.SourceItem{
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
			wsEnv.Send(messageType, mo)
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
