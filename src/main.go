package main

import (
	"attendance-api/api"
	"attendance-api/db"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	r := gin.Default()

	api.SetupRoutes(r)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
