package repository

import (
	"database/sql"
	"errors"

	"github.com/jhnthnljyng/rrf-be/internal/models"
)

var ErrEmailTaken    = errors.New("email already in use")
var ErrUsernameTaken = errors.New("username already in use")
var ErrNotFound      = errors.New("user not found")

func UpdateProfile(db *sql.DB, userID int, req *models.UpdateProfileRequest) (*models.User, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var phone, avatarURL, bio, gender, occupation, nationality, cookingFreq *string
	if req.Phone != "" {
		phone = &req.Phone
	}
	if req.AvatarURL != "" {
		avatarURL = &req.AvatarURL
	}
	if req.Bio != "" {
		bio = &req.Bio
	}
	if req.Gender != "" {
		gender = &req.Gender
	}
	if req.Occupation != "" {
		occupation = &req.Occupation
	}
	if req.Nationality != "" {
		nationality = &req.Nationality
	}
	if req.CookingFrequency != "" {
		cookingFreq = &req.CookingFrequency
	}

	_, err = tx.Exec(`
		UPDATE users SET
			full_name         = COALESCE(NULLIF($1, ''), full_name),
			username          = COALESCE(NULLIF($2, ''), username),
			phone             = COALESCE($3, phone),
			avatar_url        = COALESCE($4, avatar_url),
			bio               = $5,
			gender            = $6,
			occupation        = $7,
			nationality       = $8,
			cooking_frequency = COALESCE($9::cooking_frequency_type, cooking_frequency),
			smoking           = COALESCE($10, smoking),
			pet_owner         = COALESCE($11, pet_owner),
			pet_friendly      = COALESCE($12, pet_friendly),
			updated_at        = NOW()
		WHERE id = $13`,
		req.FullName, req.Username, phone, avatarURL, bio, gender, occupation, nationality,
		cookingFreq, req.Smoking, req.PetOwner, req.PetFriendly, userID,
	)
	if err != nil {
		return nil, err
	}

	if req.Socials != nil {
		if _, err := tx.Exec(`DELETE FROM user_socials WHERE user_id = $1`, userID); err != nil {
			return nil, err
		}
		for _, s := range req.Socials {
			if s.Platform == "" || s.URL == "" {
				continue
			}
			if _, err := tx.Exec(
				`INSERT INTO user_socials (user_id, platform, url) VALUES ($1, $2, $3)`,
				userID, s.Platform, s.URL,
			); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return GetUserByID(db, userID)
}

func UpdateAvatarURL(db *sql.DB, userID int, avatarURL string) error {
	_, err := db.Exec(`UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`, avatarURL, userID)
	return err
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(`
		SELECT id, email, full_name, username, phone, avatar_url, role,
		       smoking, pet_owner, pet_friendly, bio, gender, occupation,
		       nationality, cooking_frequency, is_active, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = TRUE`, id,
	).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Username, &user.Phone, &user.AvatarURL, &user.Role,
		&user.Smoking, &user.PetOwner, &user.PetFriendly, &user.Bio, &user.Gender, &user.Occupation,
		&user.Nationality, &user.CookingFrequency, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT platform, url FROM user_socials WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user.Socials = []models.UserSocial{}
	for rows.Next() {
		var s models.UserSocial
		if err := rows.Scan(&s.Platform, &s.URL); err != nil {
			return nil, err
		}
		user.Socials = append(user.Socials, s)
	}

	return user, nil
}

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
