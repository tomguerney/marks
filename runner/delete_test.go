package runner

import (
	"errors"
	"testing"

	"github.com/abitofoldtom/marks/marks"
	"github.com/abitofoldtom/marks/mocks"
)

func newTestDeleteRunner() *deleteRunner {
	return &deleteRunner{
		runner: newTestRunner(),
		args:   &DeleteArgs{},
	}
}

func TestDeleteSuccess(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestDeleteRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	deleteFn := func(actual string) error {
		if actual != m.Id {
			t.Fatalf("expected %v, received %v", m.Id, actual)
		}
		return nil
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.markService.(*mocks.MarkService).DeleteFn = deleteFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if !r.printer.(*mocks.Printer).MsgFnCalled ||
		!r.markService.(*mocks.MarkService).FilterFnCalled ||
		!r.markService.(*mocks.MarkService).DeleteFnCalled {
		t.Fatal("msg, filter, and delete should all be called")
	}
}

func TestDeleteConfirmationDeclined(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestDeleteRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	confirmFn := func(string) bool {
		return false
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.prompter.(*mocks.Prompter).ConfirmFn = confirmFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if r.markService.(*mocks.MarkService).DeleteFnCalled {
		t.Fatal("delete should not be called")
	}
}

func TestDeleteRunnerError(t *testing.T) {
	r := newTestDeleteRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{}, nil
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if r.markService.(*mocks.MarkService).DeleteFnCalled {
		t.Fatal("delete function should not be called")
	}
}

func TestDeleteNonRunnerError(t *testing.T) {
	r := newTestDeleteRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return nil, errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err == nil {
		t.Fatal("Run should return error")
	}
	if r.markService.(*mocks.MarkService).DeleteFnCalled {
		t.Fatal("delete function should not be called")
	}
}

func TestDeleteDeleteError(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestDeleteRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	deleteFn := func(actual string) error {
		return errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.markService.(*mocks.MarkService).DeleteFn = deleteFn
	if err := r.Run(); err == nil {
		t.Fatal("Run should return error")
	}
}
