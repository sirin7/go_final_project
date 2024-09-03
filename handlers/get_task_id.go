package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sirin7/go_final_project/models"
)

// Получение задачи по ID
func (h *Handler) GetTaskId(w http.ResponseWriter, r *http.Request) {
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

	task, err := models.GetTaskById(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "task not found"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "failed to retrieve task: %v"}`, err)))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	models.EncodeJSON(w, task)
	log.Println("JSON успешно сериализован и отправлен клиенту")
}
