package runner

import (
	"errors"
	"testing"

	"github.com/abitofoldtom/marks/marks"

	"github.com/abitofoldtom/marks/mocks"
)

func newTestAddRunner() *add {
	return &add{
		args:         &AddArgs{},
		config:       mocks.NewConfig(),
		marksService: mocks.NewMarkService(),
		printer:      mocks.NewPrinter(),
	}
}

func TestAddMark(t *testing.T) {
	a := newTestAddRunner()
	mark := mocks.DefaultMarks[2]
	msgFn := func(actualMsg string, i ...interface{}) {
		expectedMsg := "Bookmark created: %v"
		if actualMsg != expectedMsg {
			t.Fatalf("expected %v, received %v", expectedMsg, actualMsg)
		}
	}
	a.printer.(*mocks.Printer).MsgFn = msgFn
	a.args.id = mark.Id
	a.args.url = mark.Url
	a.args.tags = mark.Tags
	err := a.Run()
	if err != nil {
		t.Fatal(err.Error())
	}
	if !a.printer.(*mocks.Printer).MsgFnCalled ||
		!a.printer.(*mocks.Printer).FullMarkFnCalled ||
		!a.marksService.(*mocks.MarkService).ContainsFnCalled ||
		!a.marksService.(*mocks.MarkService).CreateFnCalled {
		t.Fatal("msg, filter, and copy should all be called")
	}
}

func TestCreateWhenMarkExists(t *testing.T) {
	a := newTestAddRunner()
	msgFn := func(actual string, i ...interface{}) {
		expected := "Bookmark with id \"%v\" already exists"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	containsFn := func(id string) (bool, error) {
		return true, nil
	}
	a.printer.(*mocks.Printer).MsgFn = msgFn
	a.marksService.(*mocks.MarkService).ContainsFn = containsFn
	err := a.Run()
	if err != nil {
		t.Fatalf("should not return error")
	}
	if a.printer.(*mocks.Printer).MsgFnCalled ||
		a.printer.(*mocks.Printer).FullMarkFnCalled ||
		a.marksService.(*mocks.MarkService).CreateFnCalled {
		t.Fatal("msg, fullMark, and create should not be called")
	}
}

func TestCreateWhenMarkExistsWithContainsError(t *testing.T) {
	a := newTestAddRunner()
	containsFn := func(id string) (bool, error) {
		return false, errors.New("error")
	}
	a.marksService.(*mocks.MarkService).ContainsFn = containsFn
	err := a.Run()
	if err == nil {
		t.Fatalf("should return error")
	}
	if a.printer.(*mocks.Printer).MsgFnCalled ||
		a.printer.(*mocks.Printer).FullMarkFnCalled ||
		a.marksService.(*mocks.MarkService).CreateFnCalled {
		t.Fatal("msg, fullMark, and create should not be called")
	}
}

func TestAddMarkWithCreateError(t *testing.T) {
	a := newTestAddRunner()
	createFn := func(*marks.Mark) error {
		return errors.New("error")
	}
	a.marksService.(*mocks.MarkService).CreateFn = createFn
	err := a.Run()
	if err == nil {
		t.Fatalf("should return error")
	}
	if a.printer.(*mocks.Printer).MsgFnCalled ||
		a.printer.(*mocks.Printer).FullMarkFnCalled {
		t.Fatal("msg and full mark should not be called")
	}
}

func TestAddMarkWithFullMarkError(t *testing.T) {
	a := newTestAddRunner()
	fullMarkFn := func(*marks.Mark) (string, error) {
		return "", errors.New("error")
	}
	a.printer.(*mocks.Printer).FullMarkFn = fullMarkFn
	err := a.Run()
	if err == nil {
		t.Fatalf("should return error")
	}
	if a.printer.(*mocks.Printer).MsgFnCalled {
		t.Fatal("msg should not be called")
	}
}
