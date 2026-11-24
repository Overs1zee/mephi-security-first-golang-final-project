package actioninfo

import (
	"fmt"
	"log"
)

// DataParser реализуют типы, которые умеют разбирать строку с данными
// и возвращать по ней готовую информацию.
type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

// Info разбирает каждую строку dataset с помощью dp и выводит информацию на экран.
// Строки, которые не удалось разобрать или обработать, пропускаются с записью ошибки в лог.
func Info(dataset []string, dp DataParser) {
	for _, v := range dataset {
		if err := dp.Parse(v); err != nil {
			log.Println(err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(info)
	}
}
