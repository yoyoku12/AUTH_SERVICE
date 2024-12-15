package sessions

import (
	"math/rand"
	"net/http"
	"time"
)

// SessionStore хранит активные сессии
var SessionStore = make(map[string]string)

// GenerateSessionID генерирует случайный Session ID
func GenerateSessionID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// CreateSession создаёт сессию для пользователя
func CreateSession(username string) string {
	sessionID := GenerateSessionID()
	SessionStore[sessionID] = username
	return sessionID
}

// GetUsername возвращает имя пользователя по sessionID
func GetUsername(sessionID string) (string, bool) {
	username, exists := SessionStore[sessionID]
	return username, exists
}

// DeleteSession удаляет сессию
func DeleteSession(sessionID string) {
	delete(SessionStore, sessionID)
}

// SessionMiddleware проверяет наличие и валидность сессии
func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем cookie с session_id
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: сессия отсутствует"))
			return
		}

		// Проверяем, существует ли сессия
		username, exists := GetUsername(cookie.Value)
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: невалидная сессия"))
			return
		}

		// Добавляем имя пользователя в заголовок для хендлера
		r.Header.Set("X-Username", username)

		// Передаём управление следующему обработчику
		next(w, r)
	}
}
