package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/login", Login)
	r.OPTIONS("/login", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/create-faculty", CreateFacultyCredentials)
	r.OPTIONS("/create-faculty", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/delete-faculty", DeleteFaculty)
	r.OPTIONS("/delete-faculty", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
