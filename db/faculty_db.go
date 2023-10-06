package db

import (
	"context"
	"log"
)

func FacultyLogin(username, password string) bool {
	const query = "SELECT COUNT(*) FROM faculty WHERE username = $1 AND password = $2"

	var count int
	err := conn.QueryRow(context.Background(), query, username, password).Scan(&count)
	if err != nil {
		log.Printf("Error during login query: %v", err)
		return false
	}

	return count == 1
}
