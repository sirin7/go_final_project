package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"

	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/handlers"
	//"github.com/sirin7/go_final_project/database"
)

func main() {

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("База данных будет создана")
			install = true
		} else {
			log.Println("Не получилось проверить файл")
			log.Fatal(err)
		}
	}
	log.Println("База данных была создана ранее")
	if install {
		fileContent, err := os.ReadFile(constants.CreateTable)
		if err != nil {
			log.Println("Не удалось прочитать файл create_table.sql")
			log.Fatal(err)
		}

		// Выполняем содержимое файла как SQL-запрос
		_, err = db.Exec(string(fileContent))
		log.Println("Таблица успешно создана")
		if err != nil {
			log.Println("Не удалось выполнить файл create_table.sql")
			log.Fatal(err)
		}
	}
	handler := handlers.NewHandler(db)

	r := chi.NewRouter()
	fs := http.FileServer(http.Dir(constants.WebDir))

	r.Get("/api/task", handler.GetTaskId)
	r.Get("/api/nextdate", handler.GetDateTask)
	r.Get("/api/tasks", handler.GetTasks)
	r.Post("/api/task", handler.PostTask)
	r.Post("/api/task/done", handler.DoneTask)
	r.Put("/api/task", handler.PutTask)
	r.Delete("/api/task", handler.DeleteTask)
	r.Handle("/*", fs)
	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("ошибка запуска сервера: %s\n", err.Error())
		return
	}
	log.Println("Сервер запущен на порту 7540")
}
