package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/pgeowng/wb-l2/develop/dev08/command"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Command interface {
	Run(args []string, vars map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int
}

type Shell struct {
	Vars          map[string]string
	Commands      map[string]Command
	TokenizeSplit func(line string) []string
}

func NewShell(vars map[string]string, commands map[string]Command) (*Shell, error) {
	return &Shell{
		Vars:     vars,
		Commands: commands,
		TokenizeSplit: func(line string) []string {
			return strings.Fields(line)
		},
	}, nil
}

type PipeItem struct {
	cmd  string
	args []string
}
type Pipeline struct {
	fork bool
	pipe []PipeItem
}

// Парсим команды разбивая & и |.
func (sh *Shell) Tokenize(line string) (actions []Pipeline, err error) {
	tokens := sh.TokenizeSplit(line)

	actions = append(actions, Pipeline{})
	idx := 0

	for i := 0; i < len(tokens); i++ {
		if actions[idx].fork {
			actions = append(actions, Pipeline{})
		}

		cmd := tokens[i]
		if cmd == "|" || cmd == "&" {
			err = fmt.Errorf("sh: parse error near %d: %v", i, cmd)
			return
		}

		i++
		args := []string{}

		for ; i < len(tokens) && tokens[i] != "|" && tokens[i] != "&"; i++ {
			args = append(args, tokens[i])
		}

		actions[idx].pipe = append(actions[idx].pipe, PipeItem{cmd, args})

		if i < len(tokens) && tokens[i] == "&" {
			actions[idx].fork = true
		}
	}

	return
}

// Игнорируем ctrl+c
func HandleInterrupt(cb func(chan os.Signal)) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		for range interrupt {
			cb(interrupt)
		}
		signal.Stop(interrupt)
	}()
}

func (sh *Shell) Run(stdin io.Reader, stdout io.WriteCloser, stderr io.Writer) {
	HandleInterrupt(func(_ chan os.Signal) {
		fmt.Fprintf(stdout, "\n%s $ ", sh.Vars["PWD"])
	})

	prog, ok := sh.Commands["help"]
	if ok {
		prog.Run([]string{}, sh.Vars, stdin, stdout, stderr)
	}
	fmt.Fprintf(stdout, "\n%s $ ", sh.Vars["PWD"])

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" {
			fmt.Fprintln(stdout, "Goodbye! :(")
			return
		}

		actions, err := sh.Tokenize(line)
		if err != nil {
			fmt.Fprintln(stderr, err)
		}

		// При использовании |, связывает cmd1.stdout -> cmd2.stdin, через io.Pipe
		for _, action := range actions {
			var dest io.WriteCloser = stdout
			var src io.Reader = stdin
			var nextSrc io.Reader = stdin
			wg := sync.WaitGroup{}

			for idx, item := range action.pipe {
				isLast := idx+1 == len(action.pipe)
				if !isLast {
					pr, pw := io.Pipe()
					nextSrc = pr
					dest = pw
				} else {
					dest = stdout
				}

				prog, ok := sh.Commands[item.cmd]
				if !ok {
					prog = sh.Commands["exec"]
					item.args = append([]string{item.cmd}, item.args...)
				}

				wg.Add(1)
				go func(args []string, stdin io.Reader, stdout io.WriteCloser, mustClose bool) {
					defer wg.Done()
					exitCode := prog.Run(args, sh.Vars, stdin, stdout, stderr)
					if exitCode != 0 {
						fmt.Fprintf(stderr, "%s: exit code %d", item.cmd, exitCode)
					}
					if mustClose {
						stdout.Close()
					}
				}(item.args, src, dest, !isLast)
				src = nextSrc

			}

			if !action.fork {
				wg.Wait()
			}
		}
		fmt.Fprintf(stdout, "\n%s $ ", sh.Vars["PWD"])
	}
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	vars := map[string]string{
		"PWD":  pwd,
		"HOME": home,
	}

	commands := map[string]Command{
		"cd":   &command.ChangeDir{},
		"pwd":  &command.ProcessWD{},
		"echo": &command.Echo{},
		"exec": &command.Exec{},
		"help": &command.Help{},
		"kill": &command.Kill{},
	}

	prog, err := NewShell(vars, commands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	prog.Run(os.Stdin, os.Stdout, os.Stderr)
}

// $ go run .
// Hello from go-shell! You can use:
//   help                  - show this message
//   exit                  - exit shell :(
//   cd [path]             - change directory
//   pwd                   - current directory
//   echo [...args]        - prints to stdout args
//   kill <pid>            - kill process by id
//   <any PATH executable> - execute file from PATH
//   <cmd1> | <cmd2>       - pipe <cmd1> stdout to <cmd2> stdin
//   <cmd1> &              - run <cmd1> in background
// /home/dt/gohigh/wb/wb-l2/develop/dev08 $
