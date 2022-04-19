package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	colly "github.com/gocolly/colly/v2"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Config struct {
	url      string
	hostname string

	depth          int
	onlySubdomains bool
}

func NewConfig() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.depth, "l", 0, "Max depth level. 0 means infinite")
	flag.BoolVar(&cfg.onlySubdomains, "s", false, "Only match subdomains")

	flag.Parse()

	if cfg.depth < 0 {
		fmt.Println("wget: depth is non-negative integer, where 0 means infinite depth")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("wget: url is not specified")
		os.Exit(1)
	}

	link := args[0]
	u, err := url.ParseRequestURI(link)
	if err != nil {
		fmt.Println("wget: bad url", err)
		os.Exit(1)
	}

	cfg.url = u.String()
	cfg.hostname = u.Hostname()

	return cfg
}

type Wget struct {
	*colly.Collector
	visited map[string]struct{}
}

func NewWget(cfg *Config) *Wget {
	opts := []colly.CollectorOption{}

	if cfg.depth > 0 {
		opts = append(opts, colly.MaxDepth(cfg.depth))
	}

	if cfg.onlySubdomains {
		reg, err := regexp.Compile(`https?://.*?` + cfg.hostname)
		if err != nil {
			fmt.Println("bad regexp hostname", cfg.hostname)
			os.Exit(1)
		}

		opts = append(opts, colly.URLFilters(reg))
	}

	return &Wget{
		Collector: colly.NewCollector(opts...),
		visited:   map[string]struct{}{},
	}
}

func (w *Wget) HandlePageLink(attr string) func(*colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr(attr))
		if _, ok := w.visited[link]; !ok {
			w.visited[link] = struct{}{}
			w.Visit(link)
		}
	}
}
func main() {
	cfg := NewConfig()
	wget := NewWget(cfg)

	wget.OnHTML("a[href]", wget.HandlePageLink("href"))
	wget.OnHTML("link[href]", wget.HandlePageLink("href"))
	wget.OnHTML("img[src]", wget.HandlePageLink("href"))

	wget.OnResponse(func(r *colly.Response) {
		u := r.Request.URL
		path := u.Path

		if filepath.Ext(path) == "" {
			path = filepath.Join(path, "index.html")
		}

		fullpath := filepath.Join(u.Hostname(), path)
		dir := filepath.Dir(fullpath)
		if _, err := os.Stat(dir); err != nil {
			os.MkdirAll(dir, os.ModePerm)
		}

		fmt.Println("loading: ", fullpath)
		r.Save(fullpath)
	})

	if err := wget.Visit(cfg.url); err != nil {
		log.Fatal(err)
	}
	wget.Wait()
}

// $ go run . https://example.com
// loading:  example.com/index.html
// loading:  www.iana.org/domains/reserved/index.html
// loading:  www.iana.org/index.html
// loading:  www.iana.org/about/index.html
// loading:  www.iana.org/domains/index.html
// ^Csignal: interrupt
