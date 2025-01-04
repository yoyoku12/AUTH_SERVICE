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

	env.LoadEnv()

	dbConn := db.Connect()
	defer dbConn.Close()

	db.MigrateToDB(dbConn)

	sessions.InitSessionCleanup()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	})

	http.HandleFunc("/profile", sessions.SessionMiddleware(handlers.ProfileHandler(dbConn)))

	http.HandleFunc("/logout", sessions.SessionMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/logout.html")
	}))

	http.HandleFunc("/login_action", handlers.LoginHandler(dbConn))
	http.HandleFunc("/register_action", handlers.RegisterHandler(dbConn))
	http.HandleFunc("/profile_action", sessions.SessionMiddleware(handlers.ProfileHandler(dbConn)))
	http.HandleFunc("/logout_action", sessions.SessionMiddleware(handlers.LogoutHandler()))

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
