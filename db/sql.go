package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

var conn *pgx.Conn
var pool *pgxpool.Pool

func InitDBS() {
	databaseURL := os.Getenv("DB_URL")

	if !strings.Contains(databaseURL, "sslmode") {
		databaseURL += "?sslmode=disable"
	}

	var err error

	conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	createFacultyTable()
}

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

// AssignFacultyToClass assigns a faculty to a class
func AssignFacultyToClass(facultyID, classID string) error {
	const query = "UPDATE classes SET faculty_id = $1 WHERE id = $2"
	_, err := pool.Exec(context.Background(), query, facultyID, classID)
	if err != nil {
		log.Printf("Error assigning faculty to class: %v", err)
		return fmt.Errorf("failed to assign faculty to class")
	}
	return nil
}

// AddStudentToClass adds a student to a class
func AddStudentToClass(studentID, classID string) error {
	const query = "INSERT INTO class_students (student_id, class_id) VALUES ($1, $2)"
	_, err := pool.Exec(context.Background(), query, studentID, classID)
	if err != nil {
		log.Printf("Error adding student to class: %v", err)
		return fmt.Errorf("failed to add student to class")
	}
	return nil
}

// RemoveStudentFromClass removes a student from a class
func RemoveStudentFromClass(studentID, classID string) error {
	const query = "DELETE FROM class_students WHERE student_id = $1 AND class_id = $2"
	_, err := pool.Exec(context.Background(), query, studentID, classID)
	if err != nil {
		log.Printf("Error removing student from class: %v", err)
		return fmt.Errorf("failed to remove student from class")
	}
	return nil
}

// GetAttendanceInfo gets attendance information for a class
func GetAttendanceInfo(classID string) (int, int, int, error) {
	const query = "SELECT COUNT(*) FROM class_students WHERE class_id = $1"
	var totalStudents, presentStudents, absentStudents int
	err := pool.QueryRow(context.Background(), query, classID).Scan(&totalStudents)
	if err != nil {
		log.Printf("Error getting attendance info: %v", err)
		return 0, 0, 0, fmt.Errorf("failed to get attendance info")
	}

	const presentQuery = "SELECT COUNT(*) FROM attendance WHERE class_id = $1 AND status = 'present'"
	err = pool.QueryRow(context.Background(), presentQuery, classID).Scan(&presentStudents)
	if err != nil {
		log.Printf("Error getting present students: %v", err)
		return 0, 0, 0, fmt.Errorf("failed to get attendance info")
	}

	absentStudents = totalStudents - presentStudents

	return totalStudents, presentStudents, absentStudents, nil
}
