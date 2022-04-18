package command

import (
	"fmt"
	"io"
)

type Help struct{}

var message string = `Hello from go-shell! You can use:
  help                  - show this message
  exit                  - exit shell :(
  cd [path]             - change directory
  pwd                   - current directory
  echo [...args]        - prints to stdout args
  kill <pid>            - kill process by id
  <any PATH executable> - execute file from PATH
  <cmd1> | <cmd2>       - pipe <cmd1> stdout to <cmd2> stdin
  <cmd1> &              - run <cmd1> in background`

func (cd *Help) Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	fmt.Fprint(stdout, message)

	return 0
}
