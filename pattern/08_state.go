package pattern

import (
	"fmt"

	"github.com/pkg/errors"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

// State - behavior pattern
// Заключается в выборе текущих действий из набора
// на основе внутреннего состояния.
// Паттерн помогает:
//   + Упростить сложную логику поведения,
//     представив её набором состояний и переходами между ними.
//     Например, кофемашина для различного вида кофе использует разные рецепты.
//     И соответственно, разные промежуточные состояния.

// - Значительно усложняет код, что может быть лишним, если
//   сложность логики не так высока.

// State - определяем все действия во время состояний
type ATMState interface {
	Balance() error
	InsertCard(card string) error
	ReturnCard() error
	Message()
}

// Context - сохраняем в embedded field, чтобы унаследовать методы.
// Также все возможные состояния в виде переменных, чтобы не перевыделять их.
type ATM struct {
	ATMState
	card string

	cardMap map[string]int

	idle ATMState
	work ATMState
}

// Так как состояния обязаны знать о полях владельца,
// необходимо сохранить указатель на него.
func NewATM(cardMap map[string]int) *ATM {
	idle := &IdleATMState{}
	work := &WorkATMState{}

	atm := &ATM{
		cardMap:  cardMap,
		ATMState: idle,
		idle:     idle,
		work:     work,
	}

	idle.ATM = atm
	work.ATM = atm

	return atm
}

// Используем базовый класс, который не поддерживает операции,
// чтобы уменьшить дублирование. Однако, это не обязательно.
type BaseATMState struct {
	*ATM
}

func (s *BaseATMState) Balance() error {
	return errors.New("access denied for balance")
}
func (s *BaseATMState) InsertCard(_ string) error {
	return errors.New("insert card can't be performed")
}
func (s *BaseATMState) ReturnCard() error {
	return errors.New("return card can't be performed")
}

// State 1
type IdleATMState struct{ BaseATMState }

// Через *ATM меняем поля и состояния по необходимым условиям.
func (s *IdleATMState) InsertCard(card string) error {
	_, ok := s.ATM.cardMap[card]
	if !ok {
		return errors.New("card not found")
	}
	s.ATM.card = card
	s.ATM.ATMState = s.ATM.work
	return nil
}

func (s *IdleATMState) Message() {
	fmt.Println("atm: Welcome! Please insert card for further actions.")
}

// State 2
type WorkATMState struct{ BaseATMState }

func (s *WorkATMState) Balance() error {
	fmt.Printf("atm: balance - %d\n", s.ATM.cardMap[s.ATM.card])
	return nil
}
func (s *WorkATMState) ReturnCard() error {
	fmt.Println("atm: Returning card...")
	s.ATM.card = ""
	s.ATM.ATMState = s.ATM.idle
	return nil
}
func (s *WorkATMState) Message() {
	fmt.Println("atm: Waiting command: balance, return")
}

func useState() {
	atm := NewATM(map[string]int{
		"4182": 100,
		"5592": 0,
	})

	var err error

	atm.Message()
	err = atm.ReturnCard()
	if err != nil {
		fmt.Println("error:", err)
	}

	atm.Message()
	err = atm.Balance()
	if err != nil {
		fmt.Println("error:", err)
	}

	atm.Message()
	fmt.Print("inserting 4211: ")
	err = atm.InsertCard("4211")
	if err != nil {
		fmt.Println("error:", err)
	}

	atm.Message()
	fmt.Print("inserting 4182: ")
	err = atm.InsertCard("4182")
	if err != nil {
		fmt.Println("error:", err)
	}

	atm.Message()
	err = atm.Balance()
	if err != nil {
		fmt.Println("error:", err)
	}

	err = atm.ReturnCard()
	if err != nil {
		fmt.Println("error:", err)
	}

}

// atm: Welcome! Please insert card for further actions.
// error: return card can't be performed
// atm: Welcome! Please insert card for further actions.
// error: access denied for balance
// atm: Welcome! Please insert card for further actions.
// inserting 4211: error: card not found
// atm: Welcome! Please insert card for further actions.
// inserting 4182: atm: Waiting command: balance, return
// atm: balance - 100
// atm: Returning card...
