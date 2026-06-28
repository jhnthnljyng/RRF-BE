package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhnthnljyng/rrf-be/internal/models"
	"github.com/jhnthnljyng/rrf-be/internal/repository"
)

func GetProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		user, err := repository.GetUserByID(db, userID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Printf("get profile error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func UpdateProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		var req models.UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.UpdateProfile(db, userID, &req)
		if err != nil {
			log.Printf("update profile error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "profile updated successfully",
			"user":    user,
		})
	}
}

func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		user, err := repository.GetUserByID(db, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Printf("get user by id error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func SearchUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter q is required"})
			return
		}

		users, err := repository.SearchUsers(db, q)
		if err != nil {
			log.Printf("search users error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not search users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users, "total": len(users)})
	}
}

func UploadProfilePicture(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		file, err := c.FormFile("picture")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no picture provided"})
			return
		}

		uploadDir := filepath.Join("uploads", "avatars")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Printf("failed to create avatar dir: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not prepare upload directory"})
			return
		}

		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("failed to save avatar: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save picture"})
			return
		}

		avatarURL := "/" + filepath.ToSlash(savePath)
		if err := repository.UpdateAvatarURL(db, userID, avatarURL); err != nil {
			log.Printf("failed to update avatar_url: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update profile picture"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "profile picture updated",
			"avatar_url": avatarURL,
		})
	}
}
