package models

import "time"

type UserSocial struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

type User struct {
	ID               int           `json:"id"`
	Email            string        `json:"email"`
	FullName         string        `json:"full_name"`
	Username         string        `json:"username"`
	Phone            *string       `json:"phone,omitempty"`
	AvatarURL        *string       `json:"avatar_url,omitempty"`
	Role             string        `json:"role"`
	Smoking          int           `json:"smoking"`
	PetOwner         int           `json:"pet_owner"`
	PetFriendly      int           `json:"pet_friendly"`
	Bio              *string       `json:"bio,omitempty"`
	Gender           *string       `json:"gender,omitempty"`
	Occupation       *string       `json:"occupation,omitempty"`
	Nationality      *string       `json:"nationality,omitempty"`
	CookingFrequency *string       `json:"cooking_frequency,omitempty"`
	Socials          []UserSocial  `json:"socials"`
	IsActive         bool          `json:"is_active"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FullName         string       `json:"full_name"`
	Username         string       `json:"username"`
	Phone            string       `json:"phone"`
	AvatarURL        string       `json:"avatar_url"`
	Bio              string       `json:"bio"`
	Gender           string       `json:"gender"`
	Occupation       string       `json:"occupation"`
	Nationality      string       `json:"nationality"`
	CookingFrequency string       `json:"cooking_frequency" binding:"omitempty,oneof=never rarely sometimes often always"`
	Smoking          *int         `json:"smoking"`
	PetOwner         *int         `json:"pet_owner"`
	PetFriendly      *int         `json:"pet_friendly"`
	Socials          []UserSocial `json:"socials"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	FullName string `json:"fullname" binding:"required"`
	Username string `json:"username"  binding:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"omitempty,oneof=owner tenant"`
}
