package main

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

func grep(cfg *GrepConfig) (string, error) {
	var prefix, postfix string
	if cfg.ignoreCase {
		prefix +=  "(?i)"
	}

	if cfg.fixed {
		prefix += "^"
		postfix += "$"
	}

	re,err  := regexp.Compile(prefix + cfg.expr + postfix)
	if err != nil {
		return "", err
	}


}

type Grep struct {
	ignoreCase bool
	fixed bool
	expr string
	re *regexp.Regexp
}

func (g *Grep) Compile() error {
	expr := ""

	if cfg.ignoreCase {
		expr += "(?i)"
	}

	if cfg.fixed {
		expr += "^"
	}

	expr += g.expr

	if cfg.fixed {
		expr += "$"
	}

	g.re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

}

func (g *Grep) Match() (match string, ok bool) {

}


func main() {
	grep := &Grep{}



	safeContext := after > 0 || before >0 || context > 0

	input := []string
	matches := []int
	count := 0
	idx := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if safeContext {
			input = append(input, line1)
		}

		match, ok := grep.Match(line)

		if invert != ok {
			matches = append(matches, line)
			count++
		}
	}

}
