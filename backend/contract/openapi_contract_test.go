package contract

import (
	"os"
	"strings"
	"testing"
)

func TestOpenAPIContainsRequiredPaths(t *testing.T) {
	b, err := os.ReadFile("../openapi.yaml")
	if err != nil {
		t.Fatalf("failed to read openapi spec: %v", err)
	}
	content := string(b)

	required := []string{
		"/menu:",
		"/menu/{id}:",
		"/categories:",
		"/search:",
		"/orders:",
		"/orders/stream:",
		"/ratings:",
		"/ratings/{menu_item_id}:",
		"/auth/login:",
		"/auth/signup:",
		"/auth/logout:",
		"/auth/refresh:",
		"/auth/reset-password-request:",
		"/auth/reset-password:",
		"/admin/menu:",
		"/admin/menu/{id}:",
		"/admin/category:",
		"/admin/category/{id}:",
		"/admin/orders:",
	}

	for _, path := range required {
		if !strings.Contains(content, path) {
			t.Fatalf("openapi spec missing required path: %s", path)
		}
	}
}
