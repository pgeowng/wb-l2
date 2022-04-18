package main

import (
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		cfg        *Config
		trueInput  []string
		falseInput []string
	}{
		{
			&Config{ignoreCase: false, fixed: false, expr: "a"},
			strings.Split("thank,another,manager", ","),
			strings.Split("therefore,keep,After", ","),
		},
		{
			&Config{ignoreCase: true, fixed: false, expr: "a"},
			strings.Split("thAnk,another", ","),
			strings.Split("there", ","),
		},
		{
			&Config{ignoreCase: false, fixed: true, expr: "her"},
			strings.Split("her", ","),
			strings.Split("there", ","),
		},
		{
			&Config{ignoreCase: true, fixed: true, expr: "her"},
			strings.Split("hER,her", ","),
			strings.Split("there,There", ","),
		},
	}

	for idx, test := range tests {
		grep, err := NewGrep(test.cfg)
		if err != nil {
			t.Logf("for %#v test expected grep, got err: %v", idx, err)
			t.Fail()
		}

		for _, line := range test.trueInput {
			if grep.Match(line) != true {
				t.Logf("for %#v test line %#v expected true, got false", idx, line)
				t.Fail()
			}
		}

		for _, line := range test.falseInput {
			if grep.Match(line) == true {
				t.Logf("for %#v test line %#v expected false, got true", idx, line)
				t.Fail()
			}
		}
	}
}

func TestProgram(t *testing.T) {

	tests := []struct {
		cfg    *Config
		input  string
		output string
		err    bool
	}{
		{
			&Config{
				after:       0,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "",
			},
			"a\nb\na\n",
			"a\nb\na\n",
			false,
		},
		{
			&Config{
				after:       0,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\na\n",
			"a\na\n",
			false,
		},
		{
			&Config{
				after:       1,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\na\nc\n",
			"a\nb\na\nc\n",
			false,
		},
		{
			&Config{
				after:       1,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\nc\na\nc\n",
			"a\nb\n--\na\nc\n",
			false,
		},
		{
			&Config{
				after:       1,
				before:      1,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\nc\na\nc\n",
			"a\nb\nc\na\nc\n",
			false,
		}, {
			&Config{
				after:       1,
				before:      1,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\nc\nd\na\nc\n",
			"a\nb\n--\nd\na\nc\n",
			false,
		}, {
			&Config{
				after:       1,
				before:      1,
				count:       true,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\nc\nd\na\nc\n",
			"2\n",
			false,
		}, {
			&Config{
				after:       1,
				before:      1,
				count:       true,
				ignoreCase:  false,
				invert:      true,
				fixed:       false,
				lineNumbers: false,
				expr:        "a",
			},
			"a\nb\nc\nd\na\nc\n",
			"4\n",
			false,
		}, {
			&Config{
				after:       0,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: false,
				expr:        "ca",
			},
			"aca\nca\ncaa\n",
			"aca\nca\ncaa\n",
			false,
		}, {
			&Config{
				after:       0,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       true,
				lineNumbers: false,
				expr:        "ca",
			},
			"aca\nca\ncaa\n",
			"ca\n",
			false,
		},
		{
			&Config{
				after:       0,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       true,
				lineNumbers: false,
				expr:        "ca",
			},
			"aca\nca\ncaa\n",
			"ca\n",
			false,
		}, {
			&Config{
				after:       1,
				before:      0,
				count:       false,
				ignoreCase:  false,
				invert:      false,
				fixed:       false,
				lineNumbers: true,
				expr:        "a",
			},
			"a\nb\nc\na\nc\n",
			"1:a\n2:b\n--\n4:a\n5:c\n",
			false,
		},
	}

	for idx, test := range tests {
		grep, err := NewGrep(test.cfg)
		if err != nil {
			t.Logf("for %#v test expected grep, got err: %v", idx, err)
			t.Fail()
		}

		input := strings.NewReader(test.input)
		output := &strings.Builder{}
		prog := NewProgram(test.cfg, grep)
		err = prog.Run(input, output)
		if test.err {
			if err == nil {
				t.Logf("for %#v test expected err, got result", idx)
				t.Fail()
			}
		} else {
			if err != nil {
				t.Logf("for %#v test expected result, got err", idx)
				t.Fail()
			}

			if test.output != output.String() {
				t.Logf("for %#v test expected result mismatch\n", idx)
				t.Logf("expected: %#v\n", test.output)
				t.Logf("got: %#v\n", output.String())
				t.Fail()
			}
		}
	}
}
