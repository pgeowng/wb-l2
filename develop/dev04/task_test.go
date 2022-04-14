package main

import (
	"fmt"
	"testing"
)

func TestHashLettersRussian(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
		err      bool
	}{
		{"english", 0, true},

		{"", 0, false},
		{"а", 1, false},
		{"я", 32, false},
		{"ё", 33, false},
		{"аа", 34, false},
		{"ая", 65, false},
		{"яа", 65, false},
		{"аё", 66, false},
		{"ёа", 66, false},
		{"яя", 96, false},
		{"ёё", 98, false},
		{"ааа", 99, false},
	}

	for _, test := range tests {
		result, err := hashLettersRussian(test.input)

		if test.err {
			if err == nil {
				t.Logf("for %#v expected error, got %d", test, result)
				t.Fail()
			}
		} else {
			if err != nil {
				t.Logf("for %#v expected result, got error %s", test, err)
				t.Fail()
			} else if result != test.expected {
				t.Logf("for %#v expected result %d doesn't match %d", test, test.expected, result)
				t.Fail()
			}
		}
	}
}

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		input  []string
		result map[string][]string
	}{
		{
			[]string{"абвг", "абв"},
			map[string][]string{},
		},
		{
			[]string{"listen", "silent"},
			map[string][]string{},
		},
		{
			[]string{"лягушка", "гуляшка", "человек"},
			map[string][]string{
				"лягушка": []string{"гуляшка", "лягушка"},
			},
		},
		{
			[]string{"гуляшка", "лягушка", "человек"},
			map[string][]string{
				"гуляшка": []string{"гуляшка", "лягушка"},
			},
		},
		{
			[]string{"ГуЛяШка", "ляГуШка", "лягуШкА"},
			map[string][]string{
				"гуляшка": []string{"гуляшка", "лягушка"},
			},
		},
		{
			[]string{"сева", "веса", "сёва", "васё"},
			map[string][]string{
				"сева": []string{"веса", "сева"},
				"сёва": []string{"васё", "сёва"},
			},
		},
		{
			[]string{"столик", "пятак", "реализме", "тяпка", "пятка", "израелем", "пятак", "листок", "слово", "тяпка", "слиток"},
			map[string][]string{
				"столик":   []string{"листок", "слиток", "столик"},
				"пятак":    []string{"пятак", "пятка", "тяпка"},
				"реализме": []string{"израелем", "реализме"},
			},
		},
	}
	for _, test := range tests {
		result := groupAnagrams(test.input)

		if fmt.Sprint(result) != fmt.Sprint(test.result) {
			t.Logf("unexpected result")
			t.Logf("expected %v", test.result)
			t.Logf("got %v", result)
			t.Fail()
		}
	}
}
