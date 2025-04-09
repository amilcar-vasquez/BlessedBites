// data/recommendations.go
package data

import (
	"database/sql"
)

type Recommendation struct {
	ID         int
	UserID     int
	MenuItemID int
	Reason     string // e.g., "frequent order", "time-based"
	Score      float64
	CreatedAt  string // time.Time preferred for real systems
}

type RecommendationModel struct {
	DB *sql.DB
}

func (r RecommendationModel) Insert(rec Recommendation) error {
	query := `INSERT INTO recommendations (user_id, menu_item_id, reason, score, created_at)
		VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.DB.Exec(query, rec.UserID, rec.MenuItemID, rec.Reason, rec.Score)
	return err
}

func (r RecommendationModel) GetByUser(userID int) ([]Recommendation, error) {
	query := `SELECT id, user_id, menu_item_id, reason, score, created_at FROM recommendations WHERE user_id = $1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recs []Recommendation
	for rows.Next() {
		var rec Recommendation
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.MenuItemID, &rec.Reason, &rec.Score, &rec.CreatedAt)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

func NewRecommendationModel(db *sql.DB) RecommendationModel {
	return RecommendationModel{DB: db}
}
