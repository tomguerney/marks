package marks

import (
	"testing"
)

const (
	mockId  string = "mockId"
	mockUrl string = "mockUrl"
)

var mockTags = []string{"mockTag1", "mockTag2", "mockTag3"}

func newTestMark() *Mark {
	return &Mark{
		Id:   mockId,
		Url:  mockUrl,
		Tags: mockTags,
	}
}

func TestMarkDoesContainsTag(t *testing.T) {
	mockMark := newTestMark()
	if !mockMark.ContainsTag(mockTags[0]) {
		t.Fatalf("mockMark does not contain tag %v", mockTags[0])
	}
}

func TestMarkDoesNotContainsTag(t *testing.T) {
	mockMark := newTestMark()
	if mockMark.ContainsTag("not a tag") {
		t.Fatalf("mockMark contains tag \"not a tag\"")
	}
}


func TestMarkDoesContainsAllTags(t *testing.T) {
	mockMark := newTestMark()
	shouldContain := []string{mockTags[0], mockTags[1]}
	if !mockMark.ContainsAllTags(shouldContain) {
		t.Fatalf("mockMark does not contain tags %v", shouldContain)
	}
}

func TestMarkDoesNotContainsAllTags(t *testing.T) {
	mockMark := newTestMark()
	shouldNotContain := []string{mockTags[0], "not a tag"}
	if mockMark.ContainsAllTags(shouldNotContain) {
		t.Fatalf("mockMark contains tags %v", shouldNotContain)
	}
}
