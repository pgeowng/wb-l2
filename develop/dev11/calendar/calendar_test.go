package calendar

import (
	"fmt"
	"testing"
	"time"
)

func Failed(t *testing.T, format string, logf ...interface{}) {
	t.Logf(format, logf...)
	t.Fail()
}

func TQuery(c *Calendar, query EventQuery, expected []Event) error {
	result := c.Query(query)
	if len(result) != len(expected) {
		return fmt.Errorf("length mismatch:\nexpected: %v\ngot: %v", expected, result)
	}

	for idx, item := range expected {
		got := result[idx]
		if item.String() != got.String() {
			return fmt.Errorf("item mismatch: expected %s, got %s", item.String(), got.String())
		}
	}

	return nil
}

func TestCalendar(t *testing.T) {
	c := NewCalendar()

	var result []Event

	// query empty
	result = c.Query(EventQuery{})
	if len(result) != 0 {
		Failed(t, "empty calendar error")
		return
	}

	t1 := time.UnixMilli(1649201584000)
	err := c.Create(1, Event{Date: t1, Msg: "hello"})
	if err != nil {
		Failed(t, "failed at adding first entry: %v", err)
		return
	}

	err = TQuery(c, EventQuery{}, []Event{{Eid: 1, Date: t1, Msg: "hello"}})
	if err != nil {
		Failed(t, "failed at querying first: %v", err)
		return
	}

	// add second
	t2 := time.UnixMilli(1649892784000)
	err = c.Create(1, Event{Date: t2, Msg: "there"})
	if err != nil {
		Failed(t, "failed at adding second entry: %v", err)
		return
	}

	err = TQuery(c, EventQuery{}, []Event{{Eid: 1, Date: t1, Msg: "hello"}, {Eid: 2, Date: t2, Msg: "there"}})
	if err != nil {
		Failed(t, "failed at querying second: %v", err)
		return
	}

	// modify time
	t3 := time.UnixMilli(1650583984000)
	err = c.Update(1, 1, Event{Date: t3})
	if err != nil {
		Failed(t, "failed at modify time: %v", err)
		return
	}

	err = TQuery(c, EventQuery{}, []Event{{Eid: 2, Date: t2, Msg: "there"}, {Eid: 1, Date: t3, Msg: "hello"}})
	if err != nil {
		Failed(t, "failed at querying modified time: %v", err)
		return
	}

	// modify msg
	err = c.Update(1, 2, Event{Msg: "another"})
	if err != nil {
		Failed(t, "failed at modify Msg: %v", err)
		return
	}

	err = TQuery(c, EventQuery{}, []Event{{Eid: 2, Date: t2, Msg: "another"}, {Eid: 1, Date: t3, Msg: "hello"}})
	if err != nil {
		Failed(t, "failed at querying modified Msg: %v", err)
		return
	}

	// add to another user
	user := 2
	err = c.Create(2, Event{Date: t1, Msg: "message"})
	if err != nil {
		Failed(t, "failed at adding to another user: %v", err)
		return
	}

	err = TQuery(c, EventQuery{User: &user}, []Event{{Eid: 3, Date: t1, Msg: "message"}})
	if err != nil {
		Failed(t, "failed at querying another user: %v", err)
		return
	}
	// query week
	err = TQuery(c, EventQuery{Date: t2, EventRange: DayRange}, []Event{{Eid: 2, Date: t2, Msg: "another"}})
	if err != nil {
		Failed(t, "failed at querying day: %v", err)
		return
	}

	// query week
	err = TQuery(c, EventQuery{Date: t2, EventRange: WeekRange}, []Event{{Eid: 2, Date: t2, Msg: "another"}})
	if err != nil {
		Failed(t, "failed at querying week: %v", err)
		return
	}

	// query month
	err = TQuery(c, EventQuery{Date: t2, EventRange: MonthRange}, []Event{{3, t1, "message"}, {Eid: 2, Date: t2, Msg: "another"}, {1, t3, "hello"}})
	if err != nil {
		Failed(t, "failed at querying month: %v", err)
		return
	}

	// delete not found
	err = c.Delete(1, 5)
	if err == nil {
		Failed(t, "failed at delete not found: %v", err)
		return
	}

	// delete one
	err = c.Delete(1, 2)
	if err != nil {
		Failed(t, "failed at modify Msg: %v", err)
		return
	}

	err = TQuery(c, EventQuery{}, []Event{{3, t1, "message"}, {1, t3, "hello"}})
	if err != nil {
		Failed(t, "failed at querying delete one: %v", err)
		return
	}
}
