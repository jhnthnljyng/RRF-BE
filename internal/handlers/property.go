package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

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

func GetListingByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
			return
		}

		property, err := repository.GetPropertyByID(db, id)
		if err != nil {
			log.Printf("get listing error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch listing"})
			return
		}
		if property == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"listing": property})
	}
}

func CreateProperty(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreatePropertyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
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

		log.Printf("images received in request: %v (count: %d)", req.Images, len(req.Images))
		if len(req.Images) > 0 {
			if err := repository.SavePropertyImages(db, property.ID, req.Images); err != nil {
				log.Printf("failed to save image records: %v", err)
			} else {
				log.Printf("saved %d images for property %d", len(req.Images), property.ID)
			}
		} else {
			log.Printf("no images to save for property %d", property.ID)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "property created successfully",
			"property": property,
		})
	}
}
