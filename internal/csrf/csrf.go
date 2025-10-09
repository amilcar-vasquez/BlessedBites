package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const (
	tokenLength = 32
	cookieName  = "csrf_token"
	fieldName   = "csrf_token"
	headerName  = "X-CSRF-Token"
)

type Config struct {
	Key    []byte
	Secure bool
	MaxAge int
}

// GenerateToken creates a new CSRF token
func GenerateToken(key []byte) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, tokenLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Create timestamp (helps with expiry)
	timestamp := time.Now().Unix()

	// Create message: random_bytes + timestamp
	message := fmt.Sprintf("%s:%d", base64.StdEncoding.EncodeToString(randomBytes), timestamp)

	// Create HMAC signature
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	signature := h.Sum(nil)

	// Final token: message + signature (base64 encoded)
	token := fmt.Sprintf("%s.%s", message, base64.StdEncoding.EncodeToString(signature))
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

// ValidateToken checks if a CSRF token is valid
func ValidateToken(token string, key []byte, maxAge int) bool {
	// Decode base64 token
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	tokenStr := string(tokenBytes)

	// Split message and signature
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 2 {
		return false
	}

	message := parts[0]
	expectedSig, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Verify HMAC signature
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	actualSig := h.Sum(nil)

	if !hmac.Equal(expectedSig, actualSig) {
		return false
	}

	// Check expiry if maxAge is set
	if maxAge > 0 {
		msgParts := strings.Split(message, ":")
		if len(msgParts) != 2 {
			return false
		}

		var timestamp int64
		if _, err := fmt.Sscanf(msgParts[1], "%d", &timestamp); err != nil {
			return false
		}

		if time.Now().Unix()-timestamp > int64(maxAge) {
			return false
		}
	}

	return true
}

// Middleware creates CSRF protection middleware
func Middleware(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip GET, HEAD, OPTIONS, TRACE
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" || r.Method == "TRACE" {
				next.ServeHTTP(w, r)
				return
			}

			// Get token from form or header
			token := r.FormValue(fieldName)
			if token == "" {
				token = r.Header.Get(headerName)
			}

			// Validate token
			if !ValidateToken(token, config.Key, config.MaxAge) {
				http.Error(w, "Forbidden - CSRF token invalid", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Token returns a CSRF token for the current request
func Token(r *http.Request, key []byte) string {
	// Try to get existing token from cookie first
	if cookie, err := r.Cookie(cookieName); err == nil {
		if ValidateToken(cookie.Value, key, 3600) { // 1 hour validity
			return cookie.Value
		}
	}

	// Generate new token
	token, err := GenerateToken(key)
	if err != nil {
		return ""
	}

	return token
}

// SetTokenCookie sets the CSRF token as a cookie
func SetTokenCookie(w http.ResponseWriter, token string, secure bool) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // Needs to be accessible by JavaScript if needed
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600, // 1 hour
	}
	http.SetCookie(w, cookie)
}

// TemplateField returns HTML for embedding in forms
func TemplateField(r *http.Request, key []byte, w http.ResponseWriter, secure bool) template.HTML {
	token := Token(r, key)

	// Set cookie for future requests
	SetTokenCookie(w, token, secure)

	return template.HTML(`<input type="hidden" name="` + fieldName + `" value="` + token + `">`)
}
