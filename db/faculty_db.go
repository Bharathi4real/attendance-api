package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func InitDB() error {
	database := os.Getenv("DB_URL")
	var err error
	db, err = sql.Open("postgres", database)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func FacultyLogin(username, password string) bool {
	query := "SELECT COUNT(*) FROM faculty WHERE username=$1 AND password=$2"
	var count int
	err := db.QueryRow(query, username, password).Scan(&count)
	if err != nil {
		log.Println("Error during login query:", err)
		return false
	}

	return count == 1
}
