package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Compare string

type Field struct {
	idx int
}

func ParseFields(val string) (result []Field, err error) {
	if len(val) == 0 {
		return
	}
	tokens := strings.Split(val, ",")

	for idx, token := range tokens {
		var n int64
		n, err = strconv.ParseInt(token, 10, 0)
		if err != nil {
			err = fmt.Errorf("fields key at %d not number: %v", idx, token)
			return
		}

		if n < 1 {
			err = fmt.Errorf("fields key at %d not positive integer: %v", idx, token)
			return
		}

		result = append(result, Field{idx: int(n) - 1})
	}

	return
}

type Config struct {
	fields []Field

	numeric bool
	unique  bool
	reverse bool

	filename string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	flag.BoolVar(&cfg.numeric, "n", false, "Use numeric sort")

	flag.BoolVar(&cfg.unique, "u", false, "Dont print same lines")
	flag.BoolVar(&cfg.reverse, "r", false, "Reverse order")

	var fields string
	flag.StringVar(&fields, "k", "", "Specify fields: -k 2,5")
	flag.Parse()

	var err error
	cfg.fields, err = ParseFields(fields)
	if err != nil {
		return nil, err
	}

	args := flag.Args()
	if len(args) > 0 {
		cfg.filename = args[0]
	}

	return cfg, nil
}
