// models/models.go
package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

func FetchTasks(db *sql.DB) ([]Task, error) {
	//Получаем задачи из таблицы
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf(`{"error": "failed to scan row: %v"}`, err)
		}
		tasks = append(tasks, task)
	}
	log.Println("Получены задачи из таблицы scheduler")
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(`{"error": "rows iteration error: %v"}`, err)
	}

	return tasks, nil
}

func GetTaskById(db *sql.DB, id int) (*Task, error) {
	var task Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	log.Println("ID задачи получен из scheduler")
	return &task, err
}

func AddTask(db *sql.DB, task Task) (int64, error) {
	// Добавляем задачи в таблицу
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	// Получаем ID последней добавленной задачи
	Id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Println("Задача добавлена в базу", task.Title)
	return Id, nil
}

func DoneTask(db *sql.DB, id int) (*Task, error) {
	_, err := db.Exec(("DELETE FROM scheduler WHERE id = ?"), id)
	if err != nil {
		return nil, err
	}

	log.Println("задача выполнена", id)
	return &Task{}, err
}

func PutTask(db *sql.DB, task Task) (Task, error) {

	// Выполнение SQL-запроса на обновление задачи
	_, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return Task{}, err
	}

	log.Printf("Задача %s обновлена: %s", task.ID, task.Title)

	// Возвращаем обновленную задачу и nil (без ошибки)
	return task, nil

}

// Десериализация и сериализация JSON
func DecodeJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}

func EncodeJSON(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)

}
