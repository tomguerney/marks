package clipper

import (
	"testing"

	"github.com/atotto/clipboard"
)

func TestClipper(t *testing.T) {
	expected := "copy test"
	clipper := NewClipper()
	clipper.Copy("copy test")
	actual, err := clipboard.ReadAll()
	if err != nil {
		t.Fatal("error returned from clipboard ReadAll")
	}
	if actual != expected {
		t.Fatalf("expected %v, received %v", actual, expected)
	}
}
