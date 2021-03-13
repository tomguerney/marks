package mocks

import (
	"strings"

	"github.com/abitofoldtom/marks/marks"
)

type Printer struct {
	MsgFn                      func(string, ...interface{})
	ErrorFn                    func(string, ...interface{})
	TabulateFn                 func([]*marks.Mark) ([]string, error)
	FullMarkFn                 func(*marks.Mark) (string, error)
	FullMarkWithFieldsFn       func(*marks.Mark) (string, error)
	IdFn                       func(string) (string, error)
	UrlFn                      func(string) (string, error)
	TagsFn                     func([]string) (string, error)
	BrowserFn                  func(string) (string, error)
	MsgFnCalled                bool
	ErrorFnCalled              bool
	TabulateFnCalled           bool
	FullMarkFnCalled           bool
	FullMarkWithFieldsFnCalled bool
	IdFnCalled                 bool
	UrlFnCalled                bool
	TagsFnCalled               bool
	BrowserFnCalled            bool
}

func NewPrinter() *Printer {
	return &Printer{
		MsgFn:                defaultMsgFn,
		ErrorFn:              defaultErrorFn,
		TabulateFn:           defaultTabulateFn,
		FullMarkFn:           defaultFullMarkFn,
		FullMarkWithFieldsFn: defaultFullMarkWithFieldsFn,
		IdFn:                 defaultIdFn,
		UrlFn:                defaultUrlFn,
		TagsFn:               defaultTagsFn,
		BrowserFn:            defaultBrowserFn,
	}
}

func (p *Printer) Msg(s string, i ...interface{}) {
	p.MsgFnCalled = true
	p.MsgFn(s, i)
}

func (p *Printer) Error(s string, i ...interface{}) {
	p.ErrorFnCalled = true
	p.ErrorFn(s, i)
}

func (p *Printer) Tabulate(mks []*marks.Mark) ([]string, error) {
	p.TabulateFnCalled = true
	return p.TabulateFn(mks)
}

func (p *Printer) FullMark(m *marks.Mark) (string, error) {
	p.FullMarkFnCalled = true
	return p.FullMarkFn(m)
}

func (p *Printer) FullMarkWithFields(m *marks.Mark) (string, error) {
	p.FullMarkWithFieldsFnCalled = true
	return p.FullMarkWithFieldsFn(m)
}

func (p *Printer) Id(s string) (string, error) {
	p.IdFnCalled = true
	return p.Id(s)
}

func (p *Printer) Url(s string) (string, error) {
	p.UrlFnCalled = true
	return p.UrlFn(s)
}

func (p *Printer) Tags(sl []string) (string, error) {
	p.TagsFnCalled = true
	return p.TagsFn(sl)
}

func (p *Printer) Browser(s string) (string, error) {
	p.BrowserFnCalled = true
	return p.BrowserFn(s)
}

var defaultMsgFn = func(s string, i ...interface{}) {
	//do nothing
}

var defaultErrorFn = func(s string, i ...interface{}) {
	//do nothing
}

var defaultTabulateFn = func([]*marks.Mark) ([]string, error) {
	return []string{}, nil
}

var defaultFullMarkFn = func(*marks.Mark) (string, error) {
	return "full mark", nil
}

var defaultFullMarkWithFieldsFn = func(*marks.Mark) (string, error) {
	return "full mark with fields", nil
}

var defaultIdFn = func(s string) (string, error) {
	return s, nil
}

var defaultUrlFn = func(s string) (string, error) {
	return s, nil
}

var defaultTagsFn = func(sl []string) (string, error) {
	return strings.Join(sl, " "), nil
}

var defaultBrowserFn = func(s string) (string, error) {
	return s, nil
}
