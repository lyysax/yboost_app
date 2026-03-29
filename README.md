# Todo List App

Application de gestion de tâches en ligne, développée en Go et déployée sur Render.

## 🌐 Lien de l'application

https://yboost-app.onrender.com

## 📖 Présentation

Yboost est une todo list multi-utilisateurs. Chaque utilisateur peut créer un compte,
se connecter, et gérer ses propres tâches. L'application affiche aussi la météo 
en temps réel via l'API wttr.in.

## Stack technique

- **Backend** : Go (net/http)
- **Base de données** : PostgreSQL via Supabase
- **Driver BDD** : database/sql + lib/pq
- **Déploiement** : Render
- **API externe** : wttr.in (météo)

## Lancer en local

### Prérequis
- Go 1.22+
- Un projet Supabase avec les tables créées

### Installation

1. Clone le repo
\```bash
git clone https://github.com/TON_USERNAME/yboost_app.git
cd yboost_app
\```

2. Crée un fichier `.env` à la racine
\```
DATABASE_URL=postgresql://...@....supabase.com:5432/postgres?sslmode=require
PORT=8080
\```

3. Lance l'application
\```bash
go run main.go
\```

4. Ouvre http://localhost:8080
