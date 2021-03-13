package runner

import (
	"errors"
	"testing"

	"github.com/abitofoldtom/marks/marks"
	"github.com/abitofoldtom/marks/mocks"
)

func newTestOpenRunner() *open {
	return &open{
		runner: newTestRunner(),
		args:   &OpenArgs{},
		opener: mocks.NewOpener(),
	}
}

func TestOpenSuccess(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestOpenRunner()
	expectedBrowser := "chrome"
	r.config.Browser = expectedBrowser
	msgFn := func(actual string, i ...interface{}) {
		expected := "Url opened in %v: %v"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	openFn := func(actualUrl, actualBrowser string) error {
		if actualUrl != m.Url {
			t.Fatalf("expected %v, received %v", m.Url, actualUrl)
		}
		if actualBrowser != expectedBrowser {
			t.Fatalf("expected %v, received %v", expectedBrowser, actualBrowser)
		}
		return nil
	}
	r.printer.(*mocks.Printer).MsgFn = msgFn
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.opener.(*mocks.Opener).OpenFn = openFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if !r.printer.(*mocks.Printer).MsgFnCalled ||
		!r.markService.(*mocks.MarkService).FilterFnCalled ||
		!r.opener.(*mocks.Opener).OpenFnCalled {
		t.Fatal("msg, filter, and open should all be called")
	}
}

func TestOpenRunnerError(t *testing.T) {
	r := newTestOpenRunner()
	msgFn := func(actual string, i ...interface{}) {
		expected := "No bookmarks found"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{}, nil
	}
	r.printer.(*mocks.Printer).MsgFn = msgFn
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if r.opener.(*mocks.Opener).OpenFnCalled {
		t.Fatal("open functions should not be called")
	}
}

func TestOpenNonRunnerError(t *testing.T) {
	r := newTestOpenRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return nil, errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err == nil {
		t.Fatal("should return error")
	}
	if r.opener.(*mocks.Opener).OpenFnCalled {
		t.Fatal("open function should not be called")
	}
}

func TestOpenOpenError(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestOpenRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	openFn := func(string, string) error {
		return errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.opener.(*mocks.Opener).OpenFn = openFn
	if err := r.Run(); err == nil {
		t.Fatal("should return error")
	}
}
