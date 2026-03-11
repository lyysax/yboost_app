package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		log.Fatal("La variable d'environnement DATABASE_URL est requise")
	}

	if strings.Contains(dsn, "[YOUR-PASSWORD]") {
		log.Fatal("DATABASE_URL contient [YOUR-PASSWORD]. Remplace ce placeholder par ton vrai mot de passe")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("DATABASE_URL invalide: %v", err)
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		log.Fatalf("DATABASE_URL invalide: scheme attendu postgres/postgresql, recu %q", u.Scheme)
	}
	if u.Host == "" {
		log.Fatal("DATABASE_URL invalide: host manquant")
	}
	if u.User == nil || u.User.Username() == "" {
		log.Fatal("DATABASE_URL invalide: user manquant")
	}

	log.Printf("Tentative de connexion DB: host=%s user=%s", u.Host, u.User.Username())

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Erreur lors de la récupération de la connexion SQL : %v", err)
	} else {
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Impossible de ping la base : %v", err)
		}

		log.Printf("✓ Connexion à la base établie")
		log.Printf("✓ Max open connections: %d", sqlDB.Stats().MaxOpenConnections)
		log.Printf("✓ Connexions actives: %d", sqlDB.Stats().OpenConnections)
	}
}
