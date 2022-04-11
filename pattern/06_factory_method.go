package pattern

import (
	"fmt"

	"github.com/pkg/errors"
)

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

// Factory Method - creational pattern
// Заключается в создании интерфейса с методом для создания объекта,
// и классов, которые реализуют интерфейс и решают как создавать классы.

// Creator Interface
type Factory interface {
	Create(id string) (Resource, error)
}

// Product Interface
type Resource interface {
	Link() string
}

// Concrete Product
type VideoClip struct {
	id       string
	width    uint
	height   uint
	duration uint
}

func (res *VideoClip) Link() string {
	fmt.Printf("calling ffmpeg for %s with %dx%d and duration %ds \n", res.id, res.width, res.height, res.duration)
	return "vc:" + res.id
}

type VideoThumb struct {
	id        string
	width     uint
	height    uint
	frameTime uint
}

func (res *VideoThumb) Link() string {
	fmt.Printf("calling ffmpeg for %s with %dx%d at %02d:%02ds \n", res.id, res.width, res.height, res.frameTime/60, res.frameTime%60)
	return "vt:" + res.id
}

type Image struct {
	id     string
	width  uint
	height uint
}

func (res *Image) Link() string {
	fmt.Printf("calling imagemagick for %s with %dx%d\n", res.id, res.width, res.height)
	return "it:" + res.id
}

type Original struct {
	id string
}

func (res *Original) Link() string {
	fmt.Printf("using original for %s\n", res.id)
	return res.id
}

// Concrete Creator

type LowQualityMedia struct{ MediaStorage }

func (f *LowQualityMedia) Create(id string) (Resource, error) {
	media, err := f.Get(id)
	if err != nil {
		return nil, err
	}

	switch media.kind {
	case "image":
		return &Image{id: id, width: 200, height: 200}, nil
	case "video":
		return &VideoThumb{id: id, width: 200, height: 200}, nil
	default:
		return nil, errors.Errorf("for %s kind %s is not implemented", id, media.kind)
	}
}

type MediumQualityMedia struct{ MediaStorage }

func (q *MediumQualityMedia) Create(id string) (Resource, error) {
	media, err := q.Get(id)
	if err != nil {
		return nil, err
	}

	switch media.kind {
	case "image":
		return &Image{id: id, width: 600, height: 600}, nil
	case "video":
		return &VideoClip{id: id, width: 600, height: 600, duration: 5}, nil
	default:
		return nil, errors.Errorf("for %s kind %s is not implemented", id, media.kind)
	}
}

type HighQualityMedia struct{ MediaStorage }

func (q *HighQualityMedia) Create(id string) (Resource, error) {
	_, err := q.Get(id)
	if err != nil {
		return nil, err
	}

	return &Original{id: id}, nil
}

type Media struct {
	id   string
	kind string
}

type MediaStorage struct {
	cache map[string]Media
}

func (s *MediaStorage) Get(id string) (media *Media, err error) {
	m, ok := s.cache[id]
	if !ok {
		return nil, errors.Errorf("media %s not found", id)
	}
	return &m, nil
}

func access(f Factory, id string) {
	fmt.Printf("getting %s...\n", id)
	res, err := f.Create(id)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		link := res.Link()
		fmt.Printf("link: %s\n", link)
	}
}

// У нас есть некоторые медиа.
// Один из них хотим загрузить в лучшем возможном качестве.
// Остальные в качестве оптимизации получаем в худшем качестве.
func useFactoryMethod() {
	storage := MediaStorage{cache: map[string]Media{
		"0": {"0", "video"}, "1": {"1", "image"}, "2": {"2", "audio"},
	}}

	hq := &HighQualityMedia{storage}
	mq := &MediumQualityMedia{storage}
	lq := &LowQualityMedia{storage}

	access(hq, "1")
	access(mq, "0")
	access(lq, "0")
	access(hq, "5")
}

// getting 1...
// using original for 1
// link: 1
// getting 0...
// calling ffmpeg for 0 with 600x600 and duration 5s
// link: vc:0
// getting 0...
// calling ffmpeg for 0 with 200x200 at 00:00s
// link: vt:0
// getting 5...
// error: media 5 not found
