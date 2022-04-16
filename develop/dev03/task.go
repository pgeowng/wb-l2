package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

func uniq(arr []string) {}

func main() {
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	sort.Strings(lines)

	// if err := prog.Perform(); err != nil {
	// 	fmt.Fprintln(os.Stderr, "sort:", err)
	// 	os.Exit(2)
	// }

	for _, line := range lines {
		fmt.Println(line)
	}
}
