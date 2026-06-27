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
	MaxOccupants  *int       `json:"maxOccupants,omitempty"`
	AvailableFrom time.Time  `json:"availableFrom"`
	GenderPref    string     `json:"genderPreference"`
	Bedrooms      *int       `json:"bedrooms,omitempty"`
	Bathrooms     *int       `json:"bathrooms,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreatePropertyRequest struct {
	Title         string   `json:"title"            binding:"required"`
	Description   string   `json:"description"`
	Location      string   `json:"location"         binding:"required"`
	Price         float64  `json:"price"            binding:"required,gt=0"`
	Type          string   `json:"type"             binding:"required,oneof=single shared studio apartment house room whole_unit looking_for_roommate"`
	Furnishing    string   `json:"furnishing"       binding:"required,oneof=furnished unfurnished partial"`
	MaxOccupants  *int     `json:"maxOccupants"`
	AvailableFrom string   `json:"availableFrom"    binding:"required"`
	GenderPref    string   `json:"genderPreference" binding:"omitempty,oneof=male female any"`
	Bedrooms      *int     `json:"bedrooms"`
	Bathrooms     *int     `json:"bathrooms"`
	Amenities     []string `json:"amenities"`
	Images        []string `json:"images"`
}
