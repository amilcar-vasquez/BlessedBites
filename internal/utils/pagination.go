package utils

// write a function for pagination
func Paginate[T any](items []T, page, itemsPerPage int) ([]T, int, int) {
	totalItems := len(items)
	if itemsPerPage <= 0 {
		itemsPerPage = 9 // Default items per page
	}
	if page <= 0 {
		page = 1 // Default to first page
	}

	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage // Calculate total pages
	if page > totalPages {
		page = totalPages // Adjust page if it exceeds total pages
	}
	start := (page - 1) * itemsPerPage
	if start >= totalItems {
		start = totalItems
	}
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}
	return items[start:end], totalItems, totalPages
}

// helper add function
func Add(a, b int) int {
	return a + b
}

// helper subtract function
func Subtract(a, b int) int {
	return a - b
}

// helper until function
func Until(n int) []int {
	var res []int
	for i := 0; i < n; i++ {
		res = append(res, i)
	}
	return res
}
