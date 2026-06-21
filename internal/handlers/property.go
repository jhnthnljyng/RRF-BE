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
	"github.com/jhnthnljyng/rrf-be/internal/models"
	"github.com/jhnthnljyng/rrf-be/internal/repository"
)

func GetListings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomType := c.Query("type")
		location := c.Query("location")

		properties, err := repository.GetAllProperties(db, roomType, location)
		if err != nil {
			log.Printf("get listings error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch listings"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"listings": properties,
			"total":    len(properties),
		})
	}
}

func CreateProperty(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreatePropertyRequest
		if err := c.ShouldBind(&req); err != nil {
			log.Printf("bind error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ownerID := c.GetInt("user_id")

		property, err := repository.CreateProperty(db, ownerID, &req)
		if err != nil {
			log.Printf("create property error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create property"})
			return
		}

		multipartForm, err := c.MultipartForm()
		if err == nil {
			log.Printf("multipart file keys received:")
			for key, files := range multipartForm.File {
				log.Printf("  key=%q count=%d", key, len(files))
			}
		}
		if err == nil && multipartForm.File["images"] != nil {
			uploadDir := filepath.Join("uploads", "properties", fmt.Sprintf("%d", property.ID))
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				log.Printf("failed to create upload dir: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save images"})
				return
			}

			var savedURLs []string
			for _, file := range multipartForm.File["images"] {
				ext := filepath.Ext(file.Filename)
				filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
				savePath := filepath.Join(uploadDir, filename)

				if err := c.SaveUploadedFile(file, savePath); err != nil {
					log.Printf("failed to save image %s: %v", file.Filename, err)
					continue
				}

				savedURLs = append(savedURLs, "/"+filepath.ToSlash(savePath))
			}

			if len(savedURLs) > 0 {
				if err := repository.SavePropertyImages(db, property.ID, savedURLs); err != nil {
					log.Printf("failed to save image records: %v", err)
				}
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "property created successfully",
			"property": property,
		})
	}
}
