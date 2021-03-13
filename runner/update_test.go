package runner

import (
	"errors"
	"reflect"
	"testing"

	"github.com/abitofoldtom/marks/marks"
	"github.com/abitofoldtom/marks/mocks"
)

func newTestUpdateRunner() *update {
	return &update{
		runner: newTestRunner(),
		args:   &UpdateArgs{},
	}
}

func TestUpdateSuccess(t *testing.T) {
	original := mocks.DefaultMarks[0]
	r := newTestUpdateRunner()
	msgFn := func(actual string, i ...interface{}) {
		expected := "Mark updated from:\n%v\nto:\n%v"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{original}, nil
	}
	updateFn := func(id string, actual *marks.Mark) error {
		expected := &marks.Mark{
			Id:   "Abc Updated News",
			Url:  "https://www.updated.com",
			Tags: []string{"current affairs", "newTag1", "newTag2"},
		}
		if original.Id != id {
			t.Fatalf("expected %v, received %v", original.Id, id)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
		return nil
	}
	r.args = &UpdateArgs{
		newId:      "Abc Updated News",
		newUrl:     "https://www.updated.com",
		newTags:    []string{"newTag1", "newTag2"},
		removeTags: []string{"news"},
	}
	r.printer.(*mocks.Printer).MsgFn = msgFn
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.markService.(*mocks.MarkService).UpdateFn = updateFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if !r.printer.(*mocks.Printer).MsgFnCalled ||
		!r.markService.(*mocks.MarkService).UpdateFnCalled {
		t.Fatal("msg update should be called")
	}
}

func TestUpdateSuccessRemoveUrl(t *testing.T) {
	original := mocks.DefaultMarks[0]
	r := newTestUpdateRunner()
	msgFn := func(actual string, i ...interface{}) {
		expected := "Mark updated from:\n%v\nto:\n%v"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{original}, nil
	}
	updateFn := func(id string, actual *marks.Mark) error {
		expected := &marks.Mark{
			Id:   "Abc News",
			Url:  "",
			Tags: []string{"current affairs"},
		}
		if original.Id != id {
			t.Fatalf("expected %v, received %v", original.Id, id)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
		return nil
	}
	r.args = &UpdateArgs{
		removeTags: []string{"news"},
		removeUrl:  true,
	}
	r.printer.(*mocks.Printer).MsgFn = msgFn
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.markService.(*mocks.MarkService).UpdateFn = updateFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if !r.printer.(*mocks.Printer).MsgFnCalled ||
		!r.markService.(*mocks.MarkService).UpdateFnCalled {
		t.Fatal("msg and update should be called")
	}
}

func TestUpdateRunnerError(t *testing.T) {
	r := newTestUpdateRunner()
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
	if r.markService.(*mocks.MarkService).UpdateFnCalled {
		t.Fatal("update function should not be called")
	}
}

func TestUpdateNonRunnerError(t *testing.T) {
	r := newTestUpdateRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return nil, errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err == nil {
		t.Fatal("Run should return error")
	}
	if r.markService.(*mocks.MarkService).UpdateFnCalled {
		t.Fatal("update function should not be called")
	}
}

func TestUpdateUpdateError(t *testing.T) {
	m := mocks.DefaultMarks[0]
	r := newTestUpdateRunner()
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{m}, nil
	}
	updateFn := func(id string, m *marks.Mark) error {
		return errors.New("error")
	}
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	r.markService.(*mocks.MarkService).UpdateFn = updateFn
	if err := r.Run(); err == nil {
		t.Fatal("Run should return error")
	}
}

func TestUpdateRemoveTagsNotContained(t *testing.T) {
	original := mocks.DefaultMarks[0]
	r := newTestUpdateRunner()
	errorFn := func(actual string, i ...interface{}) {
		expected := "Mark %v does not contain tag %v"
		if actual != expected {
			t.Fatalf("expected %v, received %v", expected, actual)
		}
	}
	filterFn := func(string, string, []string) ([]*marks.Mark, error) {
		return []*marks.Mark{original}, nil
	}
	r.args = &UpdateArgs{
		removeTags: []string{"not a tag"},
	}
	r.printer.(*mocks.Printer).ErrorFn = errorFn
	r.markService.(*mocks.MarkService).FilterFn = filterFn
	if err := r.Run(); err != nil {
		t.Fatal(err.Error())
	}
	if !r.printer.(*mocks.Printer).ErrorFnCalled {
		t.Fatal("error should be called")
	}
	if r.markService.(*mocks.MarkService).UpdateFnCalled {
		t.Fatal("update should not be called")
	}
}

func TestUpdatedIdNoNewId(t *testing.T) {
	expected := "original id"
	selected := &marks.Mark{
		Id: expected,
	}
	u := newTestUpdateRunner()
	u.args.newId = ""
	actual := u.updatedId(selected)
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestUpdatedIdNewId(t *testing.T) {
	expected := "new id"
	selected := &marks.Mark{
		Id: "original id",
	}
	u := newTestUpdateRunner()
	u.args.newId = expected
	actual := u.updatedId(selected)
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestUpdatedUrlNoNewUrl(t *testing.T) {
	expected := "original url"
	selected := &marks.Mark{
		Url: expected,
	}
	u := newTestUpdateRunner()
	u.args.newUrl = ""
	actual := u.updatedUrl(selected)
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestUpdatedUrlNewUrl(t *testing.T) {
	expected := "new url"
	selected := &marks.Mark{
		Url: "original url",
	}
	u := newTestUpdateRunner()
	u.args.newUrl = expected
	actual := u.updatedUrl(selected)
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestUpdatedTagsNoNewTags(t *testing.T) {
	expected := []string{"oldTag1", "oldTag2"}
	selected := &marks.Mark{
		Tags: expected,
	}
	u := newTestUpdateRunner()
	u.args.newTags = []string{}
	actual := u.updatedTags(selected)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestUpdatedTagsNewTags(t *testing.T) {
	expected := []string{"oldTag1", "oldTag2", "newTag1", "newTag2"}
	selected := &marks.Mark{
		Tags: []string{expected[0], expected[1]},
	}
	u := newTestUpdateRunner()
	u.args.newTags = []string{expected[2], expected[3]}
	actual := u.updatedTags(selected)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestDoesContainsRemoveTag(t *testing.T) {
	tags := []string{"tag1", "tag2", "tag3"}
	mark := &marks.Mark{
		Tags: tags,
	}
	u := newTestUpdateRunner()
	u.args.removeTags = []string{tags[1], tags[2]}
	tag, contained := u.containsRemoveTags(mark)
	if tag != "" {
		t.Fatalf("expected empty string, received %v", tag)
	}
	if !contained {
		t.Fatalf("expected true, received false")
	}
}

func TestDoesNotContainRemoveTag(t *testing.T) {
	tags := []string{"tag1", "tag2", "tag3"}
	expected := "different tag"
	mark := &marks.Mark{
		Tags: tags,
	}
	u := newTestUpdateRunner()
	u.args.removeTags = []string{tags[1], expected}
	actual, contained := u.containsRemoveTags(mark)
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
	if contained {
		t.Fatalf("expected false, received true")
	}
}

func TestIsRemoveTag(t *testing.T) {
	removeTags := []string{"removeTag1", "removeTag2"}
	u := newTestUpdateRunner()
	u.args.removeTags = removeTags
	if !u.removeTag(removeTags[1]) {
		t.Fatalf("expected true, received false")
	}
}

func TestIsNotRemoveTag(t *testing.T) {
	removeTags := []string{"removeTag1", "removeTag2"}
	u := newTestUpdateRunner()
	u.args.removeTags = removeTags
	if u.removeTag("not a remove tag") {
		t.Fatalf("expected false, received true")
	}
}

func TestRemoveTagsNoRemoveTags(t *testing.T) {
	expected := []string{"tag1", "tag2"}
	u := newTestUpdateRunner()
	u.args.removeTags = []string{}
	actual := u.removeTags(expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestRemoveTagsHasRemoveTags(t *testing.T) {
	tags := []string{"tag1", "tag2", "tag3", "tag4"}
	u := newTestUpdateRunner()
	u.args.removeTags = []string{tags[0], tags[2]}
	expected := []string{tags[1], tags[3]}
	actual := u.removeTags(tags)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestRemoveTagsHasRemoveTagsWithNotContainedTag(t *testing.T) {
	tags := []string{"tag1", "tag2", "tag3", "tag4"}
	u := newTestUpdateRunner()
	u.args.removeTags = []string{tags[0], tags[2], "not a tag"}
	expected := []string{tags[1], tags[3]}
	actual := u.removeTags(tags)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}
