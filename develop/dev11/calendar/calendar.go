package calendar

import (
	"fmt"
	"sort"
	"time"
)

type Event struct {
	Eid  int       `json:"eid"`
	Date time.Time `json:"date"`
	Msg  string    `json:"msg"`
}

func NewEvent(date time.Time, msg string) Event {
	return Event{Date: date, Msg: msg}
}

func (e *Event) Update(other Event) {
	if (time.Time{}) != other.Date {
		e.Date = other.Date
	}

	if "" != other.Msg {
		e.Msg = other.Msg
	}
}

func (e *Event) String() string {
	return fmt.Sprintf("(%d %s %v)", e.Eid, e.Date.GoString(), e.Msg)
}

func SortEvents(arr []Event) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Date.Before(arr[j].Date)
	})
}

type Calendar struct {
	storage map[int][]Event
	lastId  int
}

func NewCalendar() *Calendar {
	return &Calendar{
		storage: map[int][]Event{},
		lastId:  0,
	}
}

func (c *Calendar) Create(user int, event Event) error {
	c.lastId++
	event.Eid = c.lastId
	c.storage[user] = append(c.storage[user], event)
	SortEvents(c.storage[user])
	return nil
}

func (c *Calendar) Update(user int, eid int, event Event) error {
	for i, e := range c.storage[user] {
		if e.Eid == eid {
			c.storage[user][i].Update(event)
			SortEvents(c.storage[user])
			return nil
		}
	}

	return fmt.Errorf("not found")
}

func (c *Calendar) Delete(user int, eid int) error {
	for idx, e := range c.storage[user] {
		if e.Eid == eid {
			c.storage[user] = append(c.storage[user][:idx], c.storage[user][idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("not found")
}

type EventRange int64

const (
	All EventRange = iota
	MonthRange
	WeekRange
	DayRange
)

type EventQuery struct {
	User       *int
	Date       time.Time
	EventRange EventRange
}

func MonthFilter(pivot time.Time) func(time.Time) bool {
	y := pivot.Year()
	m := pivot.Month()
	return func(date time.Time) bool {
		return y == date.Year() && m == date.Month()
	}
}

func WeekFilter(pivot time.Time) func(time.Time) bool {
	y := pivot.Year()
	wd := pivot.YearDay() - int(pivot.Weekday())
	return func(date time.Time) bool {
		dateWd := date.YearDay() - int(date.Weekday())
		return y == date.Year() && wd == dateWd
	}
}

func DayFilter(pivot time.Time) func(time.Time) bool {
	y := pivot.Year()
	yd := pivot.YearDay()
	return func(date time.Time) bool {
		return y == date.Year() && yd == date.YearDay()
	}
}

func (c *Calendar) Query(q EventQuery) (result []Event) {
	result = []Event{}

	var userFilter func(user int) bool
	var dateFilter func(date time.Time) bool

	if q.User != nil {
		userFilter = func(user int) bool {
			return *q.User == user
		}
	}

	switch q.EventRange {
	case MonthRange:
		dateFilter = MonthFilter(q.Date)
	case WeekRange:
		dateFilter = WeekFilter(q.Date)
	case DayRange:
		dateFilter = DayFilter(q.Date)
	}

	users := []int{}
	for user := range c.storage {
		if userFilter != nil && !userFilter(user) {
			continue
		}
		users = append(users, user)
	}

	sort.Ints(users)

	for _, user := range users {
		events := c.storage[user]
		for _, event := range events {
			if dateFilter != nil && !dateFilter(event.Date) {
				continue
			}

			result = append(result, event)
		}
	}

	SortEvents(result)

	return
}
