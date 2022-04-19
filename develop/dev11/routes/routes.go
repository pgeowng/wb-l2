package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pgeowng/wb-l2/develop/dev11/calendar"
	"github.com/pgeowng/wb-l2/develop/dev11/server"
)

func ValidatePositiveInt(value string) (pint int, err error) {
	pint64, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return
	}

	if pint64 < 1 {
		err = fmt.Errorf("user id is positive integer")
		return
	}

	pint = int(pint64)
	return
}

func ValidateDate(value string) (date time.Time, err error) {
	date, err = time.Parse(time.RFC3339, value)
	if err != nil {
		return
	}

	if date.Year() < 2000 {
		err = fmt.Errorf("date can't be before 2000 January 1")
		return
	}

	return
}

type Routes struct {
	cal *calendar.Calendar
}

func NewRoutes(cal *calendar.Calendar) *Routes {
	return &Routes{cal: cal}
}

func (r *Routes) CreateEvent(ctx server.Context) {
	user, err := ValidatePositiveInt(ctx.Req.PostForm.Get("user"))
	if err != nil {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": fmt.Sprint("user field:", err),
		})
		return
	}

	date, err := ValidateDate(ctx.Req.PostForm.Get("date"))
	if err != nil {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": fmt.Sprint("date field:", err),
		})
		return
	}

	msg := ctx.Req.PostForm.Get("msg")

	err = r.cal.Create(user, calendar.NewEvent(date, msg))
	if err != nil {
		ctx.SendJSON(http.StatusServiceUnavailable, server.H{
			"error": fmt.Sprint("create:", err),
		})
		return
	}

	ctx.SendJSON(http.StatusCreated, server.H{
		"result": "created",
	})

}

func (r *Routes) UpdateEvent(ctx server.Context) {
	var err error

	user, err := ValidatePositiveInt(ctx.Req.PostForm.Get("user"))
	if err != nil {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": fmt.Sprint("user field:", err),
		})
		return
	}

	eid, err := ValidatePositiveInt(ctx.Req.PostForm.Get("eid"))
	if err != nil {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": fmt.Sprint("eid field:", err),
		})
		return
	}

	var date time.Time
	hasDate := false
	dateField := ctx.Req.PostForm.Get("date")
	if len(dateField) > 0 {
		hasDate = true

		date, err = ValidateDate(ctx.Req.PostForm.Get("date"))
		if err != nil {
			ctx.SendJSON(http.StatusBadRequest, server.H{
				"error": fmt.Sprint("date field:", err),
			})
			return
		}
	}

	msg := ctx.Req.PostForm.Get("msg")

	if len(msg) == 0 && !hasDate {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": "empty update request",
		})
		return
	}

	err = r.cal.Update(user, eid, calendar.NewEvent(date, msg))
	if err != nil {
		ctx.SendJSON(http.StatusServiceUnavailable, server.H{
			"error": fmt.Sprint("update:", err),
		})
		return
	}

	ctx.SendJSON(http.StatusOK, server.H{
		"result": "ok",
	})
}
func (r *Routes) DeleteEvent(ctx server.Context) {
	var err error

	user, err := ValidatePositiveInt(ctx.Req.PostForm.Get("user"))
	if err != nil {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": fmt.Sprint("user field:", err),
		})
		return
	}

	eid, err := ValidatePositiveInt(ctx.Req.PostForm.Get("eid"))
	if err != nil {
		ctx.SendJSON(http.StatusBadRequest, server.H{
			"error": fmt.Sprint("eid field:", err),
		})
		return
	}

	err = r.cal.Delete(user, eid)
	if err != nil {
		ctx.SendJSON(http.StatusServiceUnavailable, server.H{
			"error": fmt.Sprint("delete:", err),
		})
		return
	}

	ctx.SendJSON(http.StatusOK, server.H{
		"result": "ok",
	})
}

func (r *Routes) QueryBuilder(erange calendar.EventRange) func(ctx server.Context) {
	return func(ctx server.Context) {
		var err error

		var user int
		hasUser := false
		userField := ctx.Req.Form.Get("user")
		if len(userField) > 0 {
			hasUser = true

			user, err = ValidatePositiveInt(userField)
			if err != nil {
				ctx.SendJSON(http.StatusBadRequest, server.H{
					"error": fmt.Sprint("user field:", err),
				})
				return
			}
		}

		var date time.Time
		if erange != calendar.All {
			date, err = ValidateDate(ctx.Req.Form.Get("date"))
			if err != nil {
				ctx.SendJSON(http.StatusBadRequest, server.H{
					"error": fmt.Sprint("date field:", err),
				})
				return
			}
		}

		eq := calendar.EventQuery{Date: date, EventRange: erange}

		if hasUser {
			eq.User = &user
		}

		result := r.cal.Query(eq)
		ctx.SendJSON(http.StatusOK, server.H{
			"result": result,
		})
	}

}
