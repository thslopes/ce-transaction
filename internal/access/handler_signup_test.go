package access

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		mutateRepository   func(repo *MemoryUserRepository)
		tokenMode          string
		requestBody        map[string]string
		expectedStatusCode int
		assertResponse     func(t *testing.T, body []byte)
	}{
		{
			name:      "creates inactive user with valid signup token",
			tokenMode: "valid-signup",
			requestBody: map[string]string{
				"name":                  "Candidate",
				"email":                 "candidate@example.com",
				"phone":                 "+5511999994444",
				"password":              "candidate-password",
				"password_confirmation": "candidate-password",
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(t *testing.T, body []byte) {
				var payload map[string]any
				decodeJSON(t, body, &payload)
				if payload["status"] != string(StatusInactive) {
					t.Fatalf("expected inactive status, got %#v", payload)
				}
				profiles, _ := payload["profiles"].([]any)
				if len(profiles) != 0 {
					t.Fatalf("expected no profiles, got %#v", payload)
				}
			},
		},
		{
			name:      "rejects invalid signup token",
			tokenMode: "invalid-signup",
			requestBody: map[string]string{
				"name":                  "Candidate",
				"email":                 "candidate@example.com",
				"phone":                 "+5511999994444",
				"password":              "candidate-password",
				"password_confirmation": "candidate-password",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:      "rejects expired signup token",
			tokenMode: "expired-signup",
			requestBody: map[string]string{
				"name":                  "Candidate",
				"email":                 "candidate@example.com",
				"phone":                 "+5511999994444",
				"password":              "candidate-password",
				"password_confirmation": "candidate-password",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:      "rejects duplicate email",
			tokenMode: "valid-signup",
			mutateRepository: func(repo *MemoryUserRepository) {
				repo.CreateUserForTest("Existing", "existing@example.com", "+5511999995555", "existing-password", StatusInactive, nil)
			},
			requestBody: map[string]string{
				"name":                  "Existing",
				"email":                 "existing@example.com",
				"phone":                 "+5511999995555",
				"password":              "candidate-password",
				"password_confirmation": "candidate-password",
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:      "rejects mismatched password confirmation",
			tokenMode: "valid-signup",
			requestBody: map[string]string{
				"name":                  "Candidate",
				"email":                 "candidate@example.com",
				"phone":                 "+5511999994444",
				"password":              "candidate-password",
				"password_confirmation": "other-password",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server, tokens, signer := newTestServer(t, test.mutateRepository)
			body := cloneStringMap(test.requestBody)
			switch test.tokenMode {
			case "valid-signup":
				body["token"] = tokens["signup"]
			case "invalid-signup":
				body["token"] = "invalid-token"
			case "expired-signup":
				token, err := signer.Sign(tokenClaims{Subject: "usr_001", Expiry: time.Date(2026, 9, 23, 11, 59, 0, 0, time.UTC), Purpose: "signup"})
				if err != nil {
					t.Fatalf("failed to create expired token: %v", err)
				}
				body["token"] = token
			}

			request := httptest.NewRequest(http.MethodPost, "/signup", jsonBody(t, body))
			request.Header.Set("Content-Type", "application/json")
			response, err := server.App().Test(request)
			if err != nil {
				t.Fatalf("fiber test request failed: %v", err)
			}
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Fatalf("expected status %d, got %d: %s", test.expectedStatusCode, response.StatusCode, string(responseBody))
			}
			if test.assertResponse != nil {
				test.assertResponse(t, responseBody)
			}
		})
	}
}