package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirin7/go_final_project/helpers"
	"github.com/sirin7/go_final_project/models"
)

// Отмечаем задачу как выполненую при пустом правиле повторения, если правило повторения есть высчитываем и переносим на следующую дату
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	// Получем ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "id missing"}`))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid id format"}`))
		return
	}
	//Проверяем есть ли задача в базе
	taskFromDB, err := models.GetTaskById(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}
	if taskFromDB == nil {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	task := *taskFromDB
	// Если правило повтороения отсутствует, удаляем задачу из базы
	if task.Repeat == "" {
		doneTask, err := models.DoneTask(h.DB, id)
		if err != nil {
			http.Error(w, `{"error": "failed to delete task"}`, http.StatusBadRequest)
			return
		}
		log.Println("Задача выполнена", doneTask)

		// Формирование ответа
		w.Header().Set("Content-Type", "application/json, charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint(`{}`)))
		return
	}

	// Если правило задано, вычисляем следующую дату
	nextDate, err := helpers.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		http.Error(w, `{"error": "Failed to calculate next date"}`, http.StatusBadRequest)
		return
	}
	task.Date = nextDate

	updatedTaskID, err := models.PutTask(h.DB, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Задача обновлена", updatedTaskID)

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprint(`{}`)))

}