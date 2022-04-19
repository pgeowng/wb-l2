package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

// go run . -a example
func main() {
	var url string

	flag.StringVar(&url, "a", "0.beevik-ntp.pool.ntp.org", "Set custom ntp server address")

	flag.Parse()

	response, err := ntp.Query(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query error: %v\n", err)
		os.Exit(1)
	}

	err = response.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "response error: %v\n", err)
		os.Exit(2)
	}

	now := time.Now()
	srv := now.Add(response.ClockOffset)
	fmt.Println("Time:", now.Format("15:04:05.00000"))
	fmt.Println("NTP: ", srv.Format("15:04:05.00000"))
	fmt.Println("Offset:", response.ClockOffset)
}

// $ go run .
// Time: 11:08:58.88030
// NTP:  11:08:58.88009
// Offset: -208.97µs
