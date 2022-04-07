package pattern

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

// Builder - creational pattern
// Заключается в отделении инициализации класса от его поведения.
// Паттерн помогает, если:
//   + Инициализация объекта использует слишком много аргументов,
//     и/или их использование опционально, и/или они взаимозаменяемы.
//     Например, при работе http-клиента мы xотим кастомное поведение при redirect, timeout, forbidden.
//   + Создание объекта невозможно сразу или класс удобно наполнять поэтапно.
//     Например, мы хотим во всех дочерних элементах сохранить источник, откуда он пришел.
//       show := GetShow("...")
//       epBase := NewEpisodeBuilder().WithShow(show)
//       eps := []Episode{
//         epBase.WithEp(show.latest).Build(),
//         epBase.WithEp(show.additional).Build(),
//       }
//   + Необходимо создавать объекты на основе шаблона с небольшими изменениями.
//     Если реализовать immutable изменение данных, то шаблоны можно использовать в concurrent среде.
//     Например, создание в UI кнопок с разным текстом, поведением.
//   + Нужно управлять тем, какие объекты будут созданы во всей системе.
//     Отделение инициализации от поведения позволяет нам независимо от других модулей вносить изменения.
//     Например, если реализовать класс Director, которых хранит Builder актуальной версии.
//     При обновлении содержания достаточно изменить конфигурацию в Director,
//     при этом клиентский код будет незатронут.

// - Появление дополнительных классов, которые необходимо поддерживать и тестировать.

// Мы хотим создавать документы, предварительно валидируя поля.
// Document - валидный документ, готовый к обработке.
// DocumentBuilder - создатель документов, проверяющий их правильность.
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

type Topic string

const (
	Undefined Topic = ""
	Agreement Topic = "Agreement"
	Request   Topic = "Request"
	Statement Topic = "Statement"
	Applying  Topic = "Applying"
)

type Document struct {
	person     Person
	department string
	topic      Topic
	date       time.Time
}

func (d Document) String() string {
	result := fmt.Sprintf(
		`To %s department
From %s %s
At %s
Subject: %s`,
		d.department,
		d.person.FirstName,
		d.person.LastName,
		d.date,
		d.topic,
	)
	return result
}

// На основе переданных условий, будет производится валидация Person, Topic.
// Если конфигурация невозможна, сохраняется ошибка, чтобы в конце сообщить о ней.
type DocumentBuilder struct {
	document         Document
	personValidation func(Person) error
	topicValidation  func(Topic) error
	err              error
}

func NewDocumentBuilder() DocumentBuilder {
	return DocumentBuilder{}
}

// Запись полей реализована с созданием новых объектов.
// Поэтому в любой момент можно использовать builder в качестве шаблона для новых объектов.
func (ob DocumentBuilder) WithPerson(p Person) DocumentBuilder {
	ob.document.person = p
	return ob
}

func (ob DocumentBuilder) WithDepartment(department string) DocumentBuilder {
	if department != "finance" && department != "accounting" {
		ob.err = errors.Wrapf(ob.err, "only accepting documents to finance and accounting department: %s", department)
		return ob
	}

	ob.document.department = department
	return ob
}

func (ob DocumentBuilder) WithPersonValidation(fn func(Person) error) DocumentBuilder {
	ob.personValidation = fn
	return ob
}

func (ob DocumentBuilder) WithNoTopic() DocumentBuilder {
	ob.document.topic = Undefined
	return ob
}

func (ob DocumentBuilder) WithTopic(topic Topic) DocumentBuilder {
	ob.document.topic = topic
	return ob
}

func (ob DocumentBuilder) WithTopicValidation(fn func(Topic) error) DocumentBuilder {
	ob.topicValidation = fn
	return ob
}

func (ob DocumentBuilder) WithDate(date time.Time) DocumentBuilder {
	ob.document.date = date
	return ob
}

// Завершающий метод, который создает документ
func (ob DocumentBuilder) Build() (*Document, error) {
	if ob.err != nil {
		return nil, ob.err
	}

	if ob.personValidation != nil {
		err := ob.personValidation(ob.document.person)
		if err != nil {
			return nil, errors.Wrap(err, "document person not valid")
		}
	}

	if ob.topicValidation != nil {
		err := ob.topicValidation(ob.document.topic)
		if err != nil {
			return nil, errors.Wrap(err, "document topic not valid")
		}
	}

	if ob.document.date.IsZero() {
		ob.document.date = time.Now()
	}

	// Создаем копию, чтобы не получать указатель на тот же объект,
	// если новых полей небыло добавлено.
	doc := ob.document
	return &doc, nil
}

func useBuilder() {
	applyingTmpl := NewDocumentBuilder().
		WithDepartment("accounting").
		WithPersonValidation(func(p Person) error {
			if p.Age < 18 {
				return errors.Errorf("person %+v is underage", p)
			}
			return nil
		}).
		WithTopicValidation(func(t Topic) error {
			if t != Applying {
				return errors.New("document subject topic is not applying")
			}
			return nil
		})

	_, err := applyingTmpl.Build()
	if err != nil {
		fmt.Println("error:", err)
	}

	_, err = applyingTmpl.
		WithPerson(Person{"Cave", "Johnson", 32}).
		WithDepartment("finance").
		WithTopic(Agreement).
		Build()
	if err != nil {
		fmt.Println("error:", err)
	}

	caroline, err := applyingTmpl.
		WithPerson(Person{"Caroline", "", 28}).
		WithTopic(Applying).
		Build()
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Caroline:\n%s\n", caroline)
}

// error: document person not valid: person {FirstName: LastName: Age:0} is underage
// error: document topic not valid: document subject topic is not applying
// Caroline:
// To accounting department
// From Caroline
// At 2009-11-10 23:00:00 +0000 UTC m=+0.000000001
// Subject: Applying
