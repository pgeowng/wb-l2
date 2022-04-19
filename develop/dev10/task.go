package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type TelnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration) *TelnetClient {
	return &TelnetClient{
		address: address,
		timeout: timeout,
	}
}
func (tc *TelnetClient) Open() error {
	if tc.conn != nil {
		return fmt.Errorf("already open")
	}

	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}

	tc.conn = conn
	return nil
}

func (tc *TelnetClient) Close() error {
	if tc.conn != nil {
		if err := tc.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (tc *TelnetClient) Read(p []byte) (n int, err error) {
	n, err = tc.conn.Read(p)
	return
}

func (tc *TelnetClient) Write(p []byte) (n int, err error) {
	n, err = tc.conn.Write(p)
	return
}

// server: rm -f /tmp/f; mkfifo /tmp/f; cat /tmp/f | /bin/sh -i 2>&1 | nc -l 8241 > /tmp/f
// client: go run . localhost 8241
func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "t", time.Duration(10)*time.Second, "Connection timeout")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "tc: empty host: tc [-args] <host> [port]")
		os.Exit(1)
	}

	address := args[0]
	if len(args) == 2 {
		address += ":" + args[1]
	}

	tc := NewTelnetClient(address, timeout)

	err := tc.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tc: open error :", err)
		os.Exit(1)
	}

	go func() {
		if _, err := io.Copy(os.Stdout, tc); err != nil {
			fmt.Fprintln(os.Stderr, "tc: receive error :", err)
			os.Exit(1)
		}
	}()
	if _, err := io.Copy(tc, os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, "tc: send error :", err)
		os.Exit(1)
	}
	fmt.Println("STDIN CLOSED")
	tc.Close()
}
