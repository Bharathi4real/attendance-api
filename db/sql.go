package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func init() {
	database := os.Getenv("DB_URL")

	var err error

	db, err = sql.Open("postgres", database)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	createFacultyTable()
}

func createFacultyTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS faculty (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Error creating faculty table: %v", err)
	}
}

func UsernameExists(username string) (bool, error) {

	const query = "SELECT COUNT(*) FROM faculty WHERE username = $1"

	var count int
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func AddFaculty(username, password string) (int, error) {
	const query = "INSERT INTO faculty (username, password) VALUES ($1, $2) RETURNING id"

	var lastInsertID int
	err := db.QueryRow(query, username, password).Scan(&lastInsertID)
	if err != nil {
		log.Printf("Error inserting faculty: %v", err)
		return 0, err
	}

	return lastInsertID, nil
}

func RemoveFaculty(username string) error {
	const query = "DELETE FROM faculty WHERE username = $1"

	_, err := db.Exec(query, username)
	if err != nil {
		log.Printf("Error deleting faculty: %v", err)
		return err
	}

	return nil
}
