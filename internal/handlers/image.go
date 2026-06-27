package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhnthnljyng/rrf-be/internal/repository"
)

func UploadImages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uploadType := c.Query("type")
		if uploadType == "" {
			uploadType = c.PostForm("type")
		}
		log.Printf("image upload request: type=%q content-type=%q", uploadType, c.ContentType())

		if uploadType == "avatar" {
			file, err := c.FormFile("image")
			if err != nil {
				file2, err2 := c.FormFile("images[]")
				if err2 != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "no image provided"})
					return
				}
				file = file2
			}

			userID := c.GetInt("user_id")
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
				return
			}

			avatarURL := "/" + filepath.ToSlash(savePath)
			if err := repository.UpdateAvatarURL(db, userID, avatarURL); err != nil {
				log.Printf("failed to update avatar_url: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update profile picture"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"paths":      []string{avatarURL},
				"avatar_url": avatarURL,
			})
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse form"})
			return
		}

		files := form.File["images[]"]
		if len(files) == 0 {
			files = form.File["images"]
		}
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no images provided"})
			return
		}

		uploadDir := filepath.Join("uploads", "listings")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Printf("failed to create upload dir: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not prepare upload directory"})
			return
		}

		var paths []string
		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			savePath := filepath.Join(uploadDir, filename)

			if err := c.SaveUploadedFile(file, savePath); err != nil {
				log.Printf("failed to save image %s: %v", file.Filename, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
				return
			}

			paths = append(paths, "/"+filepath.ToSlash(savePath))
		}

		c.JSON(http.StatusOK, gin.H{"paths": paths})
	}
}
