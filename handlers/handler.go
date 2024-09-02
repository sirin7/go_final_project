package handlers

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/helpers"
	"github.com/sirin7/go_final_project/models"
)

// GetTask - структура для работы с задачами из базы данных
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
		//http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	tasks, err := models.FetchTasks(h.DB)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch tasks"}`, http.StatusInternalServerError)
		log.Println("Ошибка при извлечении задач:", err)
		return
	}

	// Устанавливаем заголовки и статус ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// Отправляем ответ клиенту

	response := models.TaskResponse{
		Tasks: tasks,
	}

	models.EncodeJSON(w, response)
	log.Println("JSON успешно сериализован и отправлен клиенту")
}

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	//добавлениу задачи в базу данных
	var task models.Task
	if err := models.DecodeJSON(r.Body, &task); err != nil {
		http.Error(w, `{"error": "failed to unmarshal JSON"}`, http.StatusBadRequest)
		log.Println("Не удалось десериализовать JSON")
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "missing title"}`, http.StatusBadRequest)
		log.Println("Заголовок задачи не может быть пустым")
		return
	}

	log.Println("задача", task.Title)

	// Проверка корректности даты

	if task.Date == "" {
		task.Date = time.Now().Format(constants.FormatDate)
	}

	checkDate, err := time.Parse(constants.FormatDate, task.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
		return
	}
	log.Println("Корректная дата", checkDate)

	if task.Date < time.Now().Format(constants.FormatDate) {
		if task.Repeat != "" {
			nextDate, err := helpers.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Failed to calculate next date"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(constants.FormatDate)
		}
	}

	// Добавление задачи в базу данных
	id, err := models.AddTask(h.DB, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func (h *Handler) GetDateTask(w http.ResponseWriter, r *http.Request) {
	nowstr := r.URL.Query().Get("now")
	if nowstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "now date missing"}`))
		return
	}

	now, err := time.Parse(constants.FormatDate, nowstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid now date format"}`))
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "date missing"}`))
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "repeat missing"}`))
		return
	}

	nextDate, err := helpers.NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	// Возвращаем следующую дату
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s", nextDate)))
}

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

func (h *Handler) PutTask(w http.ResponseWriter, r *http.Request) {
	//обнавляем задачу и сохраняем в базу данных
	var task models.Task
	if err := models.DecodeJSON(r.Body, &task); err != nil {
		http.Error(w, `{"error": "failed to unmarshal JSON"}`, http.StatusBadRequest)
		log.Println("Не удалось десериализовать JSON")
		return
	}

	//// Преобразуем ID из строки в число
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error": "invalid id format"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли задача с данным ID
	taskFromDB, err := models.GetTaskById(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}
	if taskFromDB == nil {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "missing title"}`, http.StatusBadRequest)
		log.Println("Заголовок задачи не может быть пустым")
		return
	}

	log.Println("задача", task.Title)

	// Проверка корректности даты

	if task.Date == "" {
		task.Date = time.Now().Format(constants.FormatDate)
	}

	checkDate, err := time.Parse(constants.FormatDate, task.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
		return
	}
	log.Println("Корректная дата", checkDate)

	if task.Date < time.Now().Format(constants.FormatDate) {
		if task.Repeat != "" {
			nextDate, err := helpers.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Failed to calculate next date"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(constants.FormatDate)
		}
	}

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

func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
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

	// Если Repeat задано, вычисляем следующую дату
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

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {

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

	taskFromDB, err := models.GetTaskById(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}
	if taskFromDB == nil {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	deleteTask, err := models.DoneTask(h.DB, id)
	if err != nil {
		http.Error(w, `{"error": "failed to delete task"}`, http.StatusBadRequest)
		return
	}
	log.Println("Задача выполнена", deleteTask)

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(`{}`)))
	return

}
