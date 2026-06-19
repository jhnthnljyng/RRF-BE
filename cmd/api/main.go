package main

import (
	"log"

	"github.com/jhnthnljyng/rrf-be/internal/config"
	"github.com/jhnthnljyng/rrf-be/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	handlers.RegisterRoutes(r, db)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
