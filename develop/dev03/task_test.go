package main

import "testing"

func TestSort(t *testing.T) {
	input1 := `  11   file monkey
10  Directory text
  4 Directory rose
  1 File fire
   3 file xterm
2 directory alpha
 1.4 
10 `

	tests := []struct {
		input  string
		result string
	}{
		{"b\na\nc", "a\nb\nc\n"},
		{"b\na\nc\n", "a\nb\nc\n"},
	}
}
