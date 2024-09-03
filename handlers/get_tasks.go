package handlers

import (
	"log"
	"net/http"

	"github.com/sirin7/go_final_project/models"
)

// Получаем все задачи из базы
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	tasks, err := models.GetAllTasks(h.DB)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch tasks"}`, http.StatusInternalServerError)
		log.Println("Ошибка при извлечении задач:", err)
		return
	}

	// Устанавливаем заголовки и статус ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := models.TaskResponse{
		Tasks: tasks,
	}

	// Отправляем ответ клиенту
	models.EncodeJSON(w, response)
	log.Println("JSON успешно сериализован и отправлен клиенту")
}
