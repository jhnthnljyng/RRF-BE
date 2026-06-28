package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jhnthnljyng/rrf-be/internal/repository"
)

func GetFavourites(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		listings, err := repository.GetFavourites(db, userID)
		if err != nil {
			log.Printf("get favourites error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch favourites"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"favourites": listings, "total": len(listings)})
	}
}

func AddFavourite(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		var body struct {
			ListingID int `json:"listing_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "listing_id is required"})
			return
		}

		err := repository.AddFavourite(db, userID, body.ListingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
				return
			}
			if err == repository.ErrAlreadyFavourited {
				c.JSON(http.StatusConflict, gin.H{"error": "listing already in favourites"})
				return
			}
			log.Printf("add favourite error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add favourite"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "added to favourites"})
	}
}

func RemoveFavourite(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		listingID, err := strconv.Atoi(c.Param("listing_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
			return
		}

		err = repository.RemoveFavourite(db, userID, listingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "favourite not found"})
				return
			}
			log.Printf("remove favourite error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove favourite"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "removed from favourites"})
	}
}
