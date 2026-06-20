package repository

import (
	"database/sql"
	"errors"

	"github.com/jhnthnljyng/rrf-be/internal/models"
)

var ErrEmailTaken    = errors.New("email already in use")
var ErrUsernameTaken = errors.New("username already in use")
var ErrNotFound      = errors.New("user not found")

func GetUserByEmail(db *sql.DB, email string) (*models.User, string, error) {
	var passwordHash string
	user := &models.User{}
	err := db.QueryRow(`
		SELECT id, email, password_hash, full_name, username, phone, avatar_url, role, is_active, created_at, updated_at
		FROM users WHERE email = $1 AND is_active = TRUE`, email,
	).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, &user.Username,
		&user.Phone, &user.AvatarURL, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, "", ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	return user, passwordHash, nil
}

func CreateUser(db *sql.DB, req *models.RegisterRequest, passwordHash string) (*models.User, error) {
	var exists bool

	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email).Scan(&exists); err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailTaken
	}

	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, req.Username).Scan(&exists); err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	role := req.Role
	if role == "" {
		role = "tenant"
	}

	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	user := &models.User{}
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, full_name, username, phone, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, full_name, username, phone, avatar_url, role, is_active, created_at, updated_at`,
		req.Email, passwordHash, req.FullName, req.Username, phone, role,
	).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Username,
		&user.Phone, &user.AvatarURL, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
