package main

import (
	"attendance-api/api"
	"attendance-api/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	err := db.InitDB()
	if err != nil {
		return
	}

	r := gin.Default()

	r.POST("/login", api.Login)
	r.OPTIONS("/login", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/create-faculty", api.CreateFacultyCredentials)
	r.POST("/delete-faculty", api.DeleteFaculty)

	err = r.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
