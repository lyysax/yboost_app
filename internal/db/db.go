package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Erreur ouverture BDD :", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erreur connexion BDD :", err)
	}

	fmt.Println("✅ Connecté à Supabase !")
}

// test
