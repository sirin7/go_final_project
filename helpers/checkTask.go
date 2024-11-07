package helpers

import (
	"fmt"
	"time"

	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/models"
	_ "modernc.org/sqlite"
)

// Проверяем данные задачи
func CheckTask(task *models.Task) error {

	if task.Title == "" {
		return fmt.Errorf("missing title")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(constants.FormatDate)
	}

	_, err := time.Parse(constants.FormatDate, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYYMMDD")
	}

	if task.Date < time.Now().Format(constants.FormatDate) {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("failed to calculate next date")
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(constants.FormatDate)
		}
	}

	return nil
}
