package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	// Обрабатываем кейс с первой цифрой. Так как '0'..'9' занимают один байт в utf-8.
	// То достаточно проверить первый байт строки
	if '0' <= input[0] && input[0] <= '9' {
		return "", errors.Errorf("first char can't be number")
	}

	isEscape := false
	wasLetter := false
	var result strings.Builder
	var char rune
	var count int
	for idx, ch := range input {
		isNumber := '0' <= ch && ch <= '9'
		isSlash := '\\' == ch

		if isEscape {
			if !isSlash && !isNumber {
				return "", errors.Errorf("escaped %c at %d position", ch, idx)
			}
			char = ch
			count = 0
			wasLetter = true
			isEscape = false
			continue
		}

		// Вывод происходит только когда встретили новую букву.
		// Чтобы обработать a0b и ab, сохраняем что последний раз был символ,
		// и печатаем либо указанное количество в count, либо 1, так как при ab, count == 0
		if !isNumber && idx > 0 {
			if wasLetter {
				count = 1
			}
			result.WriteString(strings.Repeat(string(char), count))
		}

		if isSlash && !isEscape {
			isEscape = true
			continue
		}

		if isNumber {
			count = count*10 + int(ch-'0')
			wasLetter = false
			continue
		}

		char = ch
		count = 0
		wasLetter = true
	}

	// Вывод происходит всегда перед обработкой,
	// поэтому необходимо вывести результат последней итерации.
	if !isEscape {
		if wasLetter {
			count = 1
		}
		result.WriteString(strings.Repeat(string(char), count))
	}

	return result.String(), nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result, err := Unpack(scanner.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, "unpack error:", err)
			os.Exit(2)
		}
		fmt.Println(result)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}

// echo "a2b3c4" | go run .
// aabbbcccc
