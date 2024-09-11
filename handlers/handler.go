package handlers

import (
	"database/sql"

	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.PostTask(w, r)
		h.DoneTask(w, r)
	case http.MethodGet:
		h.GetTasks(w, r)
		h.GetDateTask(w, r)
		h.GetTaskId(w, r)
	case http.MethodPut:
		h.PutTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		log.Println(r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
