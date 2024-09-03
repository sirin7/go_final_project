package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sirin7/go_final_project/constants"
)

// NextDate вычисляет следующую дату выполнения задачи на основе текущей даты, заданной даты и правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим переданную дату в формате YYYYMMDD
	parsedDate, err := time.Parse(constants.FormatDate, date)
	if err != nil {

		return "", fmt.Errorf("неправильный формат даты: %v", err)
	}

	parts := strings.Split(repeat, " ")
	log.Println("проверка правила", parts)

	// Проверяем тип правила повторения
	switch parts[0] {
	case "y":
		// Ежегодное повторение
		if len(parts) > 1 {

			return "", fmt.Errorf("неподдерживаемый формат повторения")
		}

		parsedDate = parsedDate.AddDate(1, 0, 0)
		for !parsedDate.After(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}
	case "d":
		// Повторение через определенное количество дней
		if len(parts) > 1 {

			days, err := strconv.Atoi(parts[1])
			if err != nil || days <= 0 || days > 400 {

				return "", fmt.Errorf("неправильный формат для дней: %v", err)
			}

			parsedDate = parsedDate.AddDate(0, 0, days)
			for !parsedDate.After(now) {
				parsedDate = parsedDate.AddDate(0, 0, days)
			}
		} else {

			return "", fmt.Errorf("не указано количество дней")
		}
	case "m", "w":
		// Неподдерживаемые форматы повторения (ежемесячное и еженедельное)
		return "", fmt.Errorf("неподдерживаемый формат повторения")
	default:

		return "", fmt.Errorf("неподдерживаемое правило повторения")
	}

	log.Println("Следующая дата", parsedDate.Format(constants.FormatDate))

	// Возвращаем следующую дату в формате YYYYMMDD
	return parsedDate.Format(constants.FormatDate), nil
}

// Десериализация и сериализация JSON
func DecodeJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}

func EncodeJSON(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)

}