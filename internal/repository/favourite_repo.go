package repository

import (
	"database/sql"
	"errors"
	"strings"
)

var ErrAlreadyFavourited = errors.New("already favourited")

func GetFavourites(db *sql.DB, userID int) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT p.id, p.owner_id, p.title, p.description, p.location, p.monthly_rent,
		       p.room_type, p.furnishing, p.max_occupants, p.available_from,
		       p.gender_preference, p.bedrooms, p.bathrooms, p.status, p.tenant_id,
		       p.is_active, p.created_at,
		       u.full_name, u.phone,
		       f.created_at AS favourited_at
		FROM favourites f
		JOIN properties p ON p.id = f.listing_id
		JOIN users u ON u.id = p.owner_id
		WHERE f.user_id = $1 AND p.is_active = TRUE
		ORDER BY f.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []map[string]interface{}
	for rows.Next() {
		var (
			id, ownerID                       int
			title, location, roomType         string
			furnishing, genderPref, status    string
			ownerName                         string
			description, ownerPhone           *string
			monthlyRent                       float64
			maxOccupants, bedrooms, bathrooms *int
			tenantID                          *int
			isActive                          bool
			createdAt, favouritedAt           interface{}
			availableFrom                     interface{}
		)
		if err := rows.Scan(
			&id, &ownerID, &title, &description, &location, &monthlyRent,
			&roomType, &furnishing, &maxOccupants, &availableFrom,
			&genderPref, &bedrooms, &bathrooms, &status, &tenantID,
			&isActive, &createdAt,
			&ownerName, &ownerPhone,
			&favouritedAt,
		); err != nil {
			return nil, err
		}

		amenities, err := getPropertyAmenities(db, id)
		if err != nil {
			return nil, err
		}
		images, err := getPropertyImages(db, id)
		if err != nil {
			return nil, err
		}

		listings = append(listings, map[string]interface{}{
			"id":               id,
			"owner_id":         ownerID,
			"owner_name":       ownerName,
			"owner_phone":      ownerPhone,
			"title":            title,
			"description":      description,
			"location":         location,
			"price":            monthlyRent,
			"type":             roomType,
			"furnishing":       furnishing,
			"maxOccupants":     maxOccupants,
			"availableFrom":    availableFrom,
			"genderPreference": genderPref,
			"bedrooms":         bedrooms,
			"bathrooms":        bathrooms,
			"status":           status,
			"tenant_id":        tenantID,
			"is_active":        isActive,
			"amenities":        amenities,
			"images":           images,
			"created_at":       createdAt,
			"favourited_at":    favouritedAt,
		})
	}

	if listings == nil {
		listings = []map[string]interface{}{}
	}
	return listings, nil
}

func AddFavourite(db *sql.DB, userID, listingID int) error {
	var exists bool
	if err := db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM properties WHERE id = $1 AND is_active = TRUE)`, listingID,
	).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	_, err := db.Exec(
		`INSERT INTO favourites (user_id, listing_id) VALUES ($1, $2)`,
		userID, listingID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return ErrAlreadyFavourited
		}
		return err
	}
	return nil
}

func RemoveFavourite(db *sql.DB, userID, listingID int) error {
	result, err := db.Exec(
		`DELETE FROM favourites WHERE user_id = $1 AND listing_id = $2`,
		userID, listingID,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
