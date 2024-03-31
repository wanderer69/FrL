package convertor

import (
	"fmt"
	"os"
	"strings"
)

type RelationItem struct {
	Relation            string
	RelationDescription string
	RelationExample     []string
}

type Relation struct {
	RelationType string
	RelationItem []RelationItem
}

func LoadRelation(file_name string) ([]Relation, error) {
	result := []Relation{}
	data, err := os.ReadFile(file_name)
	if err != nil {
		fmt.Print(err)
		return result, err
	}
	s := string(data)
	sl := strings.Split(s, "\n")
	state := 0

	var r Relation
	r = Relation{}

	for i := range sl {
		ss := strings.Trim(sl[i], " \r\n")
		if len(ss) == 0 {
			continue
		}
		if ss[0] == '#' {
			continue
		}
		ssl := strings.Split(ss, "-")
		if len(ssl) < 2 {
			// отношение
			/*
				if state == 0 {
				} else {
				}
			*/
			r.RelationItem = append(r.RelationItem, RelationItem{Relation: strings.Trim(ss, " ")})
			state = 1
		} else {
			ss_ := strings.Trim(ssl[0], " ")
			switch ss_ {
			case "вид":
				if len(r.RelationItem) > 0 {
					// сохраняем предыдущее отношение
					result = append(result, r)
					state = 0
				}

				r = Relation{}
				r.RelationType = strings.Trim(ssl[1], " ")

			case "пример":
				if state == 1 {
					state = 2
				} else {
					if state == 0 {
						r.RelationItem = append(r.RelationItem, RelationItem{Relation: r.RelationType})
						state = 2
						//} else {

					}
				}
				r.RelationItem[len(r.RelationItem)-1].RelationExample = append(r.RelationItem[len(r.RelationItem)-1].RelationExample, strings.Trim(ssl[1], " "))

			case "описание":
				if state == 1 {
					state = 2
				} else {
					if state == 0 {
						r.RelationItem = append(r.RelationItem, RelationItem{Relation: r.RelationType})
						state = 2
						//} else {

					}
				}
				r.RelationItem[len(r.RelationItem)-1].RelationDescription = strings.Trim(ssl[1], " ")

			}
			/*
				if ss_ == "вид" {
					if len(r.RelationItem) > 0 {
						// сохраняем предыдущее отношение
						result = append(result, r)
						state = 0
					}

					r = Relation{}
					r.RelationType = strings.Trim(ssl[1], " ")
				} else {
					if ss_ == "пример" {
						if state == 1 {
							state = 2
						} else {
							if state == 0 {
								r.RelationItem = append(r.RelationItem, RelationItem{Relation: r.RelationType})
								state = 2
							//} else {

							}
						}
						r.RelationItem[len(r.RelationItem)-1].RelationExample = append(r.RelationItem[len(r.RelationItem)-1].RelationExample, strings.Trim(ssl[1], " "))
					} else {
						if ss_ == "описание" {
							if state == 1 {
								state = 2
							} else {
								if state == 0 {
									r.RelationItem = append(r.RelationItem, RelationItem{Relation: r.RelationType})
									state = 2
								//} else {

								}
							}
							r.RelationItem[len(r.RelationItem)-1].RelationDescription = strings.Trim(ssl[1], " ")
						//} else {

						}
					}
				}
			*/
		}
	}
	if len(r.RelationItem) > 0 {
		// сохраняем предыдущее отношение
		result = append(result, r)
	}

	return result, nil
}
