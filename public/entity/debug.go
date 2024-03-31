package entity

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
