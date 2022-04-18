package main

import (
	"flag"
)

type Config struct {
	after  int
	before int

	count       bool
	ignoreCase  bool
	invert      bool
	fixed       bool
	lineNumbers bool

	expr     string
	filename string
}

func NewConfig() *Config {
	cfg := &Config{}

	context := 0

	flag.IntVar(&cfg.after, "A", 0, "Print +N rows after match")
	flag.IntVar(&cfg.before, "B", 0, "Print +N rows before match")
	flag.IntVar(&context, "C", 0, "Print +N rows around match")
	flag.BoolVar(&cfg.count, "c", false, "Print amount of matched lines")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "Case-insensitive matching")
	flag.BoolVar(&cfg.invert, "v", false, "Inverted matching")
	flag.BoolVar(&cfg.fixed, "F", false, "Exact line matching")
	flag.BoolVar(&cfg.lineNumbers, "n", false, "Print line number")

	flag.Parse()

	if context > cfg.after {
		cfg.after = context
	}

	if context > cfg.before {
		cfg.before = context
	}

	args := flag.Args()

	if len(args) > 1 {
		cfg.filename = args[1]
	}

	if len(args) > 0 {
		cfg.expr = args[0]
	}

	return cfg
}
