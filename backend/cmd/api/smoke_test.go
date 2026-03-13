package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// smokeSetup creates a test server connected to a real database.
// Skips if TEST_DSN is not set.
func smokeSetup(t *testing.T) (ts *httptest.Server, base string) {
	t.Helper()
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Skip("TEST_DSN not set; skipping smoke test (set TEST_DSN=<postgres-dsn> to run)")
	}

	db, err := openDB(dsn)
	if err != nil {
		t.Fatalf("smoke: open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	app := &application{
		config: config{
			addr:       ":0",
			dsn:        dsn,
			jwtSecret:  "smoke-test-jwt-secret-must-be-32bytes",
			corsOrigin: "*",
		},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		db:     db,
	}

	ts = httptest.NewServer(app.routes())
	t.Cleanup(ts.Close)
	return ts, ts.URL + "/api/v1"
}

func TestSmoke_HealthCheck(t *testing.T) {
	ts, base := smokeSetup(t)
	_ = ts

	resp, err := http.Get(base + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health: want 200, got %d", resp.StatusCode)
	}
}

func TestSmoke_AuthFlow(t *testing.T) {
	_, base := smokeSetup(t)

	// Use a unique email so this test is idempotent across re-runs on the same DB.
	email := fmt.Sprintf("smoketest%d@example.com", time.Now().UnixNano())
	password := "SmokePass123!"

	client := &http.Client{
		// Follow cookies automatically (jar-less; we handle cookies manually).
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}

	// 1. Signup ----------------------------------------------------------------
	t.Run("signup", func(t *testing.T) {
		body := map[string]string{
			"email":     email,
			"full_name": "Smoke Test User",
			"phone_no":  "+18001234567",
			"password":  password,
		}
		resp := doPost(t, client, base+"/auth/signup", nil, body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			dumpBody(t, resp.Body)
			t.Fatalf("signup: want 201, got %d", resp.StatusCode)
		}
	})

	// 2. Login → get access token + refresh cookie -----------------------------
	var accessToken string
	var refreshCookie *http.Cookie

	t.Run("login", func(t *testing.T) {
		body := map[string]string{"email": email, "password": password}
		resp := doPost(t, client, base+"/auth/login", nil, body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dumpBody(t, resp.Body)
			t.Fatalf("login: want 200, got %d", resp.StatusCode)
		}

		var result struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("login: decode response: %v", err)
		}
		if result.Token == "" {
			t.Fatal("login: empty token in response")
		}
		accessToken = result.Token

		for _, c := range resp.Cookies() {
			if c.Name == "bb_refresh" {
				refreshCookie = c
			}
		}
		if refreshCookie == nil {
			t.Fatal("login: bb_refresh cookie not set")
		}
	})

	// 3. Refresh → token rotated, new cookie issued ----------------------------
	var rotatedCookie *http.Cookie

	t.Run("refresh_rotates_token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, base+"/auth/refresh", nil)
		req.AddCookie(refreshCookie)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("refresh: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dumpBody(t, resp.Body)
			t.Fatalf("refresh: want 200, got %d", resp.StatusCode)
		}

		var result struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("refresh: decode: %v", err)
		}
		if result.Token == "" || result.Token == accessToken {
			t.Fatal("refresh: expected a new, distinct access token")
		}
		accessToken = result.Token

		for _, c := range resp.Cookies() {
			if c.Name == "bb_refresh" {
				rotatedCookie = c
			}
		}
		if rotatedCookie == nil {
			t.Fatal("refresh: new bb_refresh cookie not issued")
		}
		if rotatedCookie.Value == refreshCookie.Value {
			t.Fatal("refresh: cookie was not rotated (same value)")
		}
	})

	// 4. Replay old refresh cookie → must be rejected (revoked) ----------------
	t.Run("replay_old_refresh_rejected", func(t *testing.T) {
		if refreshCookie == nil {
			t.Skip("no refresh cookie from previous step")
		}
		req, _ := http.NewRequest(http.MethodPost, base+"/auth/refresh", nil)
		req.AddCookie(refreshCookie) // old, already consumed token
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("replay: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("replay: want 401 for revoked token, got %d", resp.StatusCode)
		}
	})

	// 5. Admin endpoint returns 403 for a regular customer ---------------------
	t.Run("admin_forbidden_for_customer", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, base+"/admin/orders", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("admin orders: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("admin orders: want 403 for customer, got %d", resp.StatusCode)
		}
	})

	// 6. Logout revokes the rotated cookie ------------------------------------
	t.Run("logout", func(t *testing.T) {
		if rotatedCookie == nil {
			t.Skip("no rotated cookie from previous step")
		}
		req, _ := http.NewRequest(http.MethodPost, base+"/auth/logout", nil)
		req.AddCookie(rotatedCookie)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("logout: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("logout: want 200, got %d", resp.StatusCode)
		}

		// Cookie after rotated token has been revoked by logout — refresh must fail.
		req2, _ := http.NewRequest(http.MethodPost, base+"/auth/refresh", nil)
		req2.AddCookie(rotatedCookie)
		resp2, err := client.Do(req2)
		if err != nil {
			t.Fatalf("post-logout refresh: %v", err)
		}
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusUnauthorized {
			t.Fatalf("post-logout refresh: want 401, got %d", resp2.StatusCode)
		}
	})
}

// doPost is a helper that POSTs JSON and returns the response.
func doPost(t *testing.T, client *http.Client, url string, headers map[string]string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("doPost marshal: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("doPost new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("doPost %s: %v", url, err)
	}
	return resp
}

func dumpBody(t *testing.T, r io.Reader) {
	t.Helper()
	b, _ := io.ReadAll(r)
	t.Logf("response body: %s", b)
}
