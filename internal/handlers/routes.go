package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhnthnljyng/rrf-be/internal/middleware"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB, jwtSecret string) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", Register(db))
		auth.POST("/login", Login(db, jwtSecret))
	}

	r.GET("/api/listings", GetListings(db))
	r.GET("/api/listings/:id", GetListingByID(db))

	api := r.Group("/api")
	api.Use(middleware.RequireAuth(jwtSecret))
	{
		api.GET("/user/profile", GetProfile(db))
		api.PUT("/user/profile", UpdateProfile(db))
		api.POST("/user/profile/picture", UploadProfilePicture(db))
		api.POST("/listings", CreateProperty(db))
		api.POST("/images/upload", UploadImages(db))
	}
}
