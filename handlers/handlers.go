package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"example.com/m/v2/repository"
	"example.com/m/v2/sessions"
)

// LoginHandler обрабатывает авторизацию и создаёт сессию
// LoginHandler обрабатывает авторизацию и создаёт сессию
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		w.Header().Set("Content-Type", "application/json") // Устанавливаем тип ответа

		if login == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "message": "Поля не могут быть пустыми"}`))
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT Password FROM Users WHERE Login = $1", login).Scan(&storedPassword)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success": false, "message": "Пользователь не найден"}`))
			return
		} else if err != nil {
			log.Println("Ошибка при запросе в базу данных:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "message": "Внутренняя ошибка сервера"}`))
			return
		}

		if storedPassword != password {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success": false, "message": "Неправильный пароль"}`))
			return
		}

		sessionID := sessions.CreateSession(login)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
		})

		w.Write([]byte(`{"success": true, "message": "Успешный вход в систему"}`))
	}
}

// RegisterHandler обрабатывает регистрацию
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if login == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Поля не могут быть пустыми"))
			return
		}

		// Вставляем пользователя в базу данных
		_, err := db.Exec("INSERT INTO Users (Login, Password) VALUES ($1, $2)", login, password)
		if err != nil {
			log.Println("Ошибка при создании пользователя:", err)
			w.Write([]byte("<p style='color: red;'>Ошибка создания пользователя. Возможно, логин уже существует.</p>"))
			return
		}

		// Сообщение об успешной регистрации
		w.Write([]byte("<p style='color: green;'>Пользователь успешно зарегистрирован! Теперь вы можете войти в систему.</p>"))
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

// ProfileData структура для хранения данных профиля
type ProfileData struct {
	Username  string
	CreatedAt string
}

// ProfileHandler выводит профиль пользователя
func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем имя пользователя из заголовка
		username := r.Header.Get("X-Username")
		if username == "" {
			http.Error(w, "Unauthorized: отсутствует имя пользователя", http.StatusUnauthorized)
			return
		}

		// Получаем дату регистрации пользователя из базы данных
		createdAt, err := repository.GetUserCreatedAt(db, username)
		if err != nil {
			log.Println("Ошибка при получении даты регистрации:", err)
			http.Error(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError)
			return
		}

		// Загружаем шаблон profile.html
		tmpl, err := template.ParseFiles("./static/profile.html")
		if err != nil {
			log.Println("Ошибка при загрузке шаблона:", err)
			http.Error(w, "Ошибка при загрузке страницы профиля", http.StatusInternalServerError)
			return
		}

		// Заполняем шаблон данными
		data := ProfileData{
			Username:  username,
			CreatedAt: createdAt,
		}

		// Отправляем HTML с данными
		if err := tmpl.Execute(w, data); err != nil {
			log.Println("Ошибка при выполнении шаблона:", err)
			http.Error(w, "Ошибка при отображении страницы профиля", http.StatusInternalServerError)
		}
	}
}
