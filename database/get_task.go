package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/models"
	_ "modernc.org/sqlite"
)

func GetTaskById(db *sql.DB, id int) (*models.Task, error) {
	var task models.Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Printf("Could not find task with ID: %d", id)
		return nil, err
	}

	return &task, err
}

func GetAllTasks(db *sql.DB) ([]models.Task, error) {
	//Получаем задачи из таблицы
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?", constants.Limit)
	if err != nil {
		log.Println("failed to get tasks", err)
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf(`{"error": "failed to scan row: %v"}`, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(`{"error": "rows iteration error: %v"}`, err)
	}

	return tasks, nil
}
