package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jhnthnljyng/rrf-be/internal/models"
)

func CreateProperty(db *sql.DB, ownerID int, req *models.CreatePropertyRequest) (*models.Property, error) {
	availableFrom, err := time.Parse("2006-01-02", req.AvailableFrom)
	if err != nil {
		return nil, err
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	furnishing := req.Furnishing
	if furnishing == "" {
		furnishing = "unfurnished"
	}

	genderPref := req.GenderPref
	if genderPref == "" {
		genderPref = "any"
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	p := &models.Property{}
	err = tx.QueryRow(`
		INSERT INTO properties
			(owner_id, title, description, location, monthly_rent, room_type, furnishing,
			 max_occupants, available_from, gender_preference, bedrooms, bathrooms)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, owner_id, title, description, location, monthly_rent, room_type,
		          furnishing, max_occupants, available_from, gender_preference,
		          bedrooms, bathrooms, is_active, created_at, updated_at`,
		ownerID, req.Title, description, req.Location, req.Price, req.Type, furnishing,
		req.MaxOccupants, availableFrom, genderPref, req.Bedrooms, req.Bathrooms,
	).Scan(
		&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Location, &p.MonthlyRent, &p.RoomType,
		&p.Furnishing, &p.MaxOccupants, &p.AvailableFrom, &p.GenderPref,
		&p.Bedrooms, &p.Bathrooms, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	for _, amenity := range req.Amenities {
		if amenity == "" {
			continue
		}
		_, err := tx.Exec(`
			INSERT INTO property_amenities (property_id, name) VALUES ($1, $2)
			ON CONFLICT (property_id, name) DO NOTHING`,
			p.ID, amenity,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
}

func GetAllProperties(db *sql.DB, roomType, location string) ([]map[string]interface{}, error) {
	query := `
		SELECT id, owner_id, title, description, location, monthly_rent, room_type,
		       furnishing, max_occupants, available_from, gender_preference,
		       bedrooms, bathrooms, is_active, created_at, updated_at
		FROM properties
		WHERE is_active = TRUE`

	args := []interface{}{}
	argIdx := 1

	if roomType != "" {
		query += fmt.Sprintf(" AND room_type = $%d", argIdx)
		args = append(args, roomType)
		argIdx++
	}
	if location != "" {
		query += fmt.Sprintf(" AND location ILIKE $%d", argIdx)
		args = append(args, "%"+location+"%")
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []map[string]interface{}
	for rows.Next() {
		p := &models.Property{}
		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Location, &p.MonthlyRent, &p.RoomType,
			&p.Furnishing, &p.MaxOccupants, &p.AvailableFrom, &p.GenderPref,
			&p.Bedrooms, &p.Bathrooms, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		amenities, err := getPropertyAmenities(db, p.ID)
		if err != nil {
			return nil, err
		}

		images, err := getPropertyImages(db, p.ID)
		if err != nil {
			return nil, err
		}

		properties = append(properties, map[string]interface{}{
			"id":               p.ID,
			"owner_id":         p.OwnerID,
			"title":            p.Title,
			"description":      p.Description,
			"location":         p.Location,
			"price":            p.MonthlyRent,
			"type":             p.RoomType,
			"furnishing":       p.Furnishing,
			"maxOccupants":     p.MaxOccupants,
			"availableFrom":    p.AvailableFrom,
			"genderPreference": p.GenderPref,
			"bedrooms":         p.Bedrooms,
			"bathrooms":        p.Bathrooms,
			"amenities":        amenities,
			"images":           images,
			"created_at":       p.CreatedAt,
		})
	}

	if properties == nil {
		properties = []map[string]interface{}{}
	}

	return properties, nil
}

func GetPropertyByID(db *sql.DB, id int) (map[string]interface{}, error) {
	p := &models.Property{}
	err := db.QueryRow(`
		SELECT id, owner_id, title, description, location, monthly_rent, room_type,
		       furnishing, max_occupants, available_from, gender_preference,
		       bedrooms, bathrooms, is_active, created_at, updated_at
		FROM properties
		WHERE id = $1 AND is_active = TRUE`, id,
	).Scan(
		&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Location, &p.MonthlyRent, &p.RoomType,
		&p.Furnishing, &p.MaxOccupants, &p.AvailableFrom, &p.GenderPref,
		&p.Bedrooms, &p.Bathrooms, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	amenities, err := getPropertyAmenities(db, p.ID)
	if err != nil {
		return nil, err
	}

	images, err := getPropertyImages(db, p.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":               p.ID,
		"owner_id":         p.OwnerID,
		"title":            p.Title,
		"description":      p.Description,
		"location":         p.Location,
		"price":            p.MonthlyRent,
		"type":             p.RoomType,
		"furnishing":       p.Furnishing,
		"maxOccupants":     p.MaxOccupants,
		"availableFrom":    p.AvailableFrom,
		"genderPreference": p.GenderPref,
		"bedrooms":         p.Bedrooms,
		"bathrooms":        p.Bathrooms,
		"amenities":        amenities,
		"images":           images,
		"created_at":       p.CreatedAt,
	}, nil
}

func getPropertyAmenities(db *sql.DB, propertyID int) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM property_amenities WHERE property_id = $1`, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amenities []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		amenities = append(amenities, name)
	}
	if amenities == nil {
		amenities = []string{}
	}
	return amenities, nil
}

func getPropertyImages(db *sql.DB, propertyID int) ([]string, error) {
	rows, err := db.Query(`SELECT url FROM property_images WHERE property_id = $1 ORDER BY display_order`, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		images = append(images, url)
	}
	if images == nil {
		images = []string{}
	}
	return images, nil
}

func SavePropertyImages(db *sql.DB, propertyID int, urls []string) error {
	for i, url := range urls {
		isPrimary := i == 0
		_, err := db.Exec(`
			INSERT INTO property_images (property_id, url, is_primary, display_order)
			VALUES ($1, $2, $3, $4)`,
			propertyID, url, isPrimary, i,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
