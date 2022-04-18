package command

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestCd(t *testing.T) {
	tests := []struct {
		args         []string
		vars         map[string]string
		expectedVars map[string]string
		exitCode     int
	}{
		{
			args:         []string{},
			vars:         map[string]string{},
			exitCode:     1,
			expectedVars: map[string]string{},
		},
		{
			args:         []string{"home"},
			vars:         map[string]string{},
			exitCode:     1,
			expectedVars: map[string]string{},
		},
		{
			args:         []string{"/home"},
			vars:         map[string]string{},
			exitCode:     0,
			expectedVars: map[string]string{"PWD": "/home"},
		},
		{
			args:         []string{"/home"},
			vars:         map[string]string{"PWD": "/anything"},
			exitCode:     0,
			expectedVars: map[string]string{"PWD": "/home"},
		},
		{
			args:         []string{".././../bin"},
			vars:         map[string]string{"PWD": "/dev/disk"},
			exitCode:     0,
			expectedVars: map[string]string{"PWD": "/bin"},
		},

		// environment dependent tests
		{
			args:         []string{"disk"},
			vars:         map[string]string{"PWD": "/dev"},
			exitCode:     0,
			expectedVars: map[string]string{"PWD": "/dev/disk"},
		},
	}

	for idx, test := range tests {
		cd := &ChangeDir{}
		exitCode := cd.Run(test.args, test.vars, strings.NewReader(""), io.Discard, io.Discard)
		if exitCode != 0 {
			if test.exitCode != exitCode {
				t.Logf("for %v expected exit(%d), got %d", idx, test.exitCode, exitCode)
				t.Fail()
			}
		} else {
			if test.exitCode != exitCode {
				t.Logf("for %v expected exit(%d), got %d", idx, test.exitCode, exitCode)
				t.Fail()
			} else if fmt.Sprint(test.vars) != fmt.Sprint(test.expectedVars) {
				t.Logf("for %v vars mismatch", idx)
				t.Logf("expected: %s", test.expectedVars)
				t.Logf("got: %s", test.vars)
				t.Fail()
			}
		}
	}
}
