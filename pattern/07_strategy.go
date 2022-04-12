package pattern

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

// Strategy - behavior pattern
// Заключается в вынесении поведения из класса для его динамической смены владельцем класса.
// Паттерн помогает:
//   + Динамически изменять поведение класса.
//   Например, для класса рекомендаций можно изменить алгоритм сортировки.
//   По дате создания, по популярности, по релевантности.
//   + Упростить тестирование.
//   Например, достаточно протестировать класс поведения, который
//   слабо связан тем, где он будет использоваться.

// - Тот, кто использует класс знает о разных поведениях, что усложняет код.
// - Появление новых классов.

type Entry struct {
	name string
	kind string
	tags []string
}

// Context - объект, на котором можно изменить поведение
// передав функцию и/или экземпляр класса, который реализует интерфейс.
type Search struct {
	data    *[]*Entry
	filter  func(*Entry) bool
	sorting Sorting
}

func NewSearch(data *[]*Entry) *Search {
	return &Search{data: data}
}

// Strategy 1 - используем передаваемую функцию как метод s.filter(entry) bool
func (s *Search) WithFilter(fn func(*Entry) bool) *Search {
	s.filter = fn
	return s
}

func (s *Search) WithSort(sort Sorting) *Search {
	s.sorting = sort
	return s
}

// Реализуем поиск, учитывая переданное поведение.
func (s *Search) Select(query func(string) bool) []*Entry {
	result := []*Entry{}
	for _, entry := range *s.data {
		if s.filter != nil && !s.filter(entry) {
			continue
		}
		if query(entry.name) {
			result = append(result, entry)
		}
	}
	if s.sorting != nil {
		sort.Slice(result, func(i, j int) bool {
			return s.sorting.Less(result, i, j)
		})
	}
	return result
}

// Strategy 2 - используем интерфейс, реализовав который можно изменить поведение.
type Sorting interface {
	Less(arr []*Entry, i, j int) bool
}

type SortingByName struct{}

func (s *SortingByName) Less(arr []*Entry, i, j int) bool {
	return arr[i].name < arr[j].name
}

type SortingByKind struct{}

func (s *SortingByKind) Less(arr []*Entry, i, j int) bool {
	ii, jj := arr[i].kind, arr[j].kind
	return ii < jj || (ii == jj && arr[i].name < arr[j].name)
}

// Strategy 3 - используем closure, чтобы основной класс,
// не знал как именно происходит выбор.
func PrefixMatch(target string) func(string) bool {
	return func(field string) bool {
		return strings.HasPrefix(field, target)
	}
}

func FuzzyMatch(target string) func(string) bool {
	return func(field string) bool {
		t := []rune(target)
		f := []rune(field)
		for len(t) > 0 && len(f) > 0 {
			if unicode.ToLower(t[0]) == unicode.ToLower(f[0]) {
				t = t[1:]
			}
			f = f[1:]
		}
		return len(t) == 0
	}
}

func PrettyPrint(list []*Entry, msg string) {
	fmt.Println(msg)
	for _, entry := range list {
		fmt.Printf("[%4s] %s %v\n", entry.kind, entry.name, entry.tags)
	}
}

// Client - определяет как будет вести поиск
func useStrategy() {
	data := &[]*Entry{
		{"ach.mp3", "file", []string{"media", "audio", "mp3"}},
		{"AudioController.h", "file", []string{"cpp"}},
		{"back.aac", "file", []string{"media", "audio", "aac"}},
		{"docs", "dir", []string{"work"}},
	}

	s := NewSearch(data)

	var result []*Entry

	result = s.Select(PrefixMatch("ach"))
	PrettyPrint(result, "prefix: ach")

	result = s.Select(FuzzyMatch("ach"))
	PrettyPrint(result, "fuzzy: ach")

	s.WithFilter(func(e *Entry) bool {
		for _, tag := range e.tags {
			if tag == "audio" {
				return true
			}
		}
		return false
	})
	result = s.Select(PrefixMatch(""))
	PrettyPrint(result, "filter: audio")

	s.WithFilter(nil).WithSort(&SortingByKind{})
	result = s.Select(PrefixMatch(""))
	PrettyPrint(result, "sort: by kind")
}

// prefix: ach
// [file] ach.mp3 [media audio mp3]
// fuzzy: ach
// [file] ach.mp3 [media audio mp3]
// [file] AudioController.h [cpp]
// filter: audio
// [file] ach.mp3 [media audio mp3]
// [file] back.aac [media audio aac]
// sort: by kind
// [ dir] docs [work]
// [file] AudioController.h [cpp]
// [file] ach.mp3 [media audio mp3]
// [file] back.aac [media audio aac]
