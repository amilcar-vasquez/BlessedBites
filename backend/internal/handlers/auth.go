package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	appjwt "github.com/amilcar-vasquez/blessed-bites/backend/pkg/jwt"
	"github.com/amilcar-vasquez/blessed-bites/backend/pkg/responses"
	"github.com/amilcar-vasquez/blessed-bites/internal/data"
	"github.com/amilcar-vasquez/blessed-bites/internal/mailer"
	"golang.org/x/crypto/bcrypt"
)

// revocationStore holds a set of revoked token JTIs.
// Entries are pruned whenever a new token is revoked.
// For multi-instance deployments, replace with a shared store (Redis, etc.).
type revocationStore struct {
	mu      sync.Mutex
	entries map[string]time.Time // jti -> expiry
}

func newRevocationStore() *revocationStore {
	return &revocationStore{entries: make(map[string]time.Time)}
}

func (s *revocationStore) Revoke(jti string, expiry time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[jti] = expiry
	now := time.Now()
	for id, exp := range s.entries {
		if now.After(exp) {
			delete(s.entries, id)
		}
	}
}

func (s *revocationStore) IsRevoked(jti string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.entries[jti]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(s.entries, jti)
		return false
	}
	return true
}

type AuthHandler struct {
	Users         *data.UserModel
	JWTSecret     []byte
	RefreshSecret []byte
	CookieSecure  bool
	Mailer        *mailer.Mailer
	revoked       *revocationStore
}

func NewAuthHandler(db *sql.DB, jwtSecret, refreshSecret []byte, cookieSecure bool, mailer *mailer.Mailer) *AuthHandler {
	return &AuthHandler{
		Users:         &data.UserModel{DB: db},
		JWTSecret:     jwtSecret,
		RefreshSecret: refreshSecret,
		CookieSecure:  cookieSecure,
		Mailer:        mailer,
		revoked:       newRevocationStore(),
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	PhoneNo  string `json:"phone_no"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordSubmission struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	user, err := h.Users.GetByEmail(input.Email)
	if err != nil {
		responses.Unauthorized(w, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		responses.Unauthorized(w, "invalid email or password")
		return
	}

	accessToken, err := appjwt.GenerateWithUse(h.JWTSecret, user.ID, user.Role, 30*time.Minute, "access")
	if err != nil {
		responses.InternalServerError(w, "failed to issue token")
		return
	}

	refreshToken, err := appjwt.GenerateWithUse(h.RefreshSecret, user.ID, user.Role, 7*24*time.Hour, "refresh")
	if err != nil {
		responses.InternalServerError(w, "failed to issue refresh token")
		return
	}

	h.setRefreshCookie(w, refreshToken)

	responses.JSON(w, http.StatusOK, map[string]any{
		"token": accessToken,
		"user": map[string]any{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var input SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		responses.InternalServerError(w, "failed to process password")
		return
	}

	user := &data.User{
		Email:    input.Email,
		FullName: input.FullName,
		PhoneNo:  input.PhoneNo,
		Password: string(hash),
		Role:     "customer",
	}
	if err := h.Users.Insert(user); err != nil {
		responses.BadRequest(w, "failed to create user")
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]any{"id": user.ID})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("bb_refresh"); err == nil && strings.TrimSpace(cookie.Value) != "" {
		if claims, err := appjwt.Parse(h.RefreshSecret, cookie.Value, "refresh"); err == nil {
			if claims.ExpiresAt != nil {
				h.revoked.Revoke(claims.ID, claims.ExpiresAt.Time)
			}
		}
	}
	h.clearRefreshCookie(w)
	responses.JSON(w, http.StatusOK, map[string]any{"message": "logged out"})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("bb_refresh")
	if err != nil || strings.TrimSpace(refreshCookie.Value) == "" {
		responses.Unauthorized(w, "missing refresh token")
		return
	}

	claims, err := appjwt.Parse(h.RefreshSecret, refreshCookie.Value, "refresh")
	if err != nil {
		responses.Unauthorized(w, "invalid refresh token")
		return
	}

	if h.revoked.IsRevoked(claims.ID) {
		responses.Unauthorized(w, "refresh token has been revoked")
		return
	}

	// Revoke the old refresh token (rotation)
	if claims.ExpiresAt != nil {
		h.revoked.Revoke(claims.ID, claims.ExpiresAt.Time)
	}

	accessToken, err := appjwt.GenerateWithUse(h.JWTSecret, claims.UserID, claims.Role, 30*time.Minute, "access")
	if err != nil {
		responses.InternalServerError(w, "failed to issue access token")
		return
	}

	newRefresh, err := appjwt.GenerateWithUse(h.RefreshSecret, claims.UserID, claims.Role, 7*24*time.Hour, "refresh")
	if err != nil {
		responses.InternalServerError(w, "failed to issue refresh token")
		return
	}

	h.setRefreshCookie(w, newRefresh)
	responses.JSON(w, http.StatusOK, map[string]any{"token": accessToken})
}

func (h *AuthHandler) ResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	var input ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	token, err := h.Users.InitiatePasswordReset(strings.TrimSpace(input.Email))
	if err != nil {
		responses.JSON(w, http.StatusOK, map[string]any{"message": "if the email exists, reset instructions were sent"})
		return
	}

	if h.Mailer != nil && h.Mailer.Host != "" && h.Mailer.Username != "" {
		body := "Use this token to reset your BlessedBites password: " + token
		_ = h.Mailer.Send(input.Email, "BlessedBites Password Reset", body)
	}

	responses.JSON(w, http.StatusOK, map[string]any{"message": "if the email exists, reset instructions were sent"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input ResetPasswordSubmission
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequest(w, "invalid JSON payload")
		return
	}

	if err := h.Users.FinalizePasswordReset(strings.TrimSpace(input.Token), input.NewPassword); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{"message": "password reset successful"})
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "bb_refresh",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "bb_refresh",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
