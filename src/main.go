package main

import (
	"attendance-api/api"
	"attendance-api/db"
	"github.com/gin-gonic/gin"
)

func main() {
	err := db.InitDB()
	if err != nil {
		return
	}

	r := gin.Default()

	r.POST("/login", api.Login)
	r.POST("/create-faculty", api.CreateFacultyCredentials) // New route for creating faculty credentials

	err = r.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
