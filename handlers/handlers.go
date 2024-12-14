package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// LoginHandler обрабатывает запросы на авторизацию
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметры из строки запроса
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		// Проверяем, что поля не пустые
		if login == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Поля не могут быть пустыми"))
			return
		}

		// SQL-запрос для получения пароля пользователя по логину
		var storedPassword string
		err := db.QueryRow("SELECT Password FROM Users WHERE Login = $1", login).Scan(&storedPassword)
		if err == sql.ErrNoRows {
			w.Write([]byte("Пользователь не найден"))
			log.Println("Пользователь не найден:", login)
			return
		} else if err != nil {
			log.Println("Ошибка при запросе в базу данных:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Сравниваем введённый пароль с сохранённым
		if storedPassword == password {
			w.Write([]byte("Авторизация успешна"))
		} else {
			w.Write([]byte("Неправильный пароль"))
		}
	}
}

// RegisterHandler обрабатывает запросы на регистрацию
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметры из строки запроса
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		// Проверяем, что поля не пустые
		if login == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Поля не могут быть пустыми"))
			return
		}

		// SQL-запрос для добавления нового пользователя
		_, err := db.Exec("INSERT INTO Users (Login, Password) VALUES ($1, $2)", login, password)
		if err != nil {
			log.Println("Ошибка при создании пользователя:", err)
			w.Write([]byte("Ошибка создания пользователя"))
			return
		}

		w.Write([]byte("Пользователь успешно создан"))
	}
}
