package command

import (
	"fmt"
	"io"
)

type Echo struct{}

func (cd *Echo) Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	for idx, arg := range args {
		if idx != 0 {
			fmt.Fprint(stdout, " ")
		}
		fmt.Fprint(stdout, arg)
	}

	return 0
}
