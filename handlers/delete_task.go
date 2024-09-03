package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sirin7/go_final_project/models"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	//Получаем ID
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
	// Проверяем существует ли задача в базе
	taskFromDB, err := models.GetTaskById(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}
	if taskFromDB == nil {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}
	// Удаляем задачу из базы
	deleteTask, err := models.DoneTask(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "failed to delete task"}`, http.StatusBadRequest)
		return
	}
	log.Println("Задача выполнена", deleteTask)

	// Формируем ответ клиенту
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(`{}`)))
	return

}
