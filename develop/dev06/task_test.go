package main

import (
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	tests := []struct{
		input string
		result string
		d []byte
		fromField int
		toField int
		silent bool
	}{
		{
			"go.mod\ntask.go\ntask_test.go",
			"\n\n",
			[]byte("t"),
			1, 1, true,
		},
		{
			"go.mod\ntask.go\ntask_test.go",
			"ask.go\nask_\n",
			[]byte("t"),
			2, 2, true,
		},
		{
			"go.mod\ntask.go\ntask_test.go",
			"\nes\n",
			[]byte("t"),
			3, 3, true,
		},
		{
			"go.mod\ntask.go\ntask_test.go",
			"\n.go\n",
			[]byte("t"),
			4, 4, true,
		},
		{
			"go.mod\ntask.go\ntask_test.go",
			"\n.go\n",
			[]byte("t"),
			5, 5, true,
		},
	}
	strings.NewReader("hello world")



	    s := newContainerStats() // Replace this the appropriate constructor
    var b bytes.Buffer
    if err := s.Process(&b); err != nil {
        t.Fatalf("s.Display() gave error: %s", err)
    }
    got := b.String()
    want := "hello world\n"
    if got != want {
        t.Errorf("s.Display() = %q, want %q", got, want)
    }
}
}
