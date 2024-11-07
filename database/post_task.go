package database

import (
	"database/sql"
	"log"

	"github.com/sirin7/go_final_project/models"
	_ "modernc.org/sqlite"
)

func AddTask(db *sql.DB, task models.Task) (int64, error) {
	// Добавляем задачи в таблицу
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("failed to add task with title: %v", task.Title)
		return 0, err
	}

	// Получаем ID последней добавленной задачи
	Id, err := res.LastInsertId()
	if err != nil {
		log.Println("failed to get the ID of the added task", err)
		return 0, err
	}
	return Id, nil
}
