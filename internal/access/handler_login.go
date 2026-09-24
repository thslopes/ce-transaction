package access

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleLogin(c *fiber.Ctx) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid_request", "Request body must be valid JSON")
	}
	token, user, err := s.service.Authenticate(request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return writeError(c, fiber.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
		case errors.Is(err, ErrInactiveUser), errors.Is(err, ErrMissingProfiles):
			return writeError(c, fiber.StatusForbidden, "forbidden", "User is not allowed to access the system")
		default:
			return writeError(c, fiber.StatusInternalServerError, "internal_error", "Internal server error")
		}
	}
	c.Set("Authorization", "Bearer "+token)
	return writeJSON(c, fiber.StatusOK, map[string]any{"user": userResponse(user)})
}