package utils

import (
	"fmt"
	"math/rand"
)

func RandomPhone() string {
	// Generate a random 4-digit number
	num := rand.Intn(1000) // 0 to 9999999
	// Format it as a string with leading zeros if necessary
	phone := "600" + formatNumber(num)
	return phone
}

// formatNumber formats an integer as a 4-digit string with leading zeros.
func formatNumber(n int) string {
	return fmt.Sprintf("%04d", n) // Format as 4 digits with leading zeros
}
