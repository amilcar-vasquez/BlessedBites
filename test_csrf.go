package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func main() {
	// Create a cookie jar to maintain sessions - this ensures fresh cookies
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	fmt.Println("Starting with fresh cookie jar (no old sessions)")

	// 1. GET the login page to get CSRF token
	fmt.Println("=== Step 1: Getting login page ===")
	resp, err := client.Get("http://localhost:4000/login")
	if err != nil {
		fmt.Printf("Error getting login page: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)

	// Print cookies
	fmt.Println("Cookies received:")
	for _, cookie := range resp.Cookies() {
		fmt.Printf("  %s = %s\n", cookie.Name, cookie.Value)
	}

	// Read page content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Extract CSRF token from form
	csrfPattern := regexp.MustCompile(`name="csrf_token" value="([^"]+)"`)
	matches := csrfPattern.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		fmt.Println("ERROR: Could not find CSRF token in form")
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500]
		}
		fmt.Printf("Page content sample:\n%s\n", bodyStr)
		return
	}

	csrfToken := matches[1]
	fmt.Printf("CSRF token found: %s\n", csrfToken[:20]+"...")

	// 2. POST login with CSRF token
	fmt.Println("\n=== Step 2: Submitting login form ===")

	formData := url.Values{}
	formData.Set("email", "admin@test.com")
	formData.Set("password", "password123") // Use appropriate test password
	formData.Set("csrf_token", csrfToken)

	postResp, err := client.PostForm("http://localhost:4000/login", formData)
	if err != nil {
		fmt.Printf("Error posting login: %v\n", err)
		return
	}
	defer postResp.Body.Close()

	fmt.Printf("Status: %s\n", postResp.Status)
	fmt.Printf("Location header: %s\n", postResp.Header.Get("Location"))

	if postResp.StatusCode == 403 {
		// Read error response
		errorBody, _ := io.ReadAll(postResp.Body)
		fmt.Printf("CSRF Error Response: %s\n", string(errorBody))
	} else {
		fmt.Println("SUCCESS: CSRF validation passed!")
	}
}
