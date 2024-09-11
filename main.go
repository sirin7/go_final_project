package main

import (
	"log"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"github.com/sirin7/go_final_project/database"
	"github.com/sirin7/go_final_project/server"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("File download error .env: %v", err)
	}

	// Инициализация базы данных
	database.InitDB()

	// Запуск сервера
	server.StartServer()

}
