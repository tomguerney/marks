package yaml

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/abitofoldtom/marks/marks"
	"github.com/abitofoldtom/marks/mocks"
	"gopkg.in/yaml.v2"
)

type mockReaderWriter struct {
	ReadFileFn  func(string) ([]byte, error)
	WriteFileFn func(string, []byte, uint32) error
}

func (rw mockReaderWriter) ReadFile(s string) ([]byte, error) {
	return rw.ReadFileFn(s)
}

func (rw mockReaderWriter) WriteFile(s string, b []byte, u uint32) error {
	return rw.WriteFileFn(s, b, u)
}

func mockReadFile(string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join("testdata", "defaultmarks.yaml"))
}

func mockWriteFile(string, []byte, uint32) error {
	return nil
}

func newMockReaderWriter() *mockReaderWriter {
	return &mockReaderWriter{
		mockReadFile,
		mockWriteFile,
	}
}

func newTestMarkService() *markService {
	return &markService{
		mocks.NewConfig(),
		newMockReaderWriter(),
	}
}

func TestRetrieveMarkSuccess(t *testing.T) {
	s := newTestMarkService()
	actual, err := s.Mark("abc nEws")
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := &marks.Mark{
		Id:   "Abc News",
		Url:  "https://www.abc.net.au/news/",
		Tags: []string{"news", "current affairs"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestRetrieveNoMark(t *testing.T) {
	s := newTestMarkService()
	actual, err := s.Mark("not a mark")
	if err != nil {
		t.Fatal(err.Error())
	}
	if actual != nil {
		t.Errorf("expected nil Mark, received %v", actual)
	}
}

func TestRetrieveMarkFileReadFail(t *testing.T) {
	errorFunc := func(string) ([]byte, error) { return nil, errors.New("read error") }
	s := newTestMarkService()
	s.readerWriter.(*mockReaderWriter).ReadFileFn = errorFunc
	actual, err := s.Mark("expecting read error")
	if err == nil {
		t.Errorf("expected error, received %T", err)
	}
	if actual != nil {
		t.Errorf("expected nil Mark, received %v", actual)
	}
}

func TestCreateMark(t *testing.T) {
	new := &marks.Mark{Id: "Github", Url: "https://github.com/", Tags: []string{"code", "repository"}}
	writeFunc := func(s string, bytes []byte, u uint32) error {
		marks := []*marks.Mark{}
		if err := yaml.Unmarshal(bytes, &marks); err != nil {
			t.Fatal(err.Error())
		}
		if len(marks) != 6 {
			t.Errorf("expected 6 marks, received %v", len(marks))
		}
		if !reflect.DeepEqual(new, marks[5]) {
			t.Fatalf("expected:\n%v\nreceived:\n%v", new, marks[5])
		}
		return nil
	}
	s := newTestMarkService()
	s.readerWriter.(*mockReaderWriter).WriteFileFn = writeFunc
	err := s.Create(new)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateExistingMark(t *testing.T) {
	s := newTestMarkService()
	m := &marks.Mark{Id: "Google"}
	err := s.Create(m)
	if _, ok := err.(marks.MarkAlreadyExistsError); !ok {
		t.Fatalf("expected MarkAlreadyExistsError, received %T", err)
	}
}

func TestUpdateMark(t *testing.T) {
	oldId := "Google"
	new := &marks.Mark{Id: "Goooogle", Url: "https://www.different.com", Tags: []string{"different1", "different2"}}
	expected := &marks.Mark{
		Id:   "Goooogle",
		Url:  "https://www.different.com",
		Tags: []string{"different1", "different2"},
	}
	writeFunc := func(s string, bytes []byte, u uint32) error {
		marks := []*marks.Mark{}
		if err := yaml.Unmarshal(bytes, &marks); err != nil {
			t.Fatal(err.Error())
		}
		updatedFound := false
		for _, actual := range marks {
			if actual.Id == new.Id {
				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("expected:\n%v\nreceived:\n%v", expected, actual)
				}
				updatedFound = true
			}
			if actual.Id == oldId {
				t.Fatal("old id still present in marks")
			}
		}
		if !updatedFound {
			t.Fatal("updated mark not found")
		}
		return nil
	}
	s := newTestMarkService()
	s.readerWriter.(*mockReaderWriter).WriteFileFn = writeFunc
	err := s.Update(oldId, new)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUpdateNonExistentMark(t *testing.T) {
	updated := &marks.Mark{Id: "Goooogle", Url: "https://www.different.com", Tags: []string{"different"}}
	s := newTestMarkService()
	err := s.Update("not a mark", updated)
	if _, ok := err.(marks.MarkDoesNotExistError); !ok {
		t.Fatal(err.Error())
	}
}

func TestDeleteMark(t *testing.T) {
	deletedId := "Google"
	writeFunc := func(s string, bytes []byte, u uint32) error {
		marks := []*marks.Mark{}
		if err := yaml.Unmarshal(bytes, &marks); err != nil {
			t.Fatal(err.Error())
		}
		if len(marks) != 4 {
			t.Errorf("expected 4 marks, received %v", len(marks))
		}
		for _, actual := range marks {
			if actual.Id == deletedId {
				t.Fatal("deleted id still present in marks")
			}
		}
		return nil
	}
	s := newTestMarkService()
	s.readerWriter.(*mockReaderWriter).WriteFileFn = writeFunc
	err := s.Delete(deletedId)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteNonExistentMark(t *testing.T) {
	s := newTestMarkService()
	err := s.Delete("not a mark")
	if _, ok := err.(marks.MarkDoesNotExistError); !ok {
		t.Fatal(err.Error())
	}
}

func TestDoesContainsMark(t *testing.T) {
	s := newTestMarkService()
	exists, err := s.Contains("goOgle")
	if err != nil {
		t.Fatal(err.Error())
	}
	if !exists {
		t.Errorf("expected true, received false")
	}
}

func TestDoesNotContainsMark(t *testing.T) {
	s := newTestMarkService()
	exists, err := s.Contains("not a mark")
	if err != nil {
		t.Fatal(err.Error())
	}
	if exists {
		t.Errorf("expected false, received true")
	}
}

func TestFilterZeroMarksById(t *testing.T) {
	s := newTestMarkService()
	id := "not a mark"
	result, err := s.Filter(id, "", []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkById(t *testing.T) {
	s := newTestMarkService()
	id := "google"
	result, err := s.Filter(id, "", []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterTwoMarksById(t *testing.T) {
	s := newTestMarkService()
	id := "nEws"
	result, err := s.Filter(id, "", []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 2 {
		t.Errorf("expected 2 marks, received %v", len(result))
	}
}

func TestFilterZeroMarksByUrl(t *testing.T) {
	s := newTestMarkService()
	url := "not a url"
	result, err := s.Filter("", url, []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByUrl(t *testing.T) {
	s := newTestMarkService()
	url := "abc.net.au"
	result, err := s.Filter("", url, []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterThreeMarksByUrl(t *testing.T) {
	s := newTestMarkService()
	url := "nEws"
	result, err := s.Filter("", url, []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 3 {
		t.Errorf("expected 3 marks, received %v", len(result))
	}
}

func TestFilterZeroMarksByOneTag(t *testing.T) {
	s := newTestMarkService()
	tags := []string{"not a tag"}
	result, err := s.Filter("", "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByOneTag(t *testing.T) {
	s := newTestMarkService()
	tags := []string{"tutorial"}
	result, err := s.Filter("", "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterTwoMarksByOneTag(t *testing.T) {
	s := newTestMarkService()
	tags := []string{"electronics"}
	result, err := s.Filter("", "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 2 {
		t.Errorf("expected 2 marks, received %v", len(result))
	}
}

func TestFilterZeroMarksByTwoTags(t *testing.T) {
	s := newTestMarkService()
	tags := []string{"electronics", "not a tag"}
	result, err := s.Filter("", "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByTwoTags(t *testing.T) {
	s := newTestMarkService()
	tags := []string{"news", "electronics"}
	result, err := s.Filter("", "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterZeroMarksByIdAndUrl(t *testing.T) {
	s := newTestMarkService()
	id := "not a mark"
	url := "news"
	result, err := s.Filter(id, url, []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByIdAndUrl(t *testing.T) {
	s := newTestMarkService()
	id := "electronics"
	url := "news"
	result, err := s.Filter(id, url, []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterZeroMarksByIdAndTag(t *testing.T) {
	s := newTestMarkService()
	id := "news"
	tags := []string{"not a tag"}
	result, err := s.Filter(id, "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByIdAndTag(t *testing.T) {
	s := newTestMarkService()
	id := "news"
	tags := []string{"uk"}
	result, err := s.Filter(id, "", tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterZeroMarksByUrlAndTag(t *testing.T) {
	s := newTestMarkService()
	url := "news"
	tags := []string{"not a tag"}
	result, err := s.Filter("", url, tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByUrlAndTag(t *testing.T) {
	s := newTestMarkService()
	url := "news"
	tags := []string{"electronics"}
	result, err := s.Filter("", url, tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterZeroMarksByIdAndUrlAndTag(t *testing.T) {
	s := newTestMarkService()
	id := "news"
	url := "not a url"
	tags := []string{"uk"}
	result, err := s.Filter(id, url, tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 0 {
		t.Errorf("expected 0 marks, received %v", len(result))
	}
}

func TestFilterOneMarkByIdAndUrlAndTag(t *testing.T) {
	s := newTestMarkService()
	id := "news"
	url := "bbc"
	tags := []string{"uk"}
	result, err := s.Filter(id, url, tags)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("expected 1 mark, received %v", len(result))
	}
}

func TestFilterAllMarksByNoInput(t *testing.T) {
	s := newTestMarkService()
	result, err := s.Filter("", "", []string{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(result) != 5 {
		t.Errorf("expected 5 marks, received %v", len(result))
	}
}
