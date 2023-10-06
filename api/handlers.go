package api

import (
	"attendance-api/db"
	"attendance-api/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	var token string
	if user.Username == adminUsername && user.Password == adminPassword {
		token, err = generateToken(user.Username, "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Admin logged in successfully",
			"token":   token,
		})
		return
	}

	facultyAuthenticated := db.FacultyLogin(user.Username, user.Password)
	if facultyAuthenticated {
		token, err = generateToken(user.Username, "faculty")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Faculty logged in successfully",
			"token":   token,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password for the provided username"})
}

func CreateFacultyCredentials(c *gin.Context) {

	verifyToken(c)

	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can perform this action"})
		return
	}

	var facultyInfo models.User
	if err := c.ShouldBindJSON(&facultyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	exists, err := db.UsernameExists(facultyInfo.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking username existence"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username exists already, try another username"})
		return
	}

	facultyID, err := db.AddFaculty(facultyInfo.Username, facultyInfo.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create faculty credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty credentials created successfully", "faculty_id": facultyID, "faculty_username": facultyInfo.Username})
}

func DeleteFaculty(c *gin.Context) {

	verifyToken(c)

	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can perform this action"})
		return
	}

	var facultyInfo models.User
	if err := c.ShouldBindJSON(&facultyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	exists, err := db.UsernameExists(facultyInfo.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking username existence"})
		return
	}

	if !exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username doesn't exist"})
		return
	}
	err = db.RemoveFaculty(facultyInfo.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete faculty credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty credentials deleted successfully", "faculty_id": facultyInfo})
}
