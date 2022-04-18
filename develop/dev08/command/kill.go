package command

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type Kill struct{}

func (cd *Kill) Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	if len(args) == 0 {
		fmt.Fprintln(stderr, "kill: not enough arguments")
		return 1
	}

	pid, err := strconv.ParseInt(args[0], 0, 0)
	if err != nil {
		fmt.Fprintln(stderr, "kill:", err)
		return 1
	}

	process, err := os.FindProcess(int(pid))
	if err != nil {
		fmt.Fprintln(stderr, "kill:", err)
		return 1
	}

	err = process.Kill()
	if err != nil {
		fmt.Fprintln(stderr, "kill:", err)
		return 1
	}

	return 0
}
