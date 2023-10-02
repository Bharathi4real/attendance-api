# attendance-api
attendance-api/
|-- api/
|   |-- handlers.go           # Contains API request handlers
|-- config/
|   |-- config.json           # Stores admin login credentials
|-- db/
|   |-- postgres.go           # Handles PostgreSQL database interactions
|   |-- sql.go                # Contains functions for SQL queries and table creation
|-- models/
|   |-- user.go               # Defines the user model
|-- src/
|   |-- main.go               # Entry point of the application






# jwt generator

package main

import (
"crypto/rand"
"encoding/base64"
"fmt"
)

func generateRandomString(length int) (string, error) {
bytes := make([]byte, length)
if _, err := rand.Read(bytes); err != nil {
return "", err
}
return base64.URLEncoding.EncodeToString(bytes), nil
}

func main() {
secret, err := generateRandomString(32) // Generate a 32-byte (256-bit) secret
if err != nil {
fmt.Println("Error generating secret:", err)
return
}
fmt.Println("Generated JWT secret:", secret)
}




package api

import (
"attendance-api/db"
"attendance-api/models"
"encoding/json"
"fmt"
"github.com/dgrijalva/jwt-go"
"github.com/gin-gonic/gin"
"net/http"
"os"
"path/filepath"
"time"
)

func loadAdminCredentials() (string, string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current working directory: %v", err)
	}

	configFilePath := filepath.Join(cwd, "..", "config", "config.json")
	fmt.Println("Attempting to load admin credentials from:", configFilePath)

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {

		}
	}(configFile)

	var config map[string]string
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		fmt.Println("Error decoding config file:", err)
		return "", "", fmt.Errorf("failed to decode config file: %v", err)
	}

	adminUsername, ok := config["admin_username"]
	if !ok {
		fmt.Println("Admin username not found in config")
		return "", "", fmt.Errorf("admin_username not found in config")
	}

	adminPassword, ok := config["admin_password"]
	if !ok {
		fmt.Println("Admin password not found in config")
		return "", "", fmt.Errorf("admin_password not found in config")
	}

	return adminUsername, adminPassword, nil
}

var jwtSecret = []byte("UOEapGWYMB9wa5rtNUUfFl9EBS_38JUCFl_MTb1DPSM=")

func generateToken(username string) (string, error) {
claims := jwt.MapClaims{
"username": username,
"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiry time (e.g., 24 hours)
}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Login(c *gin.Context) {
adminUsername, adminPassword, err := loadAdminCredentials()
if err != nil {
fmt.Println("Error loading admin credentials:", err)
c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load admin credentials: %v", err)})
return
}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if user.Username == adminUsername && user.Password == adminPassword {
		// Admin login successful, generate a token (e.g., JWT) here
		token, err := generateToken(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Admin logged in successfully",
			"token":   token, // Include the token in the response
		})
		return
	}

	authenticated := db.FacultyLogin(user.Username, user.Password)

	if authenticated {
		c.JSON(http.StatusOK, gin.H{"message": "Faculty logged in successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func AuthenticateUser(c *gin.Context) {
var user models.User
if err := c.BindJSON(&user); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
return
}

	// Load admin credentials from config
	adminUsername, adminPassword, err := loadAdminCredentials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load admin credentials"})
		return
	}

	// Check if the user is an admin based on the loaded credentials
	isAdmin := user.Username == adminUsername && user.Password == adminPassword

	if isAdmin {
		c.Set("role", "admin")
	} else {
		c.Set("role", "user")
	}

	c.Next()
}

func CreateFacultyCredentials(c *gin.Context) {
AuthenticateUser(c)
role, _ := c.Get("role")
if role != "admin" {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can perform this action"})
return
}

	// Parse faculty information from the request
	var facultyInfo models.User
	if err := c.ShouldBindJSON(&facultyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Save faculty credentials to the database
	facultyID, err := db.AddFaculty("faculty_username", "faculty_password")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create faculty credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty credentials created successfully", "faculty_id": facultyID})
}
