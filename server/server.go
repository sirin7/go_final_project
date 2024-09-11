package server

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/database"
	"github.com/sirin7/go_final_project/handlers"
)

func StartServer() {

	// Используем подключение к базе данных
	Handler := handlers.NewHandler(database.Db) // Передаем глобальную переменную Db

	// Создаем новый маршрутизатор с помощью chi
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(constants.WebDir))

	// Маршруты
	r.Get("/api/task", Handler.GetTaskId)
	r.Get("/api/nextdate", Handler.GetDateTask)
	r.Get("/api/tasks", Handler.GetTasks)
	r.Post("/api/task", Handler.PostTask)
	r.Post("/api/task/done", Handler.DoneTask)
	r.Put("/api/task", Handler.PutTask)
	r.Delete("/api/task", Handler.DeleteTask)

	r.Handle("/*", fs)

	// Запуск сервера
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	log.Println("The server is running on port 7540")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("server startup error: %s\n", err.Error())
		return
	}
}
