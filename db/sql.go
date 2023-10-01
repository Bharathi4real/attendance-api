package db

import (
	"fmt"
	"log"
)

func CreateTable() error {

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS faculty (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	log.Println("Faculty table created successfully")
	return nil
}
