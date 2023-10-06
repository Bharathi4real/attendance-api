package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
	"time"
)

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

func verifyToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

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

func admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if exists && role == "admin" {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin access only"})
			c.Abort()
		}
	}
}
