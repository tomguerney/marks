package arg

import "testing"

func TestPopZeroArgs(t *testing.T) {
	p := NewParser([]string{})
	_, err := p.Pop()
	if err == nil {
		t.Fatalf("pop on zero args should return error")
	}
	if len(p.Remaining()) != 0 {
		t.Fatalf("should be zero remaining args")
	}
}

func TestPopOneArg(t *testing.T) {
	first := "first"
	p := NewParser([]string{first})
	arg, err := p.Pop()
	if err != nil {
		t.Fatal("pop on one arg should not return error")
	}
	if arg != first {
		t.Fatalf("expected %v, received %v", first, arg)
	}
	if len(p.Remaining()) != 0 {
		t.Fatalf("should be zero remaining args")
	}
}

func TestPopOneArgWithOneRemaining(t *testing.T) {
	first := "first"
	second := "second"
	p := NewParser([]string{first, second})
	arg, err := p.Pop()
	if err != nil {
		t.Fatal("pop on one arg should not return error")
	}
	if arg != first {
		t.Fatalf("expected %v, received %v", first, arg)
	}
	remaining := p.Remaining()
	if len(remaining) != 1 {
		t.Fatalf("should be one remaining args")
	}
	if remaining[0] != second {
		t.Fatalf("expected %v, received %v", second, remaining[0])
	}
}
