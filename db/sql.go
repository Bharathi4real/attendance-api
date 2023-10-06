package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

var conn *pgx.Conn

//func InitDBS() {
//	databaseURL := os.Getenv("DB_URL")
//
//	if !strings.Contains(databaseURL, "sslmode") {
//		databaseURL += "?sslmode=disable"
//	}
//
//	var err error
//
//	conn, err = pgx.Connect(context.Background(), databaseURL)
//	if err != nil {
//		log.Fatalf("Error connecting to the database: %v", err)
//	}
//
//	createFacultyTable()
//}

func createFacultyTable() {
	_, err := conn.Exec(context.Background(), `
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
	err := conn.QueryRow(context.Background(), query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func AddFaculty(username, password string) (int, error) {
	const query = "INSERT INTO faculty (username, password) VALUES ($1, $2) RETURNING id"

	var lastInsertID int
	err := conn.QueryRow(context.Background(), query, username, password).Scan(&lastInsertID)
	if err != nil {
		log.Printf("Error inserting faculty: %v", err)
		return 0, err
	}

	return lastInsertID, nil
}

func RemoveFaculty(username string) error {
	const query = "DELETE FROM faculty WHERE username = $1"

	_, err := conn.Exec(context.Background(), query, username)
	if err != nil {
		log.Printf("Error deleting faculty: %v", err)
		return err
	}

	return nil
}
