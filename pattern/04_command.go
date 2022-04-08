package pattern

import (
	"fmt"
	"math/rand"
)

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

// Command - behavior pattern
// Заключается в превращении вызова функции в отдельный класс, в котором есть метод для исполнения действия.
// Паттерн помогает, если:
//   + Нужно изменять поведение во время исполнения, и/или переиспользовать код.
//     Например, поведение ui кнопки/комбинации клавиш меняется
//     в зависимости от текущего положения курсора/фокуса.
//     Или одно и тоже действие вызывается в ui и по комбинации клавиш.
//   + Операции могут быть отложены, отменены.
//     Например, программа для написания музыки с помощью ui подготавливает набор комманд,
//     которые будут интерпретированы во время воспроизведения/рендера звука.

// - Появление дополнительно слоя, который напрямую зависит от бизнес логики.
//   Необходимо менять вместе и поддерживать.

// Объявляем интерфейс который может вызваться и
// проверить возможность исполнения.
type Command interface {
	Execute()
	Test() bool
}

type Order struct {
	market Market
	target string
	price  int
}

// Command 1 ---
// Комманды используют Market в качестве цели исполнения.
type BuyOrder struct{ Order }

func NewBuyOrder(market Market, target string, price int) *BuyOrder {
	return &BuyOrder{Order: Order{market, target, price}}
}

func (o *BuyOrder) Execute() {
	o.market.Buy(o.target)
	fmt.Printf("Buy %s for %d\n", o.target, o.market.Price(o.target))
}

func (o *BuyOrder) Test() bool {
	return o.market.Price(o.target) <= o.price
}

// Command 2 ---
type SellOrder struct{ Order }

func NewSellOrder(market Market, target string, price int) *SellOrder {
	return &SellOrder{Order: Order{market, target, price}}
}
func (o *SellOrder) Execute() {
	o.market.Sell(o.target)
	fmt.Printf("Sell %s for %d\n", o.target, o.market.Price(o.target))
}
func (o *SellOrder) Test() bool {
	return o.market.Price(o.target) >= o.price
}

// Receiver ---
// Предоставляем некоторые действия с собой.
type Market interface {
	Buy(target string)
	Sell(target string)
	Price(target string) int
	Update()
}

type Steam struct {
	prices  map[string]int
	balance int
}

func NewSteam(balance int) *Steam {
	return &Steam{prices: make(map[string]int), balance: balance}
}

func (s *Steam) Buy(target string) {
	s.balance -= s.prices[target]
}

func (s *Steam) Sell(target string) {
	s.balance += s.prices[target]
}

func (s *Steam) Price(target string) int {
	return s.prices[target]
}

func (s *Steam) Balance() int {
	return s.balance
}

func (s *Steam) Update() {
	s.prices["AK-47 | Redline"] = rand.Intn(9) + 14
	s.prices["Operation Bravo"] = rand.Intn(10) + 42
}

// Invoker --
// Как-то добавляем комманды и проверяем возможность исполнения
func useCommand() {
	startBalance := 150
	steam := NewSteam(startBalance)

	fmt.Printf("Start balance: %d usd\n", startBalance)

	orders := []Command{
		NewBuyOrder(steam, "AK-47 | Redline", 15),
		NewBuyOrder(steam, "Operation Bravo", 44),
		NewSellOrder(steam, "AK-47 | Redline", 20),
		NewSellOrder(steam, "Operation Bravo", 50),
	}

	for len(orders) > 0 {
		steam.Update()
		for idx, order := range orders {
			if order.Test() {
				order.Execute()
				lastIdx := len(orders) - 1
				orders[idx] = orders[lastIdx]
				orders = orders[:lastIdx]
				break
			}
		}
	}

	fmt.Printf("Income: %d usd\n", steam.Balance()-startBalance)
}

// Start balance: 150 usd
// Sell Operation Bravo for 51
// Buy Operation Bravo for 42
// Buy AK-47 | Redline for 14
// Sell AK-47 | Redline for 21
// Income: 16 usd
