package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/sirin7/go_final_project/constants"
)

var Db *sql.DB

func InitDB() {

	dbFile := os.Getenv("TODO_DBFILE")
	// Если переменная не установлена или пуста, используем путь по умолчанию
	if dbFile == "" {
		dbFile = filepath.Join("./scheduler.db")
	}

	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Database file not found, new file created:", dbFile)
			install = true
		} else {
			log.Fatal("Error while checking the database file:", err)
		}
	}

	// Открываем базу данных (или создаем новый файл)
	Db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	//Если install равен true, создаем таблицу
	if install {
		fileContent, err := os.ReadFile(constants.CreateTable)
		if err != nil {
			log.Fatal("Failed to read file create_table.sql", err)
		}

		_, err = Db.Exec(string(fileContent))
		log.Println("Table created successfully")
		if err != nil {
			log.Fatal("Failed to execute file create_table.sql", err)
		}
	}
}
