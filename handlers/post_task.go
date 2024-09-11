package handlers

import (
	"fmt"

	"net/http"

	_ "modernc.org/sqlite"

	//"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/database"
	"github.com/sirin7/go_final_project/helpers"
	"github.com/sirin7/go_final_project/models"
)

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	//добавлениу задачи в базу данных
	var task models.Task
	if err := helpers.DecodeJSON(r.Body, &task); err != nil {
		http.Error(w, `{"error": "failed to unmarshal JSON"}`, http.StatusBadRequest)
		return
	}

	if err := helpers.CheckTask(&task); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := database.AddTask(h.DB, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}
