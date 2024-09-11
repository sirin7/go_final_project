package handlers

import (
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/sirin7/go_final_project/database"
	"github.com/sirin7/go_final_project/helpers"
	"github.com/sirin7/go_final_project/models"
)

func (h *Handler) PutTask(w http.ResponseWriter, r *http.Request) {
	//обнавляем задачу и сохраняем в базу данных
	var task models.Task
	if err := helpers.DecodeJSON(r.Body, &task); err != nil {
		http.Error(w, `{"error": "failed to unmarshal JSON"}`, http.StatusBadRequest)
		return
	}

	if err := helpers.CheckTask(&task); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	rowsAffected, err := database.PutTask(h.DB, task)
	if err != nil {
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}

	// Проверяем, было ли затронуто хотя бы 1 запись
	if rowsAffected == 0 {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprint(`{}`)))
}
