package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Основная идея заключается в присвоении номера для любой комбинации
// букв из русского алфавита. Выполняется за O(n).
// ""   = 0
// "а"  = 1
// "я"  = 32
// "ё"  = 33
// "аа" = 34
// "аё" = "ёа" = 66
// "ёё" = 98
// "ааа" = 99
func hashLettersRussian(word string) (uint64, error) {
	if len(word) == 0 {
		return 0, nil
	}
	const base = 32

	var key uint64
	var size uint64 = 0
	for idx, r := range word {
		if r == '\u0451' {
			key += 32
		} else if r < '\u0430' || r > '\u044F' {
			return 0, errors.Errorf("letter %c at %d is not from russian alphabet", r, idx)
		} else {
			key += uint64(r - '\u0430')
		}
		size++
	}
	key += (size - 1) + (size-1)*size/2*base + 1

	return key, nil
}

// Для каждого слова вычисляем его хеш.
// Сохраняем в dict map, проверяя уникальность слова.
// Затем по сортируем каждое множество
// Сложность: пусть n - кол-во всех слов, m - макс длина слова
// Если все в одной группе: O(n*m) + O(n*logn)
// Если все в разных группах: O(n*m) + O(n)
func groupAnagrams(list []string) map[string][]string {
	dict := map[uint64][]string{}
	uniq := map[string]struct{}{}

	for _, word := range list {
		if len(word) != 0 {
			lword := strings.ToLower(word)
			_, ok := uniq[lword]
			if ok {
				continue
			}
			key, err := hashLettersRussian(lword)
			if err != nil {
				continue
			}
			dict[key] = append(dict[key], lword)
			uniq[lword] = struct{}{}
		}
	}

	result := map[string][]string{}
	for _, set := range dict {
		if len(set) > 1 {
			result[set[0]] = set
			sort.Strings(set)
		}
	}

	return result
}

// $ echo "лягушка  гушялка гав ваг гва  а" | go run .
// гав (3): ваг, гав, гва
// лягушка (2): гушялка, лягушка
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	input := []string{}
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	result := groupAnagrams(input)

	if len(result) == 0 {
		fmt.Println("no groups")
		return
	}

	keys := make([]string, 0, len(result))

	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		words := result[k]
		fmt.Printf("%s (%d): %s\n", k, len(words), strings.Join(words, ", "))
	}
}
