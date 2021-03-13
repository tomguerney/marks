package runner

import (
	"reflect"
	"testing"

	"github.com/abitofoldtom/marks/marks"
	"github.com/abitofoldtom/marks/mocks"
)

func newTestRunner() *runner {
	return &runner{
		config:      mocks.NewConfig(),
		markService: mocks.NewMarkService(),
		printer:     mocks.NewPrinter(),
		prompter:    mocks.NewPrompter(),
	}
}

type mockRunner struct {
	filterFn       func(prompt, id, url string, tags []string) (*marks.Mark, error)
	filterFnCalled bool
}

func (r *mockRunner) filter(prompt, id, url string, tags []string) (*marks.Mark, error) {
	r.filterFnCalled = true
	return r.filterFn(prompt, id, url, tags)
}

var defaultFilterFn = func(prompt, id, url string, tags []string) (*marks.Mark, error) {
	return mocks.DefaultMarks[0], nil
}

func newMockRunner() *mockRunner {
	return &mockRunner{
		filterFn: defaultFilterFn,
	}
}

func TestFilterZeroMarks(t *testing.T) {
	r := newTestRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{}, nil
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	mark, err := r.filter("prompt", "id", "url", []string{"tag"})
	if mark != nil {
		t.Fatal("expected nil mark")
	}
	if _, ok := err.(*runnerError); !ok {
		t.Fatal("expected runnerError")
	}
}

func TestFilterOneMark(t *testing.T) {
	expected := mocks.DefaultMarks[0]
	r := newTestRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{expected}, nil
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	actual, err := r.filter("prompt", "id", "url", []string{"tag"})
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestFilterTwoMarks(t *testing.T) {
	m1 := mocks.DefaultMarks[0]
	expected := mocks.DefaultMarks[1]
	r := newTestRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m1, expected}, nil
	}
	selectFn := func(s string, table []string) (int, error) {
		return 1, nil
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.prompter.(*mocks.Prompter).SelectFn = selectFn
	actual, err := r.filter("prompt", "id", "url", []string{"tag"})
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}
