package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/pgeowng/wb-l2/develop/dev11/calendar"
	"github.com/pgeowng/wb-l2/develop/dev11/server"
)

// Я понимаю, что тесты написаны плохо, так как:
//  - Зависят от компоненов вне routes. - calendar, server
//  - Сильно связаны с результатом query. Любое изменение в Event{} приведет к поломке тестов.
//    Можно попробовать проверять наличие поля.
//    Но также при изменении Event{} сменится валидация полей,
//    необходимые поля в запросах, что опять же сломает тесты.

type RequestTest struct {
	handler func(server.Context)

	method  string
	path    string
	payload string

	statusCode int
	result     string
	err        bool
}

func (rt *RequestTest) Prepare() (req *http.Request, w *httptest.ResponseRecorder, err error) {
	path := rt.path
	var body io.Reader = nil
	contentLength := "0"

	if rt.method == "GET" {
		path += "?" + rt.payload
	} else {
		q, err := url.ParseQuery(rt.payload)
		if err != nil {
			return nil, nil, err
		}
		contentLength = strconv.Itoa(len(q.Encode()))
		body = strings.NewReader(q.Encode())
	}

	req = httptest.NewRequest(rt.method, path, body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", contentLength)

	if err := req.ParseForm(); err != nil {
		return nil, nil, err
	}

	w = httptest.NewRecorder()
	return
}

func (rt *RequestTest) Test(t *testing.T) {
	req, w, err := rt.Prepare()
	if err != nil {
		t.Logf("failed at prepare %v: %v", rt, err)
		t.Fail()
		return
	}

	rt.handler(server.Context{Req: req, Res: w})

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unable to parse body")
		return
	}

	if res.StatusCode != rt.statusCode {
		t.Logf("status code mismatch: expected %v, got %v", rt.statusCode, res.StatusCode)
		t.Log(rt)
		t.Log(string(data))
		t.Fail()
		return
	}

	var jsonBody map[string]interface{}

	err = json.Unmarshal(data, &jsonBody)
	if err != nil {
		t.Error("cant unmarshal response:", err)
		return
	}

	if rt.err {
		if _, ok := jsonBody["error"]; !ok {
			t.Log("expected err, got", jsonBody)
			t.Log(rt)
			t.Fail()
			return
		}
	} else {
		result, ok := jsonBody["result"]
		if !ok {
			t.Log("expected result, got", jsonBody)
			t.Log(rt)
			t.Fail()
			return
		}

		resultStr, err := json.Marshal(result)
		if err != nil {
			t.Error("cant marshal result", result)
			return
		}

		if string(resultStr) != rt.result {
			t.Errorf("result mismatch. expected %#v, got %#v", rt.result, string(resultStr))
			return
		}
	}
}

func QueryAll(result string, handler func(server.Context)) *RequestTest {
	return &RequestTest{
		handler,
		"GET",
		"/",
		"",

		http.StatusOK,
		result,
		false,
	}

}

func TestCreate(t *testing.T) {
	cal := calendar.NewCalendar()
	r := NewRoutes(cal)

	QueryAll("[]", r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"",

		http.StatusBadRequest,
		"",
		true,
	}).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1",

		http.StatusBadRequest,
		"",
		true,
	}).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"date=1",

		http.StatusBadRequest,
		"",
		true,
	}).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"date=2006-01-02T15:04:05Z",

		http.StatusBadRequest,
		"",
		true,
	}).Test(t)

	QueryAll("[]", r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2006-01-02T15:04:05Z",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2006-01-02T15:04:05Z","eid":1,"msg":""}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2004-01-02T15:04:05Z&msg=hello there",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2004-01-02T15:04:05Z","eid":2,"msg":"hello there"},{"date":"2006-01-02T15:04:05Z","eid":1,"msg":""}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2008-01-02T15:04:05Z&msg=hello there",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2004-01-02T15:04:05Z","eid":2,"msg":"hello there"},{"date":"2006-01-02T15:04:05Z","eid":1,"msg":""},{"date":"2008-01-02T15:04:05Z","eid":3,"msg":"hello there"}]`, r.QueryBuilder(calendar.All)).Test(t)
}

func TestDelete(t *testing.T) {
	cal := calendar.NewCalendar()
	r := NewRoutes(cal)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2008-01-02T15:04:05Z&msg=hello there",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2008-01-02T15:04:05Z","eid":1,"msg":"hello there"}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.DeleteEvent,

		"POST",
		"/delete_event",
		"",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	(&RequestTest{
		r.DeleteEvent,

		"POST",
		"/delete_event",
		"user=2",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	(&RequestTest{
		r.DeleteEvent,

		"POST",
		"/delete_event",
		"eid=-1",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	QueryAll(`[{"date":"2008-01-02T15:04:05Z","eid":1,"msg":"hello there"}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.DeleteEvent,

		"POST",
		"/delete_event",
		"eid=1&user=1",

		http.StatusOK,
		`"ok"`,
		false,
	}).Test(t)

	QueryAll(`[]`, r.QueryBuilder(calendar.All)).Test(t)
}

func TestUpdate(t *testing.T) {
	cal := calendar.NewCalendar()
	r := NewRoutes(cal)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2008-01-02T15:04:05Z&msg=hello there",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2008-01-02T15:04:05Z","eid":1,"msg":"hello there"}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"user=1",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"eid=1",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"user=1&eid=1",

		http.StatusBadRequest,
		``,
		true,
	}).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"user=1&eid=1&date=2009-01-02T15:04:05Z",

		http.StatusOK,
		`"ok"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2009-01-02T15:04:05Z","eid":1,"msg":"hello there"}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"user=1&eid=1&msg=another msg",

		http.StatusOK,
		`"ok"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2009-01-02T15:04:05Z","eid":1,"msg":"another msg"}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.UpdateEvent,

		"POST",
		"/update_event",
		"user=1&eid=1&msg=third msg&date=2010-01-02T15:04:05Z",

		http.StatusOK,
		`"ok"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2010-01-02T15:04:05Z","eid":1,"msg":"third msg"}]`, r.QueryBuilder(calendar.All)).Test(t)
}

func TestQuery(t *testing.T) {
	cal := calendar.NewCalendar()
	r := NewRoutes(cal)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2022-04-06T15:04:05Z&msg=first",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=2&date=2022-04-14T15:04:05Z&msg=second",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=2&date=2022-04-16T15:04:05Z&msg=third",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	(&RequestTest{
		r.CreateEvent,

		"POST",
		"/create_event",
		"user=1&date=2022-05-16T15:04:05Z&msg=fourth",

		http.StatusCreated,
		`"created"`,
		false,
	}).Test(t)

	QueryAll(`[{"date":"2022-04-06T15:04:05Z","eid":1,"msg":"first"},{"date":"2022-04-14T15:04:05Z","eid":2,"msg":"second"},{"date":"2022-04-16T15:04:05Z","eid":3,"msg":"third"},{"date":"2022-05-16T15:04:05Z","eid":4,"msg":"fourth"}]`, r.QueryBuilder(calendar.All)).Test(t)

	(&RequestTest{
		r.QueryBuilder(calendar.MonthRange),

		"GET",
		"/events_for_month",
		"date=2022-04-15T10:00:00Z",

		http.StatusOK,
		`[{"date":"2022-04-06T15:04:05Z","eid":1,"msg":"first"},{"date":"2022-04-14T15:04:05Z","eid":2,"msg":"second"},{"date":"2022-04-16T15:04:05Z","eid":3,"msg":"third"}]`,
		false,
	}).Test(t)

	(&RequestTest{
		r.QueryBuilder(calendar.WeekRange),

		"GET",
		"/events_for_week",
		"date=2022-04-13T10:00:00Z",

		http.StatusOK,
		`[{"date":"2022-04-14T15:04:05Z","eid":2,"msg":"second"},{"date":"2022-04-16T15:04:05Z","eid":3,"msg":"third"}]`,
		false,
	}).Test(t)

	(&RequestTest{
		r.QueryBuilder(calendar.DayRange),

		"GET",
		"/events_for_day",
		"date=2022-04-14T10:00:00Z",

		http.StatusOK,
		`[{"date":"2022-04-14T15:04:05Z","eid":2,"msg":"second"}]`,
		false,
	}).Test(t)

	(&RequestTest{
		r.QueryBuilder(calendar.All),

		"GET",
		"/",
		"user=1",

		http.StatusOK,
		`[{"date":"2022-04-06T15:04:05Z","eid":1,"msg":"first"},{"date":"2022-05-16T15:04:05Z","eid":4,"msg":"fourth"}]`,
		false,
	}).Test(t)
}
