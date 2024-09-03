package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/helpers"
	"github.com/sirin7/go_final_project/models"
)

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	//добавлениу задачи в базу данных
	var task models.Task
	if err := models.DecodeJSON(r.Body, &task); err != nil {
		http.Error(w, `{"error": "failed to unmarshal JSON"}`, http.StatusBadRequest)
		log.Println("Не удалось десериализовать JSON")
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "missing title"}`, http.StatusBadRequest)
		log.Println("Заголовок задачи не может быть пустым")
		return
	}

	log.Println("задача", task.Title)

	// Проверка корректности даты

	if task.Date == "" {
		task.Date = time.Now().Format(constants.FormatDate)
	}

	checkDate, err := time.Parse(constants.FormatDate, task.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
		return
	}
	log.Println("Корректная дата", checkDate)

	if task.Date < time.Now().Format(constants.FormatDate) {
		if task.Repeat != "" {
			nextDate, err := helpers.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Failed to calculate next date"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(constants.FormatDate)
		}
	}

	// Добавление задачи в базу данных
	id, err := models.AddTask(h.DB, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}
