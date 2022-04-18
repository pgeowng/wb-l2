package main

import (
	"strings"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		input string
		left  int
		right int
		err   bool
	}{
		{"-", 1, 1, true},
		{"0", 1, 1, true},
		{"1", 1, 1, false},
		{"1-", 1, -1, false},
		{"1-0", 1, -1, true},
		{"1-2", 1, 2, false},
		{"-2", -1, 2, false},
		{"-0", -1, 2, true},
		{"1--1", 0, 0, true},
	}

	for _, test := range tests {
		left, right, err := ParseFields(test.input)
		if test.err {
			if err == nil {
				t.Logf("for %v expected error, got result: %v %v", test.input, left, right)
				t.Fail()
			}
		} else {
			if err != nil {
				t.Logf("for %v expected result, got error: %v", test.input, err)
				t.Fail()
			} else if left != test.left || right != test.right {
				t.Logf("for %v expected result mismatch", test.input)
				t.Logf("expected: l:%v r:%v", test.left, test.right)
				t.Logf("got: l:%v r:%v", left, right)
				t.Fail()
			}

		}
	}

}

func TestCut(t *testing.T) {
	tests := []struct {
		cfg    *Config
		input  string
		output string
		err    bool
	}{
		{
			&Config{leftField: 2, rightField: 2, delimiter: " ", onlyDelimited: false},
			"a b\na\nc c",
			"b\na\nc\n",
			false,
		}, {
			&Config{leftField: 2, rightField: 2, delimiter: " ", onlyDelimited: false},
			"a b\na \nc c",
			"b\n\nc\n",
			false,
		},
		{
			&Config{leftField: 2, rightField: -1, delimiter: " ", onlyDelimited: false},
			"a b\na\nc g t",
			"b\na\ng t\n",
			false,
		},
		{
			&Config{leftField: -1, rightField: 2, delimiter: " ", onlyDelimited: false},
			"a b\na\nc g t",
			"a b\na\nc g\n",
			false,
		},
		{
			&Config{leftField: -1, rightField: 2, delimiter: " ", onlyDelimited: true},
			"a b\na\nc g t",
			"a b\nc g\n",
			false,
		},
		{
			&Config{leftField: 3, rightField: 9, delimiter: "t", onlyDelimited: true},
			"go.mod\ntask.go\ntask_test.go",
			"\nest.go\n",
			false,
		},
	}

	for idx, test := range tests {
		input := strings.NewReader(test.input)
		output := &strings.Builder{}
		prog := NewCut(test.cfg)
		err := prog.Run(input, output)
		if test.err {
			if err == nil {
				t.Logf("for %#v test expected err, got result", idx)
				t.Fail()
			}
		} else {
			if err != nil {
				t.Logf("for %#v test expected result, got err", idx)
				t.Fail()
			} else if test.output != output.String() {
				t.Logf("for %#v test expected result mismatch\n", idx)
				t.Logf("expected: %#v\n", test.output)
				t.Logf("got: %#v\n", output.String())
				t.Fail()
			}
		}
	}
}
