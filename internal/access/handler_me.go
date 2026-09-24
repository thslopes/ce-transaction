package access

import "github.com/gofiber/fiber/v2"

func (s *Server) handleMe(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*User)
	return writeJSON(c, fiber.StatusOK, userResponse(user))
}
