package printer

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	id := "id"
	url := "url"
	tags := "tag1 tag2"
	mark := &printMark{id: id, url: url, tags: tags}
	actual := mark.split()
	expected := []string{id, url, tags}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}
