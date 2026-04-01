package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"yboost_app/internal/db"
	"yboost_app/internal/handlers"

	"github.com/joho/godotenv"
)

// fichiers html et css intégrés dans le binaire

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

func main() {
	godotenv.Load() // variables d'environnement depuis .env
	db.Connect()

	mux := http.NewServeMux() // routeur

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, _ := templatesFS.ReadFile("templates/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		content, _ := templatesFS.ReadFile("templates/login.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		content, _ := templatesFS.ReadFile("templates/register.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	// authentificaion
	mux.HandleFunc("POST /login", handlers.Login)
	mux.HandleFunc("POST /register", handlers.Register)
	mux.HandleFunc("POST /logout", handlers.Logout)
	mux.HandleFunc("GET /health", handlers.Health)

	// les todos
	mux.HandleFunc("GET /todos", handlers.GetTodos)
	mux.HandleFunc("POST /todos", handlers.CreateTodo)
	mux.HandleFunc("PUT /todos/{id}", handlers.ToggleTodo)
	mux.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodo)

	// météo
	mux.HandleFunc("GET /meteo", handlers.Meteo)

	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Serveur lancé sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
