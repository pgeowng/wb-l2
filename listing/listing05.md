Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Программа выведет error

Как и в 3 задаче, err = test() сохраняет внутрь интерфейса (*customError, nil),
который не равен nil. Только nil interface - interface{}(nil) равен nil

Возможным решением будет изменить на func test() error.
Тогда при return nil будет возвращаться (nil, nil), который равен nil.
```
