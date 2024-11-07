package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/sirin7/go_final_project/database"
	"github.com/sirin7/go_final_project/helpers"
)

// Получение задачи по ID
func (h *Handler) GetTaskId(w http.ResponseWriter, r *http.Request) {
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

	task, err := database.GetTaskById(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "task not found"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "failed to retrieve task: %v"}`, err)))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	helpers.EncodeJSON(w, task)

}
