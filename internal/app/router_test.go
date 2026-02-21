package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_SetupRoutes(t *testing.T) {
	tests := []struct {
		name             string
		route            string
		method           string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "GET / - should return 200 and expected response",
			route:            "/",
			method:           "GET",
			expectedStatus:   200,
			expectedResponse: "Mock Get Transactions Handler - Method: GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			app := MockApp()
			req := httptest.NewRequest(http.MethodGet, "/transactions", nil)

			// Act
			app.SetupRoutes()

			// Assert
			resp, err := app.FiberApp.Test(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Optionally, you can read the response body and compare it with the expected response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			if string(body) != tt.expectedResponse {
				t.Errorf("Expected response %q, got %q", tt.expectedResponse, string(body))
			}
		})
	}
}
