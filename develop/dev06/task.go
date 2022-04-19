package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func ParseFields(fields string) (left int, right int, err error) {
	tokens := strings.Split(fields, "-")

	left, right = -1, -1
	switch len(tokens) {
	case 1:
		var tmp int
		n, _ := fmt.Sscan(tokens[0], &tmp)
		if n > 0 && tmp < 1 {
			err = fmt.Errorf("bad fields argument, expected 1<=N")
			return
		}
		left, right = tmp, tmp
	case 2:
		n1, _ := fmt.Sscan(tokens[0], &left)
		n2, _ := fmt.Sscan(tokens[1], &right)
		if n1 == 0 && n2 == 0 {
			err = fmt.Errorf("bad fields argument, expected N,N-,N-M,-M")
			return
		}
		if n1 > 0 && left < 1 ||
			n2 > 0 && right < 1 ||
			left > right && right != -1 {
			err = fmt.Errorf("bad fields argument, expected 1<=N<=M")
			return
		}
	default:
		err = fmt.Errorf("bad fields argument, expected N,N-,N-M,-M")
	}

	return
}

type Config struct {
	leftField     int
	rightField    int
	delimiter     string
	onlyDelimited bool

	filename string
}

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{}

	var fields string
	flag.StringVar(&fields, "f", "", "Select only these fields")
	flag.StringVar(&cfg.delimiter, "d", "\t", "Use delim instead TAB")
	flag.BoolVar(&cfg.onlyDelimited, "s", false, "Do not print lines not containing delimiters")

	flag.Parse()

	left, right, err := ParseFields(fields)
	if err != nil {
		return
	}

	cfg.leftField = left
	cfg.rightField = right

	if len([]rune(cfg.delimiter)) > 1 {
		err = fmt.Errorf("delimiter must be single character")
		return
	}

	args := flag.Args()

	if len(args) > 0 {
		cfg.filename = args[0]
	}

	return
}

type Cut struct {
	cfg *Config
}

func NewCut(cfg *Config) *Cut {
	return &Cut{cfg}
}

func (c *Cut) Run(r io.Reader, w io.Writer) error {

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Split(line, c.cfg.delimiter)

		if len(fields) == 1 && line == fields[0] {
			if !c.cfg.onlyDelimited {
				fmt.Fprintln(w, line)
			}
			continue
		}

		left := c.cfg.leftField - 1
		if left < 0 {
			left = 0
		}

		right := c.cfg.rightField
		if right < 0 || right > len(fields) {
			right = len(fields)
		}

		if left > right {
			left = right
		}

		fields = fields[left:right]
		fmt.Fprintln(w, strings.Join(fields, c.cfg.delimiter))
	}

	return nil
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cut:", err)
		os.Exit(1)
	}

	var src io.Reader

	if len(cfg.filename) > 0 {
		file, err := os.Open(cfg.filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cut:", err)
			os.Exit(2)
		}
		defer file.Close()

		src = file
	} else {
		src = os.Stdin
	}

	prog := NewCut(cfg)
	err = prog.Run(src, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cut:", err)
		os.Exit(3)
	}
}

// echo "linux" | go run . -d "n" -f 2
// ux
