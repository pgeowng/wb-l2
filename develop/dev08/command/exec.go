package command

import (
	"fmt"
	"io"
	"os/exec"
)

type Exec struct{}

func (cd *Exec) Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	if len(args) == 0 {
		fmt.Fprintln(stderr, "exec: empty call")
		return 1
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintln(stderr, "exec: command not found: ", err)
		return 1
	}

	cmd := exec.Command(path, args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	// fmt.Println("RUN", path, stdin, stdout)
	err = cmd.Run()
	// fmt.Println("EXIT", path)
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return exit.ExitCode()
		}

		fmt.Fprintln(stderr, "exec: command error: ", err)
		return 1
	}

	return 0
}
