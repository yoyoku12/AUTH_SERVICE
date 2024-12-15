package repository

import (
	"database/sql"
	"fmt"
)

func GetUserCreatedAt(db *sql.DB, username string) (string, error) {
	var createdAt string

	// Отладочный лог для проверки параметра
	fmt.Println("Запрос пользователя с логином:", username)

	err := db.QueryRow("SELECT CreatedAt FROM Users WHERE Login = $1", username).Scan(&createdAt)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("пользователь %s не найден", username)
	} else if err != nil {
		return "", err
	}

	return createdAt, nil
}
