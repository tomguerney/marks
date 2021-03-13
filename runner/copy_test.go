package runner

import (
	"errors"
	"testing"

	"github.com/abitofoldtom/marks/marks"
	"github.com/abitofoldtom/marks/mocks"
)

func newTestCopyRunner() *copyRunner {
	return &copyRunner{
		runner:  newTestRunner(),
		args:    &CopyArgs{},
		clipper: mocks.NewClipper(),
	}
}

func TestCopySuccess(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestCopyRunner()
	msgFn := func(actual string, i ...interface{}) {
		expected := "Url copied to clipboard: %v"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	copyFn := func(actual string) error {
		if actual != m.Url {
			t.Fatalf("expected %v, received %v", m.Url, actual)
		}
		return nil
	}
	r.printer.(*mocks.Printer).MsgFn = msgFn
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.clipper.(*mocks.Clipper).CopyFn = copyFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if !r.printer.(*mocks.Printer).MsgFnCalled ||
		!r.markService.(*mocks.MarkService).FilterFnCalled ||
		!r.clipper.(*mocks.Clipper).CopyFnCalled {
		t.Fatal("msg, filter, and copy should all be called")
	}
}

func TestCopyRunnerError(t *testing.T) {
	r := newTestCopyRunner()
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
	if r.clipper.(*mocks.Clipper).CopyFnCalled {
		t.Fatal("copy function should not be called")
	}
}

func TestCopyNonRunnerError(t *testing.T) {
	r := newTestCopyRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return nil, errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err == nil {
		t.Fatal("Run should return error")
	}
	if r.clipper.(*mocks.Clipper).CopyFnCalled {
		t.Fatal("copy function should not be called")
	}
}

func TestCopyCopyError(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestCopyRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	copyFn := func(actual string) error {
		return errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.clipper.(*mocks.Clipper).CopyFn = copyFn
	if err := r.Run(); err == nil {
		t.Fatal("Run should return error")
	}
}
