package mocks

type Prompter struct {
	SelectFn        func(string, []string) (int, error)
	SelectFnCalled  bool
	ConfirmFn       func(string) bool
	ConfirmFnCalled bool
}

func NewPrompter() *Prompter {
	return &Prompter{
		SelectFn:  defaultSelectFn,
		ConfirmFn: defaultConfirmFn,
	}
}

func (p *Prompter) Select(label string, table []string) (i int, err error) {
	p.SelectFnCalled = true
	return p.SelectFn(label, table)
}

var defaultSelectFn = func(label string, table []string) (i int, err error) {
	return 0, nil
}

func (p *Prompter) Confirm(label string) bool {
	p.ConfirmFnCalled = true
	return p.ConfirmFn(label)
}

var defaultConfirmFn = func(label string) bool {
	return true
}
