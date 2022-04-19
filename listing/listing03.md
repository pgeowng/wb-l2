Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
Программа выведет nil и false.

Проверка показывает, что nil != nil.
Интерфейсы устроены как структура c внутренним типом и значением (T, V)
Исторически сложилось, что интерфейс error c маленькой буквы.
err := Foo() Сохраняет (*os.PathError, nil)

Но только nil interface равен nil. interface{}(nil) == (nil, nil)
Поэтому проверка на nil не сработает.

```
