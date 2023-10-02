package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", "postgres://nvpqlwmr:hvpv3FW97Qd4TFIN2ECVlDFynCjhtnKj@bubble.db.elephantsql.com/nvpqlwmr?sslmode=disable")
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
