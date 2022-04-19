package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// func GetKeys(field string, left int, right int, leading bool) (string, error) {
// 	8
// }

// type Cmp func(arr []string, i, j int) bool
// type Algo interface {
// 	AddEntry(string)
// 	Perform(Cmp) error
// 	Result()
// }

// type Program struct {
// 	Algo
// }

// type Sort struct {
// 	entries []string
// }

// func (s *Sort) AddEntry(entry string) {
// 	s.entries = append(s.entries, entry)
// }

// func (s *Sort) Perform(cmp Cmp) error {
// 	sort.Slice(s.entries, func(i, j int) bool {
// 		return cmp(s.entries, i, j)
// 	})
// 	return nil
// }

// func (s *Sort) Result() {
// 	for _, entry := range s.entries {
// 		fmt.Fprintln(os.Stdout, entry)
// 	}
// }

func LexicalCmp(a, b string) bool {
	return a < b
}

func NumericCmp(a, b string) bool {
	na, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return true
	}

	nb, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return false
	}

	return na < nb
}

func PerFieldCmp(a, b string, fields []Field, fieldCmp func(string, string) bool) bool {
	aFields := strings.Fields(a)
	aSize := len(aFields)
	bFields := strings.Fields(b)
	bSize := len(bFields)

	for _, field := range fields {
		idx := field.idx
		if idx >= aSize {
			return true
		}

		if idx >= bSize {
			return false
		}

		if fieldCmp(aFields[idx], bFields[idx]) {
			return true
		}
	}

	return false
}

type Sort struct {
	cfg      *Config
	fieldCmp func(a, b string) bool
}

func NewSort(cfg *Config) *Sort {
	sort := &Sort{
		cfg:      cfg,
		fieldCmp: LexicalCmp,
	}

	if cfg.numeric {
		sort.fieldCmp = NumericCmp
	}

	return sort
}

func (s *Sort) Run(stdin io.Reader, stdout io.Writer) error {
	lines := []string{}
	scanner := bufio.NewScanner(stdin)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	isPerFieldCmp := len(s.cfg.fields) > 1
	sort.Slice(lines, func(i, j int) bool {
		less := false
		if isPerFieldCmp {
			less = PerFieldCmp(lines[i], lines[j], s.cfg.fields, s.fieldCmp)
		} else {
			less = s.fieldCmp(lines[i], lines[j])
		}

		return less
	})

	if s.cfg.unique {
		temp := []string{}
		uniq := map[string]struct{}{}
		for _, line := range lines {
			if _, ok := uniq[line]; !ok {
				temp = append(temp, line)
				uniq[line] = struct{}{}
			}
		}
		lines = temp
	}

	if s.cfg.reverse {
		size := len(lines)
		mid := size / 2
		for idx := 0; idx < mid; idx++ {
			lines[idx], lines[size-1-idx] = lines[size-1-idx], lines[idx]
		}
	}

	for _, line := range lines {
		fmt.Fprintln(stdout, line)
	}

	return nil
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "sort:", err)
		os.Exit(1)
	}

	var src io.Reader

	if len(cfg.filename) > 0 {
		file, err := os.Open(cfg.filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sort:", err)
			os.Exit(2)
		}
		defer file.Close()

		src = file
	} else {
		src = os.Stdin
	}

	prog := NewSort(cfg)
	err = prog.Run(src, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sort:", err)
		os.Exit(3)
	}
}

// echo "1\n2\n10\n53" | go run .
// 1
// 10
// 2
// 53

// echo "1\n2\n10\n53" | go run . -n
// 1
// 2
// 10
// 53
