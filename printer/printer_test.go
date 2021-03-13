package printer

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/abitofoldtom/marks/marks"

	"github.com/abitofoldtom/marks/mocks"
)

func NewTestPrinter() *printer {
	return &printer{&strings.Builder{}, mocks.NewConfig(), mocks.NewColorizer()}
}

func getMultilineString(lines []string) string {
	builder := strings.Builder{}
	for _, line := range lines {
		builder.WriteString(fmt.Sprintf("%v\n", line))
	}
	return builder.String()
}

func TestMessage(t *testing.T) {
	msg := "test message"
	expected := fmt.Sprintf("%v\n", msg)
	writeFn := mocks.NewTestWriteFn(func(actual string) {
		if expected != actual {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	})
	w := mocks.NewWriter()
	w.WriteFn = writeFn
	p := NewTestPrinter()
	p.out = w
	p.Msg(msg)
}

func TestMessageWithArgument(t *testing.T) {
	msg := "test message %v"
	arg := "arg"
	expected := fmt.Sprintf("%v\n", fmt.Sprintf(msg, arg))
	writeFn := mocks.NewTestWriteFn(func(actual string) {
		if expected != actual {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	})
	w := mocks.NewWriter()
	w.WriteFn = writeFn
	p := NewTestPrinter()
	p.out = w
	p.Msg(msg, arg)
}

func TestError(t *testing.T) {
	errMsg := "test error message"
	expected := fmt.Sprintf("Error: %v\n", errMsg)
	writeFn := mocks.NewTestWriteFn(func(actual string) {
		if expected != actual {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	})
	w := mocks.NewWriter()
	w.WriteFn = writeFn
	p := NewTestPrinter()
	p.out = w
	p.Error(errMsg)
}

func TestErrorWithArgument(t *testing.T) {
	errMsg := "test message %v"
	arg := "arg"
	expected := fmt.Sprintf("Error: %v\n", fmt.Sprintf(errMsg, arg))
	writeFn := mocks.NewTestWriteFn(func(actual string) {
		if expected != actual {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	})
	w := mocks.NewWriter()
	w.WriteFn = writeFn
	p := NewTestPrinter()
	p.out = w
	p.Error(errMsg, arg)
}

func TestTabulate(t *testing.T) {
	mks := []*marks.Mark{
		&marks.Mark{
			Id:   "Abc News",
			Url:  "https://www.abc.net.au/news/",
			Tags: []string{"news", "current affairs"},
		},
		&marks.Mark{
			Id:   "Google",
			Url:  "https://www.google.com",
			Tags: []string{"search"},
		},
		&marks.Mark{
			Id:   "BBC News",
			Url:  "https://www.bbc.com/news",
			Tags: []string{"news", "uk"},
		},
	}
	expected := []string{
		"Abc News    https://www.abc.net.au/news/    [news, current affairs]",
		"Google      https://www.google.com          [search]",
		"BBC News    https://www.bbc.com/news        [news, uk]",
	}
	p := NewTestPrinter()
	actual, err := p.Tabulate(mks)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, received: %v", expected, actual)
	}
}

func TestTabulateWithColorize(t *testing.T) {
	mks := []*marks.Mark{
		&marks.Mark{
			Id:   "Abc News",
			Url:  "https://www.abc.net.au/news/",
			Tags: []string{"news", "current affairs"},
		},
		&marks.Mark{
			Id:   "Google",
			Url:  "https://www.google.com",
			Tags: []string{"search"},
		},
		&marks.Mark{
			Id:   "BBC News",
			Url:  "https://www.bbc.com/news",
			Tags: []string{"news", "uk"},
		},
	}
	expected := []string{
		"colorized[Abc News]    colorized[https://www.abc.net.au/news/]    colorized[[news, current affairs]]",
		"colorized[Google]      colorized[https://www.google.com]          colorized[[search]]",
		"colorized[BBC News]    colorized[https://www.bbc.com/news]        colorized[[news, uk]]",
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.Tabulate(mks)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v, received: %v", expected, actual)
	}
}

func TestTabulateWithColorizeError(t *testing.T) {
	mks := []*marks.Mark{
		&marks.Mark{
			Id:   "Abc News",
			Url:  "https://www.abc.net.au/news/",
			Tags: []string{"news", "current affairs"},
		},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(string, string) (string, error) {
		return "", errors.New("error")
	}
	actual, err := p.Tabulate(mks)
	if err == nil {
		t.Fatalf("Tabulate should return error")
	}
	if actual != nil {
		t.Fatalf("expected nil, received %v", actual)
	}
}

func TestTabulateWithEmptySlice(t *testing.T) {
	mks := []*marks.Mark{}
	p := NewTestPrinter()
	actual, err := p.Tabulate(mks)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(actual) > 0 {
		t.Fatalf("expected zero length slice, received %v", len(actual))
	}
}

func TestColorize(t *testing.T) {
	m := &marks.Mark{
		Id:   "Abc News",
		Url:  "https://www.abc.net.au/news/",
		Tags: []string{"news", "current affairs"},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.colorize(m)
	expected := &printMark{
		id:   "colorized[Abc News]",
		url:  "colorized[https://www.abc.net.au/news/]",
		tags: "colorized[[news, current affairs]]",
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected\n%v\nreceived\n%v\n", expected, actual)
	}
}

func TestColorizeWithError(t *testing.T) {
	m := &marks.Mark{
		Id:   "Abc News",
		Url:  "https://www.abc.net.au/news/",
		Tags: []string{"news", "current affairs"},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return "", errors.New("error")
	}
	actual, err := p.colorize(m)
	if err == nil {
		t.Fatalf("colorize should return error")
	}
	if actual != nil {
		t.Fatalf("colorize should return nil")
	}
}

func TestColorizeAll(t *testing.T) {
	mks := []*marks.Mark{
		&marks.Mark{
			Id:   "Abc News",
			Url:  "https://www.abc.net.au/news/",
			Tags: []string{"news", "current affairs"},
		},
		&marks.Mark{
			Id:   "Google",
			Url:  "https://www.google.com",
			Tags: []string{"search"},
		},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.colorizeAll(mks)
	expected := []*printMark{
		&printMark{
			id:   "colorized[Abc News]",
			url:  "colorized[https://www.abc.net.au/news/]",
			tags: "colorized[[news, current affairs]]",
		},
		&printMark{
			id:   "colorized[Google]",
			url:  "colorized[https://www.google.com]",
			tags: "colorized[[search]]",
		},
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected\n%v\nreceived\n%v\n", expected, actual)
	}
}

func TestColorizeAllWithError(t *testing.T) {
	mks := []*marks.Mark{
		&marks.Mark{
			Id:   "Abc News",
			Url:  "https://www.abc.net.au/news/",
			Tags: []string{"news", "current affairs"},
		},
		&marks.Mark{
			Id:   "Google",
			Url:  "https://www.google.com",
			Tags: []string{"search"},
		},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return "", errors.New("error")
	}
	actual, err := p.colorizeAll(mks)
	if err == nil {
		t.Fatalf("colorize should return error")
	}
	if actual != nil {
		t.Fatalf("colorize should return nil")
	}
}

func TestId(t *testing.T) {
	id := "Abc News"
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.Id(id)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[Abc News]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestUrl(t *testing.T) {
	url := "https://www.abc.net.au/news/"
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.Url(url)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[https://www.abc.net.au/news/]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestTags(t *testing.T) {
	tags := []string{"news", "current affairs"}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.Tags(tags)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[[news, current affairs]]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestBrowser(t *testing.T) {
	browser := "chrome"
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.Browser(browser)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[chrome]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFullMark(t *testing.T) {
	m := &marks.Mark{
		Id:   "Abc News",
		Url:  "https://www.abc.net.au/news/",
		Tags: []string{"news", "current affairs"},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.FullMark(m)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[Abc News] colorized[https://www.abc.net.au/news/] colorized[[news, current affairs]]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFullMarkNoUrl(t *testing.T) {
	m := &marks.Mark{
		Id:   "Abc News",
		Tags: []string{"news", "current affairs"},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.FullMark(m)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[Abc News] colorized[[news, current affairs]]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFullMarkNoTags(t *testing.T) {
	m := &marks.Mark{
		Id:  "Abc News",
		Url: "https://www.abc.net.au/news/",
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.FullMark(m)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "colorized[Abc News] colorized[https://www.abc.net.au/news/]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFullMarkWithFields(t *testing.T) {
	m := &marks.Mark{
		Id:   "Abc News",
		Url:  "https://www.abc.net.au/news/",
		Tags: []string{"news", "current affairs"},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.FullMarkWithFields(m)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "Id: colorized[Abc News], Url: colorized[https://www.abc.net.au/news/], Tags: colorized[[news, current affairs]]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFullMarkWithFieldsNoUrl(t *testing.T) {
	m := &marks.Mark{
		Id:   "Abc News",
		Tags: []string{"news", "current affairs"},
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.FullMarkWithFields(m)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "Id: colorized[Abc News], Tags: colorized[[news, current affairs]]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFullMarkWithFieldsNoUrlOrTags(t *testing.T) {
	m := &marks.Mark{
		Id: "Abc News",
	}
	p := NewTestPrinter()
	p.colorizer.(*mocks.Colorizer).ColorizeFn = func(colorName, text string) (string, error) {
		return fmt.Sprintf("colorized[%v]", text), nil
	}
	actual, err := p.FullMarkWithFields(m)
	if err != nil {
		t.Fatalf("should not return error")
	}
	expected := "Id: colorized[Abc News]"
	if expected != actual {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}
