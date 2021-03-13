package mocks

type Clipper struct {
	CopyFn       func(string) error
	CopyFnCalled bool
}

func NewClipper() *Clipper {
	return &Clipper{
		CopyFn: defaultCopyFn,
	}
}

func (c *Clipper) Copy(s string) error {
	c.CopyFnCalled = true
	return c.CopyFn(s)
}

var defaultCopyFn = func(string) error {
	return nil
}
