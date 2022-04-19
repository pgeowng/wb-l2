package main

import (
	"strings"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		cfg    *Config
		input  string
		output string
		err    bool
	}{
		{
			&Config{unique: true},
			"a\na\nc",
			"a\nc\n",
			false,
		},
		{
			&Config{reverse: true},
			"a\na\nc",
			"c\na\na\n",
			false,
		},
		{
			&Config{numeric: true},
			"10\n\n1\n2",
			"\n1\n2\n10\n",
			false,
		},
		{
			&Config{},
			"10\n\n1\n2",
			"\n1\n10\n2\n",
			false,
		},
		{
			&Config{fields: []Field{{3}, {1}}},
			"john smith male 24\nkate brown female 18\njude cage female 14\nmike johnson male 14",
			"jude cage female 14\nmike johnson male 14\nkate brown female 18\njohn smith male 24\n",
			false,
		},
	}

	for idx, test := range tests {

		input := strings.NewReader(test.input)
		output := &strings.Builder{}
		sort := NewSort(test.cfg)
		err := sort.Run(input, output)
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
