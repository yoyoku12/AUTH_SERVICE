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
	// env load
	env.LoadEnv()

	// db connection
	dbConn := db.Connect()
	defer dbConn.Close()

	// migreate db
	db.MigrateToDB(dbConn)

	// routes
	http.HandleFunc("/login", handlers.LoginHandler(dbConn))
	http.HandleFunc("/register", handlers.RegisterHandler(dbConn))
	// middleware + routes
	http.HandleFunc("/logout", sessions.SessionMiddleware(handlers.LogoutHandler()))
	http.HandleFunc("/protected", sessions.SessionMiddleware(handlers.ProtectedHandler()))
	http.HandleFunc("/profile", sessions.SessionMiddleware(handlers.ProfileHandler(dbConn)))

	// start server
	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
