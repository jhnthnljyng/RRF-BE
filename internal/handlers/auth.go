package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhnthnljyng/rrf-be/internal/models"
	"github.com/jhnthnljyng/rrf-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, passwordHash, err := repository.GetUserByEmail(db, req.Email)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
				return
			}
			log.Printf("login error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process login"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"role":    user.Role,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Printf("jwt sign error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"token":   tokenString,
			"user":    user,
		})
	}
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process password"})
			return
		}

		user, err := repository.CreateUser(db, &req, string(hash))
		if err != nil {
			if errors.Is(err, repository.ErrEmailTaken) {
				c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
				return
			}
			if errors.Is(err, repository.ErrUsernameTaken) {
				c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
				return
			}
			log.Printf("register error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "registration successful",
			"user":    user,
		})
	}
}
