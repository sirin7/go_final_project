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
)

func main() {

	//Открываем базу scheduler.db, если ее нет, то создаем ее
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	//Проверяем существует ли в директории приложения файл scheduler.db
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
	//Если install равен true, после открытия БД выполняем SQL запрос в файле create_table.sql
	if install {
		fileContent, err := os.ReadFile(constants.CreateTable)
		if err != nil {
			log.Println("Не удалось прочитать файл create_table.sql")
			log.Fatal(err)
		}

		_, err = db.Exec(string(fileContent))
		log.Println("Таблица успешно создана")
		if err != nil {
			log.Println("Не удалось выполнить файл create_table.sql")
			log.Fatal(err)
		}
	}

	// Создаем новый обработчик, используя подключение к базе данных
	handler := handlers.NewHandler(db)

	// Создаем новый маршрутизатор (router) с помощью chi
	r := chi.NewRouter()

	// Создаем файловый сервер, который будет обслуживать статические файлы из директории, указанной в constants.WebDir
	fs := http.FileServer(http.Dir(constants.WebDir))

	// Получение задачи по ID
	r.Get("/api/task", handler.GetTaskId)

	// Получение следующей даты для задачи на основе повторения
	r.Get("/api/nextdate", handler.GetDateTask)

	// Получение списка всех задач
	r.Get("/api/tasks", handler.GetTasks)

	// Создание новой задачи
	r.Post("/api/task", handler.PostTask)

	// Отметка задачи как выполненной (с возможностью переноса на следующую дату)
	r.Post("/api/task/done", handler.DoneTask)

	// Обновление существующей задачи
	r.Put("/api/task", handler.PutTask)

	// Удаление задачи
	r.Delete("/api/task", handler.DeleteTask)

	// Обслуживание статических файлов для всех маршрутов, которые не совпадают с вышеуказанными
	r.Handle("/*", fs)

	// Запускаем сервер на порту 7540 с использованием настроенного маршрутизатора
	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("ошибка запуска сервера: %s\n", err.Error())
		return
	}

	// Логируем успешный запуск сервера
	log.Println("Сервер запущен на порту 7540")

}
