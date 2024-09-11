package helpers

import (
	"fmt"
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
		return "", fmt.Errorf("wrong date format: %v", err)
	}

	parts := strings.Split(repeat, " ")

	// Проверяем тип правила повторения
	switch parts[0] {
	case "y":
		// Ежегодное повторение
		if len(parts) > 1 {
			return "", fmt.Errorf("unsupported repeat format")
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
				return "", fmt.Errorf("wrong format for days: %v", err)
			}

			parsedDate = parsedDate.AddDate(0, 0, days)
			for !parsedDate.After(now) {
				parsedDate = parsedDate.AddDate(0, 0, days)
			}
		} else {
			return "", fmt.Errorf("number of days not specified")
		}
	case "m", "w":
		// Неподдерживаемые форматы повторения (ежемесячное и еженедельное)
		return "", fmt.Errorf("unsupported repeat format")
	default:
		return "", fmt.Errorf("unsupported repeat format")
	}

	// Возвращаем следующую дату в формате YYYYMMDD
	return parsedDate.Format(constants.FormatDate), nil
}
