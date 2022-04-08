package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

// Chain of responsibility - behavior pattern
// Заключается в создании цепочки объктов, каждый элемент которой обрабатывают запрос и/или отдают следующему.
// Паттерн помогает:
//   + Отделить зоны ответственности при обработке запроса.
//     Например, часть запросов к сервису не может быть выполнена
//     без аутентификации.
//   + Динамически менять исполняемый код в зависимости от контекста.
//     Например, в зависимости от наведённой области ui будет менятся
//     реакция на нажатие комбинаций клавиш.
//   + Легко добавлять новое поведение. Достаточно реализовать
//     интерфейc элемента цепочки, и его можно встроить.

// - Тяжелее анализировать, что выполняется в runtime.
type Message struct {
	username string
	password string
	action   string
}

// Объявляем интерфейс элемента цепочки,
// который может сохранить следующий за собой элемент
type Handler interface {
	Next(Handler)
	Handle(Message)
}

// Для удобства сделаем класс, от которого будем насследовать базовое поведение
type baseHandler struct{ next Handler }

func (bh *baseHandler) Next(h Handler) { bh.next = h }

// Реализуем обработчики

//Handler 1 ---
type AuthHandler struct {
	baseHandler
	cred map[string]string
}

func NewAuth(cred map[string]string) *AuthHandler {
	return &AuthHandler{cred: cred}
}

func (a *AuthHandler) Handle(r Message) {
	pass, ok := a.cred[r.username]
	if !ok {
		fmt.Println("auth: bad username")
		return
	}

	if pass != r.password {
		fmt.Println("auth: bad password")
		return
	}

	// Проверяем, есть ли следующий и вызываем
	if a.next != nil {
		a.next.Handle(r)
	}
}

// Handler 2 ---
type RootHandler struct{ baseHandler }

func (ra *RootHandler) Handle(r Message) {
	if r.username == "root" {
		switch r.action {
		case "rm -rf":
			fmt.Println("root: removing all files")
			return
		}
	}

	if ra.next != nil {
		ra.next.Handle(r)
	}
}

// Handler 3 ---
type UserHandler struct{ baseHandler }

func (ua *UserHandler) Handle(r Message) {
	switch r.action {
	case "ls":
		fmt.Println("listing files...\nbin/\nDocuments/\n123.mkv")
		return
	case "hostname":
		fmt.Println("hostname: p-1-42")
		return
	}

	if ua.next != nil {
		ua.next.Handle(r)
	}
}

// Handler 4 ---
type ErrorHandler struct{ baseHandler }

func (eh *ErrorHandler) Handle(r Message) {
	fmt.Println("error: unknown message")
	return
}

// Handler 5 ---
var count int = 0

type LogHandler struct{ baseHandler }

func (lh *LogHandler) Handle(r Message) {
	count++
	fmt.Println("-- request", count)
	if lh.next != nil {
		lh.next.Handle(r)
	}
}

// Client ---
// Тот класс, который начинает использование цепочки.

type Client struct {
	handler Handler
}

func NewClient(handler Handler) *Client {
	return &Client{handler}
}

func (c *Client) Send(r Message) {
	if c.handler != nil {
		c.handler.Handle(r)
	}
}

func useChain() {
	log := &LogHandler{}
	auth := NewAuth(map[string]string{
		"rattmann": "12345",
		"root":     "xqc",
	})
	root := &RootHandler{}
	user := &UserHandler{}
	error := &ErrorHandler{}

	cl := NewClient(log)
	log.Next(auth)
	auth.Next(root)
	root.Next(user)
	user.Next(error)

	cl.Send(Message{})
	cl.Send(Message{"rattmann", "idk", "ls"})
	cl.Send(Message{"rattmann", "12345", "ls"})
	cl.Send(Message{"root", "12345", "rm -rf"})
	cl.Send(Message{"root", "xqc", "rm -rf"})
}

// -- request 1
// auth: bad username
// -- request 2
// auth: bad password
// -- request 3
// listing files...
// bin/
// Documents/
// 123.mkv
// -- request 4
// auth: bad password
// -- request 5
// root: removing all files
