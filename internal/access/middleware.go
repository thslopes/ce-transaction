package access

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) authenticate(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return writeError(c, fiber.StatusUnauthorized, "unauthorized", "Authentication required")
	}
	token := strings.TrimPrefix(authorization, "Bearer ")
	user, err := s.service.AuthenticateToken(token)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, "unauthorized", "Authentication required")
	}
	c.Locals("user", user)
	return c.Next()
}