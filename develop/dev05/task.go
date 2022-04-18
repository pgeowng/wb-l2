package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Grep struct {
	re *regexp.Regexp
}

func NewGrep(cfg *Config) (*Grep, error) {
	expr := cfg.expr

	if cfg.fixed {
		expr = "^" + expr + "$"
	}

	if cfg.ignoreCase {
		expr = "(?i)" + expr
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &Grep{re}, nil
}

func (g *Grep) Match(line string) bool {
	return g.re.MatchString(line)
}

type Program struct {
	cfg  *Config
	grep *Grep
}

func NewProgram(cfg *Config, grep *Grep) *Program {
	return &Program{cfg, grep}
}

func (p *Program) Run(r io.Reader, w io.Writer) error {
	input := []string{}
	matches := []int{}
	count := 0
	idx := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		ok := p.grep.Match(line)

		if p.cfg.invert != ok {
			matches = append(matches, idx)
			count++
		}

		input = append(input, line)
		idx++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if p.cfg.count {
		fmt.Fprintln(w, count)
		return nil
	}

	size := len(input)
	uniq := map[int]struct{}{}
	prevRight := -1
	showContext := p.cfg.after > 0 || p.cfg.before > 0

	for _, i := range matches {
		j := i - p.cfg.before
		if j < 0 {
			j = 0
		}
		right := i + p.cfg.after

		if showContext && prevRight+1 < j && prevRight != -1 {
			fmt.Fprintln(w, "--")
		}

		for ; j <= right && j < size; j++ {
			_, ok := uniq[j]
			if ok {
				continue
			}

			uniq[j] = struct{}{}
			if p.cfg.lineNumbers {
				fmt.Fprintf(w, "%d:", j+1)
			}
			fmt.Fprintln(w, input[j])
		}

		prevRight = right
	}

	return nil
}

func main() {
	cfg := Cfg()

	grep, err := NewGrep(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "grep:", err)
		os.Exit(1)
	}

	var src io.Reader

	if len(cfg.filename) > 0 {
		file, err := os.Open(cfg.filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "grep:", err)
			os.Exit(2)
		}
		defer file.Close()

		src = file
	} else {
		src = os.Stdin
	}

	prog := NewProgram(cfg, grep)
	err = prog.Run(src, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "grep:", err)
		os.Exit(3)
	}
}
