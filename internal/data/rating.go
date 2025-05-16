package data

import (
	"context"
	"database/sql"
	"time"
)

type Rating struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	MenuItemID int       `json:"menu_item_id"`
	Rating     int       `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
}

type RatingModel struct {
	DB *sql.DB
}

func NewRatingModel(db *sql.DB) RatingModel {
	return RatingModel{DB: db}
}

func (r RatingModel) Insert(rating *Rating) error {
	query := `
		INSERT INTO ratings (user_id, menu_item_id, rating)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, rating.UserID, rating.MenuItemID, rating.Rating).
		Scan(&rating.ID, &rating.CreatedAt)
}

func (r RatingModel) GetAverageRating(menuItemID int) (float64, error) {
	query := `SELECT COALESCE(AVG(rating), 0) FROM ratings WHERE menu_item_id = $1`

	var avg float64
	err := r.DB.QueryRow(query, menuItemID).Scan(&avg)
	return avg, err
}

// validate rating to include only 1-5
func (r RatingModel) ValidateRating(rating int) bool {
	if rating < 1 || rating > 5 {
		return false
	}
	return true
}
