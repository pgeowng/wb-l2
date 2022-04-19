package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

*/
func or(channels ...<-chan interface{}) <-chan interface{} {
	result := make(chan interface{})
	wg := sync.WaitGroup{}

	go func() {
		wg.Wait()
		close(result)
	}()

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for v := range ch {
				result <- v
			}
		}(ch)
	}

	return result
}

func main() {
	sig := func(after time.Duration, num int) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
			c <- num
		}()
		return c
	}

	start := time.Now()

	tasks := []<-chan interface{}{}

	for i, v := range []int{2, 3, 5, 7, 11} {
		tasks = append(tasks, sig(time.Duration(i*200+100)*time.Millisecond, v))
	}

	for v := range or(tasks...) {
		fmt.Println(v, time.Since(start))
	}

	fmt.Printf("done after %v", time.Since(start))

}

// $ go run .
// 2 100.18847ms
// 3 300.497629ms
// 5 500.76971ms
// 7 701.087599ms
// 11 900.320515ms
// done after 900.336815ms%
