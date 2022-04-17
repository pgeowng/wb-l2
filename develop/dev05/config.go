package main

type Config struct {
	after int
	before int

	count bool
	ignoreCase bool
	invert bool
	fixed bool
	lineNumbers bool

	expr string
	filename string
}

func Cfg() *Config {
	cfg := &Config{}

	around := 0

	flag.InvVar(&cfg.after, "A", 0, "Print +N rows after match")
	flag.IntVar(&cfg.before, "B", 0, "Print +N rows before match")
	flag.IntVar(&around, "C", 0, "Print +N rows around match")
	flag.BoolVar(&cfg.count, "c", false, "Print amount of matched lines")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "Case-insensitive matching")
	flag.BoolVar(&cfg.invert, "v", false, "Inverted matching")
	flag.BoolVar(&cfg.fixed, "F", false, "Exact line matching")
	flag.BoolVar(&cfg.lineNumbers, "n", false, "Print line number")

	flag.Parse()

	if around > cfg.after {
		cfg.after = around
	}

	if around > cfg.before {
		cfg.before = around
	}

	if len(args) > 1 {
		cfg.file = args[1]
	}

	cfg.expr = args[0]

	return cfg
}