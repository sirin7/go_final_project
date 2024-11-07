package handlers

import (
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/sirin7/go_final_project/database"
	"github.com/sirin7/go_final_project/helpers"
	"github.com/sirin7/go_final_project/models"
)

// Получаем все задачи из базы
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	tasks, err := database.GetAllTasks(h.DB)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch tasks"}`, http.StatusInternalServerError)
		log.Println("Error when extracting tasks:", err)
		return
	}

	// Устанавливаем заголовки и статус ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := models.TaskResponse{
		Tasks: tasks,
	}

	// Отправляем ответ клиенту
	helpers.EncodeJSON(w, response)

}
