package main

import (
	"attendance-api/api"
	"attendance-api/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	db.InitDBS()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	r := gin.Default()

	api.SetupRoutes(r)
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if err := r.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
