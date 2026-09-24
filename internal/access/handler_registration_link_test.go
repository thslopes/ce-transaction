package access

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrationLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		authMode           string
		mutateRepository   func(repo *MemoryUserRepository)
		expectedStatusCode int
		assertResponse     func(t *testing.T, body []byte)
	}{
		{
			name:               "returns created for admin",
			authMode:           "admin",
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(t *testing.T, body []byte) {
				var payload map[string]string
				decodeJSON(t, body, &payload)
				if payload["registration_url"] == "" {
					t.Fatal("expected registration_url in response")
				}
				if payload["expires_at"] != "2026-09-23T12:10:00Z" {
					t.Fatalf("expected fixed expiry, got %q", payload["expires_at"])
				}
			},
		},
		{
			name:               "requires authentication",
			authMode:           "none",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "rejects non admin user",
			mutateRepository: func(repo *MemoryUserRepository) {
				repo.CreateUserForTest("Requester", "requester@example.com", "+5511999991111", "requester-password", StatusActive, []string{"requester"})
			},
			authMode:           "requester",
			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server, tokens, _ := newTestServer(t, test.mutateRepository)
			request := httptest.NewRequest(http.MethodPost, "/users/registration-link", bytes.NewBufferString(`{}`))
			request.Header.Set("Content-Type", "application/json")
			switch test.authMode {
			case "admin":
				request.Header.Set("Authorization", "Bearer "+tokens["admin"])
			case "requester":
				request.Header.Set("Authorization", "Bearer "+tokens["requester"])
			}
			response, err := server.App().Test(request)
			if err != nil {
				t.Fatalf("fiber test request failed: %v", err)
			}
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Fatalf("expected status %d, got %d: %s", test.expectedStatusCode, response.StatusCode, string(body))
			}
			if test.assertResponse != nil {
				test.assertResponse(t, body)
			}
		})
	}
}