package main

import (
	"log"
	"net/http"

	"example.com/m/v2/db"
	"example.com/m/v2/env"
	"example.com/m/v2/handlers"
)

func main() {
	// Загружаем переменные окружения
	env.LoadEnv()

	// Подключаемся к базе данных
	dbConn := db.Connect()
	defer dbConn.Close()

	// Миграция базы данных
	db.MigrateToDB(dbConn)

	// Регистрация обработчиков
	http.HandleFunc("/login", handlers.LoginHandler(dbConn))
	http.HandleFunc("/register", handlers.RegisterHandler(dbConn))

	// Запуск HTTP-сервера
	log.Println("Сервер запущен на порту 8080 и ожидает запросов")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
