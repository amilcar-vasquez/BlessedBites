// data/recommendations.go
package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

// GetTopRecommendationsByUser returns top N recommended items for a user, ordered by score
func (r RecommendationModel) GetTopRecommendationsByUser(userID, limit int) ([]int, error) {
	query := `
	SELECT menu_item_id
	FROM order_items
	JOIN orders ON order_items.order_id = orders.id
	WHERE orders.user_id = $1
	GROUP BY menu_item_id
	ORDER BY COUNT(*) DESC
	LIMIT $2`

	rows, err := r.DB.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("fetching top recommendations: %w", err)
	}
	defer rows.Close()

	var itemIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		itemIDs = append(itemIDs, id)
	}
	return itemIDs, nil
}

// recommends a menu item based on general popularity
func (r RecommendationModel) RecommendPopularItem() ([]Recommendation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Fetch popular items
	popularQuery := `SELECT id as menu_item_id, order_count
FROM menu_items
ORDER BY order_count DESC
LIMIT 5`
	rows, err := r.DB.QueryContext(ctx, popularQuery)
	if err != nil {
		return nil, fmt.Errorf("fetching popular items: %w", err)
	}
	defer rows.Close()

	var recs []Recommendation
	for rows.Next() {
		var itemID int
		var count int
		if err := rows.Scan(&itemID, &count); err != nil {
			return nil, fmt.Errorf("scanning popular item: %w", err)
		}
		rec := Recommendation{
			MenuItemID: itemID,
			Reason:     "popular item",
			Score:      float64(count),
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

// Returns just popular menu_item_ids
func (r RecommendationModel) GetPopularItemIDs(limit int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT oi.menu_item_id
	FROM order_items oi
	JOIN menu_items mi ON oi.menu_item_id = mi.id
	GROUP BY oi.menu_item_id
	ORDER BY SUM(order_count) DESC, oi.menu_item_id ASC
	LIMIT $1
	`

	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fetching popular item IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
