package access

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestServer(t *testing.T, mutateRepository func(repo *MemoryUserRepository)) (*Server, map[string]string, *TokenSigner) {
	t.Helper()
	clock := func() time.Time { return time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC) }
	repo := NewMemoryUserRepository(clock, BootstrapAdmin{
		Name:     "Admin",
		Email:    "admin@example.com",
		Phone:    "+5511999990000",
		Password: "admin-password",
	})
	if mutateRepository != nil {
		mutateRepository(repo)
	}
	signer := NewTokenSigner([]byte("secret-key"), clock)
	service := NewService(repo, signer, "https://app.example.com", clock)
	server := NewServer(service)

	adminRequest := httptest.NewRequest(http.MethodPost, "/login", jsonBody(t, map[string]string{
		"email":    "admin@example.com",
		"password": "admin-password",
	}))
	adminRequest.Header.Set("Content-Type", "application/json")
	adminResponse, err := server.App().Test(adminRequest)
	if err != nil {
		t.Fatalf("fiber login request failed: %v", err)
	}
	adminBody, err := io.ReadAll(adminResponse.Body)
	if err != nil {
		t.Fatalf("failed to read login response body: %v", err)
	}
	if adminResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected bootstrap admin login to succeed, got %d: %s", adminResponse.StatusCode, string(adminBody))
	}
	tokens := map[string]string{
		"admin": strings.TrimPrefix(adminResponse.Header.Get("Authorization"), "Bearer "),
	}

	if requester, _ := repo.FindByEmail("requester@example.com"); requester != nil {
		requesterRequest := httptest.NewRequest(http.MethodPost, "/login", jsonBody(t, map[string]string{
			"email":    "requester@example.com",
			"password": "requester-password",
		}))
		requesterRequest.Header.Set("Content-Type", "application/json")
		requesterResponse, err := server.App().Test(requesterRequest)
		if err != nil {
			t.Fatalf("fiber requester login request failed: %v", err)
		}
		if requesterResponse.StatusCode == http.StatusOK {
			tokens["requester"] = strings.TrimPrefix(requesterResponse.Header.Get("Authorization"), "Bearer ")
		}
	}

	signupToken, err := signer.Sign(tokenClaims{Subject: "usr_001", Expiry: clock().Add(10 * time.Minute), Purpose: "signup"})
	if err != nil {
		t.Fatalf("failed to create signup token: %v", err)
	}
	tokens["signup"] = signupToken

	return server, tokens, signer
}

func cloneStringMap(source map[string]string) map[string]string {
	copy := make(map[string]string, len(source))
	for key, value := range source {
		copy[key] = value
	}
	return copy
}

func jsonBody(t *testing.T, value any) *bytes.Buffer {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	return bytes.NewBuffer(data)
}

func decodeJSON(t *testing.T, data []byte, target any) {
	t.Helper()
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
}
