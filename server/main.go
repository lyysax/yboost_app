package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"yboost_app/internal/db"
	"yboost_app/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	root, _ := os.Getwd()

	godotenv.Load(filepath.Join(root, ".env"))

	db.Connect()

	mux := http.NewServeMux()

	// Pages
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, _ := os.ReadFile(filepath.Join(root, "templates", "index.html"))
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		content, _ := os.ReadFile(filepath.Join(root, "templates", "login.html"))
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		content, _ := os.ReadFile(filepath.Join(root, "templates", "register.html"))
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	// Auth
	mux.HandleFunc("POST /login", handlers.Login)
	mux.HandleFunc("POST /register", handlers.Register)
	mux.HandleFunc("POST /logout", handlers.Logout)

	// Santé de la BDD
	mux.HandleFunc("GET /health", handlers.Health)

	// CRUD Todos
	mux.HandleFunc("GET /todos", handlers.GetTodos)
	mux.HandleFunc("POST /todos", handlers.CreateTodo)
	mux.HandleFunc("PUT /todos/{id}", handlers.ToggleTodo)
	mux.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodo)

	// Météo
	mux.HandleFunc("GET /meteo", handlers.Meteo)

	// Fichiers statiques
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(root, "static")))))

	port := os.Getenv("PORT")
	fmt.Println("🚀 Serveur lancé sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
