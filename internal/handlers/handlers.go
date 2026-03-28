package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"yboost_app/internal/db"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

// GET /health — état de la BDD
func Health(w http.ResponseWriter, r *http.Request) {
	err := db.DB.Ping()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "BDD non joignable"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "BDD connectée ✅"})
}

// GET /todos — liste les todos de l'utilisateur connecté
func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query(
		"SELECT id, title, done, created_at FROM todos WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
		todos = append(todos, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// POST /todos — créer une todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	var t Todo
	json.NewDecoder(r.Body).Decode(&t)

	if t.Title == "" {
		http.Error(w, "Le titre est requis", http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow(
		"INSERT INTO todos (title, user_id) VALUES ($1, $2) RETURNING id, title, done, created_at",
		t.Title, userID,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)

	if err != nil {
		http.Error(w, "Erreur création", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// PUT /todos/{id} — marquer une todo comme faite/pas faite
func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	var t Todo
	err = db.DB.QueryRow(
		"UPDATE todos SET done = NOT done WHERE id = $1 AND user_id = $2 RETURNING id, title, done, created_at",
		id, userID,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Todo introuvable", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// DELETE /todos/{id} — supprimer une todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	result, err := db.DB.Exec(
		"DELETE FROM todos WHERE id = $1 AND user_id = $2",
		id, userID,
	)
	if err != nil {
		http.Error(w, "Erreur suppression", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Todo introuvable", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /meteo
func Meteo(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://wttr.in/Bordeaux?format=j1")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "API météo indisponible",
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	current := result["current_condition"].([]interface{})[0].(map[string]interface{})
	temp := current["temp_C"].(string)
	desc := current["weatherDesc"].([]interface{})[0].(map[string]interface{})["value"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"temperature": temp + "°C",
		"description": desc,
		"ville":       "Bordeaux",
	})
}
