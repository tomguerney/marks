package mocks

type Opener struct {
	OpenFn       func(string, string) error
	OpenFnCalled bool
}

func NewOpener() *Opener {
	return &Opener{
		OpenFn: defaultOpenFn,
	}
}

func (c *Opener) Open(s1, s2 string) error {
	c.OpenFnCalled = true
	return c.OpenFn(s1, s2)
}

var defaultOpenFn = func(string, string) error {
	return nil
}
