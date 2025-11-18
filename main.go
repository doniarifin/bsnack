package main

import (
	db "bsnack/database"
	"bsnack/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	DB, err := db.ConnectDB()
	if err != nil {
		log.Fatal("error", err)
	}

	app := gin.Default()
	trusted := []string{
		"127.0.0.1",
	}
	err = app.SetTrustedProxies(trusted)
	if err != nil {
		log.Fatal("Proxy config error:", err)
	}

	routes.SetupRoutes(app, DB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port " + port)
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
