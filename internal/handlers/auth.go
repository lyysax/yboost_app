package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"yboost_app/internal/db"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GET /login — page de connexion
func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

// GET /register — page d'inscription
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/register.html")
}

// POST /register — créer un compte
func Register(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	if u.Email == "" || u.Password == "" {
		http.Error(w, "Email et mot de passe requis", http.StatusBadRequest)
		return
	}

	// Hash du mot de passe
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Insertion en BDD
	err = db.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		u.Email, string(hash),
	).Scan(&u.ID)

	if err != nil {
		http.Error(w, "Email déjà utilisé", http.StatusConflict)
		return
	}

	// Crée la session
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(u.ID),
		Path:  "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Compte créé !"})
}

// POST /login — se connecter
func Login(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	if u.Email == "" || u.Password == "" {
		http.Error(w, "Email et mot de passe requis", http.StatusBadRequest)
		return
	}

	// Cherche l'utilisateur en BDD
	var stored User
	err := db.DB.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		u.Email,
	).Scan(&stored.ID, &stored.Email, &stored.Password)

	if err != nil {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Vérifie le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(u.Password))
	if err != nil {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Crée la session
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(stored.ID),
		Path:  "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Connecté !"})
}

// POST /logout — se déconnecter
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Déconnecté !"})
}

// GetUserID — récupère l'ID utilisateur depuis le cookie
func GetUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// test
