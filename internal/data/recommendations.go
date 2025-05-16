// data/recommendations.go
package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
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

// analyzes past orders and creates recommendations based on frequency of a certain user's orders
func (r RecommendationModel) RecommendUserItem(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Fetch order history for the user
	orderQuery := `SELECT id FROM orders WHERE user_id = $1`
	orderRows, err := r.DB.QueryContext(ctx, orderQuery, userID)
	if err != nil {
		return fmt.Errorf("fetching orders: %w", err)
	}
	defer orderRows.Close()

	var orderIDs []int
	for orderRows.Next() {
		var id int
		if err := orderRows.Scan(&id); err != nil {
			return fmt.Errorf("scanning order ID: %w", err)
		}
		orderIDs = append(orderIDs, id)
	}
	if len(orderIDs) == 0 {
		return nil // no orders, no recommendations
	}

	// Step 2: Aggregate item frequency
	freqMap := make(map[int]int) // key: menuItemID, value: total quantity
	itemQuery := `SELECT menu_item_id, quantity FROM order_items WHERE order_id = ANY($1)`
	itemRows, err := r.DB.QueryContext(ctx, itemQuery, pq.Array(orderIDs))
	if err != nil {
		return fmt.Errorf("fetching order items: %w", err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var itemID, qty int
		if err := itemRows.Scan(&itemID, &qty); err != nil {
			return fmt.Errorf("scanning order item: %w", err)
		}
		freqMap[itemID] += qty
	}

	// Step 3: Insert recommendations based on frequency
	for menuItemID, count := range freqMap {
		rec := Recommendation{
			UserID:     userID,
			MenuItemID: menuItemID,
			Reason:     "frequent order",
			Score:      float64(count), // simple frequency as score for now
		}
		err := r.Insert(rec)
		if err != nil {
			// optional: log and continue instead of failing
			return fmt.Errorf("inserting recommendation: %w", err)
		}
	}

	return nil
}

// recommends a menu item based on general popularity
func (r RecommendationModel) RecommendPopularItem() ([]Recommendation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Fetch popular items
	popularQuery := `SELECT menu_item_id, COUNT(*) as order_count FROM order_items GROUP BY menu_item_id ORDER BY order_count DESC LIMIT 5`
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
	SELECT menu_item_id
	FROM order_items
	GROUP BY menu_item_id
	ORDER BY COUNT(*) DESC
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
