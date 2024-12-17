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

	// Маршрут для страницы входа
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	// Маршрут для страницы регистрации
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	})

	// Маршрут для динамического отображения профиля
	http.HandleFunc("/profile", sessions.SessionMiddleware(handlers.ProfileHandler(dbConn)))

	// Маршрут для страницы выхода
	http.HandleFunc("/logout", sessions.SessionMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/logout.html")
	}))

	// Регистрация обработчиков для запросов
	http.HandleFunc("/login_action", handlers.LoginHandler(dbConn))       // Обработчик логина
	http.HandleFunc("/register_action", handlers.RegisterHandler(dbConn)) // Обработчик регистрации
	http.HandleFunc("/profile_action", sessions.SessionMiddleware(handlers.ProfileHandler(dbConn)))
	http.HandleFunc("/logout_action", sessions.SessionMiddleware(handlers.LogoutHandler()))

	// Запуск сервера
	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
