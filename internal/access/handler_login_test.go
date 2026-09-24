package access

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		email              string
		password           string
		mutateRepository   func(repo *MemoryUserRepository)
		expectedStatusCode int
		expectAuthHeader   bool
	}{
		{
			name:               "returns authorization header for bootstrap admin",
			email:              "admin@example.com",
			password:           "admin-password",
			expectedStatusCode: http.StatusOK,
			expectAuthHeader:   true,
		},
		{
			name:               "rejects invalid credentials",
			email:              "admin@example.com",
			password:           "wrong-password",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:     "rejects inactive user",
			email:    "inactive@example.com",
			password: "inactive-password",
			mutateRepository: func(repo *MemoryUserRepository) {
				repo.CreateUserForTest("Inactive", "inactive@example.com", "+5511999992222", "inactive-password", StatusInactive, []string{"requester"})
			},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:     "rejects user without profiles",
			email:    "noprofile@example.com",
			password: "noprofile-password",
			mutateRepository: func(repo *MemoryUserRepository) {
				repo.CreateUserForTest("No Profile", "noprofile@example.com", "+5511999993333", "noprofile-password", StatusActive, nil)
			},
			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server, _, _ := newTestServer(t, test.mutateRepository)
			request := httptest.NewRequest(http.MethodPost, "/login", jsonBody(t, map[string]string{
				"email":    test.email,
				"password": test.password,
			}))
			request.Header.Set("Content-Type", "application/json")
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
			gotAuthHeader := response.Header.Get("Authorization") != ""
			if gotAuthHeader != test.expectAuthHeader {
				t.Fatalf("expected auth header %v, got %v", test.expectAuthHeader, gotAuthHeader)
			}
		})
	}
}
