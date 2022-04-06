package pattern

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/pkg/errors"
)

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

// Паттерн Facade
// Заключается в объявлении класса, который объединяет несколько классов и использует их методы.
// Паттерн удобен, если нужно:
//   + Упростить взаимодействие, объединяя сложную логику обращений в себе.
//     Например, для получения данных используется cache, db, cloud service.
//     В зависимости от доступности и типа ресурса используется соответсующий компонент.
//   + Вынести зависимости и логику в отдельный модуль.
//     Таким образом, мы понижаем связность модулей между собой.
//     И объеденим дублирующийся код.
//     Например, у нас есть web, app, api клиенты.
//     В итоге они будут использовать унифицированный код при обращении к списку пользователей.
//   + Организовать взаимодействие внутри программы через слои.

// Неправильно использование может привести к тому,
// что фасад класс связан со всей системой, что затрудняет внесение изменений.

// Мы хотим получить все данные по профилю пользователя.
// Первый сервис отвечает за социальную состовляющую - фолловеры и т.п.
// Второй хранит историю прослушиваний.
type SocialData struct {
	followers int
	following int
}

type SocialSrv struct{}

func (srv *SocialSrv) GetUser(uid string) (*SocialData, error) {
	if rand.Intn(10) < 1 {
		return nil, errors.New("social service: internal error")
	}
	return &SocialData{followers: rand.Intn(100), following: rand.Intn(100)}, nil
}

type Track struct {
	id int
}

type HistorySrv struct{}

func (srv *HistorySrv) GetUser(uid string) []Track {
	size := rand.Intn(5) + 1
	result := make([]Track, 0, size)
	for i := 0; i < size; i++ {
		result = append(result, Track{rand.Intn(100000)})
	}
	return result
}

type ProfileFacade struct {
	social *SocialSrv
	hist   *HistorySrv
}

// Объявляем фасад, который подготовит объединенные данные.
// Инициализация дочерних классов внутри конструктора важна,
// если именно хотим вынести зависимости из классов, которые будут использовать facade.
func NewProfileFacade() (*ProfileFacade, error) {
	return &ProfileFacade{
		social: &SocialSrv{},
		hist:   &HistorySrv{},
	}, nil
}

type Profile struct {
	SocialData
	history []Track
}

// Метод объеденяющий сложную логику получения данных.
func (pf *ProfileFacade) GetProfile(uid string) (*Profile, error) {
	social, err := pf.social.GetUser(uid)
	if err != nil {
		return nil, err
	}

	history := pf.hist.GetUser(uid)
	return &Profile{*social, history}, nil
}

func useFacade() {
	pf, err := NewProfileFacade()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := pf.GetProfile("xhxwemt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", profile)
}

// &{SocialData:{followers:87 following:47} history:[{id:2081} {id:41318} {id:54425} {id:22540} {id:40456}]}
