package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func migrateToDB(db *sql.DB) {
	log.Println("Starting migrate to database...")
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS User (
				Login TEXT NOT NULL,
            	Password TEXT NOT NULL
        )
    `)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		panic(err)
	}
	log.Println("Database migrated successfully.")
}

func main() {

	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}

	migrateToDB(db)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		getLogin := r.URL.Query().Get("login")
		getPassword := r.URL.Query().Get("password")

		if getLogin == "" || getPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("Поля не могут быть пустыми"))
			if err != nil {
				log.Println("w.write error", err)
			}
			return
		}
		query := "SELECT * FROM User WHERE Login = ?"
		rows, err := db.Query(query, getLogin)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		if rows.Next() == true {
			var storedPassword string
			err := rows.Scan(&getLogin, &storedPassword)
			if err != nil {
				log.Fatal(err)
			}
			if storedPassword == getPassword {
				w.Write([]byte("Авторизация успешна"))
				log.Println("Авторизация успешна")
			} else {
				w.Write([]byte("Неправильный пароль"))
				log.Println("Неправильный пароль")
			}

		} else {
			w.Write([]byte("Пользователь не найден"))
			log.Println("Пользователь не найден")
		}

	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		getLogin := r.URL.Query().Get("login")
		getPassword := r.URL.Query().Get("password")

		if getLogin == "" || getPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("Поля не могут быть пустыми"))
			if err != nil {
				log.Println("w.write error", err)
			}
			return
		}

		_, err = db.Exec("INSERT INTO User (Login, Password) VALUES (?, ?)", getLogin, getPassword)
		if err != nil {
			log.Println("Ошибка создания пользователя")
			w.Write([]byte("Ошибка создания пользователя"))
		}
		log.Println("Пользователь создан")
		w.Write([]byte("Пользователь создан"))

	})

	serverStart()
}

func serverStart() {
	log.Println("Server starting...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil { // Проверяем на ошибку.
		log.Println("Server error", err)
	}
}
