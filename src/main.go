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

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
