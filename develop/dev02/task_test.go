package main

import "testing"

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"", "", false},
		{"a", "a", false},
		{"ab", "ab", false},

		{"a0", "", false},
		{"a00", "", false},
		{"a0b0", "", false},

		{"abc", "abc", false},
		{"a4bc", "aaaabc", false},
		{"a12bc", "aaaaaaaaaaaabc", false},
		{"a12b5c", "aaaaaaaaaaaabbbbbc", false},

		{"\\11", "1", false},
		{"\\18", "11111111", false},

		{"\\\\3", "\\\\\\", false},

		{"\\a8", "anything", true},
		{"45", "anything", true},
		{"\\", "", false},
		{"\\\\", "\\", false},
		{"\\\\\\", "\\", false},
		{"\\\\\\\\", "\\\\", false},

		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},

		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"qwe\\\\5", "qwe\\\\\\\\\\", false},

		{"a\\a0", "", true},
		{"\\00", "", false},
		{"\\06", "000000", false},
		{"\\020", "00000000000000000000", false},
		{"f0001", "f", false},
	}

	for _, test := range tests {
		result, err := Unpack(test.input)
		if test.err {
			if err == nil {
				t.Logf("for %s expected err, got result: %s", test.input, result)
				t.Fail()
			}

		} else {
			if err != nil {
				t.Logf("for %s expected %s, got error: %s", test.input, test.expected, err)
				t.Fail()
			} else if result != test.expected {
				t.Logf("for %#v expected %#v, got result: %#v", test.input, test.expected, result)
				t.Fail()
			}
		}
	}
}
