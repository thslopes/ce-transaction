package access

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleRegistrationLink(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*User)
	url, expiresAt, err := s.service.GenerateRegistrationLink(user)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			return writeError(c, fiber.StatusUnauthorized, "unauthorized", "Authentication required")
		case errors.Is(err, ErrForbidden):
			return writeError(c, fiber.StatusForbidden, "forbidden", "User does not have permission to perform this action")
		default:
			return writeError(c, fiber.StatusInternalServerError, "internal_error", "Internal server error")
		}
	}
	return writeJSON(c, fiber.StatusCreated, map[string]any{
		"registration_url": url,
		"expires_at":       expiresAt.Format(time.RFC3339),
	})
}
