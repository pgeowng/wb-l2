package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ChangeDir struct{}

func (cd *ChangeDir) Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	var target string

	if len(args) == 0 {
		target = vars["HOME"]
	} else {
		target = args[0]
	}

	if !filepath.IsAbs(target) {
		target = filepath.Clean(filepath.Join(vars["PWD"], target))
	}

	if !filepath.IsAbs(target) {
		fmt.Fprintln(stderr, "cd: cant enter", target)
		return 1
	}

	info, err := os.Stat(target)
	if os.IsNotExist(err) {
		fmt.Fprintln(stderr, "cd: no such file or directory", target)
		return 1
	}

	if !info.IsDir() {
		fmt.Fprintln(stderr, "cd: no a directory", target)
		return 1
	}

	os.Chdir(target)
	vars["PWD"] = target
	return 0
}
