package main

import (
	"io"
	"log"
	"os"
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

type Cut struct {
	delimiter rune
	fieldCond func(idx int) bool
	separated bool

	r io.Reader
}

func NewCut(r io.Reader) *Cut {
	return &Cut{r: r}
}

func (c *Cut) Read(p []byte) (n int, err error) {
	c.r.Read()
}

func main() {
	cut := NewCut(os.Stdin)

	if _, err := io.Copy(os.Stdout, cut); err != nil {
		log.Fatal(err)
	}
}
