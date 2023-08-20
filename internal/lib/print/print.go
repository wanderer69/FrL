package frl

type Output struct {
	Print func(fmt string, args ...any)
}

func NewOutput(print func(fmt string, args ...any)) *Output {
	o := Output{}
	o.Print = print
	return &o
}

func (o *Output) Printf(fmt string, args ...any) {
	if o.Print != nil {
		o.Print(fmt, args)
	}
}
