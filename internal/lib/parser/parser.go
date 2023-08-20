package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	print "github.com/wanderer69/FrL/internal/lib/print"
	"github.com/wanderer69/debug"
)

func LevelShift(tab int) string {
	res := ""
	for i := 0; i < tab; i += 1 {
		res = res + "\t"
	}
	return res
}

func GetSlice(text string, beginPos int, endPos int) string {
	result := ""
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		w = width
		s1 := string(runeValue)
		if i >= beginPos {
			if i < endPos {
				result = result + s1
			}
		}
	}
	return result
}

type Item struct {
	Type         string
	Data         string
	LineNumBegin int
	LineNumEnd   int
}

// type ( [ {
func LoadLevel(levelAttribute string, posBegin int, posEnd int, lent int, text string, level int, flag string, lineNumBegin int) ([]Item, int, int, int) {
	prev := posBegin
	iPrev := posBegin
	debug.Alias("parser.LoadLevel.39").Printf("<-- text '%v'\r\n", text[posBegin:posEnd])
	flagS := false
	gSFlag := false
	itemsA := []Item{}
	lenText := len(text)

	fmt.Printf("LoadLevel text %v\r\n", text)

	lineNum := lineNumBegin
	for i, w := posBegin, 0; i < posEnd; {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		debug.Alias("parser.LoadLevel.42").Printf("%v%#U starts at byte position %d %v level %v\n", LevelShift(level), runeValue, i, string(runeValue), level)
		w = width
		s1 := string(runeValue)
		if !gSFlag {
			// если это скобка - ищем ответную
			if s1 == "(" {
				debug.Alias("parser.LoadLevel.40").Printf("i %v\r\n", i)
				if prev < i {
					ci := Item{Type: "symbols", Data: text[prev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				levelG := 1
				prev = i
				i = i + w
				flagSpace1 := false
				for {
					runeValue, width := utf8.DecodeRuneInString(text[i:])
					debug.Alias("parser.LoadLevel.42").Printf("!%#U starts at byte position %d %v\n", runeValue, i, string(runeValue))
					if w == 0 {
						debug.Alias("parser.LoadLevel").Printf("Error i %v len(text) %v\r\n", i, len(text))
						return itemsA, prev, 0, 1
					}
					w = width
					s1 := string(runeValue)
					if s1 == "\"" {
						if flagSpace1 {
							flagSpace1 = false
						} else {
							flagSpace1 = true
						}
					}
					if !flagSpace1 {
						if s1 == "(" {
							levelG = levelG + 1
						} else {
							if s1 == ")" {
								levelG = levelG - 1
								if levelG > 0 {

								} else {
									break
								}
							}
						}
					}
					i = i + w
					if i > lenText {
						debug.Alias("parser.LoadLevel").Printf("Error i %v len(text) %v\r\n", i, len(text))
						panic("Error!")
					}
				}
				debug.Alias("parser.LoadLevel.41").Printf("text '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
				debug.Alias("parser.LoadLevel.41").Printf("prev %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
				ci := Item{Type: "(", Data: text[prev+1 : i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
				itemsA = append(itemsA, ci)
				i = i + w
				prev = i
				iPrev = i
			} else if (s1 == " ") || (s1 == "\t") || (s1 == "\r") || (s1 == "\n") {
				// это разделитель! если до этого были отличные символы - строим строку.
				if s1 == "\n" {
					lineNum += 1
					if len(itemsA) == 0 {
						lineNumBegin = lineNum
					}
				}
				debug.Alias("parser.LoadLevel.40").Printf("space flag_s %v\r\n", flagS)
				if flagS {
					flagS = false
					debug.Alias("parser.LoadLevel.41").Printf("i_prev %v, i %v, w %v\r\n", iPrev, i, w)
					debug.Alias("parser.LoadLevel.40").Printf("text separator %v\r\n", text[iPrev:i])
					if iPrev < i {
						ci := Item{Type: "symbols", Data: text[iPrev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
						itemsA = append(itemsA, ci)
					}
				} else {
					// нет ничего это повторный разделитель
				}
				i = i + w
				iPrev = i
				prev = i
			} else if s1 == ")" {
				// скобка закрывающая завершаем работу и выходим
				debug.Alias("parser.LoadLevel.40").Printf("close bracket flag_s %v\r\n", flagS)
				if flagS {
					flagS = false
					debug.Alias("parser.LoadLevel.41").Printf("i_prev %v, i %v\r\n", iPrev, i)
					debug.Alias("parser.LoadLevel.40").Printf("text %v\r\n", text[iPrev:i])
					ci := Item{Type: ")", Data: text[iPrev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				return itemsA, i + 1, posEnd, 0
			} else if s1 == "[" {
				debug.Alias("parser.LoadLevel.41").Printf("i %v\r\n", i)
				if prev < i {
					ci := Item{Type: "symbols", Data: text[prev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				levelG := 1
				prev = i
				i = i + w
				for {
					runeValue, width := utf8.DecodeRuneInString(text[i:])
					debug.Alias("parser.LoadLevel.43").Printf("[>> !%#U starts at byte position %d %v\r\n", runeValue, i, string(runeValue))
					w = width
					s1 := string(runeValue)
					if s1 == "[" {
						levelG = levelG + 1
					} else {
						if s1 == "]" {
							levelG = levelG - 1
							if levelG > 0 {

							} else {
								break
							}
						}
					}
					i = i + w
					if i > len(text) {
						debug.Alias("parser.LoadLevel").Printf("Error i %v len(text) %v\r\n", i, len(text))
						panic("Error!")
					}
				}
				debug.Alias("parser.LoadLevel.41").Printf("text [ '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
				debug.Alias("parser.LoadLevel.41").Printf("prev [ %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
				ci := Item{Type: "[", Data: text[prev+1 : i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
				itemsA = append(itemsA, ci)
				i = i + w
				prev = i
				iPrev = i
			} else if s1 == "]" {
				// скобка закрывающая завершаем работу и выходим
				debug.Alias("parser.LoadLevel.40").Printf("close bracket flag_s %v\r\n", flagS)
				if flagS {
					debug.Alias("parser.LoadLevel.41").Printf("i_prev %v, i %v\r\n", iPrev, i)
					debug.Alias("parser.LoadLevel.40").Printf("text ] %v\r\n", text[iPrev:i])
					ci := Item{Type: "]", Data: text[iPrev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				return itemsA, i + 1, posEnd, 0
			} else if s1 == "{" {
				debug.Alias("parser.LoadLevel.41").Printf("i %v\r\n", i)
				if prev < i {
					ci := Item{Type: "symbols", Data: text[prev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				levelG := 1
				prev = i
				i = i + w
				for {
					runeValue, width := utf8.DecodeRuneInString(text[i:])
					debug.Alias("parser.LoadLevel.40").Printf("[>>> !%#U starts at byte position %d %v\r\n", runeValue, i, string(runeValue))
					w = width
					s1 := string(runeValue)
					if s1 == "{" {
						levelG = levelG + 1
					} else {
						if s1 == "}" {
							levelG = levelG - 1
							if levelG > 0 {

							} else {
								break
							}
						} else {
							if s1 == "\n" {
								lineNum += 1
							}
						}
					}
					i = i + w
					if i > len(text) {
						debug.Alias("parser.LoadLevel").Printf("Error i %v len(text) %v\r\n", i, len(text))
						panic("Error!")
					}
				}
				debug.Alias("parser.LoadLevel.41").Printf("text { '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
				debug.Alias("parser.LoadLevel.41").Printf("prev { %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
				ci := Item{Type: "{", Data: text[prev+1 : i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
				itemsA = append(itemsA, ci)
				i = i + w
				prev = i
				iPrev = i
				fmt.Printf("itemsA %v\r\n", itemsA)
			} else if s1 == "}" {
				// скобка закрывающая завершаем работу и выходим
				debug.Alias("parser.LoadLevel.40").Printf("close bracket flag_s %v\r\n", flagS)
				if flagS {
					debug.Alias("parser.LoadLevel.41").Printf("i_prev %v, i %v\r\n", iPrev, i)
					debug.Alias("parser.LoadLevel.40").Printf("text ] %v\r\n", text[iPrev:i])
					ci := Item{Type: "}", Data: text[iPrev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				fmt.Printf("!!!! itemsA %v\r\n", itemsA)
				return itemsA, i + 1, posEnd, 0
			} else if s1 == "\"" {
				// надо проверить, что
				if (i - iPrev) > 1 {
					ci := Item{Type: "symbols", Data: text[prev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				iPrev = i
				if flagS {
					flagS = false
				} else {
					flagS = true
				}
				gSFlag = true
				i = i + w
			} else {
				// это символ!!
				if !flagS {
					iPrev = i // + w
					flagS = true
				}
				i = i + w
			}
		} else {
			if s1 == "\"" {
				// это разделитель! если до этого были отличные символы - строим строку.
				debug.Alias("parser.LoadLevel.40").Printf("space 2 flag_s %v\r\n", flagS)
				i = i + w
				if flagS {
					flagS = false
					debug.Alias("parser.LoadLevel.41").Printf("i_prev %v, i %v, w %v\r\n", iPrev, i, w)
					debug.Alias("parser.LoadLevel.40").Printf("text %v\r\n", text[iPrev:i])
					ci := Item{Type: "string", Data: text[iPrev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				} else {
					// нет ничего это повторный разделитель
				}
				gSFlag = false
				iPrev = i
			} else {
				i = i + w
				flagS = true
			}
		}
		if s1 == levelAttribute {
			if !gSFlag {
				if iPrev < i-w {
					ci := Item{Type: "symbols", Data: text[iPrev : i-w], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
					itemsA = append(itemsA, ci)
				}
				// возврат значения
				return itemsA, i + 1, posEnd, 0
			}
		}
		if i >= posEnd {
			// все закончилось....
			if flagS {
				flagS = false
				debug.Alias("parser.LoadLevel.41").Printf("i_prev %v, i %v\r\n", iPrev, i)
				debug.Alias("parser.LoadLevel.40").Printf("text all %v\r\n", text[iPrev:i])
				ci := Item{Type: "symbols", Data: text[iPrev:i], LineNumBegin: lineNumBegin, LineNumEnd: lineNum}
				itemsA = append(itemsA, ci)
			}
		}
	}
	return itemsA, 0, 0, 0
}

type GrammaticItem struct {
	Type            string
	Mod             string // "", "[0]" - нулевой элемент
	Attribute       int    // 0 none, 1 ==,
	Value           string
	GrammRuleIdList []string // список идентификатор граматического правила
	Ender           string   // здесь лежит окончатель
}
type GrammaticRule struct {
	ID                string // идентификатор
	GrammaticItemList []GrammaticItem
}

type CurrentEnv struct {
	PiCnt          int
	State          int
	NextState      int
	ResultGenerate string
	Result         string
	I              int
	IntVars        map[string]int
	StringVars     map[string]string
}

type Env struct {
	mapGR          map[string]GrammaticRule
	baseGRArray    []GrammaticRule
	highLevelArray []string
	exprArray      []string
	parseFuncDict  map[string]func(pi ParseItem, env *Env, level int) (string, error)
	ErrorsList     []string
	CE             *CurrentEnv
	Struct         interface{}
	Debug          int
	Output         *print.Output
}

type ParseFuncDict struct {
	Dict map[string]ParseFunc
}

type ParseFunc func(pi ParseItem, env *Env, level int) (string, error)

func NewEnv() *Env {
	env := Env{}
	env.parseFuncDict = make(map[string]func(pi ParseItem, env *Env, level int) (string, error))
	env.mapGR = make(map[string]GrammaticRule)
	return &env
}

func (env *Env) SetEnv(mgr map[string]GrammaticRule, gr []GrammaticRule, hla []string, ea []string) {
	env.mapGR = mgr
	env.baseGRArray = gr
	env.highLevelArray = hla
	env.exprArray = ea
}

func (env *Env) SetBGRAEnv() {
	for _, v := range env.mapGR {
		env.baseGRArray = append(env.baseGRArray, v)
	}
}

func (env *Env) SetHLAEnv(hla []string) {
	env.highLevelArray = hla
}

func (env *Env) SetEAEnv(ea []string) {
	env.exprArray = ea
}

func (gr *GrammaticRule) AddRuleHandler(pf func(pi ParseItem, env *Env, level int) (string, error), env *Env) {
	_, ok := env.mapGR[gr.ID]
	if !ok {
		panic(fmt.Sprintf("Rule %v not exist", gr.ID))
	}
	env.parseFuncDict[gr.ID] = pf
}

func MakeRule(rule_name string, env *Env) *GrammaticRule {
	gr, ok := env.mapGR[rule_name]
	if ok {
		panic(fmt.Sprintf("Rule name %v exist", rule_name))
	}
	gr = GrammaticRule{ID: rule_name, GrammaticItemList: []GrammaticItem{}}
	env.mapGR[rule_name] = gr
	return &gr
}

func (gri *GrammaticRule) AddItemToRule(Type string, mod string, attribute int, value string, ender string, gRIdList []string, env *Env) {
	gr, ok := env.mapGR[gri.ID]
	if !ok {
		panic(fmt.Sprintf("Rule %v not exist", gr.ID))
	}
	gi := GrammaticItem{Type: Type, Mod: mod, Attribute: attribute, Value: value, Ender: ender, GrammRuleIdList: gRIdList}
	gr.GrammaticItemList = append(gr.GrammaticItemList, gi)
	env.mapGR[gr.ID] = gr
}

func LoadGrammaticRule(env *Env, name string, o *print.Output) (map[string]GrammaticRule, []GrammaticRule, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		o.Print("%v", err)
		return nil, nil, err
	}
	var gra []GrammaticRule
	err = json.Unmarshal(data, &gra)
	if err != nil {
		o.Print("error: %v", err)
		return nil, nil, err
	}
	map_gr := make(map[string]GrammaticRule)
	for i := range gra {
		map_gr[gra[i].ID] = gra[i]
	}
	env.mapGR = map_gr
	env.baseGRArray = gra
	return map_gr, gra, nil
}

func SaveGrammaticRule(env *Env, name string, o *print.Output) error {
	gr := env.baseGRArray
	data, err := json.MarshalIndent(&gr, "", "  ")
	if err != nil {
		o.Print("error: %v", err)
		return err
	}
	_ = ioutil.WriteFile(name, data, 0644)
	return nil
}

type FuncDictItem struct {
	Name string
	Args int
}

type Variable struct {
	Name  string
	Value string
}
type ParseItem struct {
	Items     []Item
	Gr        *GrammaticRule
	PI        [][]ParseItem
	Variables []Variable
}

type Condition_mask struct {
	Base   string
	Length int
	Type   int
}

func ParseArgList(si string, o *print.Output) []string {
	ender := ","

	sIn := si
	sError := []string{}
	level := 0
	resOut := true
	ia := []string{}
	lineNum := 0
	fmt.Printf("ParseArgList si %v\r\n", si)
	for {
		// читаем по грамматическим элементам
		items, posBeg, posEnd, err := LoadLevel(ender, 0, len(sIn), len(sIn), sIn, 0, "", lineNum)
		if err != 0 {
			o.Print("err %v\r\n", err)
			sError = append(sError, fmt.Sprintf("ParseArgList: error %v in level %v in process %v", err, level, sIn))
			break
		}
		// ищем подходящий грамматический элемент
		if len(items) > 0 {
			debug.Alias("parser.ParseArgList.20").Printf("l_1 %v\r\n", items)
			if len(items) > 1 {
				ss := ""
				for i := range items {
					ss = ss + items[i].Data
				}
				ia = append(ia, ss)
			} else {
				ia = append(ia, items[0].Data)
			}
			if !resOut {
				break
			}
			if posBeg >= posEnd {
				break
			}
			sIn = sIn[posBeg:]
			sIn = strings.Trim(sIn, " \r\n\t")
			if len(sIn) == 0 {
				break
			}
		}
	}
	if len(sError) > 0 {
		debug.Alias("parser.ParseArgList").Printf("s_error %v\r\n", sError)
	}
	return ia
}

func ParseArgListFull(si string, o *print.Output) [][]string {
	ender := ","

	sIn := si
	sError := []string{}
	level := 0
	resOut := true
	ia := [][]string{}
	lineNum := 0
	fmt.Printf("ParseArgListFull si %v\r\n", si)

	for {
		// читаем по грамматическим элементам
		items, posBeg, posEnd, err := LoadLevel(ender, 0, len(sIn), len(sIn), sIn, 0, "", lineNum)
		if err != 0 {
			o.Print("err %v\r\n", err)
			sError = append(sError, fmt.Sprintf("ParseArgListFull: error %v in level %v in process %v", err, level, sIn))
			break
		}
		// ищем подходящий грамматический элемент
		if len(items) > 0 {
			debug.Alias("parser.ParseArgListFull.1").Printf("l_1 %#v\r\n", items)
			if len(items) > 1 {
				for i := range items {
					ia = append(ia, []string{items[i].Type, items[i].Data})
				}
			} else {
				ia = append(ia, []string{items[0].Type, items[0].Data})
			}
			if !resOut {
				break
			}
			if posBeg >= posEnd {
				break
			}
			sIn = sIn[posBeg:]
			sIn = strings.Trim(sIn, " \r\n\t")
			if len(sIn) == 0 {
				break
			}
		}
	}
	if len(sError) > 0 {
		debug.Alias("parser.ParseArgListFull").Printf("sError %v\r\n", sError)
	}
	return ia
}

func LoadItems(dataIn string, env *Env, o *print.Output) (string, error) {
	s := dataIn

	// делим на элементы
	itemsList := strings.Split(s, "\n") // \r
	itemsListNew := []string{}
	for j := range itemsList {
		ls := strings.Trim(itemsList[j], " \r\n\t")
		if len(ls) > 0 {
			if (ls[0] == ';' && ls[1] == ';') || (ls[0] == '#') {
				// это комментарий
				// по идее можно добавить комментарий как отдельный оператор
			} else {
				// ищем в строке подстроку
				ls_l := strings.Split(ls, "//")
				if len(ls_l) > 1 {

					ls = ls_l[0]
				}
				// по идее надо добавлть признак строки а еще признак остановки при пошаговой отладке
				itemsListNew = append(itemsListNew, ls)
			}
		}
	}
	s = strings.Join(itemsListNew, "\r\n")

	checkRule := func(itemsInt []Item, grammRuleNameArray []string) (*GrammaticRule, bool) {
		var res *GrammaticRule
		flag := false
		for _, grammRuleName := range grammRuleNameArray {
			grammRuleItem := env.mapGR[grammRuleName]
			if len(grammRuleItem.GrammaticItemList) == len(itemsInt) {
				n := 0
				debug.Alias("parser.loadItems.2").Printf("grammRuleItem.ID %v grammRuleItem %v\r\n", grammRuleItem.ID, grammRuleItem)
				for i := range grammRuleItem.GrammaticItemList {
					gi := grammRuleItem.GrammaticItemList[i]
					item := itemsInt[i]
					tl := strings.Split(gi.Type, "|")
					flagBrake := false
					for ii := range tl {
						if tl[ii] == item.Type {
							debug.Alias("parser.loadItems.2").Printf("--- gi %v item %v\r\n", gi, item)
							switch gi.Attribute {
							case 0:
								n = n + 1
							case 1:
								switch gi.Mod {
								case "":
									if gi.Value == item.Data {
										n = n + 1
									}
								case "[0]":
									if gi.Value == item.Data[0:1] {
										n = n + 1
									}
								}
							}
							flagBrake = false
							break
						} else {
							flagBrake = true
						}
					}
					if flagBrake {
						break
					}
				}
				if n == len(grammRuleItem.GrammaticItemList) {
					res = &grammRuleItem
					flag = true
					break
				}
			}
		}
		return res, flag
	}

	var printParseItem func(pi ParseItem, level int)

	printParseItem = func(pi ParseItem, level int) {
		pi_cnt := 0
		st := ""
		for k := 0; k < level; k++ {
			st = st + "\t"
		}
		o.Print("%v%v:\r\n", st, pi.Gr.ID)
		for i := range pi.Items {
			if len(pi.Gr.GrammaticItemList) > 0 {
				if len(pi.Gr.GrammaticItemList[i].GrammRuleIdList) > 0 {
					for j := range pi.PI[pi_cnt] {
						printParseItem(pi.PI[pi_cnt][j], level+1)
					}
					pi_cnt = pi_cnt + 1
				} else {
					o.Print("%v\t%v \r\n", st, pi.Items[i])
				}
			} else {
				o.Print("%v\t %v\r\n", st, pi.Items[i])
			}
		}
		o.Print("\r\n")
	}

	var translate func(ender string, s_in string, gr_list []string, level int, lineNum int) ([]ParseItem, []string, bool)

	var generateParseItem func(pi ParseItem, envItem *Env, level int, lineNum int) (string, bool, error)

	generateParseItem = func(pi ParseItem, envItem *Env, level int, lineNum int) (string, bool, error) {
		result := ""
		st := ""
		for k := 0; k < level; k++ {
			st = st + "\t"
		}
		state := 0
		fn, ok := envItem.parseFuncDict[pi.Gr.ID]
		if !ok {
			return "", false, fmt.Errorf("handler for rule %v not defined", pi.Gr.ID)
		}
		flag := false
		result0 := ""
		for {
			switch state {
			case 0:
				envItem.CE.State = state
				res, err := fn(pi, envItem, level)
				if err != nil {
					return "", false, fmt.Errorf("error when translate %v: %v", pi.Gr.ID, err)
				}
				result0 = res
				state = envItem.CE.State
			case 100:
				resStr := ""
				for j := range pi.PI[envItem.CE.PiCnt] {
					ceOld := env.CE
					ce := CurrentEnv{}
					ce.PiCnt = envItem.CE.PiCnt
					ce.State = -1
					ce.NextState = -1
					ce.ResultGenerate = ""
					ce.Result = ""
					ce.I = 0
					ce.IntVars = make(map[string]int)
					ce.StringVars = make(map[string]string)
					envItem.CE = &ce
					res, status, err := generateParseItem(pi.PI[envItem.CE.PiCnt][j], envItem, level+1, lineNum)
					envItem.CE = ceOld
					if err != nil {
						return "", false, err
					}
					if status {
						resStr = resStr + res
					}
				}
				state = envItem.CE.NextState
				envItem.CE.ResultGenerate = resStr
			case 200:
				// список аргументов
				ParseItemCnt := envItem.CE.PiCnt
				argLst := ""
				pi1 := pi.PI[ParseItemCnt][0]
				for _, it := range pi1.Items {
					res := it.Data
					if len(argLst) > 0 {
						argLst = argLst + ", " + res
					} else {
						argLst = argLst + res
					}
				}
				itemR := strings.Trim(argLst, " \r\n\t") + ","
				pia, errList, ok := translate(",", itemR, envItem.exprArray, level, pi.Items[0].LineNumBegin)
				if ok {
					if false {
						for i := range pia {
							printParseItem(pia[i], 0)
						}
					}
					resultStr := ""
					for i := range pia {
						ceOld := env.CE
						ce := CurrentEnv{}
						ce.PiCnt = 0
						ce.State = -1
						ce.NextState = -1
						ce.ResultGenerate = ""
						ce.Result = ""
						ce.I = 0
						ce.IntVars = make(map[string]int)
						ce.StringVars = make(map[string]string)
						env.CE = &ce
						resOut, status, err := generateParseItem(pia[i], envItem, 0, pia[i].Items[0].LineNumBegin)
						env.CE = ceOld
						if err != nil {
							o.Print("Error while generate %v\r\n", pia[i])
							return "", false, err
						}
						if status {
							resultStr = resultStr + " " + resOut
						}
					}
					envItem.CE.ResultGenerate = resultStr
				} else {
					o.Print("Error while translate %v\r\n", itemR)
					envItem.ErrorsList = append(envItem.ErrorsList, errList...)
					envItem.ErrorsList = append(envItem.ErrorsList, []string{fmt.Sprintf("Error while translate %v\r\n", itemR)}...)
					return "", false, fmt.Errorf("error while translate %v", itemR)
				}
				state = envItem.CE.NextState
			case 1000:
				if len(result0) > 0 {
					result = result0
				}
				flag = true
			default:
				envItem.CE.State = state
				res, err := fn(pi, envItem, level)
				if err != nil {
					return "", false, fmt.Errorf("error when translate %v: %v", pi.Gr.ID, err)
				}
				result0 = ""
				result = res
				state = envItem.CE.State
			}
			if flag {
				break
			}
		}
		return result, true, nil
	}

	translate = func(ender string, sIn string, grammRuleList []string, level int, lineNum int) ([]ParseItem, []string, bool) {
		resOut := true
		sError := []string{}
		pia := []ParseItem{}
		//lineNum := 0
		for {
			// читаем по грамматическим элементам
			items, posBeg, posEnd, err := LoadLevel(ender, 0, len(sIn), len(sIn), sIn, 0, "", lineNum)
			if err != 0 {
				o.Print("err %v\r\n", err)
				sError = append(sError, fmt.Sprintf("error %v in level %v in process %v", err, level, sIn))
				resOut = false
				break
			}
			max := 0
			for i := range items {
				//fmt.Printf("---> %v %v %v %v\r\n", i, items[i].LineNumBegin, items[i].LineNumEnd, items[i].Data)
				if items[i].LineNumEnd > max {
					max = items[i].LineNumEnd
				}
			}
			lineNum = max
			// ищем подходящий грамматический элемент
			if len(items) > 0 {
				debug.Alias("parser.loadItems.20").Printf("l_1 %v\r\n", items)
				gr, res := checkRule(items, grammRuleList)
				if res {
					debug.Alias("parser.loadItems.20").Printf("ID %v\r\n", gr.ID)
					pi := ParseItem{Items: items, Gr: gr}
					if len(gr.GrammaticItemList) > 0 {
						for i, grammRuleItem := range gr.GrammaticItemList {
							if len(grammRuleItem.GrammRuleIdList) > 0 {
								str := items[i].Data
								flagGRBreak := false
								if len(grammRuleItem.GrammRuleIdList) == 1 {
									ia := []Item{}
									id := ""
									switch grammRuleItem.GrammRuleIdList[0] {
									case "список":
										item := Item{Type: "symbols", Data: strings.Trim(str, " \r\n\t"),
											LineNumBegin: items[i].LineNumBegin, LineNumEnd: items[i].LineNumEnd}
										ia = append(ia, item)
										id = grammRuleItem.GrammRuleIdList[0]
										flagGRBreak = true
									case "список_аргументов":
										al := ParseArgList(str, o)
										for _, v := range al {
											item := Item{Type: "symbols", Data: strings.Trim(v, " \r\n\t"),
												LineNumBegin: items[i].LineNumBegin, LineNumEnd: items[i].LineNumEnd}
											ia = append(ia, item)
										}
										id = grammRuleItem.GrammRuleIdList[0]
										flagGRBreak = true
									}
									if flagGRBreak {
										piNew := ParseItem{Items: ia, Gr: &GrammaticRule{ID: id}}
										pi.PI = append(pi.PI, []ParseItem{piNew})
									}
								}
								if !flagGRBreak {
									piNew, errList, res := translate(grammRuleItem.Ender, str, grammRuleItem.GrammRuleIdList, level+1, items[i].LineNumBegin)
									if res {
										pi.PI = append(pi.PI, piNew)
									} else {
										resOut = false
										sError = append(sError, errList...)
										break
									}
								}
							}
						}
					}
					pia = append(pia, pi)
				} else {
					sError = append(sError, fmt.Sprintf("error %v in level %v in process %v", err, level, items))
					resOut = false
					break
				}
				if !resOut {
					break
				}
				if posBeg >= posEnd {
					break
				}
				sIn = sIn[posBeg:]
				//sIn = strings.Trim(sIn, " \r\n\t")
				if len(sIn) == 0 {
					break
				}
			} else {
				sError = append(sError, fmt.Sprintf("error %v in level %v in process %v", err, level, sIn))
				if len(pia) == 0 {
					resOut = false
				}
				break
			}
		}
		return pia, sError, resOut
	}
	pia, errList, res := translate(";", s, env.highLevelArray, 0, 0)
	if res {
		if env.Debug > 0 {
			for i := range pia {
				printParseItem(pia[i], 0)
			}
		}
		result := ""

		for i := range pia {
			ce := CurrentEnv{}
			ce.PiCnt = 0
			ce.State = -1
			ce.NextState = -1
			ce.ResultGenerate = ""
			ce.Result = ""
			ce.I = 0
			ce.IntVars = make(map[string]int)
			ce.StringVars = make(map[string]string)
			env.CE = &ce
			resOut, status, err := generateParseItem(pia[i], env, 0, pia[i].Items[0].LineNumEnd)
			if err != nil {
				o.Print("Error while translate %v\r\n", err)
				return "", err
			}
			if status {
				result = result + "\r\n" + resOut
			}
		}
		return result, nil
	} else {
		ss := ""
		for i := range errList {
			s := fmt.Sprintf("%v\r\n", errList[i])
			ss = ss + s
		}
		return "", errors.New(ss)
	}
}

func (env *Env) ParseFile(inFileName string, outFileName string, o *print.Output) (string, error) {
	data, err := ioutil.ReadFile(inFileName)
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	result, err := LoadItems(string(data), env, o)
	if err != nil {
		return "", err
	}
	if len(outFileName) > 0 {
		err = ioutil.WriteFile(outFileName, []byte(result), 0644)
		if err != nil {
			panic(err)
		}
	}
	return result, nil
}

func (env *Env) ParseString(inData string, o *print.Output) (string, error) {
	result, err := LoadItems(string(inData), env, o)
	if err != nil {
		return "", err
	}
	return result, nil
}
