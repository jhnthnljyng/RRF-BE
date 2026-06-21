package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jhnthnljyng/rrf-be/internal/config"
	"github.com/jhnthnljyng/rrf-be/internal/handlers"
)

func main() {
	cfg := config.Load()

	if err := config.RunMigrations(cfg); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	handlers.RegisterRoutes(r, db, cfg.JWTSecret)

	log.Println("Registered routes:")
	for _, route := range r.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	r.NoRoute(func(c *gin.Context) {
		log.Printf("404 - no route matched: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{"error": "route not found"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
