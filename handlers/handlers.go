package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"example.com/m/v2/repository"
	"example.com/m/v2/sessions"
)

// LoginHandler обрабатывает авторизацию и создаёт сессию
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		if login == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Поля не могут быть пустыми"))
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT Password FROM Users WHERE Login = $1", login).Scan(&storedPassword)
		if err == sql.ErrNoRows || storedPassword != password {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Неправильный логин или пароль"))
			return
		} else if err != nil {
			log.Println("Ошибка при запросе в базу данных:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sessionID := sessions.CreateSession(login)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
		})
		w.Write([]byte("Авторизация успешна\n"))
		w.Write([]byte("Ваш токен активен\n"))
	}
}

// RegisterHandler обрабатывает регистрацию
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		if login == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Поля не могут быть пустыми"))
			return
		}

		// Вставляем пользователя в базу данных
		_, err := db.Exec("INSERT INTO Users (Login, Password) VALUES ($1, $2)", login, password)
		if err != nil {
			log.Println("Ошибка при создании пользователя:", err)
			w.Write([]byte("Ошибка создания пользователя"))
			return
		}

		w.Write([]byte("Пользователь успешно создан"))
	}
}

// ProtectedHandler - защищённый маршрут
func ProtectedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("X-Username")
		w.Write([]byte("Привет, " + username + "! Это защищённый маршрут."))
	}
}

// LogoutHandler удаляет сессию
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			sessions.DeleteSession(cookie.Value)
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
		w.Write([]byte("Вы вышли из системы"))
	}
}

// ProfileHandler выводит профиль пользователя, включая дату регистрации
func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем имя пользователя из заголовка, установленного middleware
		username := r.Header.Get("X-Username")

		// Получаем дату регистрации пользователя из базы данных
		createdAt, err := repository.GetUserCreatedAt(db, username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка при получении данных пользователя"))
			return
		}

		// Формируем и отправляем ответ
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Добро пожаловать в ваш профиль, %s!\n", username)))
		w.Write([]byte(fmt.Sprintf("Дата регистрации: %s\n", createdAt)))
	}
}
