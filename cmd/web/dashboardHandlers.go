// file: cmd/web/dashboardHandlers.go
package main

import (
	"net/http"
	"time"
)

// dashboard page handler
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	data := app.addDefaultData(NewTemplateData(), w, r)
	data.Title = "Dashboard"
	data.HeaderText = "Dashboard"

	// Get the total number of orders for today
	totalOrders, err := app.Order.Count()
	if err != nil {
		app.logger.Error("Error getting total orders", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.TotalOrders = totalOrders

	// Get the daily sales list for today by calling the DailySales method from the OrderItem model (uses slice of OrderItem)
	today := time.Now().Format("2006-01-02")
	dailySales, err := app.OrderItem.DailySales(today)
	if err != nil {
		app.logger.Error("Error getting daily sales", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.DailySales = dailySales

	// Fetch last 7 days sales data
	last7DaysSales, err := app.OrderItem.Last7DaysSales()
	if err != nil {
		app.logger.Error("Error getting last 7 days sales", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.Last7DaysSales = last7DaysSales

	var chartLabels []string
	var chartData []float64
	for _, sale := range data.Last7DaysSales {
		chartLabels = append(chartLabels, sale.Date)
		chartData = append(chartData, sale.Amount)
	}
	data.ChartLabels = chartLabels
	data.ChartData = chartData

	// Fetch top 5 popular menu items
	topItems, err := app.MenuItem.GetTopPopularItems()
	if err != nil {
		app.logger.Error("Error fetching top 5 menu items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.Top5MenuItems = topItems

	// Render the dashboard template with the data
	err = app.render(w, http.StatusOK, "dashboard.tmpl", data)
	if err != nil {
		app.logger.Error("Error rendering dashboard template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Log the successful rendering of the dashboard
	app.logger.Info("Dashboard rendered successfully", "totalOrders", totalOrders, "dailySales", dailySales)
}
