package models

import "time"

type Property struct {
	ID            int        `json:"id"`
	OwnerID       int        `json:"owner_id"`
	Title         string     `json:"title"`
	Description   *string    `json:"description,omitempty"`
	Location      string     `json:"location"`
	MonthlyRent   float64    `json:"price"`
	RoomType      string     `json:"type"`
	Furnishing    string     `json:"furnishing"`
	MaxOccupants  int        `json:"maxOccupants"`
	AvailableFrom time.Time  `json:"availableFrom"`
	GenderPref    string     `json:"genderPreference"`
	Bedrooms      *int       `json:"bedrooms,omitempty"`
	Bathrooms     *int       `json:"bathrooms,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreatePropertyRequest struct {
	Title         string   `form:"title"            binding:"required"`
	Description   string   `form:"description"`
	Location      string   `form:"location"         binding:"required"`
	Price         float64  `form:"price"            binding:"required,gt=0"`
	Type          string   `form:"type"             binding:"required,oneof=single shared studio apartment house room whole_unit"`
	Furnishing    string   `form:"furnishing"       binding:"required,oneof=furnished unfurnished partial"`
	MaxOccupants  int      `form:"maxOccupants"     binding:"required,min=1"`
	AvailableFrom string   `form:"availableFrom"    binding:"required"`
	GenderPref    string   `form:"genderPreference" binding:"omitempty,oneof=male female any"`
	Bedrooms      *int     `form:"bedrooms"`
	Bathrooms     *int     `form:"bathrooms"`
	Amenities     []string `form:"amenities[]"`
}
