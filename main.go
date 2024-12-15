package main

import (
	"log"
	"net/http"

	"example.com/m/v2/db"
	"example.com/m/v2/env"
	"example.com/m/v2/handlers"
	"example.com/m/v2/sessions"
)

func main() {
	// Загружаем переменные окружения
	env.LoadEnv()

	// Подключаемся к базе данных
	dbConn := db.Connect()
	defer dbConn.Close()

	// Запускаем миграцию
	db.MigrateToDB(dbConn)

	// Инициализация очистки сессий
	sessions.InitSessionCleanup()

	// Регистрация обработчиков
	http.HandleFunc("/login", handlers.LoginHandler(dbConn))
	http.HandleFunc("/register", handlers.RegisterHandler(dbConn))
	http.HandleFunc("/logout", sessions.SessionMiddleware(handlers.LogoutHandler()))
	http.HandleFunc("/protected", sessions.SessionMiddleware(handlers.ProtectedHandler()))

	// Запуск сервера
	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
