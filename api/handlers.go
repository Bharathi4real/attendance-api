package api

import (
	"attendance-api/db"
	"attendance-api/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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
		c.JSON(http.StatusOK, gin.H{"message": "Admin logged in successfully"})
		return
	}

	authenticated := db.FacultyLogin(user.Username, user.Password)

	if authenticated {
		c.JSON(http.StatusOK, gin.H{"message": "Faculty logged in successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
