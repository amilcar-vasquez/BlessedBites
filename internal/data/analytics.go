// data/analytics.go
package data

import (
	"database/sql"
)

type Analytics struct {
	ID         int
	MenuItemID int
	Timestamp  string // Can change to time.Time if preferred
	Action     string // e.g., "view", "click", "order"
	Meta       string // JSON or encoded context info
}

type AnalyticsModel struct {
	DB *sql.DB
}

func (a AnalyticsModel) Insert(analytic Analytics) error {
	query := `INSERT INTO analytics (menu_item_id, timestamp, action, meta) VALUES ($1, NOW(), $2, $3)`
	_, err := a.DB.Exec(query, analytic.MenuItemID, analytic.Action, analytic.Meta)
	return err
}

func NewAnalyticsModel(db *sql.DB) AnalyticsModel {
	return AnalyticsModel{DB: db}
}
