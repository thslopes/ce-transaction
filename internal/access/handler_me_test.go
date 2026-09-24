package access

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		authMode           string
		expectedStatusCode int
		assertResponse     func(t *testing.T, body []byte)
	}{
		{
			name:               "returns authenticated user",
			authMode:           "admin",
			expectedStatusCode: http.StatusOK,
			assertResponse: func(t *testing.T, body []byte) {
				var payload map[string]any
				decodeJSON(t, body, &payload)
				if payload["email"] != "admin@example.com" {
					t.Fatalf("expected authenticated admin in response, got %#v", payload)
				}
			},
		},
		{
			name:               "requires authentication",
			authMode:           "none",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server, tokens, _ := newTestServer(t, nil)
			request := httptest.NewRequest(http.MethodGet, "/me", nil)
			if test.authMode == "admin" {
				request.Header.Set("Authorization", "Bearer "+tokens["admin"])
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