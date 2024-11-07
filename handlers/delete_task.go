package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirin7/go_final_project/database"
	_ "modernc.org/sqlite"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	//Получаем ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "id missing"}`))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "invalid id format"}`))
		return
	}

	// Удаляем задачу из базы
	_, err = database.DoneTask(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "failed to delete task"}`, http.StatusBadRequest)
		return
	}

	// Формируем ответ клиенту
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprint(`{}`)))
	return

}
