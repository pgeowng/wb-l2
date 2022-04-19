package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/pgeowng/wb-l2/develop/dev11/calendar"
	"github.com/pgeowng/wb-l2/develop/dev11/routes"
	"github.com/pgeowng/wb-l2/develop/dev11/server"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func main() {
	PORT := os.Getenv("PORT")

	if port, err := strconv.ParseInt(PORT, 10, 0); err != nil || port > 65535 || port < 0 {
		fmt.Println("srv: bad PORT value:", PORT)
		os.Exit(1)
	}

	srv := server.New(":" + PORT)

	cal := calendar.NewCalendar()
	routes := routes.NewRoutes(cal)

	logger := server.LoggerMW

	srv.Get("/", routes.QueryBuilder(calendar.All), logger)
	srv.Get("/events_for_day", routes.QueryBuilder(calendar.DayRange), logger)
	srv.Get("/events_for_week", routes.QueryBuilder(calendar.WeekRange), logger)
	srv.Get("/events_for_month", routes.QueryBuilder(calendar.MonthRange), logger)

	srv.Post("/create_event", routes.CreateEvent, logger)
	srv.Post("/update_event", routes.UpdateEvent, logger)
	srv.Post("/delete_event", routes.DeleteEvent, logger)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		signal.Stop(interrupt)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	if err := srv.Listen(); err != nil {
		fmt.Println(err)
	}
}
