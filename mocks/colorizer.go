package mocks

type Colorizer struct {
	ColorizeFn       func(string, string) (string, error)
	BlackFn          func(string) string
	RedFn            func(string) string
	GreenFn          func(string) string
	YellowFn         func(string) string
	BlueFn           func(string) string
	MagentaFn        func(string) string
	CyanFn           func(string) string
	WhiteFn          func(string) string
	ColorizeFnCalled bool
	BlackFnCalled    bool
	RedFnCalled      bool
	GreenFnCalled    bool
	YellowFnCalled   bool
	BlueFnCalled     bool
	MagentaFnCalled  bool
	CyanFnCalled     bool
	WhiteFnCalled    bool
}

func (c *Colorizer) Colorize(colorName string, text string) (string, error) {
	c.ColorizeFnCalled = true
	return c.ColorizeFn(colorName, text)
}

func (c *Colorizer) Black(s string) string {
	c.BlackFnCalled = true
	return c.BlackFn(s)
}

func (c *Colorizer) Red(s string) string {
	c.RedFnCalled = true
	return c.RedFn(s)
}

func (c *Colorizer) Green(s string) string {
	c.GreenFnCalled = true
	return c.GreenFn(s)
}

func (c *Colorizer) Yellow(s string) string {
	c.YellowFnCalled = true
	return c.YellowFn(s)
}

func (c *Colorizer) Blue(s string) string {
	c.BlueFnCalled = true
	return c.BlueFn(s)
}

func (c *Colorizer) Magenta(s string) string {
	c.MagentaFnCalled = true
	return c.MagentaFn(s)
}

func (c *Colorizer) Cyan(s string) string {
	c.CyanFnCalled = true
	return c.CyanFn(s)
}

func (c *Colorizer) White(s string) string {
	c.WhiteFnCalled = true
	return c.WhiteFn(s)
}

var defaultColorFunc = func(s string) string {
	return s
}

var defaultColorizeFunc = func(colorName, text string) (string, error) {
	return text, nil
}

func NewColorizer() *Colorizer {
	return &Colorizer{
		ColorizeFn: defaultColorizeFunc,
		BlackFn:    defaultColorFunc,
		RedFn:      defaultColorFunc,
		GreenFn:    defaultColorFunc,
		YellowFn:   defaultColorFunc,
		BlueFn:     defaultColorFunc,
		MagentaFn:  defaultColorFunc,
		CyanFn:     defaultColorFunc,
		WhiteFn:    defaultColorFunc,
	}
}
