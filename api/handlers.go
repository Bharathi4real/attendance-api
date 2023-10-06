package api

import (
	"attendance-api/db"
	"attendance-api/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
	"time"
)

func loadAdminCredentials() (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", fmt.Errorf("failed to load .env file: %v", err)
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" {
		return "", "", fmt.Errorf("admin_username not found in .env")
	}

	if adminPassword == "" {
		return "", "", fmt.Errorf("admin_password not found in .env")
	}

	return adminUsername, adminPassword, nil
}

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_TOKEN"))

}

func generateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
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

func verifyToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

	// Extract the token from the Authorization header
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.Abort()
		return
	}

	tokenString = parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
		c.Abort()
		return
	}

	c.Set("role", role)
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

// make a func which enables access to all faculty function to admin but restricts faculty to access admin func.
