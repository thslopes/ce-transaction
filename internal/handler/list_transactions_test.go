package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	pkgTests "github.com/thslopes/ce-transactions/pkg/tests"
)

func TestClient_GetTransactionsHandler(t *testing.T) {
	tests := []struct {
		name              string
		descriptionFilter string
		limit             int
		offset            int
		expectedStatus    int
		expectedResponse  string
	}{
		{
			name:              "GetTransactionsHandler success",
			descriptionFilter: "desc",
			limit:             32,
			offset:            2,
			expectedStatus:    http.StatusOK,
			expectedResponse:  "tests/responses/router/listTransactions.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			app := fiber.New()
			c := MockClient()
			app.Get("/", c.GetTransactionsHandler)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?description=%s&limit=%d&offset=%d", tt.descriptionFilter, tt.limit, tt.offset), nil)

			// Act
			resp, err := app.Test(req)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			pkgTests.AssertBody(t, resp.Body, tt.expectedResponse)
		})
	}
}
