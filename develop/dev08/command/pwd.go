package command

import (
	"fmt"
	"io"
	"os"
)

type ProcessWD struct{}

func (cd *ProcessWD) Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, "pwd:", err)
	}

	fmt.Fprintln(stdout, path)

	return 0
}
